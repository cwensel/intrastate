# Verification — RDR 0001 Resolution Kernel

Stage-8 implementation verification evidence.

## Phase 3c — Fixup resolution notes

Both Phase 3 verifiers returned BLOCK on the same three defects. All three are
now closed. **No test was weakened, deleted, skipped, or edited** — the
adversarial tests were failing because the implementation was wrong, and the
implementation is what changed. `internal/resolve/resolve.go` is the only
source file touched; `internal/resolve/fixup_test.go` is new.

Suite after fixup: **139 passing assertions across 56 top-level tests, 0
failures**, race-clean, `golangci-lint run ./...` → `0 issues`. REQ-MVV re-run
end-to-end and still passing, with its recorded output materially unchanged
(see `coverage.md`).

---

### Resolution A — closes FAIL-2 and ADV-1, ADV-1b

**Defect:** `missingOwned()` ran at `resolve.go:333`, *before* the guard loop at
`:341`, over every row matching on tag pattern alone. A row the seam decided
`GuardFalse` still contributed its `RequiresOwned`.

**Change:** the pre-guard `missingOwned` call is gone. Gating now happens inside
a single `gate(rows, seam, view)` helper that evaluates each row's guard
**first**. A row decided `GuardFalse` is *pruned* — the predicate is decided and
does not hold, so the row is not an edge at all and its owned-state obligation
is not part of this resolution's evaluation. Only rows the guard does not prune
can raise a blocking condition.

**Why guard-before-owned, reversing the previously documented per-row order:**
REQ-15 counts matches "after guard evaluation", so a guard-FALSE row is by
definition not among the edges whose evaluation the disposition rests on.
Keeping the old order made a dead row's requirement decide a live row's fate
(ADV-1) and, worse, masked an escapable `no_match` behind a never-escapable
`owned_state_unavailable`, silently deleting a modeled table capability
(ADV-1b). Among rows the guard does *not* prune, the documented
owned-before-guard precedence is **preserved**: absent owned state is the more
precise diagnosis and is frequently why the seam could not decide the predicate.
No Phase 1 fixture pins the collision — `missingOwnedTable` has no guard and
`unevaluableGuardTable` has its owned key present — so this reordering breaks no
frozen assertion. Recorded as deviation **D8**.

---

### Resolution B — closes FAIL-1 and ADV-2, ADV-2b (incl. widened sub-cases 1d, 1e)

**Defect:** `escapeOrRefuse` (`resolve.go:404-424`) selected escape rows on
`rescues(kind)`, `Outcome`, and `Match` only, then emitted a `Plan`. It called
neither `evaluateGuard` nor `missingOwned`, on either rescuable class, live seam
or nil.

**Change:** `escapeOrRefuse` now routes its escape candidates through the **same
`gate`** before the exact-one count is taken. A guard-FALSE escape is pruned and
does not rescue; an escape missing owned state or carrying an undecidable guard
raises that typed refusal instead of emitting a plan whose `Writes` the kernel
cannot justify.

This is deliberately the shared-gate shape both verifiers recommended, adopted
on its merits: A and B were one defect seen from two directions — gating applied
at inconsistent pipeline points — and the failure mode they warned about was a
fix that patched one path and left the other. Routing both phases through one
helper makes that class of divergence a compile-time impossibility rather than a
review obligation. `TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates`
pins the property directly.

Sub-cases **1d** (guarded escape, nil seam) and **1e** (ambiguous-class escape)
were uncovered by Phase 3b; both fall out of the shared gate and are now pinned
by `TestFixup1d_*` and `TestFixup1e_*`, each verified red against the pre-fix
kernel.

**Consistency with the D1 narrow reading:** unchanged. Escapes still rescue only
`no_match` and `ambiguous_match`; `owned_state_unavailable` and
`guard_unevaluable` are still never escapable. The fix does not widen what an
escape can rescue — it makes an escape edge earn its rescue on the same terms as
any other edge. If anything it *strengthens* the narrow reading: an escape row
that itself hits an unavailable/unevaluable condition now surfaces that
condition rather than papering over it. No change to the D1 calculus.

