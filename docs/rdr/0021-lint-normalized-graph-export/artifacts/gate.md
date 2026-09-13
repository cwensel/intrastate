# Finalization Gate — RDR 0021, lint's normalized-graph export

- **Record**: `cli/0021` (`0021-lint-normalized-graph-export`)
- **Date**: 2026-09-13
- **Verdict**: READY — Gate PASS, locked to Final (re-lock after the Stage 8
  route-back over REQ-20's `domain` presence rule).

Mechanical pre-sweep: `evidence/tooling-pass/iter-4/tooling-pass.md` — PASS
(`rdr lint --locking` exit 0, `blocking=0 resolution=0 placeholder=0`, one
advisory answered under item 2). The Stage-8 re-entry note's defect verifies
closed in live text (that report carries the item-by-item table); the note is
deleted at this lock. Item 4, Cross-Cutting Concerns, is authored in the record
at `0021:G-cross-cutting` and is deliberately not copied here.

## 1. Contradiction Check

No contradictions between Research Findings and the Proposed Solution, and
none between planned features and stated principles.

- Findings read `0005:C1` as carving `lint`/`dump`/`parse`-class command
  groups out of the flow contract, "owned by the RDR that names them". C1
  claims a root `graph` verb under exactly that carve-out, and A4 (Method:
  Peer RDR) independently verifies no 0005 envelope amendment is needed,
  re-verified against the post-0029 envelope.
- Findings read 0002 as fixing the dump's field list and row order while
  declining to define a dump grammar. C2 owns the JSON schema on that seeded
  ground and reuses 0002's field vocabulary and row order rather than minting
  a second row contract, amending no 0002 text, and closes
  `0002:§round-trip-inverse-invariants`'s lossy set-literal rendering for
  this document only, claiming no export→load inverse (RT3).
- Findings record a negative corpus result: no opened peer wraps DOT in a
  JSON envelope. C5 claims no peer precedent for that cell; it defines it
  normatively as `data.dot`, one required string member, unwrapped
  `jq -r .data.dot`. The negative result is carried as a constraint, not
  silently overridden.

Principles vs planned features: the neutrality principle ("an export is never
a lint pass") is not merely asserted — C4 makes it normative and MVV step 5
gives it a mechanism-independent oracle, byte-comparing lint's refusal with
and without the export present. C2's "no verdict or finding field, and no
per-node terminal marking" applies the same principle to the document
surface, leaving the dead-end quantifier to RDR 0015 (JDR 0001 §JD-23).

**This round's narrowing introduces no contradiction.** C2's `domain` rule now
reads over AUTHORED members and states the carve-out explicitly — "Presence is
NOT `guard.AssignmentCount` finiteness" — which is the opposite of what the
pre-route-back text said. Checked against every sibling that could disagree:
A3's Evidence now derives the same two-predicate split from `domainSize`'s
per-kind arms; `kind` continues to carry the type for the omitting kinds, so
no information is lost; the "no derived value vocabulary is minted here" clause
forecloses the `bool`-literals / `Elements` / powerset / `{min..max}` readings
the old rule would have demanded. C1, C3, C4 and C5 are byte-unchanged and
none of them predicates anything on `domain` presence. The tier declarations
from the prior round still agree with the additive-within-`/1` rule.

## 2. Assumption Verification

All eight Critical Assumptions are internally consistent and terminal.

- **Status**: 8/8 `Verified`; zero `Pending`, `Unverified`, or placeholder
  (`ca=all-terminal`). No settled-fact prose leans on an unverified property.
- **Method**: every label sanctioned — `Spike` (A1, A2, A5, A6),
  `Source Search` (A3, A7, A8), `Peer RDR` (A4). Zero off-vocabulary members.
  No `Docs Only` record exists, so the load-bearing Docs-Only bar is vacuous
  rather than waived.
- **Evidence**: Status, Method and Evidence agree on every row; each
  "If wrong" is non-empty.
- **A3, this route-back's `re-verify` target**: discharged in place. Its
  Evidence now reads `domainSize`'s arms to the point where the two predicates
  come apart — `bool` (`spread(2, …), true`, "declared nowhere"), a bounded
  `int` (`Min`/`Max` width, `cardinalityCeiling` the overflow fallback) and
  `set` (`spread(len(d.Elements), false)`) all report FINITE while populating
  no `decl.Domain`; only an `enum` with a non-empty domain authors members.
  That is exactly the rule C2 now states, and exactly what the shipped
  exporter does (`finite && len(decl.Domain) > 0`,
  `internal/cli/graph_document.go:185`). Record and code agree.
- **Self-reference**: none. The three `Source Search` rows resolve into the
  product tree, never to this record or its artifact directory. A3's
  reconcile-artifact citation sits beside resolving source anchors, which is
  where a reconcile artifact belongs.
- **Symbol resolution**: 31 `source-anchor` edges, every one `resolved: true`
  under a bound `RDR_SOURCE_REPO` — none false, none absent, so the lookups
  genuinely ran. A5's Evidence is written `clierr.go:174::WriteJSONLine`; the
  symbol resolves and the stale line component is a documented non-finding.
- **Evidence-field budget (the one advisory)**: A3's field is 43 lines against
  a soft cap of 30. The author's question, answered: the load-bearing anchors
  remain findable — four `path::Symbol` anchors (`reach.go::Reach`,
  `declaration.go::AssignmentCount`, `product.go::Groups`,
  `reach.go::heldValues`), each resolving — and the balance is not padding but
  the per-kind `domainSize` derivation that IS the re-verification this
  route-back demanded, the content a later grounding sweep reads. Kept in
  place, not truncated and not relocated; Profile is `large`, not
  `foundational`, so the check is advisory here.
- **Two residuals are the implementation's, not this record's.** The doc
  comment at `graph_document.go:51-54` still asserts the refuted finite-only
  claim, and `TestReq20And28` leaves `bool`, bounded `int` and member-less
  `enum` unasserted. Both were verified still open on the unmerged branch
  `worktree-rdr-0021`. Neither is a record defect: the RDR's rule is now
  correct and the shipped code already implements it. A3 names both as owed by
  the implementation pass, and S2's coverage requirement is what will close the
  test gap.

## 3. Scope Verification

The Minimum Viable Validation is in scope and executed during implementation,
not deferred. It is Phase-1 work, gated on one authored fixture pair rather
than on any later phase.

The specific proof: author a two-owned-state / one-terminal / one-escape-row
state-machine fixture plus a decision-table fixture, then assert four
invocations — (a) `intrastate graph --model <fixture>` run twice, stdouts
BYTE-identical, parsing as JSON and carrying every C2 field including the
`reach` abstraction marker; (b) `--emit dot | dot -Tsvg` renders and its
node/edge id set equals the JSON `reach` block; (c) `--as=json | jq .data`
equals the document value-for-value; (d) `intrastate lint` over a fixture with
blocking findings refuses byte-identically against a pre-change capture, while
`graph` over that same model succeeds with an asserted document. That fourth
invocation is the neutrality oracle, and C4 binds it mechanism-independently.

Blast radius: not applicable. This record is unclustered as it locks
(`clustered=false`, `cluster=[]`, no `impact_families`) — the `0021-0029`
pairing rests on reciprocal `cross-cutting-owner` edges, not a declared
`Cluster:` field — so it retires and renames no peer's literals and no
`impact.md` projection is owed.

Joint-decision fence: `op = none` (`fence-clear`) after one real blocker was
fired and cleared in this pass. `index --literal-intersect` reported
`0012 0021 UNCITED 5 shared: bool enum int scalar set` — a NEW pair created by
this round's own narrowing, which pulled the five kind tokens into C2's
presence rule. Grounded per §ground-before-ask (`--outcome ground` → `code`),
the check fires CLEAR: `0003:C2` (Implemented, terminal) OWNS the vocabulary —
"the value kinds are exactly five, spelled with these tokens wherever a kind is
named" — and `0012:C1` inherits it by explicit citation as "RDR 0003's
spelling". The two records then read it at disjoint seams: 0012 carries the
declared kind to the VALUE-COMPARISON seam (`grammar.go::Evaluator`'s
`Evaluate(atom, value)`, constructed at `flow_resolve.go::guardSeam`,
`product.go::valueSatisfies`, `reach.go::atomAdmitsValue`), deciding what
`eq`/`in` answer per kind; `0021:C2` decides when a DOCUMENT field is emitted.
Confirmed in source: the exporter constructs no evaluator, and 0012 anchors
neither `AssignmentCount` nor `decl.Domain` —
`AssignmentCount`'s non-test callers are cardinality consumers (`coverage.go`,
`product.go`, `lint.go`, `reach.go`), never the comparison seam. Recorded as a
citation on the Joint-check line naming 0012 and the four shared symbols, not
as a synced copy; both intersect arms now read `cited`, 0 uncited. JC1's own
fire (→ 0029, home `cli/0029 §Normative Contracts` C4) is `homed`, and C-7's
two-marker read-order question stays at that home. `rulings_open=0` — all nine
author rulings are marked absorbed.

## 5. Proportionality

Right-sized; nothing flagged to trim before locking.

**Contract count, not word count.** C1–C5 are five labelled clauses of ONE
independent load-bearing contract: the export of the normalized graph. C1
fixes its surface (verb, arm set, flag order), C2 the document it emits, C3
that document's byte stability, C4 the neutrality boundary against lint, C5
how the two `--as` modes carry it. None is separately adoptable — a consumer
cannot take the document without the surface that emits it — so this is one
seam stated five ways, not five seams locked together. No split, and the
author ruling Q1 records why collapsing to a single `**C1**` is blocked:
`0029` (Final) holds a resolved `cross-cutting-owner` edge into `0021:C2` and
cites `0021:C5` by label, so relabelling would break inbound citations from a
no-amend record. `contracts_transient=0`, so no lifespan disposition distorts
the count.

**Profile re-validated.** The Metadata field reads `large`, and that still
matches the contracts just counted: one contract, user-facing yes (a new root
verb and a documented output document), locking the `intrastate.graph/1`
field list and marker, which `0029` consumes by a resolved edge into
`0021:C2`. `--outcome floor` returns `none` (`floor-below-two`), so no
accretion floor raises it. The lens battery `large` demands did run:
grounding, 3amigo, critique (two models, differing stamps, diff written) and
repeatability-lite (`--outcome repeatability` → `none`, rule
`repeatability-lite-complete`). `lens_stale=none` — the qualifier is cleared
and no lens folder predates it.

**Growth from this round is negative.** The narrowing rewrote C2's `domain`
clause in place and extended A3's Evidence with the per-kind derivation that
justifies it; this lock then removes a 46-line re-entry note, so the record
comes out shorter than it went in. The density sits in C2, where each field
spelling is normative and therefore load-bearing at implementation, and the
one over-budget Evidence field is the verification content answering the very
defect that caused the route-back. Proportionate.
