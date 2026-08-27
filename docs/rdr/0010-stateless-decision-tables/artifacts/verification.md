# Verification — RDR 0010 Stateless Decision Tables

Phase 3a Chain-of-Verification. Each entry below is an **observed** failure:
an adversarial witness was constructed, run against the built binary or the
production packages, and the behaviour recorded. Suspicions are not entries.

The verifier did not read the Phase 1 test files, `coverage.md`, or
`deviations.md`. Every probe drove `internal/table.Load`,
`internal/table.Dump`, `graphlint.Run`/`Reach`/`Fingerprint`, or the built
`intrastate` binary over hand-authored TOML.

---

## FAIL-1 | REQ-104, REQ-105, REQ-106 — the Phase 3/4 documentation obligations did not land

**REQs.**

- REQ-104 — Phase 3: "`emit` on `resolvePayload`, joined by rule id after
  selection, rendered in text mode; **the output contract doc gains the field
  and a decision-table invocation**." (IP)
- REQ-105 — Phase 4: "a generic decision-table fixture under
  `internal/*/testdata`, **`docs/cli-output-contract.md` gaining the `emit`
  field and a decision-table invocation**, and **the model authoring docs
  gaining the class**." (IP)
- REQ-106 — Phase 4 doc obligation: "it states that a decision table's
  discriminating dimensions are authored as `[rule.guard.all.<key>]` atoms and
  says why (match atoms scope the group and contribute no dimension,
  `0006:C7`), and it documents the escape-row \"otherwise\" idiom (A9)." (IP)

The record makes the author-facing half load-bearing in its own words: "Those
two are the author-facing half of the guarantee C5 enforces at lint — the
contract catches the mistake, the doc prevents it."

**Exact failing input.** Run against the branch at `433ba08`:

```
git diff --name-only main...HEAD -- docs/ | grep -v 0010-stateless-decision-tables
grep -n "emit" docs/cli-output-contract.md
grep -rn "decision-table" docs/*.md
```

**Observed.** The first command prints nothing: **no file under `docs/`
outside the RDR's own `0010-stateless-decision-tables/artifacts/` directory
was changed by this branch.** `docs/cli-output-contract.md` contains no
occurrence of the payload field `emit` — its seven `emit` hits are all the
English verb ("what that gateway emits", "emit a one-element `findings`"),
none a payload member. The repo's only two docs are
`docs/cli-output-contract.md` and `docs/README.md`, and neither contains the
string `decision-table` or `state-machine` anywhere, so no authoring doc gained
the class, the `[rule.guard.all.<key>]` guidance, or the escape-row
"otherwise" idiom.

Concretely, `docs/cli-output-contract.md` still documents the payload surface
without the new field — the `flow next` candidate list at line 164 enumerates
its members, and the `flow resolve` invocation block at line 136 shows only a
state-machine call with `--artifact state=./state.json --outcome advance`.
There is no decision-table invocation anywhere in the file.

**Expected per spec.** `docs/cli-output-contract.md` gains the `emit` payload
field and a decision-table invocation (REQ-104, REQ-105), and the model
authoring docs gain the class plus the two REQ-106 clauses. All three REQs are
unmet.

**Scope note.** This is a documentation gap only. The corresponding *code*
obligations of Phases 3 and 4 were verified present and correct: the `emit`
field is on `resolvePayload` immediately after `Gates`, it is joined by rule id
after selection, it renders in text mode through the generic renderer, and a
generic decision-table fixture exists. Only the doc half of REQ-104/105/106 is
missing.

---

## Probed and holding

The following were probed with adversarial witnesses and **passed**; they are
recorded so a later reader knows the surface was exercised rather than assumed.

