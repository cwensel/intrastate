# Recommendation 0002: Transition Table As Reviewable Data

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Implemented
- **Type**: Architecture
- **Profile**: foundational — cross-RDR producer: the sparse TOML wire format,
  the normalization semantics that mint kernel rows (outcome binding,
  match-block expansion, block-keyed `Match`/guard routing, `RequiresOwned`,
  next-state tags and writes, existence constants, per-atom block), the
  expanded-table dump format and its total ordering, and the validation
  category taxonomy — consumed by RDRs 0003, 0004, 0006, 0007, 0008, and 0009.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel; JDR 0001
  (`docs/jdr/0001-resolve-kernel-seam.md`) §D2, §D4, §D5, §D6, §D7, §JD-3,
  §JD-9, §JD-10, §JD-15, §JD-16, §JD-17
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A flow author needs every legal edge to live in one reviewable artifact, so "is this transition legal" has one answer instead of being reconstructed from scattered prose. The system-internal requirement is to choose a table representation that can encode match predicates, multi-tag writes, and accessor references for both RDR and kata flows.

## Context

### Background

The transition table is the flow's design and must be hand-authored from the legal graph audits, not generated. The RDR table and kata table are the test workload, including multi-tag writes that a scalar state token cannot represent.

The real design fork is representation: TOML, JSON, CSV, or a small DSL, with trade-offs among diff reviewability, parse simplicity, expressiveness, and lintability.

### Technical Environment

intrastate is a Go CLI wired through `internal/cli`. The table format feeds the resolver kernel and the static lint, and must remain practical to parse, validate, and review in normal code review. The kernel this RDR produces rows for is implemented (`internal/resolve/resolve.go::Row`, `::Resolve`); it has no production consumer yet, so the normalizer is being specified against test-covered code rather than a shipped surface.

## Research Findings

### Investigation

The proposal is shaped by RDR 0001's stateless resolver split, the current CLI
output contract, and the transition-model prior that treats state as a tag-set
matched by reviewable rows. The design has to keep table authors in normal code
review instead of making them read generated code or scattered prose.
Sibling-path check for an existing table/selection signal:

```sh
rg -n "resolve|resolver|transition|state|guard|tag|predicate|recognized|outcome|next legal|illegal|refus" internal cmd docs
```

The search found no implemented transition table or normalizer package under
`internal/`; it found the kernel (`internal/resolve`), the existing CLI refusal
plumbing, and the peer RDR drafts. Prior-art reading favors a data table that
looks like "match predicates -> tag writes" and reserves external FSM libraries
for graph validation, not runtime orchestration.

Direct `arc` searches over `StateMachineOS`, `StateMachineLit`, and `DevRef`
added four design constraints. Sismic keeps nested state source separate from
rendered PlantUML output and carries transition contracts as preconditions,
postconditions, and invariants. `transitions` models guards as positive
`conditions` and negative `unless` lists evaluated as one predicate set.
Stateless builds a symbolic `StateGraph` for diagramming from machine metadata,
including superstates, stay transitions, and decision nodes, instead of making
the graph renderer the runtime. The literature search surfaced statecharts as
the formal answer to dimensional state explosion: hierarchy and extended state
factor common context instead of enumerating every Cartesian row.

### Key Discoveries

- **Documented** — RDR 0001 delegates the reviewable transition-table contract
  to this RDR and requires enough structure for deterministic single-edge
  selection.
- **Documented** — the CLI output contract already establishes refusal-first
  behavior; table parse and lint failures can flow through the existing
  structured error gateway rather than inventing table-specific output.
- **Documented** — the kernel row (`internal/resolve/resolve.go::Row`) carries
  `RuleID`, `SourceLocator`, `Outcome`, `Match`, `RequiresOwned`, a guard,
  `Writes`, and `Escape`; `Resolve` refuses `unmodeled_outcome` before any row
  is consulted and filters candidates on `Row.Outcome`, so every normalized row
  must bind one outcome from the alphabet.
- **Verified** — Go's TOML tooling can preserve a sparse authoring schema's
  ergonomics for representative malformed-row diagnostics; the Resolve spike
  identified root-key placement as the exact field-layout constraint to lock.
- **Verified** — the RDR and kata legal graphs can be expressed as sparse rules
  that normalize to explicit row candidates, including explicit escape rows for
  modeled `no_match` behavior, without requiring host-code callbacks or an
  embedded expression language.
- **Verified** — hierarchical/shared contexts plus positive/negative guard lists
  are enough factoring to avoid RDR's status/profile/prelock Cartesian explosion
  without importing a full statechart runtime.

### Critical Assumptions

- **A1 TOML can represent sparse transition rules with nested match predicates,
  multi-tag writes, per-capability accessor tables with `keys` bindings, root
  `[initial]`/`terminal` declarations, the seven type-model keys, and
  `[model.metadata]` under strict decoding without ambiguous decoding.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: The §D7 closed layout is witnessed end to end. `cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/iter-2 && GOCACHE=/private/tmp/intrastate-rdr0002-gocache GOFLAGS=-mod=mod GOPROXY=off go run . rdr-fixture.toml kata-fixture.toml` decodes both fixtures with `github.com/pelletier/go-toml/v2` v2.3.1 under strict decoding and normalizes them to 9 and 2 rows (`evidence/spikes/iter-2/output.txt`, the normative fixture approved this stage). Each element the claim names is exercised: inherited match contexts (`large-prelock` → `prelock` → `draft`); multi-tag writes plus an explicit clear (`reconcile-rewind`); the three capability tables with `keys`, `rdr-status` appearing under both `[read.*]` and `[write.*]` since RDR 0004's identity is `(flow, name, capability)`; root `[initial]` (three owned assignments) and `terminal` resolving to a context; all seven type-model keys across the two fixtures; and `[model.metadata]` carried through uninterpreted. Strict decoding is real, not assumed: `-strict-only` reports `STRICT-OK` on both, `neg/neg-unknown-field.toml` refuses `unknown schema field: model.flavor`, and the pre-§D7 fixtures are now themselves refused (`unknown schema field: accessors.*, tags.*.accessor`) — the contracts reject the shapes they author. The API is confirmed in the module cache rather than from memory: `(*toml.Decoder).DisallowUnknownFields()` (`unmarshaler.go:52`) yields `*toml.StrictMissingError` whose `DecodeError.Key()` names the offending key (`errors.go:69`), `toml.Unmarshal` is permissive by default (`unmarshaler.go:21`) — which is why strictness is an obligation on this format, not a library property — a two-pass decode over one `[]byte` is re-runnable (so the version gate can precede strict decoding), and the strict check descends only into struct branches, leaving a `map[string]any` field untouched, which is what makes `[model.metadata]` free-form while the rest stays strict. 41 negative fixtures refuse one category each, 2 further `probe-*` fixtures witness the unauthorable-provenance enumeration (not category controls), and 2 positive fixtures load, alongside 2 `-strict-only` checks (`evidence/spikes/iter-2/negative-cases.txt`, regenerated at Stage 6). The 41 is the count Testing Strategy scenario 3 asserts; the two probes are deliberately excluded from it. It was 35 before Stage 6 added six guard-block controls (A15); those six deepen existing categories rather than adding new ones, so the **category** coverage below is unchanged at 18 of 25.
  - **If wrong**: The chosen carrier either loses table semantics or forces a
    custom parser earlier than intended.
