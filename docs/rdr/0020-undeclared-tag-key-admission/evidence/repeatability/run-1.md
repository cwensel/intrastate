model: claude-fable-5-1
variant: lite (profile: mid)

Spans read: 0020:C1, 0020:MVV, 0020:S1..S4 first. Widened past the contract
set to: 0020:D-identity, D-wire-byte-format, D-naming, D-selection-predicate,
§technical-design (incl. authority census, disposition, desk-trace,
illustrative-code), §proposed-solution (approach, infra audit, decision
rationale), §testing-strategy, and A1..A5. Widening was needed because C1
names `parseTags`, `canonicalValue`, `m.Tags[key]`, `resolve.Input.Observed`
but fixes no signatures, no argument-splitting grammar, and no error-struct
shape; the illustrative code and A1/A3 evidence supplied most of those.

## 1. Public API (module: `internal/cli/flow_input.go`, the `--tag` admission seam)

The RDR changes no exported symbol. The "public" surface is the CLI flag
grammar plus the refusal JSON on stderr and the `observed` echo on stdout.

```go
// parseTags admits every --tag name=value argument against the loaded,
// normalized model m. owned is the set of keys the CLI itself writes
// (REQ-27). It returns the admitted tags in argument order, or the first
// refusal encountered.
//   GUESS: exact parameter list/order and return type; the RDR only fixes
//   that parseTags is the sole caller of the zero-decl arm and that its
//   lookup becomes two-value.
func parseTags(args []string, m *table.Model, owned []string) ([]resolveTag, *cliError)

// canonicalValue conforms a value under a REAL declaration. After this
// change it is never called for an undeclared key; the "was given an
// empty value" arm is hoisted out of it (or is kept but unreachable —
// implementation latitude per D-selection-predicate).
//   GUESS: the `what` discriminator ("tag" | "write") and the return
//   shape; the illustrative code shows canonicalValue(key, value, decl, "tag").
func canonicalValue(key, value string, decl table.TagDecl, what string) (string, *cliError)
```

Refusal codes reachable from `parseTags` (all pre-existing; none minted —
D-naming, `flow-tag-undeclared` does not exist), in precedence order
(D-selection-predicate):

| step | code | condition | message (fixed where the RDR pins bytes) |
| --- | --- | --- | --- |
| 1 | `flow-tag-invalid` | argument is not `name=value` (grammar) | GUESS: existing HEAD text, unchanged |
| 2 | `flow-tag-reserved` | `key == table.RecognizedTagKey` (REQ-26) | unchanged |
| 3 | `flow-tag-owned` | key is CLI-owned (REQ-27) | unchanged |
| 4 | `flow-tag-duplicate` | `seen[key]` already set (REQ-28) | unchanged |
| 5 | `flow-tag-invalid` | `value == ""`, declared or not | ``the tag `X` was given an empty value`` (F4, byte-pinned) |
| 6 | (none — admit) | key ABSENT from `m.Tags` | value carried verbatim |
| 7 | `flow-tag-invalid` | declared key, kind/shape/domain conformance | e.g. ``the tag `tier` is not set-valued; got the array literal ["a","b"]`` (F3); ``... is set-valued and takes a JSON array literal; got `` (declared set, empty after canonicalisation — g fixture) |

Exit code for every refusal above: 2. Success: 0.

Wire shape of a refusal (F3/F4, byte-pinned by the fixtures):
`{"code":"flow-tag-invalid","message":"...","param":"<key>"}`.

Wire shape of the echo: the resolve payload's `observed` object,
`map[string]string` keyed by tag name, value = raw argument text after the
first `=`, e.g. `"extras":"[\"a\",\"b\"]"`. Declared set values keep the
canonical array form of `docs/cli-output-contract.md` §Set values on the
wire (D-wire-byte-format).

GUESS (RDR silence): how the `name=value` grammar splits — I take the
FIRST `=` as the separator so a value may itself contain `=`; the RDR only
says "`--tag` takes `name=value`".

GUESS (RDR silence): whether a declared SET key given an empty value is
refused at step 5 with the scalar-shaped message or falls through to step
7's set-specific message. C1 says the set-specific text "stays
declared-only, inside the arms below" and F4 pins only the undeclared
case; the illustrative code checks `value == ""` BEFORE the declared
lookup, which would make the set-specific empty-value message unreachable
for a literally-empty argument. I resolve it as: step 5 fires on the raw
empty string for every key; step 7's set-specific message remains for a
declared set whose array literal canonicalises to nothing (e.g. `[]`).

## 2. Three most important internal helpers

1. `parseTags` — the admission loop. Owns the precedence order above,
   the `seen` map, the two-value `decl, declared := m.Tags[key]` lookup
   (D-identity: byte-exact presence, no folding, no sentinel), and the
   carrier branch that appends the verbatim value and `continue`s
   without touching `canonicalValue`.

