Model: claude-opus-5[1m]

# Finalization Gate — RDR 0004: Accessor Execution Safety Model

- **RDR**: `0004-accessor-execution-safety-model`
- **Date**: 2026-08-21
- **Verdict**: **READY — Gate PASS**
- **Lock type**: re-lock. Locked once at `c149031`; demoted at `76a9e67` when
  JDR 0001's decisions were applied across the 0002–0009 cluster; re-entered
  through refine → resolve → grounding → 3amigo → critique → reconcile.
- **Mechanical pre-sweep**: PASS —
  `evidence/tooling-pass/tooling-pass.md` (one mechanical finding fixed in-pass,
  sweep re-run clean).

## 1. Contradiction Check

**No contradiction.** Two pairs were specifically tested rather than waved through,
because both are places where this RDR's own clauses could collide.

**Pair 1 — the read disjunction vs. the completeness guarantee.** "A read accessor
MUST return typed tag values or a typed refusal" sits alongside "MUST return the
tag set for exactly the keys it was asked for … or take the refusal branch."
These are jointly satisfiable: completeness constrains *which* branch the
disjunction takes, it does not add a third branch. The Desk Trace step 2 exercises
both with three witnesses (`output.txt:11`, `:12`, `:13`) and produces no
CONTRADICTION row. This pair is also the one the *previous* lock got wrong — the
older clause let a partial read take the typed-values branch, contradicting this
RDR's own rule that an indeterminate result is a refusal rather than a guess.
JDR 0001 §D3 caught it and resolved **(b)**; the clause now states the same rule
at both layers.

**Pair 2 — the cross-layer pair completeness creates.** Step 2's genuine-absence
*success* against the kernel behavior A8 cites (`missingOwned` deciding on map
presence via `TagSet.has`). These collide only if absence reaches the resolver as
a *present* tag — which the seam clause forbids: absence crosses as omission from
the owned snapshot, so the accessor's success branch and the kernel's
`owned_state_unavailable` agree rather than compete.

**Research Findings vs. Proposed Solution.** The callback-based FSM prior art
(Stateless `ExecuteEntryActions`) is cited *as a contrast* and explicitly rejected
in Alternative 5 — it is not evidence for the selected design. The OpenTofu and ADO
prior art support the design they are cited for (explicit checked state mutation;
typed external failure classification). No finding is enlisted against its own
sense.

**Planned features vs. stated principles.** The RDR declines undo/transactional
semantics (Alternative 4, Day 2 Operations) and the post-mutation clause respects
that: it names what an already-owned refusal class *means* and forbids
compensation. Declining to say what a refusal means is not the same decision as
declining to undo — the critique pass tested exactly this and it collapsed on the
evidence.

## 2. Assumption Verification

**Internally consistent — every record's Status/Method/Evidence agree, every
"If wrong" is non-empty.** Ten records, verified mechanically (CHECK 2/3/4/6) and
read individually.

**A1–A8 Verified**, each against evidence that supports *its specific claim*:

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

A8 is the one re-verified this re-entry, and it is the strongest record in the set:
its claim is that the absent/unread distinction *dies at the seam type*, and the
sweep independently confirmed `resolve.go::Tag` (`:20`) has exactly two fields with
no third state and `Input.Owned` (`:224`) is `[]Tag` with no error channel. The
claim is structural, not anecdotal.

**No Source Search self-reference** (CHECK 3: A3 is the only such record; it cites
production source). **No load-bearing `Docs Only`** (CHECK 4: zero `Docs Only`
records exist). All 29 cited symbols resolve on `main` (CHECK 5).

**A9 and A10 are `Pending — DOWNGRADED at Stage 6` and are carried into lock.**
This is the one judgment the gate must actually make, so it is made explicitly
rather than inherited from the reconcile.

The no-MVV-critical-defer rule bars deferring an assumption **the MVV depends
on** — a byte-parity reference, a fixture the MVV consumes, load-bearing external
behavior. A9 and A10 are the converse: they are properties the MVV **proves**.
A9's four rules are asserted by Scenario 1 arm 8, Scenario 6, Scenario 7, and
Scenario 3; A10's by Scenario 8. None is an input any scenario consumes. `MVV Test`
is defined as "pending implementation at lock time" — this is the state that Method
exists to express, not a deferral of the MVV itself.

Why a spike extension is not the alternative: both bind at the accessor→resolver
boundary the implementation builds, and the existing fixture double structurally
cannot exercise them — its `write` re-reads by cloning the tag map, so its re-read
cannot fail independently of the write; its `read` derives the key set from the
artifact and holds absence as an in-map sentinel. Extending the spike to witness
these would mean building the real binding, which is the implementation.

Both are survivable, with named fallbacks rather than hand-waving: A9 → move the
absence encoding into RDR 0001's `Input` shape (JDR 0001 §D3 option (c), which
§D3 explicitly keeps available for this purpose — verified at
`docs/jdr/0001-resolve-kernel-seam.md:165-169`); A10 → either the caller treats
every write refusal as possibly-applied, or the RDR claims the transactional
semantics Alternative 4 declined. A9's four rules are independent and retire
per-rule, with an explicit instruction to split A9 rather than flip it as a unit.

Neither is refuted, and no prose treats their properties as settled fact (CHECK 6).

**Precedent, verified not asserted**: RDR 0006 is `Final` while carrying A5 as
`Status: Pending` / `Method: MVV Test`, downgraded at its own Stage 6 on the same
reasoning. This is established practice in this repo.

## 3. Scope Verification