- **A2 Row order is not part of successful edge selection.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: The shipped kernel already implements gate-then-count, so
    the ordering is verified against code rather than merely decided.
    `internal/resolve/resolve.go::Resolve` runs
    `selected, blocked := gate(candidates, in.Guards, view)` and returns on
    `blocked != nil` **before** `switch len(selected)`;
    `::gate` prunes `GuardFalse`, then reports `missingOwned`, then vetoes on
    an undecidable survivor; `::escapeOrRefuse` delegates to the same `gate`
    before counting `viable`. `::missingOwned` sorts with `slices.Sort` and
    documents that "Row order is a normalization detail RDR 0002 owns and
    must not reach the reported diagnosis";
    `adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`
    freezes the property. The Normative Contracts below restate this flow and
    explicitly reject first-match semantics and count-first pruning
    (JDR 0001 §D2, whose "Lands in 0002" clause this discharges; RDR 0007
    states the same ordering as the kernel's and verifies it at its A20).
  - **If wrong**: Reviewers would have to reason about hidden priority, and
    reordering rows could silently change resolver behavior.
- **A3 RDR and kata flow edges can be encoded sparsely with fixed predicate
  operators rather than host-code callbacks or Cartesian-product row
  enumeration.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/iter-2 && GOCACHE=/private/tmp/intrastate-rdr0002-gocache GOFLAGS=-mod=mod GOPROXY=off go run . rdr-fixture.toml kata-fixture.toml` parsed and normalized the RDR and kata fixtures under the §D7 layout. Between them the fixtures cover status, profile, stage, prelock iteration, equality, set membership (`in` on `recognized` and on `profile`), integer comparison (`lt`), existence (`exists = false`), a boolean guard, a set-valued write (`labels`), a terminal-context rule, rewind with an explicit clear, positive and negative guards, multi-tag writes, gate references, and an explicit `escape = ["no_match"]` rule that itself expands over two alphabet members — all without host-code callbacks or an embedded expression language. `iter-2/output.txt` captures the eleven expanded candidate rows.
  - **If wrong**: RDR 0003 must expand the predicate grammar or this table format
    becomes too weak for the target flows.
- **A4 The model data can carry enough provenance to separate owned, observed,
  and recognized tags for lint and accessor binding.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0003 `Technical Design` requires declared tag provenance for predicate lint and owned-tag read-before-write checks; RDR 0004 `Technical Design` makes RDR 0002 responsible for table-carried accessor references while accessor execution validates capabilities and writes only owned tags.
  - **If wrong**: The resolver cannot prove read-before-write or accessor safety
    from the model alone.
- **A5 Standard parse and validation failures can be surfaced through the
  existing CLI error envelope.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr/clierr.go::CLIError` already carries stable `Code`, `Message`, optional `Param`, diagnostic `Detail`, `Hint`, exit-code `Group`, and `Cause`; `internal/cli/respond/respond.go::Fail` emits that envelope in text/json modes. The envelope is the verified half. **No parse-validation code has yet travelled it:** `internal/cli/config/config.go::Load` emits only `config-not-found` and `config-read-error`, and its `config-invalid` exists solely as a doc comment and a `TODO` — `Load` reads bytes, stamps `SchemaVersion` unconditionally without reading the file's own version, and returns a nil error. The codes emitted anywhere in the CLI are exhaustively `config-not-found`, `config-read-error`, `flag-invalid-value`, and `command-error`. This RDR's table parse and lint categories would be the **first** parse-validation codes in the envelope, so their code names, `Group` assignments, and exit mapping are new design work at implementation, not a pattern to copy.
  - **If wrong**: Table loading would need a separate user-facing error contract
    owned by this RDR or RDR 0005.
- **A6 Shared contexts and positive/negative guard lists are sufficient to keep
  the RDR model sparse without hiding ambiguity.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/iter-2/rdr-fixture.toml` encodes `draft -> prelock -> large-prelock` inherited contexts plus `all.iter.lt = 3`, `all.finalized_at.exists = false`, and `unless.profile.eq = "small"` guards; `iter-2/output.txt` shows `continue-prelock` normalizing to rows carrying the inherited status/stage/profile predicates alongside the combined `all`/`unless` predicates, each atom tagged with its authored block, so no constraint is hidden by the factoring. The `continue-prelock` / `continue-prelock-cluster` sibling pair demonstrates the same at the ambiguity boundary: two rules share every inherited predicate and bind one outcome, differing only in their guards, and both survive normalization as distinct rows for RDR 0006 to adjudicate rather than being silently merged.
  - **If wrong**: The model either needs a richer statechart-like hierarchy or
    the table becomes too repetitive for reliable human review.
- **A7 Deterministic expanded-table ordering is a format contract, not an
  implementation accident.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: Normative Contracts define the expanded table as the
    normalized candidate-row value and fix the row order as the identity
    tuple `(model id, rule id, expansion suffix)` compared field by field,
    byte-lexicographically, with the source locator excluded; the same field
    list is what the Round-Trip invariant preserves. The normalizer spike
    witnesses it, including the generalized match-block expansion and its
    sequence-valued suffix (re-verified this stage with A1):
    `cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/iter-2 &&
    GOCACHE=/private/tmp/intrastate-rdr0002-gocache GOFLAGS=-mod=mod GOPROXY=off
    go run . rdr-fixture.toml kata-fixture.toml`
    emits every field of the row value — identity with `#suffix` on expanded
    rows, locator, kind, lifted outcome, per-atom `@block`, next tags, writes
    including `<clear>`, derived `requires_owned`, gate ids, escape classes —
    with atoms sorted by (key, block, operator, literal). Expansion is
    witnessed on a non-`recognized` tag: `continue-prelock` expands on the
    *inherited* `profile.in = ["large","foundational"]` into
    `#large` / `#foundational`, so the suffix is a sequence and compares
    element by element (`slices.Compare`), the empty suffix sorting first.
    Three consecutive runs produced byte-identical output despite Go's
    randomized map iteration (SHA-256
    `6ccfe9012b0705ef4d4b3d1c620daffd69523436175120be1bea8a05df9c55dd`,
    `evidence/spikes/iter-2/output.txt`), and three semantics-preserving
    permutations of the RDR fixture — TOML key order, `[[rule]]` declaration
    order, and `eq = "x"` versus `in = ["x"]` — all dumped byte-identically
    (`evidence/spikes/iter-2/negative-cases.txt`). Rule-order permutation is
    the control that excludes a positional tiebreak. Byte-lexicographic comparison is
    a Go spec guarantee ("Two string values are compared lexically
    byte-wise"), so the ordering is locale- and platform-independent. This
    explicitly rejects source-order, map-iteration, and renderer-specific
    ordering, and rejects ordering on the source locator, which may carry
    line/column detail that changes when unrelated source text is edited.
  - **If wrong**: Golden tests and review dumps could churn across machines or
    refactors even when the transition semantics are unchanged.
- **A8 Every normalized row can bind exactly one outcome from the alphabet, so
  the kernel's `Row.Outcome` filter admits the row.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::Resolve` skips any row whose
    `Outcome` differs from `Input.Recognized` (`if row.Outcome !=
    in.Recognized { continue }`), before the match check;
    `::escapeOrRefuse` applies the same filter to escape rows; `::assemble`
    binds the outcome under `const recognizedTagKey = "recognized"` only when
    `in.Recognized != ""`; and `::Resolve` refuses `unmodeled_outcome` via
    `!in.Table.models(in.Recognized)` before any row is consulted. So a row
    carrying the lifted outcome is admitted, and with a non-empty
    `Input.Recognized` an outcome-less row never matches on either path.
    **Residual, and the reason this RDR is load-bearing for it:** the kernel
    validates no alphabet well-formedness — `::Table.models` is a plain
    `slices.Contains` — so with `Input.Recognized == ""` *and* the empty
    string in the alphabet, a row with an empty `Outcome` matches and emits a
    plan against a view carrying no `recognized` key. The recognized-totality
    clause below is the sole barrier, and JDR 0001 §JD-10 ratifies it: the
    `recognized` tag is total over matching.

    **The barrier's reach is bounded, and this is a known residual.** The
    alphabet ban is a *load-time* check in the normalizer; the kernel does not
    import the normalizer and validates no alphabet well-formedness itself, so
    the barrier covers exactly one producer — tables this RDR's loader built.
    Any other construction of a `resolve.Table` (a test helper, a future
    programmatic producer, a caller assembling rows directly) reaches the same
    unguarded path, and so does a caller passing `Input.Recognized == ""` from
    the CLI edge, which is RDR 0005's boundary and not reachable from here.
    Closing it for all producers means a kernel-side check, which is RDR 0001's
    to make. This RDR closes what it produces and states the rest as residual
    rather than implying table-wide coverage it cannot enforce.
  - **If wrong**: Normalized rows never match (every resolve refuses
    `no_match`), or the kernel needs an eventless-row concept this RDR does not
    model.
- **A9 The normalized candidate row can carry one unified, block-retaining atom
  set and still produce a conforming `resolve.Row` at the kernel handoff.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Owed at implementation; the block it waits on is re-verified
    at source. The block-keyed `Match`/guard split (JDR 0001 §D6) is pinned to
    the kernel handoff so the normalized value keeps each atom's authored block
    (`match`/`all`/`unless`), which `Guard string` cannot carry — and
    `internal/resolve/resolve.go::Row` still carries `Guard string` and
    `Match []Tag`, with no `Atom` type, no operator-token type, and no per-atom
    block field anywhere in non-evidence Go source. So the field this assumption
    targets does not yet exist, and the block is structural rather than a
    scheduling preference. The normalizer half is witnessed: the spike carries
    each atom's authored block through normalization and renders it
    (`profile.eq=small@unless`, `cluster_ready.eq=true@all`,
    `stage.eq=prelock@match` in `evidence/spikes/iter-2/output.txt`).
    Verification: the Testing Strategy scenario 2 normalization assertions plus
    scenario 4's `Resolve` run must both pass against one normalized value —
    the first reading the unified atom set, the second consuming the split row.
    Blocked on RDR 0007's kernel reshape (Prerequisites), like the rest of
    Phase 2.
  - **If wrong**: Either block retention is lost (breaking RDR 0003 downstream
    and the Round-Trip field list) or the kernel handoff needs a shape this RDR
    does not model.
- **A10 `duplicate model id` is decidable only across documents, so no
  single-document loader can or should report it.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: The decidability half is verified at source and survives the
    §D7 re-authoring: `evidence/spikes/iter-2/main.go::Model` declares `Model`
    as a **singular struct field** of type `ModelMeta`, not a slice, mirroring
    `[model]` as a singular TOML table (`evidence/spikes/iter-2/rdr-fixture.toml`
    authors one `id = "rdr"` per document; `[[rule]]` is the only array-of-tables
    in the closed layout), so a single-document load structurally cannot observe
    two model ids. The
    category is therefore scoped to a multi-document dump/lint invocation. The
    behavioral half is owed: the scenario 3 fixture for this category is a
    **pair** of documents sharing a `model id`, and a single-document load of
    either one must not report it.
  - **If wrong**: The category is unreachable as specified, or a multi-document
    surface exists that this RDR has not named.
- **A11 An escape rule binding its outcome with an `in` atom expands into one
  rescuing row per member, each rescuing only its own outcome.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: Both halves now hold. The **kernel half** is Source Search:
    `internal/resolve/resolve.go::escapeOrRefuse` filters escape candidates on
    `row.Outcome != in.Recognized || !view.matches(row.Match)` exactly as
    `::Resolve` does for ordinary rows, so a per-outcome escape row is what the
    kernel admits and there is no table-wide catch-all. The **normalizer half**
    was owed and is now witnessed: the critique lens introduced the clause
    against a fixture whose escape rule bound `eq = "round-clean"`, and the
    §D7 re-authoring changed it to bind `in` over two alphabet members.
    `evidence/spikes/iter-2/rdr-fixture.toml` `draft-no-match-escape` authors
    `[rule.match.recognized] in = ["round-clean", "reconcile-block"]` and
    normalizes to exactly two escape rows —
    `rdr.draft-no-match-escape#round-clean` and
    `…#reconcile-block` — each carrying a distinct expansion suffix, its own
    bound outcome, `write=[]`, `next=[]`, `requires_owned=[]`, and the modeled
    class list `escape=[no_match]` (`evidence/spikes/iter-2/output.txt`). So an
    escape rule expands like any other rule, and one authored rule covers
    several outcomes. The MVV assertion in Testing Strategy scenario 2 remains
    as the implementation-side control.
  - **If wrong**: One authored escape rule cannot cover several outcomes, so
    covering an alphabet of N outcomes requires N hand-authored escape rules
    and the "expand like any other rule" clause overstates the format.
- **A12 The kernel handoff routes every atom to exactly one of `Match` /
  `Guard`, exhaustively and disjointly.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Owed at implementation; its discriminating witnesses now
    exist in the fixture, which they did not before. A9 asserts the unified
    atom set *can* produce a conforming `resolve.Row`; it does not assert the
    routing is total. Under the block-keyed rule (JDR 0001 §D6) totality is
    structural — every atom carries exactly one of three blocks — but the
    handoff code is the enforcement point and lands in no assertion yet.
    Verification: a scenario 2 assertion that for the RDR fixture the union of
    the constructed row's `Match` and guard atoms equals the normalized atom
    set and the intersection is empty. The discriminating witnesses are both
    authored and normalized: `evidence/spikes/iter-2/rdr-fixture.toml`'s
    `continue-prelock-cluster` carries `[rule.guard.all.cluster_ready] eq =
    true` — an `eq` atom under `guard.all` over an optional (`required` unset)
    owned key — and the kata fixture carries `status.eq=closed@unless`; both
    render as guard-block atoms in `iter-2/output.txt`, routing by block rather
    than operator. Blocked on RDR 0007's kernel reshape (Prerequisites), like
    A9.
  - **If wrong**: An atom is silently dropped at the handoff or double-counted,
    which changes which rows match without any dump or normalization test
    observing it.
- **A13 Merging is idempotent on identical atoms: one atom contributed by a
  rule and by one or more inherited contexts collapses to exactly one — and no
  member sequence, atom literal or set-valued write value, is compared by a
  joined rendering.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: Closed at Stage 6 by fixing the spike and adding the
    discriminating witnesses it lacked; transcript in
    `evidence/spikes/iter-2/negative-cases.txt` (§Delimiter-bearing controls)
    and `evidence/reconcile/iter-2/a13-a15-close.md`. Until this pass the spike
    was keyed on a joined rendering (`main.go::Atom.identity` → `literalString`,
    `strings.Join(a.Literal, ",")`) and joined set-valued writes
    (`::renderValue`), and **every** existing witness used a single-member
    literal — under which a joined rendering and a member sequence are the same
    string, so none of them could tell the two candidate keys apart.
    **The fix**: `Atom.identity` now encodes the literal as a length-prefixed
    member sequence (`::memberKey`, `%d:%s` per member, injective for any
    member content), `Write.Value` became a `[]string` member sequence with a
    matching `Write.identity`, and `literalString` is retained for **display
    only**, with a comment forbidding its use in any identity. **The
    discriminating witnesses**, both of which fail against the pre-fix spike:
    two contexts contribute
    `in = ["a<d>b", "c"]` and `in = ["a", "b<d>c"]` on one key and block (on
    `finalized_at`, `kind = "string"` with no declared domain, so the members
    are authorable) — **two atoms survive** and the rule stays the dead rule
    the author wrote, where the pre-fix spike collapsed them to one and
    normalized a **live** expanding row. The write arm spells one set write as
    `["x<d>y","z"]` and `["x","y<d>z"]`, which now normalize to **different**
    values. Both arms are **parameterized over five delimiters** (`,` `;` `|`
    space, and the empty string), as Testing Strategy 2 requires, and that
    parameterization earned its keep twice: the pre-fix spike joined on `,` and
    collapsed **only** the comma case — all four other delimiters yielded two
    atoms while the defect was live, so a single-delimiter control would have
    passed it; and the `|` case then caught a residual in the *render* (an
    unquoted `|` separator spelled `["x|y","z"]` and `["x","y|z"]` alike, so a
    test asserting on rendered text would have passed the collision even with a
    correct identity), fixed by quoting each member (`::renderMembers`). All
    fifteen fixtures are minted by `gen-cases.py`, so they belong to the
    promoted set. **No regression**: the baseline dump is byte-identical to
    the recorded digest `6ccfe901…` across three runs, all 43 negative fixtures
    still refuse one category each, all three permutations still digest
    identically, and both prior merge properties hold (idempotent → one atom,
    distinct → two, `use`-order independent). The fix is behavior-preserving on
    every fixture whose members do not contain the delimiter, which is why the
    digest is unchanged. **Idempotence half:**
    `evidence/spikes/iter-2/merge-idempotent.toml` has `reconcile-rewind`
    author `[rule.match.status] eq = "Draft"` while its inherited `draft`
    context contributes the byte-identical atom; the normalized row carries
    `atoms=[status.eq=Draft@match]` — count one.
    **Distinct-literal mirror:** `merge-distinct.toml` adds a second context
    contributing `status.eq=Final@match`, and the same row carries
    `atoms=[status.eq=Draft@match; status.eq=Final@match]` — both survive; the
    reversed `use` order (`merge-distinct-rev.toml`) digests identically
    (`3813472915…`), so the result does not depend on `use` iteration order.
    A merge keyed on `(block, key, operator)` would collapse the second pair
    and fail; one keyed on nothing would inflate the first. A merge keyed on a
    joined rendering passes both, which is why the delimiter-bearing witness is
    owed. Transcript in `evidence/spikes/iter-2/negative-cases.txt`. The
    scenario 2 assertions remain as the implementation-side controls.
    **Widened at 3amigo iter-3** to cover set-valued *write* values: the spike
    joins them too (`main.go::renderValue` → `strings.Join(parts, ",")`), and
    because RDR 0004's read-back compares the held value for equality, that
    collision decides whether a write reports as applied. Same defect, same
    owed witness shape (two spellings must stay two values), one assumption —
    the subject is "a member sequence is never compared as a joined string",
    not "the merge map key".
  - **If wrong**: Two failure arms, of unequal severity. A merge keyed too
    *narrowly* carries duplicates, inflating the dump and the Round-Trip field
    list without changing match semantics — benign. A merge keyed too **widely**
    — the joined-rendering key, which the spike used until Stage 6 closed it —
    silently drops a distinct atom, and the surviving rule normalizes to a **live
    expanding row where the author wrote a self-contradictory (dead) rule**.
    That arm changes match semantics: the artifact answers "this edge is legal"
    for an edge the author did not write, which inverts the user outcome this
    RDR exists to deliver. It is also the arm that trips **zero** load
    categories, because the defect is absorbed in the merge before any check
    sees it. The desk trace's `merge atoms` row records this arm concretely.
- **A14 The version gate precedes strict field decoding, so a v2 document
  refuses as an unsupported version rather than an unknown schema field.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: Witnessed by a **v2-shaped** fixture and a paired control
    that isolates the ordering. `evidence/spikes/iter-2/neg/neg-v2-shaped.toml`
    is the RDR fixture with `version = 2` plus a v2-only root table
    (`[schema_v2]`) that the strict decoder would reject; it refuses
    `unsupported version 2`, **not** `unknown schema field`. The
    discriminating control is the same document with `version = 1` restored
    and the v2-only key left in place: it refuses `unknown schema field:
    schema_v2`. One key, two outcomes, and only the gate ordering separates
    them — a loader decoding strictly first would name `schema_v2` in both
    cases. Scenario 3's plain `unsupported version` fixture (a v1-shaped
    document with a bad version value) trips the category on either ordering
    and is retained as the category control, not the precedence one. The
    two-pass shape this requires is confirmed feasible against the library:
    `toml.Unmarshal` into a probe struct reading only `[model].version`, then
    a second strict `toml.NewDecoder(bytes.NewReader(data))` pass over the same
    `[]byte` (re-runnable — see A1).
  - **If wrong**: A future v2 file fails with a diagnostic naming an arbitrary
    field and never mentioning the version, which is the confusion the two-pass
    gate exists to prevent.

- **A15 Atom-level validation is enforced identically in `match`, `guard.all`,
  and `guard.unless` — operator membership, `<clear>`, domain/kind conformance,
  `#` reservation, and tag-key declaration.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: Closed at Stage 6 in the same pass as A13; transcript in
    `evidence/spikes/iter-2/negative-cases.txt` and
    `evidence/reconcile/iter-2/a13-a15-close.md`. **The gap, as found**: the
    §D7 spike enforced these in match blocks only — its guard walk checked
    solely that the key was declared and was not `recognized` — so
    `frobnicate = "small"`, `eq = "<clear>"`, and `eq = "NOT_A_PROFILE"` each
    loaded clean under `[rule.guard.unless.profile]` and emitted a normalized
    atom (reproduction preserved at
    `evidence/reconcile/iter-2/a15-guard-block-repro.txt`). No transcript line
    witnessed it because `gen-cases.py` mutated only match blocks, so the gap
    was invisible to `negative-cases.txt` by construction — which is why two
    prior lenses passed over it. **The fix**: `checkMatchBlock` was generalized
    to `checkAtomBlock(m, owner, blockName, block)` and the guard walk now
    routes through it, so every atom rule is enforced in all three blocks. Two
    rules stay legitimately match-only and are gated on `isMatch`: the `eq`/`in`
    operator restriction (§D6 routing) and the `#` reservation in `in` members
    (the expansion-suffix separator — only match blocks expand). Guard blocks
    admit RDR 0003's closed set, cited by member list as `guardOperators`.
    Domain/kind conformance is applied **per operator**, because operators do
    not all take a domain member as their right-hand side: `eq`/`in`/`contains`
    take members and are domain-checked, `exists` takes a bool literal, and
    `lt`/`lte`/`gt`/`gte` take an ordered *bound* that is kind-checked but not
    domain-checked (a bound need not itself be an authorable value). That
    per-operator scoping is a real design refinement this pass forced: the
    clause's "domain/kind conformance" reads as unqualified, but applying it
    literally to a `lt` bound breaks the RDR fixture's own `iter.lt = 3`.
    **The controls**: `gen-cases.py` now mints six guard-block negatives
    (37 → 43), one per atom rule per block —
    `neg-guard-{unless,all}-bad-operator`, `neg-guard-{unless,all}-clear`,
    `neg-guard-unless-out-of-domain`, `neg-guard-unless-unknown-tag` — and each
    refuses exactly the category it names. **No regression**: baseline digest
    unchanged at `6ccfe901…`, all 43 negatives refuse one category each, all
    positives load, permutations unchanged.
  - **If wrong**: The narrower reading ships (it did, until this pass). A
    misspelled operator or an
    out-of-domain literal in an `unless` block silently disables the negation,
    so a rule the author explicitly carved out applies — the artifact answers
    "this edge is legal" for an edge the author forbade, with no load refusal
    and no lint finding. The `<clear>` arm is worse: it is an unsatisfiable
    guard atom, and RDR 0007's veto on selecting a decided sibling while an
    unevaluable candidate exists turns it into a non-escapable deadlock on
    every resolve touching that outcome.

## Proposed Solution

### Approach

Use a sparse, hand-authored TOML transition model that normalizes to explicit
candidate rows. Authors write shared tag declarations, reusable match contexts,
outcome groups, guarded rules, and tag writes; the tool normalizes that source
into a full row set for lint, resolver lookup, and table dumps.
The rendered table is review support, not the source authors maintain.

Each flow has a named model with flow-level tag declarations, per-capability
accessor tables, a declared initial owned state and terminal contexts, a legal
recognized-outcome alphabet, and sparse rules that encode the match predicates
and tag writes produced by legal edges. The source model must provide enough
structure for RDR 0001's exact-one resolver contract under the gate-then-count
ordering JDR 0001 §D2 fixes (cited, not restated — see the selection clause in
Normative Contracts). Unknown outcomes and unavailable accessors remain
resolver refusals; malformed TOML or schema-invalid rules are load/lint
failures before the resolver sees a table.

The model is data, not generated code and not a runtime FSM engine. RDR 0001's
resolver consumes the normalized representation, RDR 0003 owns the fixed
predicate operators and the tag type model, RDR 0004 owns accessor execution and
read-back safety, RDR 0005 exposes the CLI, RDR 0006 owns graph lint, RDR 0007
owns the guard seam and domain rule, RDR 0008 owns the reserved recognized tag
key, and RDR 0009 owns escape-row shape conformance. This RDR owns the on-disk
sparse representation, the normalizer that produces kernel rows from it, and the
normalized expanded-table view those peers consume.

### Technical Design

The source schema has eight conceptual parts:

1. Flow metadata: table id, version, human description, and a free-form
   `[model.metadata]` extension table the loader carries untouched.
2. Tag declarations: tag name, provenance (`owned`, `observed`, `recognized`),
   and the RDR 0003 type model spelled as the wire keys `kind`, `domain`,
   `min`, `max`, `elements`, `single_valued`, `required`. Provenance,
   authoring location, and key spelling are this RDR's; what the keys mean is
   RDR 0003's and is cited, not restated. The recognized-provenance
   declaration is named `recognized` (RDR 0008). A declaration carries no
   accessor reference — the binding runs the other way (part 4).
3. Recognized outcome alphabet: the closed set of outcome tags that the
   recognizer may emit for the flow.
4. Accessor tables: one table per capability — `[read.<id>]`, `[write.<id>]`,
   `[gate.<id>]` — each naming the tag `keys` it serves; `keys` is the only
   binding between tags and accessors. RDR 0004 owns what the entries mean and
   how they execute.
5. Root and stop set: `[initial]`, a table of owned `tag = value` assignments
   fixing the root node, and root `terminal`, a list of context ids naming the
   stop set — consumed by RDR 0006's reachability and dead-end checks.
6. Shared match contexts: named predicate blocks for dimensions such as RDR
   `Status`, `Profile`, stage, prelock iteration, cluster eligibility, or kata
   lifecycle. Contexts may inherit from another context to model statechart-like
   hierarchy without adopting a statechart runtime.
7. Sparse transition rules: reviewable rule ids, optional context references, a
   local match block, positive `all` guards, negative `unless` guards, an
   optional list of gate accessors the plan must pass, and either a write block
   containing one or more tag assignments plus an optional explicit clear list,
   or an explicit `escape` list naming the resolver failure class the row
   models.
8. Render settings: deterministic ordering and field selection for the expanded
   table dump.

The parser turns TOML into typed source data, then normalizes it into explicit
candidate rows. Load-time validation rejects malformed tag declarations,
unsupported model versions, malformed predicate atoms (unknown operator; a
literal ill-formed for its operator, including any existence literal other than
the kernel's two boolean forms; a literal outside the tag's declared domain; or
an operator other than `eq`/`in` under a match block), the reserved `<clear>`
sentinel authored as a tag value, writes to non-owned tags, rules that match on
missing tag declarations, unresolvable or cyclic context references, unknown
accessor ids, accessor `keys` that bind a tag to zero or several readers or
writers, malformed escape declarations, a rule that binds zero or more than one
outcome, a malformed `[initial]` assignment, and a recognized-provenance
declaration not named `recognized`.

Load and lint split on **arity, not severity**: a check decidable from one rule
plus the model's declarations is load-time and this RDR's; a check that must
compare normalized rows against each other — overlap, gap, dead row,
read-before-write — is graph lint and RDR 0006's. Ambiguous overlap between
candidate rows is therefore a lint finding, not a load failure.

Runtime matching is deliberately priority-free and follows the kernel's
gate-then-count ordering (JDR 0001 §D2; RDR 0007 is its normative home). This
RDR's part is the row shape that ordering evaluates: a match pattern the kernel
tests as equality, a guard predicate handed to the seam, one bound outcome, and
— for escape rows — the failure class the row rescues. Multi-tag writes are
first-class because RDR rewinds and kata lifecycle moves need to set both the
next stage/state and side-channel scope tags in one edge.

Guard factoring follows the `transitions` prior art: positive predicates and
negative predicates are authored separately and normalized into one predicate
set, in which each atom still records the block it was authored in. Contract
factoring follows the Sismic prior art: entry preconditions, postconditions,
and invariants are different validation classes, not free-form comments.
Rendered dumps follow the Stateless/Sismic export pattern: they are symbolic
views derived from model metadata and must carry enough source ids to send
diagnostics back to the authored sparse rule.

The normalized candidate row is the kernel row up to the predicate split — the
kernel's `Match`/guard routing is applied at the handoff (Normative
Contracts), not stored in the normalized value. It carries the source identity
(`(model id, rule id)` plus any expansion suffix, and the source locator), the
single outcome the row responds to, the predicate atoms (key, operator token,
literal, block), the gate accessor ids the rule names, the next-state tags and
the writes, both including rendered `<clear>` entries, the required-owned key
set derived from those writes, and the escape failure-class list. The
identity tuple `(model id, rule id, expansion suffix)` is what the dump sorts
on; the source locator is diagnostic and does not participate in that order.
The internal representation may be
indexed as a decision tree, trie, or decision DAG for efficient lookup, but that
is an implementation detail; the normative semantic object is the normalized
candidate-row set plus its source locator back to the sparse TOML rule. The
locator must identify at least the model id and rule id, and it is **derived,
not authored**: the loader composes it from `(model id, rule id)`, which the
identity tuple already makes unique, optionally extended with byte line/column
coordinates as diagnostic detail. An authored `[[rule]].source` is a provenance
*annotation* carried alongside it, not the locator itself — it is neither
required to be unique nor consulted by any comparison. (The iter-2 spike
populates `Locator` from `rule.Source` directly, which is why both
`kata-fixture.toml` rules carry one identical locator; that is a spike defect
against this clause, not a permitted reading.) The Round-Trip comparison of the
locator's "rule-identifying part" is therefore well-defined: it is the
`(model id, rule id)` projection, excluding any line/column suffix.
**Both consumers were checked at Stage 6 and a derived locator serves them.**
RDR 0006 takes the locator only as diagnostic payload and sorts findings on
model id / invariant code / **source rule id**, never on the locator
(`0006-graph-lint-authority-and-guarantees.md:983-985`; its A1 wants the locator
so lint can "locate the authored rule", `:108-120`) — all satisfied by a
`(model id, rule id)` derivation. RDR 0004 names no locator at all and keys no
accessor diagnostic on one, so it imposes no constraint here.

