Model: claude-fable-5

# Critique dual-model diff + dispositions — pass B vs pass A

Pass A: `critique.md` (claude-opus-5[1m], C-1..C-15, resolved in
`dispositions.md` iteration 1). Pass B: `critique-modelB.md`
(claude-fable-5, fresh context, B-1..B-15). Diffed by passage anchor —
IDs do not correspond across files. This file completes the dual-model
obligation `dispositions.md` recorded as unsatisfied.

## Anchor diff

Convergent (model B independently re-derived a pass-A finding whose
disposition already stands — strong signal the fix or dismissal was
right; several B-rows quote the post-fix text itself):

| B row | A row | Standing disposition |
|---|---|---|
| B-1 | C-5 | fixed (Risks residual → UNMITIGATED; enforcement point = 0003's implement stage) |
| B-2 | C-4 | fixed (Consequences status-vs-defect distinction) |
| B-3 | C-2 / A10 | dismissed-with-cite (A10 Pending blocks lock; Phase 2 encoding declared vector-local/non-normative) |
| B-8 | C-7 | dismissed-with-cite (asymmetry stated in Failure Modes, routed Phase 3) |
| B-9 | C-6 / A13 | fixed (A13 Pending; Failure Modes *Refusal defeated by caller-supplied state*) |
| B-12 | C-13 | fixed (ordering pinned to behavior, rationale retired, two-round-trip cost stated) |
| B-13 | C-3 + C-10/C-11 | dismissed re-raise + fixed (matrix scores Operability; SCXML scope statement; unknown veto frequency recorded as accepted risk) |
| B-14 | C-9 / A9 | fixed (strike-out condition moved out of the normative block; A9 blocks lock) |
| B-15 | C-14 | charted-to-successor (`Charted.md` → finalize Proportionality) |

Divergent / net-new (model B only) — dispositions this session:

| B row | Disposition | Section touched |
|---|---|---|
| B-5 | **fixed** — CONFIRMED net-new: the blessed two-row pattern self-vetoes. With X absent, a bare `X eq v` value row is GuardUnevaluable, survives (SURVIVOR MEMBERSHIP), and the aggregation veto refuses instead of selecting the absence row — contradicting Scenario 8's "exactly one row selected". Grounded against the aggregation clause + existence clause + `guard-fixture.toml` atom shape. Fix: value row must be existence-guarded; A5 claim rewritten; **A15 added (Pending)** for the guarded shape's grammar legality; Scenario 8 rewritten with a negative control; wording of the strong-Kleene clause, the Approach summary, and the D8 ratification corrected ("present tags alone decide" → "decided atoms alone determine": an existence atom legitimately decides F from absence, which the old wording forbade) | A5, A15 (new), existence-clause prose, strong-Kleene block, Approach, D8 block, Scenario 8, A12 note, Risks sentinel mitigation, Phase 3 |
| B-4 | **fixed** — the mapping-failure "surface outside the verdict channel" MUST named no channel that exists at the seam (`Evaluate` has no error return; `Resolve`'s error return is unreachable from inside the seam — grounded `resolve.go:91-94`, `:318`). Clause now names the panic concretely and forbids the kernel recovering it into a verdict | Normative Contracts, SCOPE block |
| B-6 | **fixed** — the existence clause asserted A14-gated expressibility (`all … exists = false`) as settled inside a normative block, violating the status-consistency rule. Expressibility + disjunction sentences moved to an A14-gated prose note (C-9 precedent) | Existence clause + new prose note |
| B-7 | **split** — the closed-world match behavior is dismissed-with-cite (Background's exclusion clause; match semantics owned by RDR 0001 — a decided call, re-raise). The actionable residue is fixed: Phase 3 authoring guidance gains the predicate-placement rule (state that must not mask belongs in guard atoms, never the match pattern) | Phase 3 |
| B-10 | **fixed** — A6b's routing named an unamendable Final peer as destination; Prerequisites line now requires a destination that exists (0004's implement stage or successor RDR), mirroring the A10 handoff fix | Prerequisites |
| B-11 | **dismissed-with-cite** — claims the unevaluable-escape behavior is "documented only inside an assumption's evidence note"; wrong on the RDR's own text: the aggregation normative block states it ("An unevaluable ESCAPE row yields guard_unevaluable in place of the candidate-set refusal"), it is frozen by TestAdv2/TestFixup1d, and the A3 correction is recorded openly. Re-litigates a decided + frozen call |

## Review-gate note

Pass B healthy: concrete passages, named functions (`TagSet.matches`,
`escapeOrRefuse`, `evaluateGuard`), origins span §1/§2/§3/premortem/AT.
Ten of fifteen rows converge with pass A — high cross-model agreement on
the defect map — and the five net-new rows all grounded (four fixed, one
dismissed on the RDR's own text). B-5 is the pass's load-bearing catch:
an internal contradiction both the opus pass and every prior lens missed.

## Needs (re)verification — carried to Stage 6

- **A15 (new, Pending)** — one-row existence-guarded value atom is
  grammar-legal under RDR 0003 (two operator keys on one tag-table, or
  same tag in `all` and `unless`). Method: Source Search vs 0003
  grammar / 0002 normalization. Blocks lock alongside A12; gates
  Scenario 8's primary leg.
- **A12 (Pending, unchanged claim)** — evidence note extended: the
  overlap proof must also accommodate the value row's new existence atom.
- **A5 (Pending, unchanged status)** — pattern reshaped
  (existence-guarded value row); still gated on A12, now also A15.
- **A2 (stays Verified)** — derivation untouched; evidence prose
  clarified that existence atoms enter the tables as decided verdicts
  (sanctioned absence-decided F). No load-bearing change.
- **A9 / A10 / A13 / A14 / A6b** — unchanged by this pass.

## Tiebreakers

None escalated. B-5's fork (change the veto vs change the pattern)
collapsed on evidence: narrowing the veto was already rejected in
*Briefly Rejected* as absence-as-false at resolution scope, so the
pattern side had to move; the existence-guarded value row achieves the
absent-leg selection with no unevaluable survivor, at the cost of one
new Pending grammar assumption (A15) — strictly cheaper than reopening
the adjudicated veto.

## Convergence

All 15 B-rows dispositioned (9 convergent with standing pass-A
dispositions, 4 fixed, 1 dismissed, 1 split fixed/dismissed). No open
ledger entries; edits were targeted, not a rewrite — no delta re-run
owed. The foundational dual-model obligation for the critique lens is
now satisfied: two base models, ledgers diffed by anchor, divergences
resolved.
