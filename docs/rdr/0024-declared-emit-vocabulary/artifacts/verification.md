# Verification — RDR 0024 Declared Emit Vocabulary

Phase 3a CoVe (chain-of-verification) artifact. Every input below was
designed from the RDR clause and the REQ list ALONE — the Phase 1 test
files, `coverage.md`, and `deviations.md` were not read — and then run
against the built binary (`make build`) or against `internal/table`'s
package API through a throwaway program outside the committed tree.

Fixtures were authored in a temp directory outside the repo and are
reproduced inline where a FAIL entry depends on them.

---

## What was probed

**C1 — the grammar's ten refusal arms (REQ-13…REQ-19, ASSUMPTION-1).**
One single-defect fixture per arm, each loaded on its own:

| arm | authored defect | observed |
| --- | --- | --- |
| unknown `kind` | `kind = "enumm"` | `malformed_emit_declaration` |
| no usable domain (a) | `kind = "enum"`, no `domain` key | `malformed_emit_declaration`, "carries no usable domain" |
| no usable domain (b) | empty `[emit.<k>.domain]` sub-table | same category, **same detail string** |
| no usable domain (c) | `domain = []` | same category, same detail string |
| empty-string member | `domain = ["alpha", ""]` | `malformed_emit_declaration` |
| duplicate member (flat) | `domain = ["alpha","alpha"]` | `malformed_emit_declaration` |
| duplicate member (across dispositions) | `route = ["alpha"]`, `stop = ["alpha"]` | `malformed_emit_declaration` |
| `domain` on non-enum (flat) | `kind = "bool"` + `domain = [...]` | `malformed_emit_declaration` |
| `domain` on non-enum (sub-table) | `kind = "scalar"` + `[emit.<k>.domain]` | `malformed_emit_declaration` (ASSUMPTION-8 leg holds) |
| empty disposition token | `"" = ["alpha"]` | `malformed_emit_declaration` |
| `domain` neither array nor table | `domain = "alpha"`, `domain = 42` | `malformed_emit_declaration` |
| disposition value not an array | `route = "alpha"` | `malformed_emit_declaration` |
| nesting below the disposition level | `[emit.<k>.domain.route]` with `sub = [...]` | `malformed_emit_declaration` |
| non-string member | `domain = ["alpha", 7]` | `malformed_emit_declaration` |

REQ-19/REQ-95 held on the discriminating point: the three no-usable-domain
authorings return one category AND one identical detail string, so no test
could accidentally pin a discrimination the decoder cannot make.

**C1's decoder-owned boundary (REQ-8, REQ-12, REQ-20).** A `domaim` typo
returns `unknown_schema_field`, not `malformed_emit_declaration`; a bare
`verdict = 42` returns `malformed_toml`; `single_valued = true` on an
`[emit.<key>]` declaration returns `unknown_schema_field`. No hand-written
arm re-reports any of the three. A declared emit key used in
`[rule.guard.all.<k>]` or `[rule.write]` still returns `unknown_tag` —
declaring an emit key does not promote it to a tag key.

**C1's permissiveness inheritance (REQ-9, ASSUMPTION-9).** Under
`kind = "int"`, `"03"`, `"+5"` and `"-0"` all load clean (exit 0);
`"4.2"`, `" 5"`, `""` and a 22-digit overflow all refuse
`emit_value_out_of_domain`. Under `kind = "bool"`, `"True"` refuses and
`"true"`/`"false"` pass. Both track `ConformValue`'s existing arms exactly.

**C2 — position, opt-in, and fail-fast (REQ-23…REQ-28, REQ-38).** With no
`[emit]` table and with a bare `[emit]` table, `Model.EmitDecls` is a
non-nil empty map and no refusal is reachable; `lint` exits 0 on the same
adversarial rule bodies that refuse once a declaration is added. A fixture
carrying a declaration defect AND two rule-side defects returns the
declaration defect, proving `loadEmitDecls` precedes `checkRuleEmit`
(REQ-25/REQ-55). Every refusal returned by `table.Load` was checked for
`Unwrap() []error`: none carries it, so the fail-fast, one-refusal-per-run
contract (REQ-38) holds.