#### Normative Contracts

```normative
The transition model MUST be authored as sparse TOML data, not generated code
and not a fully expanded Cartesian-product table.
```

```normative
The source schema is the closed layout JDR 0001 §D7 fixes: root `outcomes`,
root `terminal`, `[model]` (with the free-form sub-table `[model.metadata]`),
`[initial]`, `[tags.<tag>]`, `[read.<id>]`, `[write.<id>]`, `[gate.<id>]`,
`[context.<id>]`, `[[rule]]`, and `[dump]`. Context predicates live under
`[context.<id>.match.<tag>]`; rule predicates live under `[rule.match.<tag>]`,
`[rule.guard.all.<tag>]`, and `[rule.guard.unless.<tag>]`; writes live under
`[rule.write]`; explicit clears live in a rule-level `clear` list; gate
accessor references live in a rule-level `gate` list of `[gate.<id>]` ids;
modeled escape rows live in a rule-level `escape` list. `[model]` additionally
carries an optional human `description`, and each `[[rule]]` an optional
`source` — a provenance *annotation* carried alongside the locator, which is
derived from `(model id, rule id)` and never from `source` (below).
Both are admitted keys; the enumeration is exhaustive over what the normative
fixtures author, because a layout that refuses its own fixtures is not the
layout. No other root key or table is admitted (strict decoding, below).

**`[model.metadata]` is the one sanctioned extension namespace.** The loader
MUST decode it as a free-form table, carry it through to the normalized model
untouched, and MUST NOT interpret any key in it; strictness applies everywhere
else. It is a **model-level field and reaches no candidate row**, so it is
outside the dump's field list and outside the Round-Trip invariant, which
compares candidate-row sets: two documents differing only in `[model.metadata]`
normalize to identical row sets while carrying different normalized models, and
that is intended — the invariant is over rows, and metadata is carried for the
caller, not for row semantics. Nothing here forbids a dump from rendering it as
a header; that is presentation, like `[dump]`. Its internal shape is
deliberately unconstrained (arbitrary keys, values, and nesting): constraining
it would make it a schema, which is the opposite of an extension namespace.

**It has no consumer in this cluster today, and that is the point of reserving
it.** No peer RDR reads it (RDR 0006 consumes the normalized model and does not
reference metadata); JDR 0001 §D7(v) sanctions the namespace without naming a
tool. Reserving a namespace before its first consumer is what keeps that
consumer from inventing a root key that strict decoding would refuse. Two
obligations follow, so "carried through untouched" is a testable claim rather
than a decode-and-discard: the normalized model MUST expose the decoded table as
a field (Phase 2 deliverable), and scenario 1 MUST assert **value and nesting
equality** against the authored table, not merely that its top-level key names
survive — a loader discarding every value passes a key-name oracle. Because the
table is outside the dump field list and outside the Round-Trip invariant, this
assertion is the only comparison that observes it.

The cost of an unconstrained namespace is drift: it is the one place in a strict
schema where anything loads, so it attracts content that should have been
modeled — ownership, feature flags, a homegrown status field — and no load
refusal or lint finding will ever object. That is accepted rather than solved
(constraining it would make it a schema), so the control is scope, stated here
for the reviewer rather than the loader: **`[model.metadata]` is for data no
consumer in this cluster reads.** Anything a rule, lint, or the resolver must
act on is a tag, not metadata. A reviewer seeing a metadata key that some tool
has begun to branch on should read that as the signal to model it properly — and
the first such consumer is the trigger to revisit this clause, since a read key
is no longer an extension namespace.

**Accessor tables.** Each capability table decodes to its own shape, so
"exactly one capability per accessor" (RDR 0004) holds structurally rather than
by validation. Every entry carries `role`, `path`, `keys` (a non-empty list of
declared tag keys), and `timeout` (a Go duration string). **Each of the four is
required, and any of them absent, empty, or ill-formed is a `malformed accessor
declaration`** — an absent or empty `role` or `path`, an absent or empty `keys`
list, a `keys` member naming an undeclared tag (`unknown tag`), an absent
`timeout`, a `timeout` that does not parse as a Go duration, or one that parses
non-positive. The rule is stated per field because "carries" alone left three of
the four with no refusal. A `[write.<id>]` entry additionally carries
`read_back = true`; `read_back = false` or its absence on a write entry is the
same category, since RDR 0004's read-back is not optional for a writer. The same id MAY appear under two
capability tables — RDR 0004's identity is `(flow, name, capability)`. RDR 0004
owns what `role`, `timeout`, and `read_back` mean and how an accessor executes;
this RDR owns the layout and the **binding validations**, which are
provenance-scoped:

- every key in a `keys` list MUST be a declared tag (`unknown tag`), and MUST
  NOT be `recognized` — the kernel supplies that key and no accessor reads or
  writes it (RDR 0008);
- every **owned** tag MUST be served by exactly one reader; an **observed** tag
  MAY be served by at most one reader, and zero is legal because an observed
  key may arrive from the caller (JDR 0001 §JD-9) — refusing it would make
  every `--tag`-supplied key unauthorable. The asymmetry is deliberate and the
  owned arm admits no write-only tag: `RequiresOwned` carries every key a row
  writes or clears, and the kernel evaluates it against the view
  (`resolve.go::missingOwned`), so an owned tag written on one transition must
  be readable on the next or that resolve refuses `owned_state_unavailable`. A
  write-only owned tag — an audit stamp, a one-way flag — is therefore not
  merely unbound but unusable, and the load refusal is the early, diagnosable
  form of a failure that would otherwise surface as a refusal at resolve time.
  Genuinely write-only state does not belong in the owned tag set;
- every key named by any rule's write block or clear list, and every key
  assigned in `[initial]`, MUST be served by exactly one writer
  (`malformed accessor binding`);
- every key in a `[write.<id>]` entry's `keys` list MUST be an **owned** tag
  (`write to non-owned tag`);
- every id in a rule's `gate` list MUST resolve to a `[gate.<id>]` entry
  (`unknown accessor`).

A violation of the second or third bullet is a `malformed accessor binding`;
the fourth carries `write to non-owned tag`. The writer-arity and
writer-provenance obligations are stated as **separate bullets with separate
categories on purpose**: they are independent predicates, and a single bullet
naming both left it undecidable which category a violation carries — which
matters because the Testing Strategy asserts on the category, not the message.
This settles the provenance scope §D7(ii) left open at this RDR.

**Gate references.** A rule's `gate` list names the gate accessors whose
`allow` the executor requires before applying that rule's plan; the list is
carried on the normalized row and is part of its value. When it runs, and what
`deny` and `indeterminate` mean, are RDR 0004's. An escape rule MUST NOT carry
a `gate` list (it yields no plan to gate) — `malformed escape declaration`.
This fills the site §D7 left blank for this RDR.

**Root and stop set.** `[initial]` is a table of `tag = value` assignments;
every key MUST be a declared owned tag and every value MUST be well-formed for
that tag's declared kind and domain (`malformed initial declaration` otherwise;
the `<clear>` sentinel is refused under the reserved-value rule).

An **observed or recognized key in `[initial]` has no authorable form**, so this
category's provenance arm is unreachable by construction. The writer-binding
bullets above already require every `[initial]` key to be served by exactly one
writer, and every writer key to be owned; a non-owned key therefore refuses as
`malformed accessor binding` when no writer serves it, or as `write to
non-owned tag` when one does, before the owned-tag predicate here is consulted.
Those two authorings are exhaustive. `malformed initial declaration` is
consequently witnessed by its **value** arm — a literal ill-formed for the
declared kind or domain — and the Testing Strategy carries no observed-key
control for it. This is Principle 8 working (JDR 0001: the layout makes illegal
states unrepresentable), not a gap: the state is refused at load either way, and
the same drafting applies as to the absent-`recognized` rule below, which is
likewise unauthorable. Root
`terminal` is a list of context ids, each of which MUST resolve (`unknown
context`); a rule MAY `use` a terminal context, which is how an escape
self-loop on a stop state is authored. **A terminal context id is a reference,
not a state name: normalization MUST dereference each to the context's explicit
predicate set over owned tags**, in the same atom shape rules use, and the
normalized value carries the predicate sets — never the bare ids. This is what
satisfies RDR 0006 A6's shape requirement ("neither is a bare identifier a row
references by name"), which leaves the spelling to this RDR but fixes that lint
receives tag predicates: its invariant 1 checks terminals as tag keys and
values, and invariant 2 evaluates a terminal as a predicate over a node.
Naming a context is the authoring convenience; the predicate set is the
contract. Whether a model *omits* `[initial]` or
`terminal`, and whether a terminal context matches a non-owned tag, are RDR
0006's blocking findings, not load failures: the loader accepts the absence and
lint refuses to certify the graph without a root or stop set.
```

```normative
`[model]` MUST contain `id` and `version`. Version `1` is the only version this
RDR accepts; any other version MUST be refused before normalization.

**The version check MUST run before strict field validation, not merely before
normalization.** Loading MUST therefore proceed in two passes: read `[model]`
permissively enough to obtain `version`, refuse on any value but `1`, and only
then decode the document strictly. Ordering these the other way makes a future
v2 file fail as `unknown schema field` on whichever v2-only key the decoder
reaches first — a diagnostic that names an arbitrary field and never mentions the
version, which is precisely the confusion the version gate exists to prevent.

**Syntax precedes the gate.** The two-pass shape above requires parsing the
document to reach `[model]`, so a document that is both syntactically malformed
and version-2 necessarily refuses `malformed TOML`. The version gate's
precedence is therefore over **strict field validation and every later
category**, not over the parse that makes the gate reachable at all: the
ordering is `malformed TOML` → version gate → strict decoding → the remaining
categories. No lexer-level version scrape is required or wanted. A14's
precedence oracle is scoped accordingly — it discriminates the gate against
`unknown schema field`, which is the confusion it exists to prevent, and says
nothing about a document that will not parse.

**Load is fail-fast: the first category a document trips is the refusal.**
Beyond the two fixed precedences above, the order in which independent defects
are checked is deliberately **unspecified** — an implementation MAY check in any
order, and a document tripping two categories MAY be refused with either. Only
the parse-then-gate precedence above is normative. This is stated so the freedom is
visible rather than accidental: the Testing Strategy's one-defect-per-fixture
rule makes order unobservable by construction, so a later implementation MUST
NOT be read as bound to whichever order the first one happened to use, and an
accumulating loader that returns a list is a different contract than this one.

**An accumulating diagnostic mode is reserved, not foreclosed.** Fail-fast is
chosen for the *load contract* because it keeps each category independently
assertable and keeps the loader's return type a single refusal. It is a poor fit
for the hand-authoring loop this RDR exists to serve — one edit-reload cycle per
defect, with no guarantee about which defect surfaces first. That cost is
accepted here and paid elsewhere: a batch "report every problem in this file"
mode is explicitly permitted as a **lint-side** surface (RDR 0006) or a distinct
RDR 0005 verb built over repeated loads, and adding one later is not a change to
this contract. What is forbidden is the single `load` entry point returning a
list instead of a refusal — that, and only that, is the different contract.
Neither reserved surface is scheduled today, so the authoring cost is real and
unmitigated on arrival: the first author of a large model pays it in full. That
is survivable because the fixtures are small and the categories are precise, but
it is the strongest candidate for the first follow-up this format earns, and it
should be read as a known debt rather than a solved problem.

That argument is scoped to defects that **reach a check**. A defect consumed
earlier in the pipeline — merged away, de-duplicated, or otherwise absorbed
before any category is decided — trips *zero* categories rather than one, so a
one-defect fixture for it asserts nothing and the order-unobservability argument
does not cover it. Such absorption is a normalization defect in its own right
(the merge clause's set-over-full-identity obligation, the literals clause's
ban on joined identities, and the write-value sequence rule exist to prevent
exactly this), not a case the fail-fast freedom licenses.

Because this class trips no category, it is **asserted positively, not by
refusal**: the oracle is a count over the normalized value — the atom count for
the merge case (Testing Strategy 2's distinct-literal control) and the write
value's member sequence for the write case — each asserted to survive
normalization intact. A category-based test cannot cover an absorbed defect by
construction, so any future clause added to prevent one MUST arrive with such a
positive assertion or it is unenforced.

**Strict decoding is an obligation on this format, not a property of a library.**
The decoder MUST reject unmapped keys so an unknown schema field is a stable
refusal rather than a silent no-op. This is stated normatively because it is not
the default behavior of TOML decoders generally, and a parser swap that silently
lost it would retire the `unknown schema field` category without any contract
appearing to change.
```

```normative
Rule ids MUST be unique within a model, compared by exact byte equality; a
duplicate rule id is a load failure. Row identity, and therefore the dump's
total row ordering, rests on this.

Each transition rule MUST contain a stable rule id, zero or more shared-context
references, and a local match block. An ordinary transition rule MUST contain a
write block, MAY contain a rule-level explicit clear list, and MAY contain a
rule-level `gate` list. A write block MAY assign more than one tag. An escape
rule MUST contain an `escape` list and MUST NOT contain a write block, clear
list, or `gate` list, even an empty one (RDR 0009 binds the write-free
obligation at the kernel boundary).

**A write replaces.** A write assigns a tag's whole value and supplants
whatever was held; for a `set` kind the array literal is the whole new set.
There is no accumulate form. RDR 0004's read-back therefore compares the held
value for equality (JDR 0001 §D7(iv)).

**A set-valued write value is a member sequence, on the same terms as an atom
literal.** The literals clause's obligation is not scoped to predicates: a
`set`-kind write value MUST normalize to an ordered, member-sorted sequence and
MUST NOT be joined into a single delimiter-separated string, in the stored value
or in any comparison derived from it. The reason is identical and the stakes are
higher: joining makes `["a,b", "c"]` and `["a", "b,c"]` one value, and because
read-back compares the held value **for equality**, the collision decides
whether RDR 0004 reports a write as applied. A joined write value is the banned
identity of the literals clause reached through the write field, and the
Round-Trip invariant compares writes, so it must compare them as sequences.
```

```normative
An `escape` list MUST contain only resolver failure classes that RDR 0001
allows the table to model: `no_match` and `ambiguous_match`. Normalization MUST
render an escape rule as a candidate row carrying its normal predicate set,
outcome, source rule id, source locator, and modeled failure class list.

**Escape rescue is per-outcome, and an escape rule binds its outcome by the same
rule as any other.** An escape row rescues only resolves carrying the outcome it
binds: `internal/resolve/resolve.go::escapeOrRefuse` filters escape candidates
with `row.Outcome != in.Recognized` before gating them, exactly as `::Resolve`
filters ordinary rows. There is no table-wide catch-all: covering an alphabet of
N outcomes requires N escape rows. An escape rule MAY bind its outcome with an
`in` atom and expand like any other rule, which is how one authored rule covers
several outcomes — each expansion is a separate row rescuing its own outcome, and
the expansion suffix distinguishes them. This is a consequence of the shipped
kernel, not a choice this RDR makes; it is stated here because the escape clause
is where an author looks, and an author who expects one unbound escape row to
backstop the whole table gets silent unrescued refusals on every other outcome.

Row kind (`transition` / `escape`) is a **derived view property, never a row
field**: escape identity is discriminated solely by a non-empty escape class
list, which RDR 0009 fixes as the single discriminator and which forbids a
row-kind field at the kernel boundary. The dump renders the kind as a column
computed from that list; normalization introduces no field to carry it. The
prohibition binds the **normalized value and the kernel row**; a render-time
view struct MAY materialize the computed column, since the dump contract
requires the column. What it MUST NOT do is let that view feed back into the
normalized value or the kernel row — the discriminator stays the escape class
list, and no code may branch on a stored kind.
```

```normative
Shared contexts MAY inherit from other contexts, but inheritance MUST normalize
to an explicit predicate set before lint or resolution.

The combined predicate set is a **set over the full atom identity** — the tuple
`(key, block, operator token, literal)`, the same tuple the dump sorts on.
Merging MUST NOT key on any proper prefix of it: two atoms agreeing on
`(key, block, operator)` but differing in literal are **distinct atoms** and both
survive the merge. Keying the merge on `(block, key, operator)` alone silently
drops one constraint and makes the survivor depend on map iteration order over
the `use` list, which the dump's sort then certifies as deterministic — the loss
happens before the sort can see it, so ordering determinism cannot detect it.

Merging is idempotent on identical atoms: the same atom contributed by a rule and
by one or more inherited contexts collapses to one. Inheritance therefore never
overrides — it only accumulates. A narrowing "override" of an inherited
constraint is not expressible, and is not silently approximated: authoring two
atoms on one key that no view can satisfy together yields a dead rule, which is
RDR 0006's unreachable-rule finding, not a load failure here.

**The accepted cost is context forking.** With no override, a rule needing a
different value on an inherited key cannot narrow the parent — it must not
inherit it, so the chain forks and the shared atoms are restated in each fork.
At the fixture's scale (three contexts, five rules) this is invisible; it grows
with the number of dimensions that vary independently, and the Decision
Rationale names seven. This is the repetition Alternative 2 was rejected for,
relocated from rows to contexts, and it is worse there in one respect: the dump
shows expanded rows, so row-level repetition is visible to a reviewer, while
context-level repetition is only visible in the source. Accepted because the
alternative — an override rule — reintroduces order dependence into the merge,
which is the defect the full-atom-identity key above exists to prevent; a
narrowing override would make the survivor depend on inheritance order exactly
as prefix-keying did. If forking becomes the dominant authoring pattern, the
answer is a successor RDR on context composition, not an override bolted onto
this merge.
```

