Model: claude-opus-5[1m]

# A4 (a) — lint's held values are canonical renderings of the declared domain

## The chain

`decidableProduct` -> `productOver` -> `valueAssignments` -> `spreadValues`,
then each rendered held value is decided by `valueSatisfies` against the
same runtime `Evaluator`.

### internal/guard/product.go:334 — the render call

```go
		values, ok := valueAssignments(m.Tags[key])
		if !ok {
			return AssignmentSet{}
		}
		ranges = append(ranges, values)
```

The range of every value dimension is `valueAssignments(m.Tags[key])` — a
function of the DECLARATION, never of any caller-supplied string.

### internal/guard/assignment.go:294-310 — the int arm renders canonically

```go
	case "int":
		if d.Min == nil || d.Max == nil {
			return nil, false
		}
		width, ok := intWidth(*d.Min, *d.Max)
		...
		domain := make([]string, 0, width)
		for n, i := *d.Min, 0; i < width; n, i = n+1, i+1 {
			domain = append(domain, strconv.Itoa(n))
		}
		return spreadValues(domain, d.SingleValued), true
```

`strconv.Itoa(n)` is the canonical decimal spelling: no leading zero, no
`+`, no whitespace. `"07"`, `"+7"`, `" 7"` are NOT in the range of this
function for any `min`/`max`.

### internal/guard/product.go:541-555 — lint reuses the runtime seam

```go
func valueSatisfies(atom table.Atom, held string) resolve.GuardResult {
	seam := Evaluator{}
	literal := strings.Join(atom.Literal, "")
	if atom.Operator == "in" || atom.Operator == "contains" {
		literal = renderSet(atom.Literal)
	}
	return seam.Evaluate(resolve.GuardAtom{...}, held)
}
```

Its doc comment states the claim A4 rests on verbatim:

> It is the same decision the runtime evaluator makes, over the same
> rendered value form, which is what keeps the lint claim about the runtime
> it describes rather than about a second semantics

## Verdict on leg (a)

CONFIRMED for the HELD side. Lint never hands the seam a non-canonical
held value, so the held-side flip class is unreachable from lint.

## Caveat A4's own text anticipates — the LITERAL side

The literal is NOT rendered from the declaration: it is
`strings.Join(atom.Literal, "")` — the AUTHORED bytes. An authored
`n eq 07` reaches the seam as `"07"` on BOTH sides (lint and runtime).
Under today's raw-string `eq`, `held="7"` vs `literal="07"` is
GuardFalse for every rendered held value, so the atom denotes the EMPTY
set. Under typed comparison it denotes `{7}`. That is a real lint-verdict
flip, and it is the authored-literal leg (c) had to test empirically.
