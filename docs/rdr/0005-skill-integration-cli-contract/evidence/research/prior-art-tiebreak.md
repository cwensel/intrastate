# RDR 0005 — prior art for the two open Stage-4 items

Research run 2026-08-24 against `$RDR_RESOURCES` corpora, the `langref` peer-CLI
checkout set, and the `../state-machines` sibling. Spikes under `../spikes/`.

## OPEN-1 — the `clierr.Finding` field set

### The real constraint set: five producers, not two

`JDR 0001 §D10` item 3 names four carriers; RDR 0009 is a fifth
(`0009::A9` + `§JD-8`, "envelope carrier for row identities"). Rendered in
`../spikes/finding-producer-matrix.out`:

| Producer | Authority | Must carry |
| --- | --- | --- |
| RDR 0002 load categories | JDR §D10 item 3 | `code` = category slug, `locator` = file:line |
| RDR 0005 kernel refusals | `0005::Failure Modes` | rows considered / matching rows / keys / atoms |
| RDR 0005 gate results | JDR §D9 example | `code` = `gate_denied`/`gate_allowed`, `param` = gate id |
| RDR 0006 lint findings | `0006::Normative Contracts` | code, model id, **severity**, message, rule/context id or source span, atom (key/operator/literal/block), escape failure **class** |
| RDR 0008 `reserved_tag_key` | JDR §D10 item 3 | per-key code, near-miss advisory in `hint` |
| RDR 0009 row identities | `0009::A9`, `§JD-8` | offending `RowRef` identities rendered to a wire field |

RDR 0005's five (`Code, Message, Param, Locator, Hint`) are a correct **subset**,
not a wrong set. The only question is where 0006's extra five live.

### The shipped type already states the project's convention

`internal/cli/clierr::CLIError` doc comment: "**Extend with new optional fields
as needed — keep them `omitempty` so the envelope stays append-only and stable
for tools.**" `0009::A9` closed the same question by source search: "no consumer
or test reads an exact field set". Widening is additive and safe.

### RDR 0006 contradicts itself on the atom's shape

- `0006::Technical Design`: "`Key`, `Operator`, `Literal`, `Block`, and `Class`
  are `string`" — five FLAT fields.
- `0006` illustrative JSON (`:1174`): `"atom": {"key":…, "operator":…,
  "literal":…, "block":…}` — a NESTED object.
- `0006:661` prose: "the atom (`Key`, `Operator`, `Literal`, `Block`)" — reads
  as a grouped unit.

So "follow 0006" does not by itself pick flat or nested; 0006 is internally
inconsistent and a repair is needed either way.

### 0006's sort key forbids flattening into the message

`0006::Normative Contracts`: findings sort "by model id, invariant code, source
rule/context id or graph element id, then normalized predicate/write
fingerprint." Under a five-field record, `model` and `severity` have no field
and `rule` fuses into `locator` — the sort key must be re-parsed out of a
string, and collides whenever a rule id contains the delimiter
(`../spikes/finding-shape-options.out`).

### Prior art (langref checkouts)

