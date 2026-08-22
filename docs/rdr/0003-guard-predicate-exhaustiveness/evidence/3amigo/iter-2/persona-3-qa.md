Model: claude-opus-5[1m]

# 3 Amigos iteration 2 — Persona 3 (QA / Tester)

Delta scope: A5, A7, A8, and the passages that state or depend on the narrowed
lint promise (Technical Design narrowing paragraph, Normative Contracts
exhaustiveness-narrowing clause, Decision Rationale, Contradiction Check,
Prerequisites, Minimum Viable Validation, Validation / Testing Strategy).

Read alongside: `evidence/spikes/guard-fixture.toml`, `evidence/spikes/check.sh`,
`evidence/spikes/output.txt`, `evidence/research/iter-2-projection-derivation.md`,
`docs/rdr/0006-graph-lint-authority-and-guarantees.md`,
`docs/rdr/0007-guard-predicate-totality.md`, `docs/jdr/0001-resolve-kernel-seam.md`.

---

## HIGH

### P3-1 — "refuse or downgrade" is a disjunction with no discriminator, and "downgrade" has no observable form

**Anchored passage**: Normative Contracts —

> "If a guard dimension lacks a finite declared domain, lint MUST refuse or
> downgrade an exhaustiveness claim for that dimension rather than treating the
> covered examples as complete."

and the sibling clause —

> "If the finite product is too large for that proof, lint MUST refuse or
> downgrade the exhaustiveness claim rather than silently capping enumeration."

