model: Claude Sonnet 5
variant: lite (profile: large)

# Repeatability-Lens Generation Run — RDR 0003 (Guard Predicate Exhaustiveness)

Scope note: RDR 0003 defines a grammar/semantics contract, not a shipped
package. `internal/resolve` exists (RDR 0001) and declares the
`GuardEvaluator` seam the RDR's evaluator must satisfy, but no implementation
of that seam exists under `internal/` yet, and the kernel-fixed atom shape
(JDR 0001 §D1 / RDR 0007) is specified but unshipped — the kernel still
carries `Row.Guard string`. Everything below is therefore inferred from the
RDR's normative clauses, not read from existing Go source. Every signature,
field, and helper name is a GUESS unless the RDR states the exact identifier.

---

## 1. Public API

The RDR is explicit that this package owns exactly one seam: **value
semantics over a present value**. It never sees the tag view, never decides
presence, and never decides `exists`. Diagnostics, parsing, and lint coverage
proof are named as responsibilities but the RDR does not fix their package
location — GUESS that lint proof lives in a sibling `lint` sub-package (RDR
0006 territory) and this package (`internal/guard`, GUESS name) owns
evaluation, well-formedness, and the atom/domain types shared with lint.

```go
// Package guard implements the closed, symbolic guard-predicate atom
// grammar and value-only runtime evaluation described in RDR 0003.
// GUESS: package path internal/guard.
package guard

// Operator is the closed, typed operator vocabulary. RDR 0003 fixes this
// set; no other operator is legal.
type Operator string

const (
	OpEq       Operator = "eq"
	OpIn       Operator = "in"
	OpLT       Operator = "lt"
	OpLTE      Operator = "lte"
	OpGT       Operator = "gt"
	OpGTE      Operator = "gte"
	OpExists   Operator = "exists" // decided by the kernel; never reaches this evaluator
	OpContains Operator = "contains"
)

// ValueKind is the declared tag value kind an operator must accept.
// GUESS: exact identifier and enumerators not named in the RDR.
type ValueKind int

const (
	KindEnum ValueKind = iota
	KindBoolean
	KindInteger
	KindStringLike
	KindSet
)

// Atom is this RDR's normalized guard atom. The RDR fixes the shape at the
// kernel seam (JDR 0001 §D1, normative in RDR 0007) as four fields: Key,
// Operator, Literal, Block. Source identity (RuleID, SourceLocator) is
// explicitly NOT an atom field — it lives on the enclosing Row.
type Atom struct {
	Key      string   // must resolve to a declared tag
	Operator Operator // must be allowed by the operator/kind matrix
	Literal  Literal  // must parse to the operator's literal shape
	Block    Block    // GUESS enum: BlockAll | BlockUnless
}

// Block distinguishes which conjunctive list (`all` vs `unless`) an atom
// came from. GUESS type/name — RDR only requires the field survive
// normalization (A14), not its Go representation.
type Block int

const (
	BlockAll Block = iota
	BlockUnless
)

// Literal is a parsed literal: one typed scalar (eq/lt/lte/gt/gte) or a
// non-empty typed scalar/element set (in/contains). GUESS concrete shape;
// RDR only fixes "typed scalar" / "typed scalar set" / "typed element set".
type Literal struct {
	Scalar any   // GUESS: could be a sum type instead of `any`
	Set    []any // non-empty when Operator is in|contains
}

// Domain is the declared finite domain for a tag dimension, supplied by
// RDR 0002's (not-yet-existing, A11) tag declaration. Enum/boolean values,
// a declared set element universe, or a bounded integer {min..max} range.
// GUESS shape — A11/A7/A9 are Pending producer requests; RDR 0002 does not
// yet declare a domain field of any kind.
type Domain struct {
	EnumValues []any // enum/boolean finite value set
	ElementSet []any // set-valued tag's declared element universe (A9)
	IntMin     *int  // bounded integer range; nil means unbounded
	IntMax     *int
}

// Evaluator evaluates a single atom's value semantics over one present
// value. It is the seam internal/resolve.go::GuardEvaluator names and
// currently has no implementation. It MUST NOT be given, and MUST NOT
// consult, the tag view — presence/absence and existence atoms are
// decided upstream by the kernel (JDR 0001 §D4).
type Evaluator interface {
	// EvaluateAtom decides a value atom against one already-present typed
	// value. GUESS exact signature; the RDR fixes only the contract
	// ("value semantics over a present value"), not the Go shape.
	EvaluateAtom(atom Atom, value TypedValue) (bool, error)
}

// TypedValue is a present tag value of a declared kind. GUESS name/shape.
type TypedValue struct {
	Kind  ValueKind
	Value any
}

// ParseAtom performs well-formedness checks: Key resolves to a declared
// tag, Operator is allowed by the operator/kind matrix for that tag's
// kind, and Literal parses to the operator's literal shape. Returns a
// predicate semantic kind error (unknown operator, type mismatch, literal
// parse failure, unsupported operator/tag-kind pairing) on failure — this
// RDR names the kinds; RDR 0006 maps them to lint finding codes and RDR
// 0005 maps those to the CLI envelope (RDR 0003 does not itself define a
// CLI-facing error type).
// GUESS exact signature and error type.
func ParseAtom(key string, op Operator, rawLiteral string, block Block, tagDecl TagDeclaration) (Atom, error) {
	panic("GUESS — RDR does not name this function")
}

// TagDeclaration is RDR 0002's tag declaration, consumed here only for
// value kind and (once A11 lands) declared domain. GUESS shape; owned by
// RDR 0002, not this RDR.
type TagDeclaration struct {
	Name       string
	Provenance Provenance // owned | observed | recognized
	Kind       ValueKind
	Domain     *Domain // absent today (A11 Pending) — nil until RDR 0002 adds it
}

type Provenance int

const (
	ProvenanceOwned Provenance = iota
	ProvenanceObserved
	ProvenanceRecognized
)

// --- Lint-facing surface (RDR 0006 consumes this; boundary GUESS) ---

// AtomDomain projects one atom onto the scoped product as the subset of
// assignments it accepts, per A2's derivation. Returns an error/refusal
// signal (not shown) when the participating dimension lacks a finite
// declared domain — see Verdict below.
// GUESS exact signature; the RDR states the math (A2), not a Go API.
func AtomDomain(atom Atom, domain Domain) (Subset, error) {
	panic("GUESS")
}

// Subset is a symbolic/bitset-equivalent representation of the accepted
// assignments within one scoped product. The RDR requires this be
// deterministic and non-enumerative for large products (no naive
// powerset materialization). GUESS concrete representation — the RDR is
// explicit that the representation is an implementation choice as long as
// it is deterministic and "symbolic or bitset-equivalent."
type Subset struct {
	// opaque; GUESS
}

// Verdict is the exhaustiveness/overlap outcome for one scoped row group.
// GUESS name/shape — RDR fixes only that the outcome is one of: exhaustive
// (green), coverage gap, overlap, or the single blocking
// "inability-to-prove" outcome (refuse/downgrade are one outcome, never
// two), which subsumes: non-finite dimension, product too large to prove,
// and a withheld claim under the runtime-veto narrowing.
type Verdict struct {
	Exhaustive       bool
	Gaps             []Finding
	Overlaps         []Finding
	InabilityToProve []Finding // the one blocking outcome; never split into hard/soft
}

// Finding names the contributing source rule/context id(s) for every
// diagnostic, per the RDR's normative "MUST name the source rule id or
// context id" clause. GUESS shape — RDR 0006 owns the actual finding code
// contract; this is this package's internal representation before RDR
// 0006 maps it.
type Finding struct {
	RuleIDs []string
	Atom    *Atom // present when the finding must name a specific refusing atom
	Reason  string
}
```

