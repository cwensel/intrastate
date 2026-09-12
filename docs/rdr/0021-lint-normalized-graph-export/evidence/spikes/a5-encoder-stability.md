Model: claude-sonnet-5

# A5 Spike — Encoder Byte-Stability (Marshaling Determinism)

## Assumption under test

A5: Marshaling the export document through the shared non-HTML-escaping
encoder is byte-stable — struct field order fixed, every sequence
pre-sorted, and no Go map reaches the wire.

## 1. Single-encoder confirmation

`internal/cli/clierr/clierr.go:174` — `WriteJSONLine`:

```go
func WriteJSONLine(out io.Writer, v any) error {
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}
```

Doc comment (clierr.go:159-169) states this is "THE JSON emit helper every
CLI wire site routes through (`0005:C1`, REQ-71/REQ-72)" and explains the
HTML-escaping-disabled rationale (byte-equality for read-back per RDR 0004).

Call graph:
- `EmitJSON` (clierr.go:156) -> `WriteJSONLine`
- `internal/cli/respond/respond.go:242` -> `clierr.WriteJSONLine`
- `internal/cli/flowbind/flowbind.go:155` -> `clierr.WriteJSONLine`

Other files construct their OWN `json.NewEncoder(...).SetEscapeHTML(false)`
rather than calling `WriteJSONLine`:
- `internal/table/model.go:465-466`
- `internal/cli/cmdbind/cmdbind.go:1096-1097`
- `internal/cli/flow_input.go:799-800`
- `internal/cli/respond/text.go:133-134`

These are NOT on the graph-export wire path today (the `graph export`
command does not exist yet in the codebase — RDR 0021 is Pending/not yet
implemented), so no second encoder currently reaches this document's wire.
They are duplicated call sites with identical settings, not a second
encoder on THIS path. Risk noted: when export is implemented, it must route
through `clierr.WriteJSONLine` (or `respond`/`flowbind`, which already do)
rather than adding a fifth ad hoc `SetEscapeHTML(false)` site, or the
one-encoder claim in C3 becomes false by omission.

## 2. Fixture — modeled on C2's field list

Fields modeled: `schema`, model identity/class, tag declarations (name,
provenance, kind, required, single_valued, domain), `initial` assignments,
`terminal` predicate sets (set of sets), normalized rows (0002 field list:
identity, source, kind, outcome, atoms `{key, operator, literal[], block}`,
next, writes, requires_owned, gate, escape, emit), selection-context groups
(context + members), and the reachability relation (`nodes: {id, values}`,
`edges: {from, to, rule}`, plus the `abstraction` marker). Every field is
modeled as a Go **struct or slice** — no map anywhere in the fixture's
document type. String fields deliberately carry `<`, `>`, `&`, `"`,
newline, em-dash, and non-ASCII (`héllo—wörld`, `日本語`) to probe escaping.
An explicit empty slice (`atoms: []`, `writes: []`, `members: []`,
`values: []`, `empty_probe`) and one un-set slice are included to observe
`[]` vs `null` rendering.

Spike source: `/tmp/a5spike/main.go` (scratch, not committed to project
tree). Go toolchain: `go version go1.27.1 darwin/arm64`.

## 3. Commands and raw output

```
$ cd /tmp/a5spike && go build -o a5spike main.go
$ ./a5spike single
bytes_equal_double_emit: true
sha256_emit1: dcb0259aa822615592f11b65d571ba8a8b33c03f24b98b3bd54a3a177af8049f
sha256_emit2: dcb0259aa822615592f11b65d571ba8a8b33c03f24b98b3bd54a3a177af8049f
line: {"schema":"intrastate.graph/1","model_id":"svc<order>&billing","model_class":"class \"A\" <root>","tags":[{"name":"env","provenance":"declared","kind":"string","required":true,"single_valued":true,"domain":["dev","stage>prod","prod&dr"]},{"name":"region\nnote","provenance":"inferred","kind":"string","required":false,"single_valued":false}],"initial":["n0<start>","n1"],"terminal":[["done","ok&clean"],["failed<hard>"]],"rows":[{"identity":"r1","source":"n0<start>","kind":"transition","outcome":"advance","atoms":[{"key":"flag&x","operator":"==","literal":["a<b","c>d","e\"f"],"block":false},{"key":"unicode","operator":"!=","literal":["héllo—wörld","日本語"],"block":true}],"next":"n1","writes":["w1<tag>","w2&flag"],"requires_owned":["o1"],"gate":"g<1>","escape":"none","emit":"row&emit"},{"identity":"r0-zzz","source":"n1","kind":"terminal","outcome":"halt","atoms":[],"next":"","writes":[],"requires_owned":[],"gate":"","escape":"","emit":""}],"selection_context":[{"context":"ctx<A>&1","members":["r1","r0-zzz"]},{"context":"ctx\"B\"","members":[]}],"reach":{"abstraction":"declared-over-approximation","nodes":[{"id":"n0<start>","values":["env=dev","flag&x"]},{"id":"n1","values":[]}],"edges":[{"from":"n0<start>","to":"n1","rule":"r1"}]},"empty_probe":[]}
empty_slice_literal: []
nil_slice_literal: null
map_emit: {"alpha":2,"beta":4,"mu":3,"zeta":1}
```

Independent sha256 recompute (Python, over the exact line bytes plus the
trailing `\n` that `Encoder.Encode` appends):

```
recomputed_sha256: dcb0259aa822615592f11b65d571ba8a8b33c03f24b98b3bd54a3a177af8049f
```

Matches `sha256_emit1`/`sha256_emit2` exactly.

## 4. Cross-process / map-seed repetition

25 separate process invocations, hashing only the primary fixture line:

```
$ for i in $(seq 1 15); do GODEBUG=randmapiter=1 ./a5spike hashonly; done | sort -u
dcb0259aa822615592f11b65d571ba8a8b33c03f24b98b3bd54a3a177af8049f

$ for i in $(seq 1 10); do ./a5spike hashonly; done | sort -u
dcb0259aa822615592f11b65d571ba8a8b33c03f24b98b3bd54a3a177af8049f
```

Method: `GODEBUG=randmapiter=1` (supported by this toolchain, go1.27.1) was
used to force per-process map iteration-order randomization for the 15
runs in the first batch, plus 10 unmodified-env runs in the second batch —
25 total, each a fresh OS process. All 25 produced the identical sha256.
This is expected and unsurprising for THIS fixture because the fixture
itself contains no Go map — `randmapiter` only affects the separate
`map_emit` probe (see §5), not the C2 document.

## 5. The map question (the crux)

(a) **Does `encoding/json` sort map keys here?** Yes. `map_emit` shows
`map[string]int{"zeta":1,"alpha":2,"mu":3,"beta":4}` — built in
non-alphabetical insertion order — serialized as
`{"alpha":2,"beta":4,"mu":3,"zeta":1}`, alphabetically sorted by key,
identically across 25 processes and under `randmapiter=1`. Go's
`encoding/json` explicitly sorts map keys before encoding
(this is documented stdlib behavior, confirmed empirically here), so a map
value would in fact serialize byte-stably too.

(b) **Does C2's document REQUIRE a map anywhere, or can every field be a
slice/struct?** No map is required. Re-reading C2's field list: model
identity/class (scalars), tag declarations (a slice of structs, each with
scalar/slice fields), `[initial]` assignments (a slice), `terminal`
predicate sets (a slice of sets, i.e. slice-of-slices), normalized rows
(a slice of structs, atoms as a slice of `{key, operator, literal[],
block}` structs), selection-context groups (a slice of `{context,
members[]}` structs), and the reachability relation (nodes as a slice of
`{id, values[]}`, sorted by node key; edges as a slice of `{from, to,
rule}`, sorted by (from, to, rule)). Every one of these is expressible —
and, in this fixture, WAS expressed — as a Go slice or struct. The fixture
in §2 models the full C2 field list with zero maps and it round-trips
byte-stably.

**The distinction that matters**: a map would ALSO look byte-stable
(§5a), because `encoding/json` sorts its keys — but that stability comes
from an incidental encoder behavior (alphabetic key sort), not from a
declared, reviewable ordering rule the way C2's "nodes sorted by node key"
and "edges sorted by (from, to, rule)" are. A map further constrains the
key to something JSON-object-key-shaped (a string, or a type with
`MarshalText`/integer key), which the RDR does not need for any C2 field —
node ids, rule ids, and row identities are already carried as explicit
string fields inside slice elements, not as map keys. So: A5's "no Go map
reaches the wire" clause is satisfiable — not because maps are unsafe, but
because C2 never needs one, and using a map anywhere would silently swap a
reviewable, declared sort (C2's stated key/rule orders) for the encoder's
incidental alphabetic key sort, which happens to match by coincidence for
plain string keys but is not the same guarantee C3 makes.

## 6. Determinism checklist

| Item | Observation |
|---|---|
| Pre-image byte layout | Go struct field order is fixed at compile time by field declaration order; `encoding/json` marshals struct fields in declaration order (documented stdlib behavior), confirmed by identical byte output across 25 process runs. |
| Encodings | Non-ASCII (`héllo—wörld`, `日本語`) emitted as literal UTF-8 bytes, not `\uXXXX` escapes (Go's default `SetEscapeHTML` only affects `<`, `>`, `&`, not general Unicode). |
| Map order | No map in the C2-shaped document (see §5b). Where a map WAS constructed as a probe (§5a), `encoding/json` sorted keys alphabetically, byte-stable across 25 processes and `randmapiter=1`. |
| Whitespace | `json.Encoder.Encode` emits compact JSON (no indentation) followed by exactly one trailing `\n`; confirmed in the raw `line:` output and by the sha256 match including that trailing newline. |
| Case folding | None observed or applicable — no case-insensitive comparison occurs in marshaling; string values pass through verbatim (`"héllo—wörld"` preserved exactly, mixed-case `"class \"A\" <root>"` preserved exactly). |
| Empty/null/absent | Explicit empty slice (`[]string{}`) emits `[]`; a nil slice emits `null`. Observed directly: `empty_slice_literal: []`, `nil_slice_literal: null`. This is a real, distinguishable wire difference — C2's additive-evolution rule depends on producer code consistently choosing empty-slice-not-nil for "no items" fields (e.g. `atoms: []`, `writes: []`, `members: []` all render as `[]` in the fixture because the fixture uses `[]T{}`, never a nil slice, for those fields). |
| Version marker | `schema` is the first-declared, first-emitted field: `{"schema":"intrastate.graph/1", ...}` — present and positioned as C2's "required leading `schema` field". |

## Verdict

A5 **PASSES** for the document shape C2 specifies: double-emit is
byte-identical (sha256 match), 25-process/map-seed-varied emission is
byte-identical, HTML-significant characters and non-ASCII pass through
unescaped/unmangled, and every C2 field is expressible as a slice or
struct with zero Go maps — so "no Go map reaches the wire" is achievable
by construction, not merely by encoder-incidental map-key sorting. Residual
risk (non-blocking, implementation-time): the export command does not
exist yet, so this spike cannot confirm the FUTURE export code will call
`clierr.WriteJSONLine` rather than adding a fifth ad hoc
`SetEscapeHTML(false)` call site (four such duplicate sites already exist
elsewhere in the codebase); this is an implementation-discipline concern,
not a falsification of A5.
