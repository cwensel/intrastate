# Cluster Reconcile Report — 0003 · 0006 · 0007

Model: claude-opus-5

Date: 2026-08-22
Iteration: **N = 1** (no prior `reconcile-report.md` under this cluster's
evidence dir; the `0001-0006` and `0002-0009` reports are different clusters).

## Membership

| RDR | Status | Implemented? | Profile | Relatedness |
| --- | --- | --- | --- | --- |
| 0003 | **Draft** [revised from Final 2026-08-12; re-verify A5] | No | large | Records the §JD-4 narrowing; owns guard atom grammar, finite-domain exhaustiveness, and (since 2026-08-21) the tag declaration model. |
| 0006 | **Draft** [demoted from Final 2026-08-21; re-verify A2, A5] | No | large | Emits the lint finding; owns graph-lint authority and the invariant taxonomy. |
| 0007 | **Final** | No | foundational | Holds the runtime veto that triggers the narrowing; owns the guard-evaluation domain rule. |

**Membership caveat, stated rather than papered over.** Stage 7.1's cluster
definition is *Final-and-unimplemented*. Only 0007 is Final here. The gate is
nevertheless the right venue, and this is not a stretch: JDR 0001 §JD-4 is a
three-way joint decision naming exactly these documents, and **both Drafts
name this pass as their closure venue by name** — 0003 A8 ("Venue: cluster
reconcile (0003 · 0006 · 0007)") and 0006 Refinement Context Direction 1.
Running the set-scoped checks against two Drafts costs nothing and is the only
way to close the entry.

The caveat *changes the mechanics*: a defect against 0003 or 0006 needs no
demotion, because neither is locked. Only 0007 was demotable, and the single
finding against it is a citation repair.

### Member revisions recorded (for the next iteration's delta-scope)

| Document | Revision at this pass |
| --- | --- |
| 0003 | `158164b` |
| 0006 | `158164b` |
| 0007 | `cb275c6` |
| **Home — JDR 0001** | `17fc26d` (pre-edit; this pass advances it) |

## The set-level finding that governs this report

The prior cluster gate found *silence* — obligations flowing one way into peers
that never received them. This cluster's defect is different and worse:
**everybody agrees, and the agreement cannot be written down.**

Three documents plus the JDR entry itself all recommend the *same* §JD-4 arm —
0003 records, 0006 cites — and not one of them takes it:

- JDR §JD-4: "**Recommended arm** (2026-08-21): 0003 records it and 0006 cites it."
- 0003 A8 `Plan`: "**Close on the assignment, not on matching prose.**"
- 0006 Direction 1: "Recommended: RDR 0003 records it … Whichever arm is taken,
  exactly one document states it."

The reason is structural and is the whole finding: **A8's closure requires an
edit to `docs/jdr/0001-resolve-kernel-seam.md`, and no per-RDR stage touches a
JDR.** Each document routed closure to "the cluster-reconcile venue" — this
pass. Individually rational, collectively paralysing. The whole-set critique's
premortem predicted precisely this pass producing a fourth document that
recommends the same arm and returns `NEEDS_DECISION`; that outcome is the
failure mode, so this report closes the entry instead.

Corroboration that the pattern self-perpetuates rather than converges: the
row-group gap was already logged by the prior cluster critique as C-20 and has
now survived **two** reconcile passes unresolved.

## Whole-set critique and pairwise scans

Both prompts ran. Outputs beside this report:

- `critique-set.md` — whole-set critique (single-model fallback, recorded as
  such per `2-critique.md`; no alt model reachable this session).
- `pairwise-0003-0006.md` — 10 findings, 4 joint blocking.
- `pairwise-0007-0006.md` — 5 findings.
- `pairwise-0007-0003.md` — see file.

Pairs run (3): `0003-0006`, `0007-0006`, `0007-0003` — the complete pair set for
a three-member cluster.

## Findings table

| Pair | Finding | TYPE | SEVERITY | Ownership | Disposition |
| --- | --- | --- | --- | --- | --- |
| 0003-0006 | F1 §JD-4 narrowing stated in 0003, not cited in 0006's body | gap | blocks-impl | joint | **JOINT-DECISION → §JD-4 CLOSED** |
| 0003-0006 | F2 warning tier / `graph-unprovable-coverage` reuse | contradiction | risks-impl | joint | **JOINT-DECISION → §JD-4 CLOSED** |
| 0003-0006 | F3 0006 finding contract has no atom-level field | gap | blocks-impl | joint | **JOINT-DECISION → §JD-4 consequent duty** |
| 0003-0006 | F4 row group defined in 0003, unacknowledged in 0006 | gap | risks-impl | joint | JOINT-DECISION → tolerance (0003 A10) |
| 0003-0006 | F5 no predecessor-reachability contract in 0006 | gap | blocks-impl | joint | JOINT-DECISION → tolerance (0003 A12) |
| 0003-0006 | F6 0006 appendix defect 4 stale against its own body | contradiction | cosmetic | single-RDR (0006) | NO CONFLICT — already closed in body |
| 0003-0006 | **F7 `single-valued grouping` attributed to a model that never declares it** | gap | blocks-impl | joint | **JOINT-DECISION → §JD-13 (new)** |
| 0003-0006 | **F8 escape rows: overlap exemption in 0006 vs union participation in 0003** | contradiction | blocks-impl | joint | **JOINT-DECISION → §JD-14 (new)** |
| 0003-0006 | F9 "claims closed coverage" undefined; 0003 reads it default-on | gap | blocks-impl | joint | JOINT-DECISION → §JD-14 rider |
| 0003-0006 | F10 round-trip across the seam | — | cosmetic | n/a | NO CONFLICT (negative result) |
| 0007-0006 | F1 `0007:159` asserts 0006 is `Final` | contradiction | cosmetic | single-RDR (0007) | **CITATION REPAIR — done this pass** |
| 0007-0006 | F2 0006's body never reaches 0007's veto case | gap | risks-impl | joint | **JOINT-DECISION → §JD-4 CLOSED** |
| 0007-0006 | F3 no lint category for runtime unevaluability | gap | cosmetic | joint | **CLOSED — reuse, no new code (§JD-4)** |
| 0007-0006 | F4 atom-level field absent (lint-side mirror of F3) | gap | risks-impl | joint | **JOINT-DECISION → §JD-4 consequent duty** |
| 0007-0006 | F5 0007 has no capability dependency on 0006 | gap | cosmetic | single-RDR (0007) | NO CONFLICT — stale A12 prose only |
| 0007-0003 | Peer-status sweep: 1 false claim of 17 | contradiction | cosmetic | single-RDR (0007) | **CITATION REPAIR — done this pass** |
| 0007-0003 | Declaration-model rehoming conflicts with nothing in 0007 | — | cosmetic | n/a | NO CONFLICT |

## Dispositions

### JOINT-DECISION — §JD-4 CLOSED (F1, F2, and 0007×0006 F2/F3)

This gate's central act. §JD-4's substance was never in dispute; only the
recording assignment and the warning-category question were open. Both are now
decided **in the home**, `docs/jdr/0001-resolve-kernel-seam.md` §JD-4:

- **RDR 0003 is the recording document.** Its clause already exists
  (`0003:781-786`). 0006 cites, never restates.
- **The two triggers genuinely differ and both survive.** 0006 fires on a
  *non-finite dimension*; 0003 on a *fully-finite product whose participating
  row can still refuse*. Verified independently: 0003's own disposition table
  keeps them as separate input classes (`0003:927` vs `0003:930`). A model whose
  every dimension is finitely declared but whose key is declared optional fires
  0003's clause and not 0006's. Not a wording artifact.
- **0003 is the recorder because the clause must name the refusing atom**, and
  `atom` occurs once in 0006 (`0006:70`), only to delegate atoms to 0003.
- **Warning category: no.** Lint mints no code and gains no non-blocking tier
  for this class. `graph-unprovable-coverage` already sits in 0006's *mandatory
  blocking* table (`0006:314`) **and** its MVV fixture matrix (`0006:639`), so
  reuse costs no new test scaffolding — only a widened trigger. 0006's existing
  advisory tier stays scoped to "redundant rows or unreachable rules"
  (`0006:302-304`) and MUST NOT absorb the withheld claim.

Why this is a decision and not a deferral: the home is the consumer's
umbrella-decision record, the two candidate documents are both open, and every
party had already converged on the arm. Deferring again would have made the
deadlock permanent by construction.

**Consequent duty recorded at the home (F3 / 0007×0006 F4).** 0006's finding
contract (`0006:363-367`) carries no atom-level field, so the atom-naming half
of 0003's clause has no producer. That is a payload extension, not a
restatement, and it is now written into §JD-4 as 0006's duty at refine.

### JOINT-DECISION — §JD-13 and §JD-14 (F7, F8, F9) — new at this gate

Neither finding appears on 0006's four-defect Refinement Context list nor on any
0003 assumption record. Both are `blocks-impl`. Both were verified independently
of the scan that raised them. Rather than book them as `Pending` peer
assumptions — the mechanism that produced the §JD-4 deadlock — each is decided
at the home now:

- **§JD-13 — `single-valued grouping` has no producer.** 0006 attributes it to
  "the RDR 0003 tag declaration model" (`0006:261-264`) and gates a mandatory
  blocking code on it (`0006:315`, MVV `0006:640`), but `single-valued` occurs
  **zero** times in 0003 and **zero** times in 0002. This is the same
  unhomed-producer defect the declaration-model rehoming was performed to close,
  surviving in one field the rehoming did not sweep. **Decided: 0003's
  declaration model gains the field**, beside value kind, finite domain,
  optionality, and set-element universe.
- **§JD-14 — escape rows.** 0006's invariant 3 grants an overlap exemption
  (`0006:286-288`) that 0003 normatively forbids (`0003:771-779`); the same
  fixture earns a blocking `graph-overlap` under one document and passes under
  the other. **0006 is split against itself** — its Load-Bearing Decision
  (`0006:412-414`) sides with 0003 against its own invariant 3. **Decided:
  0003's reading governs**; 0006 repairs invariants 3 and 4 to match its own
  Load-Bearing Decision. F9 ("claims closed coverage" undefined) rides here:
  0003's default-on reading (`0003:614-627`) governs, since an opt-in annotation
  would let the exhaustiveness guarantee be silently skipped exactly where it
  matters.

### Standing tolerances — F4, F5 (0003 A10, A12)

Not re-decided here; they already have a home and a scheduled producer. Both
are recorded obligations on 0006's refine, and 0006's A2 already carries them as
"Re-verify at refine" (`0006:124-129`). They stay open as tolerances against
§JD-4's neighbourhood, with the open questions named:

| Sibling | Open question | Home | Answered? |
| --- | --- | --- | --- |
| 0003 (A10) | Does 0006 read the same row-group division of labour? | 0006's refine | No — 0006 body still says "state/outcome pair", undefined |
| 0003 (A12) | Does 0006 publish a citable predecessor-reachability contract? | 0006's refine | No — invariant asserted, relation never published |

Note on A12: 0006 frames the check as *graph reachability* while 0003 requires a
*syntactic decision over declarations* (`0003:832-833`). Those are different
decision procedures over the same predicate; 0006's refine must state which.

### CITATION REPAIR — `0007:159` (0007×0006 F1)

Per the prompt's FACTS OF RECORD rule, a peer figure disagreeing with the
artifact of record is a citation repair — **no demotion, no tolerance, no
deviation entry**. Applied this pass.

The full peer-status sweep of 0007 adjudicated 17 assertions: **exactly one was
false.** Line 159 called 0006 "`Final`, tolerance §JD-4"; 0006 has been Draft
since 2026-08-21. 0007 already contradicted itself internally at lines 599 and
612, which correctly say 0006 is Draft. 0006 had tracked the staleness itself
(`0006:875-880`) as "not this RDR's to fix". It is fixed here, and 0007 **stays
Final** — a Context sentence naming a peer's status is not a contract, and
nothing in 0007's implementation depends on 0006 (its Capability Dependencies
table names 0006 zero times).

