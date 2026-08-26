# Critique dual-model diff — RDR cli/0011

Pass A: `critique.md` (Model: claude-opus-5, 18 rows)
Pass B: `critique-modelB.md` (Model: claude-sonnet-5, 10 rows)

Reconciled by PASSAGE ANCHOR — `C-N` ids do not correspond across files.
Barrier: both passes landed before this diff was written; this diffing context
authored neither (rdr-common §auto-fanout).

## Agreement (both models, same passage) — strongest signal

| Passage | A | B | Reconciled claim |
|---|---|---|---|
| `0011:A8` / C1 carrier | A-7, A-18 | B-1 | The `Refusal.Undecided` carrier does not work as A8 describes. Both reached it independently, by different routes. |
| `0011:C1` compound clause | A-10 | B-3 | A conflicted key reported as `absent` names a remedy that cannot work. |
| `0011:§consequences` dual migration | A-15 (adjacent) | B-4 | Predicate flip + payload rename ship together; `--all` undoes only the first. |
| `0011:A7`/`A8` Pending vs stated-as-fact | A-4 | B-8 | Largest behaviour change rests on two Pending assumptions whose census already grew once. |
| `0011:§problem-statement` success criterion | A-15 | B-6 | Success is defined so the residual wall-of-candidates case cannot count as failure. |
| `0011:A5`/`A6`/`S1` flagship narrowing | A-14, A-15 | B-5 | The 21→3 result holds only because `stage` is owned+reader-served; unsupplied match keys do not narrow. |
| `0011:C3` conflicted fixture | A-11 | B-7 | The mandated fixture is not shown being built; A proves it is *unbuildable*. |
| `0011:D-naming` `--all` boolean | — | B-10 | B only. Weak — the RDR already states the `gh` migration path. |
| `0011:§decision-rationale` caller audit | — | B-9 | B only. Partly answered by the RDR's own three-call-site census. |

## A-only findings (B missed) — all confirmed against `main` by the diffing context

These are the pass's whole value; B's ledger stops at the design surface, A's
reaches the shipped code.

- **A-1 / A-16 — the resolution-level veto.** `0007:C10` is cited **0×** in the RDR
  (`0007:C11` 5×). `gate` is table-scoped: `resolve.go:558` `if len(undecidable) > 0
  { return nil, ... }` discards `selected`. `0007:C10` states it normatively:
  "if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST
  refuse `guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be selected".
  **Verified blast radius**: all 22 rules in `models/rdr.toml` carry a
  `[rule.match.stage]` block, so an unbound `stage` makes every row indeterminate →
  whole-table refusal. C1's "SURVIVES into `gate` as a non-selectable row" has no
  mechanism and contradicts the locked peer contract.
- **A-2 — no partition for "survives but non-selectable".** `gate`'s partition is
  exactly `{GuardFalse→pruned, GuardTrue→selected, GuardUnevaluable→veto}`. The RDR
  names no third channel. Every available invention triggers A-1.
- **A-3 — JC1 asserts the negation.** `resolve.go:478-481` returns on `blocked != nil`
  **above** the `switch len(selected)` that reaches `escapeOrRefuse`. A
  `guard_unevaluable` refusal is unreachable by an escape row, so 0010's "otherwise"
  catch-all stops rescuing. JC1 calls the coupling "strictly widening" and
  "cannot turn a 0010 acceptance into a refusal".
- **A-6 — S5 re-expects the positive control.** `TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues`
  is by its own docstring the control for `...ConflictedKeyIsNotEscapableByNamingIt`.
  Re-expecting it deletes the only oracle proving a clean escape still rescues.
- **A-7 — `isGuardBlock` is not a payload filter.** Its `continue` gates three things
  (seam call, payload append, Kleene operand). Its own comment forbids the edit A8
  proposes, naming `BlockMatch` as the reachable vector. Two unnamed files go red:
  `guard_fixup_0007_test.go:115` ("a match atom is never handed to the seam") and
  `guard_adversarial_0007_test.go:180` (`"§D12 match block"` arm).
- **A-8 — the `unresolved`→`unknown` rename is not mechanical.** All six reads go
  through `stringsAt`, which returns `nil,false` on objects; two sites discard the
  ok bool (`flow_adversarial_0005_test.go:446`, `flow_mvv_0005_test.go:66`) and the
  adversarial one is a NEGATIVE assertion that passes **vacuously** on the empty
  slice. C3 pre-declares a green suite as the pass criterion for this class.
- **A-9 — "no new field on `Plan`" needs two unnamed signature changes.**
  `summarize(row, view, gatesRan)` receives no kernel output; `excluded(...) bool`
  discards the `Result`. Both must change. The Infrastructure Audit marks them
  "Extend" and names neither.
- **A-11 / A-17 — A3 is REFUTED, and it is stamped Verified.**
  `internal/table/load.go::checkAccessorBindings` refuses at load: `owned tag %s is
  served by %d readers; want exactly one`. The two-reader route cannot load; the
  `--tag` route is refused by `parseTags` (the RDR says so itself). `merge` conflicts
  only within one provenance. So `tv.conflicted` is **unreachable from the CLI**.
  A3 is the sole ground for kernel placement (Approach, `D-placement`, ALT3
  rejection, §Decision Rationale) — the whole `large` Profile is bought to
  distinguish a case no user can reach.
- **A-13 — `--all` negative pinned only by non-registration**, asserted against
  cobra's third-party error class.
- **A-5 / A-12 — the oracle count drifts three ways.** S5 + `oracle` mini-check say
  five; MVV 8, trace step 8, illustrative-code, A7-If-wrong say four; Consequences
  says three. (Independently found by the dispatcher before either pass returned.)

## B-only findings

- **B-2 — `unknown` folds four causes into three reason tokens.** Real but weaker than
  A's framing: the RDR's `D-undecided-vocabulary` argues the three-token set is
  deliberate. Dispositioned as re-raise unless it survives the decided-text gate.

## Divergence assessment

No contradiction between the passes — B is a strict subset of A on every shared
passage, at lower resolution. The disagreement is DEPTH, not direction: B accepted
the RDR's own account of the mechanism (A8's `isGuardBlock` story, the carrier, A3's
reachability) and critiqued the design on top of it; A opened the files and found
the account itself false. That is the dual-model draw paying off in the direction it
is supposed to.

**Unhealthy-pass check**: neither pass is generic; both anchor to real passages and
named symbols. No rerun warranted.

**Verdict**: the lens found structural defects, not polish. Four findings
(A-1/A-16, A-3, A-11/A-17, A-18/B-1) each independently invalidate a load-bearing
claim, and two of them refute assumptions currently stamped **Verified** or route
through a locked peer contract the RDR never cites. This exceeds what the resolve
half may absorb by editing.
