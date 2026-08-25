
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

---

## Phase 3a — CoVe Verification

Independent chain-of-verification against the SPEC. Inputs were the record and
`artifacts/req-list.md` only; the Phase 1 and Phase 3b test files and
`coverage.md` were not read. Every witness below was authored from the REQ text
and then **executed** against the built binary (`intrastate lint --model <path>
--as=json`) or, where the RDR 0002 loader refuses the shape, against the engine
entry point `graphlint.Run(graphlint.NewRequest(m))`.

Probe models were written to a scratch directory and deleted; none is committed.

### REQs probed and how

**Authority / placement / CLI contract** — REQ-2 (`internal/resolve` imports no
`graphlint`), REQ-4/6/9 (`lint` registered at root, not under `flow`), REQ-5 (no
`os.Stdout`/`os.Stderr`/`fmt.Print` in `lint.go` or `graphlint`), REQ-7/8 (one
exported `Run`, one `NewRequest`, one call site), REQ-11 (`--flow` + `--model`
together → `flag-mutually-exclusive`, exit 2), REQ-13 (`model` field is
`[model].id`, not the path), REQ-14 (`--model <path>` instantiable), REQ-116 (no
waiver/ignore/allowlist surface), REQ-118 (no encode/decode surface).

**Invariant taxonomy** — REQ-32/44 (missing `[initial]` → `graph-dangling-edge`
naming `element: model`, never exit 0 on an empty reachable set), REQ-33/129
(terminal predicate over an observed tag → `graph-dangling-edge` carrying the
atom fields), REQ-34/103 (node whose only exit is an escape row is still a dead
end), REQ-35/36/125 (converging terminal: split node, terminated branch clean,
live branch takes `graph-dead-end`), REQ-40 (per-row syntactic single-valued
write, via the engine API), REQ-41 (two-path merge with a multi-valued node emits
nothing — no node-level check), REQ-42/113 (`[initial]` counts as a write; a row
preserves a tag it neither writes nor clears), REQ-67/68 (SC-17: owned
always-present key flagged, observed always-present key silent), REQ-69/70
(`[initial]` omitting the key → finding against `element: initial`).

**Coverage, product, withholding** — REQ-38/50/51/53 (union vs scoped product
over guard dimensions only), REQ-45/46 (`scalar` dimension →
`reason: dimension-not-finite`), REQ-48/49/123 (SC-11: `exists` over an optional
key contributes a `{absent}` dimension — product 4, two uncovered; an
always-present key contributes none — product 2), REQ-54/56/80/81 (withheld claim
→ `graph-unprovable-coverage`, `reason: row-can-refuse`, naming row and atom,
message phrased as a withholding; fires for an **escape** row too), REQ-57
(12 bool dims → product 4096 > published bound 2048 → `graph-product-too-large`),
REQ-58 (atom over a non-single-valued tag → `reason: tag-not-single-valued`),
REQ-59/60/126 (both constants published in `lint --help`: product bound 2048,
node ceiling 4096), REQ-115 (SC-12: group both bare-escape-closable and carrying
a refusing row resolves to blocking `graph-unprovable-coverage` with
`graph-coverage-closed-by-escape` **absent** — withholding dominates closure).

**Escape-row scoping** — REQ-28 (SC-9a: escape overlapping an ordinary row is
silent), REQ-29/30 (SC-9b: two escape rows sharing two classes → one
`graph-overlap` per shared class), REQ-61/62/63 (SC-19a/b: `ambiguous_match` arm
demanded only for a group carrying an ordinary `graph-overlap`, vacuously closed
otherwise), REQ-65/66 (SC-9c: bare escape closure reported as
`graph-coverage-closed-by-escape`, exit 0).

**Grouping** — REQ-19/20/21/22/23/24/52/124 (SC-20: two rows with different match
patterns whose unions would jointly close produce **two** groups and
`graph-coverage-gap` on **each**, exit 2; the paired same-match-pattern positive
control lints clean).

