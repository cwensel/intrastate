Model: claude-opus-5[1m]

# A4 (c) — before/after lint diff under typed int comparison

## Method

Two binaries built from the same commit in a throwaway worktree:

- `intrastate-base`   — HEAD, raw-string `eq`/`in`
- `intrastate-typed`  — HEAD + a minimal scaffold patching the `eq`/`in`
  int arms to parse BOTH sides (`strconv.Atoi` on held value and on each
  literal member) and answer GuardUnevaluable when either fails. A
  package-level key->kind map stands in for the declared-kind carrier,
  installed at `graphlint.NewRequest` and at the `flow` CLI seam.

## Result 1 — the committed corpus does NOT flip

Corpus directory: `models/` (5 authored example models) plus every
lintable `*.toml` under `internal/*/testdata` — 123 models total.

```
for m in $(find models internal -name '*.toml' -not -path '*/rdr.toml' | sort); do
  ./bin/intrastate-base  lint --model "$m" --as json
  ./bin/intrastate-typed lint --model "$m" --as json
done
```

```
models linted: 123
NO DIFF across 123 models
```

`go test ./...` against the typed patch: fully green, zero regressions.

### Why the corpus is immune — and why that is not reassuring

`rg -n 'kind\s*=\s*"int"' models/` returns exactly ONE hit:
`models/examples/release-grammar.toml:72` (`[tags.risk]`, min 0, max 2).
Every guard over `risk` uses `gte` / `lt` / `gt` / `lte` — operators that
ALREADY parse both sides. The corpus contains ZERO `eq`/`in` atoms over
an `int` tag.

`rg -n 'eq\s*=\s*"?\s*[+-]?0[0-9]|in\s*=.*"0[0-9]' models/` -> no matches.

So the corpus does not exercise the changed code path at all. Its
no-flip result is a statement about corpus coverage, not about the
change's safety.

## Result 2 — an authored non-canonical int literal DOES flip lint

A4 predicts "no lint verdict flips" because lint's held values are
canonical. That holds for the HELD side (see a4-render-path.md). It does
NOT hold for the LITERAL side: the literal reaches the seam as the
AUTHORED bytes, and the loader's literal check is `strconv.Atoi`, which
ACCEPTS `"00"`, `"01"`, `"+1"`.

Under raw-string `eq`, an authored `n eq "00"` matches no canonical
rendering, so the atom denotes the EMPTY set. Under typed comparison it
denotes `{0}`. That changes the coverage union, and therefore the
verdict.

### Fixture model (`cover07`)

`n` declared `int`, `min = 0`, `max = 1`, single-valued, required;
two rows over outcome `step`: `r-zero` guarded `n eq <LITERAL>` and
`r-one` guarded `n eq 1`.

### Flip 1 — BLOCKING becomes CLEAN (literal `"00"`)

base (exit 2):
```json
{"code":"graph-lint-failed","message":"the model carries blocking graph-lint findings","schema_version":"0.1","findings":[{"code":"graph-coverage-gap","message":"group cover07/step over rows [r-one r-zero] leaves 1 of 2 assignments in its scoped product uncovered for the no_match arm; the coverage union must equal the scoped product","model":"cover07","severity":"blocking","rule":"r-one","element":"cover07/step","dimension":"n","class":"no_match"}]}
```
typed (exit 0):
```json
{"type":"ok","schema_version":"0.1","data":{"findings":[]}}
```

This is the worst direction: a model lint REJECTS today is ACCEPTED
after the change.

### Flip 2 — new blocking findings appear (literal `"01"`)

base (exit 2): one `graph-coverage-gap` (`no_match`, 1 of 2).
typed (exit 2): THREE findings — `graph-coverage-gap` (`ambiguous_match`,
2 of 2), `graph-coverage-gap` (`no_match`, 1 of 2), and a NEW
`graph-overlap`:

```json
{"code":"graph-overlap","message":"rows \"r-one\" and \"r-zero\" are enabled together by some assignment in cover07/step; lint rejects the ambiguity rather than selecting by source order","model":"cover07","severity":"blocking","rule":"r-one","span":"cover07:r-one","element":"r-zero","fingerprint":"n|all|eq|1,;#|n|all|eq|01,;#"}
```

Same exit code, different finding SET. Whether this counts as a "verdict
flip" depends on A4's granularity; the finding list is observably
different.

### Flip 3 — same as flip 2 for literal `"+1"`

Identical shape, fingerprint `n|all|eq|1,;#|n|all|eq|+1,;#`.

### Not a flip — literals the loader already rejects

`"07"` against `min=0,max=1` -> `malformed_predicate_atom: n.eq: 7 is
above max 1` at LOAD, both binaries. `" 1"` -> `malformed_predicate_atom:
" 1" is not an int` at load, both binaries. So leading/trailing
whitespace is caught by the loader; leading ZEROS and a leading `+` are
NOT.

## Verdict on leg (c)

A4 is FALSIFIED as written. Its held-side reasoning is sound and its
corpus claim is true, but its scope claim — "changes no lint verdict" —
fails on the authored-literal leg A4 itself flagged as needing testing.