**C1 — class declaration and strict decoding (REQ-1..13, 70, 89).** An unknown
`class = "bogus"` and a *present* `class = ""` both refuse
`malformed_model_declaration` naming both admitted values — the pointer-typed
`sourceModel.Class` correctly distinguishes absence from the empty string, so
the zero-value-reads-as-state-machine rule (REQ-7) is not reachable by
authoring. A `decision-table` model declaring two owned tags refuses with the
detail carrying the literal token `owned=2` (REQ-4). The one-directional
check holds: a zero-owned model with `class` omitted *and* one with
`class = "state-machine"` both load and reach lint, where a rootless machine
carrying owned scaffolding takes `graph-dangling-edge` at `element = model`
(REQ-5, REQ-53, REQ-67). The doubly-malformed determinacy case (REQ-11) — a
decision table declaring both an owned tag and `[initial]` — refuses on the
class with `owned=1`, never on the writer-arity diagnostic, confirming the
check sits before `checkAccessorBindings` (REQ-9, REQ-10).

**C2 — class-conditioned rule shape (REQ-14..18, 90).** A state-machine
ordinary rule with no write block still refuses `malformed_rule_shape`; the
same rule under `class = "decision-table"` loads with empty
`Writes`/`NextTags`/`RequiresOwned` (REQ-16). Every C2 prohibition is carried
by a *pre-existing* arm with no new refusal added (REQ-15): a write block on a
decision-table row refuses `write_to_non_owned_tag` (observed tag) or
`unknown_tag` (undeclared key), a `clear` list refuses
`write_to_non_owned_tag`, and `[initial]` over a non-owned tag refuses
`malformed_accessor_binding`.

**C3 — `[rule.emit]` grammar, normalization, dump (REQ-19..37, 91, 92).** An
integer, boolean, array, and a nested `[rule.emit.sub]` value each refuse
`malformed_toml` through the decoder's type error with no hand-written check
(REQ-20, REQ-21). `[rule.match.<emit-key>]` refuses `unknown_tag` (REQ-27). A
present-but-empty `[rule.emit]` loads and normalizes to the empty sequence
rather than refusing (REQ-28). `DumpColumns()` returns `emit` appended last
after `escape`, leaving every existing column at its index (REQ-29); the cell
renders key-sorted `key=value` pairs joined by `; ` and bracketed, with the
value as the raw authored string — `Zupper=A<B & C>D` and `ampers=x&y` pass
through unquoted and unescaped, confirming `renderValue` is bypassed (REQ-30)
— and an empty block renders `[]` (REQ-31). A `[dump]` `order` list omitting
`emit` refuses `malformed_dump_declaration: column emit is omitted` (REQ-32,
REQ-74). The expansion witness is decisive: a `[rule.match.a] in = ["x","y","z"]`
atom mints three rows and **every** row carries the authored block, with a
`reflect.ValueOf(...).Pointer()` comparison showing all three share one
backing array — the carry-through is real and no defensive clone was added
(REQ-24, REQ-25, REQ-92). A `flow resolve` landing on the *third* expanded row
(`--tag a=z`) returns the block. An emit-only edit leaves every
`graphlint.Fingerprint` byte-identical (REQ-37); inspection confirms
`Fingerprint` reads only `Atoms` and `NextTags`.

**C4 — resolve payload (REQ-38..51, 81, 98, 99).** JSON key order is observed
as `...,"outcome","rule","gates","emit","next","writes","clear","escaped"` —
`emit` sits immediately after `gates`, with `Gates` not `Rule` as its
predecessor and nothing else displaced (REQ-39, REQ-40, REQ-43). A row
authoring no block renders `"emit":{}`, never `null` (REQ-38). Keys serialize
in byte order (`0num` < `Zupper` < `_under` < `ampers` < `zlower`) with HTML
escaping disabled (REQ-35, REQ-36). An escaped plan carries the *escape row's
own* block alongside `"escaped":true,"escape_class":"no_match"` with no
escape-specific arm (REQ-44). Text mode emits `emit.alpha: first` one leaf per
line and `emit: (none)` for an unauthored block, entirely through the generic
renderer — `internal/cli/respond/text.go` is byte-unchanged on this branch
(REQ-45, REQ-46). `--outcome` is still required (REQ-48); an `--artifact`
bound to an uninvoked role is ignored, not refused (REQ-47); `flow next`
carries no `emit` and every candidate has empty `required` (REQ-50, REQ-82).