**C2 — both rule classes (REQ-28).** An ordinary rule and an
`escape = ["no_match"]` rule each produce `unknown_emit_key` and
`emit_value_out_of_domain` from their own single-defect fixtures.

**C3 — the carrier (REQ-44…REQ-54, REQ-100, REQ-101).** A model declaring
five keys (a three-disposition partitioned enum, a flat enum, `bool`,
`int`, `scalar`) was loaded 200 times and its `Model.EmitDecls` compared by
`reflect.DeepEqual` against the first load. Result: **stable across all
200 loads**, with `Domain` for the partitioned key equal to the bytewise
union `[alpha bravo delta kilo mike yankee zulu]` regardless of the
disposition sub-table's map iteration order. `Domain` is `nil` (not an
empty non-nil slice) for `bool`, `int` and `scalar`; `Dispositions` is
`nil` for the flat-array spelling. `reflect.TypeOf(table.EmitDecl{})`
reports exactly three fields, all exported.

**C3 — the reuse's safe-by-omission dependency (REQ-53, REQ-55).**
`conformDomain`'s enum arm is guarded by `len(decl.Domain) > 0`, so an
empty domain would admit every value. C1's no-usable-domain arm makes that
unreachable, which is exactly the dependency REQ-55 names; the ordering in
`load.go::run` is what keeps it true.

**C4 — the payload join (REQ-57…REQ-68, REQ-102, REQ-103).** Against
`models/examples/routing-decision-table.toml` the JSON payload is:

    …,"gates":[],"emit":{"next":"park"},"dispositions":{"next":"stop"},"next":{},…

Fourteen wire keys, `dispositions` at index 9, immediately after `emit`.
Text mode renders `dispositions.next: stop` through the generic renderer
with no per-verb case. Against a five-key model whose selected row emits a
partitioned enum, a flat enum and a `scalar`, `dispositions` carries
**exactly one** entry — the partitioned key — proving the map is keyed off
the selected row's authored emit and is never padded to the declared key
set (REQ-60). An undeclared model, a flat-domain-only model and an empty
`[rule.emit]` block each render `"dispositions":{}` — present, never
`null`, never omitted (REQ-57/REQ-63). An escape-row rescue
(`escaped:true, escape_class:"no_match"`) carries the escape row's OWN
authored value's token, `{"verdict":"rescue"}` (REQ-66). `flow next` and
`flow next --all` payloads carry no `dispositions` key at all (REQ-67).

**S9 — the lexical negative control (REQ-7, REQ-108).** Two models
identical except that one declares `[emit.verdict] kind = "enum"` with
members containing a space and a semicolon (`"a b"`, `"c;d"`) plus
`[emit.n] kind = "int"` with values `"07"`/`"+5"`, and the other declares
nothing. `table.Dump` output is **byte-identical**. `renderEmit` still
bypasses `renderValue`; no value is parsed, canonicalized or converted.

**S6 — category registration and witness exclusivity (REQ-71, REQ-104,
REQ-105).** `table.Categories()` carries all three new slugs. Each of the
three checked-in `internal/table/testdata/neg/neg-emit-*.toml` witnesses
trips its own category and no other.

**Negative REQs (REQ-11, REQ-33, REQ-43, REQ-50, REQ-52, REQ-56, REQ-62).**
`internal/graphlint/` and `internal/resolve/` carry zero diff against
`main`; `internal/table` imports no `internal/guard`; no decoder-strictness
toggle or decoder-position plumbing was introduced; `roundtrip_test.go`
(`TestReq146`) is unedited; `*Model` gains no payload-shaped method
(`KernelTable` is pre-existing and carries no declarations).

**PH4 (REQ-85…REQ-88, REQ-107).** Both `models/examples/` tables lint
exit 0. `pricing-decision-table.toml` declares `plan ∈ {basic, pro}` and
`dpa ∈ {required, none}` — the values it actually authors, not the
Illustrative Code's `waived`. `docs/cli-output-contract.md` mentions
`dispositions` 13 times; the two stale `docs/model-authoring.md` sentences
("nothing declares", "loads and answers with that literal text") are gone.
The JDR 0002 §D1 note is present at the `Dispositions` field declaration.

`go build ./...`, `go vet ./...` and `go test ./...` are all clean.

---