---

### Resolution C — closes FAIL-3 and ADV-3, ADV-3b (incl. the widened `ambiguous_match` surface)

**Defect:** three diagnosis payloads were functions of `Table.Rows` position
rather than of the input tuple. REQ-2 identifies a resolution by table
*revision*, not row sequence, so two orderings of the same row set at the same
revision are the same input.

**Change, one per surface:**

- `missingOwned` sorts its keys with `slices.Sort` before returning.
- `rowRefs` sorts by source identity (`RuleID`, then `SourceLocator`) via a new
  `compareRefs`. This fixes the `ambiguous_match` surface — **new relative to
  ADV-3** — at its single call site, and `rowRefsRequiring` was rewritten to
  funnel through `rowRefs` so it inherits the ordering rather than duplicating
  it.
- `guard_unevaluable` no longer returns on the first undecidable row in slice
  order. `gate` collects **all** undecidable rows, reports them via the sorted
  `rowRefs`, and picks the reported `Guard` by lowest rule identity
  (`slices.MinFunc` over `compareRefs`) instead of by slice position. Reporting
  all of them is also the better diagnosis: a table with two undecidable
  predicates has two problems, and the old payload disclosed only one.

The degraded escape-ambiguity payload is a distinct `rowRefs` call site and is
covered separately by `TestFixup3c_DegradedEscapeAmbiguityPayloadMustNotDependOnRowOrder`.

---

### Explicit non-defects — confirmed not "fixed"

The two Phase 3a non-violations were left exactly as they were: the unreserved
`"recognized"` tag key (deterministic and refusal-shaped, not a breach), and the
silently-unreachable row whose `Escape` list holds only a non-rescuable class
(correct kernel behaviour; an RDR 0006 lint obligation). The escape path's
correct deference to the outcome-alphabet gate is likewise untouched and still
verified by the suite.

---

## Phase 3b — Adversarial Review (independent lens)

Reviewer stance: assume the implementation is wrong and prove it. Scope is the
RDR 0001 **Failure Modes** clause:

> "Visible failures are typed refusal values: `no_match`, `ambiguous_match`,
> `owned_state_unavailable`, `guard_unevaluable`, or `unmodeled_outcome`. Silent
> failure would mean the kernel guessed a transition or executed persistence
> directly; both are prohibited by the normative contract. Diagnosis starts with
> the input tuple, the refusal kind, and the transition table revision used for
> that resolution."

Tests added in `internal/resolve/adversarial_test.go`. No existing test or
source file was modified. The Phase 1 suite (`TestReq*`, `TestMVV`) still passes
in full; every failure below is new signal, not a regression.

**6 of 8 added tests currently FAIL against the implementation — all six are
real defects.** Two pass and are retained as regression guards.

---

### ADV-1 — a guard-FALSE row's `RequiresOwned` poisons an otherwise legal resolution

**Status: FAILS (defect).**
`TestAdv1_GuardFalseRowsRequiresOwnedMustNotPoisonAnExactOneMatch`

