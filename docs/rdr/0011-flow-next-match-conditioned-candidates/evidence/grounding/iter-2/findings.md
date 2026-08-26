Model: claude-opus-5[1m]

# Grounding sweep — cli/0011 (Stage 5, lens `grounding`, iter-2)

Re-entry: iter-1 swept the pre-route-back draft, whose approach was a KERNEL
three-valued match seam. The critique route-back and re-propose (`fc93bfc`,
"CLI presence-split, kernel untouched") replaced that approach wholesale —
A7–A10 were dropped, A11–A14 minted — so iter-1's findings (F1 consumer count,
F2 escape site `:614`, F3 payload breakage) do not carry: F1's blast-radius
sentence is gone with Alternative 3, and F2 dissolves because the kernel is now
untouched (`internal/resolve` diff empty is an asserted end-state, S5/MVV 8).

Scope: `Ground-sweep: clean (58 anchors)` at the re-propose scopes this sweep to
claims added or edited after it (`git diff fc93bfc HEAD`, 63 insertions across
refine `2d3fb77` and resolve `2ca962b`), PLUS the unanchored prose claims the
anchor-only ledger never covers. All 78 `source-anchor` edges project
`resolved:true`; the anchor half needed no re-verification and is not restated.

Widened past the scoped Evidence spans into: C3's re-homing census (L158-178),
C2's DEPARTURE paragraph (L101-110), the Existing Infrastructure Audit (L295-307),
the Mini-check tables (L208-265), Prerequisites (L551-552), and Phase 3 (L178) —
each carries a codebase claim naming no symbol, which is where the propose ledger
could not look.

## Findings

### F1 — REFUTED. C3's "six TEST reads" is five, and its 4/1/1 breakdown is 3/1/1.

**Claim** (C3, L162-165): "six TEST reads on this build:
`internal/cli/flow_next_0005_test.go` (4, at the presence check and the gate-id
assertions), `internal/cli/flow_adversarial_0005_test.go` (1), and
`internal/cli/flow_mvv_0005_test.go` (1)."

**Executed on `main`** — every test occurrence of the literal `"unresolved"`:

```
internal/cli/flow_next_0005_test.go:103   if _, present := c["unresolved"]; !present
internal/cli/flow_next_0005_test.go:163   unresolved, ok := stringsAt(c, "unresolved")
internal/cli/flow_next_0005_test.go:166           c["unresolved"])
internal/cli/flow_next_0005_test.go:404   unresolved, _ := stringsAt(c, "unresolved")
internal/cli/flow_adversarial_0005_test.go:446  unresolved, _ := stringsAt(c, "unresolved")
internal/cli/flow_mvv_0005_test.go:66     unresolved, _ := stringsAt(c, "unresolved")
```

Six literals, but `:166` is **not a read**. It is the `%#v` argument inside the
`t.Fatalf` that reports the failure of the read one line above at `:163`:

```go
163:  unresolved, ok := stringsAt(c, "unresolved")
164:  if !ok {
165:      t.Fatalf("candidate `unresolved` is not an array of ids: %#v",
166:          c["unresolved"])
```

Distinct read sites = **5** (`flow_next_0005_test.go` 3, adversarial 1, mvv 1).

The descriptive half is wrong independently of the count. C3 places the four in
`flow_next_0005_test.go` "at the presence check and the gate-id assertions".
`:103` is the presence check and `:163` the gate-id assertion, but `:404` is an
**adversarial guard-fact** test over the key `flag` — its own comment reads "The
rows guarded on `flag` must report it as UNRESOLVED rather than being silently
decided either way" — not a gate id. So even under the literal-count reading, one
of the four is misattributed to the wrong oracle class.