```normative
Guard predicates MUST be represented as positive `all` predicates and negative
`unless` predicates. Normalization MUST combine the match atoms and both guard
blocks into one candidate-row predicate set before ambiguity checks, and each
atom in that set MUST retain the key, operator token, literal, and the block
it was authored in — a three-valued domain, `match`, `all`, or `unless`. The
two-valued atom shape (`block ∈ {all, unless}`) is what JDR 0001 §D1 fixes and
what RDR 0007 and RDR 0003 spell; the widening to a third `match` member is
**this RDR's**, on JDR 0001 §D6's authority, because §D6 routes the kernel
handoff on the authored block and so needs match atoms distinguishable from
guard atoms. An implementer building the atom type from RDR 0007 alone gets two
members and MUST widen to three here. The block is load-bearing,
not merely carried: it is the key the handoff routes on (JDR 0001 §D6, below),
RDR 0003's atom identity tuple and `unless` semantics read it downstream, and
RDR 0006 serializes it. Normalization MUST NOT fold `unless` atoms into `all`
or `match` atoms into either guard block.

Because the block is part of the atom identity, one atom authored in **both**
`all` and `unless` is two distinct atoms and both survive normalization; this is
not a load failure. The rule is self-contradictory — no view satisfies a
predicate and its negation — so it is a dead rule, reported by RDR 0006's
unreachable-rule check. It MUST NOT be silently pruned at load, because a load
that dropped one block would turn an authoring mistake into a rule that
matches, which is the more dangerous failure.
```

```normative
For an existence atom the normalizer MUST emit the kernel's exported constants
verbatim — operator token `OpExists` and literal `LiteralTrue` or
`LiteralFalse` (RDR 0007, existence-operator clause; JDR 0001 §D4) — and MUST
reject at load, as a malformed predicate atom, any existence literal that is
not one of those two forms. The kernel treats a foreign token as a value atom
and a foreign literal as unevaluable; this load rejection is upstream of that
fail-closed backstop, not a substitute for it.
```

```normative
Tag-key identity is exact byte equality on the post-parse key string at every
stage — declaration lookup, context and rule predicate references, write and
clear targets, and the keys a normalized row carries. The normalizer performs
no case folding, trimming, or namespace rewriting (matching RDR 0008's
reserved-key comparison, which the same rule governs). The canonical spelling
of a tag is its `[tags.<tag>]` declaration key; every reference resolves to a
declaration by exact match or fails `unknown tag`, so a normalized row carries
only declared spellings and an atom's key agrees byte-for-byte with the key the
kernel assembles (RDR 0007 A22: the kernel canonicalizes nothing). The shipped
kernel already makes this decision and is the confirming sibling:
`internal/resolve/resolve.go::TagSet.matches`, `::Lookup`, and `::has` index by
raw map lookup with no folding, and the repo's only `strings.ToLower` is on the
`--as` flag value in `internal/cli/respond/respond.go::ModeOf`.
```

```normative
**Literals carry the same byte-exact identity as keys, and a set literal is a
sequence, not a joined string.** A set-valued literal MUST normalize to an
ordered sequence of its members, each compared byte-exactly; it MUST NOT be
rendered into a single delimiter-joined string as its normalized value. Joining
on any delimiter makes membership ambiguous whenever a member contains that
delimiter — space-joining collapses `["needs work"]` and `["needs", "work"]` into
one literal, so two different predicates become one atom and the dump's atom sort
certifies the collision as canonical. Members sort byte-lexicographically so two
authored orderings of one set are one literal; the sort is over members, after
which the sequence — not a joined rendering — is the value that atom identity and
the round-trip invariant compare.

**The obligation binds every comparison of an atom, not only its stored value.**
Wherever atoms are compared, keyed, deduplicated, or sorted — in particular the
**merge key** the contexts clause makes a set over `(key, block, operator token,
literal)` — the literal field MUST be compared as the member sequence,
element by element. An implementation MUST NOT derive that key, or any other
identity, by joining members into a string.

**Display is not identity.** RDR 0006's `clierr.Finding` deliberately carries
`Key`, `Operator`, `Literal`, and `Block` as flat `string` fields so `clierr`
takes no dependency on this RDR's vocabulary; that is a **transport and
serialization** boundary, and rendering a set literal into it is permitted. What
this clause forbids is deriving *identity* — merge keys, dedup, sort order,
round-trip comparison — from such a rendering. A finding string that reads
`needs work` is acceptable diagnostic text; a merge that treats it as one atom
is the defect. Producers rendering into that field MUST spell a set visibly as
a set **with each member delimited unambiguously** — quoting each member is the
cheap way — because an unquoted separator is forgeable by a member containing
it. Stage 6 found this concretely: a renderer separating members with an
unquoted `|` spelled `["x|y", "z"]` and `["x", "y|z"]` identically, so a test
asserting on the rendered value would have passed the exact collision this
clause forbids, even though the *identity* was already correct. A rendering
that cannot distinguish two distinct values is not merely imprecise for the
reader; it silently re-admits the defect through any oracle that reads it.

Keeping `Literal` a sequence in the
normalized value while keying the merge on a joined rendering satisfies the
letter of the paragraph above and still collapses two distinct atoms: the loss
happens in the merge, before the sort the paragraph above relies on can observe
it — the same blind spot the contexts clause names for a prefix key, reached by a
different route. Any delimiter is wrong here, not merely a delimiter a member
might contain, because no delimiter is banned from a tag value: the
expansion-suffix separator `#` is refused only in rule ids, alphabet members, and
match-block `in` members, so every other string a member can hold is authorable.

**Both downstream consumers were checked at Stage 6 and neither depends on a
joined rendering — they reinforce this clause rather than merely tolerating it.**
RDR 0003 independently requires the same thing from the other end
(`0003-guard-predicate-exhaustiveness.md:1381-1388`: a set literal "has one
canonical spelling: an unordered set of typed elements, duplicate-free, and
compared as a set … Implementations MUST canonicalize before the literal enters
the identity tuple"), and it stores no identity on the atom
(`0003:851-855`). RDR 0006 defers its own fingerprint to this RDR's atom sort
joined to 0003's canonical set form (`0006-graph-lint-authority-and-guarantees.md:983-990`,
"a canonical sortable serialization, never a hash"), and its flat-string
`Finding` fields are scoped by 0006 itself to transport (`0006:794-801`), which
is the same display/identity split drawn above.
```

```normative
**Rule ids, outcome literals, and the members of any match-block `in` literal
MUST NOT contain the expansion-suffix separator `#`.** The identity tuple
`(model id, rule id, expansion suffix)` is total only if its fields cannot
bleed into one another. The rendered row identity joins the rule id and each
suffix element with `#`, so a `#` inside any of them makes the rendered
identity ambiguous: rule `a#x` and rule `a` expanding on member `x` render the
same string. Rejecting `#` at load in every string that can reach the rendered
identity closes this at the source, and is what makes the rendered identity
recoverable — splitting on `#` is unambiguous. Violations are a malformed rule
id, a malformed recognized outcome alphabet, and a malformed predicate atom
respectively.

**The `<clear>` sentinel is reserved (JDR 0001 §D5).** The string `<clear>`
MUST be refused at load wherever a tag value is authored — a write-block value,
an `[initial]` assignment, or a predicate literal — as a `reserved tag value`.
Rendering an authored clear as a `<clear>` write is therefore unambiguous, and
RDR 0006 reads it as removal. The accessor-side half is **owed, not settled**:
JDR 0001 §D5 (§JD-15) makes read-back assert *absence* for a `<clear>` write,
but RDR 0004 does not yet carry that clause — its read-back prose asserts
presence with "lossy exemptions: none", and its own Status line records the
gap. This RDR therefore states the sentinel's authoring and rendering contract
only; a clearing rule is not end-to-end until RDR 0004 lands §D5. Do not read
this clause as a claim that it already has. A predicate literal of `<clear>` is refused
rather than admitted as a dead atom because no view can ever hold it and a
silent dead rule is the more dangerous failure.
```

```normative
A tag declaration with provenance `recognized` MUST be named `recognized`, and
no owned or observed declaration may take that name; violations fail in RDR
0008's `reserved_tag_key` category, which participates in this RDR's data-level
category set. The `recognized` tag is total over matching: the kernel refuses
`unmodeled_outcome` before any row is consulted unless the resolve carries an
outcome in the declared alphabet, so every view that reaches a row binds
`recognized` (JDR 0001 §JD-10). Consequently the `outcomes` alphabet MUST be
non-empty, duplicate-free, and MUST NOT contain the empty string. A rule
requiring `recognized` to be absent has **no authorable form**: the only
spelling — a `recognized` atom in a guard block — is refused at load as a
malformed outcome binding (below), so it never reaches lint. The dead rules
that *do* reach RDR 0006's unreachable-rule check are the authorable ones: one
atom in both `all` and `unless`, or two non-identical `eq` atoms on one key
accumulated through inheritance (the contexts clause).
```

```normative
Every rule — ordinary or escape — MUST bind exactly one outcome. **Outcome
binding reads the match blocks only** — the rule's local `match` block plus the
`match` blocks of its inherited contexts. That set MUST contain exactly one atom
on `recognized`, using `eq` or `in`, whose literal(s) are members of the
`outcomes` alphabet. Normalization lifts that atom out of the predicate set into
the row's outcome field; an `in` atom expands into one candidate row per member
under the general match-block expansion rule below, the chosen outcome literal
being that atom's contribution to the expansion suffix. A rule binding zero
outcomes, more than one `recognized` atom, or a literal outside the alphabet is
a load failure.

"Combined predicate set" is used in two extents in this document and they are
not interchangeable: **outcome binding** scans the match blocks (above), while
**ambiguity checking and row carriage** scan match plus `guard.all` plus
`guard.unless`. Lifting must use the narrower extent. A `recognized` atom
authored under `guard.all` or `guard.unless` MUST be refused at load as a
malformed outcome binding — never lifted. Lifting it from `unless` would inflict
the exact inversion the author guarded against, turning "this rule does not
apply to outcome X" into "this rule binds outcome X," with no diagnostic and a
normalized row that reads as intentional.

**Match-block expansion is general (JDR 0001 §D6).** Every `in` atom in a
rule's match blocks — local `match` and inherited context `match`, on
`recognized` or on any other declared tag — expands into one candidate row per
member, and the rule's rows are the Cartesian product of its expanding atoms.
In each expanded row the `in` atom becomes an `eq` atom on the chosen member,
still carrying block `match`, so every match-block atom the kernel receives is
an equality the `Match` pattern can test. The **expansion suffix** is the
sequence of chosen members, one per atom with **more than one member**, taken
in the atoms' sort order `(key, block, operator token, literal)`; it is empty
when no such atom exists. Because RDR 0003 rejects a repeated `in` element at
parse, no two rows of one rule share a suffix, which is what keeps the identity
tuple total.
A product row whose expanded atoms cannot be satisfied together (two `eq` on
one key, different literals) is a dead row for RDR 0006, not a load failure.

**Expansion suffixes attach only where a rule expands.** A single-member `in`
expands to one row whose suffix is empty, so `eq = "x"` and `in = ["x"]` are
one spelling of one edge and cannot mint two identities. A suffix is
non-empty exactly when the rule produced more than one row.
```

```normative
`Row.RequiresOwned` has no authored form. The normalizer MUST derive it for
every row as the sorted, duplicate-free set of tag keys named by the rule's
write block and clear list; by the write-to-non-owned-tag rule every such key
is an owned tag, and by RDR 0008 none is `recognized`. An escape rule carries
neither a write block nor a clear list, so a normalized escape row MUST carry
an empty set (RDR 0007 A21).

This is a **producer obligation, not a kernel guarantee**: the kernel checks
`RequiresOwned` on escape rows identically to ordinary ones — `escapeOrRefuse`
passes escape candidates through the same gate, and
`adversarial_test.go::TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement`
freezes that a populated set on an escape row yields
`owned_state_unavailable`. Nothing downstream rejects a malformed escape row
for this; the normalizer is the only enforcement point, which is why the
obligation is stated as a MUST on it rather than derived.

What the field means is RDR 0007's (post-guard write dependencies); this RDR is
its producer (JDR 0001 §JD-3) and does not add guard-read keys to it.
```

```normative
The kernel row carries the next state and the accessor-facing writes as two
distinct fields, and both have this RDR as their producer. Normalization MUST
populate the next-state tags with the tag values the rule's write block and
clear list produce, and the writes with the owned-tag writes the accessor
layer applies — the same rendered set, including `<clear>` entries. RDR 0009
A4 settles that its escape-conformance predicate stays `Writes`-only and does
**not** widen to `NextTags`; it therefore fixes no rule about these two fields'
distinctness, and this RDR does not rest on it for one. An escape row carries
neither a write block nor a clear list, so it normalizes to a row with both
empty — which is this RDR's own authoring-surface consequence, stricter than
A4's kernel-side predicate requires.

The two fields hold **equal sets under this RDR's authoring surface**, because
every write a rule can author targets an owned tag and therefore reaches the
accessor layer. They stay distinct fields rather than collapsing to one because
the kernel declares them as two fields (JDR 0001 §D2) and only the writes cross
the accessor boundary; a later RDR that admits a next-state tag which is *not* an owned write
would separate them by value. Normalization MUST populate both explicitly from
the rule and MUST NOT populate one by aliasing the other — an alias would make
that future divergence a silent behavior change rather than a compile-time one.
```

```normative
**The unified predicate set splits across the kernel's two predicate fields by
authored block (JDR 0001 §D6), and this RDR owns the split — at the kernel
handoff, not in the normalized value.** The normalized candidate row carries
the atoms as **one unified set with each atom's authored block retained**; that
single field is what the dump contract lists, what the Round-Trip invariant
compares, and what RDR 0003 reads downstream. The split is applied when a
`resolve.Row` is constructed for the kernel. It cannot be applied earlier
without loss: `Guard string` carries no per-atom block, so a normalized value
already split into `Match`/`Guard` could not satisfy the block-retention
obligation above.

The kernel row exposes `Match []Tag` (an equality pattern the kernel tests
directly) and a guard predicate (handed to the guard seam), and the handoff
MUST route each atom to exactly one of them by its block: every atom carrying
block `match` populates `Match`; every atom carrying block `all` or `unless`
belongs to the guard predicate, **regardless of operator** — an `eq` atom
under `guard.all` is a guard atom. The routing is exhaustive and disjoint by
construction because the block domain is exactly those three values.

The block is the author's declared intent about masking — "select on" versus
"gate on". A guard `eq` over an optional key that is absent reaches the seam
and refuses `guard_unevaluable` instead of silently dropping the row as a
non-candidate, which is the masking the kernel forbids; that is the wanted
behavior, not a deferral to avoid.

**Authoring rule for owned tags: gate on them, do not select on them.** An owned
tag placed in a match block is the format's most consequential authoring
mistake, and this RDR owns the surface where it is made. When the accessor
serving that key fails, a match atom makes every candidate drop out as a
non-candidate; `escapeOrRefuse` then runs and a `no-match` escape row returns a
**plan**. Missing owned state is laundered into a successful transition — the
outcome JDR 0001 §D2 exists to forbid. The same atom in `guard.all` reaches the
seam and refuses `guard_unevaluable`, which is loud and correct. RDR 0007 states
that "predicate PLACEMENT decides whether an absent key refuses non-escapably or
escapes" and routes the authoring guidance to RDR 0003; the guidance is repeated
here because this is the document an author has open while typing the block, and
a rule stated only where the author is not looking is not a control. Observed
tags carry no such hazard and may be selected on freely.

For the same reason a **match block admits
only `eq` and `in`** (`in` by expansion into per-member `eq` rows): a
comparison, existence, `contains`, or any other operator under `[rule.match]`
or `[context.<id>.match]` is refused at load as a malformed predicate atom,
because `Match` can only test equality and an atom the kernel cannot evaluate
there would silently fail to match rather than reach the guard seam. Which
surface the guard predicate itself presents — the atom shape and its encoding
— is RDR 0007's, and the reshape named in Prerequisites replaces the `Guard
string` field; this clause fixes only which atoms are the guard's, which is
stable across that reshape.

**The admitted operator set, guard blocks included, is RDR 0003's closed set
`{eq, in, lt, lte, gt, gte, exists, contains}`.** This RDR does not mint
operators and does not widen that set; it cites it by member list so
`unknown operator` is decidable here. A token under `[rule.guard.all.<tag>]` or
`[rule.guard.unless.<tag>]` outside the set is a `malformed predicate atom`
(`unknown operator`) refused at load — the guard-side mirror of the match-block
restriction above, which without this was absent. Operator *semantics* and the
literal/kind rules stay RDR 0003's; only the membership test is discharged
here, and it is the one this RDR's category floor already asserts on. Should
RDR 0003 re-lock with a different set, this list is the amendment site.
Re-checked at Stage 6: RDR 0003 is `Draft [revised from Final 2026-08-24]`, so
this pins a peer that is not currently Final — but the pinned content is
unchanged (`0003-guard-predicate-exhaustiveness.md:821-830` enumerates exactly
these eight tokens; `:154` restates the closed set), and 0003's revision reason
(rejection-rule routing, §JD-18) does not touch the operator vocabulary.
```

```normative
**Atom-level validation is block-agnostic.** Every rule this RDR states about an
authored atom — operator membership (above), the `<clear>` reserved value,
domain and kind conformance of a literal, `#` reservation in members, and tag-key
declaration — MUST be enforced identically in all three atom blocks: `match`,
`guard.all`, and `guard.unless`. A validation implemented for match blocks only
is a defect against this clause, not a permitted reading of it.

**Conformance is per operator, not per literal.** Operators do not all take a
domain member as their right-hand side, so "domain/kind conformance" binds by
operator: `eq`, `in`, and `contains` take members and are domain-checked;
`exists` takes a bool literal (`<true>`/`<false>`), never a member;
`lt`/`lte`/`gt`/`gte` take an ordered **bound**, which is kind-checked but not
domain-checked, because a bound need not itself be an authorable value — the
RDR fixture's own `iter.lt = 3` against `min = 0, max = 9` is the witness. The
`<clear>` ban and the tag-key declaration rule bind every atom regardless of
operator. Two rules are legitimately match-only and are **not** block-agnostic:
the `eq`/`in` operator restriction (§D6 routing) and the `#` reservation in `in`
members, since only match blocks expand.

This clause was stated because the §D7 spike enforced exactly the narrower
reading, and that was reached by *running* the spike rather than reading its
transcript (`gen-cases.py` derived every negative case from the RDR fixture by
one mutation and never mutated a guard block, so no transcript line witnessed
the gap — which is why two lenses passed over it). The dangerous case was
`<clear>`: an unsatisfiable `unless` atom is a silent dead rule, the failure the
reserved-value clause exists to prevent, and an unsatisfiable *guard* atom
triggers RDR 0007's veto on selecting a decided sibling while an unevaluable
candidate exists — so the flow deadlocks rather than refusing loudly.