**C5 — ∅-rooted lint (REQ-52..69, 94..97).** `graphlint.Reach` over a
decision table returns exactly one node holding nothing; over a rootless
*state machine* it returns zero nodes — the seeding predicate augments rather
than replaces the `len(Initial)` test (REQ-52, REQ-53). Coverage runs from
that root: the three-cell partial table reports one `graph-coverage-gap`
naming the uncovered cell at exit 2, and the fourth ordinary rule closes it to
`[]` at exit 0 (REQ-78, REQ-79). The match-only control reports exactly
`graph-unprovable-coverage` carrying
`"reason":"no-participating-dimension"` (REQ-55, REQ-56, REQ-57, REQ-95) with
its own message stating the guard-atom remedy rather than
`unprovableMessage`'s "declare the domain" (REQ-61). The A13 double-report
probe is decisive: a zero-dimension group carrying a *bare* escape row yields
only `graph-unprovable-coverage` findings — `graph-coverage-closed-by-escape`
is absent, confirming the early `return` delivers precedence (REQ-59,
REQ-60). The class binding holds: a *state-machine* zero-dimension group lints
`[]` at exit 0, so the arm does not redden the existing corpus (REQ-62). The
stray-`terminal` control reports `graph-dangling-edge` at `element = terminal[0]`,
distinguished from the silenced root arm by element rather than code (REQ-65,
REQ-97). The two-outcome escape control reports the gap for `dt/review` only,
pinning per-outcome rescue scoping (REQ-96). The escape-"otherwise" variant
exits 0 with only the `graph-coverage-closed-by-escape` advisory (REQ-80). No
finding over any decision-table fixture names the class or the string
`decision-table` in `Code`, `Element`, `Message`, or `Detail` — the one
apparent hit is the pre-existing `"class":"no_match"` escape-class field, not
the model class (REQ-64).

**Negative surface (REQ-33, 34, 46, 51, 68, 85).** `models/rdr.toml`,
`internal/cli/respond/text.go`, `internal/resolve/resolve.go`, and
`internal/guard/product.go` are byte-unchanged against `main`. `rowByID` is
untouched apart from a comment reference (REQ-34). Exactly two of the four
`len(Initial) == 0` sites became class-aware: `reach` and `checkDanglingEdge`'s
root arm (`analysis.go:121`); `checkAlwaysPresentOwned` (`:241`) and
`checkUnreachableRules` (`:523`) keep the bare test (REQ-68).

**Fixture census (REQ-71, REQ-87(iii)).** 103 fixtures under
`internal/table/testdata/` carry a `[dump]` `order` list; 100 name `emit`. The
three that do not are `neg/neg-dump-unknown-column.toml`,
`neg/neg-dump-omitted-column.toml`, and `neg/neg-dump-repeated-column.toml` —
negative fixtures whose purpose is to be malformed dump declarations, so
correctly left alone.

**Regression sweep (REQ-100).** `make check` is green end to end: `go vet`
clean, `golangci-lint` 0 issues, `bin/intrastate lint --model models/rdr.toml
--as=json` returns `{"type":"ok","data":{"findings":[]}}`, and
`go test -race` passes every package.

---

## Undecidable

- **REQ-42** (gate deny never computes `emit`) — no decisive *runtime* witness
  was constructed: authoring a denying gate on a decision-table row within this
  session's fixtures produced load refusals rather than a gate-deny path. The
  clause is structurally satisfied by inspection — gates are evaluated and the
  refusal returns before `resolvePayload` is built — but that is inspection,
  not an executed witness, so it is recorded here rather than as a pass.
- **REQ-101** (no two-toolchain test) and **REQ-86..88** (the licensed-diff
  rule) are scope obligations over the test corpus, which this phase is
  forbidden from reading. Not assessed.


## Phase 3b — Adversarial review (independent)