This matters because C3 makes the number normative ("The build MUST mechanically
re-home every shipped read ... six TEST reads on this build") and S6/MVV 8 assert
against it ("only C3's six mechanical `unresolved`→`unknown` re-homings"). An
implementer re-homing 5 sites and finding the 6th is a `Fatalf` argument has to
decide whether the contract is unmet — the count should be the executed result.

### F2 — REFUTED. Three of the reads discard the comma-ok, not two.

**Claim** (C3, L174-177): "Two of the six discard the comma-ok of a helper that
returns `nil,false` on a non-string element and one of those is a NEGATIVE
assertion that passes vacuously on the empty slice."

**Executed on `main`** — the four `stringsAt(c, "unresolved")` reads:

```
internal/cli/flow_next_0005_test.go:163   unresolved, ok := …   ASSERTED (if !ok { t.Fatalf })
internal/cli/flow_next_0005_test.go:404   unresolved, _  := …   DISCARDED
internal/cli/flow_adversarial_0005_test.go:446  unresolved, _ := …  DISCARDED
internal/cli/flow_mvv_0005_test.go:66     unresolved, _  := …   DISCARDED
```

Discards = **3**, asserted = 1. `internal/cli/flow_mvv_0005_test.go:66` is the
one C3 does not account for.

The helper is `internal/cli/flow_harness_0005_test.go::stringsAt`, signature
`func stringsAt(m map[string]any, key string) ([]string, bool)`, returning
`nil, false` on a missing key, a non-array value, AND a non-string element —
CONFIRMED as C3 describes.

The negative-assertion half is CORRECT and correctly identified:
`internal/cli/flow_adversarial_0005_test.go:446` is the vacuous one —
`if containsString(unresolved, "flag")` never fires when `stringsAt` fails and
returns `nil`. The other two discards are positive assertions
(`if !slices.Contains(…)` at `:404`; a `gateListed` accumulator at mvv `:66`) and
would fail correctly on an empty slice. Only the count is wrong.

### F3 — REFUTED (incomplete). C3's production-rename enumeration omits the shipped help text.

**Claim** (C3, L165-170): "the production rename additionally touches the
`candidate.Unresolved` field and its `json:"unresolved"` tag, the field's doc
comment, the file header comment naming `unresolved`, and `summarize`'s
`slices.Contains(c.Unresolved, …)` uses — none of which are assertions and none
of which the six counts."

**Found on `main`** — `internal/cli/flow_next.go` also carries, unnamed:

- **`:70-71` — the cobra `Long` help string**, shipped user-facing CLI output:
  "Gates run only under --evaluate-gates. Without it their ids are reported as
  unresolved facts, and the command invokes no gate accessor at all." This is
  not a comment; it is the text `flow next --help` prints.
- `:173` — `Unresolved: []string{}`, the zero-value initializer in the
  `candidate` literal.
- `:188`, `:197`, `:206` — the three `c.Unresolved = append(c.Unresolved, …)`
  assignments (C3 names only the `slices.Contains` guards beside them).

The help-text omission is the substantive one, and it is a *contradiction*, not
just an under-count: C3's FIRST paragraph already mandates a full help rewrite
("The command's Short and Long help MUST state that flow next reports the
candidates the supplied state can take …"), and the shipped `Long` opens "next
ENUMERATES rather than selects" — the exact sentence C1 overrides. So the help is
doubly owed. The enumeration's purpose is to reassure that nothing outside the
counted assertions is load-bearing; a shipped string describing a key the command
would no longer emit is exactly what it should have caught.

### F4 — REFUTED. The Prerequisites checkbox contradicts the assumption statuses.

**Claim** (Implementation Plan › Prerequisites, L552):
"- [ ] A3, A11, A12, A13, A14 verified — the reachability invariant, the presence
agreement, the probe independence, the guard-payload read, and the `--all`
negative."

**Found in the record**: all ten Critical Assumptions project `Status: Verified`
(`rdr inspect --filter elements`: A1, A2, A3, A4, A5, A6, A11, A12, A13, A14 →
`Verified`), flipped by the resolve commit `2ca962b` ("resolve cli/0011 —
A3/A11/A12/A13/A14 verified"). The checkbox was never ticked. L551 above it is
`- [x] A1, A2, A4, A5, A6 verified`, so the two lines now disagree about the same
record's state. The Finalization Gate's Assumption Verification item reads
Status consistency explicitly, so this is a lock-blocking inconsistency left by
the resolve pass.

## Inverse check (new rule vs. existing sibling)

C1 adds a CLI-side presence discriminator. Searched `internal/cli` for a sibling
path already making a view-presence decision: **found, and the record already
cites it as reuse rather than minting** — `internal/cli/flow_next.go::summarize`'s
`if _, known := view[atom.Key]; known { continue }` (`flow_next.go:184`) over
`::assembledView` (`flow_next.go:232-241`) is the package's one presence test, and
the Approach's step-5 sibling-path note plus the Existing Infrastructure Audit row
("View presence test … Reuse") name it. No second presence test exists in
`internal/cli`. CONFIRMED, no finding.

The `unknown` reason vocabulary reuses `internal/resolve/guard.go`'s
`ReasonAbsent`/`ReasonUncomparable` (`guard.go:71`/`:76`, closed by `Reasons()`
at `:82`) rather than minting — the audit row says so. CONFIRMED. The one minted
token, `not-evaluated`, is scoped to the gate class the seam never sees, and the
record states the mint explicitly. No finding.

## Spot-confirmed (post-propose claims, no finding owed)

- `TagSet.matches` is the three-legged disjunction `!ok || tv.conflicted ||
  tv.value != w.Value`, `resolve.go:174-181`; an empty `want` returns `true`
  (A12's empty-match-pattern claim).
- `matches` has exactly two non-test call sites, `resolve.go:472` (`Resolve`) and
  `:614` (`escapeOrRefuse`) — both left untouched by this design, which is how the
  iter-1 F2 finding dissolves.
- `Resolve` returns on `if blocked != nil` (`:479-481`) ABOVE `switch len(selected)`
  (`:483`) — Alternative 4's rejection ground and the Key Discovery both hold.
- `gate` reads `row.Guard` via `evaluateAtoms(row.Guard, seam, view)` and
  `RequiresOwned` via `missingOwned`; it never reads `row.Match` (A12).
- `missingOwned` returns the `owned_state_unavailable` refusal BEFORE the
  `undecided` collection loop (`resolve.go` gate body) — A13's precedence limit
  and C1's stated fidelity limit are real.
- `Table.CheckValid` inspects only `len(row.Escape)`/`len(row.Writes)`,
  `resolve.go:791` — never `Match` (A12).
- `RefusalKinds()` closes at five, `resolve.go:67-75` (C1's "no sixth kind").
- `atomsFromBlock` refuses any match operator but `eq`/`in`,
  `normalize.go:78` ("a match block admits only eq and in") — A12.
- `outcomeBinding` lifts the `recognized` match atom into `Row.Outcome`,
  `normalize.go:467` — A5's exception and the disposition table's row.
- `Table.models` gates on `in.Recognized` at `resolve.go:458`, so C1's
  `Recognized: row.Outcome` pairing is load-bearing (A1).
- `summarize` builds `Unresolved` from all THREE sources C1 names: the atom walk
  with NO `Block` test (`flow_next.go:183-190`), `row.RequiresOwned`
  (`:192-199`), and un-run gate ids (`:200-209`). C2's "today" DEPARTURE half
  CONFIRMED.
- `table.Atom` exports `Block` (`model.go:151-156`) and `table.BlockMatch` is a
  re-exported alias (`model.go:31-39`), so C2's mandated `BlockMatch` filter is
  expressible in `internal/cli` with no new export.
- `excluded` is `func(row, owned, observed) bool`, nils `probe.Match` AND
  `probe.Escape`, and discards the `Result` — `flow_next.go:276-295`.
- `registerSelectionFlags` carries exactly `--model`/`--flow`/`--artifact`,
  `flow.go:92-97`; the literal `"all"` appears NOWHERE in `internal/cli`
  (`rg '"all"' internal/cli` → 0 matches) — C2's flag name is free and the trace
  witness holds as an absolute.
- `runFlowNext` fixes readers (`:97`), view (`:102`), and the `evaluate` bool
  (`:103`) before the row loop opens at `:114` — C1's once-per-invocation clause.
- `runFlowNext` skips `len(row.Escape) != 0` at `:115` — the Joint-check's
  "`flow next` never lists escape rows" claim, on which 0010's rescue reasoning
  rests.
- Named oracles all resolve: `TestReq49_ReasonIsAClosedNamedStringSetWithAnEnumerator`
  (`guard_atoms_test.go:768`), `TestReq47_ReasonSetStaysAtExactlyTwoMembers` (`:803`),
  `TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch` (`:930`),
  `TestReq69_NoPlanFlagShipsOnAnyVerb` (`flow_input_0005_test.go:452`),
  `TestReq12_FlowGroupDoesNotRedefineTheRootAsFlag` (`flow_surface_0005_test.go:369`),
  `TestReq3_EachVerbRegistersItsNormativeFlagSpellings` (`:128`).
- C3's three pins resolve: `checkAccessorBindings`' message "owned tag %s is
  served by %d readers; want exactly one" (`load.go:373`), `codeTagDuplicate =
  "flow-tag-duplicate"` (`flow_input.go:35`), `codeTagOwned = "flow-tag-owned"`
  (`:37`).
- `models/rdr.toml`: 22 `[[rule]]`, exactly 3 `eq = "resolved"` (L212/222/232) —
  A6 and the trace table's step-3 witness.
- `0010:C4` projects verbatim as quoted, including "`flow next`, `flow read-state`,
  and `flow set-state` are unchanged by this RDR".
- `docs/cli-output-contract.md` names `flow next` only at L126-127 as an
  invocation example with no payload fields — Phase 3's "currently shows only the
  invocation grammar" and Consequences' "pins the envelope, not per-command
  payload fields" both CONFIRMED.

## Noted, no finding owed

`candidate.Required` is `slices.Clone(row.RequiresOwned)` (`flow_next.go:172`) and
the same unresolved owned key is therefore reported in BOTH `required` and
`unresolved` today. This is shipped 0005 behaviour that C1 neither changes nor
relies on, so it is not a defect in this draft — recorded because a reader of
C1's "owned keys the view lacks" source may expect one field, and because the
`{key, reason}` migration touches the field that double-reports.

## Dispositions (resolve half, same cycle)

| Finding | Disposition | Section touched |
| --- | --- | --- |
| F1 "six TEST reads" is five; 4/1/1 is 3/1/1 | **fixed** | C3 ¶2 — count replaced with the executed result, each read given its file:line, and `:166` named as the `Fatalf` argument it is |
| F2 "two discard the comma-ok" is three | **fixed** | C3 ¶2 — THREE, all three named with file:line; the negative/vacuous identification kept (it was correct) |
| F3 production enumeration omits the shipped `Long` help | **fixed** | C3 ¶2 — added the cobra `Long` at `flow_next.go:70-71`, the `[]string{}` initializer, and the three `append` assignments |
| F4 Prerequisites checkbox contradicts CA statuses | **fixed** | Implementation Plan › Prerequisites — `- [ ]` → `- [x]` for A3/A11/A12/A13/A14 |

Amendment sweep (rdr-common §amendment-sweep) — the count "six" was load-bearing
in six further places; every one updated in the same pass:

- Mini-checks `trace` step 8 — "C3 (6 mechanical re-homings…)" → 5.
- Decision Rationale premortem — "beyond C3's six reads" → five.
- MVV step 8 — "only C3's six mechanical … re-homings" → five.
- Phase 3 — "re-home the 6 shipped reads" → 5, plus the production sites and the
  cobra `Long` help now named so the phase covers what C3 ¶2 enumerates.
- Testing Strategy S6 — "the 6 test reads" → 5, and "C3's six named re-homings"
  → five.
- **A4's Evidence** — "mechanically breaks 6 assertions in the same suite" → 5
  (with line numbers). A4's own classification is untouched: its counts are
  predicate-scoped (PASSES 28 / MOVES 0 / BREAKS 0) and remain correct; only the
  cross-reference to C3's independent payload edit carried the stale number.

Tiebreakers escalated: none. Every finding is a false count or an incomplete
enumeration inside a frame the evidence confirms — no design fork opened.

Net-new scope charted: none. All four findings correct claims the draft already
makes about its own build; none adds a surface.

## Mini-checks

The RDR's first lens pass (iter-1) fired `authority`, `oracle`, `disposition`,
`trace` and wrote those tables into the RDR under Proposed Solution ›
Mini-checks; they persist through the re-propose and are live in the current
draft. This pass's fixes add no cue — the corrections are to counts and an
enumeration inside contracts the tables already cover, and the `trace` table's
step 8 was updated in the sweep above rather than re-derived. No re-read owed.

## Needs (re)verification → Stage 6

- **None owed by these fixes.** All four corrections replace claims with executed
  results read off `main`; none adds a normative signature, wire/byte format,
  external-behaviour claim, or exactness word. No assumption flips: A4 stays
  `Verified` (its predicate-scoped classification is unchanged and correct), and
  the corrected cross-reference in its Evidence is now the executed number.
- Carried forward, unchanged by this pass: the record's ten Critical Assumptions
  are all `Verified`, and the Prerequisites checklist now says so.
