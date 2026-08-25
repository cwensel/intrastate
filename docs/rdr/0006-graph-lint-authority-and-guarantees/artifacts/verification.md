
## Phase 3b — Adversarial Review

Three failure modes, each anchored in `## Trade-offs > ### Failure Modes`, and
each with a test that catches it. **All three tests currently FAIL against the
implementation.** Tests live in `internal/graphlint/adversarial_0006_test.go`.

### ADV-1 — the successor join collapses divergent edges, erasing reachable owned-states

**Failure mode.** `Failure Modes`, first silent-failure shape: *"Silent failure
would mean accepting a model with a blocking invariant defect."*
Anchored on **REQ-108** (`LBD, Join rule and termination`): *"two edges reaching
**the same successor** produce one node whose per-tag value sets are the union
of theirs (a lattice widening)."* The join is licensed only between edges that
agree on their successor.

**Defect.** `internal/graphlint/reach.go::successorOf` folds *every* edge leaving
a node into one successor value, not only the edges that reach the same one.
`joinNodes` then drops any key absent on either side ("absence dominates").
Two rows out of the root writing *different* owned keys therefore annihilate each
other: the reachable set loses both keys and stands for a concrete view no path
produces.

**Consequence (the false green).** With `goLeft` writing `status=b, left=l` and
`goRight` writing `status=b, right=r`, the reachable set is `{status=a}`,
`{status=b}` — `left` and `right` are gone. The two live downstream rows
`onlyLeft` and `onlyRight` are then reported `graph-unreachable-rule`, which is
**advisory** and never changes the success disposition. Worse, `groups.go::checkGroups`
skips `checkCoverage` for an unreachable group, so an entire live arm of the
model is never coverage-proved and lint still returns clean. A spurious
`graph-dead-end` is emitted on the collapsed node as well.

**Test.** `TestAdvSuccessorJoinCollapsesDivergentEdges` (ADVERSARIAL) — two
sub-assertions: the reachable set must retain a node holding each divergent key,
and no live row may take `graph-unreachable-rule`.

**Status: FAILS.** Both sub-assertions fail. Node set is `{status=a}`,
`{status=b}`; both `onlyLeft` and `onlyRight` take `graph-unreachable-rule`.

### ADV-2 — a withheld exhaustiveness claim suppresses the group's overlap finding

**Failure mode.** `Failure Modes`: *"reporting only the first defect per group so
a fixed model still fails"* — the section names the **complete-emission clause**
as the mechanism meant to prevent it.
Anchored on **REQ-84** (`0006:C16`): *"Withholding a group's exhaustiveness claim
MUST NOT suppress overlap, coverage, or further withholding findings for that
group."*

**Defect.** `internal/graphlint/groups.go::emitOverlaps` skips any row whose
accepted-assignment set is not projectable (`if !left.Projectable() { continue }`).
A row carrying an atom over an **optional** key is non-projectable. Adding one
optional-key guard atom to each of two genuinely-overlapping rows therefore makes
the `graph-overlap` finding vanish — while the same two rows overlap loudly
without it. The overlap is decidable independently of the withholding: it lives
on the finite, fully-declared dimension `p`. Declining to decide it because a
*different* dimension is unprovable is exactly the suppression REQ-84 forbids.

**Test.** `TestAdvWithholdingSuppressesOverlap` (ADVERSARIAL, table-driven) —
paired control and attack over a byte-identical overlap on `p`:
- *control* (no optional key): `graph-overlap` naming `(r1, r2)` — **PASSES**;
- *attack* (same overlap + one optional-key atom per row): must still emit
  `graph-overlap` — **FAILS**, the report carries only the two
  `graph-unprovable-coverage` withholdings.

The passing control is what proves the attack asserts a defect that *was*
detectable, rather than an overlap lint could never have seen.

**Status: FAILS** (attack arm; control arm passes).

### ADV-3 — dead end reads a merged node existentially for its outgoing-row test

**Failure mode.** `Failure Modes`, first silent-failure shape again: *"accepting a
model with a blocking invariant defect."*
Anchored on **REQ-36 / REQ-111** (`INV 2`; `LBD, Soundness direction`): dead end
is a **universal** check, and *"a ∀-claim evaluated against a widened domain gets
easier to satisfy, which is the false-green direction"*; the stated remedy is to
*"**split** the node on terminal-participating keys first, recovering exactness"*
because *"partial satisfaction is deliberately not enough — accepting it would let
a merged node close on a path that has not actually terminated."*

**Defect.** `internal/graphlint/analysis.go::checkDeadEnd` splits on
`terminalKeys()` and then asks `hasOutgoingOrdinaryRow`, which is **existential**:
`matchSatisfiable` → `ownedAtomSatisfiable` accepts a node when *some* held value
satisfies each atom. The split therefore recovers exactness for the *terminal*
half of the test and leaves the *outgoing-row* half reading the merged node
directly. The gap is a node multi-valued on a key participating in **no**
terminal: splitting never touches that key.

**Consequence (the false green).** Node `{status:[b], phase:[p,q]}` is reachable;
terminal is `status=c`. The concrete view `{status=b, phase=q}` meets no terminal
and has no outgoing non-escape row (`fromP` requires `phase=p`) — a genuine dead
end the runtime reaches. Lint reports **no findings at all**: the single exit for
`phase=p` rescues the whole merged node.

**Test.** `TestAdvDeadEndExistentialOnMergedNode` (ADVERSARIAL) — asserts the
fixture premise (a node really is multi-valued on `phase`) before asserting
`graph-dead-end` is emitted.

**Status: FAILS.** Report is empty.

### Phase 3b scope note

No Phase 1 test or implementation file was modified; the adversarial file is
additive. The pre-existing `internal/graphlint` suite and every other package
(`internal/cli`, `internal/guard`, `internal/resolve`, `internal/table`) still
pass — only the three tests above fail.
