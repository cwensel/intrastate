model: claude-opus-5[1m]
variant: full (profile: foundational)

# Repeatability diff — RDR 0009

Diffs `run-1.md` (claude-opus-5[1m]), `run-2.md` (claude-fable-5), and
`run-3.md` (glm-5.2:cloud) against
`docs/rdr/0009-escape-row-shape-conformance-ownership.md`. This session
authored none of the runs. Cross-model draw satisfied: three distinct base
models.

Findings are ordered by how many runs disagree. Each names the RDR section
the rewrite lands in.

## 1. Disagreements

### D-1 — `Resolve`'s signature and where the checked table comes from (2-vs-1, splits on the model boundary)

**Contract**: the entry precondition's call site — `Resolve`'s first
statement, per Load-Bearing Decisions *Placement* and the Technical Design
("called by `resolve.go::Resolve` at entry before `assemble`").

| Run | Rendering |
| --- | --- |
| run-1 (opus-5) | `func Resolve(t Table, in Input) (Result, error)` — table as a **separate first parameter**; step 0 is `t.CheckValid()` |
| run-2 (fable-5) | `Resolve(in) -> (Result, error)`; step 0 is `in.Table.CheckValid()`, explicitly marked **GUESS: field spelling `in.Table`** |
| run-3 (glm-5.2) | `func Resolve(in Input) (Result, error)`; step 0 is `in.Table.CheckValid()`; `Input` field set marked **GUESS** (ledger #1) |

**Ground truth**: `internal/resolve/resolve.go::Resolve` is
`func Resolve(in Input) (Result, error)`, and `Input` carries
`Table Table` (`resolve.go::Input`). Runs 2 and 3 guessed correctly; run-1
reconstructed a **signature change to RDR 0001's locked entry point** — a
second parameter that does not exist.

**The RDR passage that let them diverge**: the RDR never writes `Resolve`'s
signature or names `Input` anywhere. It says "called by
`resolve.go::Resolve` at entry", "the supplied table", "over every row of
the supplied table at `Resolve` entry", and "`Resolve`'s first statement,
above the existing `view := assemble(in)`". `assemble(in)` is the only hint
that the table is reached through `in`, and it is incidental — quoted as a
placement landmark, not as a contract. Two of three runs marked the
`in.Table` reach a GUESS even while getting it right.

This is the diff's sharpest finding because the error mode is not
cosmetic: the RDR states as a Positive consequence that "the five-kind
taxonomy, RDR 0005's mapping, RDR 0002's grammar, and the kernel's type
vocabulary all stand unchanged", and Alternative 1 is rejected partly for
"reworking the implemented kernel". An implementer reproducing run-1's
reading would change `Resolve`'s arity — breaking every one of the sixteen
`escapeRow` call sites and every frozen test — while believing they were
honoring a no-signature-change RDR. The RDR's own Normative Contracts
never forbid it, because they never state the signature.

**Lands in**: Technical Design → Normative Contracts (the kernel-precondition
clause), and Load-Bearing Decisions → *Placement*.

### D-2 — Where the zero-value `RowRef{"",""}` sorts in a mixed report (1 explicit, 2 silent)

**Contract**: the multi-breach report's order — "ordered by `RowRef`
identity using the kernel's existing `compareRefs` ordering".

| Run | Rendering |
| --- | --- |
| run-1 | Flags it explicitly (**G4**): "sort position of the zero-value `RowRef{"",""}` relative to populated identities … falls out of `compareRefs`'s string comparison, but no clause or scenario states it" |
| run-2 | Silent — renders `sort refs by compareRefs` with no mixed-identity case |
| run-3 | Silent — same; its worked example is single-identity |

**Ground truth**: `resolve.go::compareRefs` is
`strings.Compare(a.RuleID, b.RuleID)` then `strings.Compare(a.SourceLocator,
b.SourceLocator)`, so `RowRef{"",""}` sorts **first**. The behavior is
determined; the *RDR* is what is silent.

**Assessment — admissible but non-blocking.** Run-1 is right that no clause
or scenario pins the mixed case: scenario 10 is all-distinct-identity,
scenario 10b is all-degenerate (`RowRef{"",""}` ×3). Neither exercises a
table carrying both a degenerate row and an identified one. But the order
is fully derived from `compareRefs`, which the RDR names as the ordering
authority — this is a **test-coverage gap, not a contract silence**: an
implementer following the RDR cannot get it wrong, they merely have nothing
asserting they got it right. Recorded as such rather than as a normative
rewrite.

**Lands in**: Validation → Testing Strategy (scenario 10/10b coverage note).

### D-3 — `CheckValid`'s receiver form (1 explicit, 2 silent)

| Run | Rendering |
| --- | --- |
| run-2 | **GUESS: value receiver `(t Table)` rather than pointer** — "the RDR pins 'a METHOD ON Table taking no arguments' but not the receiver form" |
| run-1 | `func (t Table) CheckValid() error` — value receiver, unmarked |
| run-3 | `func (t Table) CheckValid() error` — value receiver, unmarked |

**Assessment — not admissible as a contract silence.** All three converged
on the value receiver; only one flagged the reasoning. `Table` is a read-only
structural scan on a stateless kernel, and the package's cited habit
(`TagSet.Lookup`, `Result.Refused` — `Refused` is a value receiver at
`resolve.go::Result.Refused`) settles it. Zero-run disagreement in the
rendered API. Noted for completeness; no rewrite.

## 2. GUESS clusters (two or more runs marked GUESS)

### C-1 — The A9 envelope carrier field: name, clierr-local type, whether `Count` serializes (**3/3 GUESS**)

All three runs marked this, and all three correctly identified it as
*already known open*:

- run-1 **G2**: "openly deferred to A9, which is **Pending** — this is a
  known-open contract, not an unnoticed silence";
- run-2: "the RDR explicitly leaves both (plus whether `Count`
  serializes) to A9";