## FAIL-1 — REQ-30, REQ-32, REQ-73 (and REQ-76's stated technique)

**The rule-side locator degrades to `:1` on three legal, decoder-accepted
authorings of the rule id.**

`0024:C2` is unqualified: "**All three categories MUST carry a source
line**, stamped through `internal/table/category.go::atLine`". `0024:MVV`
step 2 states the property as "each carrying the offending block's SOURCE
LINE, **not `:1`**", and `0024:F1` promises "offending file in `locator`".
REQ-32 fixes the technique: "The rule-side locator anchors on the rule's
`id = "<ruleID>"` line and **scans forward** to the emit block"; REQ-76
names it "a rule-id-anchored, **block-bounded forward scan**".

What is implemented (`internal/table/load.go::emitRuleLine`) is not a scan
but an exact whole-line string equality against the single literal
`id = "<ruleID>"`, delegated to `headerLine`, which returns 0 unless
**exactly one** line matches. `atLine` treats 0 as "no line" and leaves the
locator at the file's default `:1`.

Three inputs make a spec-conformant implementation visibly violate the
clause. All are valid TOML, decode successfully, and are refused for the
right category with the right message — only the locator is wrong.

**(a) Single-quoted (TOML literal-string) rule id.**

    [emit.v]
    kind = "enum"
    domain = ["a", "b"]

    [[rule]]
    id = 'r1'
    [rule.match.recognized]
    eq = "decide"
    [rule.emit]
    v = "zzz"

Command: `bin/intrastate lint --model <fixture>`
Observed: `emit_value_out_of_domain: rule r1 emits v = "zzz": "zzz" is
outside the declared domain (locator="<fixture>:1")`
Required: the locator must name the offending rule's source line (line 19
in this fixture), not `:1`.

**(b) Whitespace-padded assignment.** The same fixture with
`id   =   "r1"` — TOML permits arbitrary horizontal whitespace around `=`.
Observed: the identical refusal with `locator="<fixture>:1"`.
Required: line 19.

**(c) A rule id equal to the model id.** No load category forbids this —
`CatDuplicateRuleID` guarantees uniqueness among *rules*, and REQ-32's own
prose notes "`[model]` also carries an `id` key" as the reason to match the
VALUE. When the values collide, the value match finds two lines and
`headerLine`'s `matches != 1` guard returns 0:

    [model]
    id = "probe"
    …
    [emit.verdict]
    kind = "enum"
    domain = ["alpha", "beta"]

    [[rule]]
    id = "probe"
    [rule.match.recognized]
    eq = "decide"
    [rule.emit]
    verdict = "gamma"

Observed: `emit_value_out_of_domain: rule probe emits verdict = "gamma"…
(locator="<fixture>:1")` — line 4 and line 19 both match, so neither is
reported. Required: line 19 (the rule's id line). This is precisely the
`[model]`-collision hazard REQ-32 raises; matching the value narrows it but
does not close it.

The control confirms the arms are otherwise correct: the same defect with a
plainly-spelled, non-colliding id reports `locator="<fixture>:35"`.

**Scope note.** The DECLARATION-side locator (`emitHeaderLine`) shows an
analogous `:1` on the inline sub-table spelling
(`[emit]` … `verdict = { kind = "gauge" }`), but that is **not** a FAIL:
REQ-31 licenses reusing `tagHeaderLine`'s technique verbatim, and the same
weakness is reproducible at HEAD for `malformed_tag_declaration`
(`[tags]` … `recognized = { kind = "gauge", … }` also renders `:1`). It is
inherited behaviour the RDR explicitly adopts. `emitRuleLine` has no such
predecessor: it is new code, and REQ-32/REQ-76 describe a forward scan it
does not perform.

**Impact.** Latent rather than actively breaking: no model in
`models/` or `internal/table/testdata/` uses a single-quoted id, a padded
assignment, or a model/rule id collision. But the locator is a
consumer-authored-model surface, and the clause is unconditional.

---

*Phase 3a CoVe verification, run independently of the Phase 1 suite. This
file was staged concurrently with the Phase 3b adversarial test commit and
was swept into `18d5db5`; the verification work and FAIL-1 above are Phase
3a's, derived from the spec and confirmed against the built binary.*