| System | Anchor | Producer-specific data carried by |
| --- | --- | --- |
| golangci-lint `Issue` (~100 producers) | `golangci-lint/pkg/result/issue.go:16` | **(a) flat union**, `omitempty`. Faced this exact case with nolintlint and **promoted** `ExpectNoLint`/`ExpectedNoLintLinter` flat onto the shared record |
| beads `Issue` (~30 fields) | `beads/internal/types/types.go:19` | **(a) flat union**, banner-grouped, rigorous `omitempty` reasoning; provider-specific types are separate structs mapped into the core |
| k8s `field.Error` (internal) | `.../validation/field/errors.go:31` | **(c) untyped slot** (`BadValue interface{}`) + `Origin` discriminator |
| k8s `metav1.StatusCause` (wire) | `.../meta/v1/types.go:1057` | narrow flat core `{Type, Message, Field}`, all `omitempty`; `Status.Details` is a typed nested slot **on the envelope, not per-cause** |
| k8s `NewInvalid` | `.../api/errors/errors.go:284` | proves the lossy narrowing: `BadValue`/`Origin` **dropped**, `BadValue` survives only flattened into `Message` |
| helm lint `Message` | `helm/pkg/chart/v2/lint/support/message.go:45` | **(d) flattened into the message string** — total machine-readability loss |
| roborev `ErrorDetail` (RFC 7807) | `roborev/pkg/client/generated/types.go:396` | **(c) untyped bag** (`Value *struct{}`) |
| SARIF via golangci-lint | `golangci-lint/pkg/printers/sarif.go:98` | flat core + nested `locations[]`; drops `SuggestedFixes`/`LineRange` |
| gh-cli | `gh-cli/pkg/cmd/api/api.go:680` | parse-and-flatten; GitHub's `resource`/`field`/`code` **discarded** |
| goreleaser | — | `negative:` no shared multi-check finding record exists |

**Verdict from prior art.** Flat union with `omitempty` dominates *internal*
records; a narrow flat core dominates *wire* records; the observed pattern is
two records with an explicit lossy narrowing between them. **No system nests a
typed per-producer sub-object on a finding.** golangci-lint hit the identical
situation and chose flat promotion. k8s nests typed detail on the *envelope*,
one-per-response keyed by reason — not per-item — and still flattens into the
message for generic consumers.

**Load-bearing lesson.** Every system that kept structured per-producer detail
*also* kept a self-sufficient human-readable message beside it. Whatever shape
is chosen, the core message must stand alone without the structured fields.

## OPEN-2 — set-member JSON escaping

### The standard is decisive and points away from `json.Marshal`

