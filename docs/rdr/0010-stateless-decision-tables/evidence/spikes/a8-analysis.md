Model: claude-opus-5

# A8 spike — strict decoding as the old-binary guard (RDR cli/0010)

## Assumption under test

> **A8** Strict decoding refuses an unknown `[model]` key, so a
> `class = "decision-table"` model loaded by a binary predating this RDR
> refuses as `unknown_schema_field` rather than silently loading as a
> machine.

Cited evidence: `internal/table/source.go::sourceModel` is decoded strictly
(`0002:A1`); `CatUnknownSchemaField` is the category.

## Verdict

**A8 VERIFIED.** The claim holds exactly as written, on both the `lint` and
the `flow` load paths. The spike additionally establishes a *second*
protection the RDR's Risk section does not currently name, and *one hole*
it also does not name.

## Setup

The current `HEAD` binary **is** the pre-change binary — `class` is not yet
in the layout, so today's build is a faithful stand-in for "a binary
predating this RDR".

```
$ make build
$ bin/intrastate version
2886a69 (commit 2886a69, built 2026-08-27T00:06:38Z)
```

`a8-baseline.toml` is a minimal MVV-shaped model authored to the JDR 0001
§D7 closed layout. It loads and lints clean on that binary:

```
$ bin/intrastate lint --model a8-baseline.toml --as=json
{"type":"ok","data":{"findings":[]}}
exit=0
```

Every variant below is that baseline plus exactly one edit.

## The strictness mechanism

`internal/table/source.go:113-124`:

```go
func decodeStrict(src []byte, dst any) error {
	dec := toml.NewDecoder(bytes.NewReader(src))
	dec.DisallowUnknownFields()                       // :115
	if err := dec.Decode(dst); err != nil {
		var strict *toml.StrictMissingError
		if asStrictError(err, &strict) {
			return fail(CatUnknownSchemaField, strict.Error())  // :119
		}
		return fail(CatMalformedTOML, err.Error())
	}
	return nil
}
```

The load call sequence is `internal/table/load.go:30` `Load` →
`load.go:53-56` pass 2 → `decodeStrict(src, &doc)`. `dst` is `sourceDoc`
(`source.go:14-26`), whose `Model` field is `*sourceModel`
(`source.go:28-38`). `sourceModel` declares exactly `id`, `version`,
`description`, `metadata` — no `class`. `DisallowUnknownFields` is
**recursive** over the whole struct graph, which is why `[[rule]]` is
covered too (see below).

Category identifier: `internal/table/category.go:18`
`CatUnknownSchemaField Category = "unknown_schema_field"`.

## Results

| variant | edit | lint exit | observed code |
|---|---|---|---|
| `a8-baseline` | (none) | 0 | `{"type":"ok","data":{"findings":[]}}` |
| `a8-model-class` | `class` on `[model]` | 2 | `unknown_schema_field` |
| `a8-rule-emit-table` | `[rule.emit]` sub-table | 2 | `unknown_schema_field` |
| `a8-rule-emit-inline` | `emit = {…}` on `[[rule]]` | 2 | `unknown_schema_field` |
| `a8-rule-unknown-scalar` | `novel_key` on `[[rule]]` | 2 | `unknown_schema_field` |
| `a8-root-class` | `class` at root | 2 | `unknown_schema_field` |
| `a8-root-row-table` | whole `[[row]]` table at root | 2 | `unknown_schema_field` |
| **`a8-metadata-class`** | **`class` under `[model.metadata]`** | **0** | **none — silently ignored** |
| `a8-v2-plus-class` | `version = 2` + `class` | 2 | `unsupported_version` (preempts) |

Full transcript with stdout and exit codes: `a8-strict-decode.out`.

### 1. `class = "decision-table"` on `[model]` — REFUSES as claimed

```
$ bin/intrastate lint --model a8-model-class.toml --as=json
{"code":"model-invalid","message":"model does not conform to the transition-model schema","param":"model","detail":"unknown_schema_field: strict mode: fields in the document are missing in the target struct","findings":[{"code":"unknown_schema_field","message":"unknown_schema_field: strict mode: fields in the document are missing in the target struct","locator":".../a8-model-class.toml:1"}]}
exit=2
```

The `flow` load path surfaces the same finding under a different envelope
code — `flow-model-invalid` rather than `model-invalid` — but the inner
finding `code` is identical:

```
$ bin/intrastate flow next --model a8-model-class.toml --as=json
{"code":"flow-model-invalid","message":"the selected model could not be loaded","findings":[{"code":"unknown_schema_field",...}]}
exit=2
```

`flow resolve` is byte-identical to `flow next` here: the refusal happens
in load, before outcome selection or artifact binding.

A8's "If wrong" branch — *an old binary lints a decision table as a
rootless machine and reports `graph-dangling-edge`* — **did not occur**.
No graph lint runs at all; load fails first and fail-fast, so exactly one
finding is returned.

### 2. `[rule.emit]` ALSO refuses — a second, unrecorded protection

This was the open question. All three `[[rule]]`-level probes refuse with
the same `unknown_schema_field`:

- `[rule.emit]` as a sub-table
- `emit = { outcome = "ok", reason = "advanced" }` inline
- a plain unknown scalar `novel_key = "surprise"`

**An old binary cannot silently DROP an authored emit block.** Strictness
is recursive through `sourceRule` (`source.go:74-85`), which declares only
`id`/`use`/`source`/`clear`/`gate`/`escape`/`match`/`guard`/`write`. This
is a second protection worth recording in the Risk section: the `[model]`
key is not the only thing standing between an old binary and a
decision-table document — every rule-level novelty is caught too. The
silent-failure mode the task flagged as a risk **does not exist**.

### 3. The hole: `[model.metadata]` is NOT strict

```
$ bin/intrastate lint --model a8-metadata-class.toml --as=json
{"type":"ok","data":{"findings":[]}}
exit=0
```

`sourceModel.Metadata` is `map[string]any` (`source.go:37`), and the
comment at `source.go:33-36` says this is deliberate — "it decodes as a
free-form map so strict decoding descends no further". So
`[model.metadata] class = "decision-table"` loads clean and is carried
through untouched (`load.go:113`, `l.model.Metadata = m.Metadata`).

This does not refute A8, which speaks only of an unknown `[model]` key.
It is a **constraint on where the discriminator may live**: if RDR 0010
ever placed `class` inside the sanctioned metadata namespace instead of at
`[model]`, the entire old-binary guard would evaporate silently. Worth
naming so nobody "tidies" the key into metadata later.

### 4. The version gate preempts strict decoding

With both `version = 2` and `class`, the refusal is `unsupported_version`,
not `unknown_schema_field`:

```
{"code":"model-invalid",...,"detail":"unsupported_version: version 2; this RDR accepts version 1 only",...}
exit=2
```

This is by design (`load.go:32-50`, pass 1 before pass 2; the comment at
`source.go:97-101` states the ordering obligation). It gives RDR 0010 a
choice of forward-compat lever: bump `version` for a *loud, self-describing*
refusal naming the version, or add `class` at v1 for a *generic*
`unknown_schema_field`. The version bump produces the strictly more
informative message for a human reading an old binary's output.

## Bearing on the RDR

1. **A8 stands as written** — no amendment needed to the claim itself.
2. The Risk section's premise that the `[model]` key protects old binaries
   is correct, and **understated**: `[[rule]]`-level keys refuse too.
3. Two facts the RDR does not currently name and may want to:
   `[model.metadata]` is an explicit strictness hole, and the version gate
   preempts `unknown_schema_field` — relevant if 0010 chooses a version
   bump over (or alongside) a `class` key.