2. `canonicalValue` — declared-only conformance. Given a real
   `table.TagDecl`, canonicalises set literals to the wire array form,
   kind-checks (`!isSet && looksArray` -> "is not set-valued"), and
   domain-checks enums. After this change its `!isSet && looksArray`
   arm is truthful because `decl` is always a loaded declaration whose
   `Kind` was required at load (one of RDR 0003's five tokens).
   GUESS: `looksArray` (first-byte `[` sniff) survives unchanged as a
   helper inside this function; the RDR calls the sniff a bug only when
   applied to an UNDECLARED value.

3. `userErr` / refusal constructor (name GUESS from the illustrative
   `userErr(codeTagInvalid, key, "...")`) — builds the
   `{code, message, param}` refusal so every arm in the table above
   emits the byte-shape F3/F4 pin. The constants
   `codeTagInvalid`/`codeTagDuplicate`/`codeTagReserved`/`codeTagOwned`
   in the same file's constant block are unchanged (A5).

Runner-up (unchanged, cited for the boundary): `observedTagMap` /
`kernelTags` in `internal/cli/flow_exec.go`, the verbatim copies that
move admitted tags into `resolve.Input.Observed` (A3).

## 3. Data model across the boundary

```go
// One admitted --tag, produced by parseTags, consumed by the resolve
// assembly in flow_exec.go and the payload echo in flow_resolve.go.
//   GUESS: field names; the illustrative code uses resolveTag{Key, Value}.
type resolveTag struct {
    Key   string // byte-exact argument name; declared-ness NOT recorded
    Value string // declared: canonical form; undeclared: raw argument bytes
}

// Kernel input view (existing, unchanged): the RDR pins its type.
// resolve.Input.Observed map[string]string
//   key   -> tag name
//   value -> string exactly as parseTags emitted it

// Loaded declaration (existing, unchanged): the discriminator.
// m.Tags map[string]table.TagDecl — two-value lookup is the ONLY
//   declared/undeclared signal; a zero TagDecl never occurs for a
//   declared key because the loader requires a kind.

// Refusal on stderr (existing wire shape, byte-pinned by F3/F4):
// {"code":"flow-tag-invalid","message":"<text>","param":"<key>"}
//   exit 2

// Resolve payload echo on stdout (existing; F1 pins the shape):
// "observed": {"extra":"plain","labels":"[\"security\"]","tier":"free"}
//   sole field that changes when a carrier is added (F2: rule, emit,
//   gates, next, writes, clear, escaped all identical)
```

Nothing is persisted. No flag on `resolveTag` marks a carrier — the
RDR's disposition table shows the payload is the only diagnostic
surface, and a misspelled declared key is accepted silently (accepted
open-world cost).

GUESS (RDR silence): whether declared-ness is recorded anywhere
downstream (e.g. for a future strict mode per 0024's shape). I omit it;
the RDR names strictness as compatible later work, not part of this change.

## 4. Top-level pseudo-code of `parseTags`

```
parseTags(args, m, owned):
    out  := []
    seen := {}
    for raw in args:                                   # argument order preserved
        key, value, ok := split raw at FIRST '='       # GUESS: first-'=' split
        if !ok:
            return refuse(flow-tag-invalid, param=raw?)  # grammar arm, HEAD text (GUESS on param)
        # provenance guards — unchanged, still first, still before any accessor
        if key == table.RecognizedTagKey:
            return refuse(flow-tag-reserved, key)          # REQ-26
        if key in owned:
            return refuse(flow-tag-owned, key)             # REQ-27
        if seen[key]:
            return refuse(flow-tag-duplicate, key)         # REQ-28 — precedes empty-value
        seen[key] = true
        # hoisted empty-value arm: binds declared AND undeclared keys
        if value == "":
            return refuse(flow-tag-invalid, key,
                          "the tag `"+key+"` was given an empty value")   # F4 bytes
        decl, declared := m.Tags[key]                  # two-value, byte-exact (D-identity)
        if !declared:
            # PURE CARRIER (0020:C1): no canonicalise, no kind/shape/domain,
            # no comparison; array literals ride through as raw bytes.
            out.append(resolveTag{Key: key, Value: value})
            continue
        # declared-only conformance; decl.Kind is guaranteed non-zero by load
        canonical, ce := canonicalValue(key, value, decl, "tag")
        if ce != nil:
            return nil, ce      # "is not set-valued" (F3), enum domain, set-empty ("...got "), etc.
        out.append(resolveTag{Key: key, Value: canonical})
    return out, nil

# downstream (unchanged): out -> observedTagMap/kernelTags -> resolve.Input.Observed
#   -> kernel readers (matches/has/Lookup/conflicting) can only read rule-atom keys,
#      all load-checked, so a carrier key is unreadable (A2)
#   -> payload "observed" echo emits every entry byte-for-byte (F1)
```
