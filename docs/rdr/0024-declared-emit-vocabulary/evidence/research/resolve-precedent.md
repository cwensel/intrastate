Model: claude-opus-5[1m]

# Stage 4 — precedent sweep: RDR corpus + codebase

Question put: has this decision already been made in a prior RDR or in the
codebase, and do project-level principles decide the three open forks?

## 1. The house pattern census — declare a vocabulary, prove authored values

Every model-authored vocabulary in the loader refuses unknown members.
Emit is the sole exception, and it is exceptional by explicit design.

| # | Vocabulary | Declaration | Unknown member | Where | Anchor |
|---|---|---|---|---|---|
| 1 | Tag keys `[tags.<key>]` | Mandatory | **REFUSED** `unknown_tag` (match, guard, write, clear, `[initial]`, accessor) | LOAD | `normalize.go:28,528,586`; `load.go:438,666` |
| 2 | Tag kinds | Mandatory | **REFUSED** `malformed_tag_declaration` | LOAD | `load.go:328`; `model.go:84` |
| 3 | Tag value domains | **Optional per kind** | Refused when declared; admitted when absent | LOAD + runtime | `load.go:849-870`, `793` |
| 4 | Operators | Closed | **REFUSED** | LOAD | `model.go:93` |
| 5 | Operator/kind matrix | Closed | **REFUSED** | LOAD | `model.go:110-125` |
| 6 | Outcome alphabet | Mandatory | **REFUSED** `malformed_outcome_binding` | LOAD | `load.go:139-162` |
| 7 | Dump columns `[dump].order` | **Optional** (absent = canonical default) | **REFUSED both ways** — unknown AND omitted | LOAD | `load.go:567-595` |
| 8 | Model class | Optional field | REFUSED if not in the two tokens | LOAD | `load.go:118-128` |
| 9 | Accessor bindings | Optional tables, arity mandatory | REFUSED | LOAD | `load.go:377-563` |
| 10 | Reserved keys | Kernel-fixed | **REFUSED** `reserved_tag_key` | LOAD | `load.go:234-258` |
| 11 | Schema fields | Closed | **REFUSED** `unknown_schema_field` | LOAD | `source.go:129` |
| 12 | **Emit keys `[rule.emit]`** | **NONE** | **ADMITTED unconditionally** | shape only | `normalize.go:611-632` |
| 13 | Graph invariants | Derived | Findings — **accumulated** | LINT | `graphlint/engine.go:59-81` |
| 14 | Escape-row shape | Kernel predicate | **Aggregated**, sorted, collapsed | RESOLVE | `escape_shape_0009_test.go:1186` |

Row 12 is the gap 0024 closes. The direction of the pattern is
one-directional and overwhelming.

## 2. RDR 0020 — the apparent dual, which is in fact affirmative precedent

0020 ("Undeclared `--tag` key admission — what the zero TagDecl means") asks
the same shape of question for tags. It does NOT conflict; it names 0024:

- `0020:§decision-rationale`(c): "The house already adjudicated this shape
  once: emit keys are 'undeclared, uninterpreted, byte-compared'
  (`0010:C3`), with opt-in declarations arriving later (0024) —
  **carrier-by-default, declare-to-tighten is the established pattern.**"
- `0020:C1`: "**Declaring a key later is the opt-in tightening**: admission
  then enforces that declaration's kind and domain."
- `0020:BR1`: "Per-model opt-in strictness (**the 0024 shape**) … compatible
  later work layered on top of the carrier default."

**The seam that reconciles them.** 0020 governs the CALLER's vocabulary at
admission (runtime, per-invocation, no author present). 0024 governs the
MODEL AUTHOR's vocabulary at load (design time, one author, one file).
0020's own safety argument is explicitly asymmetric and depends on the
model side staying closed:

- `0020:C1`: "**The model-side closed world is untouched**: a rule atom,
  accessor key, write target, clear target, or `[initial]` assignment naming
  an undeclared tag still refuses at load (`unknown_tag`)."
- `0020:§decision-rationale`(b): "Safety is structural, not disciplinary:
  the loader's `unknown_tag` refusals make an undeclared key unreadable by
  any rule … ⇒ the carried bytes cannot influence selection."

Emit values are model-authored, so they sit on the closed-world side.

## 3. Finding cardinality — the project rule, stated per tier

The governing comment, `internal/table/load.go:22-29`:

> "the refusal is singular: **load is fail-fast and returns one categorized
> error, never a list** (`0002:C24`, `0002:C3`). … Beyond those, **the order
> in which independent defects are checked is deliberately unspecified.**"

`internal/table/category.go:75-76`: "Load is fail-fast, so a document yields
**exactly one** of these and never a list."

