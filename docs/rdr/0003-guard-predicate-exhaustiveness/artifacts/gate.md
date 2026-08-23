Model: claude-opus-5[1m]

# Finalization Gate — RDR 0003, Guard Predicate Exhaustiveness

- **RDR**: `0003-guard-predicate-exhaustiveness`
- **Date**: 2026-08-22
- **Verdict**: **PASS — READY. Locked to Final.**

Mechanical pre-sweep: `evidence/tooling-pass/tooling-pass.md` — **PASS** after two
MECHANICAL findings (a missing greppable `Premortem:` verdict line; four bare
`file:line` anchors in A5) were fixed in-pass and the sweep re-run. No
SUBSTANTIVE finding was raised.

## 1. Contradiction Check

**No contradiction blocks lock.** Research supports a symbolic guard-atom model
over declared tag domains, and the Proposed Solution keeps host callbacks and
free-form expression strings out of the contract, so lint can prove
finite-domain coverage and overlap. The declared tag domains the model rests on
are declared *here*, in the same document that consumes them — the circularity
an earlier draft carried (booking them as a producer request to RDR 0002 that
had no prospect of being met) is dissolved rather than deferred.

Three tensions were examined and each resolves:

- **The declaration-model rehoming is ratified bilaterally, not asserted
  unilaterally.** RDR 0002's refine (2026-08-22) replaced its non-normative
  prose schema list with a normative clause naming all five fields and stating
  that RDR 0003 is the normative home and the model MUST NOT be restated there;
  RDR 0007 (`Final`) independently records that value kinds and set universes
  "are RDR 0003's declarations". RDR 0002 keeps provenance, authoring location,
  and normalization carriage. The two documents cite rather than restate.
- **JDR 0001 §JD-14's overlap half is deliberately re-opened on evidence, and
  the RDR does not re-decide it unilaterally.** §JD-14 reasoned from "RDR 0001's
  runtime refusal of ambiguity"; the shipped kernel refutes that premise —
  `internal/resolve/resolve.go::Resolve` builds ordinary candidates from
  non-escape rows only, consulting escape rows solely to rescue a refusal. The
  coverage half stands unchanged; the overlap half is booked as **A17** for
  cross-document confirmation at RDR 0006's refine plus the next JDR 0001 touch.
  Disclosing a contested clause with a named venue is the correct handling, not
  a contradiction.
- **Two clauses were reconciled against the kernel at critique iteration 3** —
  `unless` restated in three-valued Kleene to match RDR 0007, and the coverage
  clause narrowed so a bare escape row closes coverage only for the failure
  classes it can actually rescue. Both were this RDR reading its own authority
  too loosely (false-greens reached by citing the kernel one function short of
  where control flow decides), not conflicts with a peer. Both are fixed in the
  live text.

§JD-4 is closed (2026-08-22, `0003-0006-0007` cluster gate) naming this RDR the
recording document for the `guard_unevaluable` narrowing, with RDR 0006 citing
the clause and minting no code. Guard-domain enforcement is settled at the
kernel (JDR 0001 §D4), scoping this RDR's evaluator to value semantics over a
present value.

## 2. Assumption Verification

**PASS.** 21 Evidence Records, independently recounted at this gate: **13
`Verified`, 8 `Pending`, 0 `Unverified`** — matching the RDR's own census and
its named Verified/Pending lists exactly.

- **Internal consistency**: every record's Status, Method, and Evidence agree,
  and every "If wrong" is non-empty. Verified by the sweep record-by-record.
- **Method vocabulary**: all 21 labels are in the sanctioned eight — Spike (A1),
  Derivation (A2), Design Decision (A3, A7, A9, A11, A13), Source Search (A4),
  Peer RDR (A5, A6, A8, A10, A12, A14, A16, A17, A18, A19, A20), MVV Test (A15,
  A21). The six records opened late in the rounds are all in-vocabulary.
- **No `Docs Only` record exists**, so nothing blocks lock on that clause.
- **No self-referential `Verified` stamp.** The one Source Search record (A4)
  cites three consumer-repo paths, none under this RDR or its artifact dir. Two
  stamps that were previously self-proving were corrected upstream and now hold:
  A1 was demoted off a spike that scanned operator *spellings* and re-verified
  on a Stage 6 evaluation harness that decides fixture rows against tag views
  (28 assertions pass / 0 fail); A2's derivation now runs over declarations this
  document owns rather than borrowed ones.
