Model: claude-opus-5

# A2 — does a `[initial]` value reach the same canonical wire form as a `--write` value?

**Verdict: FALSIFIED for three of twelve admitted `[initial]` value kinds.**

The normalization premise holds for the *common* kinds — the rendering agrees
byte-for-byte in the artifact. It fails not because the two routes *render*
differently, but because the loader's `[initial]` admission set is STRICTLY
LARGER than the argv route's: three loader-admitted authorings have no argv
spelling `canonicalValue` accepts, and one of those three cannot be written by
either route at all.

## Command

```
go test ./internal/cli/ -run 'TestSpikeA2' -v
```

Raw output: `a2-canonical-form-run.txt`. Test source: `a2-spike-source.go.txt`
(run as `internal/cli/zz_spike_a2_test.go`, deleted after the run — nothing
under `internal/` was left modified).

## What was compared

Two comparisons per value kind, both end-to-end:

1. **value level** — `table.Load` a model carrying `[initial] <key> = <toml>`,
   take the resulting `table.TagValue`, render it through the only render the
   write path can accept (`resolve.Tag.Value` is one `string`, `TagValue.Value`
   is `[]string`: `canonicalSet(members)` for a `set` kind, `members[0]`
   otherwise), and compare against `canonicalValue(key, argv, decl, "write")`
   for the equivalent argv string.
2. **byte level (the stronger form the RDR asks for)** — drive `flow set-state`
   twice into two FRESH artifacts, once with the seed-derived wire string and
   once with the argv string, then diff the two artifact files byte-for-byte.
   Achieved: nine kinds are byte-identical artifacts, not merely value-equal.

## Structural finding that frames the result

`Model.Initial` has **no write-path consumer today**. The only non-test readers
of `Model.Initial` are `internal/graphlint/analysis.go` and
`internal/graphlint/reach.go`. There is no `Model.Initial` -> `resolve.Tag`
conversion anywhere in the tree — the seed surface does not exist yet. The
spike therefore states the candidate render explicitly
(`canonicalFromInitial`) rather than testing an existing one, because the type
gap is real: `table.TagValue{Key string; Value []string}` vs
`resolve.Tag{Key string; Value string}`.

## Admitted `[initial]` value kinds (from the loader, verified)

`valueMembers` (`internal/table/load.go:1748`) admits `string`, `bool`,
`int64`, `float64`, and `[]any` of those; every other TOML shape (inline
table, nested array) is `malformed_initial_declaration`. `loadInitial`
(`internal/table/load.go:1679`) then applies the reserved-sentinel refusal,
`conform(decl, "eq", members)`, and an ARITY check that requires exactly one
member unless `decl.Kind == "set"`.

| # | `[initial]` authoring | loader | argv route | wire forms |
|---|---|---|---|---|
| 1 | `status = "draft"` (enum) | admitted | admitted | AGREE `draft`, artifacts byte-identical |
| 2 | `note = "hello world"` (scalar) | admitted | admitted | AGREE, byte-identical |
| 3 | `flag = true` (bool) | admitted -> `"true"` | admitted | AGREE, byte-identical |
| 4 | `count = 7` (int) | admitted -> `"7"` | admitted | AGREE, byte-identical |
| 5 | `ratio = 1.5` (float on scalar) | admitted -> `"1.5"` | admitted | AGREE, byte-identical |
| 6 | `labels = ["b","a"]` (set) | admitted | admitted | AGREE `["a","b"]` (sorted), byte-identical |
| 7 | `labels = ["x&y","a<b","a<b"]` | admitted | admitted | AGREE `["a<b","x&y"]` — sorted, deduped, HTML-unescaped; byte-identical |
| 8 | `labels = []` (empty set) | admitted -> `[]string{}` | admitted | AGREE `[]`, byte-identical |
| 9 | `labels = ["plain"]` | admitted | admitted | AGREE `["plain"]`, byte-identical |
| **10** | **`labels = "plain"` (bare scalar on a set tag)** | **admitted -> `["plain"]`** | **REFUSED** | **DIVERGE** |
| **11** | **`note = ""` (empty string on a scalar tag)** | **admitted -> `[""]`** | **REFUSED** | **DIVERGE — and unwritable by EITHER route** |
| **12** | **`note = ["a"]` (1-elem array on a scalar tag)** | **admitted -> `["a"]`** | **REFUSED** | **DIVERGE** |

