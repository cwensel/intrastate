Model: claude-opus-5[1m]

# 3 Amigos iter-2 — Persona 2: Implementer

Delta-scoped to the re-entry: A5, A7, A8, and the passages that state or depend on
the narrowed lint promise (Technical Design narrowing paragraph, Normative
Contracts exhaustiveness-narrowing clause, Decision Rationale, Contradiction
Check, Prerequisites, MVV).

Persona question: if I started coding this Monday, what would I ask in the first
hour?

---

## HIGH

### P2-1 — A5 claims per-atom source identity that the cited atom shape does not carry

**Anchored passage.** Load-Bearing Decisions › Identity: "a guard predicate is
identified by its source rule/context id plus its position within `all` or
`unless`". A5 Evidence: "a row carries parsed atoms (key, operator token,
literal, block), not an opaque predicate string, so there is no reconstruction
step that could lose identity."

**Grounding.** The cited shape is exactly four fields. RDR 0007's SEAM clause
(`docs/rdr/0007-guard-predicate-totality.md:1249-1266`) fixes them by name —
`Key`, `Operator`, `Literal`, `Block` — with `Block` a two-constant string type
(`BlockAll`, `BlockUnless`). There is no index, ordinal, or span field. Row-level
identity ships (`internal/resolve/resolve.go::Row` carries `RuleID` and
`SourceLocator`, confirmed at resolve.go:169-185), but that is *row* identity,
not *atom* identity. A5's own Evidence concedes the point it needs and does not
close it: the identity the assumption asserts is "(rule id, position in block)",
and "position in block" exists nowhere in the shape A5 cites as its proof.

