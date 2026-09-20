model: claude-sonnet-5
variant: full (profile: foundational)

Widened past the C1-C5/MVV/S1-S9 selections to read `§problem-statement`
and `§existing-infrastructure-audit` in full, to pin the surrounding
`GuardEvaluator`/`GuardAtom`/`GuardResult` shapes, the `struct{}`
zero-value evaluator, and `guard.DeclarationOf` — the contracts name
these by reference (REQ-10, JDR 0001 D1) without restating their fields.

## 1. Public API

```go
package guard // internal/guard

// NewEvaluator is the sole non-test construction path (C1). kinds maps
// tag key -> declared kind token, drawn from the five-kind vocabulary
// "enum" | "bool" | "int" | "set" | "scalar" (RDR 0003's spelling).
func NewEvaluator(kinds map[string]string) Evaluator

// DeclaredKinds is the single exported producer of that map (C1),
// because three packages (internal/cli, internal/guard,
// internal/graphlint) construct independently and an unexported helper
// cannot serve all three. TOTAL over m.Tags: every declared key gets an
// entry; a key whose Kind is "" is OMITTED (never mapped to "").
func DeclaredKinds(m *table.Model) map[string]string

// Evaluator implements resolve.GuardEvaluator (0007:C1). Unexported
// kinds field only; construction via composite literal / var / new /
// embedded zero value still compiles but yields a nil-mapping
// evaluator whose eq/in arms answer GuardUnevaluable (LOUD, C1) rather
// than reverting to raw-string comparison. lt/lte/gt/gte and contains
// are operator-inferred and unaffected by a nil mapping (S5 scope
// note) -- GUESS on exact type name "Evaluator" vs "evaluator" for the
// unexported struct; RDR only says "guard.Evaluator" as the package-
// qualified public type.
type Evaluator struct {
    kinds map[string]string // unexported; no resolve.TagSet, no runtime
                              // tag value, no view-typed state (C1)
}

// Evaluate is UNCHANGED by this RDR (0007:C1, REQ-10). Signature and
// atom shape (JDR 0001 D1: Key, Operator, Literal, Block) are frozen.
func (e Evaluator) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult

// resolve.GuardResult is three-valued: GuardTrue | GuardFalse |
// GuardUnevaluable (existing type, referenced not redefined here).

// Error / refusal modes surfaced through this change:
//   - flow-guard-unevaluable  (CLI/library refusal when Evaluate
//     returns GuardUnevaluable for an eq/in atom against an owned
//     value; reason "uncomparable", reuses 0011's vocabulary)
//   - malformed_predicate_atom (C5's canonical-int-spelling refusal at
//     load, internal/table/normalize.go::(*loader).atom; same code as
//     today's atom-build refusals, no new code minted)
//   - flow-tag-invalid (pre-existing, upstream of the seam entirely;
//     internal/cli/flow_input.go::canonicalValue; unchanged by this
//     RDR, fixture S8 pins it as a control)
```

Dispatch matrix for `eq`/`in` under a non-nil kinds map (C2, GUESS-free
-- this is stated near-verbatim in the contract):

