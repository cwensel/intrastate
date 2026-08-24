Model: claude-opus-5[1m]

# Finalization Gate — RDR 0003, Guard Predicate Exhaustiveness

- **RDR**: `0003-guard-predicate-exhaustiveness`
- **Date**: 2026-08-24
- **Verdict**: **PASS — READY. Re-locked to Final.**

**Re-lock.** This RDR locked Final 2026-08-22, was demoted to Draft by
cluster-reconcile iteration 3 (`567600f`) as **RE-LOCK-ONLY / `re-verify
none`** over a single mis-routed fenced clause, and was refined by `aa39d86`.
This record **overwrites** the 2026-08-22 responses; where a finding below is
unchanged from that lock it is restated rather than re-derived, and the deltas
are called out explicitly.

Mechanical pre-sweep: `evidence/tooling-pass/iter-2/tooling-pass.md` — **PASS**
after five MECHANICAL findings (a surviving 43-line older-template legend block
in `Critical Assumptions`; three bare RDR 0007 `file:line` anchors; two stale
RDR 0002 quotations in A14) were fixed in-pass and the sweep re-run. **No
SUBSTANTIVE finding was raised**, so no stage return was owed. No ROUTING-LOOP:
both fixes from the 2026-08-22 sweep survive intact, and this round's C5
findings land on records that round did not name.

## 1. Contradiction Check

**No contradiction blocks lock.** Research supports a symbolic guard-atom model
over declared tag domains, and the Proposed Solution keeps host callbacks and
free-form expression strings out of the contract, so lint can prove
finite-domain coverage and overlap. The declared tag domains the model rests on
are declared *here*, in the same document that consumes them.

**The defect that caused the demotion is discharged, and bilaterally.** The
fenced clause at the old `0003:1065-1067` routed this RDR's two rejection rules
onto *RDR 0006 findings*. JDR 0001 §D7(iii) decides otherwise: "Two load
categories in 0002, rules supplied by 0003: `malformed tag declaration` … and
`malformed predicate atom` … **0006 mints nothing for either**; 0005 maps them
under §JD-8." The RDR now routes exactly there — in A4's Evidence, in the
`disposition` table's three rows, and in the load-category passage — and **RDR
0002 (`Final`) reciprocally names both categories** as "rejection rules" that
"RDR 0003's rules reject". The two documents cite rather than restate, in both
directions. §D6 is cited as decided where the participation and can-refuse
clauses read the authored block, per the re-entry Direction.

Three standing tensions were re-examined and each still resolves:

- **The declaration-model rehoming is ratified bilaterally, not asserted
  unilaterally.** RDR 0002 carries a normative clause naming all five fields and
  stating that RDR 0003 is the normative home and the model MUST NOT be
  restated there; RDR 0007 (`Final`) independently records that value kinds and
  set universes "are RDR 0003's declarations". RDR 0002 keeps provenance,
  authoring location, and normalization carriage.
- **§JD-14's overlap half is settled, not contested.** It was decided
  2026-08-22 and **corrected 2026-08-23** to the two-population reading —
  overlap is checked among ordinary rows and among escape rows per declared
  failure class, never escape-vs-ordinary. JDR 0001 states "0003 A17 closes on
  it", and **A17 is now `Verified`** against RDR 0006's escape-row clause. This
  is a delta from the prior lock, where A17 was `Pending` and booked as
  contested.
- **The block domain is three-valued in RDR 0002 and two-valued here — by
  design, not by drift.** RDR 0002 widens `block` with a third `match` member on
  JDR 0001 §D6's authority, declares the widening as its own, and tells an
  implementer reading RDR 0007 alone to "widen to three here". This RDR and RDR
  0007 spell the two-valued guard domain §D1 fixes. A14's evidence was re-quoted
  from RDR 0002's current text in this pass; the guard-side carriage the record
  depends on is unaffected.

