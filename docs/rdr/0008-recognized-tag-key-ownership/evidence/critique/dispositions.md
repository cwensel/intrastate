Model: claude-opus-5[1m]

# critique lens — resolve pass (iteration 1)

Origin ledger = the reconciled pass-A ↔ pass-B diff (`diff.md`), entries
**D-1…D-22**. Reconciliation is by passage anchor, since `C-N` IDs are per-file
and do not correspond across the two passes. Every entry exits exactly once.

Grounding gate applied per entry against (1) code on `main`, (2)
`{RDR_RESOURCES}` (`docs/cli-output-contract.md`, `CLAUDE.md`), and (3) the
RDR's own decided text — including cove's F-1/F-5 and the 3amigo dispositions,
so re-raises are cited rather than re-fixed.

| # | Origin | Disposition | Section touched |
| --- | --- | --- | --- |
| D-1 | A(C-1) + B(C-3,C-4) | **fixed** | A4 (rewritten, both clauses) |
| D-2 | A(C-2) + B(C-4) | **fixed** | A4 migration inventory; Phase 2 |
| D-3 | A(C-3) + B(C-5) | **fixed** | Normative block 2 (predicate vs write position; non-coverage stated) |
| D-4 | A(C-5) + B(C-9) | **fixed** | A8 scope limit; new A10 |
| D-5 | A(C-7) + B(C-8) | **fixed** (fork collapsed, not escalated) | Normative block 4 (reason given for the unconditional clause) |
| D-6 | A(C-9) + B(C-12) | **fixed** | Testing Strategy (split consequences binding); MVV kernel half |
| D-7 | A(C-10) + B(C-13) | **dismissed-with-cite** (standing charted item) | — |
| D-8 | A(C-14) + B(C-14) | **fixed** (design change) | Normative block 5 (predicate now enforces it); scenario 9; Approach 5–6 |
| D-9 | A(C-13) + B(C-10) | **fixed** (deferral reversed) | Load-Bearing Decisions / Identity; scenario 4; Risks |
| D-10 | B(C-1) | **fixed** | Normative block 4 (0009 collision named); new A11 |
| D-11 | B(C-7) | **fixed** | Normative block 1 (unreachability claim withdrawn) |
| D-12 | B(C-6) | **fixed** | Normative block 3 (category vs rule-id keying) |
| D-13 | B(C-16) | **fixed** | A1 scope limit; new A12 |
| D-14 | B(C-11) | **fixed** | A5 re-derived scope note |
| D-15 | B(C-2) | **dismissed-with-cite** (re-raise of cove F-1) | — |
| D-16 | A(C-8) | **fixed** | Normative Contracts (closure argument added) |
| D-17 | A(C-11) | **fixed** | QOC blast-radius row; deciding-rows para; Consequences |
| D-18 | A(C-6) | **partly dismissed-with-cite**, partly fixed | A8 (evidence corrected) |
| D-19 | A(C-4) | **fixed** | Normative block 2 (lower bound stated as absent) |
| D-20 | A(C-12) | **fixed** | A5 (same repair as D-4/D-14) |
| D-21 | A(C-15) | **fixed** | Approach 5; block 5; Capability Dependencies; Phase 2; scenario 9 |
| D-22 | B(C-15) | **partly dismissed-with-cite**, partly fixed | Problem Statement (XState analogy conceded) |

## The findings that changed the design

**D-1/D-2 — A4 was `Verified` against the wrong proposition.** The reserved-key
rule has two clauses with different blast radii. A4 swept only the second (no
owned/observed tag named `recognized`) and reported "invalidates nothing at
HEAD" for the whole rule. The first clause — a recognized-provenance
declaration MUST be named `recognized` — is violated by **every** committed
recognized-provenance declaration in the repo, and A4 quoted all three as
evidence that nothing is violated. Grounded independently:
`0002…/evidence/spikes/rdr-fixture.toml:25` and `kata-fixture.toml:22`
(`[tags.outcome]`), `0003…/evidence/spikes/guard-fixture.toml:36`
(`[tags.rewind_target]`). RDR 0002 (`Final`) calls these "the canonical examples
implementation tests must promote" (`0002…md:361`). Both passes reached this
independently, which is the strongest signal the fallback produced.

A4 is rewritten to sweep both clauses, carries the migration inventory (three
declarations, five `[rule.match.outcome]` reference sites), and the Risks
section gains the day-one-breakage risk that the old "no migration needed"
mitigation had concealed.

**D-1 surfaced a genuine unknown (A9).** In the canonical fixture the
recognized-provenance tag is *also* the outcome-gating predicate:
`[tags.outcome]` is matched by `[rule.match.outcome] eq = "round-clean"` whose
values are exactly the root `outcomes` alphabet (`rdr-fixture.toml:1,25,59`).
Whether that normalizes into `Row.Match` or `Row.Outcome` is **not stated in RDR
0002 and not derivable from the kernel** — the rename is mechanical under the
first reading and a semantics change on a Final peer's evidence under the
second. I did not collapse this fork: the evidence is genuinely absent rather
than merely hard to read, which is the tiebreaker gate's own exception. Booked
`Pending` and it gates Phase 2's rename.