- run-3 ledger **#6**: "a genuinely open contract the RDR itself marks
  Pending and owes to Stage 6".

**Disposition: no rewrite owed by this lens.** The RDR's own clause says
"The clause binds the carrier CLASS only; the field's name, its
clierr-local type … and whether the per-identity `Count` serializes are
A9's to settle." A9 is `Pending` with a written Plan and is already on
Prerequisites. Unanimous GUESS on a *deliberately* deferred contract is the
lens confirming the deferral is legible, not finding a gap. Carried to the
needs-verification list for Stage 6, where it already sat.

### C-2 — Diagnostic prose: sentinel message and `EscapeShapeBreachError.Error()` format (**3/3 GUESS**)

run-1 **G1**, run-2 (both bullets), run-3 ledger **#3/#4**. All three
invented plausible strings and all three noted the RDR forbids depending on
them.

**Disposition: no rewrite owed.** The Normative Contracts require identity
be inspectable structurally and "MUST NOT be recoverable only from formatted
prose"; nothing asserts on text. Unpinned prose is the intended state —
pinning it would manufacture a contract the RDR deliberately declines. run-1
notes the operator does read it, which is true and is what the `Hint` and
`Code` exist for.

### C-3 — Collapse mechanism: map-then-sort vs sort-then-dedupe (**3/3 GUESS**)

run-1 **G6** ("output is pinned, so this is genuinely free"), run-2 ("the
accumulation structure"), run-3 ledger **#5** ("RDR fixes the outcome …
but not the algorithm").

**Disposition: no rewrite owed** — and this is the healthy signal, not a
finding. All three runs derived *identical observable behavior* (collapse
equal identities, per-identity `Count` = pre-collapse rows, `compareRefs`
order, `errors.Join`, one-level `Unwrap`) from different internal
mechanisms. The RDR pinning outcome over algorithm is working exactly as
intended. Per the lens's own admissibility rule, helper-decomposition taste
is not a finding.

### C-4 — Kernel type field sets the RDR does not restate (**2/3 GUESS**)

run-3 ledger #2/#10/#11/#12 (`Escape`'s element type, the five
`RefusalKind` names, whether `Table` carries more than `Rows`, `Plan`'s full
field set); run-1 **G5** (`Escape`'s element type); run-2 silent.

**Ground truth**: `Escape []RefusalKind` (`resolve.go::Row`) — run-3
guessed `[]EscapeEdge`, run-1 guessed "plausibly `[]RefusalKind` or a
dedicated failure-class type". `Table` does carry more (`Revision`,
`Outcomes`), and `Plan` carries `Revision` beyond the named five.

**Disposition: no rewrite owed.** These are RDR 0001's implemented types,
which this RDR explicitly leaves unchanged ("`Row`, `Table`, `Plan`,
`Result`, `RowRef` … keep their existing shapes"). An RDR is not obliged to
restate a predecessor's locked vocabulary; an implementer reads it from
source, where it is unambiguous. Only emptiness of `Escape` is load-bearing
here, and all three runs got that right.

### C-5 — File location of the new symbols (**1/3 GUESS**)

run-3 ledger #9 (`resolve.go` vs a new `breach.go`); runs 1 and 2 silent.

**Disposition: not admissible** — file layout is not a contract.

## 3. Agreement

The three runs rendered these **identically**, which confirms the RDR is
determinate on its load-bearing surface: the four pinned exported spellings
(`ErrEscapeShapeBreach`, `*EscapeShapeBreachError{Ref RowRef; Count int}`,
`Table.CheckValid() error`, `Code: "escape-row-shape-breach"`); the
per-row predicate `len(row.Escape) != 0 && len(row.Writes) != 0` as
length-based and `Writes`-only (never `NextTags`, never nil-sensitive);
whole-table scope including dormant rows; precedence over all five modeled
dispositions; the always-aggregate `errors.Join` return with a
one-level `Unwrap() []error` of `*EscapeShapeBreachError` elements and
per-element `errors.Is`; verbatim (never `fmt.Errorf`-wrapped) return from
`Resolve`; identity collapse with per-identity pre-collapse `Count`;
`compareRefs` order with no positional tiebreak; the zero-`Result` /
`Refused() == false` caller trap; `GroupInternal`/exit 2 with a `Hint` and a
non-`Detail` serialized identity carrier; and that no type split, row-kind
field, or sixth refusal kind is introduced.

Run-3 states it directly: "The RDR is otherwise unusually determinate on the
load-bearing contracts … The silences above are mostly naming/prose/mechanism,
not contract gaps."

## Verdict

One admissible contract silence (**D-1**), one test-coverage gap
(**D-2**), and one already-open deferred contract (**C-1**, A9, no new
work). Everything else is prose, mechanism, or predecessor vocabulary —
the healthy pattern the lens's *Expected signal* describes.
