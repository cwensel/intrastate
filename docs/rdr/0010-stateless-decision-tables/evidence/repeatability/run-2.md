model: claude-fable-5
variant: full (profile: foundational)

# Repeatability run 2 — RDR 0010 (stateless decision tables)

Read: C1–C5, MVV, S1–S6. Widened past the contract spans into:
`§approach`, `§technical-design` (seam list + data flow), `§authority`,
`§load-bearing-decisions` (D-identity, D-wire-byte-format, D-naming,
D-selection-predicate), `§illustrative-code`, `§problem-statement`, and the
metadata Overrides line. Widening resolved: the class accessor's zero value,
which two `len(Initial)==0` sites become class-aware, the `EmitValue` row
type, the rowByID join, and the payload field position — none of which are
fully pinned inside the C-spans alone.

## 1. Public API

The module is intrastate's existing surface; 0010 changes three seams
additively. Public here = the CLI verbs plus the exported Go seam between
`internal/table`, `internal/graphlint`, `internal/resolve`, and
`internal/cli`.

CLI (unchanged verbs, changed behavior/payload):

- `intrastate flow resolve --model <toml> --outcome <o> --tag k=v ... [--as=json|text] [--artifact role=path]`
  Success payload gains `emit`: JSON object of string values, keys in byte
  order, `{}` never `null` never omitted, declared immediately after `Rule`
  in `resolvePayload`'s thirteen-field struct (fourteen after the change —
  field order is the JSON key order). `next`/`writes`/`clear`/`owned` keep
  0005:C1 shapes and are empty over a decision table; `readers` is `[]`.
  `--artifact` for a role no invoked reader needs is ignored (empty demand
  set ⇒ no reader invoked, no `flow-artifact-missing`). `--outcome` remains
  required in both classes. Text mode: generic renderer only —
  `emit.<key>: <value>` per pair, `emit: (none)` when empty; no per-verb
  special case.
- `intrastate lint --model <toml> --as=json` — exit 2 with findings / exit 0;
  over a decision table the reachability root is the ∅ owned-state node.
- `intrastate dump --model <toml>` — column vocabulary gains `emit`,
  appended last after `escape`; cell renders key-sorted `key=value` pairs
  bracketed like `writes` (`[a=x; b=y]`), raw authored strings (not through
  `renderValue`); empty block renders the empty bracket.
- `flow next`, `flow read-state`, `flow set-state` — unchanged. `flow next`
  over a decision table succeeds with every rule a candidate, empty
  `required`/`next`/`writes`; it carries no `emit` (BR4 rejected).

Go seam (exported within the repo):

- `table.Model` gains a class field with accessor — GUESS: `func (m *Model)
  Class() ModelClass` with `type ModelClass string`, constants
  `ClassStateMachine ModelClass = "state-machine"`,
  `ClassDecisionTable ModelClass = "decision-table"` (the string values are
  normative; the Go identifiers are mine). Zero value reads as
  `state-machine`, so a hand-constructed Model keeps today's behavior. The
  accessor is the single source for "what class"; no site re-derives it from
  `len(owned) == 0`.
- `table.Row` gains `Emit []EmitValue`, key-sorted;
  `type EmitValue struct { Key, Value string }` — a new row-level type, not
  `TagValue`.
- `sourceModel` gains `class string`; `sourceRule` gains
  `Emit map[string]string` (decoder-typed so non-string / nested-table
  values refuse as a decode type error).
- Kernel `internal/resolve/resolve.go::Row` is untouched — emit never
  crosses the kernel boundary.