**Closed at Stage 6 (A15 `Verified`).** The spike's guard walk now routes
through the same block-agnostic checker as match blocks, and `gen-cases.py`
mutates guard blocks: six controls, one per atom rule per block, each refusing
the category it names (`evidence/spikes/iter-2/negative-cases.txt`).
Implementation promotes them; scenario 3 continues to owe one negative control
per rule **per block** for any rule added later.
```

```normative
The tool MUST normalize the sparse source into deterministic candidate rows for
lint, resolver lookup, diagnostics, and table dumps. Each candidate row MUST
retain its source rule id and source locator.
```

```normative
The expanded table dump MUST be derived from the normalized candidate-row value
and MUST carry every field of it: row identity, source locator, outcome,
predicate atoms (with block), gate ids, next-state tags, writes including
`<clear>` entries, required-owned keys, and escape failure classes — plus the
derived row kind column. `[dump]` carries exactly one key, `order`: a list of
column identifiers. The identifiers are the **field names enumerated in the
sentence above**, snake_cased on the same rule as the category identifiers.
That closed list is the column vocabulary, so an implementer need not invent
one; it is fixed here **verbatim as the normative fixtures author it**, because
those fixtures are what implementation tests promote:

    identity, source, kind, outcome, atoms, next, writes, requires_owned,
    gate, escape

The vocabulary is the fixtures' spelling, not a re-derivation from the field
names in the sentence above: `identity` (not `row_identity`), `source` (not
`source_locator`), and `writes` (not `write`). A clause that snake_cases the
prose field names instead would refuse both normative fixtures at load under
the very rule below — the check that catches this is fixture conformance, run
as part of scenario 1, not review. `[dump]` settings MAY
reorder the rendered columns; they MUST NOT omit a field. An `order` naming an
unknown identifier, repeating one, or omitting one is a `malformed dump
declaration` refused at load — not a silently truncated dump. A dump missing any
field is not an expanded table dump and does not satisfy the Round-Trip
invariant.

`[dump]` is **presentation, and MUST NOT reach the normalized value**: editing it
changes column order in the rendered view and nothing else. Normalization MUST
ignore it entirely, so two models differing only in `[dump]` normalize to
identical candidate-row sets and the Round-Trip invariant — which compares values,
not dump text — is unaffected by it. It lives in the source file because the
column order a table wants reviewed travels with that table; the cost is that a
`[dump]` edit churns golden files asserted over rendered text. Golden tests
SHOULD therefore assert over the normalized value, which no `[dump]` edit can
move, and reserve rendered-text goldens for tests of rendering itself.

Dump ordering MUST be deterministic across source key order. Rows sort by row
identity, compared field by field as the tuple `(model id, rule id, expansion
suffix)`; `model id` and `rule id` compare byte-lexicographically on the
post-parse string, and the expansion suffix compares as a sequence —
element by element, byte-lexicographically, a shorter sequence that is a
prefix of a longer one sorting first, so the empty suffix sorts before any
non-empty one. Within each row, atoms sort by (key, block, operator token,
literal), and next-state tags, writes, required-owned keys, gate ids, and
escape classes sort by key.

The ordering is total **over one model's rows**: rule ids are unique within a
model, and two rows of one rule differ in at least one suffix element because
no `in` literal repeats a member — so no two distinct rows compare equal and no
positional tiebreak is needed (RDR 0009 A8 establishes that a tiebreak on table
position is forbidden; RDR 0001's identity rule makes two orderings of one row
set the same input). A dump spanning more than one model MUST carry a
model-unique `model id`, which is what makes the leading tuple field
discriminating; two models sharing an id is a `duplicate model id` load
failure. RDR 0006 lints exactly one model per invocation.

**One source document carries exactly one model** — `[model]` is a singular
table, not an array of tables — so `duplicate model id` is never decidable
within a single document. It is scoped to a caller that loads **several
documents into one dump or lint invocation** and MUST be checked there, over the
set of loaded models, before rows are merged for rendering. A loader handed one
document MUST NOT report it. The Testing Strategy's mutated fixture for this
category is therefore a **pair** of documents sharing a `model id`, not a single
malformed file.

The **source locator MUST NOT participate in row ordering**. It may carry
optional line/column detail, so ordering on it would make the dump reorder
when unrelated source text is edited — the churn this contract exists to
prevent. This is a deliberate divergence from RDR 0007's refusal-payload
tuple, which orders on `(RuleID, SourceLocator, …)` because a refusal payload
is not a stable artifact under review.

The dump MUST be emitted from a pre-sorted sequence of normalized rows, never
by iterating a map and never by delegating key order to an encoder. Go's map
iteration order is deliberately randomized, and encoder guarantees vary:
`encoding/json` v1 sorts map keys but `encoding/json/v2` does not without an
explicit option, and TOML encoders may sort only within key groups.
```

```normative
Source order and rendered-row order MUST NOT decide a successful transition.
This is this RDR's clause: the normalized row set is unordered as far as
selection is concerned, so reordering the authored rules or the dump cannot
change a disposition, and normalization MUST NOT emit any positional field a
consumer could tiebreak on.

The selection procedure itself — gate-then-count, and which refusals are
escapable — is **JDR 0001 §D2's and the kernel's, and MUST NOT be restated
here**; RDR 0007 is its normative home for the guard seam. This RDR's obligation
is only to produce rows that procedure can evaluate: exactly one outcome bound
per row, escape rows carrying their modeled failure classes, and no positional
dependence. Ordinary transition success requiring exactly one surviving
non-escape row, and a modeled escape requiring exactly one surviving escape row
for that class, are consequences of that procedure this RDR relies on rather
than defines.
```

```normative
The model MUST declare every tag it matches or writes, including each tag's
provenance: owned, observed, or recognized.
```

```normative
A tag declaration also carries its **type model**, spelled as the wire keys
`kind`, `domain`, `min`, `max`, `elements`, `single_valued`, and `required`
(JDR 0001 §D7(iii); `required` is the optionality marker and defaults to
optional). RDR 0003 is the normative home of that model: what those keys mean,
which kinds admit a finite domain, and how domain/kind disagreement is rejected
are stated there and MUST NOT be restated here. This RDR owns where a
declaration is authored — under `[tags.<tag>]`, beside `provenance`, with no
accessor reference — and the two load categories that carry RDR 0003's
rejection rules: `malformed tag declaration` (a declaration RDR 0003's rules
reject, decided before normalization completes) and `malformed predicate atom`
(a literal outside the tag's declared domain, decided after declarations load
and before rows are yielded). Normalization MUST carry every declared field
through to the normalized model without loss, so lint (RDR 0006) and the guard
proof (RDR 0003) read the same declaration the author wrote.
```

```normative
Clearing a tag MUST be represented by an explicit rule-level `clear` entry that
normalization renders as a `<clear>` write. Absence from both the write block and
the clear list MUST NOT imply deletion.
```

```normative
Load-time validation failures MUST retain stable data-level categories before
CLI mapping, including at minimum malformed TOML, unknown schema field, missing
recognized outcome alphabet, malformed recognized outcome alphabet (an empty
alphabet, a duplicate member, the empty string as a member, or a member
containing the expansion-suffix separator — each distinguishable), unknown tag,
unknown context, cyclic context inheritance, write to non-owned tag, unknown
accessor, malformed accessor declaration, malformed accessor binding,
malformed tag declaration, malformed initial declaration, reserved tag value,
unsupported version, malformed predicate atom (unknown operator; literal
ill-formed for its operator; literal outside the declared domain; operator not
admitted under a match block; `#` in a match-block `in` member — each
distinguishable), malformed escape declaration (a write block, clear list, or
`gate` list on an escape rule), **malformed model declaration** (an absent
`[model]`, an absent `id`, an absent `version`, or an absent `[tags.recognized]`
declaration — the shape failures that precede any version comparison; an absent
`version` MUST refuse here and MUST NOT be folded into `unsupported version`,
which would name a version the author never wrote), **malformed dump
declaration** (a `[dump]` whose column list names an unknown column, repeats
one, or omits one — "MUST NOT omit a field" is a refusal, not a silent
truncation), malformed outcome binding (including a
`recognized` atom authored in a guard block), **malformed rule shape** (an
ordinary rule carrying no write block — the mirror of `malformed escape
declaration`, which covers only the escape-side shapes), malformed rule id
(containing the expansion-suffix separator), duplicate rule id, duplicate model
id, and `reserved_tag_key` (RDR 0008).

These are the checks the load/lint arity split assigns to this RDR. All but one
are single-rule — decidable from one rule plus the model's declarations.
`duplicate model id` is the sole exception and is **not** single-rule: it is
decidable only over the set of models a caller loads into one invocation (the
clause above), so a loader handed one document MUST NOT report it. It is listed
here because it is a load-time refusal this RDR owns, not because the arity
split makes it single-rule.
Cross-row findings — overlap, gap, dead row, read-before-write — carry RDR
0006's lint categories, not these.

**"Load" names the whole source-to-candidate-rows pipeline, not one callable.**
Every category above MUST be refused by that pipeline before it yields candidate
rows, whatever internal decomposition an implementation chooses. This is stated
because several categories are only decidable during normalization — outcome
binding and the lifted literal's alphabet membership are decided while lifting,
and the write-free escape check names the normalizer as its only enforcement
point — so a reading of "load-time" as "before normalization begins" would make
them unreachable. An implementation MAY split parse, validate, and normalize
across several exported calls; it MUST NOT let a category escape because the
stage that decides it ran after the one the caller thinks of as "load".

The categories are **data-level and MUST NOT depend on the CLI envelope**: this
package MUST NOT import `internal/cli`, and CLI code names, `Group` assignments,
and exit mapping remain RDR 0005's design work (A5). `internal/cli/clierr` is a
leaf package for exactly this reason, so a later mapping needs no restructuring
here.

Category identifiers are **snake_case renderings of the prose names above**,
anchored on `reserved_tag_key` — the one member spelled as an identifier, whose
spelling RDR 0008 owns. The Testing Strategy asserts on the category, so these
identifiers are an API surface, not message text.
```

#### Load-Bearing Decisions

- **Identity** — a transition rule is identified by `(model id, rule id)`.
  Normalized candidate rows inherit that identity plus a deterministic expansion
  suffix — the sequence of chosen members, one per multi-member `in` atom in
  the rule's match blocks, in atom sort order. Rule ids are stable review
  anchors and must not be reused for a different edge. The identity tuple
  `(model id, rule id, expansion suffix)` is also the dump's row-sort key, and
  it is total; the source locator is diagnostic and deliberately excluded from
  it. Rejected: a suffix drawn from the `outcomes` alphabet alone, because
  JDR 0001 §D6 generalizes expansion to every match-block `in` and a
  non-outcome member is not an alphabet element; rejected: a positional
  suffix (`#1`, `#2`), because it would move when an unrelated member is added
  to the literal.
- **Wire / byte format** — TOML is the on-disk carrier. The exact field names
  are the JDR 0001 §D7 closed layout: root `outcomes`, root `terminal`,
  `[model]`, `[model.metadata]`, `[initial]`, `[tags.<tag>]` with `provenance`
  and the seven type-model keys, `[read.<id>]`/`[write.<id>]`/`[gate.<id>]`
  with `role`, `path`, `keys`, `timeout`, and `read_back`, `[context.<id>]`,
  `[[rule]]`, `[rule.write]`, rule-level `clear`, `gate`, and `escape`, and
  `[dump]`. `[model].version = 1` is the only accepted format version. The RDR
  and kata fixtures, re-authored to this layout, are the canonical examples
  implementation tests must promote. Rejected: a tag-side `accessor`
  reference, because `keys` on the capability table already binds the pair
  and a second copy can disagree; rejected: `[accessors.<id>]` with a `mode`
  field, because one-capability-per-entry is then a validation rule rather
  than a shape the strict decoder refuses.
- **Naming** — the canonical source artifact name is "transition model"; the
  canonical rendered view is "expanded transition table." Tag keys are exact
  byte strings with no folding, and the recognized-provenance tag is the
  reserved name `recognized`. Rejected names: "state machine config" because
  it suggests a runtime driver, and "workflow graph" because this RDR owns
  sparse transition data, not orchestration. Rejected: case-insensitive or
  trimmed key matching, because it would let two authored spellings reach the
  kernel as distinct keys and would make `Recognized` collide with the reserved
  key RDR 0008 compares byte-exactly.
- **Selection / predicate** — the only successful transition selection is
  exact-one surviving non-escape row after the kernel's gate (guard-FALSE
  pruned, owned-state and unevaluable refusals first). Count-first pruning is
  rejected because it launders missing state into an escapable `no_match`
  (JDR 0001 §D2). Lint rejects ambiguous ordinary overlaps and ambiguous escape
  overlaps before runtime.
- **Routing key** — an atom reaches `Match` or the guard by its authored
  block, never by its operator (JDR 0001 §D6). Match blocks admit `eq` and
  `in`; everything else there is a load refusal. Rejected: operator-keyed
  routing, because it makes the author's block advisory and turns an `eq`
  guard over an optional key into a silent non-candidate — the masking the
  kernel exists to forbid.
- **Outcome binding** — a rule's outcome is the single `recognized` atom in its
  **match blocks** (local `match` plus inherited context `match`), lifted into
  the row's outcome field. The extent is the narrow one the Normative Contracts
  fix; "combined predicate set" is this document's *wider* term (match plus
  `guard.all` plus `guard.unless`) and MUST NOT be read here. Rejected: an
  eventless row with no outcome, because the kernel has no such concept and
  would never match it; rejected: leaving the atom in the predicate set only,
  because the kernel filters on `Row.Outcome` before matching.

#### Round-Trip / Inverse Invariants

`parse ∘ normalize = expanded-table value identity` on valid transition model
fixtures: two authorings of one model — differing in TOML key order, rule
declaration order, or `eq`-versus-single-member-`in` spelling — must normalize
to the same candidate-row set, compared as values over the full field list the
dump contract carries (row identity, source locator, outcome, predicate atoms
with block, gate ids, next-state tags, writes, required-owned keys, escape
failure classes), with the source locator's optional line/column detail
excluded from the comparison.

The locator's **presence and its rule-identifying part are compared**, even
though its line/column detail is not: excluding the whole field would let a
normalizer that dropped locators entirely satisfy the invariant, and the locator
is what routes a lint or refusal diagnostic back to the authored rule. Only the
positional detail is exempt, because it moves when unrelated source text is
edited — the same reason the locator is excluded from row ordering.

The invariant is stated over the **normalized value, not the dump text**,
because this RDR fixes the dump's field list and row ordering but does not
define a dump grammar — no delimiter, escaping, quoting, or record separator.
`dump` therefore has no specified inverse, and a textual round-trip is not
claimed. One site makes the rendered form lossy today and would have to be
closed by any later RDR that defines a re-readable dump: set-valued and
multi-entry fields are rendered with unescaped separators. This makes the
rendered form of a **set-valued atom literal** non-recoverable — `["a,b", "c"]`
and `["a", "b,c"]` render identically — and the `#`-ban does not reach it: `#`
is refused in rule ids, alphabet members, and match-block `in` members, which is
what makes the rendered *row identity* recoverable, but an atom literal in a
guard block or on any `string`-kind tag may hold any byte. (The `<clear>`
sentinel is reserved and the suffix separator is banned from every string the
row identity joins, so those two renderings are unambiguous.) Rendering is a
review surface here, and the normalized value is the contract — which is why the
literals clause forbids a joined rendering from becoming an identity anywhere
upstream of the dump.

This is a real narrowing of the Problem Statement's promise, and it is stated
rather than left for a reviewer to discover. "Every legal edge lives in one
reviewable artifact, so 'is this transition legal' has one answer" holds at the
level the artifact answers: **which rows exist, what each binds, and in what
order** — the questions a table review is actually for, all of them recoverable
from the dump. What is not recoverable from the *rendered text alone* is the
exact member decomposition of a set-valued literal, a strictly narrower question,
and the reviewer who needs it reads the source rule the locator names. The
alternative readings were to define a grammar (which claims an inverse this RDR
declines) or to ban delimiters from literals (which constrains authored data to
serve a render). Both cost more than the residue is worth; the visible-set
requirement below keeps the ambiguity honest at the point of reading.

**Phase 2's dump is not required to carry an escaping grammar, and that is a
decision, not an omission.** Defining one would make the dump a re-readable
format — an inverse this RDR explicitly does not claim — and every consumer that
needs recoverability has the normalized value, which is unambiguous. What the
dump owes the reviewer instead is that the ambiguity be **visible rather than
silent**: a set-valued field MUST render its members delimited in a way that
shows the field is a set (for example bracketed), so a reviewer reading two rows
can see that a difference *may* be hiding, even though the rendered text alone
cannot settle it. A reviewer needing the exact members reads the source rule the
locator names. A re-readable dump grammar is a follow-up RDR, seeded from this
pass rather than assumed; nothing here depends on it landing.
Source rewrite is likewise out of scope for this RDR; a later rewrite-capability
RDR must define its own source-preservation invariant before mutating authored
TOML.

#### Conditional Mini-Checks

**`fidelity`** — the chain's invariant holds over the normalized value; the
rendered form is lossy at one named site and carries no inverse.

| Operation | Invariant | Lossy exemptions |
| --- | --- | --- |
| parse (TOML → source) | none claimed — comments, key order, and whitespace are dropped by design; `[model.metadata]` carried verbatim | source rewrite out of scope |
| normalize (source → rows) | **value identity**: key order, rule order, and `eq`/single-member-`in` spelling all normalize to one candidate-row set | contexts flattened; `recognized` atom lifted out; every multi-member match-block `in` expands 1→N (product across atoms); set literals member-sorted; `match`/`all`/`unless` merged with block retained |
| dump (rows → text) | fixed field list + total row order; every field present, columns reorderable from the closed `order` vocabulary | **no grammar by decision, not omission**: separators unescaped, so a set-valued atom literal's rendered form is non-recoverable — set-valued fields must render visibly as sets so the ambiguity is not silent (`<clear>` reserved and `#` banned from row-identity parts, so those two render unambiguously; the ban does not reach atom literals) |
| read-back | **not claimed** — no inverse is specified | no successor RDR exists yet; seeding one is a named follow-up, not a scheduled dependency (below) |

**`disposition`** — routing is by arity: single-rule → load (this RDR),
cross-row → lint (RDR 0006), evaluation-time → kernel refusal (RDR 0001).
`duplicate model id` is the one load-time check that is cross-**document**
rather than single-rule (Normative Contracts); it stays this RDR's.