Written without reading Phase 1's test files or Phase 3a's findings. Every
fixture below is authored TOML handed to the real `internal/table` loader
and driven through the real `graphlint.Run` / `table.Load`; nothing is
mocked. All four assertions were run against the current implementation and
the failing/passing state recorded is the observed one, not a prediction.

### ADV-1 — invariant 7 (declared-terminal / escape handling) fires over a decision-table model

- **Failure mode.** The ∅-root seeding admits a machine-only invariant it
  should refuse. `internal/graphlint/analysis.go::checkTerminalEscape` gates
  on `len(a.nodes) == 0` and on nothing else. Before this RDR that gate was
  true for every rootless model; after it, `reach` seeds the ∅ node for the
  decision-table class, so `len(a.nodes) == 1` and invariant 7 walks that
  node, finds no outgoing ordinary row, and reports
  `graph-terminal-escape` against a stateless table.
- **RDR anchor.** `Trade-offs / Failure Modes`, the **Silent (guarded)**
  bullet: the record guards the ∅-root regression in the direction where
  "the traversal seeds nothing, `checkGroups` skips every group" and leaves
  the opposite direction — seeding admitting too much — unguarded.
  `0010:C5` states the intended silence outright: "Declared-terminal
  handling (7) proper, and the escape arm of it, stay silent by
  construction: they key on a declared terminal, which C2 forbids."
  `checkTerminalEscape` does **not** key on a declared terminal, so that
  premise is false and the silence does not follow.
- **Aggravating.** The finding's remedy reads "declare it in the root
  `terminal` list" — the declaration `0010:C2` forbids this class. Following
  the finding's own advice yields `graph-dangling-edge` instead. The
  `Recovery` bullet names the only two real remedies ("drop `class`" or add
  owned state and an `[initial]`); the finding states neither.
- **Test added.** `internal/graphlint/adversarial_0010_test.go::TestAdv0010_TerminalEscapeStaysSilentOverADecisionTable`
  (escape-only decision table) and
  `…::TestAdv0010_TerminalEscapeSilentOnARuleFreeDecisionTable` (the
  smallest reproducer: a decision table with no rules at all — nothing to
  blame, no terminal to point at).
- **Currently fails:** YES, both.

### ADV-2 — invariant 2 (dead end) fires over a decision-table model

- **Failure mode.** The second half of the same over-admission.
  `checkDeadEnd` gates on `len(a.model.Terminal) == 0 || len(a.nodes) == 0`.
  `0010:C2` forbids a decision table from declaring `terminal`, but it
  enforces that **at lint** — "a terminal predicate over a non-owned tag is
  0006's `graph-dangling-edge` terminal arm, which C5 keeps live for exactly
  this reason" — so such a model **loads** and reaches the engine with a
  non-empty `Terminal`. With the ∅ node seeded the second disjunct is false
  too, and invariant 2 accuses the ∅ owned-state of being a dead end. The
  author reads two blocking findings for one authoring mistake, only one of
  which (`graph-dangling-edge`) `0010:C2` licenses.
- **RDR anchor.** `0010:C5`: "The missing-root arm of invariant 1
  (`0006:C18`), **dead end (2)**, always-present-owned (5),
  owned-set-before-match (6), and single-valued state are vacuous by
  construction over this class and **MUST NOT emit**." Also the `Approach`
  paragraph: "the machine-only invariants (root, dead end, always-present,
  owned-before-match, terminal handling) are vacuous by construction and
  stay silent."
- **Test added.** `internal/graphlint/adversarial_0010_test.go::TestAdv0010_DeadEndStaysSilentOverADecisionTable`.
- **Currently fails:** YES.

*Scope note.* The other invariants `0010:C5` names silent were probed and
**are** silent for the right reason, not by luck: `checkAlwaysPresentOwned`
and `checkUnreachableRules` keep the bare `len(Initial)` test (A10);
`checkOwnedBeforeMatch`'s `ownedReadSet` is empty with zero owned tags;
`checkSingleValuedState` has no write block to read; `checkNodeCeiling` is
complete at one node. Overlap (3), coverage (4), redundant-row and
vacuous-atom all run, and `graph-unprovable-coverage` /
`no-participating-dimension` was confirmed to fire only for the
decision-table class and never for a zero-dimension state-machine group.

