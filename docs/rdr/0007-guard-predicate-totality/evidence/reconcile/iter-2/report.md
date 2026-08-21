Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0007 guard predicate totality (iteration 2)

Second Stage 6 on this RDR. Iteration 1 (`../report.md`) reconciled the
pre-demotion draft and the RDR locked Final 2026-08-12 (`a2341b1`). JDR 0001
then **demoted** it (`76a9e67`) and the draft was re-proposed, re-refined and
re-resolved against a materially different design. This pass reconciles that
re-entered draft.

## Verdict: NOT RECONCILED — return to Stage 5

The open set was **not built**. The reconcile prompt gates on Stage 5
completeness before the open set is assembled, and that preflight fails: two
of the four lenses this RDR's `Profile` requires have never reviewed the
current design.

## Stage 5 completeness preflight — FAIL

`Profile: foundational` → lens row `cove → 3amigo → critique → repeatability`
(rdr-common §lens-row). Lens evidence on disk:

| Lens | iter-1 (pre-demotion) | iter-2 (current design) | Owed |
| --- | --- | --- | --- |
| cove | present | **present** | — |
| 3amigo | present | **present** | — |
| critique | present | **ABSENT** | **yes** |
| repeatability | present (full: 3 runs, 3 stamps + diff) | **ABSENT** | **yes** |

The `foundational` row is a required checklist against the *current* draft.
Its members are not discharged by evidence produced against a superseded one.

### The iteration-1 critique and repeatability evidence is stale, not merely old

This is the load-bearing finding, and it is a fact about content, not a
folder-date heuristic. The re-entry rewrote the document — `git diff
a2341b1..HEAD` on the RDR body is **1666 insertions / 2137 deletions** against
a 1927-line file — and the rewrite changed the very seam both lenses reasoned
about:

- **Critique iteration 1's central finding no longer has a referent.** C-2
  reads: "Every atom-level normative clause binds a seam that cannot see
  atoms. `Row.Guard` is `string`; `Evaluate(guard string, view TagSet)` is the
  whole seam; no RDR states the mapping." That was a finding against
  assumption **A10** (the unowned guard-structure mapping). JDR 0001 §D1 then
  put parsed atoms on the row, dissolving the problem class. **A10 does not
  exist in the current draft** (grep: 0 hits). C-7 likewise turns on
  `Refusal.Guard` carrying "an opaque guard string" — a field this draft
  deletes.
- **Repeatability iteration 1 reconstructed the wrong seam.** `run-1.md:16`
  states "It quotes the seam as `Evaluate(guard string, view TagSet)`." The
  current SEAM clause specifies `Evaluate(atom, value) GuardResult` and adds
  "It never sees the view." All three runs reconstructed the pre-§D1 design,
  so the three-way diff attests determinacy of a document that has since been
  replaced.

Neither lens has been run against the kernel-enforced, atom-carrying design
that is now the RDR's whole proposal. Their `Model:` stamps are valid
(§model-stamp: `claude-opus-5[1m]` / `claude-fable-5` for the critique pair;
three distinct stamps for repeatability) — the defect is the *subject*, not
the model draw.

Determinacy trigger: n/a as a separate obligation — at `foundational` the full
repeatability lens is a row member, so it is owed as a lens rather than as a
trigger.

## What DID pass

Recorded so the return pass does not redo it.

**Absorption audit (delegated) — PASS, no residue.** All 48 iteration-2
findings (cove F-1..F-13; 3amigo PM-1..10, IMP-1..11, QA-1..14) carry a ledger
row. Every one marked `fixed` was verified present in the RDR at the named
section; every `dismissed-with-cite` had its cite confirmed in-text; all five
`charted-to-successor` items are genuinely absent from the body and recorded
in `Charted.md`. Cove iter-2 charted nothing. The one iteration-1 re-opening
(F-2 flipping A3 Verified → Pending) is reflected in A3's Status.

**Mechanical completeness — PASS.** No `_Draft placeholder._`, no seed-skeleton
header, no surviving template bracket. `## References` carries 13 real
citations (JDR 0001, ISO/IEC 9075, SCXML, PostgreSQL 16, RDRs 0001–0009,
`internal/resolve` symbols, the prior-art and spike caches).

