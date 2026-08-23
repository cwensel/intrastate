model: Claude Sonnet 5
variant: lite (profile: large)

# Repeatability Lens — RDR 0003 Guard Predicate Exhaustiveness

Source: only `docs/rdr/0003-guard-predicate-exhaustiveness.md`. No evidence
directory, no prior run, no diff was read. Everything below is reconstructed
from the RDR text; every place the RDR is silent on an implementation detail
(package path, exact Go type names, function signatures, file layout) is
marked **GUESS**.

## 1. Public API

The RDR describes two seams normatively (grammar/semantics and declaration
model) but does not specify a package path or Go signatures anywhere — it
only names the kernel seam it plugs into (`internal/resolve::GuardEvaluator`,
owned by RDR 0001/0007) and says this RDR's evaluator satisfies the
value-only half of it. All concrete signatures below are **GUESS**, shaped
to match the normative constraints (four-field atom shape from JDR 0001 §D1;
evaluator "MUST NOT read the tag view"; declaration model's five fields).

```go
// GUESS: package path. RDR only says "no guard-predicate package under
// internal/" exists yet; it does not name where this RDR's code should live.
package guard // GUESS — could equally be internal/guard or internal/predicate

// --- Declaration model (this RDR owns; RDR 0002 owns authoring/carriage) ---

// Kind is one of the five value kinds. Spelled exactly these five tokens
// per Normative Contracts.
type Kind string

const (
	KindEnum   Kind = "enum"
	KindBool   Kind = "bool"
	KindInt    Kind = "int"
	KindSet    Kind = "set"
	KindScalar Kind = "scalar"
)

// TagDecl is the tag declaration model: value kind plus, as the kind
// admits, a finite domain, optionality, single-valuedness, and (set only)
// an element universe. GUESS on field names/types; the five fields
// themselves are normative ("Declaration model" Load-Bearing Decision).
type TagDecl struct {
	Kind          Kind
	Domain        []Literal // enum value set, or int {min,max} — GUESS shape
	Min, Max      *int64    // int bound, inclusive both ends — GUESS representation
	Elements      []Literal // set element universe (KindSet only)
	Optional      bool      // default true (no marker = optional, conservative default)
	SingleValued  bool      // only valid on enum/bool/int; MUST NOT be set on set/scalar
}

// Validate enforces the domain/kind agreement clause (rejects e.g. a
// {min..max} bound on an enum, an element universe on non-set, or
// single-valued on set/scalar) before normalization completes.
// GUESS signature; RDR only requires "rejected before normalization
// completes" by "the declaration loader."
func (d TagDecl) Validate() error // GUESS

// --- Guard atom grammar (this RDR owns; atom shape itself fixed by
// JDR 0001 §D1 / RDR 0007, cited not restated here) ---

type Operator string

const (
	OpEq       Operator = "eq"
	OpIn       Operator = "in"
	OpLt       Operator = "lt"
	OpLte      Operator = "lte"
	OpGt       Operator = "gt"
	OpGte      Operator = "gte"
	OpExists   Operator = "exists"
	OpContains Operator = "contains"
)

type Block string

const (
	BlockAll    Block = "all"
	BlockUnless Block = "unless"
)

// Atom mirrors the kernel-fixed shape: Key, Operator, Literal, Block.
// Source identity (RuleID, SourceLocator) is NOT an atom field — it is
// carried by the enclosing normalized row (RDR 0002 / resolve.go::Row).
type Atom struct {
	Key      string
	Operator Operator
	Literal  Literal // typed scalar or typed scalar/element set — GUESS
	Block    Block
}

// Literal is a typed value or typed set of values. GUESS representation;
// RDR fixes only that set literals canonicalize to an unordered,
// duplicate-free typed set before entering the identity tuple.
type Literal struct {
	Kind  Kind
	Value any   // GUESS
	Set   []any // GUESS, used when Operator is in/contains
}

// AtomIdentity is the total identity tuple used for diagnostics/dedup:
// (RuleID, SourceLocator, Key, Block, Operator, Literal). Semantic
// equality for domain computation is the narrower (tag, operator, literal).
type AtomIdentity struct {
	RuleID       string // GUESS type — RDR names RuleID/SourceLocator as
	SourceLocator string // existing fields on resolve.go::Row; types not given
	Key          string
	Block        Block
	Operator     Operator
	Literal      Literal
}

// --- Well-formedness / parse (this RDR owns) ---

// ParseAtom validates Key resolves to a declared tag, Operator is allowed
// by the operator/kind matrix for that tag's kind, and Literal parses to
// the operator's literal shape and lies inside the tag's declared domain.
// GUESS signature. Error modes named by the RDR (not as Go error types,
// but as "predicate semantic kinds"):
//   - unknown operator
//   - type/operator-kind mismatch (unsupported operator/tag-kind pairing)
//   - literal parse failure
//   - literal-outside-declared-domain (distinct from parse failure)
//   - declaration/kind disagreement (declaration-time, not atom-time)
func ParseAtom(key string, op Operator, rawLiteral string, block Block, decls map[string]TagDecl) (Atom, error) // GUESS

// --- Runtime evaluator (this RDR owns; kernel decides presence/existence
// and never hands this evaluator an absent-key case per JDR 0001 §D4) ---

// Evaluator decides value semantics over a PRESENT value only. It MUST NOT
// read the tag view (normative). This is the concrete type intended to
// satisfy internal/resolve's GuardEvaluator seam (name/signature there is
// RDR 0001/0007's, not restated here — GUESS shown for shape only).
type Evaluator struct{}

// EvaluateAtom evaluates one atom's operator against one already-present
// value. Two-valued (bool), never three-valued — three-valued combination
// (Kleene, ¬U = U for unless) is the kernel's job, not this evaluator's.
func (Evaluator) EvaluateAtom(atom Atom, value Literal) (bool, error) // GUESS

// --- Static lint proof (consumed by RDR 0006; grammar/semantics owned here) ---

// RowGroup is the set of normalized candidate rows sharing one selection
// context (same source state, same recognized outcome) — defined
// normatively in this RDR, not RDR 0006.
type RowGroup struct {
	Rows []Row // GUESS — Row here is RDR 0002/kernel's normalized row type
}

// ScopedProduct is the cartesian product of every participating guard
// dimension's assignment count (per the per-kind assignment-count table).
type ScopedProduct struct {
	Cardinality uint64 // GUESS type; must handle up to 2^N for set kinds
	Dimensions  []Dimension // GUESS
}

// ProveExhaustiveness computes coverage (union of accepted assignments ==
// scoped product) and overlap (pairwise intersection, split into the
// ordinary-row population and the escape-row population) for one RowGroup.
// Returns one Verdict, never stops at the first defect (normative: "MUST
// report every defect it can decide in one pass").
// GUESS signature — RDR states the proof obligation, not a Go API for it;
// this may in fact live inside RDR 0006's lint package rather than here.
func ProveExhaustiveness(g RowGroup, decls map[string]TagDecl, bound uint64) Verdict // GUESS

// Verdict carries: Covered bool, Gaps []Gap, Overlaps []Overlap,
// Withheld []WithholdReason (unprovable dimension / too-large product /
// runtime-veto narrowing), and — per the bare-escape-row observability
// clause (A19, still open) — which bare escape row (if any) closed
// coverage. GUESS shape entirely; the RDR fixes only the semantic content
// each finding MUST carry (source rule/context ids, witness assignment for
// gaps, refusing atom for withheld claims — the atom-naming half is A20,
// open, no producer on either surface yet).
type Verdict struct {
	Exhaustive bool
	Gaps       []Gap
	Overlaps   []Overlap
	Withheld   []WithholdReason
}
```

**Error modes (named normatively, not as Go types):**
- unknown operator
- unknown tag
- operator/tag-kind mismatch
- literal parse failure
- literal-outside-declared-domain
- declaration/kind disagreement (rejected before normalization; includes
  `{min..max}` on `enum`, element universe on non-`set`, `single_valued` on
  `set`/`scalar`)
- (lint-only) non-exhaustive finite domain / `graph-coverage-gap`
- (lint-only) overlapping candidate rows (two populations: ordinary vs.
  escape-for-one-failure-class)
- (lint-only) blocking inability-to-prove: non-finite dimension, too-large
  finite product, or withheld claim under the runtime-veto narrowing —
  **all three are the same single blocking outcome**, never a soft/advisory
  tier
- (runtime, not this RDR's) `guard_unevaluable` — RDR 0007's per-row/per-atom
  payload, only cited here

All of the above route through `internal/cli/clierr::CLIError` /
`internal/cli/respond::Fail` per A4 (GUESS on exact call sites — RDR only
names the existing gateway, not how this package invokes it).

## 2. Three most important internal helpers

1. **Declaration validator / domain-kind agreement check.** Rejects a
   `TagDecl` whose fields disagree with its `Kind` (e.g., `set` carrying
   `single_valued`, `enum` carrying an element universe, `scalar` carrying
   any finite domain) before normalization completes. This is the gate that
   makes every downstream exhaustiveness claim trustworthy — A2's whole
   derivation depends on domains being declared, not invented. **GUESS**
   name (`TagDecl.Validate` above); normative behavior is explicit.

2. **Scoped-product builder / participation resolver.** For one `RowGroup`,
   determines which keys participate (any row in *that* group carries a
   guard atom over the key, in `all` or `unless`, regardless of whether it
   discriminates), excludes match keys (the grouping context, never a
   dimension), and computes each participating dimension's assignment count
   per the per-kind table (`|domain|` for single-valued enum/int, `2` for
   single-valued bool, `2^|domain|` for unmarked finite kinds, `2^|universe|`
   for `set`, and ×2 per optional key's presence dimension). This is the
   function whose correctness the critique iteration explicitly flagged
   (cardinality previously understated for `set` and unmarked kinds).
   **GUESS** name/signature; behavior is fully normative.

3. **Coverage/overlap prover over the product.** Projects each atom onto
   the product (`all` atoms as intersection of allowed subsets; `unless`
   block as one conjunctive excluded intersection, subtracted only when
   fully decided — never per-atom negation, never subtracted when any
   participating `unless` atom is unevaluable, per the two-valued-
   subtraction-is-a-lint-projection-of-Kleene clause), unions accepted
   assignments across rows for coverage, and intersects pairwise for
   overlap in two separate populations (ordinary rows; escape rows against
   each other for one failure class only). Must also apply the runtime-veto
   narrowing (withhold the claim if any participating row can refuse
   `guard_unevaluable` — decided syntactically from the optional-declared
   field alone, in either block) and the bare-escape-row observability rule
   (a bare escape row that closes coverage must be named in the verdict,
   not silently certified green). **GUESS** decomposition into a single
   function versus several; the RDR states these as one connected
   obligation ("Lint MUST report every defect... in one pass").

## 3. Data model (persisted / crossing the boundary)

Two things cross the RDR 0002 ↔ RDR 0003 ↔ RDR 0006/0007 boundary:

**a) Tag declaration** (authored in RDR 0002's TOML, meaning owned here):

```toml
[tags.<name>]
provenance = "owned" | "observed" | "recognized"   # RDR 0002's, cited here
kind       = "enum" | "bool" | "int" | "set" | "scalar"
domain     = [...]              # enum: required to claim exhaustiveness
min = 0                         # int: both bounds inclusive
max = 3
elements   = [...]              # set: element universe
single_valued = true|false      # enum/bool/int only; illustrative spelling —
                                 # authoring location still open (A16)
optional   = true|false         # default true (no marker = optional)
```

**b) Guard atom**, the kernel-fixed four-field shape (JDR 0001 §D1),
cited not restated: `Key` (string, resolves to a declared tag), `Operator`
(one of the eight closed tokens), `Literal` (typed scalar or canonicalized
unordered duplicate-free typed set), `Block` (`all` | `unless`). Authored
under `[rule.guard.all]` / `[rule.guard.unless]` in RDR 0002's TOML. Source
identity (`RuleID`, `SourceLocator`) is **not** an atom field — it lives on
the enclosing row (shipped today as `internal/resolve/resolve.go::Row`,
fields `RuleID`, `SourceLocator`, and currently `Guard string` pending RDR
0007's reshape to an atom slice — **not this RDR's field to change**).

Identity tuple for diagnostics/dedup: `(RuleID, SourceLocator, key, block,
operator, literal)`. Semantic equality for domain math only:
`(tag, operator, literal)`.

Nothing here is persisted at runtime beyond the normalized model RDR 0002
already carries — this RDR introduces no encode/decode pair and no new
storage (explicit in "Round-Trip / Inverse Invariants").

## 4. Top-level pseudo-code of the main operation

The "main operation" is the static exhaustiveness/overlap proof over one
row group — the operation this RDR's Normative Contracts and desk trace
walk through step by step (GUESS on decomposition into functions; the
control flow itself is normative).

```
function proveRowGroup(group, declarations, bound):
    # 1. Well-formedness (should already be done at parse time, restated
    #    here since lint depends on it)
    for atom in group.allAtoms():
        assert atom.key resolves to a declaration
        assert atom.operator allowed for that declaration's kind
        assert atom.literal parses to operator's literal shape
        assert atom.literal within declared domain          # GUESS ordering

    # 2. Determine participating dimensions (guard keys only; match keys
    #    are the grouping context, excluded)
    dims = {}
    for atom in group.allAtoms():
        if atom.block in {all, unless}:               # every guard atom
            dims[atom.key] = true                      # participates even
                                                         # if only one row
                                                         # constrains it
    if atom.key optional:
        dims[atom.key].addPresenceDimension()          # {present, absent}

    # 3. Check every participating dimension has a finite declared domain
    unprovable = []
    for key in dims:
        decl = declarations[key]
        if decl has no finite domain:
            unprovable.append(InabilityToProve(dimension=key))

    # 4. Build scoped product cardinality (per-kind assignment-count table)
    cardinality = product(assignmentCount(decl) for key, decl in dims)
    if cardinality > bound:
        unprovable.append(InabilityToProve(product=cardinality, bound=bound))

    if unprovable is non-empty:
        emit each unprovable finding                   # loud, one per cause
        return Verdict(exhaustive=false, reasons=unprovable)  # but continue
                                                         # to steps 5-7 too —
                                                         # withholding MUST NOT
                                                         # suppress other findings

    # 5. Project each atom onto the product; compute each row's accepted
    #    assignment set
    accepted = {}
    for row in group.rows:
        allSet = intersect(atom.acceptedSubset() for atom in row.allAtoms)
        if row.unlessAtoms all decided:
            unlessSet = intersect(atom.acceptedSubset() for atom in row.unlessAtoms)
            accepted[row] = allSet - unlessSet
        else:
            # any unless atom over an optional key is potentially
            # unevaluable -> ¬U = U -> whole row unevaluable
            mark row as "can refuse"

    # 6. Coverage: union of ordinary-row accepted sets vs. scoped product
    #    (bare escape row with no guard atoms = whole product, closes
    #    coverage alone but must be named as such, not silently green)
    covered = union(accepted[row] for row in ordinaryRows + escapeRows)
    gap = scopedProduct - covered
    if gap non-empty:
        emit gap finding(context, all rule ids in group, one witness assignment)

    # 7. Overlap, two separate populations
    for (r1, r2) in pairs(ordinaryRows):
        if accepted[r1] ∩ accepted[r2] non-empty:
            emit overlap finding(r1.ruleID, r2.ruleID)
    for (e1, e2) in pairs(escapeRowsForSameFailureClass):
        if accepted[e1] ∩ accepted[e2] non-empty:
            emit overlap finding(e1.ruleID, e2.ruleID)
    # escape-vs-guarded overlap is NOT checked — not a runtime ambiguity

    # 8. Runtime-veto narrowing: withhold claim if any participating row
    #    "can refuse" (declaration-only test, either block)
    for row in group.rows:
        if row "can refuse" (per step 5 marking):
            emit withheld-claim finding(row, refusing atom)  # atom-naming
                                                               # has no
                                                               # producer yet
                                                               # (A20)

    # 9. Verdict is green ONLY if: no unprovable dims, no gap, no overlap,
    #    no withheld-claim reason. All findings from steps 3-8 are emitted
    #    together, not short-circuited at the first one.
    return Verdict(...)
```

## Where guesses were heaviest

Almost entirely in Section 1 (concrete Go API) — the RDR fixes the *grammar,
semantics, and normative behavior* exhaustively but never states a package
path, exported type/function names, or signatures for anything except by
reference to code that belongs to peer RDRs (`internal/resolve::Row`,
`::GuardEvaluator`, `internal/cli/clierr::CLIError`). Every struct/function
name and shape in Section 1 and the pseudo-code's decomposition into
discrete functions is a best-guess fit to the normative constraints, not a
reconstruction of an actual API the RDR specifies.