### NO CONFLICT — F6, F10, 0007×0006 F5

- **F6**: 0006's appendix defect 4 (finite-domain metadata unhomed) is stale
  against 0006's own body — `0006:261-264` and A2's Evidence already cite 0003's
  declaration model by name. The appendix is *behind* the document. Deleted at
  re-lock like the rest of the Refinement Context; no action owed.
- **F10**: the 0003→0006 seam is producer/consumer, not encode/decode. No
  inverse invariant is available to break. Identity-tuple stability across the
  seam was checked and is safe.
- **0007×0006 F5**: 0007 has no capability dependency on 0006. Its A12 evidence
  prose ("RDR 0003 … is silent on how an existence atom projects") is stale —
  true when written, false since 0003's rehoming answered it. Cosmetic; 0007's
  own re-lock sweeps it.

### DEFER-TO-IMPLEMENTATION — none

No finding qualified. F7 and F8 both change a clause's meaning and both touch
fenced normative text, which the prompt's guardrails exclude categorically. The
citation repair was cheaper to simply do than to defer.

## Observations outside this gate's authority

Recorded, not dispositioned — these are code-vs-spec or process facts no
cross-RDR pair owns. All independently verified this pass.

- **JDR §D1 is `Resolved` and has zero code effect.** `internal/resolve/resolve.go:185`
  still declares `Guard string`; the atom-carrying `Row` the entire cluster is
  written against does not exist. The last commit touching `internal/` or `cmd/`
  is `1c9c0ca` (2026-08-09); **129 commits have landed since, all specification.**
  Every atom-based clause in all three RDRs is written against a type nobody has
  written. This is a finalize/implement-gate question, not a cross-RDR one, but
  it is the largest risk in the set.