§JD-4 is closed naming this RDR the recording document for the
`guard_unevaluable` narrowing, with RDR 0006 citing the clause and minting no
code. **§JD-16 is answered** (decided 2026-08-24 by §D6; JDR 0001 records "0003
consistent, qualifier cleared") and was correctly dropped from the Metadata
qualifier by the refine. Guard-domain enforcement is settled at the kernel
(§D4), scoping this RDR's evaluator to value semantics over a present value.

## 2. Assumption Verification

**PASS.** 21 Evidence Records, independently recounted at this gate: **17
`Verified`, 4 `Pending`, 0 `Unverified`**. This is the material delta from the
2026-08-22 lock (13/8): `aa39d86` closed **A10, A12, A17 and A19** against RDR
0006, which is `Final`.

- **Internal consistency**: every record's Status, Method, and Evidence agree,
  and every "If wrong" is non-empty (21/21). The four newly-`Verified` records
  each dropped their `Evidence needed` / `Plan` scaffolding cleanly and now
  carry `Evidence` citing a fenced clause in a `Final` peer; each cited clause
  was confirmed to exist in RDR 0006 (row-group division of labour; the
  reachability relation; the escape-row two-population overlap MUST together
  with `graph-coverage-closed-by-escape`).
- **Method vocabulary**: all 21 labels are in the sanctioned eight — Spike (A1),
  Derivation (A2), Design Decision (A3, A7, A9, A11, A13), Source Search (A4),
  Peer RDR (A5, A6, A8, A10, A12, A14, A16, A17, A18, A19, A20), MVV Test (A15,
  A21). The refine relabeled no Method and added or removed no record.
- **No `Docs Only` record exists**, so nothing blocks lock on that clause.
- **No self-referential `Verified` stamp.** The one Source Search record (A4)
  cites three consumer-repo paths, none under this RDR or its artifact dir.
- **Cited anchors resolve on `main`**: `clierr.go::CLIError`, `respond.go::Fail`,
  `config.go::Load`, and in `internal/resolve/resolve.go` — `Row` (carrying
  `RuleID` and `SourceLocator`), `Row.Guard`, `Row.rescues`, `GuardEvaluator`,
  `Refusal`, `assemble`, `Resolve`. The negative claims verify too (`conform`
  absent from the non-test kernel). Every JDR 0001 section the refine introduced
  resolves — §D6 and §D7(iii) — as do §D1, §D4, §JD-4, §JD-8, §JD-13, §JD-14,
  §JD-16, §JD-18. Four bare `file:line` anchors and two stale peer quotations
  were rewritten to stable form in this pass.
- **Status consistency**: no `Pending` record has settled-fact prose resting on
  it. A15's bound adequacy is stated open inside the clause it gates; A18's
  conformance premise is an explicit conditional carrying a "booked gap, not a
  silence" note; A20's MUST carries a consequent-duty note recording that no
  surface can meet it today; A21's marker-coverage question is disclosed inline.
  Prerequisites checkboxes agree with the record statuses item for item.

**The four `Pending` records each carry a named plan and none blocks lock.** The
Stage 6 hard rule is the test applied: does the missing input pin the property
the MVV proves? **A18 and A20** are producer fields routed to JDR 0001 §JD-18
(open; the conforming-view enforcer and the runtime atom carrier); their absence
weakens Scenario 8's diagnostics without removing a verdict. **A15 and A21** are
measurements the MVV *performs* rather than inputs it consumes, discharged by
MVV Scenarios 3 and 4 at implementation. Each leaves a provable case standing.

The one remaining unchecked Prerequisite — RDR 0007's kernel reshape — is
explicitly scoped in the document as gating *implementation sequencing, not
lock*: RDR 0007 is `Final`, so the specification is settled, while the shipped
kernel still carries `Row.Guard` as a `string` (verified at this gate,
`resolve.go:185`).

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
  rather than blanket.

**Authorability is closed**, including for the green cases: every case is
writable against this RDR's own declaration model, and the `single_valued`
marker the positive half depends on is writable because RDR 0002's schema
carries it (A16, `Verified`). The remaining dependency is *sequencing*, not
authorability — RDR 0007's kernel reshape must land before Phase 1 builds
against the atom-slice shape. The phase-gating table records this per phase, and
it now reads "Nothing open on peers" for Phase 2, since A10, A12, A17 and A19
closed against RDR 0006.

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
  determinacy trigger fires. Successful matching and lint findings MUST NOT
  depend on source order; the canonical set-literal spelling is fixed
  normatively (unordered, duplicate-free, canonicalized before entering the
  identity tuple, repeats rejected at parse); and the same model MUST receive
  the same verdict on every conforming implementation. This RDR claims no
  byte-identical output, content-addressed identity, or replay-stable hash, so
  the hash/pre-image sub-checklist does not apply. The obligation is discharged
  by evidence: `evidence/repeatability/iter-2/` and `iter-3/` each carry
  `run-1.md` + `diff.md` + `dispositions.md`, `variant: lite (profile: large)`,
  all five iteration-3 findings dispositioned `fixed`.
- **Peer-owned policy this RDR conforms to** — exact-one selection is RDR 0001's;
  source-row identity, authoring location, normalization carriage and the two
  load categories carrying this RDR's rejection rules are RDR 0002's; CLI
  envelope mapping is RDR 0005's; blocking graph-lint findings are RDR 0006's;
  the kernel seam, domain rule, and `guard_unevaluable` payload are RDR 0007's.

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
Alternative placements were tested and rejected on evidence; the model was
homeless, not homed elsewhere.

**Profile re-validated: `large` holds.** Contract axis: this RDR locks a
grammar, which is the `large` trigger. It is not `foundational` — the
declaration model has cross-RDR *consumers*, but RDR 0002 owns authoring
location and normalization carriage, so this RDR is not the module-spanning
producer that trigger names. Accretion axis does not floor it: `Seam Lineage`
reads "no prior accretion". The Metadata field's form is correct — value plus
one clause naming the contract, with no matrix or provenance prose.

**The lenses the Profile requires all ran on the revised text.** `grounding`,
`3amigo`, and `critique` (dual-model, two distinct stamps diffed behind a
barrier) are complete through iteration 3, plus `repeatability` (lite) because
the Determinacy trigger fired. The re-entry scope was **RE-LOCK-ONLY /
`re-verify none`** — a cross-reference fix disturbing no assumption — so no lens
re-run was owed and the latch's backstop finds no disagreement.

**Trim check:** `aa39d86` net-shrank the document (450 insertions vs 632
deletions), and this pass removed a further 43 lines of older-template
instructional text. No Evidence field exceeds the 30-line advisory budget
(longest: A14 at 26 after the re-quote, then A17/A16/A12 at 21). Nothing further
to trim before locking.

## Lock actions

- `Status` → **Final**.
- Finalization Gate section body replaced with the pointer line to this file.
- README index row for 0003 flipped to `Final`.