FM clause: "Visible failures are typed refusal values … `owned_state_unavailable`",
read with REQ-15 ("the only successful selection is exactly one matching edge
**after guard evaluation**") and the req-list ASSUMPTION binding
`owned_state_unavailable` to a row "whose *evaluation* requires an owned tag
absent from the accessor-produced owned snapshot".

`Resolve` calls `missingOwned(candidates, view)` at resolve.go:333 — **before**
the guard loop at :341. `candidates` is every row matching on tag pattern alone.
A row the seam decides `GuardFalse` is not a surviving candidate and its
evaluation never needed its owned tags, yet its `RequiresOwned` entry still
converts a clean exact-one match into `owned_state_unavailable`.

Observed: kind `owned_state_unavailable`, `MissingOwned=[never-present]`.
Expected: a plan for `rdr.live`.

This is a **wrong-refusal** defect, not a missing refusal — the kernel refuses a
transition it should have planned, and names a condition that does not hold for
any surviving edge. It is the more damaging direction: the caller sees a
diagnosis pointing at owned state that is irrelevant to the edge they wanted.

**ADV-1b (`TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`) — also FAILS.**
Same root cause, sharper consequence. When *every* ordinary row is guard-FALSE,
the honest condition is zero-match → `no_match`, which RDR 0002 makes an
**escapable** class. Because the premature `missingOwned` check fires first and
returns `owned_state_unavailable` (never escapable), the table author's modeled
`no_match` escape edge is silently unreachable. A correctness bug in the
precedence order removes a modeled table capability.

**Fix direction:** move the `missingOwned` check to run over the rows that
*survive* guard evaluation, not over all pattern-matching candidates. Note the
ordering interaction: a row whose `RequiresOwned` is missing may also be the row
whose guard cannot be evaluated, so the fixup must pick and document a
deterministic precedence between the two.

---

### ADV-2 — an escape edge bypasses the guard seam entirely

**Status: FAILS (defect), both sub-cases.**
`TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam`

FM clause: "Silent failure would mean the kernel **guessed a transition**",
read with REQ-15 ("exactly one matching edge **after guard evaluation**") and
REQ-23 ("guard evaluation **delegation**").

`escapeOrRefuse` (resolve.go:404-424) filters escape rows on `rescues(kind)`,
`Outcome`, and `Match` — then emits a `Plan` at :418. It **never calls
`evaluateGuard`**. Consequences:

- *guard FALSE*: the kernel emits a plan for an edge whose own predicate says it
  does not hold. Observed: plan `rdr.escape.guarded`, `escaped=true`. Expected:
  `no_match` refusal. This is literally "guessed a transition".
- *guard UNEVALUABLE*: the kernel emits a plan for an edge it cannot know holds.
  Observed: plan `rdr.escape.guarded`. Expected: `guard_unevaluable`.

The escape path is the one place the kernel is most permissive, and it is
exactly where the guard seam is skipped — a guarded escape edge is strictly
more dangerous than a guarded ordinary edge, because it fires only when the
table has already failed to resolve normally.

**ADV-2b (`TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement`) — also FAILS.**
The same path skips `RequiresOwned`. Observed: a plan carrying
`Writes=[{status Escaped}]` derived from owned state the accessor snapshot never
produced. Per FM ("or executed persistence directly") and REQ-19, `Writes` is the
persistence description handed to the RDR 0004 accessor layer — describing a
write off absent owned state is the kernel guessing at persistence.

**Fix direction:** the escape candidate set must pass through the same
`RequiresOwned` and guard gates as ordinary candidates before the exact-one
count is taken. Note this changes the rescue arithmetic: an escape row rejected
by its guard should not count toward the "matches exactly once" test at :416.

---

### ADV-3 — refusal diagnosis payload depends on `Table.Rows` order

**Status: FAILS (defect).**
`TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`

FM clause: "**Diagnosis starts with the input tuple, the refusal kind, and the
transition table revision** used for that resolution", read with REQ-1 (same
input → same disposition) and REQ-3 (value-level replay determinism).

The guard loop returns on the **first** `GuardUnevaluable` row in slice order
(resolve.go:348-354), so `Refusal.Guard` and `Refusal.Rows` name whichever
undecidable row happens to sit earlier in `Table.Rows`:

```
rows [a,b] -> guard="unknown-a" rows=[{rdr.unevaluable.a flows/rdr.toml:10}]
rows [b,a] -> guard="unknown-b" rows=[{rdr.unevaluable.b flows/rdr.toml:20}]
```