- **`intrastate lint` does not exist.** The member with the only user-visible
  surface (0006) claims blocking CI authority; grep for a `lint` verb in
  `internal/` and `cmd/` returns nothing.
- **Specification-to-code ratio.** 5,037 lines of RDR body and 19,125 lines of
  evidence for this three-RDR set — 24,162 lines — against 3,971 lines of Go
  (1,236 non-test), whose product surface is a `version` subcommand. The prior
  cluster's critique raised the same alarm; it has grown, not shrunk.
- **0003 is a Draft acting as the cluster's normative center of gravity.** It
  carries 5 of 15 assumptions `Pending` and 24 normative blocks, and is cited as
  the floor by a Final peer. The critique predicts its rehomed declaration model
  is the section most likely rewritten within six weeks of shipping, because its
  real consumers are 0002's normalizer and 0006's lint input contract — not
  0003. Worth watching; not a contradiction this gate can dispose of.
- **0007 is `Final` carrying a DOWNGRADED assumption (A12) whose home was an
  open JDR entry.** That entry is now closed, which retires the concern, but the
  general pattern — a Final RDR with a downgraded assumption homed in an open
  entry — is a finalize-gate question worth adopting as a checklist item.

## Verdict

**RECONCILED WITH TOLERANCES.**

- **No demotions.** 0003 and 0006 are already Draft; 0007's only finding is a
  citation repair applied in place, and it stays Final.
