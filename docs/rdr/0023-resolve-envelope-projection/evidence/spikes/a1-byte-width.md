# A1 spike: shipped projection byte width

Date: 2026-08-29

Projection measured: delete the ECHO group (`model`, `observed`, `owned`,
`readers`, `outcome`) from the success payload; keep the full PLAN group
(`revision`, `rule`, `gates`, `emit`, `next`, `writes`, `clear`, `escaped`,
plus `escape_class` when present). The success payload sits under `.data`
of the `{"type":"ok","data":{...}}` envelope, so the deletion is applied at
`.data.*`:

```
jq -jc 'del(.data.model,.data.observed,.data.owned,.data.readers,.data.outcome)'
```

jq key deletion preserves key order; this approximates the shipped
key-deletion projection (byte-exactness is A2's spike). Default byte counts
are the raw JSON line without the trailing newline; projected byte counts
are the `jq -jc` output.

## Byte table

| call | shape | default | projected | saved | saved % |
|------|-------|--------:|----------:|------:|--------:|
| S1  | decision table, 2 facts (pricing 2x2)           | 290 | 152 | 138 | 47.6% |
| S1b | decision table, 48 facts (synthetic, motivating) | 708 | 146 | 562 | 79.4% |
| S2  | state machine, writes+clear (release ship-clean) | 392 | 191 | 201 | 51.3% |
| S2b | state machine, gate+writes (release begin)       | 412 | 238 | 174 | 42.2% |

## S1 — pricing decision table (models/examples/pricing-decision-table.toml)

```
intrastate flow resolve --model models/examples/pricing-decision-table.toml \
    --outcome decide --tag tier=paid --tag region=eu --as json
```

Default:

```json
{"type":"ok","data":{"model":"models/examples/pricing-decision-table.toml","revision":"","observed":{"region":"eu","tier":"paid"},"owned":{},"readers":[],"outcome":"decide","rule":"paid-eu","gates":[],"emit":{"dpa":"required","plan":"pro"},"next":{},"writes":{},"clear":[],"escaped":false}}
```

Projected:

```json
{"type":"ok","data":{"revision":"","rule":"paid-eu","gates":[],"emit":{"dpa":"required","plan":"pro"},"next":{},"writes":{},"clear":[],"escaped":false}}
```

## S1b — synthetic 48-fact decision table (pricing-48.toml, scratch)

Synthetic model, authored for this spike on the pricing example's grammar:
48 required observed enum tags `f01`..`f48` (each `domain = ["a", "b"]`),
two rows discriminating on `f01` only. The per-tag stanza is identical for
all 48 tags.

```toml
outcomes = ["decide"]

[model]
id = "pricing48"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

# ... one stanza per tag, f01 through f48, all identical to:
[tags.f01]
provenance = "observed"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[[rule]]
id = "low"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.f01]
eq = "a"
[rule.emit]
plan = "basic"
dpa = "none"

[[rule]]
id = "high"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.f01]
eq = "b"
[rule.emit]
plan = "pro"
dpa = "required"
```

```
intrastate flow resolve --model pricing-48.toml --outcome decide \
    --tag f01=a --tag f02=b ... --tag f47=a --tag f48=b --as json
```

(48 `--tag` flags: `a` for odd-numbered tags, `b` for even-numbered.)

Default:

```json
{"type":"ok","data":{"model":"pricing-48.toml","revision":"","observed":{"f01":"a","f02":"b","f03":"a","f04":"b","f05":"a","f06":"b","f07":"a","f08":"b","f09":"a","f10":"b","f11":"a","f12":"b","f13":"a","f14":"b","f15":"a","f16":"b","f17":"a","f18":"b","f19":"a","f20":"b","f21":"a","f22":"b","f23":"a","f24":"b","f25":"a","f26":"b","f27":"a","f28":"b","f29":"a","f30":"b","f31":"a","f32":"b","f33":"a","f34":"b","f35":"a","f36":"b","f37":"a","f38":"b","f39":"a","f40":"b","f41":"a","f42":"b","f43":"a","f44":"b","f45":"a","f46":"b","f47":"a","f48":"b"},"owned":{},"readers":[],"outcome":"decide","rule":"low","gates":[],"emit":{"dpa":"none","plan":"basic"},"next":{},"writes":{},"clear":[],"escaped":false}}
```

Projected:

```json
{"type":"ok","data":{"revision":"","rule":"low","gates":[],"emit":{"dpa":"none","plan":"basic"},"next":{},"writes":{},"clear":[],"escaped":false}}
```

## S2 — release grammar, ship-clean (models/examples/release-grammar.toml)

Owned state seeded to `phase=building`, `build-id=pending` via:

```
intrastate flow set-state --model models/examples/release-grammar.toml \
    --artifact release=release-state.artifact \
    --write phase=building --write build-id=pending --as json
```

```
intrastate flow resolve --model models/examples/release-grammar.toml \
    --artifact release=release-state.artifact \
    --outcome ship --tag 'checks=["tests","signoff"]' --tag risk=0 --as json
```

Default:

```json
{"type":"ok","data":{"model":"models/examples/release-grammar.toml","revision":"","observed":{"checks":"[\"signoff\",\"tests\"]","risk":"0"},"owned":{"build-id":"pending","phase":"building"},"readers":["release-state"],"outcome":"ship","rule":"ship-clean","gates":[],"emit":{},"next":{"build-id":"<clear>","phase":"shipped"},"writes":{"phase":"shipped"},"clear":["build-id"],"escaped":false}}
```

Projected:

```json
{"type":"ok","data":{"revision":"","rule":"ship-clean","gates":[],"emit":{},"next":{"build-id":"<clear>","phase":"shipped"},"writes":{"phase":"shipped"},"clear":["build-id"],"escaped":false}}
```

## S2b — release grammar, begin (gate on the selected row)

Owned state seeded to `phase=idle`, `build-id=none` via `flow set-state`
(same form as S2).

```
intrastate flow resolve --model models/examples/release-grammar.toml \
    --artifact release=release-idle.artifact \
    --outcome build --tag risk=0 --tag 'checks=[]' --as json
```

Default:

```json
{"type":"ok","data":{"model":"models/examples/release-grammar.toml","revision":"","observed":{"checks":"[]","risk":"0"},"owned":{"build-id":"none","phase":"idle"},"readers":["release-state"],"outcome":"build","rule":"begin","gates":[{"id":"change-window","result":"allow"}],"emit":{},"next":{"build-id":"pending","phase":"building"},"writes":{"build-id":"pending","phase":"building"},"clear":[],"escaped":false}}
```

Projected:

```json
{"type":"ok","data":{"revision":"","rule":"begin","gates":[{"id":"change-window","result":"allow"}],"emit":{},"next":{"build-id":"pending","phase":"building"},"writes":{"build-id":"pending","phase":"building"},"clear":[],"escaped":false}}
```

## Reading

The seed's 72% figure was the echo SHARE of a 48-fact call; the honest
shipped reduction on that shape is 79.4% (708 -> 146 bytes), larger because
the echo group also carries `owned`, `readers`, and `outcome`. On small
payloads the floor is set by the retained PLAN group: 42-51% on the
state-machine shapes and 47.6% on the minimal 2x2 table. The reduction
grows with fact count and is material on every shape measured.
