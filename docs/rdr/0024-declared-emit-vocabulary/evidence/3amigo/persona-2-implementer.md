Model: claude-opus-5[1m]

# Persona 2 — Implementer: clarification requests

Owned set: `0024:C1`–`C4`, `0024:D-identity`/`D-naming`/`D-selection-predicate`,
`edges[] source-anchor`. **Widened** to `§implementation-plan` (Phase 1/3/4),
`§testing-strategy`, `0024:A4`, `0024:A7`, `§mini-checks`, and to source under
`internal/table/` and `internal/cli/` — several first-hour questions are about
what the contracts do NOT say (a silence has no line range), and the Phase text
is where the answer would live if it existed.

---

## HIGH

### I1 — `[emit]` key ordering is unspecified where the whole model is emit-key-sorted
**Anchor**: `0024:C1`, `0024:D-identity`
`C1` fixes `[emit.<key>]` as a TOML table and `D-identity` fixes byte-exact key
identity, but neither says whether the CARRIED declaration set is key-sorted the
way every sibling normalized surface is. The repo's invariant is explicit:
`internal/table/normalize.go::emitSequence` sorts (`slices.Sorted(maps.Keys)`)
precisely so "map and authoring order never reach the normalized value"
(`sortedList` doc, same file), and `internal/table/model.go::Model.Rows` is
"pre-sorted by the identity tuple". A `map[string]EmitDecl` carrier satisfies
`C3`'s lossless clause and silently breaks the determinism `roundtrip_test.go::
TestReq122_KeyOrderRuleOrderAndRepeatedDumpDeterminism` exists to hold.
**Ask**: is the carried declaration set a sorted slice (like `Rows`/`Emit`) or a
map (like `Model.Tags`)? If a map, state that no reader iterates it in map order.
**Blocks**: the `C3` carrier type declaration — the first line of Phase 2 — and
whether the Phase 1 refusal for duplicate members can report deterministically.

### I2 — `C2` places the checks "beside `loadTags`" but the ordering constraint is unstated and one placement is wrong
**Anchor**: `0024:C2`
`C2` says the two checks "run as a step ahead of [`normalizeRules`], beside
`loadTags`", justified as "forced rather than preferred". The actual pipeline
(`internal/table/load.go::run`) is an eleven-step slice; the window between
`loadTags` and `normalizeRules` holds four other steps (`loadAccessors`,
`loadDump`, `loadContexts`, `loadInitial`, `loadTerminal`). `C2` fixes neither a
position within that window nor a rule for it, yet the window is not free:
`malformed_emit_declaration` (`C1`, a declaration-shape check) and
`unknown_emit_key`/`emit_value_out_of_domain` (`C2`, cross-checks of
`sourceRule.Emit` against declarations) are two different steps with a mandatory
order between them — checking rule values against a malformed declaration set is
undefined. `C1` and `C2` never say the declaration load precedes the rule
cross-check.
**Ask**: name the two steps and their order in `run`'s slice (e.g.
`loadEmitDecls` immediately after `loadTags`, `checkEmitUse` immediately before
`normalizeRules`), or state explicitly that the RDR delegates placement within
the window.
**Blocks**: the Phase 1 diff to `loader.run`, and Testing Strategy scenario 1's
premise that a malformed declaration refuses `malformed_emit_declaration`
rather than an incidental `unknown_emit_key`.

### I3 — `C1` mandates `ConformValue` reuse (`A4`) but `C3` forbids the type it takes
**Anchor**: `0024:C1`, `0024:C3`, `0024:A4`
`A4` verifies the reuse target as `ConformValue(decl TagDecl, member string)` and
concludes "no adapter or extraction is required". `C3` requires the carried
declaration be "a new model-level type, deliberately NOT `TagDecl`". The two are
only jointly satisfiable by constructing a throwaway `TagDecl{Kind:…, Domain:…}`
at each check site — which is exactly an adapter, and which quietly reintroduces
`TagDecl`'s unset `Min`/`Max`/`Elements`/`Provenance` into the emit path
(harmless today per `A4`, load-bearing if `TagDecl` ever gains a defaulted
field). Phase 1 hedges differently again: "reuse `conformKind`/`conformDomain`
through a narrow helper" — naming the two UNEXPORTED functions, not
`ConformValue`, i.e. a third shape.
**Ask**: fix one of (a) call `ConformValue` with a synthesized `TagDecl`,
(b) call `conformKind`/`conformDomain` directly in-package with a synthesized
`TagDecl`, or (c) add a sibling `conformEmitValue(EmitDecl, string)`. If (a) or
(b), say the `TagDecl` synthesis is licensed and does not violate `C3`.
**Blocks**: the Phase 1 conformance-helper signature; also whether `C3`'s
carrier can share `Domain []string` shape with `TagDecl` at all.

### I4 — `C4` join for a declared key the selected row does not emit is unspecified
**Anchor**: `0024:C4`, `0024:D-selection-predicate`
`C4` fixes `dispositions` as "emit key → the disposition token the declaration
assigns the SELECTED ROW'S authored value", and `D-selection-predicate` says
"when the selected row authors a value for a declared enum key". Neither states
the converse: a key declared under `[emit]` that the selected row does NOT emit.
The obvious reading (no entry) is not the only one — `C2` explicitly makes
"a declared key no rule emits" a non-finding "authoring headroom", so a reader
could take declarations as a projected key set. The `{}`-never-`null` clause
enumerates only whole-payload emptiness cases, not per-key absence.
**Ask**: state that `dispositions` keys are a subset of the selected row's
`emit` keys, never of the declared key set.
**Blocks**: the Phase 3 join loop (iterate `row.Emit` vs. iterate declarations)
and Testing Strategy scenario 5's "join matrix", which lists no such case.

---

## MEDIUM

### I5 — `C1` gives no source-schema shape for `[emit]`, and `A7`'s closure claim is still Pending
**Anchor**: `0024:C1`, `0024:A7`
`C1` states `domain` must be `any`-typed so `DisallowUnknownFields` "stops
descending at it", and enumerates five hand-written arms. It never writes the
`sourceDoc`/`sourceEmitDecl` struct, so an implementer must infer
`Emit map[string]sourceEmitDecl` with `Kind string` + `Domain any`. `A7` is
carried **Pending** on exactly the claim that decides whether those five arms are
complete ("the arms were derived from the type structure rather than enumerated
exhaustively"), and Prerequisites lists only A1–A5 as verified — A7 is not
mentioned there at all, so its Pending status has no stated gate.
**Ask**: give the `sourceDoc` field and sub-struct as `C1` gives `sourceRule.Emit`
its type-enforced contract at `internal/table/source.go:102`; and state whether
A7 must be discharged before Phase 1 merges or is discharged BY scenario 1.
**Blocks**: the Phase 1 schema edit and the scenario-1 fixture table's row count.

### I6 — Two behavior-defining doc sentences are not in Phase 4's amendment list
**Anchor**: `0024:C1`, Phase 4 (`§phase-4-example-and-docs`)
Phase 1 names three Go comments to amend and pre-clears two 0010 spec tests by
name. Phase 4 names only `docs/cli-output-contract.md` and the pricing example.
It misses `docs/model-authoring.md:641-644`, which states as a positive rule that
"[t]he reservation does not reach `[rule.emit]`. Emit keys are not tags —
**nothing declares them** … so `plan = "<clear>"` in an emit block loads and
answers with that literal text." Under a declaration, `<clear>` is
`emit_value_out_of_domain` unless listed as a member — the sentence reads false
by the same test Phase 1 applies to the three Go comments. The same file at
lines 757-758 ("Its keys are not tags — nothing declares them") repeats it.
**Ask**: add these to Phase 4's list, or state that they read as scoped to the
undeclared leg and need no edit.
**Blocks**: Phase 4 scope; and whether `<clear>` in a DECLARED emit key is a
refusal (I read `C1` as yes — it is just a member string — but the doc says
otherwise and no contract adjudicates).

### I7 — Phase 4's pricing-example declaration would refuse the shipped example
**Anchor**: `0024:C1`, `§illustrative-code`, Phase 4
Phase 4 says "declare the pricing example's emit keys". The example
(`models/examples/pricing-decision-table.toml`) authors `dpa = "required"` and
`dpa = "none"`; the RDR's own `§illustrative-code` shows
`[emit.dpa] domain = ["required", "waived"]` — `"none"` is not a member, so
transcribing the illustration refuses `emit_value_out_of_domain`. The section is
labelled "Illustrative only — shapes, not fixtures", so this is a trap rather
than a contradiction, but Phase 4 gives no domain and no instruction to read the
example's actual values.
**Ask**: confirm Phase 4 derives the domain from the four `[rule.emit]` blocks
(`plan ∈ {basic, pro}`, `dpa ∈ {required, none}`) and that the illustration is
not the fixture. Also: does the example get the flat form or the disposition
sub-table? A pricing model has no routes or stops, so the flat form is the only
sensible read, but Phase 4 does not say — and if it takes the flat form, the
example ships zero coverage of the disposition grammar.
**Blocks**: the Phase 4 edit and whether the example doubles as a disposition
witness.

### I8 — `emit_value_out_of_domain` under `kind = "int"` will fire from `conformKind` and `conformDomain` at two different truths
**Anchor**: `0024:C1`, `0024:A4`
`C1` says `int` "constrains the value to a base-10 integer literal" and `A4`
settles non-canonical literals permissively via `strconv.Atoi`. `ConformValue`
composes `conformKind` then `conformDomain`; for `Kind: "int"`, `conformDomain`
runs `strconv.Atoi` AGAIN and then tests `decl.Min`/`decl.Max`, which `C1` gives
an emit declaration no way to set. That is harmless-but-dead only if the
synthesized `TagDecl` leaves `Min`/`Max` nil — which is a property of the
synthesis site I3 leaves unfixed, not of any contract.
**Ask**: state that emit conformance never sets `Min`/`Max`, so the `int` arm is
kind-only. (Or, if only `conformKind` is called for `int`, say so — `C1`'s "not
an `int` literal" wording and `A4`'s `ConformValue` wording disagree on which
function is authoritative for `int`.)
**Blocks**: the Phase 1 helper body and the scenario-2 int fixture's expected
category.

### I9 — `C1` does not say whether a `[emit]` table with zero sub-tables is "zero declarations"
**Anchor**: `0024:C1`, `0024:C2`
`C2`'s opt-in predicate is "ZERO `[emit.*]` declarations" vs. "ONE OR MORE
declarations present" — a count of SUB-TABLES. A model authoring a bare `[emit]`
line and nothing under it decodes to an empty map. `C4` mentions "empty emit
block" in its `{}` enumeration, but that is about `[rule.emit]`, not `[emit]`.
The repo has a live precedent that this distinction is contract-level:
`sourceRule.Emit` is deliberately non-pointer because "`emit` keys on LENGTH
rather than on key presence" (`internal/table/source.go:99-102`), whereas
`Write`/`Clear`/`Gate`/`Escape` are pointers precisely so presence is keyed.
`C1` picks neither convention for `[emit]`.
**Ask**: state that `[emit]` with no sub-tables is the zero-declaration leg
(length, matching `sourceRule.Emit`) — or that it refuses.
**Blocks**: the `sourceDoc` field's pointer-ness and MVV step 4's re-run
(delete the table vs. empty it).

---

## LOW

### I10 — Phase 3's `NumField() != 14` → `15` edit lands on a test whose CURRENT assertion is 14
**Anchor**: Phase 3 (`§phase-3-envelope-surfacing`), `0024:C4`
Phase 3 says `TestReq39`'s "`NumField() != 14` becomes 15". At HEAD
`internal/cli/decision_table_0010_test.go:435` already reads `rt.NumField() != 14`
with a failure message "thirteen to fourteen" — so the constant matches, but the
adjacent message string does not, and the 13-key `slices.Equal` list must move to
14 in the same edit. Phase 3 names the key-list change and the NumField change
but not the message string, which will read false after the edit.
**Ask**: include the failure-message wording in the licensed diff, or state
that message strings are outside the enumerated set.
**Blocks**: nothing structural — but Phase 3's closing rule is "any assertion
outside this list going red is a signal the append did more than C4 licenses",
so an implementer needs to know a stale message is not such a signal.

### I11 — `D-naming` fixes three category slugs but not their Go constant names or `flow-model-invalid` mapping site
**Anchor**: `0024:D-naming`, `0024:C2`
`D-naming` gives the wire slugs; `C2` says they "map to `flow-model-invalid`
under `0005:C1`". The repo's constants live in `internal/table/category.go` under
a `Cat*` convention (`CatReservedTagKey`, `CatMalformedTagDeclaration`), and the
0005 mapping is a separate site. Neither is named.
**Ask**: confirm `CatMalformedEmitDeclaration` / `CatUnknownEmitKey` /
`CatEmitValueOutOfDomain` and that the `flow-model-invalid` mapping is
table-driven (so appending three entries is the whole change) rather than an
explicit switch.
**Blocks**: a 5-minute lookup, not a decision — listed because
`reserved_key_0008_test.go` and `roundtrip_test.go::
TestReq119_OneFixturePerCategoryAssertedByCategory` both enumerate categories by
constant and will need the three appended.

### I12 — `C4`'s text-mode rendering claims "no per-verb special case" but names a two-line format
**Anchor**: `0024:C4`
`C4` fixes text output as `dispositions.<key>: <token>` / `dispositions: (none)`
and attributes it to "the generic payload renderer". If the renderer is genuinely
generic over a `map[string]string`, that spelling is a consequence, not a choice
— and `emit` (the same Go type, the adjacent field) should already render the
same way. Worth confirming the two spellings match rather than being separately
authored.
**Ask**: confirm `emit` renders as `emit.<key>: <value>` / `emit: (none)` today,
so `dispositions` needs no renderer change at all.
**Blocks**: whether Phase 3 has a text-mode diff or zero text-mode diff.