Same flow, same revision, same owned snapshot, same observed tags, same
recognized outcome — different reported diagnosis. REQ-2 defines the input tuple
by the transition table **revision**, not by row sequence, so two orderings of
the same row set at the same revision are *the same input* and must replay the
same disposition. Row order is a normalization detail RDR 0002 owns; it must not
reach the reported diagnosis.

**ADV-3b (`TestAdv3b_MissingOwnedPayloadMustNotDependOnTableRowOrder`) — also FAILS.**
Same class, `owned_state_unavailable` payload:

```
rows [a,b] -> missing=[alpha beta]
rows [b,a] -> missing=[beta alpha]
```

`missingOwned`'s doc comment claims it returns keys "in stable order". It is
stable with respect to *row order*, which is not the same as stable with respect
to the *input tuple* — the property REQ-1 and REQ-10 actually require.

Phase 1's `TestReq3_ValueLevelReplayIsIndependentOfInputSliceOrder` covers
`Owned`/`Observed` slice order only. `Table.Rows` order is the surface that
actually leaks, and it was untested.

**Fix direction:** sort the diagnostic collections (`MissingOwned`, and `Rows`
by `RuleID`/`SourceLocator`) before they enter a `Refusal`; for
`guard_unevaluable`, either report all undecidable rows rather than the first,
or select deterministically by rule identity instead of slice position.

---

### Regression guards (pass today — implementation is correct here)

These two were written expecting a defect and did not find one. Retained so a
later refactor cannot silently regress the provenance boundary.

- **ADV-4** `TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot` — **PASSES.**
  `assemble` (resolve.go:145-161) applies `Observed` before `Owned`, so owned
  wins on key collision. A caller supplying `status:Final` as an observed tag
  cannot steer the kernel off the artifact's real `status:Draft` owned state.
  Correct as written; the ordering in `assemble` is load-bearing and undocumented
  as such at the call site.
- **ADV-5** `TestAdv5_ObservedTagCannotSatisfyAnOwnedStateRequirement` — **PASSES.**
  `missingOwned` checks `view.has(key, ProvenanceOwned)`, so an observed tag
  carrying a required owned key does **not** satisfy the requirement; the kernel
  still refuses `owned_state_unavailable` and names the key. Correct: caller
  context cannot stand in for accessor-produced artifact state.

---

### Note for the fixup phase

ADV-1/1b and ADV-2/2b are the same underlying shape from two directions: the
gating checks (`RequiresOwned`, guard) and the candidate set are applied at
inconsistent points in the pipeline. Ordinary candidates get `RequiresOwned`
checked *too early* (before guards prune them); escape candidates get both
checked *not at all*. A fix that only patches one path will leave the other.
The cleanest repair is a single `viable(row) → (ok, refusal)` gate applied
uniformly to ordinary and escape candidates alike, with the exact-one count
taken over the survivors in both phases.

ADV-3/3b are independent of that and can be fixed separately — they are payload
ordering, not selection semantics.

---

## Phase 3a — Chain-of-Verification (independent lens)

Method: for each REQ in `artifacts/req-list.md` (37 REQs + REQ-MVV), name an
input that would make a correct implementation visibly violate the clause, then
run it. Probes were derived **from the RDR and req-list only** — the Phase 1
test files (`resolve_test.go`, `mvv_test.go`, `fixtures_test.go`,
`boundary_test.go`) and the Phase 3b section above were not read before the
probe set was designed, so this is a genuinely independent lens.

Harness: throwaway `main` programs under `.run-cove/` calling the exported
package API. **77 probes across 12 classes; 3 defects.** All findings below were
observed by running code and reproduce deterministically (100–200 replays each).

**Convergence note:** this pass independently rediscovered Phase 3b's ADV-1,
ADV-2, and ADV-3 from spec-derived inputs. FAIL-1/2/3 below are **the same three
defects**, not new ones — two independent lenses landing on the same three
findings raises confidence that the defect set is complete, and each entry adds
sub-cases 3b did not record. The fixup phase should treat ADV-n and FAIL-n as one
work item apiece.

