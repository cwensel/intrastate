Model: claude-opus-5[1m]

# Cove Findings — RDR 0002 (Transition Table As Reviewable Data)

Chain-of-Verification lens, iteration 1. Step 0 grounding sweep + 12
independent verification questions + the fired structural mini-check cues
(authority census, test-discriminability, round-trip/fidelity, disposition,
desk trace).

Findings are the origin ledger for this lens's resolve pass.

## Step 0 — Grounding sweep

Every `path::Symbol` in the RDR's Critical Assumptions was read on the current
branch. **A2, A8, and the Key Discoveries row-shape claim are CONFIRMED
field-by-field and in the exact order the RDR states** (see Appendix). Two
claims did not survive.

- **CV-001 REFUTED — A5's parse-code precedent does not exist.**
  A5's Evidence says `internal/cli/config/config.go::Load` "already uses stable
  config load/read error codes (`config-not-found`, `config-read-error`) and
  marks `config-invalid` as the planned parse-validation path." The first half
  is CONFIRMED. But `config-invalid` is emitted **nowhere** — it appears only in
  a doc comment and a `TODO`. `Load` reads bytes, stamps
  `SchemaVersion: CurrentSchemaVersion` unconditionally without reading the
  file's own version, and returns a nil error. The codes actually emitted
  anywhere in the CLI are exhaustively `config-not-found`, `config-read-error`,
  `flag-invalid-value`, and `command-error`. **No parse-validation code has ever
  been emitted through the envelope.** A5's load-bearing half (the envelope
  itself — `clierr::CLIError` + `respond::Fail`) stands; the *precedent* half
  does not. This RDR's table parse/lint codes would be the first of their kind,
  which is a materially different risk posture than "already uses."

- **CV-002 REFUTED (scope) — the Prerequisites bullet understates the missing
  kernel surface.**
  The bullet says the kernel "still carries `Row.Guard` as a string and exports
  no `OpExists`/`LiteralTrue`/`LiteralFalse`." Both true (zero hits repo-wide).
  But there is **no atom type at all**: no `Atom`, no operator-token type, no
  `Block` field. `Row.Match` is `[]Tag` — a flat `{Key,Value}` equality list —
  and `TagSet::matches` is pure value equality with no operator dispatch. So the
  normative clauses requiring each atom to "retain the key, operator token,
  literal, and the block" and to emit the existence constants "verbatim" target
  a row shape with **no field to carry them**. Phase 2's four deliverables
  (outcome lifted, per-atom block retained, `RequiresOwned` derived, existence
  constants emitted) are three-quarters unbuildable against the shipped `Row`;
  only `Outcome` and `RequiresOwned` have real fields. The Prerequisites gate is
  wider than the RDR scopes it.

### Inverse sweep (new rule vs. existing sibling)

- **CV-003 — exact-byte tag-key identity has a confirming in-repo sibling the
  RDR does not cite.** `resolve.go::TagSet.matches` / `::Lookup` / `::has` index
  by raw Go map lookup — byte-exact, no folding; the only `strings.ToLower` in
  the repo is on the `--as` flag value in `respond.go::ModeOf`. The RDR grounds
  this rule on RDR 0007 A22 alone. The shipped kernel already makes the same
  decision and is the stronger cite.
- Outcome lifting, `RequiresOwned` derivation, and load-time existence-literal
  rejection: **searched, none exists** — no normalizer, no `clear` concept, no
  load path anywhere in `internal/`.
- Deterministic row sort: a sibling exists and **deliberately diverges** —
  `resolve.go::compareRefs` orders on `(RuleID, SourceLocator)`. The RDR flags
  this divergence correctly, but attributes it to RDR 0007's spec; the concrete
  sibling is a shipped, tested symbol.

## Internal contradictions

