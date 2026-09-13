Model: claude-sonnet-5

# A9 verification — opaque sentinel unforgeability and terminal-predicate node shape

---

## Limb (a): sentinel is unforgeable and reaches the wire as an ordinary value

### (a.1) `OpaqueValue` is a const equal to `"<opaque>"`

Confirmed.

`internal/graphlint/reach.go:28`
```go
const OpaqueValue = "<opaque>"
```

### (a.2) A tag with no finite declared domain yields `[]string{OpaqueValue}` into `Node.Values`

Confirmed. `internal/graphlint/reach.go:436-441`:

```go
func heldValues(m *table.Model, key string, value []string) []string {
	if _, finite := guard.AssignmentCount(m.Tags[key]); !finite {
		return []string{OpaqueValue}
	}
	return canonicalValues(value)
}
```

`heldValues` is the sole introduction point for held values, called from the two sites the comment at reach.go:433-435 names as "the two sites that INTRODUCE a held value": the root (`reach.go:120`, `root.Values[tv.Key] = heldValues(m, tv.Key, tv.Value)`) and a row's write (`reach.go:317`, `out.Values[w.Key] = heldValues(m, w.Key, w.Value)`). `Node.Values` is declared `map[string][]string` at `reach.go:22`. `guard.AssignmentCount` (`internal/guard/declaration.go:80`) is "the one authority on whether a declaration carries a finite domain" (reach.go:428-430) — a `scalar` or an unbounded `int` reports `ok=false`, so any such tag's value is discarded and replaced with the literal sentinel. Mechanism confirmed.

### (a.3) CRUX — can an AUTHORED tag value produce the literal string `<opaque>`?

**REFUTED.** No validation path anywhere in the loader rejects, escapes, or otherwise treats the literal string `<opaque>` specially. The only reserved *tag value* in the entire codebase is the clear sentinel `<clear>` (`table.ClearSentinel`), and it is the one value ever checked against.

`internal/table/model.go:11-14`:
```go
// ClearSentinel is the reserved tag value an explicit `clear` entry
// normalizes to. It is refused wherever a tag value is authored
// (`0002:C11`).
const ClearSentinel = "<clear>"
```

Every value-rejection site in the loader tests only against `ClearSentinel`:

- `internal/table/normalize.go:61-66` (predicate/atom values):
  ```go
  members, err := valueMembers(raw)
  if err != nil {
      return badAtom(err.Error())
  }
  if slices.Contains(members, ClearSentinel) {
      return Atom{}, fail(CatReservedTagValue, where+" authors the reserved value "+ClearSentinel)
  }
  ```
  (comment above, normalize.go:50-56: "The `<clear>` ban binds every atom regardless of operator or block (`0002:C17`)" — no mention of, or provision for, `<opaque>`.)

- `internal/table/load.go:1692-1695` ([initial] declaration values):
  ```go
  // The reserved-value rule takes precedence over the value arm: the
  // sentinel is refused wherever a tag value is authored (`0002:C11`).
  if slices.Contains(members, ClearSentinel) {
      return fail(CatReservedTagValue, "[initial] "+key+" authors the reserved value "+ClearSentinel)
  }
  ```

- `internal/table/normalize.go:619` (row write values): same pattern, `ClearSentinel` only.

A full-repo, non-test grep for `opaque`/`Opaque` (case-insensitive) turns up exactly one production definition (`graphlint.OpaqueValue`, reach.go:28) and its consumers inside `graphlint`; nothing in `internal/table` (the loader), `internal/guard`, or `internal/resolve` reserves, escapes, rejects, or special-cases the string `<opaque>`. `internal/table/category.go` enumerates exactly one `reserved_tag_value` category and it is anchored on `ClearSentinel`, not on a second sentinel.

Consequently: a `scalar` tag (or an unbounded `int`, i.e. exactly the class `guard.AssignmentCount` reports non-finite, i.e. exactly the class `heldValues` abstracts) can be authored with the literal value `"<opaque>"` — e.g. `write = { free = "<opaque>" }` in a row, or `[initial] free = "<opaque>"` — and it passes `valueMembers`, the `ClearSentinel` check, kind conformance (a scalar has no charset restriction — declaration.go:146, "`scalar` is the opaque scalar: it carries no finite domain"), and lands in `Node.Values["free"] = []string{"<opaque>"}` — a value string-identical to, and thus indistinguishable at the wire from, the sentinel `heldValues` manufactures for a genuinely unconstrained tag.

