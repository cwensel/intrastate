Model: claude-opus-5[1m]

# Tooling Pass — RDR 0004 (mechanical adherence sweep), iteration 2

RDR: `0004-accessor-execution-safety-model`
Status line at entry: `Draft [revised from Final 2026-08-24; re-verify none — JDR
0001 §D5 (§JD-15) makes read-back assert absence for a `<clear>` write; … §D7
(§JD-17) retires the "cannot be rebound" sentence]`
Date: 2026-08-24

## Scope of this iteration

This is the **second** re-lock of 0004. The RDR locked at `c149031`, was demoted
at `76a9e67`, re-locked at `48cb2a4` (iteration 1 of this evidence dir), and was
demoted **again** by cluster-reconcile iteration 3 (`567600f`) as
**RE-LOCK-ONLY / `re-verify none`** over JDR 0001 §D5 (JD-15) and §D7 (JD-17).
Refine `1a5ddca` landed those answers. Stage 3 → Stage 7 is the sanctioned path;
no lens re-run is owed.

**LOOP-BREAKER check.** Prior report read: `evidence/tooling-pass/tooling-pass.md`
(iteration 1). Its single mechanical finding was the 43-line older-template
"Method vocabulary" block — deleted then, confirmed still absent now. **No finding
in this report is a re-report of an iteration-1 finding**, so no routing loop.

## Verdict

**PASS** — after three MECHANICAL citation repairs applied in-pass and the sweep
re-run. No SUBSTANTIVE finding; no stage return owed.

## CHECK 1 — Template section coverage

**PASS.** All 35 Required (spine) sections present and substantive. No verbatim
template bracket (`[Conditional`, `[Required`, `[Resource]`, `[Capability]`), no
`_Draft placeholder._`, no `TBD`, no "see above", no seed-skeleton header. The
older-template Method-vocabulary block deleted in iteration 1 has not returned.

`Phase 2: Operational Activation` remains a **Conditional** section cleanly
omitted — PASS per the false-positive guard, not Missing. Instance phases
(Phase 1 Accessor Model, Phase 2 Executor Boundary, Phase 3 Write Read-Back
Verification, Phase 4 CLI Integration Hook) are legitimate instance content.

Stage-5 mini-check tables (Disposition, Oracle Discriminability, Fidelity, Desk
Trace) present and substantive.

The `### Critical Assumptions` nesting under `## Research Findings` is unchanged
and is **again not raised** — the iteration-1 report records it as a repo-wide
authoring convention across all nine RDRs, not a regression these rounds
introduced. Unchanged disposition.

## CHECK 2 — Method-label vocabulary

**PASS.** Eleven Critical Assumption records (A1–A11); every Method is exactly
one sanctioned label, no compound labels and no glosses to strip:

| Method | Count | Records |
| --- | --- | --- |
| Spike | 5 | A1, A2, A6, A7, A8 |
| MVV Test | 4 | A4, A9, A10, **A11** |
| Source Search | 1 | A3 |
| Design Decision | 1 | A5 |

A11 is the record this re-entry added (JDR 0001 §D5 clear-as-removal). Its
`MVV Test` label is sanctioned and correct for a property the MVV proves.

## CHECK 3 — Source Search self-reference

**PASS.** A3 is the only `Source Search` record. Its anchors resolve to
production source — `internal/cli/clierr::CLIError` (`clierr.go:48`),
`internal/cli/clierr::ExitCodeFor` (`:110`), `internal/cli/respond::Fail`
(`respond.go:133`) — not to this RDR and not under its artifact directory. The
spike transcript is cited as corroboration alongside the source anchors, not as
the anchor itself. No offender. (Verified by delegated read.)

## CHECK 4 — Docs Only on load-bearing claims

**PASS.** Zero records carry `Method: Docs Only`. No `Docs Only` string survives
anywhere in the document.

## CHECK 5 — Symbol resolution of anchors

**PASS.** 29 `path::Symbol` anchor instances / 24 distinct symbols; **all resolve
in their cited file**. No phantoms, no moves, no bare `file:line` standing in for
a symbol. (Verified by delegated read.)

- 14 spike symbols in `evidence/spikes/main.go` — `main`, `read`, `write`,
  `refusal`, `replay`, `formatMap`, `validateDefinitions`, `validationCases`,
  `absentValue`, `newArtifacts`, `newPartialArtifacts`, `newSparseArtifacts`,
  `expectedTagKeys`, `withTimeout`.
- 9 kernel symbols in `internal/resolve/resolve.go` — `Tag` (`:20`, two fields,
  corroborating A8), `Input.Owned` (`:227`, `[]Tag`), `missingOwned` (`:444`),
  `TagSet.has` (`:125`), `TagSet.matches` (`:132`), `Refusal.MissingOwned`
  (`:267`), `Row.RequiresOwned` (`:180`), `escapeOrRefuse`, `RefusalKinds`.
