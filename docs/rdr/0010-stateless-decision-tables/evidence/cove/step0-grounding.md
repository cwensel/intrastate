Model: claude-opus-5[1m]

# cove Step 0 — GROUNDING sweep, 0010 (stateless decision tables)

Scope: semantic verification of the anchors carried by the record (42 resolved
`source-anchor` edges; the propose-time micro-sweep covered 33, so the 9-anchor
delta and all prose claims were treated as unswept), the unanchored prose
claims, and the inverse check on the four new rules this RDR introduces.

Every check below was made by opening the cited source, not by re-reading the
record. Symbol *existence* was not re-verified (already `resolved: true`).

## CONFIRMED roll-up

- **A1 / C5 — `nodeSatisfiesMatch` consults only owned atoms.** `internal/graphlint/analysis.go::nodeSatisfiesMatch` `continue`s on `a.model.Tags[atom.Key].Provenance != table.ProvenanceOwned`, so an owned-atom-free context is satisfied by any node, the ∅ node included. `contextReachable` is the only caller. CONFIRMED.
- **A1 / A10 / C5 — the rootless-model behaviour.** `internal/graphlint/reach.go::reach` returns `nil, true` on `len(m.Initial) == 0`; `internal/graphlint/groups.go::checkGroups` `continue`s past overlap/redundancy/coverage for any group whose context is unreachable; `internal/graphlint/analysis.go::checkDanglingEdge` fires its root arm unconditionally on `len(a.model.Initial) == 0`. The record's characterization — one `graph-dangling-edge`, silence elsewhere — matches source exactly. CONFIRMED.
- **A2 / C4 — demand set keyed on owned provenance.** `internal/cli/flow_exec.go::invokedReaders` unions `row.RequiresOwned` with `guardOwnedKeys`, which filters `m.Tags[atom.Key].Provenance != table.ProvenanceOwned`; `runReaders` checks `r.artifacts[def.Accessor.Role]` only inside `for _, name := range names`. Zero owned tags ⇒ empty demand ⇒ zero binding checks ⇒ no `flow-artifact-missing`. CONFIRMED.
- **A2 — supplied bindings are never swept.** `internal/cli/flow_input.go::parseArtifacts` does `strings.Cut(binding, "=")` into `out[role] = path` and never touches the model. Every consumer repo-wide is lookup-by-need: `flow_exec.go:186` (readers), `flow_exec.go:528` (gates), `flow_state.go:177` (writers), plus `artifactMap()` which copies the map verbatim into `accessor.Artifacts`. No site iterates the supplied bindings to validate them against declared roles. The "ignored rather than refused" clause CONFIRMED, including the "nowhere" part.
- **A3 — no kernel precondition on emptiness.** `internal/resolve/resolve.go::missingOwned` iterates `RequiresOwned` only. `Resolve`'s two preconditions are `in.Table.CheckValid()` and `precondition.go::CheckInput`, and `CheckInput`'s whole body is three `recognizedTagKey` scans — no `len(owned)`, no empty-view, no initial-state, no non-empty-`Writes` guard. `TagSet.matches` returns `true` on an empty `want`. `planOf` copies `NextTags`/`Writes` unconditionally. CONFIRMED.
- **A4 — closed 10-column vocabulary, exhaustive `order =`.** `internal/table/dump.go::dumpColumns` is exactly ten members. `internal/table/load.go::loadDump` refuses both directions; the omitted arm loops `vocabulary` and fails `CatMalformedDumpDeclaration` with `"column "+col+" is omitted"`. Absent `[dump]`, `l.model.DumpOrder = DumpColumns()`. CONFIRMED.
- **A4 census — independently re-counted.** `grep -rl '^\[dump\]' internal/table/testdata/ | wc -l` → **103**, and **103 of 103** carry an explicit `order =`. The 75 further repo-wide matches are all under `docs/rdr/0002-*/evidence/`. `models/rdr.toml` declares no `[dump]` and is the only file in `models/`. The record's "exactly 103" is CONFIRMED at the number I actually got.
- **A7 — `Fingerprint` hashes two fields.** `internal/graphlint/engine.go::Fingerprint` is a hand-written `strings.Builder` serialization over `row.Atoms` (`Key`, `Block`, `Operator`, `Literal`) and `canonicalTags(row.NextTags)` (`Key`, `Value`). No reflection, no marshal, no whole-`Row` hash; `Writes`, `RuleID`, `SourceLocator`, `Outcome`, `RequiresOwned`, `Gate`, `Escape` are all unread. A new `Emit` field cannot perturb it. CONFIRMED.
- **A8 — strict decode.** `internal/table/source.go::decodeStrict` calls `dec.DisallowUnknownFields()` and maps `*toml.StrictMissingError` to `fail(CatUnknownSchemaField, …)`. `sourceModel` admits exactly `id`, `version`, `description`, `metadata` — `class` is unknown today and refuses. `Metadata map[string]any` is the documented single point where strictness stops descending, which is why strictness *is* recursive elsewhere (including `[[rule]]`). CONFIRMED.
- **A9 / C4 — escape and payload shape.** `internal/graphlint/coverage.go::bareEscapeFor` filters `slices.Contains(row.Escape, class)` on `row.Kind() == table.KindEscape` with `len(guardAtomsOf(row)) == 0`. `CodeCoverageClosedByEscape` is emitted unconditionally when `closedBy != ""`, commented "a bare green does not satisfy the clause" — fires by design, never a silent default. `Plan.Escaped` is surfaced as `json:"escaped"` on `internal/cli/flow_resolve.go::resolvePayload`. The escape prohibitions in `normalize.go` key on `Write`/`Clear`/`Gate` presence only; nothing forbids an emit block. CONFIRMED.
- **A10 — fourteen codes.** `internal/graphlint/taxonomy.go` declares ten blocking + four advisory. The code is spelled `CodeOwnedBeforeWrite = "graph-owned-before-write"` (the record's correction is right). `CodeProductTooLarge` is one code under two call sites — `coverage.go` (group product) and `analysis.go::checkNodeCeiling` (traversal). CONFIRMED including both corrections the record folded in.
- **A10 silent bucket.** `checkSingleValuedState` iterates `writesOf(row)` per row (empty under C2); `checkOwnedBeforeMatch`/`checkAlwaysPresentOwned` filter `ProvenanceOwned`; `checkDeadEnd`/`checkTerminalEscape` depend on declared terminals. CONFIRMED.
- **C2 — the write-block obligation.** `internal/table/normalize.go:363-366`: `} else if rule.Write == nil {` → `fail(CatMalformedRuleShape, "ordinary rule "+id+" carries no write block")`, commented "An ordinary transition rule MUST contain a write block (`0002:C4`)". This is precisely the arm C2 conditions on class. CONFIRMED.
- **Risks — the two constants.** `internal/guard/product.go`: `const bound = 2048`, `func Bound() int`. `internal/graphlint/taxonomy.go`: `const nodeCeiling = 4096`, `func NodeCeiling() int`. Both blocking. CONFIRMED at the exact values the record states.
- **Key Discoveries — guard-only dimensions.** `internal/guard/product.go::Dimensions` iterates `guardAtoms(row)`, which admits `BlockAll`/`BlockUnless` only, commented "Match keys never enter". The record's guard-block correction is right. CONFIRMED.
- **Audit — `Model.Metadata`.** `internal/table/model.go:370-372` — `Metadata map[string]any`, "the one sanctioned extension namespace, carried through untouched and never interpreted (`0002:C2`)". Model-level, `any`-typed, and absent from `dumpColumns`, exactly as the audit's "Known Limit" states. CONFIRMED.

## INVERSE CHECK

### A model-level class discriminator — searched, none exists

Greps: `\bclass\b|\bKind\b|\bMode\b|\bVariant\b` across `internal/table`,
`internal/graphlint`, `internal/resolve`; `len(owned)|OwnedTags|ProvenanceOwned`
across all of `internal/`; the `sourceModel` and `Model` struct definitions read
in full.

Result: **no model-level "what kind of model is this" decision exists.** The
`Model` struct (`internal/table/model.go:364-393`) has no class/kind/mode field.
`sourceModel` (`internal/table/source.go:29-39`) admits four keys, none of them a
class. Every existing `class` identifier in the tree is a *refusal* class
(`resolve.RefusalKind` — `no_match`/`ambiguous_match`), an unrelated axis.
Every `ProvenanceOwned` site is a per-tag/per-atom filter; nothing anywhere tests
`len(owned) == 0` or branches on the owned set being empty. The nearest thing to
a shape branch is `len(m.Initial) == 0` in `reach` and `checkDanglingEdge`, which
is a root test, not a class test.

The one genuine sibling discriminator is per-*row*, not per-model, and the record
already names it: `internal/table/model.go::Row.Kind` — see Finding 1 for the
one place the record mis-cites it.

### `[rule.emit]`, an `emit` dump column, ∅-rooted seeding — searched, none exists

- `emit`: `grep -rn "\bemit\b|Emit"` across `internal/table`, `internal/resolve`,
  `internal/cli` returns only `clierr.EmitJSON`/`EmitText`/`EmitFindingsText`,
  `ExecuteAndEmit`, and `analysis.emit` — all output-writer verbs, unrelated to
  row data. `sourceRule` has no `emit` field; `Row` has no `Emit` field;
  `dumpColumns` has no `emit` member. No sibling.
- ∅-rooted seeding: `reach` has exactly one root construction path
  (`root := Node{Values: map[string][]string{}}` populated from `m.Initial`),
  guarded by an early `return nil, true` when `m.Initial` is empty. There is no
  second seeding path and no caller that supplies a root. No sibling.
- The uninterpreted-literal *treatment* does have a sibling — `Model.Metadata` —
  which the Existing Infrastructure Audit already names and correctly scopes as
  "Reuse the treatment, not the field". Not a finding.

## FINDINGS

### 1. `Row.Kind` is anchored to the wrong symbol

- **Element**: `0010:§approach` (also echoed at `0010:ALT2` Cons, unanchored)
- **Claim as the RDR states it**: "The sibling discriminator in this codebase —
  `Row.Kind` (ordinary vs escape, dump column `kind`) — is inferred from the
  *presence* of an `escape` list
  (`internal/table/normalize.go::normalizeRule`)."
- **What the source actually says**: the discrimination is not in
  `normalizeRule`. It is a derived method in a different file:

  ```
  internal/table/model.go::Row.Kind
  // Kind is the derived row-kind view property, discriminated solely by a
  // non-empty escape class list (`0002:C5`).
  func (r Row) Kind() Kind {
      if len(r.Escape) > 0 {
          return KindEscape
      }
      return KindTransition
  }
  ```

  `normalizeRule` does compute `isEscape := rule.Escape != nil` (line 336) and
  stores `Escape: sortedList(rule.Escape)` (line 460), but that is *populating the
  field*, not inferring the kind. The inference — and the `kind` dump column
  (`internal/table/dump.go:81` → `string(r.Kind())`) — lives on
  `model.go::Row.Kind`.
- **Why it matters**: this is the record's single load-bearing piece of prior-art
  evidence for "declared, not inferred", and it is the rebuttal ALT2 rests on. The
  substance of the claim is **correct** (a presence-inferred, derived,
  never-stored discriminator with a dump column, precisely as described) — only
  the anchor is wrong, and it points at the symbol C2 *also* names as the site it
  changes, which conflates two distinct claims about `normalizeRule`. An
  implementer following the anchor to justify the ALT2 rejection lands in the
  wrong function. Correct anchor: `internal/table/model.go::Row.Kind`.

### 2. C2's `terminal` clause is unenforceable as written at the site it names

- **Element**: `0010:C2`
- **Claim as the RDR states it**: "A decision-table model MUST NOT declare
  `[initial]`, `terminal`, … None of these is a new refusal: each is already
  unauthorable once the owned set is empty — … and a terminal predicate over a
  non-owned tag is 0006's `graph-dangling-edge` terminal arm."
- **What the source actually says**: the terminal arm is real —
  `internal/graphlint/analysis.go::checkDanglingEdge` loops `a.model.Terminal` and
  emits `CodeDanglingEdge` for any atom whose tag is not `ProvenanceOwned`. But
  `graph-dangling-edge` is exactly the code C5 suppresses for this class:

  ```
  0010:C5 — "The missing-root arm of invariant 1 (`0006:C18`) … MUST NOT emit"
  ```

  C5 scopes its override to the *missing-root arm* only, and the Technical Design
  says `checkDanglingEdge`'s "missing-root arm keys on class", which leaves the
  terminal arm live. So the two contracts are consistent as written. The gap is
  narrower: C2 asserts the terminal prohibition is "already unauthorable", but a
  decision-table model that declares `terminal` is caught **only at lint**, not at
  load — while the sibling prohibitions in the same sentence (`[initial]`, accessor
  keys, write/clear keys) are load failures. C2's own framing ("each is already
  unauthorable") flattens a load-time refusal and a lint-time blocking finding into
  one category.
- **Why it matters**: an author who runs `flow resolve` without `lint` on a
  decision table declaring `terminal` gets no refusal at all — the model loads,
  and `Model.Terminal` is simply never consulted by `internal/resolve`. The
  prohibition is real but its enforcement point differs from the three it is
  bundled with, and C2 does not say so. Phase 1's "grammar and load" plan does not
  list a terminal check, so this is a decision the implementer will have to make
  unaided: enforce at load (a genuinely new refusal, contradicting C2's "none of
  these is a new refusal") or accept lint-only coverage.

### 3. C3's dump-column addition has an unstated blast radius the plan does not carry

- **Element**: `0010:C3` (census in `0010:A4`; plan in `0010:§phase-1-grammar-and-load`, `0010:S6`)
- **Claim as the RDR states it**: A4 — "every checked-in model or fixture carrying
  an explicit `[dump]` column list is updated in the same change… **Count
  confirmed exact: 103**."
- **What the source actually says**: the 103 count is right and I reproduced it
  (103 files, 103 with `order =`). But A4's next sentence — "the 75 further
  repo-wide matches are all under `docs/rdr/0002-*/evidence/`, i.e. 0002's spike
  fixtures" — is used to *exclude* them from the update, and that exclusion is
  where the claim gets thin. I confirmed the location breakdown:

  ```
  docs/rdr/0002-…/evidence/spikes/iter-2/neg   43
  docs/rdr/0002-…/evidence/spikes/iter-2/delim 15
  docs/rdr/0002-…/evidence/spikes/iter-2       10
  … (75 total, all under docs/rdr/0002-*/evidence/)
  ```

  These are RDR evidence artifacts, so leaving them stale is defensible under the
  no-amend rule. The unstated part is that they become **non-loadable** the moment
  `emit` joins `dumpColumns` — `loadDump`'s omission arm fails every one of them —
  so any future replay of 0002's spikes against a current binary refuses. A4 asserts
  the exclusion without noting that consequence.
- **Why it matters**: the record's `S6` regression sweep is scoped to "every
  checked-in model and fixture", which reads as the 103. If a 0002 replay harness
  exists or is later written, it breaks silently at load with
  `malformed_dump_declaration`, and the diagnosis path points at 0002, not at 0010.
  This is a one-line note in C3's ⇒ or in `S6`, not a design change.

### 4. C4's `emit: {}` normalization has no precedent in `resolvePayload`'s existing never-null handling

- **Element**: `0010:C4`
- **Claim as the RDR states it**: "The `flow resolve` success payload MUST carry
  `emit`: a JSON object of string values, keys in byte order, present as `{}` —
  never `null`, never omitted."
- **What the source actually says**: `internal/cli/flow_resolve.go::resolvePayload`
  does enforce never-null, but by two *different* mechanisms depending on the
  field, and neither is a general rule:

  ```
  Writes:  map[string]string{},   // initialized at literal-construction
  Clear:   []string{},            // initialized at literal-construction
  ...
  if payload.Gates == nil { payload.Gates = []gateResult{} }   // post-hoc
  if payload.Next  == nil { payload.Next  = map[string]string{} } // post-hoc
  ```

  `Next` needs the post-hoc guard because `tagMap(plan.NextTags)` can return nil;
  `Writes` does not because it is built by an append loop. Which mechanism `emit`
  needs depends on whether the join-by-`RuleID` helper returns a nil map on a row
  with no emit block — undetermined by the record.
- **Why it matters**: minor, but it is exactly the class of bug that produces
  `"emit": null` in one arm and `{}` in another, and C4 states the invariant
  without naming the mechanism. `0010:S5` asserts the payload shape; it should
  assert it specifically over a decision-table row that authors **no** emit block,
  which is the only arm where the nil can surface. The Testing Strategy's S5 as
  written ("`flow resolve` and `flow next` over the fixture") does not pin that.

## Verdict rationale

No REFUTED anchor. All 42 anchors carry semantically accurate claims at the sites
they name, with the single exception of Finding 1 (right claim, wrong symbol). The
A4 census reproduces exactly. The inverse check is the strongest result here: the
class discriminator, `[rule.emit]`, the `emit` dump column, and ∅-rooted seeding
each have **no** existing sibling — the RDR is not reinventing a discriminator the
code already holds, and its one true sibling (`Row.Kind`, per-row, presence-inferred)
is already cited as the contrast case rather than as a reusable mechanism.

Findings 2–4 are specification-completeness gaps, not grounding refutations: each
is a place where the record's prose is true but under-determines what the
implementer must do. None blocks lock.