**Limb (a) is REFUTED.** The permitting path is the absence of any `<opaque>`-checking arm alongside the `ClearSentinel` arms at `internal/table/normalize.go:64-66`, `internal/table/normalize.go:619`, and `internal/table/load.go:1694-1695` — every one of which checks only `ClearSentinel`. A consumer of `Node.Values` cannot distinguish "tag abstracted because its domain is infinite" from "author literally wrote the string `<opaque>`."

---

## Limb (b): the two terminal-evaluation predicates and where each applies

### (b.1) `ownedAtomSatisfiable` — existential, opaque-admits reading

Confirmed. `internal/graphlint/reach.go:464-489`:

```go
// ownedAtomSatisfiable reports whether SOME value the node holds for the
// atom's key satisfies it. This is an existential test over a merged node,
// which is the safe direction: a merged node admits a superset of concrete
// views, so a witness real at a concrete view is still visible here.
func ownedAtomSatisfiable(n Node, a table.Atom) bool {
	held := n.Values[a.Key]
	if len(held) == 0 {
		// The key is absent in this node. An `exists = false` atom is
		// satisfied by exactly that; every other atom needs a held value.
		return a.Operator == resolve.OpExists && !existsWantsPresent(a)
	}
	if a.Operator == resolve.OpExists {
		return existsWantsPresent(a)
	}
	if slices.Contains(held, OpaqueValue) {
		// A tag with no finite declared domain abstracts to held/absent, so
		// any value atom over it is satisfiable wherever it is held.
		return true
	}
	for _, v := range held {
		if atomAdmitsValue(a, v) {
			return true
		}
	}
	return false
}
```

Confirmed as claimed: `slices.Contains(held, OpaqueValue)` short-circuits to `true` (satisfiable) at reach.go:478-481, an unconditional existential-with-opaque-admits reading. Called from `matchSatisfiable` (reach.go:457) and from `analysis.go:83` (`nodeSatisfiesMatch`, itself documented at analysis.go:76-77 as "an EXISTENTIAL test per atom").

### (b.2) `nodeMeetsAll` — universal, no-opaque reading

Confirmed. `internal/graphlint/analysis.go:423-444`:

```go
// nodeMeetsAll reports whether every value the node holds for each of the
// predicate's owned keys meets that predicate.
func (a *analysis) nodeMeetsAll(n Node, set []table.Atom) bool {
	for _, atom := range set {
		if a.model.Tags[atom.Key].Provenance != table.ProvenanceOwned {
			// A non-owned terminal key is invariant 1's dangling finding.
			// It cannot be met by an owned-state, so the predicate as a
			// whole is not satisfied here.
			return false
		}
		held := n.Values[atom.Key]
		if len(held) == 0 {
			return false
		}
		for _, v := range held {
			if !atomAdmitsValue(atom, v) {
				return false
			}
		}
	}
	return len(set) > 0
}
```

This is a universal quantifier over `held` (`for _, v := range held { if !atomAdmitsValue(...) { return false } }`) with **no** `OpaqueValue` short-circuit anywhere in the function — contrary reading confirmed, as claimed.

### (b.3) CRUX — is `nodeMeetsAll` confined to SPLIT (single-value-per-key) nodes?

**Confirmed — not refuted.** `nodeMeetsAll` has exactly one caller in the whole repo: `satisfiesSomeTerminal` (analysis.go:414-421, `if a.nodeMeetsAll(n, set) { return true }`). `satisfiesSomeTerminal` in turn has exactly two call sites, both inside a `splitNode` loop:

- `checkDeadEnd`, analysis.go:344-351:
  ```go
  for _, n := range a.nodes {
      for _, split := range splitNode(n, keys) {
          id := split.key()
          if seen[id] {
              continue
          }
          seen[id] = true
          if a.satisfiesSomeTerminal(split) {
  ```
- `checkTerminalEscape`, analysis.go:520-527:
  ```go
  for _, n := range a.nodes {
      for _, split := range splitNode(n, keys) {
          ...
          if a.satisfiesSomeTerminal(split) || a.hasOutgoingOrdinaryRow(split) {
  ```