**The MVV is in scope, not deferred.** It is a fixture flow with one read, one
gate, and one write accessor over caller-supplied artifact roles, and the specific
proofs are named:

- **Scenario 1** — definition validation, **eight** arms, each asserting its *own
  named code*, not merely non-empty validation. Seven are witnessed
  (`output.txt:17-23`); the eighth (missing/empty requested-key set) is A9's.
- **Scenario 2** — read/gate dispositions, with the truncation case asserted on
  the refusal **carrying an empty value set** and the requested key set pinned in
  the definition, so a derived-key-set implementation fails rather than reporting
  success over a smaller set.
- **Scenario 3** — write read-back, including `read_back_incomplete` (A9).
- **Scenario 4** — replay disposition equality.
- **Scenario 5** — no package prints, asserted **positively** (returned value
  carries the refusal class) *plus* stdout/stderr capture requiring both empty.
- **Scenario 6** — absent required key carried through to the resolver asserting
  `owned_state_unavailable` (A9).
- **Scenario 7** — timeout precedence over `incomplete_read` (A9).
- **Scenario 8** — post-mutation reporting: applied-but-unverified sense, exactly
  one write invocation, no undo/retry/re-derivation (A10).

Every scenario carries a **negative control** in the Oracle Discriminability table,
and each control is discriminating rather than decorative — e.g. Scenario 6's
control ("a key the artifact *does* carry must resolve normally") exists precisely
to stop an implementation that omits every key from passing the refusal assertion
vacuously. Scenario 5 was flagged as the one weak absence-of-error oracle and
strengthened with a named capture control that is normative for the implementation
test.

The Scope Verification restatement in the RDR body was **narrower** than the MVV
until the Stage 6 reconcile propagated the two iter-2 widenings; that gap is
closed, and the two now agree.

## 4. Cross-Cutting Concerns

Three apply; the rest are omitted rather than N/A-bulleted.

- **Secret/credential lifecycle** — intrastate owns neither credentials nor remote
  resource lifecycle. External API accessors receive caller-provided artifacts and
  environment and return typed success/refusal only. This is A5's explicit scoping
  decision, and the Normative Contracts enforce its mechanism: accessors operate on
  caller-supplied roles and MUST NOT discover artifacts from ambient process state.
- **Concurrency model** — every invocation is context-bound with a declared,
  validated-positive timeout. Write accessors verify effects through same-role
  read-back rather than fire-and-forget. Missing or non-positive timeout metadata
  fails validation before execution.
- **Determinism** — the RDR claims **stable replay disposition**, explicitly not
  byte-identical output and not replay-stable hashes. A4 and Scenario 4 verify it;
  sorted map formatting prevents map-order drift. The Fidelity Table records the
  deliberate weakening (no artifact-level fidelity — the accessor observes tag
  values, not the artifact encoding) as a recorded limit rather than an oversight.

**Peer-owned policy, correctly delegated:** RDR 0002 owns the TOML table carrier
and accessor references; RDR 0003 owns predicate semantics over bound values;
RDR 0005 owns the user-facing CLI mapping (exit codes are explicitly out of scope
here — the table stops at the structured value the accessor package returns).

**One residual is disclosed rather than claimed closed.** The seam guarantee
delivers `owned_state_unavailable` only where a surviving row declares the key in
`Row.RequiresOwned`, and JDR 0001 §JD-3 records that no layer is yet obliged to
populate that field. The RDR does not paper over this: it is an **Open** row in
Capability Dependencies, a note under the Disposition Table, and a named
not-fully-closed shape in Failure Modes, with the reasoning that omission remains
strictly safer than a placeholder. Not resolvable here — RDR 0002 owns the
normalized row. This is a correctly-routed cross-RDR obligation, not a gate
blocker.

## 5. Proportionality

**Right-sized. Nothing to trim before locking.**

The RDR owns **one** load-bearing contract: accessor execution safety for declared
read/gate/write accessors — capability, refusal classes, timeout, and write
read-back. The two additions this re-entry made are inside that contract, not
beside it:

- **Read completeness** is the *success predicate of the read capability*, not a
  separate obligation. It arrived from JDR 0001 §D3(b), which routed it here as
  "one normative clause plus an MVV scenario" — and that is exactly what landed.
- **The seam-omission clause** constrains this RDR's own output. Naming the
  resolver-visible encoding of an absent key is what makes the branch rule mean
  anything; it leaves RDR 0001's `Input` shape unchanged.

The `large` Profile is retained and re-validated: this contract governs
authoritative artifact mutation with no ledger or undo, which is where the
mini-check tables (Disposition, Oracle, Fidelity, Desk Trace) earn their cost.

The document is long (1,034 lines) but the mass is load-bearing: ten evidence
records, sixteen normative clauses, five weighed alternatives, and four analysis
tables. CHECK 9 found **zero** Evidence fields over the 30-line budget. One
genuinely non-load-bearing block — 43 lines of Method-vocabulary guidance
duplicated from the engine README — was removed by the mechanical sweep this pass.

Deliberately *not* trimmed: the disclosed witness gaps ("not witnessed — see A9",
the "(allow arm)" parenthetical on the gate-deny row). Those annotations are the
RDR being honest about what its spike does and does not prove, and deleting them
to tidy the tables would convert disclosed gaps into hidden ones.

## Verdict

**READY — Gate PASS.** No blocker on any of the five items. The mechanical sweep
passes, no cluster re-entry note survives, the determinacy trigger is discharged
with a written `n/a` disposition, and the two Pending assumptions are a legitimate,
precedented, survivable downgrade with named MVV scenarios rather than an
unverified claim smuggled past the gate.

Locking to **Final**.
