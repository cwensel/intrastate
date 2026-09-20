Model: claude-opus-5[1m]

# A4 (b) — MVV: `--tag n=07` against a guard `n eq 7`

## The model (`mvv-int`)

```toml
outcomes = ["step"]

[model]
id = "mvv-int"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.iter]
provenance = "observed"
kind = "int"
min = 0
max = 9
single_valued = true
required = true

[emit.next]
kind = "enum"
[emit.next.domain]
go = ["guarded", "fallback"]

[[rule]]
id = "guarded"
[rule.match.recognized]
eq = "step"
[rule.guard.all.iter]
eq = 7
[rule.emit]
next = "guarded"

[[rule]]
id = "fallback"
[rule.match.recognized]
eq = "step"
[rule.emit]
next = "fallback"
```

## The runs

```
intrastate flow resolve --model mvv-int.toml --outcome step \
    --tag iter=<V> --plan-only --as json
```

### `--tag iter=7` — unchanged (control)

base and typed, both exit 2:
```json
{"code":"flow-ambiguous-match","message":"more than one rule matches the recognized outcome `step`; resolve does not choose among them","schema_version":"0.1","findings":[{"code":"flow-ambiguous-match","message":"rule `fallback`: this rule matches and conflicts with the others named here","locator":"mvv-int:fallback","rule":"fallback"},{"code":"flow-ambiguous-match","message":"rule `guarded`: this rule matches and conflicts with the others named here","locator":"mvv-int:guarded","rule":"guarded"}]}
```

### `--tag iter=07` — FLIPS (GuardFalse -> GuardTrue)

base, exit 0 — the guarded row prunes, the fallback plans:
```json
{"type":"ok","schema_version":"0.1","data":{"revision":"","rule":"fallback","gates":[],"emit":{"next":"fallback"},"dispositions":{"next":"go"},"next":{},"writes":{},"clear":[],"escaped":false}}
```

typed, exit 2 — the guarded row now matches too:
```json
{"code":"flow-ambiguous-match","message":"more than one rule matches the recognized outcome `step`; resolve does not choose among them","schema_version":"0.1","findings":[{"code":"flow-ambiguous-match","message":"rule `fallback`: this rule matches and conflicts with the others named here","locator":"mvv-int:fallback","rule":"fallback"},{"code":"flow-ambiguous-match","message":"rule `guarded`: this rule matches and conflicts with the others named here","locator":"mvv-int:guarded","rule":"guarded"}]}
```

This is A4's predicted non-canonical-numeral flip, confirmed at the CLI.

### `--tag iter=+7` — FLIPS identically

base exit 0 plans `fallback`; typed exit 2 `flow-ambiguous-match`.

### `--tag iter=many` — the malformed leg NEVER REACHES THE SEAM

base and typed, IDENTICAL, exit 2:
```json
{"code":"flow-tag-invalid","message":"the value for `iter` does not conform to its declaration: \"many\" is not an int","schema_version":"0.1","param":"iter"}
```

`internal/cli/flow_input.go:751` calls `table.ConformValue` on every
DECLARED `--tag` before resolution, and `conformKind`
(`internal/table/load.go:1852`) runs `strconv.Atoi` for kind `int`. A
malformed value for a declared int tag is refused `flow-tag-invalid`
BEFORE any guard runs.

Routes checked for a malformed held value reaching the seam:

- `--tag` on a DECLARED int key -> refused `flow-tag-invalid` at input.
- `--tag` on an UNDECLARED key -> crosses verbatim as a pure carrier
  (`flow_input.go:695`), but then the kind map has no entry for it, and
  C2's own "key absent from the mapping" arm answers GuardUnevaluable
  regardless of the value.
- reader-sourced OBSERVED value -> `runReaders`
  (`internal/cli/flow_exec.go:331`) is provenance-filtered and folds only
  OWNED keys into the kernel input, so an observed reader value never
  reaches the guard; the atom is decided `absent`:
  ```json
  {"code":"flow-guard-unevaluable","message":"a guard predicate could not be decided over the assembled state","schema_version":"0.1","findings":[{"code":"flow-guard-unevaluable","message":"rule `guarded`: the atom on `iter` could not be decided (absent)","locator":"mvv-owned:guarded","rule":"guarded","key":"iter","operator":"eq","literal":"7","block":"all"}]}
  ```

So through today's CLI, the `GuardFalse -> GuardUnevaluable` leg A4 calls
"the fix" is NOT reachable via `--tag`. The seam-level arm is still worth
having as defense in depth, but the seed defect it is meant to cure is
not reproducible from the command line on a declared int tag.

## Direct evaluator probe — the full arm matrix

Ran in-package against `Evaluator{}.Evaluate` with the kind map installed
and cleared. Verdict encoding (`internal/resolve/resolve.go:82-89`):
`0 = GuardFalse`, `1 = GuardTrue`, `2 = GuardUnevaluable`.

```
key=n(int) eq  lit="7"        held="7"     BASE 1   TYPED 1
key=n(int) eq  lit="7"        held="07"    BASE 0   TYPED 1   <- flip F->T
key=n(int) eq  lit="7"        held="+7"    BASE 0   TYPED 1   <- flip F->T
key=n(int) eq  lit="7"        held=" 7"    BASE 0   TYPED 2   <- flip F->U
key=n(int) eq  lit="7"        held="many"  BASE 0   TYPED 2   <- flip F->U
key=n(int) eq  lit="07"       held="7"     BASE 0   TYPED 1   <- flip F->T
key=n(int) eq  lit="07"       held="07"    BASE 1   TYPED 1
key=n(int) in  lit=["7","8"]  held="07"    BASE 0   TYPED 1   <- flip F->T
key=n(int) in  lit=["07"]     held="7"     BASE 0   TYPED 1   <- flip F->T
key=n(int) in  lit=["7"]      held="many"  BASE 0   TYPED 2   <- flip F->U
key=s(scalar) eq lit="7"      held="07"    BASE 0   TYPED 0   (unchanged)
```

Both flip classes A4 names are real, and both are reachable from EITHER
side of the comparison, exactly as A4 states. Note the whitespace case
` 7`: A4's malformed arm covers it, but the loader already rejects a
whitespace LITERAL, so only the held side can carry it.