**D-8 — block 5 was an unenforceable MUST.** Both passes converged. The block
had no authored source form, no lint half, and no producing-side test; scenario
9 tested only the kernel refusal shape this RDR does not change. Rather than
delete the block or leave it decorative, the fix gives it the enforcement point
already being built: the exported `Input` predicate now also rejects a
`Row.RequiresOwned` entry naming `recognized`. Zero new surface — it is the same
predicate block 4 requires — and it converts an unfalsifiable clause into one
that fails the moment any derivation path emits the reserved name.

**D-10 — the peer collision pass A missed.** RDR 0009 writes a structurally
identical clause over the *same* entry point: "the conformance predicate MUST be
exported by the kernel package … and Resolve's entry precondition MUST be that
same function — one predicate, two call sites" (`0009…md:510-517`), and pins its
own precedence at `:531-538`. This RDR borrows that shape while forbidding any
symbol-name binding, and neither RDR states composition or order. Block 4 now
names the collision and binds only what this RDR owns (non-masking, disjoint
fields); the total order is booked as A11 for Stage 7.1, since neither RDR may
take it unilaterally.

**D-9 — a deferral reversed rather than carried.** "Whether lint additionally
warns on case- or whitespace-variants is a Resolve question" was wrong on its
face: Resolve had already run, and the byte-exact rule is what *creates* the
hazard. `[tags." recognized"]` is invisible in a diff, lints clean, and never
fires — the Problem Statement's failure reproduced by the fix. Settled here as a
non-blocking advisory (advisory because the name is legal and this RDR must not
reject a tag it does not own), with scenario 4 pinning both dispositions and the
advisory, including a negative case so it is targeted rather than blanket.

## The two dismissals

