Model: claude-opus-5[1m]

# 3amigo iter-2 — Persona 1 (Product Manager): does RDR 0007 deliver the user outcome?

Delta scope: A3, A16, A17, A18, A19, A21, A22, the domain rule and its kernel
enforcement site, and the refusal payload's route to the user. The user outcome
under test is the Problem Statement's: *"They want a refusal they can act on:
told plainly that the artifact state needed to decide was missing, not handed a
plan that quietly routed around it."*

Findings are severity-ranked. Severity is judged by how much of that user
outcome the passage leaves undelivered or unverifiable.

---

## PM-1 — HIGH — The user-visible half of the outcome is explicitly declared not to block lock, so the RDR can lock while "told plainly" is undelivered

**Anchored passage**: Implementation Plan → Prerequisites, third box:
"**JDR 0001 §D4 is reopened and §JD-8 amended to grant ONE `omitempty`
structured field on `CLIError` (A19).**" … "**Does not block lock**: the
direction is reversible either way (an `omitempty` field is append-only), and
A19 carries a `Detail`-flattening fallback if the reopening is declined. Blocks
RDR 0005's implementation of the JSON rendering only." Also A19 **Status**:
"Pending — blocked on a peer-document change, not on evidence. The kernel half
is Verified."

**What is wrong from the PM view**: The Problem Statement's promise is not
"the kernel holds a payload"; it is that a human is *told plainly what was
missing*. This RDR's own Decision Rationale makes that the deciding criterion —
"it is the only one whose refusal meets the Problem Statement's 'told plainly
that the artifact state needed to decide was missing' at key granularity."
Yet the only clause that carries the payload out of `internal/resolve` to a
person is deferred to a peer-document reopening, and the RDR declares the
deferral non-blocking. The stated fallback is not outcome-neutral: A19's own
Evidence says flattening "puts a second, undocumented encoding inside a string,
which defeats the sorted-by-row-then-key determinism whose only purpose is
machine consumption." So the two branches deliver materially different products
and the RDR locks without choosing. `Refusal` is a struct in a package with
**zero production importers** (Key Discoveries, "Zero production importers of
`internal/resolve`"), so at lock the payload reaches nobody at all.

**Decision it blocks**: Whether RDR 0007 can be declared Final. As written,
"Final" means the masking path is closed in the kernel but the user-facing
outcome the RDR is justified by has no committed delivery mechanism. Also
blocks RDR 0005's rendering decision, which is now a coin-flip the PM cannot
plan around.

**What would close it**: either make the §D4 reopening a hard lock gate, or
pick the `Detail`-flattening branch *in this RDR* and restate the Problem
Statement's promise honestly at the granularity that branch actually delivers.

---

## PM-2 — HIGH — The one shipped user-facing string for this refusal describes the wrong cause, and no Prerequisite fixes it

**Anchored passage**: A19 → **Evidence**: "`guard_unevaluable` already has a row
in RDR 0005's stable-code table (`flow-guard-unevaluable`, `GroupUserEnv`) —
though scoped there to 'supplied facts', which the cluster gate already flagged
against this RDR's provenance-blind view".

**What is wrong from the PM view**: RDR 0005 line 611 reads
`| guard cannot be evaluated from supplied facts | flow-guard-unevaluable |
GroupUserEnv |`. The *primary* case this RDR exists for is a tag absent from the
**artifact state** (the accessor's owned snapshot), not a fact the caller
supplied. Under `GroupUserEnv` the exit-code grouping — which JDR 0001 P4 says
"encodes remedy" — tells the operator the remedy is their own input. For the
flagship scenario the remedy is "make the artifact state readable, or fix the
accessor" (Failure Modes → **Visible break**; Risks → "operators treat an
unevaluable refusal as 'inspect the accessor,' never 'retry'"). The RDR
acknowledges the mismatch in a parenthetical inside an assumption's evidence
block and then does nothing with it: the Prerequisites list five boxes and none
of them is "RDR 0005 restates `flow-guard-unevaluable`'s description and
re-checks its exit group." Compare the treatment of `owned_state_unavailable`,
where the missing `Code` row *is* named as a §JD-8 gap.

**Decision it blocks**: Whether the refusal an operator actually sees is
actionable. A correct payload behind a message that misdirects the remedy is
the same product failure the RDR was written to prevent, one layer up. Blocks
RDR 0005's message/hint authoring and the exit-group choice (§JD-8's live
question about exit 3 / `GroupEnvUnavailable` is left un-joined to this case).

---

## PM-3 — HIGH — No scenario anywhere validates that a user is told plainly; every scenario stops at kernel struct fields

**Anchored passage**: Validation → Testing Strategy preamble: "Every scenario is
a kernel test in `internal/resolve` (package `resolve_test`)"; and Minimum
Viable Validation: "yield `Refusal.Kind == guard_unevaluable`, a nil `Plan`, the
absent key in the payload, and the row in `Refusal.Rows`". Also Testing Strategy
→ "**Not covered here, by ownership.** … The CLI rendering of the payload is RDR
0005's and is blocked on the A19 envelope-field grant".

**What is wrong from the PM view**: All 20 scenarios and all four MVV rows
assert on in-process Go values. Not one asserts that a person or a JSON consumer
can distinguish "I could not read what I needed" from "this row does not apply"
— which the Problem Statement names as the exact confusion to eliminate ("nothing
in the output distinguishes 'this row does not apply' from 'I could not read
what I needed to find out'"). The RDR ships a mechanism and defers its
observable effect, with the deferral itself blocked (PM-1). The `oracle`
mini-check reinforces this: its four rows discriminate on refusal KIND and
selection, never on what the output says.

**Decision it blocks**: Acceptance. There is no test whose failure would tell
anyone that the user outcome regressed, so "done" for this RDR cannot be
distinguished from "the kernel struct has a new field." Blocks defining a
cross-RDR acceptance scenario (0007 payload → 0005 envelope) and blocks anyone
signing off that kata `xg7p` closed "against behavior, not against a doc
comment" (Trade-offs → Consequences) at the level a user experiences behavior.

---

## PM-4 — MEDIUM — The aggregation veto's blast radius on real users is asserted as acceptable with no sizing and no escape hatch

**Anchored passage**: Normative Contracts → SURVIVOR MEMBERSHIP: "This veto is
the cost the RDR accepts: one unreadable row refuses a table whose other rows
decide cleanly. Narrowing it would decide that an undecided edge is a non-edge —
absence-as-false at resolution scope." Reinforced by Failure Modes →
"**Seam unevaluable on a present value** (A18) … one malformed observed value
whose key any guard mentions can DENY every resolution touching it,
non-escapably".

**What is wrong from the PM view**: This is the single largest behavioral
change a table author will feel, and it is stated as a principle with no
estimate of how often it fires and no recovery affordance. The correctness
argument is sound; the product argument is missing. Trade-offs →
Consequences bounds it only by "no evaluator exists yet," which bounds the
*migration*, not the steady state. The A18 case is worse than the absence case
because the trigger is caller-supplied and unconstrained (A13), so an operator
can deny an entire flow by typo — non-escapably, by design.

**Decision it blocks**: Whether the veto ships at resolution scope or needs a
per-row diagnostic mode / dry-run affordance so an author can see which rows
would veto before the veto is load-bearing. Also blocks sizing the Phase 2
"classify each authored guard in the reference fixtures as safe-or-migration"
work — no acceptance threshold is stated for what fraction of authored guards
may need migration before the design is reconsidered.

---

## PM-5 — MEDIUM — The `unless`-block footgun is a known user-outcome regression, sanctioned by an idiom, and handed off with no gate

**Anchored passage**: Failure Modes → "**`unless` over a rarely-set tag**
(premortem P-5): an `unless` block whose only atom is a value atom over an
absent key is `¬U = U` and vetoes — correct by the rule and a footgun in
practice ('skip when `legal_hold` is true' refuses on every artifact that never
had `legal_hold`). The sanctioned idiom conjoins an existence atom … Phase 4
hands this idiom to RDR 0003's authoring guidance."

**What is wrong from the PM view**: The most natural way an author expresses
the most common real policy ("skip when the hold flag is set") produces a
non-escapable refusal on every artifact that never had the flag — which is the
overwhelmingly common case. The RDR labels this "correct by the rule" and routes
the fix to another Draft RDR's authoring prose, in Phase 4, with no Prerequisite
and no lint/validation requirement that the sanctioned idiom be enforced or even
warned about. RDR 0003 is Draft; nothing compels it to land the guidance, and
the RDR's own Alternative-1 rejection reasoning ("a conformance harness that
nothing compels RDR 0003's build to run") applies verbatim here.

**Decision it blocks**: Whether the `unless`-with-value-atom shape needs a load-
time rejection or a lint warning rather than authoring prose. This is the same
"held by discipline vs. held by structure" fork the RDR resolved decisively for
the domain rule and left unresolved for the shape most likely to hit users.

---

## PM-6 — MEDIUM — A22 and A16's producer halves are unstated inferences; the user outcome depends on clauses no document carries

**Anchored passage**: A22 **Status**: "Pending — kernel half Verified; the
producer half is an UNSTATED INFERENCE against RDR 0002's current text." Its
Evidence: "The only supporting sentence is JDR 0001 §D4's landing note … an
obligation FILED on 0002, which 0002 has not received as a clause." Plus A16's
Evidence: "The PRODUCER half is not yet a duty on 0002: its Refinement Context
Direction list names the `RequiresOwned` producer and the fixture rename, and
carries neither the existence-token emission duty nor A22's canonicalization
duty. The joint-check line is an acknowledgement, not a clause." And
Prerequisites box four: "Two duties are added to RDR 0002's Refinement Context
Direction list, which currently carries neither".

**What is wrong from the PM view**: The RDR is candid, which is good — but the
*user consequence* of each fallback is understated. A22's "If wrong" says the
result is "an unevaluable refusal, never a plan (fail closed), but one that
misdiagnoses a producer defect as missing state." From the user's chair that is
precisely the outcome the Problem Statement forbids: the refusal names the
wrong thing. A user told "tag `Legal_Hold` was absent" when the real defect is
that the normalizer let two spellings through has been *misled*, not informed.
A22's fallback ("narrow A22 to kernel-side only … re-file canonicalization as
an obligation on 0002") makes the kernel correct and leaves the user misdirected
indefinitely, since 0002 is Draft with no deadline.

**Decision it blocks**: Whether the RDR can lock with a refusal payload whose
accuracy depends on a clause no document carries. Concretely: does the payload
need a third `reason` value (e.g. `unrecognized_key` when the key is absent from
the view but present in the tag declarations) so a producer defect is
distinguishable from missing artifact state? The closed reason set — "`absent`
… or `uncomparable`" (Normative Contracts, refusal payload clause) — currently
cannot express it.

---

## PM-7 — MEDIUM — A3's spike-shape gap means the one frozen test that pins user-visible nil-seam behavior must be re-decided, and the RDR does not say what it should decide

**Anchored passage**: A3 **Status**: "Pending — the BOUND is verified; the SHAPE
the spike proved is not the shape this RDR specifies. The spike retained
`Refusal.Guard` (retyped to `[]GuardAtom`) and added `UndecidedAtoms` BESIDE it
… Under the specified design that test must be RE-DECIDED, not re-encoded, which
is this assumption's own 'If wrong' branch." Paired with Normative Contracts →
"Frozen `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` keeps its
verdict but changes its REASON … To keep testing the nil-seam rule the fixture
needs a value atom over a PRESENT key (e.g. `reviews >= 3`)."

**What is wrong from the PM view**: The RDR asserts the verdict holds and only
the reason changes — but that assertion rests on a spike that proved a
*different* shape, and A3's own "If wrong" branch is "Phase 1 must reopen that
test's contract rather than re-encode it." The user-visible proposition at stake
is real and product-relevant: *does a missing evaluator rescue an escape edge?*
The RDR simultaneously (a) tells Phase 1 to change the fixture so the test keeps
testing the old proposition, and (b) leaves A3 Pending on whether the frozen
suite survives the specified shape. Those cannot both be settled by a re-run
that has not happened. The Testing Strategy preamble nevertheless claims the
baseline as established: "the frozen 154-test suite passes unchanged under the
reshape, so every row below is ADDED surface, not a re-encoding" — which is the
superset shape's result, not this RDR's.

**Decision it blocks**: Whether Phase 1 is a mechanical migration or a contract
reopening — i.e. whether shipping this RDR silently changes a frozen user-facing
guarantee about nil-seam escape behavior. Also blocks trusting the Testing
Strategy's "ADDED surface, not a re-encoding" framing, which currently
overstates what the spike proved.

---

## PM-8 — LOW-MEDIUM — The `RequiresOwned` narrowing (kata `xg7p`) is stated normatively but has no producer, so the outcome it promises is unreachable

**Anchored passage**: A21 **Status**: "Verified (composition); the PRODUCER
question stays open at §JD-3 and is carried as a Capability Dependency, not as a
claim of this RDR." Its Evidence: "§JD-3 remains NOT closed (unlike
§JD-1/§JD-2/§JD-6/§JD-7/§JD-12), and RDR 0002 carries the producer only as a
re-entry Direction ('Name the `RequiresOwned` producer'), not yet a clause. …
the deferral must survive into 0002's re-lock or the derive-from-`Writes`
definition has no implementer." Cross-check: JDR 0001 §JD-3 — "no layer is
obliged to populate it, and 0007's owned-before-guard ordering may quantify over
an always-empty set."

**What is wrong from the PM view**: The RDR's Metadata names kata `xg7p` as one
of its three coupled contracts and the Overrides line says this RDR "*fixes* the
meaning of `Row.RequiresOwned`." A definition with no producer fixes a doc
comment, not a behavior — the same failure mode the Consequences section
congratulates itself on escaping ("kata `xg7p` closes against behavior, not
against a doc comment"). If nobody populates the field, the entire
`owned_state_unavailable` diagnosis quantifies over the empty set and the user
never receives the "more precise owned_state_unavailable diagnosis" the
`Row.RequiresOwned` clause promises. Notably §JD-3 is *not* in the Prerequisite
list as a blocking item in the same terms as §D4 — Prerequisite box one says
"JDR 0001 §JD-3 landed in RDR 0002," but with no statement of what happens if
0002's re-lock does not carry it.

**Decision it blocks**: Whether kata `xg7p` may be closed on this RDR's lock, or
must stay open pending 0002's producer clause. Also blocks the honest scoping of
Testing Strategy rows 13 and 15, both of which assume a populated field.

---

## PM-9 — LOW — The "honest cost" of owned-before-unevaluable precedence hides the cause the RDR exists to reveal, and is accepted without a user check

**Anchored passage**: Normative Contracts → GATE, THEN COUNT: "The honest cost
is that a row failing both ways surfaces the write-dependency problem first."
Reinforced by the `disposition` mini-check row: "Row failing BOTH ways in one
survivor set | `owned_state_unavailable` first | … | owned payload only | loud
(the accepted cost)".

**What is wrong from the PM view**: In the both-ways case the user is told about
a *write* dependency while the guard-input absence — the thing this RDR exists
to surface — is suppressed entirely ("owned payload only"). The user fixes the
owned key, re-runs, and gets a second refusal for a different reason. That is a
two-round diagnosis loop for a single artifact state problem, and the RDR
accepts it purely on the grounds that it matches shipped `gate` ordering ("pinned
to shipped behavior, not to D8's original rationale … which this RDR's narrowing
invalidates"). The RDR itself notes the original rationale for that ordering is
now invalid, which removes the *design* reason for the precedence and leaves
only the *incumbency* reason — a weak basis for accepting a user-visible cost.

**Decision it blocks**: Whether the precedence should stay, or whether a
survivor failing both ways should carry both payloads (the refusal kind stays
`owned_state_unavailable`, but the guard payload rides along). That is a small
change that removes the second round-trip and costs nothing in the taxonomy.

---

## PM-10 — LOW — Phase 3's value-seam contract test is declared blocked, so a whole promised guarantee ships unverifiable

**Anchored passage**: Implementation Plan → Phase 3: "**Blocked on one RDR 0003
declaration.** `resolve.Tag.Value` is a bare `string`, so a set-valued tag
reaches the seam as one opaque string with no stated element encoding. The
`contains` leg of this contract test therefore cannot be written from any
current document, and neither can Testing Strategy scenario 8's present-key
half."

**What is wrong from the PM view**: Phase 3 is the RDR's only mitigation for the
first Risk ("RDR 0003's value evaluator drifts — e.g. treats an unparseable
present value as `false`"), which is a masking risk one layer down. Shipping the
mitigation with a hole in it, and routing the unblocking as a "REQUEST" in Phase
4 ("the REQUEST that 0003 declare the element encoding of a set-valued tag
value"), means the drift the mitigation targets is untested for one operator
class indefinitely. A request is not an obligation.

**Decision it blocks**: Whether `contains` may be in the shipped operator
vocabulary at RDR 0003's implement stage before the encoding is declared —
i.e. whether the vocabulary needs staging. Also blocks stating a completion
criterion for Phase 3.

---

## Not findings (checked, and the RDR handles them)

- The domain rule itself and its kernel enforcement site (§D4(b)) are stated
  clearly and are the right call for the user outcome: the Decision Rationale's
  deciding rows — "correctness-by-structure, drift enforcement, and the refusal
  payload" — are the user's rows, not the implementer's.
- A17's seam narrowing (`Evaluate(atom, value)`) does not cost the user
  anything; its "If wrong" branch keeps the domain rule kernel-side.
- A18's `uncomparable` reason genuinely prevents skew from reading as absence —
  a real user-outcome improvement, correctly identified.
- The MVV's self-caught contradiction (`trace` step 5) is exactly the kind of
  check that protects the user outcome; it is reported honestly.