In both, `keys := a.terminalKeys()` (analysis.go:342 and analysis.go:518) — the terminal-participating owned keys, computed by `terminalKeys()` (analysis.go:371-385) which walks `a.model.Terminal` and collects exactly the owned keys any terminal predicate atom names.

`splitNode` (analysis.go:391-409):
```go
func splitNode(n Node, keys []string) []Node {
	out := []Node{n.clone()}
	for _, key := range keys {
		values := n.Values[key]
		if len(values) <= 1 {
			continue
		}
		next := make([]Node, 0, len(out)*len(values))
		for _, base := range out {
			for _, v := range values {
				split := base.clone()
				split.Values[key] = []string{v}
				next = append(next, split)
			}
		}
		out = next
	}
	return out
}
```

This expands the *cross product* of value sets for exactly the keys in `keys` (the terminal-participating owned keys), producing one node per combination where each such key holds a singleton (`split.Values[key] = []string{v}`). Keys **not** in `keys` are left untouched and may remain multi-valued/merged on the resulting split node.

Because `nodeMeetsAll(n, set)` only ever reads `n.Values[atom.Key]` for `atom.Key` drawn from `set`, and every terminal predicate's atoms range only over keys that `terminalKeys()` collects from that same `a.model.Terminal` (terminalKeys walks precisely the sets `satisfiesSomeTerminal` iterates, analysis.go:373-374 `for _, set := range a.model.Terminal`), every key `nodeMeetsAll` actually evaluates on a `split` node is guaranteed singleton by construction of `splitNode`. Keys `nodeMeetsAll` never reads (because no terminal atom names them) may still be merged on `split`, but that is irrelevant to `nodeMeetsAll`'s universal loop since it never iterates them.

This is exactly what the code's own commentary states, analysis.go:450-456 (`hasOutgoingOrdinaryRow`'s doc, describing the same split node):
> "The node reaching here is ALREADY SPLIT on the terminal-participating keys (REQ-36), so on every key that decides terminal satisfaction it holds a single value and the test is exact there. On the remaining keys it stays a merged node read existentially..."

and analysis.go:309-319 (`checkDeadEnd`'s doc):
> "A terminal is a predicate over OWNED tags, and a node SATISFIES it when every value in each of the node's per-tag value sets meets it... The test runs on SPLIT nodes, not merged ones. Universal satisfaction is anti-monotone in merging... The split is exact: each split node is a set of concrete views the merged node already stood for. It is bounded by the declared domains of the terminal-participating keys ONLY, never the whole lattice."

**No path exists where `nodeMeetsAll` runs over a node multi-valued on any key it actually reads.** The two predicates therefore cannot disagree on the same relation: `ownedAtomSatisfiable`'s existential/opaque-admits reading is used for match-pattern reachability over genuinely merged nodes (`matchSatisfiable`, `nodeSatisfiesMatch`), while `nodeMeetsAll`'s universal/no-opaque reading is used only for terminal-predicate satisfaction over nodes pre-split to singleton exactly on the keys the predicate names. On a singleton value set, "some value satisfies" and "every value satisfies" coincide trivially — the two readings agree by construction, not by coincidence.

One residual caveat noted in the code itself (analysis.go:458-463, `hasOutgoingOrdinaryRow`'s doc) is a known, deliberately accepted false-negative in a *different* invariant (invariant 2's outgoing-row test, which reads merged non-terminal keys existentially) — this is an accepted imprecision, not a disagreement between the two terminal predicates, and does not touch `nodeMeetsAll`'s own correctness.

**Limb (b) holds.**

---

## Verdict

Limb (a): REFUTED. `<opaque>` is not a reserved tag value anywhere in `internal/table`'s loader/normalizer; only `ClearSentinel` (`<clear>`) is checked. An author can write `"<opaque>"` as an ordinary scalar/unbounded-int value and it reaches `Node.Values` indistinguishably from the synthesized sentinel.

Limb (b): holds. `nodeMeetsAll`'s only path to execution is via `satisfiesSomeTerminal`, called only from `checkDeadEnd` and `checkTerminalEscape`, both of which pass nodes already split (via `splitNode`) on exactly the terminal-participating keys — the only keys `nodeMeetsAll` reads. It never runs over a genuinely merged (multi-valued) key on the keys it evaluates.

Overall A9: BLOCKED by limb (a).