### ADV-3 — the class/owned-set check pre-empts the undeclared-tag refusal C1 orders ahead of it

- **Failure mode.** `checkClassAgreement` was inserted immediately after
  `loadTags` — the **earliest** point in `0010:C1`'s window — which places
  it seven steps ahead of `normalizeRules`, where a rule's undeclared-tag
  refusal is raised. A model that both declares an owned tag under
  `class = "decision-table"` and matches an undeclared tag refuses
  `malformed_model_declaration`, hiding the `unknown_tag`.
- **RDR anchor.** `Trade-offs / Failure Modes`, the **Visible** bullet, is
  what the check delivers (`flow-model-invalid`, the class and the
  `owned=<n>` token — all confirmed present). What is not delivered is the
  order that visible failure is fixed against. `0010:C1`: "The class is read
  in `loadModelHeader` … and the agreement is checked in a step at or after
  `loadTags`, **so an undeclared-tag refusal precedes a class-disagreement
  refusal under `run`'s fail-fast order.**" C1 fixes the window at both
  ends; the floor's stated consequence and the ceiling's determinacy
  requirement are jointly satisfiable **only** by a step between
  `normalizeRules` and `checkAccessorBindings`. The chosen position
  satisfies the ceiling and breaks the floor's consequence.
- **Author-visible cost.** The reported category is the reverse of what C1
  tells an author to expect: they fix the class, reload, and meet a second
  refusal C1 promised would have come first.
- **Test added.** `internal/table/adversarial_0010_test.go::TestAdv0010_UndeclaredTagPrecedesTheClassDisagreement`.
- **Currently fails:** YES.
- **Companion pin (passes, deliberately kept).**
  `…::TestAdv0010_ClassDisagreementStillWinsOverAccessorBinding` asserts
  C1's ceiling — a decision table declaring both an owned tag and
  `[initial]` refuses on the class carrying `owned=1`, never on the
  writer-arity diagnostic `0010:C2` calls "survivable but not
  self-explanatory". It passes today and is kept so a fix for ADV-3 cannot
  trade the ceiling away by moving the check past `checkAccessorBindings`.

### Attacked and found sound (no test added)

Probed against the real loader/engine/CLI and **not** defective, so no
failing test could be written for them:

- `emit` on the resolve payload: position immediately after `gates`, `{}`
  and never `null` when unauthored, byte-ordered keys, text mode via
  `flatten`'s empty-container arm rendering `emit: (none)` with no per-verb
  special case (`0010:C4`).
- The escape-rescue join: a plan rescued by an escape row carries **that
  row's** own `emit` on the `Plan.RuleID` path, with `escaped: true` and
  `escape_class` intact.
- `emit` carried onto **every** row a multi-member match block expands to,
  and onto escape rows; excluded from `graphlint.Fingerprint` per the
  Load-Bearing Decisions.
- Strict decode: a non-string `emit` value (int, bool, array) and a nested
  `[rule.emit.sub]` all refuse `malformed_toml`, never `unknown_schema_field`
  (`0010:C3`); an adjacent unknown `[model]` key and a misspelled
  `[rule.emitt]` both refuse `unknown_schema_field`.
- A present-but-empty `[rule.emit]` normalizes to the empty sequence and
  dumps as `emit=[]` — no refusal, the deliberate divergence from the
  write/clear/gate blocks.
- The dump column: appended last after `escape`, `; `-separated like
  `writes` and `atoms`, values raw; an explicit `[dump]` list omitting only
  `emit` refuses `malformed_dump_declaration`.
- `class = ""` refuses `malformed_model_declaration`; a zero-owned
  `state-machine` is never refused at load; the `owned=<n>` count is exact.
- `flow next` carries no `emit` (`0010:BR4`, deliberate).

---

## Phase 3c — Fixup resolution

