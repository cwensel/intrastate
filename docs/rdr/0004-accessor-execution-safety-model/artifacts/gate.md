Model: claude-opus-5[1m]

# Finalization Gate — RDR 0004: Accessor Execution Safety Model

- **RDR**: `0004-accessor-execution-safety-model`
- **Date**: 2026-08-24
- **Verdict**: **READY — Gate PASS.** Re-locked to Final.

**Re-lock (second).** This RDR locked at `c149031`, was demoted at `76a9e67`
when JDR 0001's decisions were first applied across the 0002–0009 cluster, and
re-locked at `48cb2a4`. Cluster-reconcile **iteration 3** (`567600f`) demoted it
again as **RE-LOCK-ONLY / `re-verify none`**, Stage 3 → Stage 7, over JDR 0001
§D5 (JD-15, `<clear>` semantics) and §D7 (JD-17, definition shape and the
"cannot be rebound" sentence). Refine `1a5ddca` landed both. This record
**overwrites** the 2026-08-21 responses; where a finding is unchanged from that
lock it is restated rather than re-derived, and the deltas are called out.

**Mechanical pre-sweep**: `evidence/tooling-pass/iter-2/tooling-pass.md` —
**PASS** after three MECHANICAL citation repairs (a stale §JD-3 "still open"
claim at three sites) were fixed in-pass and the sweep re-run. **No SUBSTANTIVE
finding was raised**, so no stage return was owed. No routing loop: iteration
1's sole finding (an older-template vocabulary block) is gone and is not
re-reported.

**Delta from the 2026-08-21 lock.** One new Critical Assumption (**A11**), one
new normative clause (`<clear>` reserved / clear-is-removal / read-back-absent /
idempotent-clear / unreadable-on-read), one new MVV scenario (**9**) with its
Oracle row and negative control, four new Disposition Table rows, an absence arm
on the Fidelity Table's owned-tag row, the Identity decision's rebind repair, and
the §JD-3 citation repairs above.

## 1. Contradiction Check

**No contradiction.** Three pairs were tested rather than waved through.

**Pair 1 — the read disjunction vs. the completeness guarantee.** Restated from
the prior lock, unchanged: "a read accessor MUST return typed tag values or a
typed refusal" alongside "MUST return the tag set for exactly the keys it was
asked for … or take the refusal branch". Jointly satisfiable — completeness
constrains *which* branch the disjunction takes; it does not add a third. Desk
Trace step 2 exercises both with three witnesses (`output.txt:11`, `:12`, `:13`)
and yields no CONTRADICTION row.

**Pair 2 — the cross-layer pair completeness creates.** Restated, unchanged:
step 2's genuine-absence *success* against the kernel behavior A8 cites
(`missingOwned` deciding on map presence via `TagSet.has`). These collide only if
absence reaches the resolver as a *present* tag, which the seam clause forbids —
absence crosses as omission, so the accessor's success branch and the kernel's
`owned_state_unavailable` agree rather than compete.

**Pair 3 — NEW this pass: clear-as-removal vs. the read-back equality clause.**
The §D5 landing makes a `<clear>` write's read-back assert **absence**, while the
§D7(iv) landing makes read-back compare the held value for **equality**. These
are one predicate, not two: the RDR states it as "an assigned value must be held
exactly, and a key written as the reserved `<clear>` sentinel must be absent"
(Approach), and the Fidelity Table carries it as "value equality (assignment) or
absence (clear) over the planned key set". A11 says so explicitly — it "extends
the expected-value predicate to absence rather than re-opening it", leaving A2
and A7 untouched. The Desk Trace step 5 assertions and the Round-Trip invariant
were both widened to match, so no site still asserts presence for a cleared key.
The two unfenced "present" wordings the iteration-3 answer-check named
(`0004:215`, `:511` at `526c481`) are the sites that were repaired; both are
verified rewritten.

**Research Findings vs. Proposed Solution.** Unchanged and still sound. The
callback-based FSM prior art (Stateless `ExecuteEntryActions`) is cited *as a
contrast* and rejected in Alternative 5 — never enlisted as evidence for the
selected design. OpenTofu and ADO prior art support the propositions they are
cited for (explicit checked state mutation; typed external failure
classification).

**Planned features vs. stated principles.** The RDR declines undo/transactional
semantics (Alternative 4, Day 2 Operations), and both post-mutation clauses
respect that: they name what an already-owned refusal class *means* and forbid
compensation. The new clear clause does not breach it either — removal is a
write, not a rollback.

## 2. Assumption Verification

**Internally consistent — every record's Status/Method/Evidence agree, every
"If wrong" is non-empty.** Eleven records, verified mechanically (CHECK 2/3/4/6)
and read individually.

**A1–A8 Verified**, each against evidence supporting *its specific claim*
(restated from the prior lock; nothing this re-entry disturbed them):

