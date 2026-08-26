Model: claude-opus-5[1m]

# Grounding sweep — cli/0011 (Stage 5, lens `grounding`, iter-1)

Scope: `Ground-sweep: clean (24 anchors)` at propose (commit `6dea213`) scopes this
sweep to claims added or edited after propose (`git diff 6dea213 HEAD`, 173 added
lines across refine `1c90ac4` and resolve `7d65cb5`), PLUS the unanchored prose
claims the propose ledger's anchor-only table never covered. All 49 `source-anchor`
edges project `resolved:true`; the anchor half needed no re-verification and is not
restated here.

Widened past the scoped Evidence spans into: Decision Rationale (L213), the
Load-Bearing Decisions list (L175, L178), Testing Strategy S5/S6 (L377, L383),
and Trade-offs/Consequences (L310) — each carries a codebase claim that names no
symbol and so mints no edge, which is where the propose ledger could not look.

## Findings

### F1 — REFUTED. `internal/resolve` does not have "a single non-test consumer".

**Claim** (Decision Rationale, L213): "The cost is a larger blast radius than the
Stage-2 estimate — recorded in the `large` Profile — bought with `internal/resolve`'s
single non-test consumer and P7's pre-release licence."

**Found on `main`**: 16 non-test files across 5 packages import
`intrastate/internal/resolve`:

```
internal/accessor/{binding,disposition,executor,model}.go
internal/cli/{flow_exec,flow_next,flow_resolve,flow_state}.go
internal/cli/flowbind/flowbind.go
internal/graphlint/reach.go
internal/guard/{grammar,lint,product}.go
internal/table/{load,model,normalize}.go
```

The narrower reading — consumers of the changing entry point `resolve.Resolve` — is
**three** non-test call sites, not one: `internal/cli/flow_next.go:281`,
`internal/cli/flow_resolve.go:106`, and `internal/cli/flow_resolve.go:217`.

This is the sole stated ground for accepting the enlarged blast radius, so the
claim is load-bearing for the Alternative-3 rejection and the `large` Profile
rationale. `TagSet` is unexported and `matches` is package-private, which bounds
the *signature* change to `internal/resolve` — that is the true and defensible
version of this claim, and it is not what the record says.

### F2 — REFUTED (incomplete). The second `matches()` call site is undisposed.

**Claim** (Technical Design, L80): "The kernel change is confined to the match seam:
… `Resolve`'s candidate loop excludes only on `no-match`, so `flow resolve`'s
behaviour is unchanged in every case where the CLI supplies a complete, unconflicted
view."

**Found on `main`**: `view.matches(row.Match)` has **two** non-test call sites in
`internal/resolve/resolve.go`:

- `:472` — the ordinary candidate loop in `Resolve`. The record disposes this one.
- `:614` — the **escape-rescue** loop in `escapeOrRefuse`:
  `if row.Outcome != in.Recognized || !view.matches(row.Match) { continue }`.
  No normative clause, design statement, or phase intent in the record disposes it.

A3's Evidence field *names* `escapeOrRefuse` as a caller of `matches` (L44), but
every prescriptive statement — C1, Technical Design, Phase 1's intent, S5 — scopes
the change to "`Resolve`'s candidate loop". C1's "no escape list" is a property of
the `flow next` **probe**, not a disposition of the kernel's escape path, and
`flow next`'s probe stripping `Escape` means `flow next` never reaches `:614` —
so the gap is invisible from the `flow next` side and lands entirely on
`flow resolve`.

The behaviour is genuinely open. Under C1 an `indeterminate` match verdict must not
exclude. Applied at `:614` an escape edge whose own match block names an absent or
conflicted key would newly become a rescue candidate, changing which refusals get
rescued and potentially converting a `no_match` refusal into an escaped plan —
which is a `flow resolve` semantics change beyond "unchanged in every case where
the view is complete". Left two-valued at `:614`, the kernel holds two different
match semantics at one seam, contradicting the "the verdict is a property of the
seam, not of who asks" ground the record uses to reject a separate
`resolve.Candidates` entry point (L177).