**Exactness-word delta — CLEAN.** Swept the `normative` blocks; each
MUST/never/every/only/exactly/sole claim traces to a named assumption or a
source pointer.

**Three source-grounded checks run this pass** (they will shorten the return
pass, and one materially strengthens an open assumption):

1. **A23 is better-supported than its Status says.** No frozen test constructs
   a `Refusal` composite literal (`grep 'Refusal{' internal/resolve/*_test.go`
   → 0 hits), so nothing pins the field set positionally. The four
   `reflect.DeepEqual` sites (`fixup_test.go:200,236`;
   `adversarial_test.go:313,375`) compare two `Refusal` values *to each other*
   for order-invariance — agnostic to which fields exist. A23's claim ("no
   shipped test asserts `Refusal`'s field set exhaustively") holds against
   source today.
2. **The SEAM clause's Fixup-1d reasoning is grounded.**
   `fixup_test.go:44` sets the guard to `iterations >= 3`; the test's input is
   `noMatchInput()`, which inherits `legalInput()`'s view — `status` (owned),
   `reviews` (observed), `recognized`. `iterations` is absent, so the kernel
   decides that atom unevaluable on absence and the nil seam is never
   consulted, exactly as the clause states.
3. **The A3 spike diverges from the specified shape on THREE axes, not one.**
   Cove F-2 caught the first; the other two are recorded here for the
   re-spike's scope. (a) It retained `Refusal.Guard` and added
   `UndecidedAtoms` beside it (`a3-reshape-probe.md:156`) where the contracts
   say *replace*. (b) Its seam is `EvaluateAtom(atom GuardAtom, value string,
   view)` (`:155`) — it **still passes the view**, which the SEAM clause
   forbids ("It never sees the view"). (c) Its existence atom is
   `{Key: "guard.missing", Op: resolve.OpExists}` with no `Literal`
   (`spike_kleene_test.go.txt:77`), which the `presence == literal` rule
   forbids and which A24's fail-closed clause now makes `uncomparable`.

## Open set — NOT BUILT

Deliberately not assembled. Building it now would force terminal dispositions
on assumptions that two required lenses have not yet examined against this
design — which is the exact failure this gate exists to prevent, one stage
early. For the return pass's orientation, the open items stand at: **A3, A19,
A22, A23, A24** (Pending) and **A12** (Deferred). A23 and A24 are net-new from
the 3amigo iteration-2 pass and have never been reviewed by any cross-model
lens.

Note that A3, A23 and A24 all resolve through **one** re-spike (the reshape
with `Refusal.Guard` actually removed), and its scope is now known to include
the two additional divergences named above.

## Hard rules — not reached

Neither hard rule was evaluated; both are downstream of the open set.

- **Refutation → BLOCKER**: no refutation is live in the current draft. A19's
  own text records a refutation ("the original 'without a new envelope field'
  reading is refuted"), but it was resolved in the honest direction — the RDR
  was corrected to match shipped `clierr.CLIError`, and the consequence routed
  to a Prerequisite. That is a closed route-back, not an open blocker.
- **No MVV-critical deferral**: flagged for the return pass, not adjudicated
  here. MVV scenario 1 asserts "the absent key in the payload," and the
  payload's named surface is **A23**, which is Pending. Whether that makes A23
  a pre-lock prerequisite rather than a Phase-1 obligation is a live question
  the next Stage 6 must answer — it is not obviously survivable, because the
  MVV cannot be written without the field it asserts on.

## Next

Run the two missing lenses against the current draft, then re-run this gate:

```
/rdr-prelock 0007 critique
/rdr-prelock 0007 repeatability
/rdr-reconcile 0007
```

`critique` is `foundational`, so it owes the dual-model draw (§model-stamp);
`repeatability` owes the full variant (three runs, three distinct stamps, plus
`diff.md`). Both land under `evidence/<lens>/iter-2/`. The alt-model roster is
in `{RDR_RESOURCES}`.