**Emission, ordering, envelope** — REQ-83/84/85 (SC-13: one run carrying
`graph-overlap` plus two independent `row-can-refuse` withholdings — no
first-failure short-circuit), REQ-86/89 (five consecutive runs byte-identical),
REQ-88 (a rule-namespace finding sorts before an element-namespace one inside the
same code bucket even when lexically greater), REQ-90/92 (aggregate
`graph-lint-failed`, exit 2, `findings` a top-level sibling of `code` with no
`{"type":"failed"}` wrapper), REQ-93/94/128 (SC-16: success emits
`{"type":"ok","data":{"findings":[]}}` — key present, list empty),
REQ-95/96 (`clierr` imports no `graphlint`; atom/class fields are `string`),
REQ-97/98/99 (text and JSON carry the identical finding-code set, same exit).

**Reachability** — REQ-101/104/106 (syntactic dataflow, no guard pruning),
REQ-109 (an ordinary self-loop terminates instantly), REQ-107 (the can-refuse
test never consults the relation).

**Advisory tier** — REQ-73/74/75/77/78/79 (SC-14: redundant row + unreachable
rule, both `info`, exit 0).

**CI gate** — REQ-119/121 (`.github/workflows/ci.yml` carries a `graph-lint` job
distinct from the golangci-lint `lint` job; it runs `make build` then
`./bin/intrastate lint --model models/rdr.toml --as=json` and asserts on the JSON
`code` field; `make check` carries the `build` edge), REQ-120 (inverting one
`gate_passed` guard in the checked-in model makes the gate fail with
`graph-lint-failed`), REQ-122 (SC-23 census: `models/rdr.toml` lints clean —
**zero** blocking findings).

### FAIL entries

Four actual violations. FAIL-2, FAIL-3, and FAIL-4 were derived independently
from the REQ text before the suite was run; they coincide with Phase 3b's ADV-1,
ADV-2, and ADV-3 above, which is corroboration from a separate starting point.
**FAIL-1 is not covered by any existing test.**

#### FAIL-1 — `graph-terminal-escape` is unreachable on any model declaring a terminal

- **REQ violated**: REQ-43 (INV 7). *"The input-observable defect this mints
  `graph-terminal-escape` for is a model that relies on such inference: a
  reachable non-terminal node with no outgoing non-escape row and no terminal
  declaration covering it (an implied terminal)."*

- **Failing input** (`p17`): a model declaring `terminal = ["fin"]` with
  `[context.fin.match.stage] eq = "x"`, root `stage=a, mark=m1`, and two rows out
  of the root — `r-x` writing `stage=x, mark=m1` and `r-y` writing
  `stage=y, mark=m2`. `stage=y` is reachable, satisfies no declared terminal, and
  is the source of no non-escape row: precisely INV 7's implied terminal.

  ```
  $ intrastate lint --as=json --model p17-implied-terminal.toml
  → findings: ['graph-dead-end']
  ```

- **Observed**: only `graph-dead-end`. `graph-terminal-escape` is never emitted.

- **Spec-required**: `graph-terminal-escape` naming the node that relies on the
  inferred terminal (alongside `graph-dead-end`, per REQ-83's complete-emission
  clause — the two invariants are separate obligations and one must not stand in
  for the other).

- **Mechanism**: `analysis.go::checkTerminalEscape` opens with
  `if len(a.model.Terminal) > 0 || len(a.nodes) == 0 { return }`. The code can
  therefore fire **only** on a model that declares no terminal *at all* — a
  global gate where the spec states a per-node condition. Confirmed by the
  positive control `p19` (same model with the `terminal` list and context
  removed), which does emit `graph-terminal-escape`. Consequence: a
  `graph-terminal-escape` fixture built the way the RDR prescribes ("a fixture is
  authored by omitting the declaration the model depends on" — omitting *one*
  terminal, not all of them) cannot mint the code, and REQ-MVV's illegal matrix
  entry for `graph-terminal-escape` can only be satisfied by the degenerate
  no-terminal-whatsoever model.

#### FAIL-2 — the successor join collapses edges reaching *different* successors