**Clarification request.** Where does an atom's position come from at
implementation time? Three candidates, and I cannot pick one from the text:
(a) the slice index in `Row.Guard []Atom` is the position, in which case this
RDR must say the slice order is authoritative and stable — which collides with
the same section's "Successful row matching MUST NOT depend on source order"
and with MVV Scenario 5 ("Reorder authored rows and guard atoms … lint findings
are unchanged"); (b) `(Key, Block)` is the identity — but then two atoms over one
key in one block (RDR 0007 A5's conjoined value row, `X exists = true` +
`X eq v`) are indistinguishable in a diagnostic; or (c) a fifth field is needed
on the atom, which is a change request to RDR 0007's `Final` SEAM clause, not a
citation of it.

**Decision blocked.** The `UndecidedAtom` mirror and every overlap/coverage
diagnostic payload. Normative Contracts requires "Overlap and coverage
diagnostics MUST name the source rule id or context id that contributed each
predicate involved in the finding" — I cannot emit "each predicate" without a
per-atom handle, and I cannot choose the handle without breaking either
Scenario 5's reorder invariance or diagnostic distinguishability.

**Note.** This is squarely in delta scope: the re-entry's stated purpose for A5
is to "cite rather than restate" the kernel-fixed atom shape, and the citation is
where the gap appears.

---

### P2-2 — A7 `Pending` leaves Phase 2 with no defined behavior for the only TOTAL operator, and the fixture is the failing case

**Anchored passage.** A7 Status `Pending`; `trace` Step 5: "**GAP — booked as A7,
not a contradiction.** Without the presence dimension this step has no defined
result." A7 Plan: "Blocked on that request: RDR 0002's tag declaration carries
name, provenance, value kind, and optional accessor reference, with no
optionality field."

**Grounding.** Verified: RDR 0002 line 217-218 declares exactly "tag name,
provenance (`owned`, `observed`, `recognized`), value kind, and optional
accessor reference for observed or owned read-back" — no optionality field, and
RDR 0002 is itself `Draft`
(`docs/rdr/0002-transition-table-as-reviewable-data.md:9`). The derivation in
`evidence/research/iter-2-projection-derivation.md` recommends P1 (per-key
presence axis) and rejects P2, but the RDR body carries no normative clause for
either. Meanwhile `evidence/spikes/guard-fixture.toml` — this RDR's own
representative fixture — has `foundational-to-cove` guarding
`cluster_eligible exists = true` where `cluster_eligible`'s declared domain is
`[true, false]`, i.e. complete with no presence element.

**Clarification request.** What do I build in Phase 2 for a row group containing
an `exists` atom while A7 is unresolved? The RDR is explicit that the step "has
no defined result", and A7's "If wrong" enumerates both branches as harmful
(dropping the atom is unsound; refusing every group with an `exists` atom
disqualifies the RDR's own fixture). Concretely: is P1 the intended clause, and
does implementation sequence *after* RDR 0002 lands optionality — or is there an
interim behavior (e.g. treat undeclared-optionality keys as never-absent, so
`exists = true` is tautologous and `exists = false` empty)?

**Decision blocked.** Phase 2 ("Finite-Domain Lint Semantics") cannot be
specified, and Phase 3's fixture cannot be run — the fixture's second rule is
precisely the unspecified case. Prerequisites already states "**A7 is Pending**",
so this is not an objection to the status, it is a request for the interim
implementation contract, which nothing in the document supplies.

---

### P2-3 — The narrowing clause puts a MUST on a `Final` sibling this RDR has no authority over

**Anchored passage.** Normative Contracts, exhaustiveness-narrowing clause:
"RDR 0006, which holds the JDR 0001 §JD-4 tolerance, MUST carry this same
narrowing for the findings it emits; the two documents MUST NOT state it
differently."

**Grounding.** RDR 0006 is `Final [joint decision → JDR 0001 §JD-4]`
(`docs/rdr/0006-graph-lint-authority-and-guarantees.md:9`). Its exhaustiveness
clause (lines 345-349) narrows only for a non-finite dimension. Grep confirms
RDR 0006 contains zero occurrences of `guard_unevaluable` and zero of "narrow" —
A8's divergence claim is accurate. A8's own Plan agrees the fix is external:
"RDR 0006 is `Final`, so closing this needs a route-back, not a silent edit
here."

**Clarification request.** A normative clause in RDR 0003 cannot make RDR 0006
true. If I implement lint against RDR 0006's clause alone (which is `Final`, and
which RDR 0006's own Predecessors list makes a downstream consumer of this RDR),
I certify domain-exhaustive groups this RDR forbids. Which document do I code
against before A8 closes, and does this clause bind *my* implementation of
RDR 0006's findings, or is it a note-to-sibling that should live in A8 rather
than in a `normative` block?

**Decision blocked.** Whether the exhaustiveness verdict emitter withholds on a
possibly-absent key. The `authority` table already marks this row
"**Contested — A8.** … Must converge before lock" — but the normative block
states it as settled MUST, so an implementer reading only Normative Contracts
gets the opposite signal from one reading only the authority census.

---

## MEDIUM

### P2-4 — "Refuse or downgrade" is stated five times without saying which, or what "downgrade" produces

**Anchored passage.** Normative Contracts: "lint MUST refuse or downgrade an
exhaustiveness claim for that dimension"; and again "lint MUST refuse or
downgrade the exhaustiveness claim rather than silently capping enumeration".
`disposition` table, three rows: "**Refuse or downgrade** the claim". A2:
"lint must refuse or downgrade the exhaustiveness claim." MVV Scenario 3:
"lint refuses or downgrades the exhaustiveness claim for that dimension/product."

**Grounding.** JDR 0001 §JD-4 records the open half as "which document records
the narrowing **and whether lint gains a warning category**"
(`docs/jdr/0001-resolve-kernel-seam.md:238-241`). The warning-category question
is the "downgrade" arm. RDR 0006 line 345-349 states only the refuse arm
("MUST emit a blocking inability-to-prove finding"). A8 carries the recording
question but is silent on the warning-category half of §JD-4.

**Clarification request.** Is "downgrade" a distinct outcome I must implement
(a non-blocking advisory finding, per RDR 0006's "Non-blocking advisories may
exist"), or is it a synonym for refuse? A disjunctive MUST with no selection rule
is untestable — MVV Scenario 3's "Expected" accepts either branch, so a build
that always refuses and a build that always downgrades both pass. And if
downgrade is real, RDR 0006's clause says blocking, which is a second
`guard_unevaluable`-shaped divergence A8 does not currently name.

**Decision blocked.** Finding severity and CLI exit behavior for the
inability-to-prove case; whether `graph-lint-failed` fires.

---

### P2-5 — A8's convergence target is under-specified: which document's *wording* wins, and what "the same narrowing" means mechanically

**Anchored passage.** A8: "Verification is a route-back to RDR 0006 confirming
it adopts the same wording, or a §JD-4 disposition assigning the recording to
one document."

**Grounding.** These two resolutions produce different implementations. "Same
wording in both" duplicates the promise across documents — the exact
single-source failure A8's own "If wrong" says §JD-4 exists to prevent ("the two
documents state one lint promise two ways"). "Assign to one document" produces a
citation from the other. The RDR asks for either without saying which it expects,
while its "If wrong" argues against the first.

**Clarification request.** If A8 closes by RDR 0006 adopting the same wording,
have we not created the duplication §JD-4 forbids? Which arm should cluster
reconcile pursue, and does this RDR's Normative Contracts clause survive that
arm (if the narrowing lands in RDR 0006, this RDR's clause becomes a restatement
and should become a citation per the engine's "one home per contract; cite, never
restate" — `docs/jdr/0001-resolve-kernel-seam.md:60`)?

**Decision blocked.** Whether I code the narrowing into the RDR 0003 evaluator's
exhaustiveness verdict or into RDR 0006's finding emitter — the two land in
different packages.

---

### P2-6 — Row-group grouping is delegated to RDR 0006, but RDR 0006 does not define it in the shape this RDR consumes

**Anchored passage.** Technical Design, narrowing-adjacent paragraph: "For lint,
RDR 0002 supplies the normalized candidate rows and RDR 0006 supplies the
graph-lint grouping context, such as one source-state/recognized-outcome
selection group."

**Grounding.** RDR 0006's invariant 4 (line 277-279) reads "for each state/outcome
pair **that claims closed coverage**". A grep across RDR 0002, 0003, and 0006 for
a mechanism by which a group *claims* closed coverage returns nothing —
`0006:696` mentions "a model claims closed coverage" in an MVV scenario, but no
document says where the claim is authored. This RDR's A2 says only "For every
exhaustiveness-eligible row group", never defining eligibility.

**Clarification request.** How does a row group become exhaustiveness-eligible?
Is every state/outcome group implicitly claiming coverage (in which case every
group with an unbounded or absent-capable dimension emits a finding, and the
"refuse or downgrade" valve fires constantly), or is it opt-in via a declaration
someone must own? The `authority` table has a row for per-tag optionality
("RDR 0002 — **not yet declared**") but no row for coverage-claim ownership.

**Decision blocked.** Whether the lint pass is opt-in or default-on; directly
determines whether the MVV's "one intentional gap" fixture even triggers.

---

## LOW

### P2-7 — Prerequisites marks RDR 0007 sequencing `[x]` while the shipped kernel still carries the superseded shape

**Anchored passage.** Prerequisites: "[x] RDR 0007 is the normative home of the
guard seam … RDR 0007 is `Final`; its kernel reshape is specified, not yet
implemented, so this RDR's implementation sequences after it." A5: "the shipped
kernel still carries `Row.Guard string`".

**Grounding.** Confirmed at `internal/resolve/resolve.go` — `Row.Guard string`
(line 185) and `GuardEvaluator.Evaluate(guard string, view TagSet) GuardResult`
(line 91-93). Both are the pre-reshape form the `authority` table calls
"superseded, not an arm".

**Clarification request.** "Sequences after it" — after RDR 0007 *locks* (done)
or after RDR 0007 *ships* the reshape (not done)? The checkbox is `[x]`, which
reads as satisfied, but the shipped seam is `Evaluate(guard string, view TagSet)`
and this RDR's evaluator is specified as `Evaluate(atom, value)` that "MUST NOT
read the tag view". I cannot write the evaluator against the interface that
exists today.

**Decision blocked.** Implementation start date — whether RDR 0003 Phase 1 can
begin before RDR 0007's kernel reshape lands, or is hard-blocked on it.

### P2-8 — A7's Plan names one blocking producer request; Capability Dependencies names two, and only one has a home

**Anchored passage.** Capability Dependencies, last row: "Per-tag optionality
declaration (which keys may be absent) | RDR 0002 | Requested | A7's presence
dimension needs it … **Requested alongside the set-valued element encoding RDR
0007 Phase 4 routes here.**"

**Grounding.** A7's Plan and the Assumption Verification block name only the
optionality declaration as blocking. The set-valued element encoding appears
only in this table cell. Yet Phase 3 and Testing Strategy both make `contains`
gating: "the implementation MVV must add at least one `contains` predicate over
a declared set-valued tag before the full operator vocabulary is accepted", and
`contains` requires "a declared element universe" (operator matrix) that the
fixture's tag declarations do not demonstrate — `guard-fixture.toml` declares no
set-valued tag at all.

**Clarification request.** Is the set-valued element encoding a second `Pending`
blocker (it gates the MVV's `contains` case and the full operator vocabulary),
or is it tracked somewhere I should read? It has no assumption ID, no status,
and no entry in the `authority` census.

**Decision blocked.** Whether Phase 3's fixture can be authored at all, and
whether the "closed operator vocabulary" normative clause is provable at lock.

---

## Out of scope

None. Every finding above anchors to a delta-scoped passage (A5, A7, A8, the
narrowing paragraph, Normative Contracts, `authority`/`disposition`/`trace`
tables, Prerequisites, Capability Dependencies, or MVV).