- **CV-004 — the same input classes are routed to both load and lint.**
  Three sites disagree:
  - The validation-category `normative` block (l.559) lists **`unknown tag`**,
    **`write to non-owned tag`**, and **`ambiguous overlap`** as load-time
    data-level categories.
  - **Failure Modes** (l.890) puts **"write to an undeclared tag"**,
    **"overlap"**, and **"ambiguous expansion"** in the *lint* bucket
    ("...is a lint failure before the model is accepted").
  - **Load-Bearing Decisions** (l.597) says "**Lint** rejects ambiguous ordinary
    overlaps and ambiguous escape overlaps"; **Technical Design** (l.322) says
    "**Validation** rejects ... ambiguous overlaps between candidate rows."
  The spike refuses undeclared writes at *load*
  (`neg-undeclared-write.toml`). RDR 0006 owns graph lint. The MVV requires a
  variant "refused as ambiguous" without saying by which layer. One class, two
  owners, no rule for which fires.

- **CV-005 — row kind is required by this RDR and normatively forbidden at the
  kernel boundary by Final RDR 0009.**
  RDR 0002 states "The normalized candidate row **is the kernel row**" (l.339)
  and then requires that row to carry "the row kind (`transition` or `escape`)"
  (l.343), and makes it a mandatory dump field (l.489, l.610). But
  `internal/resolve/resolve.go::Row` has **no kind field**, and RDR 0009's
  normative block says: "No ordinary/escape type split and no row-kind field is
  introduced at the kernel boundary. Escape identity remains discriminated
  solely by a non-empty `Escape` list." Either the normalized row is *not* the
  kernel row, or this RDR introduces a field a Final peer forbids. The RDR does
  not say which.

- **CV-006 — `[dump]` field selection contradicts the dump's "every field"
  MUST.** Schema part 6 (l.306) gives `[dump]` "deterministic ordering **and
  field selection**", and the spike fixture spells `[dump].order` with a
  nine-entry list. The dump `normative` block (l.488) says the dump "MUST carry
  **every field**" of the normalized row, and the Round-Trip invariant depends
  on that. An author omitting a field from `[dump].order` — which the schema
  plainly permits — breaks both. The RDR never says `[dump]` may reorder but not
  omit. (Note: the spike parses `[dump].order` into a `Dump` struct and then
  **never reads it** — `render()` hardcodes the field order — so the one witness
  the RDR offers for the dump contract does not exercise `[dump]` at all.)

## Silences (missing requirements)

- **CV-007 — cyclic / self-referential context inheritance has no clause and no
  category.** The RDR says contexts "MAY inherit from other contexts, but
  inheritance MUST normalize to an explicit predicate set before lint or
  resolution" (l.399) — which a cycle makes impossible. The tokens
  `cycle`/`acyclic`/`self-inherit` appear **zero times** in the RDR, and the
  14-item validation-category list has no entry for it. The spike had to invent
  the check to run at all: `main.go::mergeContext` → `"cyclic context
  inheritance at %q"`. The implementation already has behavior the contract does
  not name, so RDR 0005 cannot map its code and RDR 0006 cannot enumerate it.

- **CV-008 — a single-member `in` on `recognized` produces a different row
  identity than the equivalent `eq`.** The RDR says an `in` atom "expands into
  one candidate row per member, each identified by the rule id plus the outcome
  literal as its expansion suffix" (l.452), and that "an absent expansion suffix
  sorts before any present one" (l.496). So `in = ["a"]` yields a row with
  suffix `#a` while `eq = "a"` yields an unsuffixed row — two spellings of one
  edge, two durable identities, two sort positions. Row identity is in the
  Round-Trip invariant's preserved field list, so this is visible in the dump
  and in golden tests. The RDR is silent on whether a single-member `in` should
  normalize to the unsuffixed form or be rejected.

- **CV-009 — the alphabet's duplicate-free and empty-string rules have no
  validation category.** The recognized-totality clause (l.443) requires the
  `outcomes` alphabet to be "non-empty, duplicate-free, and MUST NOT contain the
  empty string" — three distinct obligations, and the empty-string one is the
  *sole barrier* A8 names for its residual. The category list offers only
  "missing or malformed recognized outcome alphabet." The spike emits three
  distinct refusals here. A caller cannot distinguish them, and the barrier A8
  leans on has no nameable category.

