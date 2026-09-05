Model: claude-opus-5

# SPIKE — RDR 0020 undeclared `--tag` key admission, HEAD baseline

Purpose: capture the CURRENT (pre-change) behaviour of `flow resolve --tag`
for keys the loaded model does not declare. Nothing was implemented; no
source file was edited.

## Setup

- Repo: `intrastate`, branch `main`, HEAD `b142fd9` (worktree carries one
  unrelated modified RDR doc; no Go source is dirty).
- Build: `go build -o /tmp/intrastate-spike ./cmd/intrastate` (exit 0).
- Fixture: `carrier-table.toml` (this directory) — a `decision-table` so no
  reader / `--artifact` is needed. It declares the enum scalar `tier`
  (guarded by both rows, so a resolve actually selects), the unguarded set
  `labels`, and the recognized outcome `decide`. It declares NOTHING named
  `extra` or `extras`. `intrastate lint --model carrier-table.toml` exits 0
  with no findings.

Code under test: `internal/cli/flow_input.go` — the guarded lookup
`decl := m.Tags[key]` (yields the zero `TagDecl` for an undeclared key)
feeding `canonicalValue`, whose `!isSet` arm refuses an array literal and,
separately, refuses an empty value.

## Commands (copy-pasteable, run from this directory)

```sh
go build -o /tmp/intrastate-spike ./cmd/intrastate            # from repo root

# (a) undeclared scalar + undeclared array carrier
/tmp/intrastate-spike flow resolve --model carrier-table.toml --outcome decide \
  --tag tier=free --tag 'labels=["security"]' \
  --tag extra=plain --tag 'extras=["a","b"]' --as=json

# (a2) fallback per protocol: undeclared SCALAR carrier only
/tmp/intrastate-spike flow resolve --model carrier-table.toml --outcome decide \
  --tag tier=free --tag 'labels=["security"]' --tag extra=plain --as=json

# (b) same resolve, no carrier flags
/tmp/intrastate-spike flow resolve --model carrier-table.toml --outcome decide \
  --tag tier=free --tag 'labels=["security"]' --as=json

# (c) DECLARED scalar handed an array literal
/tmp/intrastate-spike flow resolve --model carrier-table.toml --outcome decide \
  --tag 'tier=["a","b"]' --tag 'labels=["security"]' --as=json

# (d) UNDECLARED key, empty value
/tmp/intrastate-spike flow resolve --model carrier-table.toml --outcome decide \
  --tag tier=free --tag 'labels=["security"]' --tag extra= --as=json

# (e) DECLARED key, empty value
/tmp/intrastate-spike flow resolve --model carrier-table.toml --outcome decide \
  --tag tier= --tag 'labels=["security"]' --as=json

# (f) undeclared scalar alone (different row, to show it is carried, not guarded)
/tmp/intrastate-spike flow resolve --model carrier-table.toml --outcome decide \
  --tag tier=paid --tag 'labels=["docs"]' --tag extra=plain --as=json
```

## Results

### (a) undeclared scalar + undeclared array — `a-undeclared-both.txt`, exit 2

```
{"code":"flow-tag-invalid","message":"the tag `extras` is not set-valued; got the array literal [\"a\",\"b\"]","param":"extras"}
```

REFUSED at HEAD. The refusal names `extras` — the UNDECLARED key — and
speaks of a set-valuedness the model never declared. `extra=plain` never
gets a chance to be echoed because admission fails whole-request.

### (a2) fallback diff arm: undeclared SCALAR only — `a2-undeclared-scalar-only.txt`, exit 0

```
{"type":"ok","data":{"model":"carrier-table.toml","revision":"","observed":{"extra":"plain","labels":"[\"security\"]","tier":"free"},"owned":{},"readers":[],"outcome":"decide","rule":"free","gates":[],"emit":{"plan":"basic"},"dispositions":{},"next":{},"writes":{},"clear":[],"escaped":false}}
```

