Model: claude-opus-5[1m]

# 3amigo consolidation — RDR 0007 Guard Predicate Totality

Iteration 1. Three personas run independently (PM / Implementer / QA), no
findings reused between them. Consolidation names the passages that drew fire
from **two or more** personas — those are the priority rewrites.

Persona lists: `persona-1-pm.md`, `persona-2-implementer.md`,
`persona-3-qa.md`.

## Hotspot passages (2+ personas)

### H-1 — The guard-text representation gap (IMP-1 + QA-1 + IMP-5)

**Passage**: Normative Contracts, every clause quantifying over atoms —
"Value-comparing guard operators (equality, membership, bounded integer
comparison, set containment) are PARTIAL over the assembled evaluation view: an
atom whose referenced tag key is absent…"; and Technical Design, "The contract
constrains one seam — `GuardEvaluator.Evaluate`".

The whole normative core is written over **atoms, operators, `all`/`unless`
blocks, and referenced tag keys**. The seam it constrains is
`Evaluate(guard string, view TagSet)` over `Row.Guard string` — an opaque
string. Nothing in this RDR, RDR 0001, 0002, or 0003 states whether that string
is a guard *name* (which is how `fixtures_test.go::fixtureGuards` treats it, via
`decided map[string]bool` keyed on guard text) or a serialized predicate the
evaluator parses. RDR 0003's own spike fixture authors guards *structurally*
(`[rule.guard.all.profile]`, `[rule.guard.unless.prelock_iterations]`), so the
structure exists upstream and is flattened to a string somewhere unnamed.

Consequences across personas: the implementer cannot tell what `Evaluate`
receives (IMP-1); QA cannot construct fixtures for scenarios 4/6/7/8 because no
Go representation of an atom exists anywhere and there is no `testdata/`
repo-wide (QA-1); and the empty-`unless` identity has no landing site on `Row`,
which has no `all`/`unless` fields (IMP-5).

**Priority**: highest. Every other atom-level clause is unimplementable and
untestable until this is pinned, and it is the one gap that blocks both the
Phase 2 vectors and the RDR 0003 evaluator.

### H-2 — `guard_unevaluable` cannot name the absent tag, or discriminate which row it came from (PM-1 + IMP-3 + QA-2 + QA-4)

**Passage**: Problem Statement — "told plainly that the artifact state needed to
decide was missing"; Failure Modes — "`GuardResult` carries no missing-tag
payload, so 'which tag was absent' is inferable only from the guard text";
Testing Strategy scenarios 1 and 5.

All three personas hit the same observable from different sides:

- **PM-1**: the promised user outcome is a refusal naming the missing state, but
  the shipped refusal carries only guard text. The sibling
  `owned_state_unavailable` names its keys via `MissingOwned`; this one does not.
- **IMP-3**: `Refusal.Guard` is single-valued and `gate` picks it by
  `slices.MinFunc` over `compareRefs(refOf(a), refOf(b))` — a `(RuleID,
  SourceLocator)` tiebreak the RDR never states normatively.
- **QA-2**: scenario 1 and the MVV name the expected observable as
  "`guard_unevaluable` with `Escaped:false`", but `Escaped` is a field of `Plan`
  (`resolve.go` L250, set in `planOf` L509) and `Result.Plan` is nil on every
  refusal — so that criterion cannot be asserted as written.
- **QA-4**: scenario 5's two legs both expect `guard_unevaluable`, discriminated
  only by an unspecified `Refusal.Guard` string; leg (b) has no observable at all.

**Priority**: high. This is the RDR's own headline user outcome and its MVV
oracle, and the defect is partly factual (QA-2), not merely under-specified.

### H-3 — Aggregation: an unevaluable candidate vetoes a decided-true sibling (PM-4 + QA-4 partially)

**Passage**: Normative Contracts — "if any surviving candidate row's guard is
GuardUnevaluable, the resolution MUST refuse guard_unevaluable — a
decided-GuardTrue sibling row MUST NOT be selected while an unevaluable
candidate exists."

PM-4 reads this as an unexamined user-facing cost: one unreadable row disables
an otherwise-working table, and the RDR ratifies it on the grounds that the
kernel already behaves this way (A3, spike Probe B) rather than weighing it. The
tension PM-4 names with D5's anti-poisoning rationale is worth an explicit
sentence in the draft even if the reading is ultimately wrong (D5 scopes to
*matching* candidates; an unevaluable matching row is inside that set by
construction).

**Priority**: medium-high — likely resolvable by making the already-made
decision explicit rather than by changing it.

### H-4 — Migration for tables that relied on closed-world reads (PM-2 + PM-3)

**Passage**: Trade-offs / Consequences — "tables whose authors expected
closed-world reads ('x == 1 just fails when x is missing') now refuse where they
previously escaped or pruned; those guards must be rewritten with explicit
existence atoms."

A behavior flip for existing tables is disposed of in one bullet: no inventory,
no lint or detection step before production, and no migration phase in the
three-phase Implementation Plan. PM-3 stacks the related cost — non-escapable by
design plus open A6b (no retry guidance) means a degraded artifact has no
modeled recovery path — and notes the scoring matrix carries no
author-expressiveness or operator-recovery criterion.

**Priority**: medium. Single-persona in origin (PM) but two findings converge on
the same passage.

## Single-persona findings (resolved individually, not hotspots)

- IMP-2 — "present in the assembled view" is provenance-blind while the kernel's
  only presence predicate is `has(key, ProvenanceOwned)`; RDR 0003's fixture
  puts `exists` on an `observed` tag (`cluster_eligible`). Grounded and real.
- IMP-4 — the exported conformance harness has no package/signature, and
  `GuardEvaluator` lives under `internal/`.
- QA-3 — scenario 2 cites `TestAdv1b` as freezing D8 preservation, but that test
  decides on guard text with the `TagSet` discarded.
- QA-5 / IMP-5 — A9 is `Pending` while the normative block states the
  empty-block identity as settled.
- PM-5 — Consequences claims masking is closed in the present tense; nothing
  this RDR ships changes user-visible behavior until RDR 0003 builds an
  evaluator.