A shipped oracle fences this exact interaction and would adjudicate either way:
`internal/resolve/match_conflicted_test.go:195`
`TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues` — its escape row's Match
names only a non-conflicted key, so it passes today; a sibling case where the escape
row's own Match names the conflicted key is exactly the undisposed one, and is
untested.

### F3 — REFUTED. The `unresolved` → `unknown` field change breaks 6 shipped 0005
assertions, which A4 and C3 both record as unbroken.

**Claims**:
- A4 (L49): "Counts: PASSES-UNCHANGED 28, MOVES-UNDER-`--all` **0**, BREAKS **0**."
- C3 (L152): "The 0005 flow next tests in internal/cli/flow_next_0005_test.go MUST
  keep passing unchanged where they assert the alphabet, gate handling, reader
  narrowing, determinism, and non-mutation."
- S6 (L383): "all 28 oracles pass unchanged (A4 enumerated them; the MOVES and
  BREAKS classes are both empty)".

**Found on `main`**: the record decides the flat field is replaced, not supplemented
— Load-Bearing Decisions, L175: "the per-candidate field is `unknown`, a list of
`{key, reason}`. **Rejected: keeping the flat `[]string` `unresolved`**, which cannot
carry the remedy"; L178 rejects "a single opaque `unresolved` key list". Phase 2's
intent likewise carries gate ids and absent owned keys *into* `unknown`.

Six shipped assertions read the `unresolved` JSON key, and all six are gate-handling
oracles — precisely the class C3 says keeps passing unchanged:

```
internal/cli/flow_next_0005_test.go:103   if _, present := c["unresolved"]; !present
internal/cli/flow_next_0005_test.go:163   unresolved, ok := stringsAt(c, "unresolved")
internal/cli/flow_next_0005_test.go:166   c["unresolved"]
internal/cli/flow_next_0005_test.go:404   unresolved, _ := stringsAt(c, "unresolved")
internal/cli/flow_adversarial_0005_test.go:446  unresolved, _ := stringsAt(c, "unresolved")
internal/cli/flow_mvv_0005_test.go:66     unresolved, _ := stringsAt(c, "unresolved")
```

`internal/cli/flow_next.go:49` is `Unresolved []string \`json:"unresolved"\``, and the
helper is `stringsAt` — a `[]string` reader. `unknown` as `[{key, reason}]` fails
both the key lookup and the element type.

A4's classification is not wrong on its own terms — it classified the 28 oracles
against **C1's candidate predicate**, where MOVES/BREAKS are genuinely 0. The defect
is that C3 and S6 promote that predicate-scoped result into an unqualified "pass
unchanged", while a second, independent change in the same build (the payload field)
breaks six of them. These are mechanical re-homings, not predicate moves, so they
belong to neither of A4's classes and currently appear nowhere.

## Inverse check (new rule vs. existing sibling path)

C1 adds a new three-valued discriminator at the match seam. Searched for a sibling
path already making an "absent/uncomparable ⇒ undecided" decision:
**found, and the record already cites it as reuse rather than minting** —
`internal/resolve/guard.go::evaluateAtom` returns `GuardUnevaluable` with
`ReasonAbsent` (`guard.go:71`) / `ReasonUncomparable` (`guard.go:76`), closed by
`Reasons()` (`guard.go:82`), carried on `UndecidedAtom{Key, Block, Operator, Literal,
Reason}` (`guard.go:87`). The record's Existing Infrastructure Audit rows ("Reuse …
C1 reuses `absent` / `uncomparable` verbatim; mints no new vocabulary") and the
Approach's step-5 sibling-path note are CONFIRMED. No finding.

## Spot-confirmed (post-propose claims, no finding owed)