- **Cited anchors resolve on `main`**: `clierr.go::CLIError`,
  `respond.go::Fail`, `config.go::Load`, `resolve.go::Row` (carrying `RuleID`
  and `SourceLocator`), `Row.Guard`, `Row.rescues`, `GuardEvaluator`, `Refusal`,
  `assemble`, `escapeOrRefuse`, `Resolve`. The RDR's negative claims verify too
  (`conform` absent from the non-test kernel; no type satisfies
  `GuardEvaluator`). A5's four bare `file:line` anchors were rewritten to stable
  anchors in this pass.
- **Status consistency**: no `Pending` record has settled-fact prose resting on
  it. Each unsettled property is disclosed at its point of use — A12's
  reachability is explicitly carved out beside the clause that quantifies over
  it; A18's conformance premise is stated as a conditional with an inline
  "booked gap, not a silence" note; A20's MUST carries a consequent-duty note
  recording that neither surface can meet it today; A17 and A10 are flagged
  Contested in the `authority` table; A15 and A21 route to named MVV scenarios;
  A19 names RDR 0006 as the carrier owner. Prerequisites checkboxes agree with
  the record statuses item for item.

**The eight `Pending` records each carry a named plan and none blocks lock.**
The Stage 6 hard rule is the test applied: does the missing input pin the
property the MVV proves? A16 did — the projection clause makes `eq`/`in`/
comparisons single-value operators whose atoms over an unmarked dimension do not
project at all, so an unwritable `single_valued` marker left the MVV's only
passing exhaustiveness verdict (Scenario 2's complete partition) and Scenario 8's
green negative control unauthorable, degrading the MVV to refusals only. Stage 6
iteration 5 returned NOT RECONCILED on exactly that and routed to
`/rdr-refine 0002`. **That route was executed and landed**: RDR 0002 now carries
the marker in both its `[tags.<tag>]` schema enumeration and its normative
type-model clause (verified at this gate — the token occurs twice in RDR 0002,
in both required surfaces), and states per-atom `block` retention normatively
with an explicit MUST NOT against folding `unless` into `all`. **A16 and A14 are
`Verified` and the iteration-5 BLOCKER is discharged** (Stage 6 iteration 6:
RECONCILED). The remaining eight — A10, A12, A15, A17, A18, A19, A20, A21 —
were each tested against the same rule and none pins the MVV's property: A10,
A12, A17, A19 are peer agreements about this RDR's own definitions whose failure
mode is a false positive a fixture catches; A18 and A20 are producer fields whose
absence weakens Scenario 8's diagnostics without removing a verdict; A15 and A21
are measurements the MVV performs rather than inputs it consumes. Each leaves a
provable case standing.

## 3. Scope Verification

**PASS — in scope, executed during implementation, not deferred.**

The specific proof is a **production test fixture** (Phase 3, `Target-Flow
Fixture`) encoding one RDR flow slice and one kata flow slice as normalized
candidate rows, covering equality, enum membership, set containment, bounded
integer comparison, existence, and mixed `all`/`unless`. It must:

- prove one exhaustive and mutually exclusive scoped row group;
- detect one intentional gap and one intentional overlap that appear only in the
  multi-dimensional product, with source rule/context ids in the diagnostic;
- assert **positively** on a domain-exhaustive group carrying a possibly-absent
  guard key — lint emits the blocking inability-to-prove finding naming the row
  and the refusing atom. Asserting only the absence of a green result does not
  discharge this; a run that emits nothing must fail the test;
- carry a **negative control** — a row group whose keys are all declared
  always-present, which lint certifies green, proving the narrowing is tight
  rather than blanket;
- add at least one `contains` predicate over a declared set-valued tag before
  the full operator vocabulary is accepted.

**Authorability is closed.** Every case above is writable against this RDR's own
declaration model, including the green ones — the single exception (A16's
marker authoring location) closed 2026-08-22, restoring the MVV's positive proof
outcome. The remaining dependency is *sequencing*, not authorability: RDR 0007's
kernel reshape must land before Phase 1 builds against the atom-slice shape, and
Phase 2's group construction should wait on RDR 0006's refine answering A10/A17
rather than be built twice. The phase-gating table records this per phase.

## 4. Cross-Cutting Concerns

- **Incremental adoption** — addressed here. Guards over unbounded dimensions
  stay runtime-evaluable; lint refuses or downgrades the exhaustiveness claim
  for those dimensions rather than blocking all predicate use.
- **Memory management** — addressed here. Set-valued finite domains MUST use a
  deterministic symbolic or bitset-equivalent proof; naive powerset
  materialization is explicitly rejected, and "too large to prove" is a
  declared, model-independent published bound rather than an implementation
  accident.
- **Canonical-form / determinism** — addressed here, and the reason the
  determinacy trigger fires for this RDR. Successful matching and lint findings
  MUST NOT depend on source order; the canonical set-literal spelling is fixed
  normatively (unordered, duplicate-free, canonicalized before entering the
  identity tuple, repeats rejected at parse); and the same model MUST receive
  the same verdict on every conforming implementation. This RDR claims no
  byte-identical output, content-addressed identity, or replay-stable hash, so
  the hash/pre-image sub-checklist does not apply. The determinacy obligation is
  discharged by evidence: `evidence/repeatability/iter-2/` and `iter-3/` each
  carry `run-1.md` + `diff.md` + `dispositions.md`, `variant: lite (profile:
  large)`, with all five iteration-3 contract silences dispositioned `fixed`.
- **Peer-owned policy this RDR conforms to** — exact-one selection is RDR 0001's;
  source-row identity is RDR 0002's; CLI envelope mapping is RDR 0005's;
  blocking graph-lint findings are RDR 0006's; the kernel seam, domain rule, and
  `guard_unevaluable` payload are RDR 0007's.

Versioning, build-tool compatibility, licensing, deployment model, IDE
compatibility, secret/credential lifecycle, concurrency, and character encoding
do not apply and are omitted rather than N/A-bulleted.

## 5. Proportionality

**PASS — right-sized for lock; nothing to trim.**

**Contract count: one.** This RDR is the sole author of a single independent
load-bearing contract — the guard-predicate grammar and its finite-domain
exhaustiveness semantics, together with the tag declaration model that grammar
quantifies over. The 26 `normative` blocks are clauses of that one contract, not
26 contracts: the split test is independent sole-authored contracts, not block
count or word count. The declaration model is not a second seam but the input
alphabet — a grammar that cannot say what a tag may hold cannot prove anything
about it — so an implementer holds one contract in working memory, not two.
Alternative placements were tested and rejected on evidence (RDR 0002 has no
normative tag-type vocabulary, no type or value-domain validation category, and
a charter the declaration's *meaning* falls outside of; a new RDR would split one
tag's declaration across two documents and add an edge from all of
0002/0003/0006/0007). The model was homeless, not homed elsewhere — a gap closed,
not a boundary crossed.

**Profile re-validated: `large` holds.** Contract axis: this RDR locks a grammar,
which is the `large` trigger. It is not `foundational` — the declaration model
has cross-RDR *consumers*, but RDR 0002 owns authoring location and normalization
carriage, so this RDR is not the module-spanning producer that trigger names.
Accretion axis does not floor it: `Seam Lineage` reads "no prior accretion", so
there is no ≥2 point-fix history and no accretion disposition is owed. The
Metadata field's form is correct — value plus one clause naming the contract,
with no matrix or provenance prose left from the template.

**The lenses the Profile requires all ran on the revised text**, so the latch's
backstop finds no disagreement: `grounding`, `3amigo`, and `critique` (dual-model,
two distinct stamps diffed behind a barrier) are complete through iteration 3,
plus `repeatability` (lite) because the Determinacy trigger fired. Neighboring
surfaces are explicitly delegated to peer RDRs rather than locked here.

## Lock actions

- `Status` → **Final**.
- Finalization Gate section body replaced with the pointer line to this file.
- README index row for 0003 flipped to `Final`.