- **§JD-4 CLOSED** — both halves decided at the home, unblocking 0003 A8 (its
  Status line calls A8 "the one record still gating lock") and 0006 Direction 1.
- **Two new joint decisions homed** — §JD-13 and §JD-14, decided rather than
  booked as `Pending` peer assumptions.
- **Two standing tolerances** — 0003 A10 and A12, both with a scheduled producer
  in 0006's refine.

The cluster does **not** implement from here: 0003 and 0006 are Draft and must
re-lock first. What this gate removed is the reason they could not.

## Review Gate

- **Cluster membership right?** Yes, with the caveat stated in full above: only
  0007 is Final. The set is exactly the three documents §JD-4 names, and both
  Drafts named this pass as their closure venue. No peer was missed — 0002 was
  checked and its coupling here is carriage, already homed at §JD-3/§D1.
- **Iteration contract bound before the checks?** Yes. N=1 for this cluster, so
  the set-scoped checks ran in full; no delta-scoping applied. Member and home
  revisions are recorded above for the next pass to diff.
- **Every JOINT-DECISION genuinely joint?** Yes. §JD-4, §JD-13, and §JD-14 each
  sit between two documents where neither is solely wrong: in every case one
  document states a rule and the other either cannot reach it or contradicts it
  from its own taxonomy. The one finding where a document is affirmatively
  false about a fact (0007's stale status claim) is typed CITATION REPAIR, not
  joint — and not a demotion, per FACTS OF RECORD.
- **Does each tolerance name the open QUESTION, not just the home?** Yes — the
  A10/A12 table above names the question, the home, and whether it is answered.
- **Home form — a paragraph per decision, no sibling restatement?** Yes. §JD-4,
  §JD-13, §JD-14 each state the decision and the user-visible stake; derivation
  lives in this gate's evidence tree, not the home. 0006 cites 0003's clause
  rather than restating it — the restatement is itself named as a finding.
- **Nothing fenced or meaning-changing deferred?** Correct — no deferrals at
  all. F7/F8 both change clause meaning and were decided, not deferred.
- **A SPEC-DEFECT routed at the right re-entry scope?** None routed. Both
  candidate documents are already Draft with re-entry stages assigned (both
  refine), so a demotion would be a no-op; 0007's finding did not earn one.
- **Both prompts run?** Yes — `critique-set.md` plus three `pairwise-*.md`
  files, the complete pair set, all under this directory.