- **CV-010 — "normalization explosion" is routed two ways and has no
  threshold.** Failure Modes (l.890) says normalization explosions "fail at
  load/validation time with stable CLI errors." Risks and Mitigations (l.867)
  says "lint must **report** expansion counts per rule." Failure or diagnostic?
  No threshold, no code, no category in the 14-item list.

- **CV-011 — the round-trip invariant names an inverse the RDR never
  specifies.** The invariant is `parse ∘ normalize ∘ dump = expanded-table value
  identity`, and reads "dumping the normalized model and **reading the dump as a
  table view** must preserve the candidate-row set" (l.607). The RDR specifies
  the dump's *field list* and *row ordering* and nothing else — no grammar, no
  delimiter, no escaping, no quoting, no record separator. `dump` has no
  inverse because it has no syntax, so the invariant is unverifiable as written.
  Concretely lossy/ambiguous sites in the spike's own normative output:
  - **`<clear>` sentinel collides with a legitimate value.** `Write.Clear` is a
    structural `bool` in memory but renders as the literal string `<clear>` in
    the same position as a value. Nothing reserves it in the tag-value space, so
    an authored `status = "<clear>"` is indistinguishable from a clear of
    `status` on read-back.
  - **Unescaped separators.** `output.txt` renders
    `atoms=[...; ...]`, `requires_owned=[a,b]`, `k=v` — with `; `, `,`, `=`, and
    `]` all unescaped. The real output line
    `profile.in=foundational large@match` shows an `in` literal whose members
    are **space-separated and unescaped**, so a member containing a space is
    indistinguishable from two members.
  - **Expansion-suffix separator is unspecified.** `#` appears only in A7 prose
    and the spike's `identity()`; the normative blocks name only the tuple. Rule
    ids have no grammar beyond byte-uniqueness, so `(rule "a", suffix "b#c")`
    and `(rule "a#b", suffix "c")` render identically.

## Weak oracles (test-discriminability)

- **CV-012 — the MVV's most discriminating assertion cannot be witnessed by the
  fixtures the RDR names.** MVV/Testing scenario 4 requires "one tag-set in
  which **a sibling candidate's** guard is unevaluable" and expects
  `guard_unevaluable` "even though a decidable sibling and a `no_match` escape
  row exist." But in `rdr-fixture.toml` **every outcome binds exactly one
  transition row** — `round-clean`→`continue-prelock`,
  `reconcile-block`→`reconcile-rewind`, and `terminal-archive`'s two expansions
  are separate outcomes, not siblings. There is no outcome with two sibling
  candidates anywhere in either fixture, so the sibling-interaction property —
  the core of the gate-then-count claim this RDR exists to uphold — has no
  possible witness in the named canonical fixtures.

- **CV-013 — the `Resolve` handshake is entirely unwitnessed.** The spike never
  imports `internal/resolve` (`main.go` mirrors the constants locally as
  `const opExists = "exists"` etc.). So MVV's "resolves through
  `internal/resolve::Resolve` to exactly one ordinary row" and the modeled-escape
  assertion have no executable oracle, and scenario 6's "byte-for-byte against
  `resolve.OpExists`/`LiteralTrue`/`LiteralFalse`" is asserted against symbols
  that do not resolve on `main` (CV-002) and verified against a **local copy**.
  Per the RDR's own Method vocabulary ("The cited symbol must resolve on
  `main`"), that assertion is currently unsatisfiable.

- **CV-014 — the category oracle is a substring match, not a category
  assertion.** The contract names 14 categories; `negative-cases.txt` records
  free-text refusal messages for 9 and pins the *category* of none. An
  implementation collapsing all 14 into one generic code passes every witnessed
  case. Six classes have no control at all: unknown tag *matched* (vs written),
  unknown context, unknown accessor, ambiguous overlap, unknown schema field,
  and a rule with *two* `recognized` atoms.