- **REQs violated**: REQ-101 (*"an edge is a normalized non-escape row … producing
  the node with **that row's** writes applied and clears removed"*), REQ-108
  (*"**two edges reaching the same successor** produce one node whose per-tag
  value sets are the union of theirs"*), and consequentially REQ-38 and REQ-78.

- **Failing input** (`p25`): root `stage=a, g=true`; two rows out of the root —
  `r-set` writing `stage=b, note=n1`, and `r-plain` writing `stage=b` only
  (`note` is optional and absent at the root). A third row `gap-row` matches
  `note=n1` and guards `g eq true`, leaving `g=false` uncovered over a fully
  declared finite domain.

  ```
  $ intrastate lint --as=json --model p25-false-green.toml
  {"type":"ok","data":{"findings":[{"code":"graph-unreachable-rule",…,"rule":"gap-row"}]}}   # exit 0
  ```

  Reachable set: `{g:[true] stage:[a]}`, `{g:[true] stage:[b]}`,
  `{g:[true] stage:[z]}` — **`note` appears on no node**, though `r-set` writes it.

- **Observed**: exit 0 with an empty blocking list. `gap-row` is declared
  `graph-unreachable-rule` (advisory), so `checkGroups` skips `checkCoverage` for
  its group and the gap is never proved.

- **Spec-required**: `graph-coverage-gap` for `gap-row`'s group, exit 2. The
  runtime reaches `{stage=b, note=n1, g=false}` via `r-set`, where `gap-row` is a
  candidate whose guard refuses.

- **Isolating control** (`p26`): the identical model with `r-plain` deleted — one
  edge instead of two, `gap-row` untouched — correctly emits
  `graph-coverage-gap`. Adding a second, unrelated edge out of the same source is
  what erases the finding.

- **Mechanism**: `reach.go::successorOf` folds every edge leaving a node into one
  successor unconditionally (`out = joinNodes(out, produced)`), and `joinNodes`
  drops any key absent on either side. The relation therefore
  *under*-approximates presence, contradicting REQ-106's over-approximation
  contract and REQ-110's premise that a merged node admits a *superset* of
  concrete views.

#### FAIL-3 — a withheld exhaustiveness claim suppresses the group's overlap finding

- **REQ violated**: REQ-84 (`0006:C16`). *"Withholding a group's exhaustiveness
  claim MUST NOT suppress overlap, coverage, or further withholding findings for
  that group."* Also REQ-25/27.

- **Failing input** (`p42`): one group whose two ordinary rows overlap on the
  fully-declared finite enum `color` — `r1` guards `color in ["red","green"]`,
  `r2` guards `color in ["green","blue"]`, so `color=green` enables both — and
  where each row *additionally* carries `[rule.guard.all.opt] eq = "x"` over an
  observed key declared optional.

  ```
  $ intrastate lint --as=json --model p42-overlap-withheld.toml
  → findings: graph-unprovable-coverage(row-can-refuse, r1)
              graph-unprovable-coverage(row-can-refuse, r2)
  ```

- **Observed**: no `graph-overlap`.

- **Spec-required**: `graph-overlap` naming `r1` and `r2`, in addition to the two
  withholding findings.

- **Isolating control** (`p43`): the identical model with only the two `opt` atoms
  removed — which changes the `color` overlap not at all — emits `graph-overlap`
  naming `r1` and `r2`. The overlap is decidable on `color` alone; declining to
  decide it because a *different* dimension is unprovable is the suppression
  REQ-84 forbids.

- **Mechanism**: `guard.AcceptedAssignments` returns a non-projectable set for a
  row whenever *any* of the group's dimensions is unprojectable, and
  `groups.go::emitOverlaps` then does `if !left.Projectable() { continue }`,
  skipping the pair silently.

#### FAIL-4 — dead end's outgoing-row test reads a merged node existentially

- **REQs violated**: REQ-34 (INV 2), REQ-111 (*"Universal checks must not [read
  merged nodes] … A ∀-claim evaluated against a widened domain gets easier to
  satisfy, which is the false-green direction"*), and INV 2's own claim that
  *"splitting is exact … no false green is introduced"*.

- **Failing input** (`p27`): root `stage=a, kind=k1`; two rows out of the root
  writing `stage=b, kind=k1` and `stage=b, kind=k2` respectively; one row `b-out`
  leaving `stage=b` that matches `kind eq k1`. Terminal is `stage=z`.

  ```
  $ intrastate lint --as=json --model p27-nonterm-key-merge.toml
  {"type":"ok","data":{"findings":[]}}   # exit 0
  ```

  Reachable set: `{kind:[k1] stage:[a]}`, `{kind:[k1 k2] stage:[b]}`,
  `{kind:[k1 k2] stage:[z]}`.

- **Observed**: exit 0, no findings.

- **Spec-required**: `graph-dead-end` for the concrete owned-state
  `{stage=b, kind=k2}`, which is reachable via the second root row, satisfies no
  declared terminal, and is the source of no non-escape row.

- **Mechanism**: `analysis.go::checkDeadEnd` splits on `terminalKeys()` only
  (REQ-37), then calls `hasOutgoingOrdinaryRow` → `matchSatisfiable` →
  `ownedAtomSatisfiable`, which accepts when *some* held value satisfies the atom.
  `kind` participates in no terminal, so the split never separates `k1` from `k2`
  and the single `kind=k1` exit rescues the whole merged node. The split recovers
  exactness for the *terminal* half of invariant 2 and leaves the *outgoing-row*
  half resting on the over-approximation the record says it does not rest on.

### Observations (no FAIL raised)

- **`graph-single-valued-state` is unreachable through the CLI.** RDR 0002's
  loader refuses a multi-value write on every declared kind (`enum`/`bool`/`int`
  refuse the member sequence; `set` and `scalar` refuse the `single_valued`
  marker), so no TOML the loader accepts can mint invariant 5's code. The check is
  correct and fires at the engine boundary (verified by constructing a
  `table.Model` in memory), and REQ-100 routes schema non-conformance upstream —
  but REQ-MVV requires the illegal fixture matrix to assert
  `graph-single-valued-state` through "the same command shape intended for CI"
  (ASSUMPTION-7), which cannot be done. Recorded for the deviations file rather
  than as a FAIL, since the engine honours REQ-40 and REQ-41 exactly.

- **Subsumption is routed to the advisory tier, not to overlap.** Two rows where
  one's accepted assignments are a proper subset of the other's do enable two
  ordinary rows at a concrete assignment (REQ-27's literal reading), but the
  record settles the tension explicitly: overlap is *"a partial intersection
  between two rows neither of which subsumes the other"*, and the disposition
  table gives the subset case exit 0 with `graph-redundant-row`. The
  implementation matches the record. No FAIL.

- **Advisory findings are dropped from the failure envelope.** On a blocking run
  only `report.Blocking()` reaches the wire, so a `graph-coverage-closed-by-escape`
  or `graph-redundant-row` decided in the same pass is not reported. REQ-91 scopes
  the failure envelope's `findings` to "the individual blocking findings", so this
  is conformant as written; noted because SC-19a's *"Both variants emit
  `graph-coverage-closed-by-escape`"* is observable only at the engine boundary
  for the failing variant.

- **`--as=<invalid>` renders its refusal in text mode.** `respond.ModeOf` falls
  back to `ModeText` for an unrecognised value, so `lint --as=yaml` emits
  `error: flag-invalid-value …` on stderr rather than a JSON envelope. Consistent
  with the shipped `respond` contract and outside RDR 0006's surface.

### Verdict

**BLOCK** — four FAIL entries (FAIL-1 … FAIL-4). FAIL-1 is new to this phase;
FAIL-2, FAIL-3, and FAIL-4 independently corroborate Phase 3b's ADV-1, ADV-2, and
ADV-3.

---

## Phase 3c — Fixup

Four defects addressed. FAIL-2/3/4 and ADV-1/2/3 are the same three found
independently, fixed once each and verified against both the 3a-described
symptom and the 3b failing test. Every fix was committed as its own green
increment.

Suite state at close: `go test -race ./...` green across every package,
`gofmt` clean, `go vet ./...` clean, `golangci-lint run` reports 0 issues,
and `intrastate lint --model models/rdr.toml --as=json` returns
`{"type":"ok","data":{"findings":[]}}` — REQ-122's census stays zero.

### FAIL-1 (REQ-43) — invariant 7 was decided per model, not per node

- **Fixed**: `internal/graphlint/analysis.go::checkTerminalEscape` opened
  with `if len(a.model.Terminal) > 0 { return }`, a global emptiness gate
  where REQ-43 states a per-node condition ("no terminal declaration
  **covering it**"). The code could fire only on a model declaring no
  terminal at all.
- **Change**: the check walks per-node terminal coverage, testing the SPLIT
  node for the same reason invariant 2 splits (REQ-36) — terminal
  satisfaction is universal over a node's per-tag value sets.
  `satisfiesSomeTerminal` is false for an empty terminal list, so the
  no-terminal-at-all case is unchanged.
- **Regression test** (this defect had none):
  `TestReq43_ImpliedTerminalIsMintedInAModelDeclaringOtherTerminals` in
  `internal/graphlint/invariants_0006_test.go`, driving FAIL-1's `p17`
  shape — `terminal = ["fin"]` covering `stage = x`, with `r-y` reaching an
  uncovered `stage = y`. **Verified to fail against the old gate** (the
  gate was temporarily restored and the test reproduced FAIL-1's exact
  symptom) and to pass against the fix. It also asserts `graph-dead-end` is
  emitted alongside per REQ-83, and carries a control that clears the
  finding by declaring the missing terminal.
- **State**: RESOLVED. Recorded as D14.

### FAIL-2 / ADV-1 (REQ-101, REQ-108) — the successor join collapsed divergent edges

- **Fixed**: `internal/graphlint/reach.go::successorOf` folded EVERY edge
  leaving a node into one successor, and `joinNodes` drops a key absent on
  either side, so two rows writing different owned keys annihilated each
  other. Confirmed as an under-approximation of presence, against REQ-106's
  over-approximation contract and REQ-110's superset premise.
- **Change**: `successorsOf` groups a source node's edges by their
  successor's PRESENCE FOOTPRINT before joining. REQ-108 licenses the join
  only between "two edges reaching the same successor", and REQ-101 makes
  "absent" a value of a tag's dimension, so the footprint is the successor
  identity and the per-tag value sets are what the widening unions over it.
  Edges writing different values to the same tag still converge on one node
  (REQ-108's merged fixpoint); edges establishing different tags stay
  distinct. The traversal remains a fixpoint over merged nodes and never
  path-sensitive, and the lattice stays finite so termination is unchanged.
- **Regression test**: `TestAdvSuccessorJoinCollapsesDivergentEdges` (3b,
  pre-existing) now passes both sub-assertions — the reachable set retains
  a node holding each divergent key, and no live row takes
  `graph-unreachable-rule`. `TestReq108_TraversalMergesConvergingEdgesIntoOneNode`
  and `TestReq41_MergedNodeWithTwoValuesIsNotASingleValuedViolation` pin
  the merging direction and both stay green.
- **Fixture corrected**:
  `TestReq110And112_OwnedSetBeforeMatchReadsMergedFixpointNodes` asserted
  `mids == 1` over nodes carrying `status = mid`, which encoded D9's
  functional-successor reading — the reading this defect disproves. It now
  counts nodes per owned-state identity, the actual
  merged-versus-path-sensitive distinction, and asserts the `opt`-absent
  node exists. The REQ-112 obligation it exists for is untouched: the
  `graph-owned-before-write` finding naming `reads-opt` still fires, which
  is the outcome the record's scenario 10 requires.
- **State**: RESOLVED. Recorded as D11, superseding D9.

### FAIL-3 / ADV-2 (REQ-84) — a withheld claim suppressed the group's overlap

- **Fixed**: `internal/graphlint/groups.go::emitOverlaps` skipped any row
  whose accepted-assignment set is not projectable, so one optional-key
  atom per row made a genuine finite-domain overlap vanish.
- **Change**: `decidableAccepted` re-expresses a non-projectable row by
  intersecting its projectable atom denotations (`guard.Denotation`) into
  the group's scoped product — RDR 0003's own move one level down, since
  `acceptedIn` already computes over the group's decidable sub-product.
  Where a row projects whole it is `guard.AcceptedAssignments` unchanged.
  Scoped to overlap and sound only because overlap is existential
  (REQ-110): dropping an undecidable conjunct widens the row, giving false
  positives at worst (REQ-117), never a missed defect. Deliberately NOT
  reused for coverage, which is universal; an `unless` block is left
  undecided rather than partially subtracted. `internal/guard` (RDR 0003's
  surface) was not modified.
- **Regression test**: `TestAdvWithholdingSuppressesOverlap` (3b,
  table-driven) — both arms now pass, the control and the attack, and the
  attack's finding names the pair `(r1, r2)` exactly once.
- **State**: RESOLVED. Recorded as D13.

### FAIL-4 / ADV-3 (REQ-34, REQ-36, REQ-111) — dead end's outgoing-row test

- **Investigated and NOT fixed as recommended.** The reported miss is real:
  a node multi-valued on a key participating in no terminal is rescued
  whole by an exit serving only one of its values. But the recommended fix
  conflicts with the spec, so per the brief the spec is followed and the
  reason recorded.
- **Why**: REQ-37 bounds the split as a MUST — "by the declared domains of
  terminal-participating keys only, **never the whole lattice**." The wider
  quantifier was implemented and MEASURED rather than argued away: ranging
  the outgoing-row test over match-participating keys manufactures views
  the merged node never correlated. In `models/rdr.toml` every
  `stage = "dropped"` write also writes `status = "abandoned"` (seven rows,
  verified), but the merge decorrelates the keys, so the cross product
  invents `{stage=dropped, status=draft}` and mints NINE false
  `graph-dead-end` findings on the model REQ-122 requires lint clean. The
  record anticipates this exactly: "a check against it can accuse a path
  the runtime never walks. A path-sensitive reading would be exponential
  and is rejected."
- **Not separable**: the ADV-3 node and the conforming model's node are
  structurally identical after the terminal split — one singleton terminal
  key beside one multi-valued non-terminal key — so no local rule
  distinguishes the true positive from the false positives. The separation
  needs the correlation tracking the record rejects.
- **Change**: none to the quantifier. `hasOutgoingOrdinaryRow` keeps
  REQ-37's bound, with the reasoning and the accepted miss documented at
  the call site so it is not silently rediscovered.
- **Fixture corrected**: `TestAdvDeadEndExistentialOnMergedNode` now PINS
  the accepted miss — any later widening of the split fails there and
  forces the REQ-122 census to be re-run before the widening is accepted —
  and adds a positive arm asserting the bounded split still catches a dead
  half living on a terminal-participating key, so the fixture keeps
  exercising invariant 2.
- **State**: DEFERRED to the author. Recorded as **D12** with
  `Status: needs author decision`. This is the only entry from this phase
  that needs one; the suite is green and the census is zero on the reading
  implemented.

### Scope note

Edits are confined to RDR 0006's surface: `internal/graphlint/`
(`reach.go`, `analysis.go`, `groups.go`) and its test files, plus the
artifacts. `internal/guard`, `internal/cli`, `internal/resolve`,
`internal/table`, and `models/rdr.toml` were not modified. No test was
weakened, skipped, or deleted: the two fixtures whose shape changed are
recorded as `TEST-FIXTURE` deviations with RDR evidence (D11, D12), each
keeps asserting the obligation it exists for, and D12's fixture gained an
assertion rather than losing one.

### Verdict

**OK** — all four defects dispositioned, suite green, `models/rdr.toml`
census zero. D12 carries `Status: needs author decision` and does not block.