**Error modes** (all GUESS for exact Go error types; the RDR fixes the
*kinds*, not the representation):

- Parse/well-formedness failures, rejected **before resolution**: unknown
  operator, unknown tag, operator/tag-kind mismatch, literal parse failure.
- Runtime: this evaluator itself never produces `guard_unevaluable` — that
  refusal belongs to RDR 0007's kernel combinator (per-atom verdict
  combination happens in the kernel, not here). This package's
  `EvaluateAtom` is only called when the kernel has already established the
  key is present, so it should not itself need an "absent key" error path;
  GUESS that its error return is reserved for genuine internal invariant
  violations (e.g., a `TypedValue.Kind` that doesn't match `atom`'s expected
  kind, which "should never happen" if `ParseAtom` was honored).
- Lint-side: the single blocking "inability-to-prove" finding (never two
  categories, per the RDR's explicit "refuse and downgrade are one outcome"
  clause) for: non-finite dimension, too-large finite product, and withheld
  claim under the runtime-veto narrowing.

---

## 2. Three most important internal helper functions

1. **`buildScopedProduct(rows []Row, group GroupKey) (Product, error)`**
   (GUESS name/signature). Responsibility: given the rows sharing one
   selection context (RDR 0003 fixes the grouping key itself: same source
   state + same recognized outcome — "the row group is defined here"),
   assemble the multi-dimensional product of every participating guard
   dimension's declared finite domain. This is the load-bearing helper
   because the RDR is explicit that coverage/overlap **must not** be
   computed one dimension at a time — it must be one product. Fails (or
   flags a dimension) when any participating dimension lacks a finite
   declared domain.