`observed` carries `"extra":"plain"` byte-for-byte as given.

### (b) baseline, no carriers — `b-baseline-no-carriers.txt`, exit 0

```
{"type":"ok","data":{"model":"carrier-table.toml","revision":"","observed":{"labels":"[\"security\"]","tier":"free"},"owned":{},"readers":[],"outcome":"decide","rule":"free","gates":[],"emit":{"plan":"basic"},"dispositions":{},"next":{},"writes":{},"clear":[],"escaped":false}}
```

### DIFF — which pair was diffed

(a) refuses at HEAD, so per the protocol the diffed pair is **(a2) vs (b)** —
the undeclared SCALAR carrier against the same resolve without it.
`selection-diff-a2-vs-b.txt` is EMPTY: the projected selection fields
(`rule`, `emit`, `dispositions`, `gates`, `next`, `writes`, `clear`,
`escaped`, `outcome`) are byte-identical — both select rule `free`, emit
`{"plan":"basic"}`, `escaped:false`. The ONLY payload difference is the
`observed` echo gaining `"extra":"plain"`:

```
(a2) observed: {"extra": "plain", "labels": "[\"security\"]", "tier": "free"}
(b)  observed: {"labels": "[\"security\"]", "tier": "free"}
```

### (c) DECLARED scalar handed an array — `c-declared-scalar-array.txt`, exit 2

```
{"code":"flow-tag-invalid","message":"the tag `tier` is not set-valued; got the array literal [\"a\",\"b\"]","param":"tier"}
```

Confirmed: `flow-tag-invalid`, "is not set-valued", exit 2. Note this is the
BYTE-IDENTICAL message template used for the undeclared key in (a) — only
the key name differs. Nothing in the message distinguishes a real
declaration conflict from the zero-decl artefact.

### (d) UNDECLARED key, empty value — `d-undeclared-empty.txt`, exit 2

```
{"code":"flow-tag-invalid","message":"the tag `extra` was given an empty value","param":"extra"}
```

### (e) DECLARED key, empty value — `e-declared-empty.txt`, exit 2

```
{"code":"flow-tag-invalid","message":"the tag `tier` was given an empty value","param":"tier"}
```

(d) and (e) are the same code and the same message template. The empty-value
arm is **declaration-independent** at HEAD: it fires from the `!isSet` branch
before any `ConformValue` call, so the zero decl and a real scalar decl reach
it identically.

### (g) extra run — DECLARED SET, empty value — `g-declared-set-empty.txt`, exit 2

Run to test the limit of (d)/(e)'s declaration-independence:

```
{"code":"flow-tag-invalid","message":"the tag `labels` is set-valued and takes a JSON array literal; got ","param":"labels"}
```

A DECLARED **set** given an empty value does NOT take the empty-value arm at
all — `isSet` is true, so it falls to the set-shape arm and gets a different
message (note the trailing space before the closing quote: the empty value is
interpolated). The empty-value message is therefore independent of *whether*
a declaration exists, but NOT independent of a declaration's *kind*.

### (f) undeclared scalar alone — `f-undeclared-scalar-alone.txt`, exit 0

```
{"type":"ok","data":{"model":"carrier-table.toml","revision":"","observed":{"extra":"plain","labels":"[\"docs\"]","tier":"paid"},"owned":{},"readers":[],"outcome":"decide","rule":"paid","gates":[],"emit":{"plan":"pro"},"dispositions":{},"next":{},"writes":{},"clear":[],"escaped":false}}
```

Passes today, selects rule `paid` off the declared `tier`, and echoes
`"extra":"plain"` verbatim. The carrier does not perturb selection.

## Verbatim values for downstream fixtures

Read from the run, not composed:

- carried scalar echo, key `extra`: `"extra":"plain"` (JSON string `plain`).
- undeclared array refusal:
  code `flow-tag-invalid`, param `extras`,
  message ``the tag `extras` is not set-valued; got the array literal ["a","b"]``
  (the array literal appears in the message exactly as the caller spelled it,
  unsorted and un-canonicalised), exit 2.
- declared scalar array refusal: code `flow-tag-invalid`, param `tier`,
  message ``the tag `tier` is not set-valued; got the array literal ["a","b"]``, exit 2.
- empty-value refusal (undeclared): code `flow-tag-invalid`, param `extra`,
  message ``the tag `extra` was given an empty value``, exit 2.
- empty-value refusal (declared scalar): code `flow-tag-invalid`, param `tier`,
  message ``the tag `tier` was given an empty value``, exit 2.
- declared set echo is re-canonicalised into `observed` as a JSON-array
  STRING: `"labels":"[\"security\"]"`.

## Which RDR claims this run SUPPORTS

1. "an undeclared scalar passes, an undeclared array refuses" — SUPPORTED
   exactly, by (f)/(a2) and (a).
2. The refusal message points at a declaration that does not exist —
   SUPPORTED: (a) says ``the tag `extras` is not set-valued`` for a key the
   model never declares, and it is the same string (c) produces for a real
   declaration, so a caller cannot tell the two apart.
3. Undeclared keys currently reach the payload's `observed` echo unchanged —
   SUPPORTED: (a2)/(f) echo `"extra":"plain"` byte-for-byte.
4. A carrier key does not perturb selection — SUPPORTED: the (a2)-vs-(b)
   selection diff is empty. This is the runtime leg of the "only guard atoms
   make dimensions" assumption, for the SCALAR carrier that HEAD admits.
5. The empty-value arm's code and message — SUPPORTED as written: the code
   is `flow-tag-invalid` and the message is ``the tag `X` was given an empty
   value``, matching the RDR's quoted phrase verbatim.

## Which RDR claims this run does NOT establish / qualifies

- The array-carrier leg of claim 4 is **NOT** runtime-verifiable at HEAD:
  (a) refuses, so no run can show that an undeclared SET carrier leaves
  selection unchanged. That leg stays a design claim about the proposed
  change, provable only after implementation. This spike diffed the SCALAR
  carrier pair (a2 vs b) and says so explicitly.
- **LOUD FLAG.** The empty-value arm's declaration-independence holds for
  scalars but NOT across kinds. Run (g): a declared `set` given `--tag
  labels=` diverts to ``the tag `labels` is set-valued and takes a JSON array
  literal; got `` — never to "was given an empty value". Any contract text
  saying "an empty value is always refused with `was given an empty value`"
  would be wrong for declared sets. This matters directly if RDR 0020 makes
  the undeclared key a pure carrier: the carrier's empty value would then have
  to pick one of these two arms deliberately, and today's answer is the scalar
  one only because the zero decl's kind is empty, not because anyone chose it.
- Nothing here exercises `--write` or `--outcome`, which share
  `canonicalValue` via `codeInvalidFor`; the `flow-write-invalid` twin of
  each refusal was not run.

## Files

| case | invocation | capture |
| --- | --- | --- |
| a | undeclared scalar + undeclared array | `a-undeclared-both.txt` |
| a2 | undeclared scalar only (diff arm) | `a2-undeclared-scalar-only.txt` |
| b | no carrier flags | `b-baseline-no-carriers.txt` |
| c | declared scalar given an array | `c-declared-scalar-array.txt` |
| d | undeclared key, empty value | `d-undeclared-empty.txt` |
| e | declared key, empty value | `e-declared-empty.txt` |
| f | undeclared scalar alone | `f-undeclared-scalar-alone.txt` |
| g | declared SET, empty value (extra) | `g-declared-set-empty.txt` |
| diff | (a2) vs (b) selection projection | `selection-diff-a2-vs-b.txt` (empty), `sel-a2.json`, `sel-b.json` |
| fixture | model under test | `carrier-table.toml` |
