Model: claude-opus-5[1m]

# Repeatability Diff — RDR 0002 (×3, full variant)

Neutral diff pass. Runs compared:

| Run | Model | Variant |
| --- | --- | --- |
| run-1 | `claude-opus-5[1m]` | full (profile: foundational) |
| run-2 | `Claude Sonnet 5` | full (profile: foundational) |
| run-3 | `claude-fable-5` | full (profile: foundational) |

run-1's header records `variant: full`, so ×3 is the intended coverage — no
variant mismatch to flag. Two model families are represented (opus / sonnet /
fable), so a split that isolates one run against the other two is noted below
where it occurs; note that with three distinct models every 2-1 split is
nominally "along a model boundary," so that signal is only weighted where the
odd run's rendering is also *semantically* different rather than a naming
choice.

Findings are ordered by how many runs disagree (3-way first). Each names the
RDR section the rewrite lands in.

---

## Health verdict

**Healthy.**

The disagreements localize to four contract seams — pipeline decomposition and
error ownership (R-1), the multi-model / `duplicate_model_id` scope (R-2),
validation step order (R-3), and the load entry's I/O boundary (R-4) — plus a
thin tail. GUESS markers cluster rather than scatter: all three runs
independently marked the same three things (Go identifiers, the concrete error
type, the dump grammar), and two or three converged on the same *contract*
gaps (`duplicate_model_id` scope, the missing-write-block category, the
accessor-entry shape). No run was confidently wrong in silence: every run
carried an explicit GUESS index or inline GUESS annotations, and all three
reproduced the RDR's hard normative clauses (two-pass version gate, full-atom
merge identity, match-blocks-only outcome lifting, suffix-only-on-expansion,
`RequiresOwned` derivation, no-alias writes/next, operator-based Match/Guard
split, identity-tuple sort with locator excluded) *identically*. That is the
signature of an RDR that is determinate where it speaks and silent in
locatable places — edit the RDR, do not rerun on another model.

The one caveat against a clean bill: the disagreement count is inflated by the
RDR's own internal tension at R-6, which is a genuine contradiction inside the
document rather than a silence.

---

## 1. Disagreements

### R-1 — Pipeline decomposition and which stage owns which error category (3-way)

**Contract in question**: how many operations the package exposes between "TOML
bytes" and "expanded table," and which of them may return which load-time
categories.

- **run-1**: one operation. `Load(path) (*Model, error)` performs decode,
  document validation, per-rule validation, expansion, and sort; the returned
  `*Model` then exposes `Rows()`, `Table()`, `Dump(w)` as pure accessors. All
  categories are `Load`'s.
- **run-2**: three operations. `LoadModel(path) (*Model, error)` (decode +
  validate only), `Normalize(m *Model) ([]resolve.Row, error)` (expansion),
  `Dump(rows, settings) (ExpandedTable, error)`. Categories split across
  `LoadModel` and `Normalize`, and `Normalize` is given its own error return.
  `Model` here is the **parsed source document**, a different type from
  run-1's `Model`.
- **run-3**: two-plus operations, and *flags the split as unresolved in the
  RDR*: `Load(src) (*Model, error)`, `Normalize(m) (*Expanded, error)`, then
  `(*Expanded).Rows()` and `Dump(e, w)`. Its inline note states outright that
  "whether Load internally calls Normalize to surface [the categories], or the
  two stages each own some categories behind one facade, is unspecified," and
  offers `Load returning (*Model, *Expanded, error)` as an equal alternative.

**RDR passage that let them diverge**: the Technical Design says "The parser
turns TOML into typed source data, then normalizes it into explicit candidate
rows," which names two stages — but every error clause is phrased as
"MUST be refused **at load**" (`#### Normative Contracts`, load-time category
block; `Load-time validation failures MUST retain stable data-level
categories`). Outcome binding, alphabet membership of the lifted literal, and
the escape-row write-free check are only decidable *during* expansion, yet are
all listed as load-time categories. The RDR never says whether "load" is one
callable or a pipeline name, so the runs placed the seam in three different
places. Note this is not naming: run-2's `Model` and run-1's `Model` are
different *values* (source document vs. normalized model), so a consumer
written against one reconstruction would not compile against the other.