- `internal/resolve/adversarial_test.go::TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`
  (`:114`) — exact name match.
- `internal/cli/config::Load` resolves.
- External prior-art symbols (`StateRepresentation::ExecuteEntryActions`,
  `…Async::ExecuteEntryActionsAsync`, `TaintCommand::Run`, `stateMgr.WriteState`,
  `stateMgr.PersistState`, `Client::transitionWorkItem`) do not resolve in this
  repo — **expected and out of scope**: cited as Prior Art in Investigation,
  Alternative 5, and References, and confirmed **never** used as Source Search
  evidence for any assumption.

**Witness-citation accuracy.** All cited `output.txt` lines match the 23-line
transcript verbatim, including the two quotes the Desk Trace reproduces inline
(`:9`, `:10`). The seven rows annotated "not witnessed — see A9/A10/A11" are
genuinely absent from the transcript — the four new A11 clear rows, the
`read_back_incomplete` row, the timeout-with-unresolved-keys row, and the
post-mutation timeout row. **Honest disclosure, confirmed, not a finding.** The
Disposition Table's "Definition invalid (8 shapes)" row correctly cites seven
witnessed arms (`:17-23`) and flags the eighth as A9's.

**Referenced-document accuracy.** `docs/cli-output-contract.md` exists. JDR 0001
§D3 (resolved (b), option (c) kept available — exactly what A9's "If wrong"
cites), §D5 (all four landing clauses), §D7(ii) (capability tables, `keys` binds
readers and writers, `read_back = true` on writers, same id in `[read.x]` and
`[write.x]` is two identities), and §D7(iv) (write replaces the whole value;
read-back compares for **equality**) all resolve and say what the RDR claims.

### C5 findings — FIXED IN-PASS (mechanical, citation repair ×3)

