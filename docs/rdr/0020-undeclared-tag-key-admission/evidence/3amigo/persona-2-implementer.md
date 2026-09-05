Model: claude-sonnet-5

# Persona 2 — Implementer

Widening note: I widened past the owned elements (0020:C1, D-identity,
D-wire-byte-format, D-naming, D-selection-predicate) into
§technical-design, §implementation-plan, §illustrative-code, and
§critical-assumptions A1/A3, and cross-checked the named symbols
directly in `internal/cli/flow_input.go` and `internal/cli/flow_state.go`
under source. The owned elements describe a WHAT (the carrier
semantics) but the actual Monday-morning questions are about the WHERE
in the existing switch statement — a silence that has no line range in
the contract elements themselves, only in the surrounding code-shape
text and the real source.

All source anchors I checked resolve exactly as cited: `parseTags`
(flow_input.go:642), `canonicalValue` (flow_input.go:714, empty-value
check at :723-726), `codeInvalidFor` (flow_input.go:757),
`parseWrites`/`writerFor` (flow_state.go:434/623). No anchor is stale
or misnamed. The `rdr inspect --json --filter edges` pass found zero
unresolved `source-anchor` edges for 0020 — the two unresolved edges
are a `kind: issue` (kata pointer) and a `kind: artifact` (gate.md),
neither in my owned set.

### Medium: exact insertion point of the hoisted empty-value check is shown only in prose, not in the illustrative code

- **Anchor:** 0020:D-selection-predicate
- **Finding:** D-selection-predicate and §implementation-plan Step 1
  both state the empty-value arm "lifts out of `canonicalValue`'s
  `!isSet` branch onto the shared admission path, ahead of the carrier
  branch" — but the §illustrative-code snippet starts at the empty
  check and the `decl, declared := m.Tags[key]` branch; it never shows
  the actual `parseTags` body it is being spliced into, i.e. the
  reserved-key / owned-key / duplicate-key `switch` at
  `internal/cli/flow_input.go:654-667` and the `seen[key] = true` at
  line 667. Concretely: does the hoisted empty check go inside the
  `switch`'s cases (as a fourth case, `case value == "":`) or as a
  separate `if` statement between the switch and `decl := m.Tags[key]`
  at line 675? The ordering the RDR requires (grammar → reserved →
  owned → duplicate → empty-value → carrier-or-conformance) is
  satisfiable either way, but the two shapes differ in whether
  `seen[key] = true` runs before or after the empty check fires — which
  matters for whether a caller who repeats an empty `--tag k=` twice
  gets `flow-tag-invalid` (empty) both times or `flow-tag-duplicate`
  the second time. I would ask which is intended before writing the
  diff.
- **Blocks:** 0020:C1's stated refusal precedence and
  0020:D-selection-predicate's ordering guarantee.

### Medium: is the empty-value check duplicated verbatim in two functions, or shared?

- **Anchor:** 0020:C1
- **Finding:** C1's normative text says the empty-value arm "MOVES to
  the admission path ahead of the carrier branch," and
  §implementation-plan Step 1 clarifies "`canonicalValue` keeps its
  copy for the `--write` carrier, which enters by its own path" — so
  `canonicalValue`'s existing check at flow_input.go:723-726 stays
  in place for `--write`, and `parseTags` gains a NEW, separate empty
  check ahead of its carrier branch. That means the message string
  `"was given an empty value"` and its `userErr(codeTagInvalid, key,
  ...)` construction will exist in source twice. Is a second, hand-
  copied literal acceptable, or should Step 1 introduce a small shared
  helper (e.g. `checkEmpty(key, value, flag string) *clierr.CLIError`)
  to keep the two copies from drifting under a future message-wording
  change? The RDR doesn't rule either way, and "no new refusal code is
  minted" (C1) only promises the CODE stays `flow-tag-invalid`, not
  that the message text is centralized.
- **Blocks:** 0020:D-selection-predicate (empty-value arm) and
  §implementation-plan Step 1's diff shape.

### Low: codeInvalidFor's doc comment claims `--outcome` shares this refusal path, but it does not — pre-existing, adjacent to the exact function C1 edits

- **Anchor:** 0020:C1
- **Finding:** `internal/cli/flow_input.go:754-756`'s comment on
  `codeInvalidFor` says "`--tag` and `--outcome` take
  `flow-tag-invalid`," but grounding shows `--outcome` is a plain
  `cmd.Flags().GetString("outcome")` in `flow_resolve.go:237-241` with
  its own direct empty-check and its own `userErr(codeTagInvalid,
  "outcome", ...)` call — it never reaches `parseTags` or
  `canonicalValue` at all. This is a pre-existing stale comment, not
  something C1 introduces, and 0020 makes no claim about `--outcome`.
  I flag it only because Step 1 touches the exact function this stale
  comment sits beside; not a blocker, but I'd want a "not touching
  this, it's already wrong" note in the PR so a reviewer doesn't
  assume the RDR's silence on `--outcome` was an oversight.
- **Blocks:** nothing in 0020 directly; risk is reviewer confusion
  during Step 1's diff, not an implementation blocker.

### Low: MVV/testing steps don't state the fixture's fourth key's declared kind, only that it's "nothing named extra/extras"

- **Anchor:** 0020:§implementation-plan
- **Finding:** MVV step 1 says the fixture model declares "one scalar
  tag and one set tag, with no declaration for `extra` or `extras`" —
  consistent with Testing Strategy fixture F1's spike-sourced model
  (scalar `tier`, set `labels`). This is answerable from the existing
  spike evidence (`evidence/spikes/carrier-table.toml`,
  `a2-undeclared-scalar-only.txt`) so it is not a true gap — I checked
  because the desk-trace only asserts consistency for the SCALAR leg
  of F1, marking the array leg "post-change, red test adds it." First-
  hour question this answers itself: reuse `carrier-table.toml`
  directly rather than authoring a new fixture, since MVV step 1's
  fixture shape is already pinned by the spike.
- **Blocks:** none — resolved by existing evidence, included for
  completeness since I widened into §implementation-plan.