---

### FAIL-1 — the escape-rescue path skips both gating checks (= ADV-2, widened)

**REQs violated: REQ-5, REQ-12, REQ-15, REQ-23.**

`escapeOrRefuse` (resolve.go:404-424) selects escape rows on `rescues(kind)`,
`Outcome`, and `Match` only, then emits a `Plan` at :418. It never calls
`missingOwned` and never calls `evaluateGuard`. REQ-5 makes refusal on
unavailable owned state and undecidable guards an unconditional MUST with no
escape carve-out; REQ-15 requires the escape edge to "itself match exactly once",
and REQ-15 defines matching as holding "after guard evaluation".

Minimal reproducer 1a — escape row requiring absent owned state:

```go
in := resolve.Input{
    Flow: "rdr",
    Table: resolve.Table{
        Revision: "rev-1", Outcomes: []string{"successful"},
        Rows: []resolve.Row{{
            RuleID: "esc", SourceLocator: "table.toml:9", Outcome: "successful",
            Escape:        []resolve.RefusalKind{resolve.KindNoMatch},
            RequiresOwned: []string{"status"},
            NextTags:      []resolve.Tag{{Key: "status", Value: "Held"}},
            Writes:        []resolve.Tag{{Key: "status", Value: "Held"}},
        }},
    },
    Owned:      nil, // accessor snapshot does NOT carry "status"
    Recognized: "successful",
}
```

Observed: `Plan{RuleID:"esc", Escaped:true, Writes:[{status Held}]}`, nil error.
Expected: `owned_state_unavailable`. **Control:** the byte-identical row with
`Escape: nil` correctly refuses `owned_state_unavailable`, proving the `Escape`
field alone causes the bypass.

Sub-cases, all observed as `Plan{RuleID:"esc", Escaped:true}` with writes:

| # | Escape row condition | Seam verdict | Expected | Observed |
| --- | --- | --- | --- | --- |
| 1a | `RequiresOwned:["status"]`, empty snapshot | n/a | `owned_state_unavailable` | plan + writes |
| 1b | `Guard:"iterations >= 3"` | `GuardFalse` | `no_match` | plan + writes |
| 1c | `Guard:"iterations >= 3"` | `GuardUnevaluable` | `guard_unevaluable` | plan + writes |
| 1d | `Guard:"g"`, `Guards: nil` | n/a (nil seam) | `guard_unevaluable` | plan + writes |
| 1e | ambiguous-class escape, `Guard` FALSE **and** owned missing | `GuardFalse` | `ambiguous_match` | plan + writes |

1b is the sharpest: the delegated seam returned `GuardFalse` and the kernel
emitted the transition anyway — literally REQ-12's "the kernel guessed a
transition". 1d and 1e are **new relative to ADV-2**: the bypass is not specific
to a live seam (a nil seam is also bypassed, whereas an ordinary guarded row with
a nil seam correctly refuses) and not specific to the `no_match` class (the
`ambiguous_match` escape class bypasses identically). A fix must cover the escape
path for *both* rescuable classes and the nil-seam case.

Scoped correctly (verified, not a defect): the escape path does **not** bypass
the outcome-alphabet gate — an escape row for an unmodeled outcome still refuses
`unmodeled_outcome`.

---

### FAIL-2 — a guard-FALSE row's `RequiresOwned` suppresses a legal edge (= ADV-1)

**REQ violated: REQ-15.**

`missingOwned(candidates, view)` runs at resolve.go:333, **before** the guard
loop at :341, over every row matching on tag pattern alone. A row the seam
decides `GuardFalse` is not a surviving candidate — REQ-15 counts matches "after
guard evaluation" — yet its `RequiresOwned` still converts an otherwise clean
exact-one match into `owned_state_unavailable`.

Minimal reproducer:

```go
in := resolve.Input{
    Flow: "rdr",
    Table: resolve.Table{
        Revision: "rev-1", Outcomes: []string{"successful"},
        Rows: []resolve.Row{
            {RuleID: "good", SourceLocator: "table.toml:1", Outcome: "successful",
             Match: []resolve.Tag{{Key: "status", Value: "Draft"}},
             RequiresOwned: []string{"status"},
             NextTags: []resolve.Tag{{Key: "status", Value: "Final"}},
             Writes:   []resolve.Tag{{Key: "status", Value: "Final"}}},
            {RuleID: "dead", SourceLocator: "table.toml:2", Outcome: "successful",
             Match: []resolve.Tag{{Key: "status", Value: "Draft"}},
             Guard: "never", RequiresOwned: []string{"reviewer"},
             NextTags: []resolve.Tag{{Key: "status", Value: "Rejected"}}},
        },
    },
    Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
    Recognized: "successful",
    Guards:     /* seam returning GuardFalse for "never" */,
}
```

Observed: `Refusal{Kind: owned_state_unavailable, MissingOwned: ["reviewer"]}`.
Expected: `Plan{RuleID: "good"}`. **Controls:** removing the `dead` row yields the
expected plan; a row that simply fails its `Match` pattern does *not* poison the
resolution — so the defect is specific to rows pruned by the guard seam rather
than by pattern matching.

This is a wrong-refusal, and the diagnosis names an owned key (`reviewer`)
irrelevant to the edge the caller wanted.

---

### FAIL-3 — refusal diagnosis payloads depend on `Table.Rows` order (= ADV-3, widened)

**REQs violated: REQ-1, REQ-2, REQ-10.**

REQ-2 defines the input tuple by the transition table **revision**, not by row
sequence, so two orderings of the same row set at the same revision are the same
input and must replay the same disposition. Three payload surfaces leak row
order (all at `Revision: "rev-1"`, identical owned/observed/recognized inputs):

```
guard_unevaluable       rows [a,b] -> guard="unknown-a" rows=[{unevaluable.a flows/rdr.toml:10}]
                        rows [b,a] -> guard="unknown-b" rows=[{unevaluable.b flows/rdr.toml:20}]

owned_state_unavailable rows [a,b] -> missing=[alpha beta]
                        rows [b,a] -> missing=[beta alpha]

ambiguous_match         rows [a,b] -> rows=[{amb.a t:10} {amb.b t:20}]
                        rows [b,a] -> rows=[{amb.b t:20} {amb.a t:10}]
```

The third surface (`ambiguous_match` `Refusal.Rows` order, from `rowRefs` at
:452) is **new relative to ADV-3**, which recorded only the `guard_unevaluable`
and `owned_state_unavailable` payloads. Any fix that sorts diagnostic collections
must cover `rowRefs` as well as `missingOwned`.

Note the refusal *kind* is row-order-stable in all three cases; only the
diagnostic payload varies. Under a stricter reading this is REQ-10 diagnosis
instability rather than REQ-1 disposition instability — but since `Refusal` is a
single value and REQ-1 requires "the same disposition", the payload is part of
the disposition.

---

### Probe classes run (no violation observed)

Each class ran spec-derived inputs designed to break the clause; all held.

- **Disposition totality (REQ-1)** — zero-value `Input`, empty table, nil rows,
  row with all-nil fields, empty recognized outcome: every probe returned exactly
  one of `Plan`/`Refusal`, never both, never neither, always nil error.
- **Replay determinism (REQ-1/2/3/4/29)** — 200 replays each over plan, ambiguous,
  and missing-owned dispositions; independently constructed identical tuples;
  owned/observed key-clash inputs. No map-iteration leakage into output.
- **Go error path (REQ-6)** — all five refusal kinds plus the zero-value input
  return `err == nil`. No modeled refusal travels the error path.