Every FAIL-N / ADV-N above is resolved. `go test ./...` and `make check`
are green end to end (`go vet` clean, `golangci-lint` 0 issues,
`bin/intrastate lint --model models/rdr.toml --as=json` returning
`{"type":"ok","data":{"findings":[]}}`, `go test -race` passing every
package).

### FAIL-1 (REQ-104, REQ-105, REQ-106) — RESOLVED

The documentation half landed. `docs/cli-output-contract.md` gains a
`flow resolve` / `emit` section (object of strings, byte-ordered keys,
`{}` and never `null`, positioned after `gates`, the escape-rescue join,
text-mode rendering via the generic renderer, and `flow next` carrying none
by design) and a decision-table invocation beside the existing
state-machine one.

`docs/model-authoring.md` is **new** — the repo had no authoring doc, only
an index inviting one — and carries the class and its one-directional
agreement rule, the `[rule.guard.all.<key>]` guidance with the `0006:C7`
rationale, the `graph-unprovable-coverage` /
`no-participating-dimension` fence that refuses the match-discriminated
silent green, and the escape-row "otherwise" idiom including `0002:C5`'s
per-outcome rescue scoping and the `graph-coverage-closed-by-escape`
advisory. `docs/README.md` gains its index entry. Recorded as **D15**,
which states the authoring choice and the options weighed.

Every code, `reason`, message, and payload shape quoted in either document
was taken from the built binary over a hand-authored table, not from
reading source.

### ADV-1 (`graph-terminal-escape` over a decision table) — RESOLVED

`checkTerminalEscape` returns early on `table.IsDecisionTable(a.model)`,
mirroring the two sites C5 already class-keys. C5's stated premise — "they
key on a declared terminal" — is false of this function, which keys on
`len(a.nodes) == 0` alone; the **obligation** ("stay silent") is
implemented, the false premise is not honoured. Recorded as **D13**
(SPEC-DEFECT). Both red tests
(`TestAdv0010_TerminalEscapeStaysSilentOverADecisionTable`,
`TestAdv0010_TerminalEscapeSilentOnARuleFreeDecisionTable`) are green.

### ADV-2 (`graph-dead-end` over a decision table) — RESOLVED

`checkDeadEnd` gains the same class-keyed early return. Its
`len(a.model.Terminal) == 0` disjunct is defeated because C2 enforces the
`terminal` prohibition at lint rather than at load, so such a model reaches
the engine with a non-empty `Terminal`. Recorded in **D13** alongside
ADV-1. `TestAdv0010_DeadEndStaysSilentOverADecisionTable` is green.

**A10's partition is intact:** neither edit touches a `len(Initial) == 0`
site, so REQ-68's "exactly two of four became class-aware" still holds and
`checkAlwaysPresentOwned` / `checkUnreachableRules` keep the bare root test.

### ADV-3 (`checkClassAgreement` ordering) — RESOLVED

The step moved from immediately after `loadTags` to between
`normalizeRules` and `checkAccessorBindings` — the one position at which
C1's floor consequence ("an undeclared-tag refusal precedes a
class-disagreement refusal") and its ceiling (the doubly-malformed case
refuses on the class) are jointly satisfied. REQ-12 licenses the move
without an override.

`TestReq9_AnUndeclaredTagRefusalPrecedesTheClassDisagreement` had asserted
the opposite of its own name and of the C1 sentence in its own doc comment;
it now pins `CatUnknownTag`, which is a tightening. Recorded as **D14**.

**Ceiling pin confirmed green:**
`TestAdv0010_ClassDisagreementStillWinsOverAccessorBinding` and
`TestReq10_DoublyMalformedRefusesOnTheClassNotOnWriterArity` both passed
before the move and both still pass, so the floor was not fixed by trading
the ceiling away.

### Collateral

`internal/cli/flow_next_0011_test.go`'s cross-RDR diff guard fired on the
two new doc paths. Its 0010 allow-list block is extended the way D7/D9
already established — named seams appended, never a relaxed check, so a
file outside both RDRs' lists still fails. Recorded as **D16**.