**D-7 — "told why" is asserted at the data level; no verb renders it.** Both
passes raised it; it is already `3amigo/charted.md` C-2, charted as net-new
scope belonging to RDR 0005's extensible code-string table and carried as a
cluster-reconcile watch item. Grounding source 3 (the RDR's own decided text):
Consequences already states the limit explicitly — comprehension "is therefore
asserted at the data level only… Whether that suffices is a question for
whichever CLI surface renders it." Re-absorbing it here would author a contract
this RDR does not own (the scope-expansion wormhole). The standing disposition
governs; no edit.

**D-15 — a Final peer (0007) rests on this Draft's sentence.** Identical to cove
F-1, dismissed-with-cite last pass with a partial fix already landed on the
joint-check line (the inbound reference disclosed, A13's negative-existential
carried as a cluster-reconcile watch item). Pass B rediscovered it
independently, which is expected of an independent pass and is not new
information. RDR 0007 A13's verdict is *accepted exposure* built from a
five-peer absence sweep; it is unchanged or reinforced however this RDR's
enforcement locus settles. Standing dismissal governs; no edit.

## Partial dismissals

**D-18 — pass A's symptom overstated.** A claimed the first non-nil error path
is "a behavioral break for callers that ignore the error — which is every
caller." Refuted on the symptom: there are **zero non-test callers** of `Resolve`
at HEAD (only the definition at `resolve.go:318`), and every test call site
routes through `resolve_test.go::mustResolve`, which `t.Fatalf`s on a non-nil
error. A tripped predicate is a red test, not a silent nil-deref, and pass A's
premortem "nil-pointer panic in a CLI verb" cannot happen today. The underlying
contract observation is real for future CLI wiring, so A8's evidence is
corrected rather than the finding dropped — including the
`resolve_test.go:748` empty-`Input{}` nil-error pin the predicate must keep
green.

**D-22 — half re-raise, half new.** The silent-reframing half is 3amigo PM-1,
already fixed (the product stance is explicit in the Problem Statement). New and
kept: the XState analogy compares a field in an engine-owned struct (`GuardArgs`,
disjoint from the author's `context`) to a bare word carved out of the author's
own flat `[tags.*]` namespace — which is precisely why this RDR needs a
validation category and XState needs none. Also new: all three observed authors
exercised the naming freedom being withdrawn. Both conceded in the Problem
Statement; the stance itself stands on the Decision Rationale.

## Needs (re)verification (Stage 6 closes these)

Four **new** `Pending` assumptions, all introduced by fixes in this pass
(flag-as-you-go):

- **A9 (new, `Pending`; widened at iteration 2)** — which normalized row field
  each tag-predicate reference to a recognized-provenance declaration produces:
  `Row.Match` for `[rule.match.<tag>]`, the guard field for
  `[rule.guard.all.<tag>]` / `[rule.guard.unless.<tag>]`, versus `Row.Outcome`.
  All three targets are live in the committed fixtures. Load-bearing on A4's
  migration claim, block 2's reference-position argument, and Phase 2's rename,
  which is explicitly gated on it. *Not collapsible here* — RDR 0002 never maps
  a predicate to a normalized field and the kernel cannot reveal it.
- **A10 (new, `Pending`)** — whether an accessor-produced owned tag key can
  spell `recognized` from artifact data. A8's literal sweep is structurally
  incapable of clearing this. Load-bearing because block 4 sends the breach down
  the path RDR 0001 reserves for *programmer* mistakes; if user data can reach
  it, that is the wrong channel and the block needs repair before lock.
- **A11 (new, `Pending`)** — composition and total order of this RDR's
  `Resolve`-entry precondition with RDR 0009's. Widened in scope by D-8's fix:
  the predicate now reads `in.Table.Rows`, so the independence half must be
  re-checked on *fields* rather than on inputs.
- **A12 (new, `Pending`)** — whether an RDR 0002 implementer reading 0002 alone
  encounters this constraint at all, given 0002's `[tags.<tag>]` grammar block
  admits any name and points nowhere.

Status changes to existing assumptions:

- **A4** — rewritten and re-`Verified`, now against both clauses. The claim text
  changed materially; the sweep was re-run, not re-stamped.
- **A1, A5, A8** — `Verified` stands; each gained an explicit **scope limit**
  naming what its evidence does *not* establish (A1: the category fence does not
  license the naming-grammar constraint; A5: the derivation governs declaration
  space, not `Input` data; A8: a literal sweep cannot clear a data-derived key).
  No previously-`Verified` assumption was invalidated, but three were narrowed.
- **A7** — `Pending`, unchanged. Its five settled-fact dependents (Approach 5,
  block 5, Capability Dependencies, Phase 2, scenario 9) are now written as
  conditional on it, per the Finalization Gate's Status-consistency rule.

Carried watch items (unchanged): RDR 0007 A13's negative-existential, and the
user-facing-surface question (`charted.md` C-2) — both cluster-reconcile.

## Charted to successor

None new. D-7's intent-channel and surfacing items were already charted by the
3amigo pass (`charted.md` C-1, C-2) and are re-cited rather than re-absorbed.

## Iteration 2 — delta-scoped re-run (convergence check)

The iteration-1 fixes were substantial (a rewritten A4, a widened predicate, four
new assumptions), so the lens was re-run **delta-scoped to the ten revised
passages** rather than re-critiqued in full — critiquing one's own edits is the
drift this guard exists to prevent. Verdicts: `iter-2/delta-verdicts.md`.

**8 CLOSED, 1 NEW-GAP, 1 OPEN.** Both survivors were real and are now fixed.

**NEW-GAP (A4's migration inventory) — fixed.** The both-clause sweep was
correct, but the inventory the fix introduced enumerated reference sites only for
the two `[tags.outcome]` declarations and missed the third violating
declaration's own reference:
`docs/rdr/0003-…/evidence/spikes/guard-fixture.toml:72`
(`[rule.guard.all.rewind_target]`). Grounded independently: `grep` over that
fixture returns exactly two `rewind_target` hits, `:36` (the declaration) and
`:72`. RDR 0002 puts `[rule.guard.all.<tag>]` in the same tag-predicate family as
`[rule.match.<tag>]` (`0002…md:266-274`), so the inventory's own justifying
sentence covers it. Executing Phase 2 as written would have renamed the
declaration and left `:72` pointing at an undeclared tag — an `unknown tag`
regression on RDR 0003's canonical fixture, introduced *by* the migration.

Fixed: the inventory is now six sites across three fixtures, the count is
corrected in all four propagations (blast-radius row, Consequences, Risks,
Phase 2), and **A9 was widened** — the guard reference normalizes into the row's
guard field, a third target the assumption did not contemplate and the one whose
predicates RDR 0007 owns.

**OPEN (block 4's non-masking MUST) — fixed.** The collision paragraph required
that this RDR's breach "MUST NOT mask or be masked by a table-shape breach" while
also forbidding the implementer to fix an order and never licensing an aggregate
error. With `Resolve` returning a single `error` (verified: `resolve.go:318`),
non-masking *is* an ordering statement, so the two MUSTs could not both be
discharged while A11 is `Pending` — which would have blocked scenario 6 despite
Testing Strategy placing it in the runnable-at-HEAD half.

Fixed by separating detection from reporting: this RDR now requires only that its
breach be **detected whenever present, never skipped because another precondition
fired** — which is what it actually owns — and explicitly licenses either error
for a doubly-breaching input until Stage 7.1 narrows it. A doubly-breaching caller
is a producer with two programmer mistakes, not a case this RDR owes exact error
text on. A11's "If wrong" was re-scoped to match: report-order variation is no
longer a falsification, only genuine verdict suppression is.

**Convergence: reached at iteration 2.** No open ledger entry remains — every
D-entry is fixed, dismissed-with-cite, or charted, and both iteration-2 survivors
are repaired. The cap (3) was not hit and no entry churned across iterations.

## Tiebreakers escalated to the driver

None. The one fork that looked genuine — block 4's unconditional clause (D-5) —
collapsed under the tiebreaker-reduction gate with a reason rather than by
symmetry: a reservation that lapses per-call is not a reservation, since a
producer would have to know whether an outcome is in flight to know whether its
own key is legal, which is the cross-side coupling this RDR exists to remove.

**A9 is not an escalated tiebreaker** — it is an absent fact booked for
verification, which is the flag-as-you-go path rather than the tiebreaker path.
It does gate Phase 2's rename, so Stage 6 must close it before implementation.