2. **`rowAcceptedAssignments(row Row, product Product) (Subset, error)`**
   (GUESS name/signature). Responsibility: compute one row's accepted
   assignments as the intersection of its `all`-atom domains minus the
   single conjunctive assignment set matched by its full `unless` block.
   This directly implements the RDR's Load-Bearing Decision that `unless`
   is a conjunctive exclusion set, not per-atom negation, and that an
   escape row with no guard atoms denotes the whole product. Also must not
   special-case escape rows out of overlap checking.

3. **`canRefuse(row Row, tagDecls map[string]TagDeclaration) bool`** (GUESS
   name/signature). Responsibility: implement the "a row can refuse" test —
   true when the row carries a value atom whose key is not declared
   present-for-every-reachable-predecessor. The RDR requires this be
   decided **syntactically over declarations**, not from a runtime trace or
   witness input, so it must be a total, pure function over the declared
   model (blocked today on A12's predecessor-reachability contract from RDR
   0006, and A7's presence/optionality field from RDR 0002). This helper is
   what makes the runtime-veto narrowing enforceable: a row group cannot be
   certified exhaustive if `canRefuse` is true for any participating row.

(A fourth candidate, `projectExistsAtom`, implementing A7's presence
dimension, is explicitly **not yet specifiable** — the RDR marks A7
`Pending` and states the projection rule does not exist as a normative
clause yet. I list only three per the prompt, but flag that a shipped
implementation needs this one too once A7 lands — GUESS/inference, not
stated as done in the RDR.)

---

## 3. Data model (persisted / passed across the boundary)

Nothing in this RDR is persisted to disk — "This RDR introduces no
encode/decode pair. Parse/render fidelity ... belongs to RDR 0002." The
data model below is what crosses the **in-process** boundary between RDR
0002 (table/normalization), this RDR (grammar/evaluation), RDR 0001
(kernel/selection), and RDR 0006 (lint).

```go
// Row is RDR 0002's normalized candidate row, extended with this RDR's
// atom slice once RDR 0007's kernel reshape lands. Today the shipped
// kernel still carries Row.Guard as a raw string
// (internal/resolve/resolve.go::Row.Guard string) — this is the
// TARGET shape, not the shipped one. GUESS field names beyond RuleID/
// SourceLocator, which the RDR names explicitly.
type Row struct {
	RuleID       string // source identity — RDR 0002, cited not owned here
	SourceLocator SourceLocator // source identity — RDR 0002
	All          []Atom // positive conjunctive guard atoms
	Unless       []Atom // negative conjunctive exclusion atoms
	// GUESS: outcome/source-state fields used for row-group scoping,
	// name not given by the RDR ("same source state and same recognized
	// outcome").
	SourceState string
	Outcome     string
}

// SourceLocator identifies where in the TOML source a row/atom came from.
// GUESS shape; owned by RDR 0002, only referenced here.
type SourceLocator struct {
	File string
	Line int
	Col  int
}

// AtomIdentity is the Load-Bearing Decision's total identity tuple for a
// guard atom: (RuleID, SourceLocator, key, block, operator, literal).
// Never an index within all/unless. Used for diagnostics/dedup; NOT used
// for domain computation (that uses semantic equality: tag, operator,
// literal only).
type AtomIdentity struct {
	RuleID        string
	SourceLocator SourceLocator
	Key           string
	Block         Block
	Operator      Operator
	Literal       Literal
}
```

Wire/byte format: none owned here. RDR 0002 owns the TOML container; this
RDR only owns the guard-atom grammar embedded in it, illustrated (not
normative) as:

```toml
[rule.guard.all]
status.eq = "Draft"
profile.in = ["mid", "large", "foundational"]

[rule.guard.unless]
prelock_iterations.gte = 3
```

---

## 4. Top-level pseudo-code of the main operation

The RDR names two "main operations": (a) runtime evaluation feeding
exact-one selection, and (b) lint's exhaustiveness/overlap proof over a row
group. The lint path is the one this RDR most substantively owns (grammar +
finite-domain semantics), so pseudo-code below is for **exhaustiveness/
overlap proof over one scoped row group** — GUESS control flow; RDR states
the required checks and ordering constraints (must be one product, must not
be per-dimension) but not literal code.

```go
// ProveRowGroup implements RDR 0003's exhaustiveness/overlap contract for
// one scoped row group (rows sharing one source state + one recognized
// outcome — the set RDR 0001 resolves exact-one over).
// GUESS: entire function body is inferred from normative clauses, not
// read from source.
func ProveRowGroup(rows []Row, tagDecls map[string]TagDeclaration) Verdict {
	var v Verdict

	// 1. Collect every guard dimension (tag key) any row in the group
	//    references, across all/unless.
	dims := participatingDimensions(rows)

	// 2. Each dimension needs a finite declared domain (enum, boolean,
	//    set element universe, or bounded integer range) or the whole
	//    group's exhaustiveness claim is withheld for THAT dimension.
	//    "Refuse or downgrade" is ONE outcome, never a warning tier.
	for _, dim := range dims {
		decl, ok := tagDecls[dim]
		if !ok || !decl.HasFiniteDomain() { // GUESS method name
			v.InabilityToProve = append(v.InabilityToProve,
				Finding{Reason: "non-finite dimension: " + dim})
			return v // whole-group claim withheld; do not proceed to a false green
		}
	}

	// 3. Build the scoped product across ALL participating dimensions at
	//    once — never dimension-by-dimension.
	product, err := buildScopedProduct(rows, dims, tagDecls)
	if err != nil {
		// e.g. product too large for deterministic symbolic/bitset proof
		v.InabilityToProve = append(v.InabilityToProve,
			Finding{Reason: "product too large to prove deterministically"})
		return v
	}

	// 4. Project every atom onto the product; compute each row's
	//    accepted-assignments subset: intersection(all atoms' domains)
	//    minus (unless block's matched assignment set).
	accepted := make(map[string]Subset) // ruleID -> subset
	for _, row := range rows {
		sub, err := rowAcceptedAssignments(row, product)
		if err != nil {
			// e.g. an exists atom with no presence-dimension projection
			// yet (A7 Pending) — group cannot be resolved either way.
			v.InabilityToProve = append(v.InabilityToProve,
				Finding{RuleIDs: []string{row.RuleID}, Reason: err.Error()})
			return v
		}
		accepted[row.RuleID] = sub
	}

	// 5. Overlap: any pairwise non-empty intersection, using SEMANTIC
	//    equality for domain math but ALWAYS naming both rule ids —
	//    escape rows are NOT excluded from this check.
	for i := range rows {
		for j := i + 1; j < len(rows); j++ {
			if intersects(accepted[rows[i].RuleID], accepted[rows[j].RuleID]) {
				v.Overlaps = append(v.Overlaps, Finding{
					RuleIDs: []string{rows[i].RuleID, rows[j].RuleID},
				})
			}
		}
	}

	// 6. Coverage: union(accepted) == product. An escape row with no
	//    guard atoms denotes the whole product and closes coverage alone,
	//    but is still counted via the union, never as a side fact.
	if !unionEquals(accepted, product) {
		v.Gaps = append(v.Gaps, Finding{Reason: "coverage gap"})
	}

	// 7. Runtime-veto narrowing: even if steps 3–6 are clean, the claim
	//    MUST be withheld if any participating row "can refuse" —
	//    decided syntactically over declarations (A12/A7-dependent).
	for _, row := range rows {
		if canRefuse(row, tagDecls) {
			v.InabilityToProve = append(v.InabilityToProve, Finding{
				RuleIDs: []string{row.RuleID},
				Reason:  "row can refuse guard_unevaluable at runtime",
			})
			return v // green result must not be emitted alongside a withheld one
		}
	}

	// 8. Only if gaps/overlaps are empty AND no row can refuse: exhaustive.
	v.Exhaustive = len(v.Gaps) == 0 && len(v.Overlaps) == 0 &&
		len(v.InabilityToProve) == 0
	return v
}
```

Key ordering constraints the RDR makes normative (reflected above): finite-
domain check precedes product construction; the runtime-veto narrowing check
happens even after a clean coverage/overlap result (a group can be
domain-exhaustive and still have its claim withheld); "refuse" and
"downgrade" are never split into separate code paths; every emitted finding
names contributing source rule/context ids; emitting nothing is never a
valid substitute for the blocking finding.