- **CV-015 — determinism oracles pass on a no-op, and the permutation control
  is the weaker fixture.** "Repeated dumps are byte-identical" and "unchanged
  when authored with keys in a different order" are both satisfied by a
  normalizer emitting a constant or dropping every `unless` atom. They are only
  saved by pairing with scenario 2's golden, which the RDR does not state. And
  the key-order permutation control is **kata-only** — the RDR fixture, which
  carries the escape row, the `in`-expansion, `<clear>`, and inherited contexts,
  was never permuted. Separately, `exists = true` is **never exercised** (only
  `false` appears), so a normalizer hardcoding `LiteralFalse` passes everything
  in evidence; and the `Status`/`status` case-folding control named by scenario
  6 has **no fixture**, so a normalizer applying `strings.ToLower` to keys passes
  all witnessed evidence — defeating the byte-equality contract outright.

- **CV-016 — "assert the comparator is total" is satisfied by the defect it
  claims to exclude.** The MVV says: sort, then assert no adjacent pair compares
  equal, "so a future identity change cannot silently reintroduce a positional
  tiebreak." A comparator that falls back to table position when identities are
  equal produces **zero** adjacent-equal pairs — that is precisely what the
  tiebreak accomplishes. Detecting a positional tiebreak requires permuting
  *rule declaration order* and asserting byte-identical dumps; the RDR only
  specifies permuting *TOML key order*, and the spike evidence likewise attests
  key-order, not rule-order, independence.

## Authority / ownership

