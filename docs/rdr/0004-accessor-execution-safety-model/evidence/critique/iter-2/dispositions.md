Model: claude-opus-5[1m]

# Critique iter-2 — Dispositions

Dual-model pass. Origin ledgers: `critique.md` (Pass A, claude-sonnet-5, C-1…C-9)
and `critique-modelB.md` (Pass B, claude-fable-5, C-1…C-8), reconciled by passage
anchor in `diff.md` — the `C-N` ids do not correspond across files. Grounding
verdicts in `grounding.md` (G-1…G-5) are the authority for the gate.

Delta-scoped per re-entry: the RDR's Status names A8 as the re-verify target, so
rows landing on the revised read-completeness surface were dispositioned; rows
re-raising a concern the iter-2 3amigo pass already closed (T-1…T-14) were
checked against `../../3amigo/iter-2/dispositions.md` rather than re-edited.

## Fixed

- **fixed** — origin: Pass B C-1 (seam clause + Disposition Table "loud one layer
  down") — sections touched: Disposition Table absent-key row and its summary
  paragraph, Capability Dependencies, Failure Modes. Grounded G-4: the guarantee
  fires only where a surviving row declares the key in `Row.RequiresOwned`;
  `TagSet.matches` is provenance-blind and drops the row before `missingOwned`
  runs, and `escapeOrRefuse` can rescue the resulting `no_match` into a plan with
  `Escaped = true`. The table's flat "No input class exits silently" was false as
  written; the row itself already carried the conditional and the summary dropped
  it. Scoped as a stated dependency, not a fix here — see *Charted* below.

- **fixed** — origin: Pass A C-2 + Pass B C-2 (AGREED #1; absence encoding
  unpinned) — sections touched: A8 Evidence, Load-Bearing "Absence crosses the
  seam as omission". Grounded G-1/G-2. Two corrections: the distinction dies at
  `internal/resolve/resolve.go::Input.Owned` (`[]Tag`, two fields, no third
  state) *before* the kernel runs — both sites said "the kernel cannot recover
  it", which is true in effect but mislocates the constraint; and the sentinel
  prohibition is not compiler-enforceable, since no accessor package exists under
  `internal/` and `Tag{Key:…, Value:"<absent>"}` typechecks. Stated explicitly so
  the rule reads as load-bearing rather than decorative.

- **fixed** — origin: Pass B C-4 (post-mutation `read_back_incomplete` leaves the
  artifact advanced with no modeled edge) — sections touched: Normative Contracts
  (new clause), Risks and Mitigations, Disposition Table (new row), Desk Trace
  step 5, Oracle Discriminability (new row 8), MVV, Testing Strategy (new
  Scenario 8), Critical Assumptions (new **A10**, Pending). Genuinely uncovered:
  no retry/idempotency/recovery treatment existed anywhere in the draft. Resolved
  *without* re-opening the declined undo semantics — stating what the caller may
  assume is the missing half of an already-owned refusal class, not a
  transaction. Tiebreaker collapsed on that distinction; no §strong-consult
  needed.

- **fixed** — origin: Pass B C-5 (refusal vocabulary disagrees across three
  sections) — sections touched: Failure Modes. Grounded G-5: `RefusalKind` is a
  closed five-kind set owned by RDR 0001 with no artifact-unavailable member, and
  `internal/cli/clierr` has no such code (no named code constants at all). The
  Disposition Table already routed the class to `execution_failure`; Failure
  Modes listed "artifact unavailable" as a peer refusal. Removed it and stated
  the disjointness from the kernel's set.

- **fixed** — origin: Pass A C-7 (silent-shape enumeration is reactive, not
  derived) — sections touched: Failure Modes. The count had widened 2→3 across
  passes, which is the "question the unit" signal. Replaced the bare count with
  the derivation rule (a success-shaped return over state the accessor did not
  establish) so the list is closed by construction rather than by collection, and
  recorded the one shape that is guarded but not fully closed.

- **fixed** — origin: Pass A C-1 + Pass B C-8 (AGREED #3; A9 bundles four
  independent rules under one Pending) — sections touched: A9 (bundling note),
  Assumption Verification, Prerequisites. Rules retire per-rule rather than
  flipping A9 as a unit; A9 splits if any single rule blocks. The new A10 was
  added as a *separate* assumption rather than a fifth A9 clause, so the bundling
  does not worsen.

## Dismissed with cite

- **dismissed-with-cite** — origin: Pass A C-9 (MVV Scenario 6 has no negative
  control) — sections touched: none. Refuted by the draft: Oracle
  Discriminability row 6 already carries it verbatim ("A key the artifact *does*
  carry must resolve normally through the same path — an implementation that
  omits every key passes the refusal assertion vacuously"). Pass A read Scenario
  6 without its paired Oracle row.

- **dismissed-with-cite** — origin: Pass A C-6 (A4 replay rides the disowned
  fallback; "the only exercised read path, 10 of 13 lines") — sections touched:
  none. Refuted on its stated basis by G-3: `expectedTagKeys` fires only when
  `len(requested) == 0`, and all three read scenarios pass explicit keys
  (`"status", "profile"`), so the derived path is not exercised by them. The
  surviving thin kernel — that the RDR cites the fallback as "fixture
  convenience, not the contract" — is already stated in Load-Bearing "Requested
  key set", and witness provenance was closed as 3amigo T-13.

- **dismissed-with-cite** — origin: Pass A C-3 (timeout precedence presumes an
  executor observing partial progress) — sections touched: none. Re-raise of
  3amigo **T-8**, whose fix is the precedence clause and Oracle row 7 now in the
  draft. Pass B C-3 raises the adjacent determinism concern; both are the same
  already-adjudicated call, and A9 books the implementability question.

- **dismissed-with-cite** — origin: Pass A C-4 (whole-read refusal poisons
  multi-key accessors) — sections touched: none. Re-raise of 3amigo **T-9**,
  closed as a decision in Load-Bearing "Read refusal granularity" — the
  conservative choice matching JDR 0001's refuse-rather-than-guess rule. Not
  re-litigated.

- **dismissed-with-cite** — origin: Pass A C-5 / Pass B C-6 (AGREED #2;
  absent-vs-unreadable delegated to bindings with no enforcement) — sections
  touched: none. Re-raise of 3amigo **T-5**, closed by the Load-Bearing "Absent
  vs unreadable is the binding's call, defaulting to unreadable" bullet, which
  already concedes the spike proves branch machinery rather than a binding's
  classification. Pass B's sharper form (the conservative default incentivizes
  guessing below the seam) is a real observation but argues against a decided
  call; A9's If-wrong already names the fallback.

- **dismissed-with-cite** — origin: Pass A C-8 (Proportionality's "`Input`
  unchanged" is irrelevant/reads as license) — sections touched: none as argued.
  Re-raise of 3amigo **T-11**. The sentence is *true* — G-2 confirms
  `Input.Owned` is unchanged — and the substantive gap Pass A gestures at (why a
  sentinel still typechecks) is now stated directly at the Load-Bearing bullet,
  which is the load-bearing site.

- **dismissed-with-cite** — origin: Pass B C-3 (timeout entangles a
  wall-clock-dependent class with deterministic ones; A4 replay narrowed) —
  sections touched: none. The precedence rule is the decided call (T-8), and the
  flakiness concern is an implementation-fixture matter the MVV owns; A4's replay
  claim is scoped to disposition equality over *fixture* results, not over
  wall-clock races. Booked under A9 rather than re-opened.

## Charted to successor

- **charted-to-successor** — origin: Pass B C-1 residual / Pass B C-7 (nothing
  ties the declared requested-key set, or `Row.Match` / guard keys, to
  `Row.RequiresOwned`; no layer is obliged to populate it). Out of scope: RDR
  0002 owns the normalized row, and JDR 0001 **§JD-3** already records this as an
  open joint decision ("appears **zero** times in 0002 … no layer is obliged to
  populate it"), with 0007 deriving it from `Writes` on escape rows while 0009
  empties `Writes`. Absorbing it here would be the scope-expansion wormhole.
  Recorded in this RDR as an **Open** row in Capability Dependencies, and see
  `Charted:` below.

`Charted:` `Row.RequiresOwned` has no obliged producer, and no layer cross-checks
that keys used in `Row.Match` or a guard string are declared in it — so RDR
0004's absence-as-omission seam is necessary but not sufficient for
`owned_state_unavailable` to fire. Out of scope for 0004 (0002 owns the
normalized row). Suggested successor: close JDR 0001 §JD-3 in the 0002/0007/0009
cluster, or seed an RDR for match/guard→`RequiresOwned` consistency validation.
Flagged for reconcile and for `/rdr-cluster-reconcile`.

## Needs verification (Stage 6 closes these — not verified here)

- **A10 is newly Pending** (Method: MVV Test). Introduced by this pass's
  post-mutation clause. Verify by the MVV write case whose re-read fails on a
  compared key: assert the refusal carries the applied-but-unverified sense and
  that exactly one write invocation occurred with no compensating action. The
  Resolve spike cannot witness it — its `write` re-reads by cloning the tag map,
  so the re-read cannot fail independently of the write.
- **A9 unchanged in status**, but its four rules are now explicitly independent
  and retire per-rule; a single-rule refutation must not flip all four.
- **A8 stays Verified.** The corrected Evidence cites `Input.Owned` alongside
  `missingOwned` — a more precise anchor for the same claim, not a new one; the
  spike transcript lines it rests on are untouched.
- **Cross-RDR (for reconcile)**: the Capability Dependencies row naming JDR 0001
  §JD-3 is an open dependency on 0002/0007/0009, not a 0004 obligation. Route at
  reconcile / cluster-reconcile.

## Amendment sweep

New/amended clauses swept per rdr-common §amendment-sweep. `read_back_incomplete`
re-read at all 10 sites; `A9`→`A9 and A10` at Prerequisites and Assumption
Verification; the "kernel cannot recover" phrasing corrected at both sites that
carried it (A8 Evidence and Assumption Verification prose); "artifact unavailable"
removed from Failure Modes and checked against the Disposition Table row that
already routed it to `execution_failure`; the bare silent-shape count replaced by
a derivation rule so no count needs future widening. New Disposition Table row
and Oracle row 8 added as siblings to the existing write rows. Mechanical grep for
template brackets and `_Draft placeholder._`: none.

Cross-lens check: the iter-2 grounding and 3amigo passes edited A8's Evidence
citation form and the Load-Bearing absence bullet. This pass edits the same two
fields — deliberately, to correct the layer the constraint binds at — and the
edits are additive to those passes' fixes (omission-not-sentinel is preserved
verbatim; only the *why* is re-anchored). No two-lenses-pinned-one-field
conflict.

## Tiebreakers

None escalated. The one candidate was Pass B C-4: whether stating post-mutation
caller expectations re-opens the declined undo/transaction semantics. It
collapsed on the evidence — Alternative 4 and Day 2 Operations decline *undo*,
while the new clause only names what an already-owned refusal class means and
explicitly forbids compensation. Declining to say what a refusal means is not the
same decision as declining to undo. No §strong-consult was needed.

## Mini-check cue read

Not owed by this pass. The iter-2 grounding pass ran the cue read and its four
fired tables (Disposition Table, Oracle Discriminability, Fidelity Table, Desk
Trace) are in the draft; tables persist across passes. This pass added rows to
three of them rather than firing a new cue. Source-authority census remains `no`:
the executor is still the sole writer of the read result, and the post-mutation
clause names reporting semantics rather than introducing a competing source of
truth.
