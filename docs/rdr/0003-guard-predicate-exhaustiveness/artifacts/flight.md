# Flight — batch:rdr-0003

`/kata-flight --label batch:rdr-0003 --drain`, run as the `--close-and-flight`
tail of `/rdr-implement-triage 3`.

```
flight: 4 shipped, 0 stopped, 0 skipped (over 2 waves)
  shipped: x0fp=3bd7d83  9yeq=7ac1e67  0x69=e72a858  s1r1=508c263
```

## Wave 1 — 5 resolved, 3 shipped, 2 routed out

The scope-review gate paid for itself: two katas would have shipped fixes
that silently broke locked contracts.

| Kata | Verdict | Outcome |
|---|---|---|
| `x0fp` P1/high | IN-SCOPE | shipped `3bd7d83` |
| `9yeq` P3/low | IN-SCOPE | shipped `7ac1e67` |
| `0x69` P3/low | IN-SCOPE | shipped `e72a858` |
| `cq5p` P2/med | RDR-SHAPED | → `kind:rdr-seed` |
| `x2bp` P2/med | RDR-SHAPED | → `kind:rdr-seed`, blocked by `ge67` |

`fyf4` was held (`inbox:hold`) by author decision before the wave: its own
course of action defers it to RDR 0006's implementation (BUILD-ORDER run
4), and RDR 0006 is Final but unimplemented.

### Why the two were routed out

- **`cq5p`** — its prescribed fix (thread tag declarations into the
  evaluator) would falsify `0007:A17`, a **Verified** assumption that the
  seam "needs neither the view, nor provenance, nor sibling atoms, nor tag
  declarations at evaluation time". `guard.Evaluator` is `struct{}` and a
  reflection test pins that. `0007:A18` left the fork open for RDR 0003 to
  decide, and RDR 0003's own conformance clause points at the other branch
  (view-assembly). Two mutually exclusive cross-RDR contract changes.
- **`x2bp`** — grounding DISPROVED a premise in the kata's own body (written
  during triage): the claim that "the machinery already exists" for
  presence-only projection is false for the reproduced case. `Denotation`'s
  fallback is gated on `!ok` from `valueAssignments`, which SUCCEEDS for an
  unmarked enum, so the defect routes around it. The fix would also redefine
  "dimension" in REQ-58 from participating KEY to value/presence
  sub-dimension — the same predicate `ge67` was seeded to adjudicate.

## Wave 2 — 1 shipped

`s1r1` was minted by `9yeq`'s refine loop and re-swept by `--drain`.

Filed as "ADV-2 needs an instrumentation seam". Grounding **reframed** it:
`TestReq133_TheBoundsVerdictAgreesWithWhetherTheProofCompletes` was
TAUTOLOGICAL — it computed `card <= Bound()` and compared it to
`ProofCompletes`, which IS `card <= Bound()`. Since `ProofCompletes`' doc
comment names A15 as refuted "if the two ever disagree", and **A15's Status
is `Pending` with Method `MVV Test`**, the tautology meant A15 appeared
discharged while untested.

Fixed with NO new public surface: completion is now observed independently
via the already-exported `Product` and `Denotation` under a deadline.

## Defects found beyond the filed scope

The refine loops found more than the katas named:

- **`x0fp`** surfaced two further defects of the same class inside its own
  diff: the `{MaxInt..MaxInt}` enumeration hang (`n <= *Max` unfalsifiable
  once `n++` wraps) and the optional-presence factor overflow (`count *= 2`
  wrapping to `-2`, reading UNDER the bound → **green over a full-int-range
  dimension**). Both are D12's defect one operation later.
- **`9yeq`** proved its own thesis by mutation: against a real regression
  (D11's bound guard removed), the OLD wall-clock oracle PASSED (1.1x,
  0.00s) while the new structural assert FAILED in 1.41s.

## Residuals filed, not shipped

| Kata | Why it did not ship |
|---|---|
| `dh40` `kind:rdr-seed` | REQ-133's oracle is still bound-gated (`bound` 2048→64 passes). Needs an ungated completion seam — a new public surface RDR 0003 does not name. |
| `r5ja` `kind:rdr-seed` | RDR 0006 vs RDR 0003 REQ-77 prescribe opposite outcomes on an overlap-free group. |
| `ge67` `kind:rdr-seed` | REQ-93 "known finite" vs REQ-94 "provable" are unreconciled in the record. |
| `cq5p`, `x2bp` `kind:rdr-seed` | See wave 1 above. |
| `fyf4` `inbox:hold` | Deferred to RDR 0006's implementation (run 4). |

## Honest limits

- `s1r1` strictly increases discrimination but does not close A15's adequacy
  question: every exported enumeration seam consults `Bound()`, so an
  over-conservative bound still passes. Recorded in the test comment and
  filed as `dh40` rather than overclaimed.
- The enumerate-then-check cost regression remains uncovered by design — it
  is value-invisible (verdicts identical, only cost differs), so no
  completion oracle can see it. A wall-clock substitute is forbidden by
  REQ-86 and was deliberately removed from ADV-2 by `9yeq`.

Every shipped kata: rebased onto current main, `go test ./...` + `-race` +
`golangci-lint` green, ff-merged, worktree torn down, closed with typed
evidence.