**The RDR cited §JD-3 as an open decision. It was CLOSED 2026-08-23** — before
this re-entry — by the `0002-0009` iteration-2 gate
(`docs/jdr/0001-resolve-kernel-seam.md:391`): RDR 0002's normalizer is the named
producer, carrying a MUST at
`docs/rdr/0002-transition-table-as-reviewable-data.md:1125` ("`Row.RequiresOwned`
has no authored form. The normalizer MUST derive it … sorted, duplicate-free set
of tag keys named by the rule's write block and clear list"). The RDR's
parenthetical "the field appears zero times in 0002" was also false — it occurs
12 times there, including as that MUST.

Three sites repaired:

1. **Disposition Table conditional-loudness note** — "JDR 0001 §JD-3 records that
   no layer is yet obliged to populate that field" → §JD-3 (CLOSED) names 0002's
   normalizer as producer; the set is the write block plus clear list, and RDR
   0007 fixes its meaning as post-guard *write-dependency* keys.
2. **Capability Dependencies row** — status `**Open**` → `Decided, unimplemented`;
   row renamed from "populated by some layer" to "derived by the normalizer";
   the stale zero-occurrences claim removed.
3. **Failure Modes residual paragraph** — "the residual is a cross-RDR
   obligation" → the residual is a property of that peer-owned *scope*.

One `## References` line added naming §JD-3 (CLOSED) and its producer.

**Why MECHANICAL, not SUBSTANTIVE.** The repair is a stale citation of a peer
document's state, not a design change. Nothing in the contract, the refusal
classes, the seam clause, or the MVV moves. Critically, **the conclusion the
three sites draw survives intact**: the absent-key row's loudness *is* still
conditional, because `RequiresOwned` is derived as writes ∪ clears and RDR 0007
fixes it as post-guard write-dependency keys — so a key a row consumes only
through `Row.Match` or a guard remains outside the set. Only the *reason* changes,
from "no producer is yet obliged" to "the producer's derived scope excludes
match/guard-only keys". No new evidence and no design call was required, which is
the test the Stage 7 prompt applies. Routing this to Refine would return it
unchanged — conformance and citation repair are outside Refine's contract.

Sweep re-run after the repairs: clean.

## CHECK 6 — Status consistency

**PASS.** No settled-fact prose leans on an unsettled assumption. Every property
belonging to the three Pending assumptions is marked at its site: Disposition
Table rows carry "not witnessed — see A9/A10/**A11**", Oracle Discriminability
rows are tagged `(A9)`/`(A10)`/`(A11)`, and Testing Strategy scenarios 6/7/8/**9**
name their assumption. Scenario 2's prose separates what the spike witnesses from
what is "new and unwitnessed — they are A9's".

The A9/A10/A11 clauses in Normative Contracts state **requirements this RDR
legislates**, not verified empirical facts — the correct register for an
unverified normative clause, not a settled-fact contradiction.

**Checklist-vs-gate agreement.** The Prerequisites checkbox now reads "A1-A8
Verified; **A9, A10, and A11 Pending**" — it was widened to include A11 by the
refine pass, so it agrees with the Critical Assumptions section and with the gate
record. No disagreement. (The cluster-reconcile report named the unchecked
Prerequisites boxes "for the re-lock gate"; they are unchecked because the four
prerequisites are genuinely unmet — three peer RDRs unimplemented and three
assumptions Pending — which is the accurate state, not a conformance defect.)

## CHECK 9 — Evidence-field budget (advisory, never blocks)

**0 fields over budget.** Longest: A9 at 10 lines, A10 at 7, **A11 at 3**; A1–A8
are 1 line each. Budget is 30 lines. Nothing to report, no truncation proposed.

## Determinacy / repeatability trigger

**Discharged — unchanged from iteration 1.** `Profile: large` and the Normative
Contracts remain adjacent to ambiguous-field-ownership territory (absent vs
unreadable vs cleared key encoding), so the trigger was evaluated, not skipped.
`evidence/grounding/iter-2/dispositions.md:40-45` carries the written disposition:

> `Determinacy trigger: **n/a**` — the contract legislates execution safety and
> refusal classing at an I/O boundary. It carries no parse/deparse,
> import/export, compose/decompose, hashing, identity, or migration behavior, no
> data-model field whose ownership could be read more than one way … and the MVV
> rests on per-call branch dispositions rather than multi-step transformation
> fidelity.

The §D5 addition does not disturb this: clear-as-removal is one more branch
disposition at the same boundary, not a transformation with a round-trip. A
written `determinacy: n/a — <reason>` disposition is what the Stage 7 prompt
accepts in place of `repeatability/run-1.md` + `diff.md`. No repeatability row is
owed.

## Cluster re-entry note

**Absent — PASS.** `grep 'Refinement Context'` → no hit. The refine commit
`1a5ddca` dropped it by name in its subject. The cross-RDR defect that demoted
this RDR is closed, not dropped: §D5's four landing clauses and §D7's two are all
present in live text (verified individually below), and the `Row.RequiresOwned`
question is now correctly recorded as decided rather than open.

**JDR landing verification** (the defect that drove the demotion):

| Landing item | Home | Status in 0004 |
| --- | --- | --- |
| `<clear>` write removes the key | §D5 | **present** — normative clause; Approach; Technical Design |
| read-back asserts the key is absent | §D5 | **present** — normative clause; Fidelity Table; Desk Trace step 5 |
| clearing an absent key succeeds | §D5 | **present** — normative clause; Disposition Table row |
| a read yielding `<clear>` is unreadable | §D5 | **present** — normative clause; Absent-vs-unreadable decision; Disposition row |
| MVV scenario for a clearing rule | §D5 | **present** — Scenario 9, with Oracle row and negative control |
| definition is the capability-table entry | §D7(ii) | **present** — Technical Design; Load-Bearing "Definition shape" |
| writers carry `keys` | §D7(ii) | **present** — "`keys` is the binding on readers and writers alike" |
| read-back is equality | §D7(iv) | **present** — normative clause "for equality against the held value — a write replaces the whole value, so containment is not equality" |
| "cannot be rebound" prose repaired | §D7(ii) | **present** — now "the same id may appear in both `[read.x]` and `[write.x]` — two identities, not a rebinding" |

The two unfenced "present" wordings the answer-check flagged (`0004:215`, `:511`)
are both repaired: Approach now reads "an assigned value must be held exactly, and
a key written as the reserved `<clear>` sentinel must be absent"; the Fidelity
Table's owned-tag row now carries an absence arm.

**Noted, not a finding.** §D7(ii) carries an internal tension the iter-3
answer-check flagged for the gate's eye: "`keys` exists only on readers and
`read_back` only on writers" (describing per-capability decoder shapes, in
contrast to the retired unified `[accessors.<id>]`) versus "Readers and writers
both declare `keys`". The operative reading is the second — it is reinforced by
the loader rules ("a written or cleared key not in exactly one writer"), by the
`Fields:` enumeration, and by the *Lands in 0004* item "writers carry `keys`".
0004 landed that reading. Not a 0004 defect; if §D7(ii)'s first sentence is ever
tightened, it is the home's edit.

## Findings summary

| Check | Verdict | Findings |
| --- | --- | --- |
| C1 template coverage | PASS | none (Conditional omission correct; nesting convention unchanged) |
| C2 Method vocabulary | PASS | none (11/11 sanctioned, incl. new A11) |
| C3 Source Search self-reference | PASS | none |
| C4 Docs Only load-bearing | PASS | none (zero records) |
| C5 symbol resolution | PASS (after fix) | 3 mechanical citation repairs (§JD-3 stale-as-open), fixed in-pass; 24/24 symbols resolve; witnesses accurate |
| C6 status consistency | PASS | none (checklist widened to A11, agrees with gate) |
| C9 evidence budget (advisory) | PASS | 0 over budget |

**PASS — no findings outstanding; proceed to the Gate's written responses.**
