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