- `TagSet.matches` is the three-legged disjunction quoted, `resolve.go:174-181`.
- Zero occurrences of `conflicted` in `internal/cli` — confirmed, count 0.
- `RefusalKinds()` closes at five, `resolve.go:67-75`.
- `kernelTags` passes repeated `--tag` through undeduplicated, `flow_exec.go:673-679`
  — the conflicted case is user-reachable as A3 states.
- `assembledView` is last-write-wins `map[string]string` with no conflicted concept,
  `flow_next.go:232-241`.
- `summarize` walks `row.Atoms` with no `Block` test, `flow_next.go:183-188` — a
  match-block atom over a view-absent key lands in `unresolved` today (A2).
- `excluded` strips `probe.Match`/`probe.Escape`, binds `Recognized: row.Outcome`
  alongside `Outcomes: []string{row.Outcome}`, excludes on `KindNoMatch` alone,
  `flow_next.go:276-295` (A1's pairing requirement is real).
- `TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch` exists at
  `internal/resolve/guard_atoms_test.go:930` — S5's retirement target resolves.
- `models/rdr.toml`: 22 `[[rule]]`, 22 `[rule.match.stage]`, exactly 3
  `eq = "resolved"` (L212/222/232), `finalize-pass` guarded on `gate_passed` — A6
  confirmed, and the refined 21-of-22 wording is now accurate.
- `--all` is unregistered anywhere in `internal/cli`; `registerSelectionFlags`
  (`flow.go:92`) carries `--model`/`--flow`/`--artifact` only — C2's flag name is free.
- `internal/table/model.go::seamValue` canonicalizes set literals (`model.go:313`);
  match-block `in` expansion at `normalize.go:670` — A3's atom-independence holds.

## Dispositions (resolve half, same cycle)

| Finding | Disposition | Section touched |
| --- | --- | --- |
| F1 consumer count false | **fixed** | Decision Rationale — replaced "single non-test consumer" with the executed result (3 `Resolve` call sites; 16 non-test importers; unexported `TagSet`/`matches` as the real containment) |
| F2 escape site undisposed | **fixed** | C1 (new uniformity paragraph), Technical Design, Phase 1, S5, Consequences; new **A7 (Pending)** |
| F3 `unresolved`→`unknown` breaks 6 oracles | **fixed** | A4 (count qualified to the predicate), C3 (normative re-homing clause naming all 6), S6, Phase 4 |

Tiebreakers escalated: none. F2 was the load-bearing fork; it collapsed on the
record's own stated ground — "the verdict is a property of the seam, not of who
asks" (the basis for rejecting a separate `resolve.Candidates` entry point) —
plus JDR 0001 P1, plus the confirmed reachability of escape-row match blocks
(`internal/resolve/guard_fixtures_test.go:171-182` carries `Match: status=Draft`).
Leaving `:614` two-valued would reinstate the per-caller fork the record rejects.

Net-new scope charted: none. F2 is the same function's second call site, which C1
had to dispose of one way or the other to be implementable — not a new surface.

## Mini-checks fired (first lens pass owes the cue read)

`authority`, `oracle`, `disposition`, `trace` — tables written into the RDR under
Proposed Solution › Mini-checks. `fidelity` did not fire (no import/export,
parse/deparse, or serialize inverse in this contract). No CONTRADICTION row in the
desk trace.

## Needs (re)verification → Stage 6

- **A7 (new, Pending)** — the escape-site extension's blast radius on the shipped
  `internal/resolve` escape oracles. Method: Source Search + MVV Test. Classify
  `match_conflicted_test.go`, `escape_shape_0009_test.go`,
  `escape_shape_mvv_0009_test.go`, `guard_atoms_test.go` against the widened escape
  predicate, as A4 did for the 0005 `next` suite.
- **A4 (still Verified, now qualified)** — its counts are predicate-scoped; the
  payload-rename breakage is tracked in C3, not in A4's classes. No flip needed:
  the classification is correct on its own terms, the unqualified promotion in C3/S6
  was the defect and is fixed.