- **CV-017 — the dump's totality proof holds per-model, but nothing bounds a
  dump to one model.** The clause argues totality from rule ids being unique
  "within a model" plus alphabet duplicate-freedom (l.500). The tuple leads with
  `model id`, so it *is* total across models — provided model ids are unique.
  The RDR never states that a file or a dump contains one model, nor constrains
  model-id uniqueness. The spike sidesteps it by taking two *files* and emitting
  two independent `MODEL` sections, so the cross-model case the contract claims
  is unwitnessed. (RDR 0006 independently assumes "exactly one model is linted
  per invocation" — an assumption RDR 0002 never grants it.)

- **CV-018 — gate-then-count is reproduced as this RDR's own normative block
  though this RDR does not own it.** The RDR assigns the guard seam to RDR 0007
  and the ordering to JDR 0001 §D2 / the kernel, then restates the entire
  ordering — including the non-escapability rule — as a `normative` block
  (l.540). That is three normative homes for one rule; if the kernel's `gate`
  changes, this RDR silently becomes wrong. The RDR applies exactly the right
  discipline to the tag type model ("MUST NOT be restated here", l.554) and not
  to this.

- **CV-019 — escape-row emptiness is a producer obligation stated as a derived
  guarantee.** The clause reads "an escape row **therefore** carries an empty
  set (RDR 0007 A21)" (l.466). But the kernel actively tests the opposite as
  live behavior: `adversarial_test.go::TestAdv2b_EscapeEdgeMustNotBypassThe`
  `OwnedStateRequirement` builds an escape row with
  `RequiresOwned: []string{"never-present"}` **and** `Writes`, and asserts the
  kernel refuses `owned_state_unavailable`. `escapeOrRefuse` passes escape
  candidates through the identical `gate`. So a non-empty `RequiresOwned` on an
  escape row is tested kernel behavior, not an impossibility; "therefore" reads
  as a guarantee where only a normalizer-side obligation exists.

## Cross-RDR drift (route to Stage 6 — not a 0002 edit)

- **CV-020 — a Final peer normatively asserts what this RDR now forbids.**
  RDR 0008's `normative` block states: "A table whose declared `outcomes`
  alphabet contains the empty string would pass the alphabet gate ... **nothing
  in this RDR or RDR 0002 forbids that alphabet entry. It is left unforbidden
  rather than silently assumed away**." RDR 0002's recognized-totality clause now
  forbids it, and A8 leans on that clause as "the sole barrier." A8 already
  notes the 0008 text "predates that clause," and the two *decisions* are
  compatible (0008 scoped itself: "not this RDR's to rule on") — but a Final
  peer carries a stale factual assertion about this RDR's contents.

- **CV-021 — this Draft producer is locked-against by five Final consumers.**
  RDR 0002 is `Draft`; RDRs 0003, 0006, 0007, 0008, 0009 are all `Final`
  (0001 `Implemented`). RDR 0009 twice refers to "RDR 0002's **Final**
  escape-rule prohibition" and rests its two-point enforcement argument on it.
  The producer being the only unlocked document inverts the normal dependency
  order: any contract this lens changes can silently invalidate a Final peer's
  verified assumption.

## Dismissed at the grounding gate

- **`RequiresOwned` derivation "conflicts" with the kernel's `missingOwned`.**
  Raised as a High. **Dismissed — re-raise of an adjudicated decision.** RDR
  0007's Metadata says it "*fixes* the meaning of `Row.RequiresOwned`
  (post-guard write-dependency keys)" and names kata `xg7p`
  ("RequiresOwned conflates guard-input state with post-guard transition state")
  as the defect it closes. Narrowing away from guard-input state is the
  deliberate call; 0007 A7 covers the residue (every guard-referenced key is
  declared, so a typo fails at table load). RDR 0002's "does not add guard-read
  keys to it" is faithful to that settled decision.

- **`next-state tags` and `writes` being value-identical undermines RDR 0009
  A4.** **Dismissed — inverts A4.** 0009 A4 settled that the breach predicate
  stays `Writes`-only and explicitly does **not** widen to `NextTags`, because
  only `Writes` reaches the accessor layer. Two fields populated from one
  rendered set, with different downstream consumers, is coherent and matches
  `resolve.go::planOf`, which copies both.

- **"The normative spike evidence is missing from the repo."** **Dismissed —
  false against disk.** `evidence/spikes/` contains `output.txt` (1899 B),
  `negative-cases.txt` (1711 B), both fixtures, and `main.go`. All read
  directly this pass; `git status` clean.

## Appendix — grounding claims CONFIRMED

- A2: `::Resolve` runs `gate(...)` and returns on `blocked != nil` **before**
  `switch len(selected)`. `::gate` prunes `GuardFalse` → `missingOwned` →
  vetoes undecidable, in that order. `::escapeOrRefuse` delegates to the same
  `gate` before counting `viable`. `::missingOwned` ends `slices.Sort(missing)`
  and carries the documented "Row order is a normalization detail RDR 0002
  owns" comment verbatim. `adversarial_test.go::TestAdv3_GuardUnevaluable`
  `RefusalMustNotDependOnTableRowOrder` exists and asserts equality across
  permuted `Table.Rows`.
- A8: `::Resolve` skips on `row.Outcome != in.Recognized` **before**
  `view.matches`; `::escapeOrRefuse` applies the same filter; `::assemble`
  binds `recognizedTagKey` only when `in.Recognized != ""`; `unmodeled_outcome`
  is refused via `!in.Table.models(...)` as the first act after `assemble`;
  `::Table.models` is a bare `slices.Contains` with no well-formedness check.
  **The residual was reproduced end-to-end**: with `Recognized == ""` and `""`
  in the alphabet, a plan is emitted against a view carrying no `recognized`
  key. Note the barrier is load-time only — a caller building a `resolve.Table`
  directly in Go (which every current test does, there being no loader)
  bypasses it entirely.
- A5 (envelope half): `clierr::CLIError` carries all seven of `Code`,
  `Message`, `Param`, `Detail`, `Hint`, `Group`, `Cause`; `respond::Fail`
  emits it in both modes.
- Key Discoveries row shape: all nine `Row` fields present as named.
- `Row.Guard` is `Guard string`; `OpExists`/`LiteralTrue`/`LiteralFalse` return
  zero hits repo-wide.
- Existing Infrastructure Audit: `internal/` holds only `cli`, `resolve`,
  `version` — no normalizer/table package.
- Round-Trip / Inverse Invariant cross-references to RDR 0007's payload tuple
  and RDR 0009 A8's positional-tiebreak prohibition are accurate and mutually
  consistent at their own arities.