**Lands in**: Technical Design → **Normative Contracts** (a clause fixing that
"load" is the whole source-to-candidate-rows pipeline and that every listed
category is refused by it, regardless of internal stage decomposition).

---

### R-2 — Scope of `duplicate_model_id`: who holds more than one model? (3-way)

**Contract in question**: `duplicate_model_id` is a load-time category, but the
load entry point in all three reconstructions takes a single document.

- **run-1**: an in-loader step. Pseudo-code carries
  `requireUniqueModelID(src.model.id)  # duplicate_model_id` inside the
  single-document `Load`, against an unstated registry — the check has no
  visible second model to compare against.
- **run-2**: relocated to *after* normalization as
  `checkNoDuplicateRowIdentities(rows)`, commented `// duplicate rule id /
  model id -> error`. This conflates two distinct RDR categories
  (`duplicate_rule_id`, `duplicate_model_id`) into one row-identity check.
- **run-3**: kept as its own step but marked GUESS on the spot:
  `check duplicate model id (multi-model dump input) -> duplicate model id
  # GUESS: where multi-model arrives`.

**RDR passage that let them diverge**: the dump-ordering contract says "A dump
spanning more than one model MUST carry a model-unique `model id` … two models
sharing an id is a `duplicate model id` load failure. RDR 0006 lints exactly
one model per invocation." The category is asserted, and a multi-model dump is
acknowledged, but no clause names the surface that admits more than one model —
neither the source schema (`[model]` is singular, one per document) nor any
load signature. The Testing Strategy compounds it by requiring "a duplicate
model id" mutated fixture without saying whether that is one file or two.

**Lands in**: Technical Design → **Normative Contracts** (the dump-ordering
clause), plus the load-time category list.

---

### R-3 — Validation step order: batched pre-pass vs. interleaved per-rule (3-way, with a 2-1 shape)

**Contract in question**: whether all validation runs as a pass over the whole
document before any expansion, or is interleaved with expansion rule by rule —
and therefore which category a document with two independent defects reports.

- **run-1**: interleaved. Document-level checks (alphabet, tags, contexts,
  model id) first, then a single `for each rule` loop that validates *and*
  expands in one pass. First defect encountered wins.
- **run-2**: batched. A complete `for rule in model.Rules` validation loop
  (`validateRuleID`, `validateWriteTargets`, `validateEscapeShape`,
  `validateOutcomeBinding`, `validatePredicateAtoms`) runs to completion,
  followed by `if any validation error: return firstOrAggregated(errors)`, and
  only then a *second* `for rule in model.Rules` loop that normalizes. This is
  the only run that contemplates **error aggregation** at all.
- **run-3**: interleaved, like run-1, but with a different intra-rule order:
  escape/write shape is checked *before* context flattening, whereas run-1
  flattens contexts and validates atoms before reaching the escape/write
  branch.