Refused by the loader, for completeness (these need no seed route): enum value
outside `domain`; float on an `int` tag; bare bool outside a set's `elements`;
multi-member array on a non-set kind (`kind scalar holds one value, not a
member sequence`); nested array literal; inline table; the reserved
`<clear>` sentinel (`reserved_tag_value`, which outranks the value arm).

## The three divergences, exactly

**#10 `labels = "plain"` on `kind = "set"`.** The loader's arity check is
`decl.Kind != "set" && len(members) != 1` — it constrains only non-set kinds,
so a bare scalar on a set tag passes and normalizes to the one-member set
`["plain"]`. The seed render produces `["plain"]` correctly. But the
*equivalent argv token* a caller would type is `plain`, and
`canonicalValue` (`internal/cli/flow_input.go:735`) refuses a non-array for a
set tag: `` the tag `labels` is set-valued and takes a JSON array literal; got
plain ``. The wire forms would agree IF the seed route rendered first and then
re-coerced; they diverge if the seed route hands the raw authored token to the
same coercion stage `--write` uses.

**#11 `note = ""` on `kind = "scalar"` — the sharpest one.** The loader admits
the empty string (no kind or domain rule rejects it; the arity check counts one
member). `canonicalValue` refuses it unconditionally for a non-set tag
(`internal/cli/flow_input.go:723`): `` the tag `note` was given an empty
value ``. In the byte-level test BOTH routes were refused — the seed-derived
wire string is itself `""`, so there is no wire form that reaches the artifact.
A model that loads clean therefore carries an `[initial]` assignment that is
**unwritable by any route that goes through `canonicalValue`**. If the seed
verb reuses the `--write` coercion stage, it will refuse at runtime a model the
loader already accepted, and the refusal will name a `--write`-shaped
complaint about a value the caller never typed on `--write`.

**#12 `note = ["a"]` on `kind = "scalar"`.** The loader admits it — arity 1
satisfies the non-set check — and flattens it to the bare member `a`. The seed
render is `a`, which writes fine. But the argv spelling of the same authoring
is `["a"]`, which `canonicalValue` refuses for a non-set tag
(`internal/cli/flow_input.go:719`): `` the tag `note` is not set-valued; got
the array literal ["a"] ``. Two authorings that the loader treats as identical
(`note = "a"` and `note = ["a"]`) have argv spellings only one of which is
accepted.

## What this means for the RDR

The claim as written — "read-back equality holds for seeded keys" — is TRUE for
every kind the two admission sets share, and the evidence is byte-identical
artifact files, not value-level agreement. It is FALSE as a general
normalization premise: the loader's `[initial]` admission set is a proper
superset of the argv route's, so the design cannot say "just route seeds
through the same coercion stage `--write` uses" without either

- narrowing the loader's `[initial]` admission to match (refuse #10, #11, #12
  at load time as `malformed_initial_declaration`, which makes the two sets
  coincide by construction and keeps the refusal at authoring time where the
  author can act on it); or
- routing the seed from the ALREADY-NORMALIZED `TagValue.Value []string`
  rather than from an argv-shaped token, and doing the domain re-check with
  `table.ConformValue` per member instead of through `canonicalValue`'s
  argv-shape arms — which handles #10 and #12 but still leaves #11 unwritable
  and needs an explicit decision on whether an empty scalar seed is legal.

Either way an explicit choice is required; the premise does not hold on its own.
