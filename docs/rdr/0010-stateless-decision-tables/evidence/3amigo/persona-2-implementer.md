Model: claude-opus-5[1m]

# 3amigo — Persona 2 (Implementer), RDR 0010

Owned set read: `C1`–`C5`, `D-identity`, `D-wire-byte-format`, `D-naming`,
`D-selection-predicate`, and every `source-anchor` edge those land on. Source
read on `main` at `/Users/cwensel/sandbox/newcoinc/intrastate`:
`internal/table/load.go`, `internal/table/model.go`, `internal/table/dump.go`,
`internal/table/normalize.go`, `internal/table/source.go`,
`internal/graphlint/reach.go`, `internal/graphlint/analysis.go`,
`internal/graphlint/engine.go`, `internal/cli/flow_resolve.go`,
`internal/cli/flow_next.go`, `internal/cli/respond/text.go`.

---

### P2-1 — C4's `--as=text` emit rendering contradicts the respond-owned text renderer  [High]
**Anchor**: `0010:C4`, `0010:D-wire-byte-format`, `0010:S5`
**Trigger**: C4 — "Under `--as=text` the plan MUST render each emit pair as one
`key=value` line." S5's expected column adds "text renders one `key=value` line
per pair, and **no line at all** for an unauthored block."
**Question**: which code writes that line? `internal/cli/flow_resolve.go` never
prints its own payload — it calls `respond.OK`, and
`internal/cli/respond/text.go::writeTextPayload` is the sole renderer. That
renderer is deliberately payload-agnostic: it marshals the payload to JSON,
re-decodes, and emits one path-qualified leaf per line via `flatten`, whose
doc-comment states the design constraint outright ("a verb never prints its
payload itself … a hand-written text template satisfies neither for long").
Over `resolvePayload.Emit = {"next":"propose"}` that renderer produces
`emit.next: propose`, **not** `next=propose`. And `flatten` explicitly emits
`label(path) + "(none)"` for an empty map — "Absence and emptiness are
different claims in this contract … a renderer that dropped empty containers
would erase exactly the distinction REQ-66 turns on" — so an unauthored emit
renders `emit: (none)`, the exact opposite of S5's "no line at all". So C4
either (a) mandates a per-verb text special case that 0005's respond design
forbids, (b) means `emit.<key>: <value>` and the phrase `key=value` is loose
prose, or (c) intends an override of 0005's text-rendering rule that
`Overrides` does not name.
**Blocks**: whether Phase 3 touches `internal/cli/respond/text.go` at all (a
cross-RDR change to 0005's renderer) or only adds a struct field to
`resolvePayload` and accepts the generic projection. It also blocks writing the
S5 assertion — the two readings assert different bytes, and one of them fails
against a correct implementation.

---

### P2-2 — C1's agreement check has no legal home in the loader's step order  [High]
**Anchor**: `0010:C1`, `0010:§authority` (the "model class" row)
**Trigger**: the Authority table names `loadModelHeader` as the *canonical
writer* of "model class (C1 agreement check)"; C1 requires the check to compare
`class` against "the declared owned tag set".
**Question**: `internal/table/load.go::run` fixes the step order as
`loadModelHeader, loadOutcomes, loadTags, loadAccessors, loadDump,
loadContexts, loadInitial, loadTerminal, normalizeRules,
checkAccessorBindings` — `loadModelHeader` runs **before** `loadTags`, so
`l.model.Tags` is empty when it runs and the owned count is unavailable. Three
options and the RDR picks none: split the check into a new step after
`loadTags`; move the class read into `loadTags`; or have `loadModelHeader`
reach into `l.doc.Tags` (the undecoded source) rather than `l.model.Tags`. The
choice is observable, because it fixes *refusal precedence*: does a model with
both an unknown `class` value and an undeclared tag report
`malformed_model_declaration` or `unknown_tag`? Scenario 1 asserts a specific
category on the disagreement cases but says nothing about interleaving with the
other load failures, and `run` is fail-fast on the first error.
**Blocks**: where the Phase 1 check lands in `run`'s step list, and what the
loader table tests in Scenario 1 may assume about which refusal wins.

---

### P2-3 — C1's "a state-machine model MUST declare at least one owned tag" is a new refusal on existing models, and the RDR states the opposite  [High]
**Anchor**: `0010:C1`, `0010:§consequences`
**Trigger**: C1 — "a `\"state-machine\"` model MUST declare at least one
[owned tag]; disagreement in either direction is a `malformed model
declaration` load failure". Consequences — "every existing model loads
unchanged"; C1's own `⇒` — "every existing model loads unchanged".
**Question**: those two cannot both hold. Since absent `class` reads as
`state-machine`, C1 makes *zero owned tags with no `class` key* a load failure
today — a model that currently loads fine and merely lints
`graph-dangling-edge` (which is precisely MVV step 6's and Scenario 4's
"class-omitted control": "the same model with `class` omitted → lint
`graph-dangling-edge` naming the absent `[initial]`"). If C1's agreement check
fires at **load**, that control never reaches lint — it refuses with
`malformed_model_declaration` and the 0006:C18-unchanged control is
unwritable as specified. I checked the fixture corpus: all 103 `[model]`-
carrying TOML files under `internal/` and `models/` declare at least one
`provenance = "owned"` tag, so no shipped fixture breaks — but the RDR's own
negative control does. Which is it: does the state-machine direction of the
agreement check exist, and if so is the class-omitted control a load refusal
rather than a `graph-dangling-edge`?
**Blocks**: whether C1's check is one-directional (decision-table ⇒ zero owned)
or bidirectional; and the expected output of MVV step 6 / Scenario 4's
class-omitted control, which currently asserts a lint finding for a model the
loader would refuse.

---

### P2-4 — `Row.Emit`'s type is under-specified against the existing `TagValue` shape  [Medium]
**Anchor**: `0010:C3`, `0010:§technical-design` (item 1), `0010:D-wire-byte-format`
**Trigger**: Technical Design — "`Row` gains `Emit []TagValue`-shaped literals
sorted by key"; C3 — "a flat TOML table whose values MUST be strings …
key-sorted, duplicate-free sequence".
**Question**: `internal/table/model.go::TagValue` is `{Key string; Value
[]string}` — its value is a *member sequence*, not a string. "`[]TagValue`-
shaped" is either a literal instruction (in which case every emit value is a
one-member slice, and the dump's `renderValue` would bracket-and-quote it via
the `len(members) != 1` arm, and `resolvePayload.Emit` needs an unwrap to
`map[string]string`) or a loose analogy for a new
`type EmitValue struct{ Key, Value string }`. Reusing `TagValue` also drags in
`cloneTagValues`, `renderTagValues`, and the `isClear` sentinel path — none of
which C3 wants, since emit keys are explicitly not tag keys. Relatedly: `expand`
builds each expanded row's `NextTags`/`Writes` from the rendered assignments
and never copies from `base`, so a new `Emit` field must be added to *both* the
`rows := []Row{{…}}` seed literal and/or the per-row loop; the RDR does not say
which, and the seed literal currently drops any field not enumerated there.
**Blocks**: the Phase 1 type declaration and every downstream signature
(`column()`'s new `case "emit"`, `resolvePayload`'s field type, the
normalization sort/dedup helper). Also blocks whether `dump.go::renderValue`
is reused for emit (bracketing/quoting a single string) or emit gets its own
plain renderer.

---

### P2-5 — C3's dedup rule is unreachable in a TOML decoder, so "duplicate-free" pins nothing testable  [Medium]
**Anchor**: `0010:C3`, `0010:S3`
**Trigger**: C3 — normalization "MUST carry the block on the row as a
key-sorted, **duplicate-free** sequence"; Scenario 3 lists "duplicate keys" as
a case with expected "key-sorted duplicate-free sequence".
**Question**: `[rule.emit]` decodes as a TOML table into a Go map (following
`sourceModel.Metadata`'s `map[string]any` precedent), and `pelletier/go-toml/v2`
refuses a duplicate key in a table before this code sees it. So what does the
Scenario 3 "duplicate keys" test actually author, and what category does it
expect — an upstream decode error, `CatUnknownSchemaField`, or a new emit
category? A8 already warns that upstream message text is not assertable
(`0002:C24`). If the answer is "duplicates are impossible", C3's dedup clause is
dead text and the test case must be removed rather than written to an invented
refusal.
**Blocks**: whether Phase 1 writes a dedup pass at all, and what Scenario 3's
duplicate-key row asserts.

---

### P2-6 — C3 requires a non-string emit value to refuse, but names no category  [Medium]
**Anchor**: `0010:C3`, `0010:S3`, `0010:§failure-modes`
**Trigger**: C3 — "a flat TOML table whose values MUST be strings"; Scenario 3
— "a non-string value … **non-string refuses**". Failure Modes enumerates the
visible refusals (`malformed model declaration`, `write to non-owned tag` /
`unknown tag`, `malformed dump declaration`) and the non-string emit value is
not among them.
**Question**: which `Category` from `internal/table/category.go`? If `[rule.emit]`
is typed `map[string]string` in `sourceRule`, `decodeStrict` produces an
upstream type error whose category is whatever `source.go` maps decode failures
to — message text upstream-owned per A8, so the test can only assert the
category. If it is typed `map[string]any` (the `Metadata` precedent the
Infrastructure Audit cites) the refusal must be hand-written and needs a
category: a new arm on `CatMalformedRuleShape`? `CatMalformedRuleShape` is
0002's, and `0010:C2` already conditions one of its arms. Also unstated: is a
*nested* table (`[rule.emit.sub]`) a non-string value or an unknown schema
field?
**Blocks**: the `sourceRule` field type in `internal/table/source.go`, and the
category Scenario 3's non-string case asserts.

---

### P2-7 — C5 says "root the reachability relation at the empty owned-state node" but `reach`'s guard is `len(m.Initial) == 0`, and the class is not on `table.Model`  [Medium]
**Anchor**: `0010:C5`, `0010:§technical-design` (item 2), `0010:§authority` (reachability-root row)
**Trigger**: C5 — "Graph lint over a `\"decision-table\"` model MUST root the
reachability relation at the empty owned-state node"; Technical Design —
"`reach` seeds the traversal at the ∅ node when the model's class is
decision-table; `checkDanglingEdge`'s missing-root arm keys on class, not on
`len(Initial)`"; the Authority table adds "Every site that today asks
`len(Initial) == 0` … routes through one loaded-model accessor for the class,
so there is exactly one place the question is answered."
**Question**: what is that accessor's signature and where does it live?
`internal/graphlint` reads `*table.Model`, whose struct (`model.go:366`) has no
class field and no accessor. So Phase 1 must add both a field and a method
(`func (m *Model) Class() Class`? a `Class` string type? exported constants?) —
none named. Two sub-questions that decide the shape: (a) `graphlint.Run` takes
a `*table.Model` and returns `Report{}` for nil, but nothing forces that model
through the loader — a caller constructing a `table.Model` literal gets the
zero value for the new field, and the zero value must therefore be
`state-machine` (an empty-string default) rather than an unset third state,
or lint over a hand-built model silently flips class; (b) `reach`'s guard is
currently `m == nil || len(m.Initial) == 0` — under C5 a *decision-table* model
skips the guard, but what about a state-machine model that reached lint with an
empty `Initial`? Its `graph-dangling-edge` root arm still fires
(`0006:C18` untouched), but `reach` must still return `nil, true` for it, so
the new guard is `class == decision-table || len(Initial) > 0`, not a
replacement of the `len(Initial)` test — the RDR's "keys on class, not on
`len(Initial)`" reads as a replacement.
**Blocks**: the accessor's name/type/zero-value on `table.Model`, and the exact
new predicate in `reach` and in `checkDanglingEdge`'s root arm — including
whether the two predicates are the same expression or differ.

---

### P2-8 — C2 leans on `[initial]` being refused for a zero-owned model, but the loader's `[initial]` path refuses on a *different* rule than C2 names  [Medium]
**Anchor**: `0010:C2`, `0010:A10`, `0010:A12`
**Trigger**: C2 — "an `[initial]` key or a write/clear key that is not a
declared owned tag refuses as `unknown tag` (undeclared) or, once declared and
non-owned, as `malformed accessor binding` (an `[initial]` key is a written tag
demanding exactly one writer, which a decision table declares none of)".
**Question**: I traced it and it holds, but by a longer route than C2 states,
and the route matters for the negative-control test. `loadInitial`
(`load.go:540`) refuses only `CatUnknownTag` for an undeclared key; for a
*declared observed* key it passes (its own comment defers: "Every [initial] key
MUST be a declared OWNED tag. A non-owned key is caught by the writer
binding"). The actual refusal comes from `checkAccessorBindings`
(`load.go:360`) — the **last** step in `run` — via
`written[t.Key] = true` for every `Initial` entry, then
`writerCount[key] != 1` ⇒ `CatMalformedAccessorBinding "written tag %s is
served by 0 writers; want exactly one"`. That message names *writers*, not
`[initial]`, and it fires only because the decision table declares no writer.
So: is a decision-table `[initial]` over an observed tag a `malformed accessor
binding` (correct today, but a confusing diagnostic that never mentions the
class), or does C2 want a first-class refusal? And the corner C2 does not
cover: a decision table declaring `[initial]` over an observed tag **plus** a
writer accessor for that tag — the writer arity check then passes, so what
refuses? (`write to non-owned tag` fires in `renderWrites`, which `[initial]`
never enters.) If nothing refuses, C2's "already unauthorable" claim has a
hole.
**Blocks**: whether Phase 1 adds an `[initial]`-under-decision-table check
(which would be the *second* new refusal, against C2's "None of these is a new
refusal"), and what MVV step 6 / Scenario 4's stray-`[initial]` control
asserts. Coupled to A12, which is still `Pending`.

---

### P2-9 — C3's `emit` dump column lands in a closed vocabulary whose renderer has no emit case and whose 103 fixtures are hand-edited  [Medium]
**Anchor**: `0010:C3`, `0010:A4`, `0010:§phase-1-grammar-and-load`
**Trigger**: C3 — "`emit` MUST join the closed dump column vocabulary
(`0002:C19`), rendered as `key=value` pairs in key order"; D-wire-byte-format —
"in the dump it is one column, `key=value` pairs in key order, **bracketed like
`writes`**"; Phase 1 — "with every `[dump]`-carrying fixture updated".
**Question**: two mechanical questions the contract stops short of. (1) *Where*
in `dumpColumns` does `emit` go? `internal/table/dump.go:13` fixes the order
verbatim (`identity, source, kind, outcome, atoms, next, writes,
requires_owned, gate, escape`) and the comment says it is "fixed VERBATIM as
the normative fixtures author it". Appending `emit` at the end vs inserting it
after `writes` changes 103 golden files differently, and the *default* order
(`loadDump` sets `DumpOrder = DumpColumns()` when `[dump]` is absent) changes
every dump-golden regardless — Consequences claims "existing models load, lint,
and resolve byte-identically except for the new `emit` payload field and dump
column", which for dumps means every golden changes. (2) `bracket()` joins with
a separator; `writes` uses `"; "` and `renderTagValues` puts `key=value` where
value goes through `renderValue` (which quotes and brackets a non-1-member
sequence). "Bracketed like `writes`" gives `[a=x; b=y]` — but does the emit
value go through `renderValue` (quoted when set-kinded — and emit keys are
*undeclared*, so `m.Tags[key].Kind` is `""` and `len(members)==1` ⇒ bare) or
straight through as a raw string?
**Blocks**: the exact `dumpColumns` insertion index and the `column()` case
body — both of which must be fixed before the 103 fixture edits are made, since
redoing them is the expensive part of Phase 1.

---

### P2-10 — `Plan.RuleID` is not unique across expanded rows, and `rowByID` already relies on it  [Low]
**Anchor**: `0010:C3`, `0010:§technical-design` (item 1), `0010:§authority` (the-answer row)
**Trigger**: Technical Design — "the CLI joins it back on `Plan.RuleID`, which
`0002:C4` makes unique"; the Authority table — "joined by `Plan.RuleID` after
selection".
**Question**: `0002:C4` makes the *rule id* unique among rules, not among
normalized rows. `internal/table/normalize.go::expand` mints one row per
match-block `in` member, distinguished by `Suffix`, and
`internal/table/model.go::KernelRow` carries only `RuleID` — no suffix — so
several kernel rows share one `RuleID` and `internal/cli/flow_resolve.go::
rowByID` returns the first. For `emit` this is benign (emit is authored
per-rule, so every expansion carries the same block) — but the RDR asserts a
uniqueness property that is false as stated, and A9's "one `in`-atom rule
expanding to N" is exactly the shape that exercises it. Worth pinning as "emit
is rule-scoped, so the existing first-match join is sound" rather than as a
uniqueness claim, so the next reader does not build a per-row join on it.
**Blocks**: nothing structural — but it decides whether Phase 3 leaves
`rowByID` alone (correct) or a reviewer flags it as a latent
wrong-row-selected bug.

---

### P2-11 — `flow next` carries no `emit`, and C4's silence is ambiguous  [Low]
**Anchor**: `0010:C4`, `0010:A11`
**Trigger**: C4 — "`flow next`, `flow read-state`, and `flow set-state` are
unchanged by this RDR."
**Question**: `internal/cli/flow_next.go::candidate` already previews
`Next`, `Writes`, `Clear` "the preview the normalized row exposes without
evaluating anything (REQ-40, A-5)". Over a decision table all three are empty
by C2, so `flow next` reports candidates carrying **no** distinguishing
information at all — the answer lives only in `emit`, which C4 declines to add.
Is that deliberate (the caller must `flow resolve` to see the answer, and
`next` is enumeration only), or an oversight that a follow-up will have to
override C4 to fix? MVV step 5 asserts `flow next` exits 0 with empty
`required` and stops there, so nothing pins the intent.
**Blocks**: whether Phase 3 touches `flow_next.go` at all, and whether the
`candidate` struct's preview trio gains a fourth field now or never.

---

## Widening

The owned set (`C1`–`C5`, `D-*`) could not answer P2-1, P2-3, P2-5, P2-6, or
P2-11 on its own — each turns on what the contracts do **not** say. What sent
me:

- **P2-1** — C4 names a text rendering; nothing in the contracts names a
  renderer, so I read `internal/cli/respond/text.go` (reached from the
  `resolvePayload` source-anchor) and then `0010:S5` for the assertion.
- **P2-3** — C1's `⇒` and `§consequences` both claim "every existing model
  loads unchanged"; the contradiction is only visible against `0010:MVV` step 6
  and `0010:S4`'s class-omitted control, plus a sweep of the 103 fixture files.
- **P2-5 / P2-6** — C3 states obligations ("duplicate-free", "MUST be
  strings") with no category and no mechanism; `0010:S3` is the only place the
  test shape appears, and it is silent on both.
- **P2-8** — C2's "none of these is a new refusal" is a claim about existing
  code, verifiable only in `load.go::loadInitial` +
  `load.go::checkAccessorBindings`; the still-`Pending` `0010:A12` is the
  companion.
- **P2-9 / P2-11** — the `§phase-1-grammar-and-load` and `§implementation-plan`
  spans, reached because C3/C4 fix *what* changes but not *where*.