**RDR passage that let them diverge**: the RDR fixes the order of exactly one
pair — the version gate before strict decoding ("**The version check MUST run
before strict field validation**") — and is silent on every other ordering. The
Testing Strategy's "each has one mutated fixture that trips it **and no other**"
sidesteps the question by construction: with single-defect fixtures, order is
unobservable. Nothing states whether the loader is fail-fast or accumulating.
run-2's `firstOrAggregated` is a visible contract difference — an aggregating
loader and a fail-fast loader are not substitutable behind the same error
return.

**Lands in**: Technical Design → **Normative Contracts** (a fail-fast /
accumulate clause; and either a stated check order or an explicit statement
that order among independent defects is unspecified and untestable by design).

---

### R-4 — Load entry point: filesystem path vs. bytes (2-1, splits on the alt-model run)

**Contract in question**: whether the format layer touches the filesystem.

- **run-1**: `Load(path string)` plus a `LoadBytes(name, data)` variant.
  Marks as GUESS (G10) that file reading might live outside, but its
  pseudo-code opens with `bytes := readFile(path)`.
- **run-2**: `LoadModel(path string)` plus `LoadModelBytes(data, sourceID)`.
  Pseudo-code opens `raw := readFile(path)`. Same shape as run-1.
- **run-3** (fable, the odd run): `Load(src []byte)` **only**, with a reasoned
  GUESS: "`[]byte` rather than a path — config discovery / path binding is
  explicitly deferred to CLI integration (RDR 0005), so the format layer should
  not touch the filesystem. A LoadFile convenience wrapper is plausible but not
  required."

**Model-boundary note**: this is the one 2-1 split where the odd run reached a
*different architectural conclusion* from the same text rather than a different
name — run-3 read the RDR's deferral as prohibitive, runs 1 and 2 read it as
permissive. Worth flagging, though the underlying signal is the same silence
both readings exploit.

**RDR passage that let them diverge**: the Existing Infrastructure Audit says
"Config discovery | `internal/cli/config::Load` | … | **Reuse later** | Table
path binding belongs with CLI integration, not this RDR's data format." That
sentence says path *binding* is not this RDR's, which run-3 read as "the loader
takes bytes" and runs 1–2 read as "the loader takes a path, discovery is
elsewhere." The Capability Dependencies row for RDR 0005 ("CLI parse/lint
output") does not settle it either.

**Lands in**: Technical Design → **Normative Contracts**, or the Existing
Infrastructure Audit row for config discovery.

---

### R-5 — `Match`/`Guard` routing: is the outcome atom removed before or after the split? (2-1)

**Contract in question**: the sequencing of outcome lifting relative to the
operator-based `Match`/`Guard` routing, and whether the routing runs once per
rule or once per expanded row.

- **run-1**: routing runs **inside** the per-member expansion loop —
  `for each m in members: … match, guard := splitByOperator(predicate)` — so
  `splitByOperator` is called once per expanded row over the identical
  predicate set.
- **run-2**: routing runs inside the expansion loop too, but as a field
  initializer (`Match, Guard: splitByOperator(predicateAtoms)`), on the atom set
  `liftOutcome` returned.
- **run-3**: no `splitByOperator` step in the pseudo-code at all. The
  `CandidateRow` carries a single flat `Atoms []Atom` field, and the split to
  `resolve.Row.Match` / guard is deferred entirely to `Rows()`, described as
  "Mapping to the kernel row" — i.e. the normalized value is *unsplit* and the
  split happens at the kernel handoff.

**Why it matters beyond naming**: run-3's normalized row and runs 1–2's
normalized rows are different values. The RDR's Round-Trip invariant is stated
over "the full field list the dump contract carries … predicate atoms with
block," which matches run-3's unsplit `Atoms`; but the split clause says
"normalization MUST route each atom to exactly one of them," which reads as the
normalizer's output being already split. Both readings are defensible from the
text.

**RDR passage that let them diverge**: the split clause ("**The unified
predicate set splits across the kernel's two predicate fields by operator, and
this RDR owns the split**") says normalization owns the routing but never says
*which artifact* carries the split — the normalized candidate row, or only the
`resolve.Row` handed to the kernel. The dump contract lists "predicate atoms
(with block)" as one field, not two, which pulls toward run-3; the sentence
"The normalized candidate row **is** the kernel row" pulls toward runs 1–2. The
Prerequisites note that `Row.Guard` is still a `string` on `main` and the
atom-shaped reshape is unimplemented makes the target shape genuinely
underdetermined today.

**Lands in**: Technical Design → **Normative Contracts** (the operator-split
clause, and its relationship to the dump field list).

---

### R-6 — Outcome-binding extent: "combined predicate set" vs. "match blocks only" (0-way among runs; an RDR-internal contradiction)

**Contract in question**: which atom extent outcome binding scans.

All three runs rendered this **identically and correctly** — match blocks only
(local `match` plus inherited context `match`), with a `recognized` atom under
`guard.all`/`guard.unless` refused as `malformed_outcome_binding` and never
lifted. No disagreement between runs.

It is recorded here because the RDR contradicts itself on it, and the runs
agreed only because the Normative Contracts block is more emphatic than the
prose:

- **Normative Contracts**: "**Outcome binding reads the match blocks only** …
  Lifting must use the narrower extent."
- **Load-Bearing Decisions → Outcome binding**: "a rule's outcome is the single
  `recognized` atom in its **combined predicate set**, lifted into the row's
  outcome field" — and the same document defines "combined predicate set" as
  the *wider* extent (match + `all` + `unless`).

A reader who reaches Load-Bearing Decisions first gets the wrong extent. This
is a live contradiction, not a silence, and the runs' unanimity is luck of
reading order rather than evidence the document is clear.

**Lands in**: **Load-Bearing Decisions** (the Outcome binding bullet must be
restated to cite the narrow extent, or deleted in favor of the Normative
Contracts clause).

---

### R-7 — Category for an ordinary rule missing a write block (2-1, GUESS-flagged by one)

**Contract in question**: which category refuses an ordinary (non-escape) rule
with no `[rule.write]`.

- **run-1**: `fail(malformed_escape_declaration)` with an inline
  `# GUESS: category`, and lists it in the GUESS index as G8: "Category
  assigned to an ordinary rule missing a write block (RDR says a write block is
  required but does not name the category)."
- **run-2**: no check at all in the pseudo-code. `validateEscapeShape(rule)`
  is commented "escape XOR write/clear, never both" — it validates the escape
  direction only; a rule with neither escape nor write passes.
- **run-3**: `else: rule must have write block; writes/clears target declared
  owned tags only -> unknown tag / write to non-owned tag` — the check exists
  but the arrow assigns it the *target* categories, leaving the missing-block
  case with no category.

**RDR passage that let them diverge**: "An ordinary transition rule MUST contain
a write block" is stated in the rule-shape contract, but the load-time category
list has no entry for a violated write-block requirement.
`malformed_escape_declaration` is defined in the Testing Strategy purely by its
escape-side shapes ("an empty write block on an escape rule, and separately an
escape rule carrying a `clear` list"), which does not cover an ordinary rule
with no write block at all.

**Lands in**: Technical Design → **Normative Contracts** (the load-time
category list).

---

### R-8 — Next-state tags and writes: shared value vs. independently constructed (2-1)

**Contract in question**: the RDR's MUST NOT-alias clause, and whether the
reconstructions honor it.

- **run-1**: two independent `render(...)` calls, with the intent spelled out
  in a comment: `writes := render(rule.write, rule.clear)  # populated
  independently, NOT aliased to next`.
- **run-3**: same — `Writes: same rendering, populated independently (no
  alias)`.
- **run-2**: violates it. `writes := renderWrites(rule.Write, rule.Clear)`
  then `row.NextStateTags = writes; row.Writes = writes` — one map value
  assigned to both fields. Its prose *states* the contract correctly ("kept as
  two distinct fields per RDR 0009 A4") while its pseudo-code aliases them.

**RDR passage that let them diverge**: the clause is unambiguous
("Normalization MUST populate both explicitly from the rule and MUST NOT
populate one by aliasing the other"), so this is not an RDR silence — it is one
run failing to carry a stated MUST into its pseudo-code. Recorded as a
disagreement because the three renderings differ, but it is **not a candidate
rewrite**: the RDR already says the right thing. The only actionable note is
that Go semantics make the defect invisible to a value-equality test (both
fields compare equal either way), so the MVV cannot catch it — which argues for
a note in the Testing Strategy that this clause is enforced by review, not by
assertion.

**Lands in**: Validation → **Testing Strategy** (note that the no-alias clause
is unassertable by value comparison), not the contracts.

---

### R-9 — Where the derived row-kind column is computed (2-1, minor)

- **run-1** and **run-3**: kind is computed by the dump renderer from the
  escape-class list; the normalized row carries no kind.
- **run-2**: `ExpandedRow` carries a `Kind RowKind` field with the comment
  "derived: transition | escape," i.e. the derived value is **materialized on
  the expanded row type**, even though the same run's prose says row kind is "a
  *derived view property*, never a stored field."

**RDR passage**: "Row kind (`transition` / `escape`) is a **derived view
property, never a row field**: … The dump renders the kind as a column computed
from that list; normalization introduces no field to carry it." The clause
binds *normalization* and the *kernel row*; it does not explicitly say whether
an intermediate dump-model type may carry the computed column. run-2 exploited
that gap. Low severity — a one-clause clarification at most.

**Lands in**: Technical Design → **Normative Contracts** (escape/row-kind
clause).

---

## 2. GUESS clusters

Contracts two or more runs independently marked GUESS. These are the candidate
RDR rewrites, in descending cluster strength.

### C-1 — All Go identifiers: package path, function names, type names, receivers (3/3)

- run-1: G1, G2, G6 — "Package path/name `internal/table`; all Go identifiers,
  signatures, receivers"; guesses `Load`/`LoadBytes`/`Rows`/`Table`/`Dump` and
  `flattenContexts`/`bindOutcome`/`emitRow`/`sortRows`.
- run-2: guesses `internal/transitiontable`, `LoadModel`/`Normalize`/`Dump`,
  `readVersionOnly`/`liftOutcome`/`mergeContexts`.
- run-3: guesses `internal/transition`, `Load`/`Normalize`/`Rows`/`Dump`,
  `flattenContexts`/`liftOutcome`/`sortExpanded`, and says outright that
  "alternatives like `internal/table` or `internal/model` fit equally."

The three runs produced three different package names and three different
function-name sets, and all three cited the same source of silence — the
Existing Infrastructure Audit's "New internal package can own sparse source
structs and normalization," which names a responsibility and no identifier.

**Assessment: likely a leave-non-normative.** This is the largest cluster by
volume and the weakest by signal. An RDR of this class (a data/normalization
contract spec) legitimately leaves Go naming to implementation taste; the
divergence here costs nothing because no downstream RDR cites an identifier
from this package. Record it and move on — do **not** spend a Normative
Contracts clause minting names. The *structural* half of this cluster (how many
functions, and which owns which error) is a real gap, but it is already
captured as R-1 and should be fixed there, in terms of stages and category
ownership, not names.

### C-2 — The concrete error type and how categories surface in Go (3/3)

- run-1 (G4): "Single `*LoadError` type with `Category` + `Locator`, plus
  per-category sentinels," justified from the Testing Strategy's assert-on-
  category requirement.
- run-2: "GUESS: these surface as a Go `error` implementing
  `internal/cli/clierr.CLIError` … most plausibly via a `LoadError` type
  carrying `Category`, `Detail`, and a `SourceLocator`."
- run-3: a full `LoadError` struct with `Category`/`Model`/`Rule`/`Detail`/
  `Locator`, marked "GUESS: struct field set beyond Category."

All three landed on a `LoadError`-shaped thing and all three marked it GUESS.
Crucially they differ on **whether the load error is a `clierr.CLIError`**:
run-2 says yes (citing A5), run-1 and run-3 say the CLI mapping is explicitly
downstream design work (also citing A5). A5 supports both readings — it says
the envelope exists and that "their code names, `Group` assignments, and exit
mapping are new design work at implementation."

**Assessment: partly actionable.** The Go type is naming (leave it). The
*boundary* question — does the loader emit a CLI-shaped error, or a data-level
category that RDR 0005 maps? — is a real contract gap worth one sentence,
because it determines whether this package imports `internal/cli`. The RDR
already says "stable data-level categories **before** CLI mapping," which
favors the run-1/run-3 reading; making that clause say "the loader MUST NOT
depend on the CLI envelope" would close it.

### C-3 — Category-constant spelling (3/3)

- run-1 (G3): "snake_case spellings of the load-time category constants — the
  RDR only quotes `reserved_tag_key` and `unknown tag` / `unknown schema field`
  as prose."
- run-2: "GUESS: exact string spellings (e.g. whether it's
  `unknown_schema_field` or `unknown-schema-field`) — the RDR names categories
  in prose only."
- run-3: renders the whole set in prose form rather than as identifiers,
  implicitly declining to fix spellings, and uses `CategoryUnsupportedVersion`
  / `CategoryUnknownSchemaField` style in its API comments — a third convention.

**Assessment: naming, but with a wrinkle.** The RDR spells exactly one category
as an identifier (`reserved_tag_key`, inherited from RDR 0008) and the rest as
prose, so an implementation gets one snake_case anchor and no rule. Since the
Testing Strategy requires assertion *by category*, the categories are a real
API surface even though the strings are taste. Cheapest fix: state in the
category-list clause that the identifiers are snake_case matching the prose
names, anchored on `reserved_tag_key`. One line, and it stops three
implementations from disagreeing.

### C-4 — Dump output grammar (3/3)

- run-1 (G11): "Dump output syntax entirely — the RDR fixes the field list and
  row order but explicitly defines no grammar."
- run-2: "GUESS: concrete output format (text table vs JSON vs CSV) is left to
  implementation."
- run-3: "GUESS: io.Writer and the concrete text grammar — the RDR explicitly
  defines NO dump grammar."

**Assessment: not a gap — a deliberate scope exclusion.** All three runs cited
the *same* RDR sentence as the reason ("this RDR fixes the dump's field list
and row ordering but does not define a dump grammar — no delimiter, escaping,
quoting, or record separator"), and the Round-Trip section charts the
re-readable dump to a successor RDR with three named lossy sites. This cluster
is a false positive: the runs marked GUESS *because the RDR said so*. Leave
non-normative; the RDR is already explicit.

### C-5 — `duplicate_model_id` scope / multi-model input (2/3 explicit, 3/3 divergent)

run-3 marks it GUESS inline ("GUESS: where multi-model arrives"); run-1's G-index
does not list it but its pseudo-code check has no visible comparand; run-2
silently merged it into row-identity dedup. Combined with R-2's three-way
rendering split, this is the strongest *genuine* contract cluster in the set.

**Assessment: actionable.** Name the surface that admits more than one model,
or restate the category as scoped to whatever surface RDR 0006 lints over.

### C-6 — Contents of an `[accessors.<id>]` entry (2/3)

- run-2: `type AccessorRef struct { ID string } // GUESS shape; RDR only says
  "accessor references," execution is RDR 0004's`.
- run-3: "GUESS: contents of an accessor entry (the RDR never shows one)."
- run-1: carries `Accessors map[string]AccessorRef` with the whole struct set
  marked GUESS (G5), so it is inside its cluster without being called out.

**RDR silence**: the source-schema list names `[accessors.<id>]` as a root key
and the Capability Dependencies row says "This RDR may reference accessors but
does not execute them" (RDR 0004 owns execution) — but the RDR never shows an
accessor entry's body, and neither do the illustrative TOML nor the normative
field-layout clause. Two runs noticed the hole independently.

**Assessment: actionable but small.** Either show one accessor entry in the
Illustrative Code, or state that the entry body is RDR 0004's and this RDR
requires only that the id resolves (which is what `unknown_accessor` already
implies).

### C-7 — Whether the returned model retains declarations for RDR 0006 (1/3, noted for completeness)

run-1 alone raises it (G12): "That the alphabet, tag declarations, and accessor
map survive into the returned model as addressable state for RDR 0006 lint."
Below the two-run threshold, but it is the same seam as R-1 and Phase 5's "Expose
the parsed representation needed by RDR 0006" — which names no shape. Fold into
R-1's rewrite rather than tracking separately.

---

## 3. Agreement

All three runs rendered these identically, confirming the RDR is determinate
there: the two-pass version gate (permissive `[model]` read → refuse any
version but `1` → strict decode rejecting unmapped keys, in that order); strict
decoding as a format obligation rather than a library property; the full-atom
merge identity `(key, block, operator, literal)` with the explicit rejection of
any proper-prefix key, idempotent-on-identical, accumulate-never-override, and
contradictory atoms surviving as a dead rule for RDR 0006; outcome binding over
the match extent only with a guard-block `recognized` atom refused rather than
lifted; `in`-expansion to one row per member with the suffix attaching only when
the rule yields more than one row (single-member `in` ≡ `eq`, unsuffixed);
`RequiresOwned` as the sorted duplicate-free union of write-block and clear-list
keys, empty on escape rows, with the normalizer named as the sole enforcement
point; escape rules carrying neither write block nor clear list and row kind as
a derived property discriminated solely by a non-empty escape-class list;
existence atoms emitting `OpExists`/`LiteralTrue`/`LiteralFalse` verbatim with
any other existence literal refused; byte-exact key and literal identity with no
folding or trimming, and set literals normalizing to a member-sorted sequence
rather than a joined string; `[dump]` as presentation-only that never reaches
the normalized value; the total row order on the identity tuple `(model id,
rule id, expansion suffix)` compared byte-lexicographically with absent suffix
first and the source locator excluded, emitted from a pre-sorted sequence and
never from map iteration; the within-row sorts (atoms by `(key, block,
operator, literal)`; next tags, writes, required-owned, escape classes by key);
the arity-based load/lint split routing overlap, gap, dead row, and
read-before-write to RDR 0006 and `unmodeled_outcome`,
`owned_state_unavailable`, `guard_unevaluable`, and zero/multiple survivors to
the RDR 0001 kernel; the value-level (not textual) round-trip invariant with the
locator's line/column detail excluded but its presence compared; and the
seventeen-entry load-time category list itself, which all three reproduced with
the same membership.