| Input class | Bucket | Owner | Silent vs loud |
| --- | --- | --- | --- |
| malformed TOML, unknown schema field, unsupported/duplicate version or model id, malformed model declaration (absent `[model]`/`id`/`version`/`[tags.recognized]`), malformed dump declaration | load | this RDR | loud, stable category |
| unknown/cyclic context, unknown tag, unknown accessor, malformed accessor declaration/binding, write to non-owned tag, malformed initial declaration | load | this RDR (0004 supplies accessor semantics) | loud, stable category |
| malformed tag declaration, malformed predicate atom (incl. literal outside domain, non-`eq`/`in` under a match block), reserved tag value, escape declaration, outcome binding, alphabet; `reserved_tag_key` | load | this RDR (0003 supplies the declaration/domain rules; 0008 the reserved key) | loud, stable category |
| missing `[initial]` / `terminal`, terminal context on a non-owned tag | lint | RDR 0006 | loud, blocking |
| overlap, gap, dead row, read-before-write, ambiguous expansion | lint | RDR 0006 | loud, blocking |
| expansion counts per rule | dump (derived: rows per source rule id) | this RDR | reviewer-visible, non-blocking; no threshold here. Not an RDR 0006 finding — its advisory tier is closed at four members |
| `unmodeled_outcome`, `owned_state_unavailable`, `guard_unevaluable` | kernel refusal | RDR 0001 | loud; the last two are never escapable |
| zero/multiple survivors | kernel refusal, or modeled escape when exactly one escape row survives the gate | RDR 0001 | loud |

**`oracle`** — every MVV assertion carries a failing control; the rail and the
per-item controls are in the Minimum Viable Validation and Testing Strategy.
The three that previously passed against a wrong implementation — the category
oracle (message text, not category), the determinism oracles (a constant
passes), and "the comparator is total" (satisfied by the positional tiebreak it
excludes) — are replaced by category-level assertions, RDR-fixture permutation,
and rule-order permutation respectively. A fourth is added at cove iter-2: the
**merge distinct-literal control**, which passed against a merge keyed on a
joined rendering because every witness used a literal containing no delimiter.
Its replacement asserts the **atom count** and parameterizes over the plausible
delimiters, since a black-box test cannot read the delimiter a wrong
implementation chose (Testing Strategy 2). Two more are added at 3amigo iter-3:
the **write-value control** (same defect through the write field, asserted by
value inequality) and the **`[model.metadata]` deep-equality control**, which
replaces a top-level-key-name oracle a value-discarding loader passes. The
pattern across all six is the same — a control whose inputs cannot distinguish
the right implementation from the wrong one — so a new assertion is not accepted
here until a wrong implementation that fails it is named. Two of the six share a
further property worth naming: an **absorbed** defect trips no category, so its
oracle must be a positive assertion over the normalized value, never a refusal.

**`trace`** — desk trace of the MVV against the clauses in force. One row per
step; witness values from `evidence/spikes/iter-2/output.txt`.

| Step | Assertions in force | Witness / verdict |
| --- | --- | --- |
| load fixture | closed §D7 layout under strict decoding; alphabet non-empty, dup-free, no empty string; version = 1 before field decoding | `STRICT-OK`; `outcomes=round-clean,verdict-flapping,reconcile-block,finalized terminal=archived metadata_keys=owner,review_cadence` — holds |
| normalize | outcome lifted; per-atom block retained; `RequiresOwned` = writes ∪ clears | `continue-prelock … atoms=[…profile.eq=small@unless…] requires_owned=[iter,stage]` — holds |
| expand `in` | one row per member of every multi-member match-block `in`; suffix only where the rule expands | holds on both arms: `terminal-archive#finalized` / `#verdict-flapping` for the `recognized` atom, and `continue-prelock#large` / `#foundational` for the **inherited non-`recognized`** `profile.in`, so the expansion is general and the suffix is a sequence. `eq`/single-member-`in` equivalence is witnessed by `perm/rdr-eq-as-in.toml` digesting identically to the baseline |
| escape row | no write block or clear list; `RequiresOwned` empty; expands like any other rule | `draft-no-match-escape#round-clean` / `#reconcile-block … next=[] write=[] requires_owned=[] escape=[no_match]` — holds, and covers A11 |
| merge atoms | the merged set is a set over the full atom identity `(key, block, operator token, literal)`, compared as a member sequence | **holds** — closed at Stage 6. The spike had keyed the merge on a comma-joined rendering (`main.go::Atom.identity` → `literalString`), collapsing `in = ["a,b", "c"]` and `in = ["a", "b,c"]` on one key and block into one atom and normalizing a live expanding row where the author wrote a dead rule; scalar witnesses could not see it. `Atom.identity` now encodes a length-prefixed member sequence (`::memberKey`) and `merge-delim-atom.toml` witnesses both atoms surviving. A13 `Verified` |
| dump order | rows by identity tuple; locator excluded; no positional tiebreak | three runs byte-identical, and key-order **and rule-order** permutations digest identically — holds |
| resolve, sibling gate | unevaluable survivor refuses despite a decidable sibling and an escape row | inputs now present: `continue-prelock` and `continue-prelock-cluster` both bind `round-clean`, the latter guarding on an optional owned key (`cluster_ready.eq=true@all`), alongside a `round-clean` escape row. The `Resolve` run itself is scenario 4's and waits on RDR 0007's reshape |

| write values | a `set`-kind write value is a member sequence, member-sorted, never joined | **holds** — closed at Stage 6. The spike had joined write values (`main.go::renderValue`), colliding `["a,b","c"]` with `["a","b,c"]` — the equality RDR 0004's read-back compares on. `Write.Value` is now a `[]string` member sequence with a matching `Write.identity`, and `write-delim-a/b.toml` witness the two spellings normalizing to different values (`labels=[x,y|z]` vs `labels=[x|y,z]`) |

**No CONTRADICTION row survives.** The merge row and the write-value row were
the same joined-rendering defect reached through two fields; Stage 6 fixed the
spike on both and minted the delimiter-bearing controls that discriminate them,
so A13 and A15 are `Verified` rather than deferred. The earlier gaps this table
recorded — the missing sibling pair and the outcome-only expansion — were closed
by the §D7 re-authoring. Every row now holds against executed evidence, except
`resolve, sibling gate`, whose `Resolve` run is scenario 4's and waits on RDR
0007's kernel reshape (Prerequisites) — a sequencing block, not a contradiction.

#### Illustrative Code

Illustrative sparse source shape, not a locked schema:

```toml
outcomes = ["round-clean", "verdict-flapping"]
terminal = ["archived"]

[model]
id = "rdr"
version = 1

[initial]
status = "Draft"
stage = "propose"

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["Draft", "Final"]

[tags.stage]
provenance = "owned"
kind = "enum"
domain = ["propose", "prelock", "archive"]

[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large", "foundational"]

[tags.iter]
provenance = "owned"
kind = "int"
min = 0

[tags.finalized_at]
provenance = "observed"
kind = "string"

[tags.recognized]
provenance = "recognized"
kind = "enum"

[read.rdr-status]
role = "rdr"
path = "rdr.status"
keys = ["status", "stage", "profile", "iter"]
timeout = "2s"

[write.rdr-status]
role = "rdr"
path = "rdr.status"
keys = ["status", "stage", "iter"]
timeout = "2s"
read_back = true

[gate.rdr-lock]
role = "rdr"
path = "rdr.lock"
keys = ["status"]
timeout = "2s"

[context.draft.match.status]
eq = "Draft"

[context.prelock]
inherits = "draft"
[context.prelock.match.stage]
eq = "prelock"
[context.prelock.match.profile]
in = ["large", "foundational"]    # expands: one row per member

[context.archived.match.stage]    # the stop set `terminal` names
eq = "archive"

[[rule]]
id = "prelock-flapping-cap"
use = ["prelock"]
gate = ["rdr-lock"]
[rule.match.recognized]
eq = "verdict-flapping"
[rule.guard.all.iter]
lt = 3
[rule.guard.all.finalized_at]
exists = false
[rule.guard.unless.profile]
eq = "small"
[rule.write]
stage = "prelock"

[[rule]]
id = "draft-no-match-escape"
use = ["draft"]
escape = ["no_match"]
[rule.match.recognized]
eq = "round-clean"
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Stateless exact-one resolution | RDR 0001 | Implemented | This RDR must provide normalized candidate rows, outcome binding, and tag writes the kernel can evaluate. |
| Fixed predicate operators and tag type model | RDR 0003 | Pending | This RDR names predicate slots and carries declarations but does not own the operator grammar or declaration semantics. |
| Accessor execution and safe read-back | RDR 0004 | Pending | This RDR authors the capability tables, the `keys` binding, and rule-level gate references, and validates the bindings; RDR 0004 owns execution, `role`/`timeout`/`read_back` semantics, and what a `<clear>` write and its read-back mean. |
| CLI parse/lint output | RDR 0005 plus existing respond gateway | Pending | Failures must map to the CLI output contract. |
| Graph lint over normalized rows | RDR 0006 | Pending | This RDR must expose enough structure for determinism and reachability checks, including the declared `[initial]` root and `terminal` stop set (RDR 0006 A6). |
| Guard seam, atom shape, existence constants, gate ordering | RDR 0007 | Final, kernel reshape unimplemented | Normalizer emits per-atom block and the kernel's `OpExists`/`LiteralTrue`/`LiteralFalse`; derives `RequiresOwned`. The constants and the atom-shaped row do not exist in `internal/resolve` yet, so Phases 2–3 sequence behind that reshape (Prerequisites). |
| Reserved recognized tag key | RDR 0008 | Final | Declaration named `recognized`; `reserved_tag_key` joins this RDR's category set. |
| Escape-row shape conformance | RDR 0009 | Final | Escape rules carry no write block or clear list; normalized escape rows render write-free. |
| Expanded table dump | This RDR | Introduced | Reviewers can inspect the full table without maintaining it by hand. |
| Graph/render export | RDR 0006 | Pending | This RDR exposes normalized rows and source ids; graph-specific rendering remains with graph lint. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Structured CLI failures | `internal/cli/clierr::CLIError` | Envelope exists; no parse-validation code has ever travelled it | Extend | Add stable parse/lint refusal codes later — these are the first of their kind, so names, `Group`, and exit mapping are new design work. |
| Text/json output gateway | `internal/cli/respond::Fail` | Gateway is CLI-only, not kernel behavior | Reuse | Parser/lint commands must report through existing gateway. |
| Kernel row and resolver | `internal/resolve::Row`, `::Resolve` | Consumes rows; has no loader or normalizer | Reuse | The normalizer targets this row shape; no new kernel surface. |
| Normalizer/table package | `internal/` search | Not implemented | Introduce | New internal package can own sparse source structs and normalization. |
| Config discovery | `internal/cli/config::Load` | Discovery and file read exist; no parsing or validation is wired | Reuse later | Table path binding belongs with CLI integration, not this RDR's data format. The load entry therefore takes **already-read bytes plus a source id** for the locator, not a filesystem path — this package performs no file I/O and no path resolution. A path-taking convenience wrapper MAY live at the CLI seam. |

### Decision Rationale

Joint-check: fired → 0007 (home: JDR 0001 §D4 / §JD-12) — the kernel enforces the guard domain; the normalizer obligations §D4 places on this RDR (existence constants, exact key identity, per-atom block) are Normative Contracts above, citing RDR 0007 rather than restating it.

Sparse TOML is the best fit because the source is meant to be reviewed and
edited by humans, while the expanded table is a mechanical view for lint,
debugging, and documentation. A fully expanded table would make the RDR model's
dimensions multiply: status, profile, stage, prelock iteration, cluster
eligibility, rewind scope, and guards would force authors to copy the same
predicate fragments across rows. That is exactly the DX failure this RDR must
avoid.

The resource corpus points to the same split. `BUILD-SEEDS.md` says the table is
the flow's design and is hand-authored, but also says the resolver is table +
thin CLI + lint, not a runtime. `ANALYSIS-kernel-vocabulary.md` supplies the
vocabulary: recognized outcome is a Deferred Choice, guards are an Exclusive
Choice layer, and ambiguous enabled edges should be lint-rejected rather than
resolved by document order. The Ragel POC shows why flat compiled topology is
insufficient for the RDR model: the hard parts are the coupled status register,
cap-3 counter, guards, and readable current state. Therefore the source should
be sparse and semantic, while normalization can render explicit rows for tools.

The direct `arc` corpus checks sharpen the sparse shape. Sismic demonstrates a
source model with nested states, transition guards, and contract classes that can
be exported to another view; that supports inherited contexts plus separate
precondition/invariant/postcondition diagnostics. `transitions` demonstrates
positive `conditions` and negative `unless` guard factoring; that supports
`all`/`unless` blocks instead of forcing every guard into one expression string.
Stateless demonstrates a symbolic graph object built from machine metadata; that
supports an expanded table/graph dump as a rendered view rather than the
authoring source. The statechart literature search supports hierarchy and
extended state as the known way to contain dimensional state explosion, while
the project still rejects adopting a statechart runtime.

Choosing exact-one candidate-row matching aligns with RDR 0001's deterministic
kernel and deliberately diverges from first-match FSM engines: priority order is
convenient in code, but it makes review harder and lets source or dump
reordering change behavior. Gate-then-count is the kernel's ordering and the
only one under which missing artifact state cannot hide behind an escapable
refusal. Runtime FSM libraries are kept out of the core because they would
either drive orchestration or hide the contract in host callbacks; their useful
role is vocabulary, validation, and visualization.

Premortem: this could fail if the sparse source becomes a verbose
pseudo-language or if expansion hides surprising implicit rows. The
recommendation survives only if Resolve proves the RDR and kata examples with a
parse-normalize-dump spike. If that fails, the source schema must be reduced
before lock.

## Alternatives Considered

### Alternative 1: Sparse TOML Model With Expanded Table Dump

**Description**: Hand-authored TOML files declare tags, recognized outcomes,
shared contexts, accessor references, and sparse transition rules. The tool
normalizes rules into explicit candidate rows and can dump that expanded table.

**Pros**:

- Avoids Cartesian-product authoring for dimensional models such as RDR status
  plus profile plus prelock iteration plus guards.
- Keeps the source reviewable while still giving lint and resolver code a fully
  explicit row set.
- Supports mechanical dumps, diagnostics tied back to stable source rule ids,
  and future graph export by RDR 0006.
- Fits Go CLI implementation with ordinary typed decoding and validation.

**Cons**:

- Normalization is a real contract, not just parsing.
- Exact field layout needs spikes before lock.
- TOML is not a formal state-machine standard, so graph validation must be built
  over the parsed representation.

**Reason for selection**: Best balance of reviewability, parse simplicity, and
structured expressiveness for the target RDR and kata models.

### Alternative 2: Fully Expanded TOML Row Table

**Description**: Authors maintain one TOML row per candidate edge, with all
dimensions repeated inline.

**Pros**:

- Simplest parser and easiest mental model for tiny graphs.
- The source file is already the table the resolver sees.

**Cons**:

- Explodes for the RDR model once status, profile, prelock iteration, cluster
  gates, rewind scope, and guard dimensions interact.
- Repetition makes edits risky: changing one shared condition requires finding
  every copied row.
- The source becomes a generated-looking artifact even though humans are
  expected to own it.

**Reason for rejection**: It is acceptable as a dump format, not as the
authoring format.

### Alternative 3: JSON Model

**Description**: Store the same sparse model as JSON.

**Pros**:

- Simple to parse and strict about data types.
- Easy for tools to generate and consume.

**Cons**:

- Poor hand-review ergonomics: comments are unavailable, trailing comma churn is
  common, and nested objects become noisy for prose-heavy transition data.
- Encourages machine-generated artifacts, which conflicts with the goal that the
  legal graph be authored and reviewed directly.

**Reason for rejection**: It optimizes interchange over the human review loop
that is the central user outcome.

### Alternative 4: CSV / Matrix Table

**Description**: Store transitions as rows with columns for current state,
outcome, guard columns, and writes.

**Pros**:

- Compact and easy to scan for small state machines.
- Familiar representation for simple transition matrices.

**Cons**:

- Weak fit for nested predicates, typed values, accessor references, and
  multi-tag writes.
- Escaping and comments become awkward exactly where RDR/kata examples need
  explanation.

**Reason for rejection**: Too scalar for the required tag-set model.

### Alternative 5: Small DSL

**Description**: Define a custom text grammar such as `match -> writes` with
inline predicates.

**Pros**:

- Can be concise and domain-specific.
- Could make graph-like edges visually obvious.

**Cons**:

- Requires custom parsing, error recovery, formatting, editor support, and
  long-term grammar ownership.
- Pushes this RDR into language design before the target tables are proven.

**Reason for rejection**: The parser and tooling burden is not justified while
TOML can carry the same semantics.

### Briefly Rejected

- **Generated Go tables**: Fast and type-safe, but the reviewable artifact would
  be code, not the legal graph as data.
- **SCXML/XState as the source format**: Strong standard/tooling story, but a
  poor fit for tag provenance and the project's non-orchestrating resolver
  boundary.
- **Embedded host predicates**: Expressive, but defeats static lint and makes
  graph review depend on reading arbitrary code.

## Trade-offs

### Consequences

- The legal graph becomes a first-class reviewed model with stable rule ids.
- The expanded table becomes a generated diagnostic view, not hand-maintained
  source.
- Parser and lint errors become part of the CLI surface even though this RDR is
  primarily an internal data-format decision.
- Some expressiveness is intentionally deferred to RDR 0003 so this table stays
  statically checkable.
- Every rule must name its outcome; the model has no eventless rows because the
  kernel has none.

### Risks and Mitigations

- **Risk**: Sparse TOML rules become too verbose or too magical for large
  graphs.
  **Mitigation**: Resolve must encode representative RDR and kata fixtures and
  reject the shape if reviewers cannot scan the source or explain the dump.
- **Risk**: Sparse contexts hide an accidental Cartesian product.
  **Mitigation**: The dump must show normalized candidate rows with source rule
  ids. Expansion count per rule is then **derivable from the dump itself** —
  rows sharing a source rule id, which the `identity` and `source` columns
  already carry — so this mitigation rests on this RDR's own deliverable and
  needs no lint output. It is deliberately *not* routed to RDR 0006: that RDR's
  advisory tier is normatively closed at `graph-coverage-closed-by-escape`,
  `graph-redundant-row`, `graph-unreachable-rule`, and `graph-vacuous-atom`, so
  an expansion-count report is not a finding it can mint.
- **Risk**: Future contributors treat row order as priority.
  **Mitigation**: Normative exact-one semantics require explicit escape rows for
  modeled no-match or ambiguous-match behavior; graph lint rejects overlapping
  ordinary or escape rows instead of picking the first match.
- **Risk**: Accessor references pull execution semantics into the table.
  **Mitigation**: The table only names accessor bindings; RDR 0004 owns execution
  and read-back behavior.
- **Risk**: Normalizer output drifts from the kernel's exported constants or
  key identity, so a valid-looking atom fails closed at runtime.
  **Mitigation**: The MVV asserts emitted existence atoms byte-for-byte against
  `resolve.OpExists`/`LiteralTrue`/`LiteralFalse` and asserts row keys equal
  declaration keys.
- **Risk**: Dump order tracks byte positions in the authored source, so adding
  a comment or reordering unrelated rules churns every downstream golden test.
  **Mitigation**: The source locator is excluded from the row-sort key, and
  the identity tuple that replaces it carries no positional component. The MVV
  permutes rule declaration order and asserts the normalized value is
  unchanged, so a later identity change cannot reintroduce a positional
  tiebreak unnoticed.
- **Risk**: This RDR is the cluster's wire-format producer and locks after
  several of its consumers (0003, 0004, 0006, 0007, 0008, 0009 are `Final`),
  so a layout or routing change here can leave a peer's cited fact stale.
  **Mitigation**: The layout, routing key, and sentinel are fixed by JDR 0001
  §D5–§D7 and this RDR cites them rather than deriving them; cross-RDR drift
  is checked by `/rdr-cluster-reconcile` before implementation, and peers
  absorb by citation at their own re-lock.
- **Risk**: A `keys` binding that must cover every declared key would make
  `recognized` and caller-supplied `--tag` keys unauthorable.
  **Mitigation**: The binding validation is provenance-scoped — exactly one
  reader per owned key, at most one per observed key, none for `recognized` —
  and scenario 3 carries a fixture for each arm.

### Failure Modes

Malformed TOML, unknown schema fields, unresolvable or cyclic context
references, malformed tag or accessor declarations, an accessor binding that
leaves an owned key unread or a written key unwritten, malformed predicate
atoms (including a non-`eq`/`in` operator under a match block or a literal
outside its domain), the reserved `<clear>` value authored as a tag value,
malformed escape declarations, malformed outcome bindings, a malformed
`[initial]` assignment, a misnamed recognized declaration, or a write to an
undeclared or non-owned tag fail at load/validation time with stable CLI
errors — each is decidable from one rule plus the declarations. A missing root
or stop set is RDR 0006's blocking finding. A row gap, overlap, dead row,
read-before-write condition, or ambiguous expansion requires comparing
normalized rows and is therefore an RDR 0006 lint failure before the model is
accepted. Expansion count per rule is derived from the dump — rows sharing a source rule
id — not a lint diagnostic and not a refusal: this RDR sets no expansion
threshold, and RDR 0006's advisory tier is closed against adding one. At runtime, unknown outcomes,
unavailable accessor inputs, absent owned state, undecidable guards, and
unmodeled zero/multiple survivors are typed resolver refusals rather than
guessed edges; modeled zero/multiple-match behavior must come from exactly one
surviving escape row, and the gate refusals are never modeled.

## Implementation Plan

### Prerequisites

- [x] All Critical Assumptions verified **except A9, A10, and A12**. A1 was
  re-verified at Stage 4 by re-authoring both spike fixtures to the §D7 layout
  and re-running the spike under strict decoding; A7, A11, and A14 were
  promoted to `Verified` by that same run (general match-block expansion and
  the sequence suffix; escape expansion under `in`; the
  version-gate precedence, witnessed by a v2-shaped fixture against a v1
  control). The three that remain `Pending` carry named MVV assertions:
  **A9/A12** are blocked on RDR 0007's kernel reshape (below) — their fixture
  inputs now exist, but no field exists to route into; and **A10**'s decidability
  half is source-verified and only its paired-document behavior is owed, since
  no single-file fixture can express two model ids. **A13 and A15 were the two
  carve-outs this list previously carried, and Stage 6 closed both rather than
  deferring them.** Each had a wide arm that changes match semantics while
  tripping **zero** load categories — meaning the owed fixture was the *only*
  detector, and deferring the fixture deferred the only detector. Both were
  fixed in the spike and given discriminating controls in the same pass
  (`gen-cases.py` now mutates guard blocks and mints the delimiter-bearing
  cases), with the baseline digest unchanged and all 43 negatives still
  refusing. The three that remain are not load-bearing for a *pre-lock* MVV
  assertion and are survivable per their "If wrong".
- [ ] RDR 0007 is the normative home of the atom shape and the existence
  constants this RDR's normalizer emits. RDR 0007 is `Final`, so the
  *specification* is settled — but its kernel reshape is unimplemented, and the
  gap is the whole atom surface, not only the constants. On `main`
  `internal/resolve/resolve.go::Row.Guard` is a `string`, `Row.Match` is a flat
  `[]Tag` equality list matched by `::TagSet.matches` with no operator
  dispatch, and there is no `Atom` type, no operator-token type, no per-atom
  block field, and no `OpExists` / `LiteralTrue` / `LiteralFalse` anywhere in
  the repo. Of Phase 2's four deliverables only outcome lifting and
  `RequiresOwned` derivation have fields to land in; per-atom block retention
  and verbatim existence constants have none. Phases 2 and 3 therefore sequence
  behind the reshape in their entirety, and the Testing Strategy assertions that
  cite `resolve.OpExists` / `LiteralTrue` / `LiteralFalse` are unsatisfiable
  until it lands. This item gates implementation sequencing, not lock. It does not
  invert the cluster's ordering: RDR 0008 and RDR 0009 both name this RDR's
  normalizer as the enforcement point their own checks land inside, and this
  RDR's normalized rows target RDR 0001's shipped kernel row. **The reshape has
no owner**, and that is the open item this bullet leaves behind: RDR 0007 is
`Final` and therefore will not itself schedule the work, no kata tracks it, and
no phase in this RDR's plan claims it. Sequencing behind an unowned prerequisite
is indefinite, not merely ordered. Implementation MUST NOT open Phase 2 by
mirroring the constants locally to get moving — the Risks section names that
drift and its only mitigation is that the real constants compile once exported.
Landing the reshape, or explicitly assigning it, is the first implementation
action and belongs to whoever opens Stage 8 on this cluster.
- [x] RDR 0001 remains aligned on exact-one stateless resolution.
- [x] RDR 0003 is coherent enough for this RDR to defer the fixed predicate
  operator set and tag type model without importing that grammar.
- [x] RDR 0007, 0008, 0009 are Final; their producer obligations on this RDR are
  bound above.

### Minimum Viable Validation

Parse two hand-authored sparse TOML fixtures, one for a representative RDR flow
slice and one for a representative kata flow slice, into typed source data;
normalize them into candidate rows; dump the expanded table; validate tag
declarations, context references, predicate references, recognized outcomes,
supported model version, multi-tag writes, outcome binding, and escape rows;
then prove by unit test that one sample tag-set resolves through
`internal/resolve::Resolve` to exactly one ordinary row, one unmatched sample
tag-set resolves to exactly one modeled escape row when the table declares one,
one unsupported-version variant is refused before normalization, and one
deliberately overlapping variant is reported by lint as an ambiguous overlap.

The overlap item is the one MVV assertion this RDR cannot discharge alone: the
ambiguous-overlap check is cross-row and therefore RDR 0006's by the arity
split, and RDR 0006 is unimplemented. It is **deferred to the Phase 5 lint
handshake rather than counted as satisfied here** — the rest of the MVV is
executable against this RDR's own load and normalize paths. What this RDR owes
before that handshake is that the normalized rows *carry enough structure* for
the check: overlapping rows must be distinguishable by row identity and
comparable by predicate set, which the identity tuple and the atom-set contract
above already fix. Marking the overlap assertion green before RDR 0006 lands
would be asserting on a stub.

The RDR fixture must cover at least `Status`, `Profile`, prelock iteration, one
guard that actually carries a predicate, one existence atom, one `in` atom on
`recognized` that expands into more than one row, and one explicit `no_match`
escape row. It must additionally carry **two sibling candidate rows binding the
same outcome** — without them the gate's ordering is untestable, since a
one-candidate-per-outcome fixture can never exercise "an unevaluable survivor
refuses even though a decidable sibling and a modeled escape both exist," which
is the property this RDR's row shape exists to support.

**The fixtures now carry those siblings; the re-author-and-extend was done at
Stage 4 (A1).** The pre-§D7 fixtures authored `[accessors.<id>]`, a tag-side
`accessor` reference, and no `[initial]`/`terminal`, and bound `round-clean` on
one ordinary rule and one escape rule, so no two *ordinary* candidates ever
contended — which is why the desk trace used to record "no witness possible in
the current fixture" for the sibling-gate row. Those fixtures are now refused
outright by this RDR's own contracts (`unknown schema field` on
`accessors.*` and `tags.*.accessor`), which is the sharpest statement of why
they could not be promoted.

`evidence/spikes/iter-2/` is the promoted set: both fixtures re-authored to the
closed layout, plus `continue-prelock-cluster` — a second ordinary rule binding
`round-clean`, which `continue-prelock` already binds, whose `guard.all`
carries an `eq` atom over the optional (`required` unset) owned key
`cluster_ready`. That single addition supplies the MVV's two-sibling
requirement and A12's discriminating witness at once. The Load-Bearing
Decisions' "canonical examples" clause fixes the authoring idiom those fixtures
demonstrate, not their row census; a fixture may not be narrowed on promotion,
only extended.

**Every oracle must have a failing control.** An assertion that passes by
absence-of-error, exit-0, or fixture-name match does not discharge an MVV item:

- Each load-time validation category is asserted by **category**, not by
  message text, and each has one mutated fixture that trips it and no other —
  otherwise an implementation collapsing all categories into one code passes.
- Determinism and key-order assertions are paired with the value assertion
  above; alone they are satisfied by a normalizer emitting a constant or
  dropping every `unless` atom. The permutation control runs on the RDR
  fixture, which carries the escape row, the `in`-expansion, `<clear>`, and
  inherited contexts — not only on the simpler kata fixture.
- The existence-atom test exercises `exists = true` **and** `exists = false`; a
  fixture carrying only one lets a hardcoded literal pass.
- The byte-equality contract is asserted by a negative control — a reference
  spelling a declared `status` as `Status` must fail `unknown tag` — since a
  normalizer applying case folding passes every positive test.
- Unknown schema fields are asserted by a fixture carrying one; the decoder
  must be configured to reject unmapped keys, which is not the TOML library's
  default and is otherwise a silent no-op.

The dump ordering is MVV-covered, not left to implementation: assert that the
normalized value is unchanged when the same model is authored with its TOML
keys in a different order, **and when its rules are declared in a different
order**, and that repeated dumps of one model are byte-identical. Row-order
permutation is the assertion that actually excludes a positional tiebreak;
sorting and checking that no adjacent pair compares equal does not, because a
comparator falling back to table position produces no equal adjacent pair
either — it is satisfied by the defect it was meant to catch.

### Phase 1: Fixture and Schema Spike

**Discharged at Stage 4** — `evidence/spikes/iter-2/` holds the RDR and kata
fixtures authored to the closed §D7 layout (capability tables with `keys`,
`[initial]`, `terminal`, `[model.metadata]`, the type-model keys), keeping
shared contexts, a rewind with an explicit clear, a gate reference, a
profile-dependent branch that expands under an inherited match-block `in`, an
escape rule that itself expands under `in`, a multi-tag write, and the two
same-outcome sibling rows the MVV requires. Phase 1 remains listed because
implementation promotes these fixtures into the production test tree; it does
not re-derive them.

### Phase 2: Normalizer and Dump

Introduce typed source structures, normalization to `internal/resolve::Row`
values (outcome lifted, per-atom block retained, `RequiresOwned` derived,
existence constants emitted), and an expanded-table dump with deterministic
ordering and source rule ids.

### Phase 3: Parser and Validation Skeleton

Add validation rules for tag and accessor declarations, the provenance-scoped
`keys` bindings, the reserved `recognized` name and `<clear>` value, rule ids,
context references, `[initial]` assignments and `terminal` ids, predicate
atoms (including the match-block operator restriction and domain membership),
outcome bindings, gate references, write targets, explicit clears, escape
declarations, and expansion-count diagnostics.

### Phase 4: Resolver Handshake

Connect the normalized rows to RDR 0001's gate-then-count exact-one contract
without adding runtime ordering or host-code predicate callbacks. Internal
indexes may be decision trees or tries, but they must preserve the normalized
semantics.

### Phase 5: Lint Handshake

Expose the parsed representation needed by RDR 0006 for graph determinism,
reachability, and read-before-write checks.

### Day 2 Operations

| Resource | List | Info | Delete | Verify | Backup |
| --- | --- | --- | --- | --- | --- |
| Transition model files | In scope via normal repository listing | In scope via parse/lint/dump output | N/A; source-controlled files | In scope via lint | N/A; source control is backup |
| Expanded table dumps | In scope through dump command | In scope through source rule ids | Delete/regenerate | In scope via dump tests | N/A; generated from source |

No runtime persistent resource is introduced by this RDR.

### New Dependencies

Use `github.com/pelletier/go-toml/v2` as the TOML parser candidate. The Resolve
spike ran against v2.3.1 from the local module cache, and the module license is
MIT. No production dependency is added until implementation.

The strict-decoding obligation this format states is satisfiable with that
library, confirmed against its source rather than its documentation:
`(*Decoder).DisallowUnknownFields()` returns `*StrictMissingError`, whose
`DecodeError.Key()` and `Position()` name the offending key and location — so
`unknown schema field` can carry a key, not just a message. Two properties the
design depends on were checked at the same time: the strict check descends only
into struct branches, so a `map[string]any` field passes untouched (this is what
makes `[model.metadata]` free-form while the rest of the document stays strict),
and a two-pass decode over one `[]byte` is re-runnable (this is what makes the
version gate expressible ahead of strict decoding). `toml.Unmarshal` is
permissive by default, which is precisely why strictness is stated as an
obligation on this format rather than assumed from the parser.

When the dependency lands it also settles the open TOML-library choice at
`internal/cli/config/config.go::Load`, which carries a `TODO` for one; the two
consumers should share a single library rather than each picking independently.

## Validation

### Testing Strategy

Implementation tests must promote the Resolve spike into production fixtures:

1. **Scenario**: Parse the RDR and kata sparse TOML fixtures from `docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/iter-2/` into typed source structs.
   **Expected**: Tag declarations (including `[tags.recognized]` and the seven type-model keys), root recognized-outcome alphabets, `[initial]` and root `terminal`, `[read.*]`/`[write.*]`/`[gate.*]` tables with `keys`, rule-level `gate` lists, shared-context inheritance, positive/negative guards, explicit clears, escape declarations, and multi-tag writes decode without ambiguous field placement under strict decoding. **`[model.metadata]` is asserted by deep equality against the authored table** — every value and every nesting level, not its top-level key names: a loader that decodes the table and discards its values passes a key-name oracle, and since metadata sits outside the dump field list and the Round-Trip invariant, this is the only assertion that observes it. The `source` and `description` keys decode onto the rule and model structs rather than refusing as unknown schema fields.
2. **Scenario**: Normalize the RDR fixture's `continue-prelock`, `continue-prelock-cluster`, `reconcile-rewind`, `terminal-archive`, and `draft-no-match-escape` rules and the kata fixture's `review-accepted` and `review-needs-work` rules.
   **Expected**: Candidate rows retain source rule ids/source locators, inherited predicates are expanded, each atom reports its authored block (`match`/`all`/`unless`), the single `recognized` atom is lifted into the row's outcome field and absent from the predicate set, `RequiresOwned` equals the sorted write-plus-clear key set (empty on the escape row), gate ids are carried, escape rows retain their modeled failure class list — from which the derived `escape` kind is computed — and carry neither writes nor next-state tags, and writes are deterministic. `terminal-archive`'s `recognized.in` expands to two rows whose suffix is the one-element sequence of the outcome literal; `continue-prelock`, inheriting `large-prelock`'s `profile.in = ["large", "foundational"]`, expands to two rows carrying `profile.eq=<member>@match` and the member as their suffix; a rule with two multi-member match `in` atoms yields their product with a two-element suffix in atom sort order. The normative fixture for this scenario is `evidence/spikes/iter-2/output.txt` (approved Stage 4, A1; SHA-256 `6ccfe9012b0705ef4d4b3d1c620daffd69523436175120be1bea8a05df9c55dd`), whose nine RDR rows and two kata rows are the expected **row set and row identities** — not the expected bytes. **This SHA MUST NOT be asserted as a golden hash by any implementation test**, and the reason is *not* that the bytes are wrong. Stage 6 fixed the joined-rendering defect (A13) and the digest came out **unchanged** — `6ccfe901…` before and after — because these two fixtures carry no set literal or write value whose members contain the joining delimiter, which is exactly why they could not discriminate the defect in the first place. The ban stands on the oracle's *character*, not its current value: a rendered-text hash is satisfied by any renderer that happens to agree on these inputs, so it witnesses nothing about the identity rules the clauses actually bind, and a legitimate rendering change (spelling a set visibly as a set, as the display clause encourages) breaks it for no semantic reason. Assert over the normalized value, per the dump clause's own "assert over the normalized value … reserve rendered-text goldens for tests of rendering itself." The SHA is retained as provenance for what Stage 4 reviewed and what Stage 6 re-ran. Note the RDR fixture yields **nine** rows, not seven: `continue-prelock` and `continue-prelock-cluster` each expand on the inherited `profile.in`, and `draft-no-match-escape` expands on its `recognized.in`.
   **The no-alias obligation on next-state tags and writes is not assertable by value comparison** — the two fields hold equal sets under this RDR's authoring surface, so an aliased pair and an independently built pair compare equal. On `main` both are `[]Tag` slices, so an alias also shares a backing array — which is exactly what makes the obligation assertable after all. **The control is a mutation test**: normalize a row, mutate one field in place, and assert the other is unchanged. An aliased pair fails it; an independently built pair passes. Enforcement is therefore by this test, not by review of the normalizer — this RDR's own `oracle` rule refuses to accept an assertion until a wrong implementation that fails it is named, and "review the normalizer" names none. Value comparison alone remains insufficient, which is why the control mutates rather than compares.
   **Additionally**, atom-set semantics are asserted positively, since these shapes must *survive* load rather than be refused: two inherited contexts contributing the same key and operator with **different literals** yield **two** atoms in the normalized set, not one — asserted by count and by value, with the contexts listed in both `use` orders to prove the result is order-independent (a merge keyed on `(block, key, operator)` drops one and passes a count-only check on a single ordering); a set-valued literal whose member contains a space normalizes to a member sequence, so `["needs work"]` and `["needs", "work"]` are **distinct** atoms; and a rule authoring one atom in both `all` and `unless` loads successfully as two atoms and is reported by lint as a dead rule, not refused at load.
   **The distinct-literal control asserts on the atom count, and covers the delimiter space rather than guessing one.** Two contexts contributing `in = ["a,b", "c"]` and `in = ["a", "b,c"]` on one key and block MUST yield **two** atoms in the normalized set. That count is the whole oracle: it is a property of this RDR's own normalized value, readable without RDR 0006, and it fails against any joined-rendering key. The earlier phrasing — "members containing the implementation's own joining delimiter" — was not assertable, because a black-box test cannot read the delimiter a wrong implementation chose; the control instead **parameterizes over the plausible delimiters** (`,` `;` `|` space, and the empty string) and asserts two atoms for each, so no single choice of join survives. Do **not** state the assertion as "the rule is a dead rule for lint": that is an RDR 0006 verdict, and the MVV explicitly refuses to count assertions against an unimplemented lint (the overlap item, above). Deadness is the *consequence* that makes the defect matter; the atom count is the test. A control using only scalar or space-bearing members does not discriminate: it passes against a merge keyed on a comma-joined rendering, which is what the spike did before Stage 6 (`main.go::Atom.identity` → `literalString`), collapsing the pair to one atom and normalizing the rule to a live expanding row. **The parameterization is not belt-and-braces**: that pre-fix spike collapsed *only* the comma case, so the four other delimiters returned two atoms while the defect was live — a single-delimiter control chosen unluckily would have certified it. Stage 6 ran all five on both arms and A13 is `Verified`; the fixtures live in `evidence/spikes/iter-2/delim/` and implementation promotes them.

   **The same control runs on the write side.** A `set`-kind write authoring `["a,b", "c"]` and one authoring `["a", "b,c"]` MUST normalize to two distinct write values, parameterized over the same delimiter set. The spike joined write values before Stage 6 (`main.go::renderValue` → `strings.Join(parts, ",")`), so the two set spellings collided — the literals-clause defect reached through the write field, where RDR 0004's read-back compares for equality. Closed by giving `Write.Value` a member sequence and `Write.identity` a length-prefixed key; A13 `Verified`. Asserted by value inequality, a positive assertion, since the defect trips no category. **The assertion must read the normalized value, not the rendered one** — Stage 6's `|` case showed a render can still collide after the identity is correct, so an oracle reading rendered text can pass a defect the identity has already fixed; the display clause now requires each member to be delimited unambiguously for exactly this reason.
   **The idempotence mirror of that assertion is required too (A13)**: a rule and an inherited context contributing a **byte-identical** `(block, key, operator, literal)` atom collapse to **exactly one** atom, asserted by count on that key. Without it a normalizer that de-duplicates nothing passes every control above, since they only ever count the two-distinct-literal case. Both halves are witnessed by the spike on scalar literals — `evidence/spikes/iter-2/merge-idempotent.toml` (count one) and `merge-distinct.toml` / `merge-distinct-rev.toml` (count two, both `use` orders digesting identically). Those witnesses answer the `(block, key, operator)` question only; the delimiter-bearing control above is what closes the full-identity question, and Stage 6 ran it — parameterized over five delimiters, both arms, in `evidence/spikes/iter-2/delim/` — which is why A13 is now `Verified`. Implementation promotes all fifteen fixtures.
   **Escape expansion under `in` is asserted here (A11)**: the fixture's `draft-no-match-escape` binds its outcome with `in` over two alphabet members, and normalization must yield **two** escape rows carrying distinct expansion suffixes, each with an empty write set and each retaining the modeled failure-class list. Witnessed in `evidence/spikes/iter-2/output.txt` as `rdr.draft-no-match-escape#round-clean` and `…#reconcile-block`.
   **Handoff routing is asserted as total and disjoint (A12)**: constructing the `resolve.Row` for each normalized RDR-fixture row, the union of `Match` and the guard's atoms equals the normalized atom set and their intersection is empty, and every `Match` atom carries block `match`. The discriminating cases are now authored in the fixtures: `continue-prelock-cluster`'s `[rule.guard.all.cluster_ready] eq = true` — an `eq` atom under `guard.all` over an optional (`required` unset) owned key — and the kata fixture's `status.eq=closed@unless`; both are equality operators that route to the guard because of their block, and both appear as guard-block atoms in `evidence/spikes/iter-2/output.txt`. Like the rest of Phase 2 the assertion itself is unsatisfiable until RDR 0007's reshape lands (Prerequisites), but its inputs are no longer owed.