| declared kind | comparison | unparseable held value | unparseable literal | key absent |
| --- | --- | --- | --- | --- |
| `int` | parsed integer equality (`strconv.Atoi`, platform-native width -- **GUESS**: no explicit width pin exists per C2 itself, so 32/64-bit divergence is an accepted open edge, not this RDR's to close) | GuardUnevaluable (in: any unparseable member poisons the whole list) | GuardUnevaluable | GuardUnevaluable |
| `bool` | token equality on `"true"`/`"false"` only | GuardUnevaluable | GuardUnevaluable | GuardUnevaluable |
| `enum`/`scalar` | exact string equality (unchanged) | n/a (every string parses) | n/a | GuardUnevaluable |
| `set` | GuardUnevaluable (defensive; eq/in not admitted by the operator/kind matrix) | -- | -- | GuardUnevaluable |
| unknown token (`default` arm) | GuardUnevaluable | -- | -- | -- |

## 2. Three most important internal helpers

1. **`guard.DeclaredKinds(m *table.Model) map[string]string`**
   (`internal/guard`, new per C1). Walks `m.Tags`, emitting one
   `key -> kind` entry per declared tag whose `Kind` is non-empty.
   Responsibility: be the SOLE translation from the model's declaration
   shape to the seam's plain-map carrier, so the three construction
   sites (CLI, guard/product.go, graphlint/reach.go) cannot hand-roll
   three divergent loops. Reuses the existing lookup surface
   `internal/guard/declaration.go::DeclarationOf` internally per the
   infra audit ("Reuse, behind the exported guard.DeclaredKinds(m)").

2. **The per-kind dispatch arm inside `Evaluator.Evaluate`** (private,
   unexported switch/dispatch — GUESS at exact internal shape; RDR
   describes it as "dispatch over the kind token is EXHAUSTIVE across
   the five-kind vocabulary" with a `default` meaning GuardUnevaluable,
   C2). Responsibility: own the eq/in kind-typed comparison rules in
   one place so lint (`conformKind`) and the runtime evaluator cannot
   silently diverge — this is the RDR's central correctness property
   (`D-canonical-spelling-at-load`).

3. **The canonicality check folded into
   `internal/table/normalize.go::(*loader).atom`**, beside its existing
   `conform(decl, operator, members)` call at `normalize.go:172` (C5).
   Responsibility: for `int`-declared predicate (guard/match) atom
   literals only, require `strconv.Itoa(strconv.Atoi(s)) == s` and
   refuse non-canonical spellings (`"00"`, `"+1"`, `"-0"`) under the
   existing `malformed_predicate_atom` code, with a diagnostic naming
   the site and the canonical rewrite. Deliberately does NOT touch
   `conformKind` itself, which stays lexical/permissive for
   `[emit]`/`[initial]`/`[rule.write]`/CLI values per `0024:REQ-7`/
   `REQ-9`.

(A fourth near-equal-importance helper the RDR names but I am not
promoting past 3: the typed check backstop for S4/S5, a `go/ast` or
`x/tools/go/packages` pass wired into `make check` that fails on any
non-`NewEvaluator` construction of `guard.Evaluator` in non-test code.
GUESS: exact package/function name for this checker is not given in
the RDR; I have not invented one above to avoid manufacturing a false
API surface.)

## 3. Data model (persisted / cross-boundary)

- **`map[string]string`** — the kind-declaration carrier itself
  (tag key -> kind token). This is the ONLY new persisted/passed shape;
  it crosses the `internal/table` -> `internal/resolve`/`internal/guard`
  boundary one-way (A5, verified import direction), so no
  `internal/table` type (e.g. `TagDecl`) crosses with it. Not
  serialized to disk — constructed fresh per request/evaluator
  lifetime from the loaded `*table.Model`.
- **`resolve.GuardAtom`** (JDR 0001 §D1, unchanged): `{ Key, Operator,
  Literal, Block }` — the atom shape the seam already receives; this
  RDR adds no field.
- **`resolve.GuardResult`** (existing, unchanged shape): three-valued
  enum `GuardTrue | GuardFalse | GuardUnevaluable`.
- **Envelope for the new refusal path**: `flow-guard-unevaluable`
  payload names the guarded row and the offending atom with
  `reason: "uncomparable"` (MVV step 2) — GUESS at exact JSON field
  names beyond "names the guarded row and the atom with reason
  uncomparable"; the RDR states the semantic content, not the literal
  envelope schema, so I widened to the existing envelope convention
  (`code`/`message`/`schema_version`/`param`, per fixture S8's
  `flow-tag-invalid` example) as the shape to reuse rather than
  inventing new fields.
- **C5's refusal fixture shape** (normative, S7): JSON envelope with
  `code: "malformed_predicate_atom"`, message
  `"<site>: \"00\" is not the canonical spelling of int tag n; write 0"`
  — no new envelope field, no sixth refusal kind.

## 4. Top-level pseudo-code of the main operation

```text
// Evaluator.Evaluate(atom, value) — the guard-seam comparison,
// C1 (carrier) + C2 (typed eq/in dispatch), ~constructed pseudocode
// from the normative prose; the RDR gives rules, not a body, so
// control flow/ordering below is a GUESS at sequencing consistent
// with every stated verdict.

func (e Evaluator) Evaluate(atom GuardAtom, value string) GuardResult {
    switch atom.Operator {
    case "lt", "lte", "gt", "gte":
        // operator-inferred int parse; kind lookup NOT consulted (C2)
        return evalOperatorInferredInt(atom, value)

    case "contains":
        // unchanged set-membership arms (§D13); kind lookup NOT consulted
        return evalContains(atom, value)

    case "exists":
        // never reaches this seam (0007:C1) — unreachable here
        panic("unreachable: exists resolved upstream")

    case "eq", "in":
        kind, declared := e.kinds[atom.Key]   // nil-map lookup is safe;
                                                // nil eval -> always !declared
        if !declared {
            // covers BOTH the true-absent key (defensive, A1 makes it
            // unreachable when C4 holds) AND a declared-but-empty-Kind
            // key (LIVE arm, C1 omits these from the map deliberately)
            return GuardUnevaluable
        }

        literalOK, parsedLiteral := parseUnderKind(kind, atom.Literal)
        if !literalOK {
            return GuardUnevaluable // defense in depth; load should
                                      // already have rejected this (0003:C8)
        }

        switch kind {
        case "int":
            n, err := strconv.Atoi(value)
            if err != nil {
                return GuardUnevaluable // present-but-malformed: never GuardFalse
            }
            if atom.Operator == "in" {
                // parse ALL members first; one bad member poisons the
                // whole list (parse-all-then-compare, not per-member)
                members, allOK := parseAllMembers(atom.Literal)
                if !allOK {
                    return GuardUnevaluable
                }
                return boolToVerdict(contains(members, n))
            }
            return boolToVerdict(n == parsedLiteral.(int))

        case "bool":
            if value != "true" && value != "false" {
                return GuardUnevaluable // canonicalValue upstream should
                                          // already prevent this on the
                                          // CLI door; owned door does not
            }
            return compareToken(atom, value) // eq: token equality;
                                                // in: membership

        case "enum", "scalar":
            return compareToken(atom, value) // exact string equality,
                                                // unchanged behavior

        case "set":
            return GuardUnevaluable // defensive; matrix should already
                                      // reject eq/in over set at load

        default:
            return GuardUnevaluable // exhaustive dispatch; unknown
                                      // token falls through here
        }
    }
}

// Construction obligation (C4), threaded not rebuilt:
//   CLI:       guardSeam(model) -> Evaluator, called once/request,
//              threaded down through probeRow (NOT rebuilt per row)
//   guard:     valueSatisfies(model, ...) constructs once per call
//   graphlint: matchSatisfiable(model, ...) constructs once, threads
//              down through atomAdmitsValue (NOT rebuilt per atom)
// Every non-test site: kinds := DeclaredKinds(sameModelRowsCameFrom)
//                       ev := NewEvaluator(kinds)
```
