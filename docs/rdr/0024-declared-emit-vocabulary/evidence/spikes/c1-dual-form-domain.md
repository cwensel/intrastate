Model: claude-opus-5[1m]

# Spike — does C1's two-shaped `domain` decode, and what does strictness still catch?

`0024:C1` makes two `domain` forms normative for one emit key:

```toml
[emit.next]
kind = "enum"
domain = ["a", "b"]          # flat array

[emit.other]
kind = "enum"
[emit.other.domain]          # sub-table: disposition -> members
route = ["a"]
stop  = ["b"]
```

Two questions the Draft left open: (1) can one struct field decode both
shapes under the repo's decoder with `DisallowUnknownFields` set, and (2)
what does strict decoding still catch once it can?

Decoder: `github.com/pelletier/go-toml/v2 v2.2.4` (the repo's, per `go.mod`).

## Probe

The decl's `domain` typed as `any`, decoded with `DisallowUnknownFields()`
set — the same strictness `internal/table/source.go::decodeStrict` applies:

```go
type decl struct {
	Kind   string `toml:"kind"`
	Domain any    `toml:"domain"`
}
type doc struct {
	Emit map[string]decl `toml:"emit"`
}
```

## Q1 — both forms decode

```
flat        err=<nil>  domain=[]interface{}{"a", "b"}
subtable    err=<nil>  domain=map[string]interface{}{"route":[]interface{}{"a"},
                                                     "stop":[]interface{}{"b"}}
```

**Both decode.** C1's grammar is implementable as written. The Draft's Risk
entry worrying that the sub-table "proves awkward to … decode against real
TOML tooling" is answered: it is not.

The cost is typing. A concretely-typed field (`[]string`, or a
`map[string][]string`) accepts one shape and refuses the other, so the field
must be `any` or carry a custom unmarshaller. Every other field in
`sourceDoc` is concretely typed; this would be the sole exception, alongside
`sourceModel.Metadata`, which is `map[string]any` for the same reason and
documents it (`source.go:39-42`): "a free-form map so strict decoding
descends no further".

## Q2 — what strictness still catches, and what it stops catching

```
typo at decl level      (domaim = ["a"])                err=strict mode: fields
                                                            in the document are
                                                            missing in the
                                                            target struct
scalar under domain     (route = "not-an-array")        err=<nil>
deep junk under domain  ([emit.next.domain.route.oops]) err=<nil>
```

- A misspelled key at the **declaration** level still refuses. Strictness is
  intact down to that level.
- Inside `domain`, strictness catches **nothing**: a disposition whose value
  is a scalar rather than a member array, and arbitrary nesting below the
  disposition level, both decode without error.

## Verdict

C1's dual-form `domain` is decodable, and the strictness carve-out it
requires is bounded to the domain sub-tree — not to the declaration as a
whole. The consequence for the contract is that
`malformed_emit_declaration` must refuse by hand, inside that sub-tree,
what `DisallowUnknownFields` refuses free everywhere else in the schema:

- a `domain` that is neither a flat array of strings nor a table of
  disposition keys;
- a disposition whose value is not an array of strings;
- any nesting below the disposition level;
- a non-string member;
- a disposition sub-table carrying no keys.

Folded into C1's refusal list and into Testing Strategy scenario 1.