**RFC 8785 (JSON Canonicalization Scheme) §3.2.2.2**: "If the Unicode value is
outside of the ASCII control character range, it MUST be serialized 'as is'
unless it is equivalent to U+005C (`\`) or U+0022 (`"`)". Escaping is mandated
only for control chars U+0000–U+001F, backslash, and quote. `<`, `>`, `&` MUST
be serialized as-is; `\uXXXX` for non-control characters is **prohibited**, not
optional. (Not in any arc corpus — `negative: RFC 8785 / 8259 / ECMA-404 /
OLPC Canonical JSON absent from DevRef, PapersFast, StateMachine*`; read from
rfc-editor.org.)

### The divergence surface is exactly three characters

Verified in `../spikes/escaping-surfaces.out`: only `<`, `>`, `&` differ between
`json.Marshal` and `SetEscapeHTML(false)`. U+2028/U+2029, quote, backslash, tab,
newline, NUL, non-ASCII, and emoji are **byte-identical under both**. So
`SetEscapeHTML(false)` is minimal over ASCII but not strictly JCS-conformant —
a residual deviation worth stating if JCS is ever claimed.

### Correctness is not at stake; stability and legibility are

Both forms are valid RFC 8259 and unmarshal to the identical Go string — the
escaped and unescaped spellings denote the same character. What breaks is
**byte-equality when two call sites disagree**: `resolve` emits one spelling,
`set-state` re-renders the other, and read-back reports
`flow-write-readback-mismatch` on a write that actually succeeded. The failure
is data-dependent, so ASCII-clean fixtures pass and it surfaces only in the
field. KERI (Smith 2019, arXiv, p.128) names this class "the canonicalization
problem": "round trip serializations and deserializations may not be identical
across all implementations of JSON serializers."

### Go ecosystem practice on canonical / machine-readable paths

| System | Anchor | Choice |
| --- | --- | --- |
| gh-cli `--json` (primary output funnel) | `gh-cli/pkg/cmdutil/json_flags.go:228` | `SetEscapeHTML(false)` |
| gh-cli `jsoncolor` | `gh-cli/pkg/jsoncolor/jsoncolor.go:125` | "works like json.Marshal but with HTML-escaping disabled" |
| beads `marshalCanonical` | `beads/cmd/bd/protocol/corpus.go:343` | no HTML-escaping, "so the bytes are stable and minimal" |
| beads `marshalIndentNoEscape` | `beads/cmd/bd/metrics.go:237` | "so `<`, `>`, and `&` render as themselves in human-facing output" |
| prometheus golden fixtures | `prometheus/cmd/prometheus/features_test.go:89` | disabled for byte comparison |
| helm, kubebuilder, goreleaser, roborev, semble | — | zero uses; they emit YAML/templates or have no byte-equality contract |

Go's own stdlib frames it: `json.Marshal` escapes "so that the JSON will be safe
to embed inside HTML `<script>` tags"; `SetEscapeHTML` doc: "In non-HTML
settings where the escaping interferes with the readability of the output".
This output is never embedded in a `<script>` tag.

### Surface-by-surface cost (`../spikes/escaping-surfaces.out`)

The escaped form's damage compounds where the value is stored: a set written to
a TOML artifact becomes `labels = "[\"a\\u003cb\",…]"` — double-escaped, and a
round trip rewrites the author's source with a semantically-null diff. In
`--as=text` a user cannot recognize, grep, or copy back their own value. It also
changes 0006's finding sort fingerprint, since `\` (0x5c) and `<` (0x3c) order
differently.

### Counter-argument, stated fairly

Byte-equality does not *require* the unescaped form — either policy is
internally consistent, and `json.Marshal` is the shorter call with no buffer, no
encoder, and no trailing-newline quirk (`Encoder.Encode` appends `\n` the helper
must strip). "Always `json.Marshal`, never `Encoder`" is mechanically easier to
enforce than "always route through the helper". That argument loses only on
JCS conformance, TOML source corruption, and `--as=text` legibility — but those
are three of the four axes, and the escaping buys nothing here.

Note the cost is not free: adopting the unescaped form means both shipped output
paths (`internal/cli/clierr/clierr.go:135`, `internal/cli/respond/respond.go:183`)
change, and that touches RDR 0006's envelope, which is Final.

---

## Dispositions — author's round, 2026-08-24

Both items approved on the recommendations above.

**OPEN-1 → `Finding` is one flat record with `omitempty` fields.** Fields:
`Code`, `Message` (both always present), plus `omitempty` `Param`, `Locator`,
`Hint`, `Severity`, `Model`, `Rule`, `Key`, `Operator`, `Literal`, `Block`,
`Class`. Each producer populates the subset it owns; the guard atom's four
fields sit flat, not in a nested `atom` object. Carried obligation from the
prior art: `Finding.message` MUST be self-sufficient — readable with no
structured field consulted (k8s `NewInvalid` precedent).

*Consequence for RDR 0006:* its `Technical Design` ("`Key`, `Operator`,
`Literal`, `Block`, and `Class` are `string`") **agrees** with this shape. Its
illustrative JSON showing a nested `"atom": {…}` does not, and is the side
needing a citation repair — illustrative, not normative, so no reopening of a
Final RDR is required.

**OPEN-2 → HTML escaping disabled on the canonical set literal.** `<`, `>`, `&`
serialize as themselves. Normative fixtures approved: `["cli","final"]` (common
case, identical under either encoder) and `["a<b","p>q","x&y"]` (divergent
case). Recorded deviation: not strictly JCS-conformant, since U+2028/U+2029
escape under both encoders — acceptable because JCS conformance is not claimed.

*Consequence for implementation:* `internal/cli/clierr::EmitJSON` and
`internal/cli/respond::writeJSONLine` both move off bare `json.Marshal` onto one
shared helper. Bare `json.Marshal` must not remain on any path a set value
crosses. Booked as a Prerequisite. MVV scenario 5 now asserts a set member
containing `<` and one containing `&` round-tripping byte-identically, so the
rule is covered by validation rather than review; new scenario 7 covers the
`Finding` record's `omitempty` and self-sufficient-message rules.