| ID | Method | Basis |
| --- | --- | --- |
| A1 | Spike | read/gate/write bound as declared capabilities; `output.txt:1-8` |
| A2 | Spike | read-back detects owned-tag divergence; `output.txt:8-9` |
| A3 | Source Search | `clierr::CLIError`, `clierr::ExitCodeFor`, `respond::Fail` — all resolve |
| A4 | MVV Test | replay disposition equality; `output.txt:14-15` |
| A5 | Design Decision | explicit scoping — credentials/remote lifecycle outside intrastate |
| A6 | Spike | seven validation arms, each its own named code; `output.txt:16-23` |
| A7 | Spike | collateral non-owned mutation caught; `output.txt:10` |
| A8 | Spike + source | three read dispositions (`output.txt:11-13`), grounded on `Input.Owned`/`Tag`/`missingOwned` |

**No Source Search self-reference** (CHECK 3: A3 is the only such record; it
cites production source). **No load-bearing `Docs Only`** (CHECK 4: zero such
records). All 24 distinct cited symbols resolve in their cited files (CHECK 5).

**A9, A10, and A11 are Pending and are carried into lock.** This is the judgment
the gate must actually make, so it is made explicitly rather than inherited.

The no-MVV-critical-defer rule bars deferring an assumption **the MVV depends
on** — a fixture it consumes, a byte-parity reference, load-bearing external
behavior. All three are the converse: properties the MVV **proves**. A9's four
rules are asserted by Scenario 1 arm 8, Scenario 6, Scenario 7, and Scenario 3;
A10's by Scenario 8; **A11's by Scenario 9**. None is an input any scenario
consumes. `MVV Test` is defined as "pending implementation at lock time" — this
is the state that Method exists to express.

**A11 specifically, since it is new.** Its claim is that a `<clear>` write is
verifiable as removal. The spike structurally cannot witness it, and the record
says why in terms that check out on inspection: `spikes/main.go::write` applies
planned tags with `maps.Copy` and re-reads by `clone` — there is no `delete` call
anywhere in the spike and no fixture carries the sentinel, so a removal has
nothing to assert and the re-read cannot fail independently of the write.
Extending the spike to witness it would mean building the real binding, which is
the implementation. Its "If wrong" names a concrete, pre-existing fallback: JDR
0001 §D5 option **(c)**, a typed absent value on the seam — verified present in
§D5 as a weighed-and-rejected option that remains available.

Both other Pending records keep their named fallbacks: A9 → JDR 0001 §D3 option
(c), which §D3 explicitly keeps available for this purpose; A10 → either the
caller treats every write refusal as possibly-applied, or the RDR claims the
transactional semantics Alternative 4 declined. A9's four rules are independent
and retire per-rule, with an explicit instruction to split A9 rather than flip it
as a unit.

None of the three is refuted, and no prose treats their properties as settled
fact (CHECK 6) — each site carries "not witnessed — see A9/A10/A11" or an
`(A9)`/`(A10)`/`(A11)` tag.

**Precedent, verified not asserted**: RDR 0006 is `Final` while carrying A5 as
`Status: Pending` / `Method: MVV Test`, downgraded at its own Stage 6 on the same
reasoning. Established practice in this repo.

**Checklist agreement.** The Prerequisites box was widened this re-entry to
"A1-A8 Verified; **A9, A10, and A11 Pending**", so it agrees with the Critical
Assumptions section and with this record. The four Prerequisites boxes remain
unchecked — correctly: three peer RDRs are unimplemented and three assumptions
are Pending. That is the accurate state, not a conformance defect. The
cluster-reconcile report named them "for the re-lock gate"; the gate's finding is
that they are honest.

## 3. Scope Verification

**The MVV is in scope, not deferred.** A fixture flow with one read, one gate,
and one write accessor over caller-supplied artifact roles, with named proofs:

- **Scenario 1** — definition validation, **eight** arms, each asserting its own
  named code. Seven witnessed (`output.txt:17-23`); the eighth is A9's.
- **Scenario 2** — read/gate dispositions; truncation asserted on the refusal
  **carrying an empty value set**, with the requested key set pinned in the
  definition so a derived-key-set implementation fails rather than reporting
  success over a smaller set.
- **Scenario 3** — write read-back, including `read_back_incomplete` (A9).
- **Scenario 4** — replay disposition equality.
- **Scenario 5** — no package prints, asserted **positively** plus stdout/stderr
  capture requiring both empty.
- **Scenario 6** — absent required key carried through to the resolver asserting
  `owned_state_unavailable` (A9).
- **Scenario 7** — timeout precedence over `incomplete_read` (A9).
- **Scenario 8** — post-mutation reporting: applied-but-unverified sense, exactly
  one write invocation, no undo/retry/re-derivation (A10).
- **Scenario 9 — NEW (A11)** — a clearing write against three artifacts (holds
  the key; lacks the key; binding stores the literal) plus a read of an artifact
  holding the literal. Asserts success with the key **absent**, idempotent
  success, `read_back_mismatch` for a stored literal, and `incomplete_read` for
  the read. The MVV paragraph and Phase 3 were both widened to carry it, so the
  restatement and the scenario agree.