Contract-level, `0008:C4`: "the predicate reports **the first breach it finds
and returns a single non-nil error**; it does not aggregate … **a producer
holding a programmer mistake is not owed an exhaustive list.**"

`0018:C1` settles coexisting breaches by PRECEDENCE, not accumulation:
"When an Input breaches both entry preconditions … Resolve MUST return the
escape-shape breach error … reordering the entry checks is a breaking change
to this contract." `0018:ALT1` rejected `errors.Join`: "aggregation serves
user-facing declarative validation; this is the programmer-mistake path."

The opposite tier, `0003:C19` (restated by `0006:C16`): "Lint MUST report
every defect it can decide in one pass over a row group, **not the first one
it encounters** … so the emitted set does not depend on row or dimension
iteration order."

**The rule that reconciles all three tiers: accumulate only where a TOTAL
ORDER over findings exists.**

| Tier | Behavior | Licensing condition |
|---|---|---|
| `load` | fail-fast, one error | defect order *deliberately unspecified* — no total order |
| `resolve` (0009) | aggregate, sorted, collapsed | `compareRefs` total order over `RowRef` |
| `graphlint` | accumulate all | `sortFindings` total order over findings |

Emit checks land in `load` by `0002:C24`'s arity split, so they inherit
fail-fast. Accumulating there would expose a non-deterministic SET, which is
exactly what the unspecified order forbids.

## 4. Payload append — JDR 0002 §D1 already decided that `dispositions` lands

`docs/jdr/0002-success-envelope-projection.md` §D1, settled 2026-08-28 at
the 0023/0024 Stage-2 joint check:

> "**Resolved:** lands in 0023 — `flow resolve`'s concrete ECHO/PLAN
> assignments … and in **0024 — `dispositions` is a PLAN-group field** (its
> C4 join reads plan-group inputs only), consistent by citation."

So the existence of the field is settled at decision level, not open. §D1
also supplies two principles bearing on its shape:

- "**Over-reporting is recoverable; silent omission is not**" — the tiebreak
  supporting `{}`-always-present over `omitempty`.
- "**Partitions are enforced by structure, not prose** — an unassigned field
  is a test failure, never a silent default."

The reflective partition oracle §D1 mandates does not exist in code yet
(0023 is Final but unimplemented), so whichever of 0023/0024 lands second
registers its field with the other's oracle.

`0010:C4` is the executed precedent for the mechanical move — it appended
`emit` after `Gates`, fixing the position in-contract "for the same reason
C3 fixes the dump column's … an unfixed position is an unlicensed diff
waiting on whichever implementer guesses differently."

## 5. Stale code comments the implementation must amend

Three comments state emit's openness as a positive design commitment and
read false after C1:

- `internal/table/model.go:159-166` — "It is deliberately NOT a TagValue …
  it is **undeclared, uninterpreted**, and compared by exact byte equality."
- `internal/table/normalize.go:617-619` — "they are **undeclared and
  uninterpreted**, so no case folding, trimming, or value coercion."
- `internal/table/dump.go:124-128` — "An emit key has neither: **it is
  undeclared, so there is no kind to consult.**"

The reconciliation: "undeclared" in `0010:C3` means *not a tag* — no
provenance, no accessor, no match/guard participation — every clause of
which C1 preserves. `[emit]` is a separate namespace from `[tags]`, so
`m.Tags[k]` stays absent and REQ-27's assertion keeps passing.

## 6. Principle statements — the conceptual-integrity anchors

- `internal/table/load.go:784` — "**Letting a caller write what a rule may
  not author would make a declaration advisory.**" (The most on-point line
  in the tree; it is the project's stated reason declarations must bind.)
- `docs/model-authoring.md:50-53` — "**Declaring the class rather than
  deriving it is what makes the refusal possible.** A tool that inferred …
  would silently accept a state machine whose owned declaration an author
  had dropped."
- `docs/model-authoring.md:64` — "**A misspelled facet is a stable refusal
  rather than a silent no-op.**" (Verbatim the defect 0024 exists to fix.)
- `internal/table/source.go:129` — "an unknown schema field is a stable
  refusal, never a silent no-op."
- `internal/table/load.go:585` — "a truncated dump is a refusal, never a
  silent narrowing."
- `0020:§decision-rationale`(e) — "an **interpretation of bytes the system
  pledges not to interpret**." (The genericity pledge.)
- `0006:C1` — "Coupling the resolver to a lint result would put a
  design-time proof on the runtime path, which this RDR's authority split
  rejects." (Independently forbids ALT3's shape.)
- `0021:C4` / `0020` — "structural, not disciplinary" — a recurring house
  value, and the argument for a declaration over a prefix convention.