Same disjunction repeated in A2's Evidence, the Technical Design paragraph
("lint must also refuse or downgrade the exhaustiveness claim instead of
silently capping the proof"), the `disposition` table rows "Guard dimension
lacks a finite declared domain" and "Finite product too large to prove
deterministically" (both "**Refuse or downgrade**"), Prerequisites-adjacent
Phase 2, Cross-Cutting Concerns, and Testing Strategy Scenario 3.

**Untestable criterion**: The RDR never says *which* of the two outcomes applies
under which condition, and never defines what "downgrade" produces. "Refuse" is
observable — RDR 0006 gives it a code, `graph-unprovable-coverage`
(`docs/rdr/0006-graph-lint-authority-and-guarantees.md:302`), blocking, exit
`GroupUserEnv`. "Downgrade" has no code, no severity, no exit behavior, and no
payload anywhere in either document. RDR 0006's advisory channel is described
only as "Non-blocking advisories, if implemented, are warnings and do not change
exit behavior" (line 316) — conditional on an implementation choice this RDR
does not make. JDR 0001 §JD-4 explicitly leaves "whether lint gains a warning
category" open (`docs/jdr/0001-resolve-kernel-seam.md:238-241`).

**Test it prevents**: Testing Strategy Scenario 3 as written
("**Expected**: Runtime evaluation remains available, but lint refuses or
downgrades the exhaustiveness claim for that dimension/product") cannot be
turned into an assertion. A test must assert one exit code and one finding code.
With a free disjunction and one undefined arm, an implementation that emits
nothing observable but internally "downgraded" a claim passes the sentence, and
an implementation that blocks also passes. The MVV's stated refusal/downgrade
cases ("refusal/downgrade cases for unbounded or too-large finite domains",
Scope Verification) inherit the same defect. **Fix shape**: pick one outcome per
input class, or name the discriminator and give `downgrade` a stable finding
code and exit behavior, so Scenario 3 asserts `graph-unprovable-coverage` with a
named exit.

### P3-2 — MVV Scenario 6's central assertion ("lint withholds") has no observable form and no code

**Anchored passage**: Minimum Viable Validation —

> "The fixture must also include one row group that is domain-exhaustive yet
> contains a possibly-absent guard key, and assert that lint withholds the
> exhaustiveness claim there rather than contradicting the runtime refusal."

and Testing Strategy Scenario 6 —

> "**Expected**: The value atom is unevaluable rather than false, resolution
> refuses `guard_unevaluable` under RDR 0007's veto, and lint does not certify
> that row group exhaustive"

and the `disposition` rows "Value atom over an absent key → Exhaustiveness claim
**withheld** for that group … Loud" and "Row group is domain-exhaustive but a
participating row can refuse → Claim withheld — the narrowing … RDR 0007
payload; no green certification | Loud".

**Untestable criterion**: "withheld" and "does not certify" are stated as
negatives. Both `disposition` rows mark the outcome **Loud**, which by the
table's own column semantics means something is emitted — but the "Diagnostic
minted" cell for the second row says "RDR 0007 payload; no green certification",
i.e. the only named artifact belongs to the *runtime* refusal, not to lint. Lint
in this project is a command that emits findings and an exit code (RDR 0006
Normative Contracts, `graph-lint-failed` aggregate). RDR 0006's blocking code
table has no code for a withheld exhaustiveness claim; the nearest,
`graph-unprovable-coverage`, is defined for "Required finite-domain proof
unavailable" — and this case is precisely one where the finite-domain proof *is*
available and succeeds. So the RDR asks for an assertion on a lint run without
saying whether that run passes with no finding, fails with a finding, or emits
an advisory.

**Test it prevents**: The single test the entire §JD-4 narrowing exists to pin —
"run lint on a domain-exhaustive row group whose participating row carries a
value atom over a possibly-absent key; assert X". X is undefined. A test asserting
`exit != 0` and a test asserting `exit == 0 && no coverage-gap finding` are both
consistent with "withholds the claim", and they are opposite tests. This is the
MVV's own added clause, so the RDR cannot lock a promise it cannot pin.

### P3-3 — Scenario 6 requires a "possibly-absent guard key", which the fixture format cannot express

**Anchored passage**: Testing Strategy Scenario 6 — "one participating row
carries a value atom over a key that can be absent" — and MVV — "contains a
possibly-absent guard key". Cross-anchored to A7's **Plan**:

> "Blocked on that request: RDR 0002's tag declaration carries name, provenance,
> value kind, and optional accessor reference, with no optionality field, so no
> producer can today declare which keys may be absent."

and Capability Dependencies row "Per-tag optionality declaration (which keys may
be absent) | RDR 0002 | **Requested**".

**Untestable criterion**: The MVV mandates a fixture property — a key declared
possibly-absent — that the RDR itself states no declaration surface exists to
express. `evidence/spikes/guard-fixture.toml` confirms this: all seven tag
declarations carry `provenance` / `kind` / `domain` (or `min`/`max`) and none
carries an optionality marker. A test author writing Scenario 6 must invent
either the declaration syntax (which is RDR 0002's to own, per the same
Capability Dependencies row) or an out-of-band mechanism (omit the key from the
runtime input tuple and rely on it), and lint proves over *declarations*, not
over input tuples — so the second route gives lint nothing to see.

**Test it prevents**: Scenario 6's lint half. Its runtime half is writable
(supply an input tuple missing the key, assert `guard_unevaluable`); its lint
half cannot be written until the optionality declaration lands, because lint has
no way to know the key "can be absent". The Prerequisites checkbox is honest
that A7 blocks lock, but the MVV states Scenario 6 as an implementation
obligation without flagging the same block — Scope Verification says "The
Minimum Viable Validation is in scope for implementation, not deferred", which
is not true of the Scenario 6 fixture as specified.

---

## MEDIUM

### P3-4 — A7's `exists` projection is `Pending`, but the operator is already in the closed vocabulary the MVV must exercise

**Anchored passage**: A7 Status `Pending`; `trace` step 5 —

> "5. Project the `exists` atom | A7 presence dimension — **Pending** |
> `foundational-to-cove` carries `cluster_eligible exists = true`;
> `cluster_eligible`'s declared domain `{true,false}` is complete, so no element
> selects presence | **GAP — booked as A7** … Without the presence dimension
> this step has no defined result."

against MVV — "Encode one RDR flow slice and one kata flow slice as normalized
candidate rows, including equality, enum membership, set containment, bounded
integer comparison, and mixed `all`/`unless` guards" — and Testing Strategy
Scenario 2, which requires "one complete partition" over "finite enum, boolean,
set-universe, and bounded-int tag domains".

**Untestable criterion**: The lint-side semantics of `exists` are undefined for
any row group containing one. Note the MVV's operator list omits `exists`, while
Scenario 1 includes it ("equality, membership, set containment, bounded integer
comparison, existence, and mixed `all`/`unless` guards") — so the runtime
scenario exercises the operator and the lint scenario silently does not. The
draft's own representative fixture row `foundational-to-cove` sits in the same
selection group as `profile-to-grounding` (both `match status eq Draft`, per
`trace` step 2), so the RDR's *only* multi-row representative group is exactly
the one with no defined coverage result.

**Test it prevents**: Any coverage/overlap assertion over a row group containing
an `exists` atom — including the coverage assertion on the fixture's own primary
group. `trace` step 6 concedes this: "Step 5 unresolved for this group; groups
with no `exists` atom compute normally". The RDR is transparent that this is a
`Pending` gap, so this is not a hidden defect; it is recorded here because it
directly narrows what MVV Scenario 2 can assert, and no passage states which
fixture group Scenario 2's "one complete partition" is to be built on given the
constraint.

### P3-5 — A8 states a cross-document MUST with no test surface

**Anchored passage**: Normative Contracts, exhaustiveness-narrowing clause,
final sentence —

> "RDR 0006, which holds the JDR 0001 §JD-4 tolerance, MUST carry this same
> narrowing for the findings it emits; the two documents MUST NOT state it
> differently."

and A8's verification route: "Verification is a route-back to RDR 0006
confirming it adopts the same wording, or a §JD-4 disposition assigning the
recording to one document."

**Untestable criterion**: A normative MUST whose subject is another *document's*
prose. Every other clause in this section constrains code (a guard predicate, an
evaluator, lint, a diagnostic); this one constrains RDR 0006's wording. There is
no pass/fail criterion an implementation test can evaluate — "MUST NOT state it
differently" has no mechanical form, and "same wording" is not defined (identical
text? equivalent semantics? RDR 0006's clause is scoped to non-finite dimensions
and this RDR's to a finite product with a refusing row, so identical text is not
even coherent across the two scopes).

**Test it prevents**: The conformance test for the clause. Note the asymmetry:
if the narrowing's *behavioral* content were stated once as a lint outcome (see
P3-2), a single test on the lint command would pin it regardless of which
document records it — which is exactly what §JD-4 is for. As written, the clause
routes a testable behavior into an untestable document-agreement obligation.
`trace` step 10 books this honestly as GAP/A8, but the normative clause is
already in the contract block and would ship as a MUST an implementer cannot
discharge.