3. **Scenario**: Validate one malformed variant per load-time category — unknown tag written, unknown tag matched, unknown context, cyclic context inheritance, writes to non-owned tags, unknown accessors, unsupported versions, unknown schema field, malformed predicate atoms (unknown operator; `exists` with a non-boolean literal), malformed escape declarations (an empty write block on an escape rule, an escape rule carrying a `clear` list, and one carrying a `gate` list — the normalizer is the only enforcement point for a write-free escape row, so each shape needs its own **fixture**; note the three share one category, so the category oracle alone cannot tell them apart. Each fixture's discriminating assertion is therefore that it refuses **and** that the other two shapes are absent from it — the one-mutation rule (`gen-cases.py`) supplies that by construction. A finer split into three categories is deliberately not minted: the shapes are one authoring error — a plan-bearing field on a plan-free rule — and RDR 0009 binds the write-free obligation as a single kernel-boundary obligation), a rule with zero or two `recognized` atoms, a `recognized` atom authored under `guard.all` and one under `guard.unless`, a rule id containing `#` and an alphabet member containing `#`, an outcome literal outside the alphabet, an alphabet containing the empty string, a duplicate alphabet member, an empty alphabet, an owned declaration named `recognized` and a recognized declaration named `outcome`, a duplicate rule id, a duplicate model id (a **pair** of documents sharing a `model id`, loaded into one invocation — the only surface on which the category is decidable), an ordinary rule carrying no write block, a missing root outcome alphabet, a `lt` atom and separately an `exists` atom under `[rule.match]` (match-block operator restriction), a `#` inside a match-block `in` member, a predicate literal outside the tag's declared `domain`, a tag declaration RDR 0003's rules reject (`domain` on a `kind` that admits none), `<clear>` authored as a write value, as an `[initial]` value, and as a predicate literal, an `[initial]` key that is undeclared, an owned key served by zero readers and separately by two, a written key with no writer, a writer whose `keys` names an observed tag, `recognized` in a reader's `keys`, an accessor entry missing `timeout`, a rule whose `gate` names an id declared only under `[read.*]`, an escape rule carrying a `gate` list, and a `terminal` entry naming no context. Positive controls in the same fixture set: an observed key served by no reader loads (JD-9), and a model with no `[initial]` loads and is left to lint.
   **No control is listed for an `[initial]` key that is *observed*, because that document has no authorable form** (Normative Contracts, Root and stop set): the writer-binding bullets refuse it as `malformed accessor binding` or `write to non-owned tag` before the owned-tag predicate is reached, so a fixture for it would witness a different category than the one it names. `malformed initial declaration` is controlled by its **value arm only** — a literal ill-formed for the declared kind or domain (`neg-initial-bad-value`, `neg-initial-int-out-of-range`). The undeclared-key variant is *not* one of its controls: an `[initial]` key naming no declared tag refuses as `unknown tag` (`neg-initial-unknown` — "unknown tag \"nosuchtag\" written"), which is the correct category, since the defect is the undeclared tag rather than the `[initial]` declaration's shape. Scenario 3 lists that variant under `unknown tag`, not here; assigning it to `malformed initial declaration` would fail the category-not-message oracle against the fixture that exists. Verified by enumeration this stage: both authorings were built and each refused under the writer-binding category (`evidence/spikes/iter-2/neg/probe-a-initial-observed-no-writer.toml`,
   `neg/probe-b-initial-observed-with-writer.toml`).
   **Expected**: Each variant is refused with the **one** category its mutation targets and no other — the assertion is on the category, not the message text, so an implementation collapsing several categories into one code fails. Ambiguous overlap is deliberately absent from this scenario: it is cross-row and therefore an RDR 0006 lint finding (scenario 4), not a load failure. **Forty-one of these controls are witnessed by the §D7 spike**, one refusal per mutated fixture, in `evidence/spikes/iter-2/negative-cases.txt`, covering **18 of the 25 categories** named above. (Thirty-five before Stage 6; the six guard-block controls A15 owed deepen three categories already covered — `malformed predicate atom`, `reserved tag value`, `unknown tag` — so the fixture count rises and the category count does not.) The **seven** still owed at implementation are `malformed TOML`, `missing recognized outcome alphabet`, `cyclic context inheritance`, `malformed tag declaration` (RDR 0003 supplies its rejection rules), `duplicate model id` (cross-document, so it needs the paired-document surface no single-file fixture can provide), and the two categories the 3amigo pre-lock pass added — `malformed model declaration` (absent `[model]`, `id`, `version`, or `[tags.recognized]`; the spike today refuses the first three off-contract and folds an absent `version` into `unsupported version 0`) and `malformed dump declaration` (an `order` naming an unknown, repeated, or omitted column; the spike decodes `[dump]` and never reads it). Both are owed a fixture **and** an implementation. That is seven: five pre-existing plus the two added. Two of the seven are owed only a **fixture**, not an implementation: `cyclic context inheritance` and `malformed tag declaration` are both already decided by the spike and were confirmed to fire against mutated fixtures (`cyclic context inheritance at "prelock"`; `malformed tag declaration: "cluster_ready" has no kind`), so `gen-cases.py` should mint their single-mutation cases rather than leaving them to implementation discovery. (Stage 6 extended `gen-cases.py` for the guard-block and delimiter controls; these two remain unminted and are the cheapest outstanding fixture work.) Regenerate the whole set with `python3 evidence/spikes/iter-2/gen-cases.py`, which derives each fixture from the RDR fixture by a single mutation so the one-defect-per-fixture rule holds by construction. Implementation promotes these fixtures; it may extend the set but must not narrow it.

   **Guard-block controls are owed for every atom-level rule (A15).** `gen-cases.py` mutates match blocks only, so no existing control witnesses the block-agnostic clause and the coverage figures above overstate what is proven: an atom rule counted "witnessed" is witnessed in one of three blocks. Each of operator membership, `<clear>`, domain/kind conformance, `#` reservation, and tag-key declaration owes a `guard.all` and a `guard.unless` variant. The three reproduced today — `frobnicate = "small"`, `eq = "<clear>"`, `eq = "NOT_A_PROFILE"` under `[rule.guard.unless.profile]` — are the minimum starting set and all three currently load clean.

   **Fixture conformance is a control, not a review step.** Both normative fixtures MUST load clean under the finished loader, asserted as part of scenario 1. This is the check that would have caught the `[dump]` column vocabulary being fixed to spellings the fixtures do not use, and it is cheap: the fixtures are already parsed by every other scenario.
   **The version gate's precedence needs its own fixture (A14), separate from the `unsupported version` category fixture above.** That fixture is a v1-shaped document with a bad version value, which trips the category under *either* check ordering and therefore witnesses nothing about precedence. The precedence control is a **v2-shaped** document — `version = 2` plus a v2-only key the strict decoder would reject as an unknown schema field — asserted to refuse `unsupported version` and **not** `unknown schema field`. This is the only ordering this RDR fixes normatively, so it is the only ordering that gets an oracle.
4. **Scenario**: Run `internal/resolve::Resolve` over one matching ordinary tag-set in which every sibling candidate's guard is decidable, one tag-set with no ordinary match but one matching `no_match` escape row, and one tag-set in which a sibling candidate's guard is unevaluable because its `guard.all` carries an `eq` atom over an optional owned key the view does not hold — all three drawn from the fixture's two same-outcome sibling rows. Separately, run RDR 0006's lint over a deliberately overlapping variant.
   **Expected**: The matching ordinary tag-set resolves to one transition row; the no-match tag-set resolves to the modeled escape disposition; the unevaluable-sibling tag-set refuses `guard_unevaluable` even though a decidable sibling and a `no_match` escape row exist; zero or multiple survivors without exactly one surviving escape row are refusals and never fall back to row order. The overlapping variant is reported by lint as an ambiguous overlap before the model is accepted, naming both rows.
5. **Scenario**: Normalize and dump semantically identical variants of the **RDR** fixture — one with its TOML keys authored in a different order, one with its `[[rule]]` blocks declared in a different order, and one spelling a single-outcome binding `in = ["x"]` where the original spells `eq = "x"` — and dump the same fixture repeatedly in one process and across processes.
   **Expected**: The normalized value is identical in every case because rows sort by the identity tuple and atoms/gate ids/next tags/writes/required-owned keys/escape classes sort by key, and because a single-member `in` contributes nothing to the expansion suffix — on `recognized` and on any other match-block tag alike; repeated dumps are byte-identical despite Go's randomized map iteration. Rule-order permutation is what excludes a positional tiebreak. Fully witnessed by the §D7 spike in `evidence/spikes/iter-2/negative-cases.txt`: all three permutations of the **RDR** fixture — `perm/rdr-keyorder.toml`, `perm/rdr-ruleorder.toml`, and `perm/rdr-eq-as-in.toml` — digest identically to the baseline (`2ff26b0add37…`), and three consecutive runs over both fixtures are byte-identical (`6ccfe9012b07…`). The RDR fixture is the one carrying the escape row, the `in`-expansions, `<clear>`, and inherited contexts, so the permutation control runs on the demanding fixture rather than the simpler kata one.
6. **Scenario**: Normalize a rule carrying `exists = true` and one carrying `exists = false`, one carrying a non-boolean existence literal, and a variant spelling a declared tag `Status` where the declaration is `status`.
   **Expected**: Emitted atoms carry `resolve.OpExists` and `resolve.LiteralTrue`/`LiteralFalse` byte-for-byte, referencing the kernel constants directly rather than a local mirror — so this assertion compiles only once RDR 0007's reshape exports them (Prerequisites), and is owed at that point rather than satisfiable today. The non-boolean literal is refused at load as a malformed predicate atom; the mis-cased reference fails `unknown tag` rather than folding. The RDR fixture's `continue-prelock` guard carries the `finalized_at.exists=false@all` atom the spike emits; the spike mirrors the constants locally because the kernel does not yet export them.

### Performance Expectations

Resolve evidence is functional rather than throughput-oriented. The spike
(§D7 closed layout, `evidence/spikes/iter-2/`) normalizes representative RDR
and kata sparse fixtures into eleven deterministic rows — nine RDR, two kata —
comprising two escape rows from one escape rule expanding on a `recognized`
`in`, and two pairs of ordinary rows expanding on an inherited
non-`recognized` `profile.in`. Three consecutive runs produced byte-identical
output with SHA-256
`6ccfe9012b0705ef4d4b3d1c620daffd69523436175120be1bea8a05df9c55dd`, and three
semantics-preserving permutations of the RDR fixture — different TOML key
order, different `[[rule]]` declaration order, and `in = ["x"]` for
`eq = "x"` — each produced a matching digest.

Determinism checklist, as run: source key order is neutralized by sorting every
emitted sequence; Go's map iteration is never the emission order (rows, atoms,
next tags, writes, required-owned keys, and escape classes are all sorted
slices before rendering); string comparison is byte-lexicographic per the Go
spec, so ordering is locale- and platform-independent; set-valued literals are
rendered with their members sorted so two authored orderings of one set are one
literal; no case folding, trimming, or namespace rewriting is applied at any
stage; the expansion suffix is absent or an alphabet member, never the empty
string, so the identity tuple is total. No hash is part of the contract: the
SHA is evidence for this spike output only, and production golden tests must
assert the normalized expanded-table value the normative contract defines.
Runtime lookup may index rows later, but that optimization must preserve the
normalized candidate-row semantics.

## Finalization Gate

Responses: 0002-transition-table-as-reviewable-data/artifacts/gate.md (Gate PASS 2026-08-24)

## References

- RDR 0001, Resolution Kernel Contract; `internal/resolve/resolve.go::Row`,
  `::Resolve`, `::assemble`.
- JDR 0001, `docs/jdr/0001-resolve-kernel-seam.md` — §D2 (gate-then-count),
  §D4 (kernel-enforced guard domain; normalizer emits existence constants),
  §D5 (reserved `<clear>`), §D6 (block-keyed routing; general match-block
  expansion), §D7 (closed layout: root/stop set, capability tables with
  `keys`, type-model keys, write-replaces, `[model.metadata]`), §JD-3
  (`RequiresOwned` producer), §JD-9 (`--tag` keys are observed), §JD-10
  (recognized-tag totality), §JD-15/§JD-16/§JD-17 (interface record).
- RDR 0007, Guard Predicate Totality — atom shape, `OpExists`/`LiteralTrue`/
  `LiteralFalse`, A21/A22, gate ordering.
- RDR 0008, Recognized Tag Key Ownership — reserved `recognized` name and
  `reserved_tag_key` category.
- RDR 0009, Escape-Row Shape Conformance — write-free escape rows.
- `docs/cli-output-contract.md`.
- Resource index: `.rdr/resources.md`.
- Seed prior: `../state-machines/BUILD-SEEDS.md`, especially Seed 2
  transition-table representation and Seed 3 guard expression.
- Transition model prior: "The transition model — inputs, outputs, error
  conditions."
- Kernel vocabulary prior: `../state-machines/ANALYSIS-kernel-vocabulary.md`,
  especially Deferred Choice, Exclusive Choice, and lint-rejecting
  document-order ambiguity.
- Direct `arc` corpus checks:
  - `StateMachineOS`: `sismic/sismic/io/datadict.py::import_from_dict` and
    `export_to_dict` for nested source models and contracts.
  - `StateMachineOS`: `transitions/transitions/core.py::Transition` for
    positive `conditions` and negative `unless` guard lists.
  - `StateMachineOS`: `stateless/src/Stateless/Graph/StateGraph.cs::StateGraph`
    for symbolic graph generation from machine metadata.
  - `StateMachineOS`: `sismic/sismic/io/plantuml.py::PlantUMLExporter` for
    rendered graph output as a view over model data.
  - `StateMachineLit`: statechart hierarchy / extended-state literature hits as
    the prior-art answer to dimensional state explosion.
- RDR Ragel POC contrast: `../state-machines/contrast/poc-rdr-ragel/REVIEW.md`,
  especially the coupled status register, cap-3 counter, guard, and readable
  current-state limitations.
- RDR and kata flow audits from the state-machine prior-art corpus.
- Tool-fit assessment for FSM libraries as validation/visualization tools, not
  runtime orchestrators.