- `graphlint.Reasons()` closed set gains a fourth member,
  `no-participating-dimension` (append-only; 0006's assent, A14).

Error modes (categories, asserted on category never upstream text):

- `malformed model declaration` — `class` outside the two-value set; or a
  `"decision-table"` model with a non-zero owned count, detail carrying the
  literal token `owned=<n>` and naming the class. One-directional: a
  zero-owned `"state-machine"` model loads and lints as
  `graph-dangling-edge` (0006:C18 untouched).
- `unknown_schema_field` — a `class` key seen by a pre-0010 binary (A8,
  0002 strict decoding; not a new arm here).
- `malformed TOML` — non-string emit value or nested `[rule.emit.sub]`
  (decoder type arm; `CatUnknownSchemaField` does not apply to a well-named
  key of the wrong type).
- `malformed rule shape` — ordinary rule with no write block, now
  conditioned: fires for `state-machine` only, MUST NOT fire for
  `decision-table`.
- `unknown tag` / `malformed accessor binding` / `write to non-owned tag` —
  the pre-existing arms that make `[initial]`, `terminal`-over-owned,
  writers, and clears unauthorable for a decision table (C2: no new
  refusal).
- `malformed dump declaration` — explicit `[dump]` list omitting `emit`.
- `flow-ambiguous-match` / `no_match` — kernel exact-one is the hit policy
  in both classes (A5); an escape "otherwise" row rescues `no_match`,
  reports `escaped = true`, and may carry its own emit.

## 2. Three most important internal helpers

1. **Class/owned-set agreement check** — GUESS on name:
   `checkClassAgreement` (or an arm appended to an existing loader step);
   the RDR fixes its *position*, not its name: a step in
   `internal/table/load.go::run` at or after `loadTags`, never in
   `loadModelHeader` (the tag table is empty while the header loads), so an
   undeclared-tag refusal precedes a class-disagreement refusal under run's
   fail-fast order. Responsibility: read the class off the loaded model
   (parsed in `loadModelHeader`), count tags of provenance `owned`, and
   refuse `malformed model declaration` iff class is `decision-table` and
   the count is non-zero — detail names the class and renders `owned=<n>`.
   Never refuses the converse.
2. **Emit normalization and carry-through in
   `internal/table/normalize.go`** — inside `normalizeRule` (GUESS: a small
   `normalizeEmit(map[string]string) []EmitValue` helper): sort pairs by
   key; absent block → empty sequence; present-but-empty block → the same
   empty sequence *without* refusing (deliberately opposite the
   write/clear/gate blocks, which refuse even an empty one under 0002:C4);
   no dedup pass (TOML decoder already refuses duplicate keys). Also
   conditions the no-write-block `CatMalformedRuleShape` arm on the class.
   Critically, `expand` must copy `Row.Emit` into every row it mints from a
   multi-member match atom — expand builds each row from a seed literal, not
   by copying `base`, so an uncarried Emit drops silently and
   `rowByID`'s first-match join then returns `{}` on a rule that authored
   an answer (S3 makes this arm mandatory).
3. **Class-keyed lint arms in `internal/graphlint`** —
   `reach.go::reach` seeds the traversal at the ∅ owned-state node under the
   predicate `class == decision-table || len(Initial) > 0` (augmenting, not
   replacing, the `len(Initial)==0` early return: a rootless state machine
   still traverses nothing); `checkDanglingEdge`'s missing-root arm keys on
   class the same way while its terminal-predicate arm stays unconditional
   (A12) — that arm is C2's enforcement of the `terminal` prohibition; and
   `coverage.go::checkCoverage` gains the class-keyed zero-dimension arm:
   inside the `len(dims) == 0` branch, *before* `emitCoverageArms` (whose
   `bareEscapeFor` path would close the group first), emit
   `graph-unprovable-coverage` naming the group with reason
   `no-participating-dimension`; where both would apply,
   `graph-coverage-closed-by-escape` MUST NOT be reported. The other two
   `len(Initial)==0` sites (`checkAlwaysPresentOwned`,
   `checkUnreachableRules`) keep the bare test and never learn the class.

## 3. Data model (persisted / cross-boundary)

Authored TOML (persisted, the reviewable artifact per 0002):

- `[model]` MAY carry `class`, admitted values `"state-machine"` (default
  when absent) and `"decision-table"`; third optional key beside
  `id`/`version` (0002:C2/C3 overridden additively). Class is model data,
  carried on the loaded model, never inferred from the owned set.
- Any rule, either class, MAY carry `[rule.emit]`: a flat TOML table,
  values MUST be strings. Emit keys are not tag keys — undeclared,
  uninterpreted, byte-compared, never matched/guarded/written/read
  (`[rule.match.<emit-key>]` refuses `unknown tag`).
- A decision-table model: zero owned tags, no `[initial]`, no `terminal`,
  no accessor whose keys name an owned tag, no write block or clear list on
  ordinary rules. Discriminating dimensions are authored as
  `[rule.guard.all.<key>]` atoms (match atoms scope the group and
  contribute no product dimension — 0006:C7, 0003:C13).

Normalized model (in-memory, crosses table → lint/CLI):

- `Row.Emit []EmitValue` (`{Key, Value string}`), key-sorted; every row
  expanded from one rule carries the same block (emit is rule-scoped).
- Model class on `table.Model`, zero value state-machine.
- Kernel row keeps `RequiresOwned`/`Writes`/`NextTags` — all empty for
  decision-table rows (the shape an escape row already has, 0002:C14/C15);
  no `Emit` on the kernel row.

Wire (JSON payload, `flow resolve` success):

- `emit`: object of string→string, keys byte-ordered, `{}` when unauthored,
  positioned immediately after `rule`; HTML escaping disabled like every
  other payload map (0005:C1 — stated in D-wire-byte-format). Joined back
  from table data by `Plan.RuleID` after selection via the existing
  first-match `flow_resolve.go::rowByID`.
- Escaped selection: rescued plan reports `escaped = true` (A9).

Dump (persisted-ish fixture surface):

- `emit` column appended last after `escape`; cell `[k=v; k2=v2]` in key
  order, raw values; empty bracket when empty; default column set includes
  it for both classes; explicit `[dump]` must name it.

Lint findings:

- No new finding code; `graph-unprovable-coverage` gains reason
  `no-participating-dimension` (group-scoped). No finding's
  Code/Element/Message/Detail may name the class or the string
  `decision-table`. Fingerprint (`engine.go::Fingerprint`) excludes emit —
  editing an emit block never changes finding identity (A7).

## 4. Pseudo-code — `flow resolve` over a decision-table model (load → answer)

    resolve(modelPath, outcome, tagBindings, artifactBindings, mode):
      src   = decodeStrict(modelPath)             # unknown key -> unknown_schema_field
                                                  # bad emit value type -> malformed TOML
      # load.go::run fail-fast step order:
      m     = loadModelHeader(src)                # reads id, version, class
                                                  # class not in {state-machine, decision-table}
                                                  #   -> malformed model declaration
      loadOutcomes(src, m); loadTags(src, m)      # undeclared-tag refusals land first
      checkClassAgreement(m):                     # at/after loadTags (C1)
        if m.Class() == decision-table and countOwned(m.Tags) > 0:
          fail malformed model declaration        # detail: class + "owned=<n>"
      for rule in src.rules:                      # normalizeRule
        if rule.ordinary and no write block and m.Class() == state-machine:
          fail malformed rule shape               # arm suppressed for decision-table
        rule.EmitSeq = sortByKey(rule.emit or {}) # empty/absent -> [], no refusal
      rows = expand(rules)                        # each minted row carries rule.EmitSeq
      loadInitial/loadAccessors/checkAccessorBindings  # C2: [initial]/writers unauthorable
                                                       # with zero owned tags (existing arms)

      demand = RequiresOwned(candidates) ∪ guardOwnedKeys   # == ∅ for decision-table
      readers = invokedReaders(demand)            # ∅ -> none run; stray --artifact ignored
      view = assemble(tagBindings, outcome)       # observed dims + recognized outcome

      plan = kernel.Resolve(view, kernelRows)     # exact-one hit policy (A5)
        0 matches: escape row rescuing no_match? -> plan{escaped: true} : fail no_match
        2+ matches: fail flow-ambiguous-match
      row  = rowByID(rows, plan.RuleID)           # first match; sound: emit is rule-scoped
      payload = {
        rule:  plan.RuleID,
        emit:  toObject(row.Emit),                # {} never null; after rule in field order
        next: {}, writes: {}, clear: [], owned: {}, readers: [],  # 0005:C1 shapes, empty
        ... remaining 0005 fields unchanged ...
      }
      render(payload, mode)                       # json: field order; text: generic
                                                  #   flatten -> "emit.k: v" / "emit: (none)"

(Lint path, for contrast: `reach` seeds at ∅ when class is decision-table;
groups all reachable; overlap/coverage run over guard dimensions;
`len(dims)==0` group → `graph-unprovable-coverage` /
`no-participating-dimension` before `emitCoverageArms`.)

Residual GUESSes beyond those marked above: the loader step name
(GUESS: `checkClassAgreement`), the emit-normalization helper's existence
as a separate function (GUESS), the Go identifiers for the class type and
constants (GUESS — string values are normative, identifiers are not), and
the exact name of the dump cell renderer for emit (GUESS: inline in the
dump column code rather than a named helper — the RDR fixes rendering, not
a function).