Every scenario carries a **negative control** in the Oracle Discriminability
table, and each is discriminating rather than decorative. Scenario 9's control is
the scenario-3 assignment fixture, which "must still assert presence and equality
— a read-back that asserts absence for every write fails it". That is the right
control: it is exactly the failure mode a naive absence-checking implementation
would have.

## 4. Cross-Cutting Concerns

Three apply; the rest are omitted rather than N/A-bulleted.

- **Secret/credential lifecycle** — intrastate owns neither credentials nor
  remote resource lifecycle. External API accessors receive caller-provided
  artifacts and environment and return typed success/refusal only (A5's explicit
  scoping). The Normative Contracts enforce the mechanism: accessors operate on
  caller-supplied roles and MUST NOT discover artifacts from ambient process
  state.
- **Concurrency model** — every invocation is context-bound with a declared,
  validated-positive timeout. Writes verify effects through same-role read-back
  rather than fire-and-forget. Missing or non-positive timeout metadata fails
  validation before execution.
- **Determinism** — the RDR claims **stable replay disposition**, explicitly not
  byte-identical output and not replay-stable hashes. A4 and Scenario 4 verify
  it. The Fidelity Table records the deliberate weakening (no artifact-level
  fidelity) as a recorded limit. The determinacy/repeatability trigger is
  discharged by the written `n/a` disposition at
  `evidence/grounding/iter-2/dispositions.md:40-45`; the §D5 addition does not
  disturb it, being one more branch disposition at the same boundary rather than
  a transformation with a round-trip.

**Peer-owned policy, correctly delegated:** RDR 0002 owns the TOML carrier, the
capability tables, and the binding validations; RDR 0003 owns predicate semantics
over bound values; RDR 0005 owns the user-facing CLI mapping (exit codes
explicitly out of scope — the table stops at the structured value the accessor
package returns). RDR 0007 owns `Row.RequiresOwned`'s meaning.

**One residual is disclosed rather than claimed closed — and its basis was
corrected this pass.** The seam guarantee delivers `owned_state_unavailable` only
where a surviving row declares the key in `Row.RequiresOwned`. The prior lock
justified that hedge by citing §JD-3 as an *open* decision; §JD-3 was CLOSED
2026-08-23, naming RDR 0002's normalizer as the producer. **The residual
survives the correction, on a better basis**: the normalizer derives the set as
the rule's write block plus clear list, and RDR 0007 fixes its meaning as
post-guard *write-dependency* keys — so a key a row consumes only through
`Row.Match` or a guard is outside the set **by construction**, not by an unclosed
decision. Omission remains strictly safer than a placeholder. Correctly routed to
0002/0007, not a gate blocker.

## 5. Proportionality

**Right-sized. Nothing to trim before locking.**

The RDR owns **one** load-bearing contract: accessor execution safety for
declared read/gate/write accessors — capability, refusal classes, timeout, and
write read-back. This re-entry's addition is inside that contract, not beside it:
clear-as-removal is the **expected-value predicate of the write capability**
extended from "held exactly" to "held exactly, or absent". It arrived from JDR
0001 §D5, which routed it here as remove-key / read-back-absent /
idempotent-clear / unreadable-on-read clauses plus an MVV scenario — and that is
exactly what landed, no more. Six lines of normative text, one assumption, one
scenario, four disposition rows.

Why it belongs here rather than in 0002: §D5 resolved (a) precisely because the
alternatives (a `Clears []string` row field, a typed absent value on `Tag`)
reopen RDR 0001's `Row`/`Tag` and RDR 0009's fenced predicate for one field. This
RDR is where the sentinel's *meaning at execution* lives, and its read-back is
the only place a stored-literal clear is detectable.

The `large` Profile is retained and re-validated: the contract governs
authoritative artifact mutation with no ledger or undo, which is where the
mini-check tables (Disposition, Oracle, Fidelity, Desk Trace) earn their cost.

The document is 1,016 lines and the mass is load-bearing: eleven evidence
records, sixteen normative clauses, five weighed alternatives, four analysis
tables. CHECK 9 found **zero** Evidence fields over the 30-line budget (longest:
A9 at 10, A10 at 7, A11 at 3).

Deliberately *not* trimmed: the disclosed witness gaps ("not witnessed — see
A9/A10/A11", the "(allow arm)" parenthetical on the gate-deny row). Those are the
RDR being honest about what its spike does and does not prove; tidying them away
would convert disclosed gaps into hidden ones. Scenario 9's four rows are all so
annotated.

## Verdict

**READY — Gate PASS.** No blocker on any of the five items. The mechanical sweep
passes, no cluster re-entry note survives, every §D5 and §D7 landing item is
present in live text, the determinacy trigger is discharged with a written `n/a`
disposition, and the three Pending assumptions are a legitimate, precedented,
survivable downgrade with named MVV scenarios and named fallbacks rather than
unverified claims smuggled past the gate.

Locking to **Final**.