- **Closed kind set (REQ-7/8)** — `RefusalKinds()` returns exactly the five
  contracted constants with the contracted string identifiers; `Refusal.Kind` is a
  distinct package-local comparable named type (not `error`, not a bare string);
  mutating the returned slice does not corrupt package state.
- **Diagnosis echo (REQ-10)** — every refusal kind echoes `Flow`, `Revision`, and
  `Recognized`; plans echo `Revision`.
- **Prohibitions (REQ-11/13/30/32)** — structurally, `Refusal` carries no
  `Writes`/`NextTags`/`Plan` field, `Input` carries no `io.Writer`/`*os.File`/
  persistence collaborator, and every refusal disposition has `Plan == nil`.
- **Statelessness (REQ-14/27)** — `Resolve` does not mutate the caller's `Table`,
  `Owned`, or `Observed`.
- **Aliasing (REQ-4/24)** — mutating a returned `Plan.NextTags`/`Writes`,
  `Refusal.Rows`, or `Refusal.MissingOwned` does not corrupt the caller table or
  a later replay.
- **Selection semantics (REQ-15/16/31)** — two pattern-matching rows with exactly
  one guard TRUE correctly yields exact-one success (ambiguity is judged after
  guards for *ordinary* rows); two matching escape rows correctly do not rescue;
  cross-class escape does not rescue (a `no_match` escape does not rescue
  ambiguity and vice versa); an escape row does not create ambiguity alongside a
  real match; rescued plans are marked `Escaped`.
- **Tag-set assembly (REQ-17)** — owned, observed, and recognized tags all reach
  the seam's view with correct provenance; a row can match on the recognized tag;
  owned wins over observed on key collision, stably across 100 replays.
- **Guard delegation (REQ-23)** — the seam is called once per guarded row and not
  at all for unguarded rows; a syntactically nonsense guard the seam decides TRUE
  still selects the row, confirming the kernel implements no operator semantics.
- **Outcome isolation (REQ-15)** — rows for a *different* recognized outcome
  contribute no owned obligations, no ambiguity, and trigger no guard calls.

Doc-only REQs (REQ-19, REQ-20, REQ-21, REQ-22, REQ-25, REQ-26, REQ-28, REQ-34,
REQ-35, REQ-36, REQ-37) were checked by inspection of the package surface rather
than by probe: the package is `internal/resolve`, imports only `slices` from the
standard library (REQ-26, REQ-36, REQ-37: no third-party dep, no encode/decode,
no hashing or canonical serialization), exposes only kernel values, and the
package name is "resolve" rather than an orchestrator/runner name (REQ-25).

---

### Observations for the fixup phase (not violations)

Two probes produced surprising behaviour that is nonetheless contract-compliant.
Recorded so the fixup phase does not "fix" them by accident, and so RDR 0002/0006
can pick them up as lint obligations.

1. **`recognizedTagKey` is an unreserved namespace.** A caller-supplied owned or
   observed tag literally keyed `recognized` shadows the freshly recognized
   outcome in the assembled view (`assemble` writes recognized first, then
   observed, then owned), turning a would-be plan into `no_match`. The RDR fixes
   no reserved key name, and the disposition stays deterministic and refusal-shaped
   rather than a guess, so this is not a REQ-17 breach. It is a collision hazard
   worth either reserving the key or documenting at the boundary.

2. **A row whose `Escape` list contains only a non-rescuable class is silently
   unreachable.** `Resolve` excludes any row with `len(row.Escape) != 0` from
   ordinary candidacy, but `escapeOrRefuse` only ever consults rows rescuing
   `no_match`/`ambiguous_match`. A row with `Escape: ["unmodeled_outcome"]` or a
   garbage class is therefore dead: if it is the only matching edge, the kernel
   reports `no_match`. RDR 0002 forbids such tables, so refusing rather than
   guessing is correct kernel behaviour — but it is invisible at runtime and
   belongs in RDR 0006 graph lint. (`Escape: []` non-nil-but-empty correctly
   stays an ordinary candidate.)