### P3-6 — "can refuse" in the narrowing clause is a modal with no decision procedure

**Anchored passage**: Normative Contracts —

> "lint MUST NOT certify a row group exhaustive when a participating row **can
> refuse** `guard_unevaluable` under RDR 0007's aggregation veto."

Restated in Technical Design ("Lint therefore cannot certify a row group
exhaustive when any participating row could refuse `guard_unevaluable` at
runtime") and in `trace` step 8 ("withhold if a participating row can refuse").

**Untestable criterion**: "can refuse" is a possibility claim over *all* runtime
inputs. The RDR gives no procedure for deciding it statically. Two readings with
opposite test outcomes: (a) syntactic — any row carrying a value atom over a key
declared possibly-absent; (b) semantic — any row for which some assignment in the
scoped product reaches the veto. Reading (a) is decidable but over-withholds;
reading (b) requires the presence dimension A7 has not yet defined. RDR 0007's
veto is resolution-scoped ("if any surviving candidate row's guard is
GuardUnevaluable", `docs/rdr/0007-guard-predicate-totality.md:1649-1651`),
depending on which rows *survive* a given input — a runtime property this RDR's
static product proof has no stated way to over-approximate.

**Test it prevents**: A negative-control test — a row group that is
domain-exhaustive and that lint *should* certify green. Without a decision
procedure for "can refuse", there is no way to construct a fixture and assert
"lint certifies this one", which means the narrowing clause cannot be shown to
be *tight* rather than a blanket refusal. The MVV asks only for the withholding
case; the complementary green case is untestable and unasked-for.

---

## LOW

### P3-7 — A5's Evidence mixes a specified shape with a shipped one, and names no verification command

**Anchored passage**: A5 Evidence —

> "The normalized atom shape is fixed at the kernel seam by JDR 0001 §D1; RDR
> 0007 is the landing document and states it normatively — a row carries parsed
> atoms (key, operator token, literal, block), not an opaque predicate string,
> so there is no reconstruction step that could lose identity. RDR 0007 is
> `Final`, not yet implemented: the shipped kernel still carries
> `Row.Guard string`, and 0007's A3 scopes the reshape."

**Untestable criterion**: The assumption is `Verified` by Method `Peer RDR`, and
the passage is careful and correct about the split (`evidence/spikes/iter-2/atom-shape.md`
confirms `Row.Guard string` ships at `resolve.go:185` while `Row.RuleID` /
`Row.SourceLocator` ship at `resolve.go:170-172`). What it does not give is a
regression criterion: no test name, no assertion, nothing that would fail if the
reshape lands *without* per-atom source identity. The only source-identity
assertion the RDR asks for is at row granularity — Normative Contracts,
"Overlap and coverage diagnostics MUST name the source rule id or context id
that contributed each predicate involved in the finding" — while the Load-Bearing
"Identity" entry defines predicate identity as "its source rule/context id plus
its position within `all` or `unless`", a finer granularity nothing asserts.

**Test it prevents**: A diagnostic-precision test — given a row with three `all`
atoms where one causes the overlap, assert the finding names *that atom*, not
just the row. `trace` step 9 checks only "`RuleID` + `SourceLocator` ship on
`internal/resolve/resolve.go::Row`", which is row-level. A5's "If wrong"
("Lint may detect an error but fail to point reviewers at the guard to fix") is
about the atom, but no criterion in the document tests at that granularity.
Severity low because the row-level assertion is testable today and the atom-level
gap is a precision limit, not a soundness one.

### P3-8 — Prerequisites checkbox conflates two blocks with different closure routes and no verifiable done-condition

**Anchored passage**: Prerequisites —

> "- [ ] All Critical Assumptions verified — A1-A6 hold (A5 re-verified against
> the kernel-fixed atom shape, JDR 0001 §D1); **A7 is Pending**, blocked on RDR
> 0002's optionality declaration; **A8 is Pending**, blocked on a route-back to
> `Final` RDR 0006 to agree one lint promise"

**Untestable criterion**: One unchecked box gates on two conditions whose
done-criteria differ in kind: A7 closes when a *peer document adds a field* and
this RDR *states a normative clause*; A8 closes when a `Final` peer is amended or
a JDR entry is dispositioned. Neither has a stated verifier. A7's Plan says
"state the projection rule as a normative clause before lock" — but the clause
is not drafted anywhere in this document, so there is no text to check against.
A8's Plan says "raise the divergence at cluster reconcile before either document
locks", naming a process step rather than an artifact.

**Test it prevents**: A pre-lock mechanical check — no script or reviewer can
evaluate "is this box now checkable?" against a stated artifact. Severity low
because both blocks are named explicitly, honestly, and repeated in the
Contradiction Check and Assumption Verification; the defect is the absence of a
done-condition, not concealment.

---

## Not findings (checked, testable as written)

- Testing Strategy Scenario 5 (reorder rows and atoms) — has a clean pass/fail:
  "Successful matching and lint findings are unchanged". Writable today against
  the fixture.
- Testing Strategy Scenario 4 (malformed atoms) — four named input classes, each
  with a stated rejection point ("before resolution") and a stated kind. Writable.
- Phase 3 / Testing Strategy `contains` obligation — "the implementation MVV must
  add at least one `contains` predicate over a declared set-valued tag before the
  full operator vocabulary is accepted" is a checkable gate. Note it depends on
  RDR 0007 Phase 4's set-element-encoding request, which Capability Dependencies
  records; the RDR does not overclaim it.
- Decision Rationale's §JD-4 paragraph and the Contradiction Check's third
  paragraph — narrative about the divergence, not criteria; consistent with A8.
