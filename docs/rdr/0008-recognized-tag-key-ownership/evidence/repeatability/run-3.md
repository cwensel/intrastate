model: glm-5.2:cloud
variant: full (profile: foundational)

# Repeatability run-3 — RDR 0008 reconstruction

Reconstructed from the RDR at
`docs/rdr/0008-recognized-tag-key-ownership.md` only. Where the RDR
is silent I made a best guess and marked it `GUESS`. The RDR's
enforcement spans two surfaces — the `internal/resolve` kernel
(RDR 0001, implemented) and the RDR 0002 normalizer (Final,
unimplemented) — so the "module" below is the union of the
contract this RDR authors for each side.

## 1. Public API (function signatures, types, error modes)

### Kernel surface — `internal/resolve`

```go
// Existing (RDR 0001), unchanged signature. Now applies the
// reserved-key precondition at entry.
func Resolve(in Input) (Result, error)
```
- On a reserved-key breach of `in`: returns a non-nil `error` and
  a zero-valued `Result` (no `Plan`, no `Refusal` disposition).
  This is the kernel's first non-nil error path; it travels the
  Go error path RDR 0001 REQ-6 reserves for *programmer mistakes*,
  not a new `RefusalKind`.
- On conforming `in` (including the empty `Input{}` already pinned
  at `resolve_test.go:748`): unchanged behavior, nil error.

```go
// New, exported by the kernel package. One predicate, two call
// sites: producers call it at construction; Resolve calls it at
// entry. Carries BOTH reserved-key obligations — the Input
// owned/observed keys AND each row's RequiresOwned reservation.
func <ExportedPredicateName>(in Input) error
```
- The exported **name is implementation latitude** — the RDR
  forbids binding it to any symbol name from RDR 0009 (still
  Draft). `GUESS: func ValidateInput(in Input) error` — the
  RDR names neither the symbol nor its package path.
- Breach: non-nil `error`, same programmer-mistake semantics as
  `Resolve`. Conforming: nil.
- It reads `in.Owned`, `in.Observed`, and `in.Table.Rows[*].RequiresOwned`
  (block 5 widens it onto the rows, so it is not decidable without
  the table). `GUESS`: returns a single aggregate error for the
  first breach; the RDR does not specify whether one error or a
  list is returned.

```go
// Existing, unexported constant (RDR 0001 D2). This RDR ratifies
// it from implementation latitude to normative contract.
const recognizedTagKey = "recognized"
```
- Whether this is **exported** for the normalizer to reference is
  implementation latitude — the RDR requires only a behavioral
  conformance test, and explicitly leaves "whether the normalizer
  additionally references an exported kernel constant for the
  spelling" to the implementer. `GUESS`: left unexported, matching
  D2's "no public surface" posture; the normalizer hardcodes the
  literal under the behavioral test.

### Types

```go
type Input struct {
    Recognized string          // the freshly recognized outcome; singleton per resolve (RDR 0001)
    Owned    map[string]TaggedValue  // GUESS: accessor-produced owned tag snapshot; key = tag name
    Observed map[string]TaggedValue  // GUESS: caller-supplied observed context; key = tag name
    Table    Table                      // carries Rows
}
```
- `Owned`/`Observed` map shapes are `GUESS` — the RDR cites the
  package doc ("an accessor-produced owned tag snapshot") but
  not the concrete Go type. Whether they are maps, slices, or a
  `TagSet`-like wrapper is not stated. `GUESS: map[string]TaggedValue`.

```go
type Result struct {
    Plan    *Plan      // GUESS: field name; the RDR says "no Plan" on breach
    Refusal Refusal    // GUESS: field name; the RDR says "no Refusal" on breach
}
```
- Field names are `GUESS`. The RDR only asserts the zero value
  carries "no `Plan`, no `Refusal`."

```go
type Provenance int          // GUESS: underlying type
const (
    ProvenanceOwned     Provenance = iota // GUESS ordering
    ProvenanceObserved
    ProvenanceRecognized
)
```
- The three provenance values are named in the RDR; their Go
  representation and ordering are `GUESS`.

```go
// Assembled evaluation view. Exports the two observables the
// conformance test reads (scenario 3).
func (t TagSet) Lookup(key string) (value string, prov Provenance, ok bool)  // GUESS: exact arity
func (t TagSet) Len() int
```
- The RDR writes the test as
  `Lookup("recognized") == (in.Recognized, ProvenanceRecognized, true)`,
  implying a 3-tuple return `(value, provenance, found)`. The
  value type (`string` vs a `TaggedValue` struct) is `GUESS` —
  the equality is written against `in.Recognized` (a string), so
  `GUESS: value string`. `Len()` is named by scenario 3.
- `TagSet` wraps a map (pass-by-value shares one backing map, per
  A2), so `Lookup` reads the shared backing map.

```go
type Row struct {
    Outcome      string         // first-class non-tag gate; unaffected by this RDR
    Match        []KeyVal       // GUESS: tag pattern over the assembled view
    Guard        GuardSpec      // GUESS: guard predicate spec; RDR 0007 owns its domain
    RequiresOwned []string      // derived normalizer output; no authored TOML form
    Escape       []KeyVal       // GUESS: RDR 0009's field
    Writes       []KeyVal       // GUESS: RDR 0009's field
}
```
- `RequiresOwned []string` is the one field the RDR pins by name
  and type (block 5, scenario 9: `RequiresOwned: []string{"recognized"}`).
  The others' shapes are `GUESS`.

### Normalizer surface — RDR 0002 (this RDR authors the contract)

```go
// New additive data-level validation category. The Go type,
// package, and field names are RDR 0002's to choose — this RDR
// constrains the VALUES and their distinctness, not the struct.
type ReservedTagKeyFailure struct {
    OffendingName string // as authored (post-parse key string)
    RequiredName   string // literal "recognized"
    RuleID         string // literal "reserved-tag-key/kernel-owned"
}
```
- The category discriminator value is the token `reserved_tag_key`
  (snake_case, matching RDR 0001's `RefusalKind` grammar). The
  spaced form `reserved tag key` is prose for this document only.
- The rule identifier is `reserved-tag-key/kernel-owned`, a
  comparable token asserted byte-for-byte by a golden test.
- `reserved_tag_key` (category) and `reserved-tag-key/kernel-owned`
  (rule id) are **not alternatives** — a consumer MUST NOT choose
  between them. The category is what RDR 0005's exit-code map
  keys on; the rule id is for golden assertions / remediation
  lookup inside the payload.
- Struct shape is `GUESS` (RDR 0002's to choose); only the three
  fields and their literal values are fixed.

```go
// Non-blocking advisory for near-miss spellings. NOT a failure;
// MUST NOT change the load/lint verdict.
type NearMissAdvisory struct {
    AuthoredName string // the near-miss spelling as authored
    ReservedName string // "recognized"
}
```
- Fires when the post-parse key is not `recognized` but becomes
  `recognized` under Unicode-simple case folding **and** trimming
  of leading/trailing whitespace. Both conditions must hold
  (case-folded-only or trimmed-only is `GUESS` — the RDR says
  folding *and* trimming, conjoined; I read that as both
  required for the advisory).
- The advisory's surface type and how it is carried alongside a
  clean verdict is `GUESS` — the RDR fixes the trigger and the
  "non-blocking" property, not the data shape.

### Error modes (summary)

| Origin | Channel | Shape |
|---|---|---|
| Declaration named `recognized` under provenance `owned`/`observed` | RDR 0002 load/lint | `reserved_tag_key` data-level failure (not a Go error, not a `RefusalKind`) |
| Recognized-provenance declaration under any other name | RDR 0002 load/lint | `reserved_tag_key` data-level failure |
| Near-miss (`Recognized`, `" recognized"`) declaration | RDR 0002 load/lint | non-blocking advisory; verdict unchanged |
| `Input.Owned`/`Input.Observed` keyed `recognized` | kernel Go error path | non-nil `error`, zero `Result` (producer programmer mistake) |
| `Row.RequiresOwned` naming `recognized` (via predicate) | kernel Go error path | non-nil `error`, zero `Result` |
| `Row.RequiresOwned` naming `recognized` (bypassing predicate, package-internal path) | kernel refusal | `owned_state_unavailable` with `MissingOwned` containing `recognized` |
| Conforming input | unchanged | nil error, normal disposition |

## 2. Three most important internal helper functions

### `assemble(in Input) TagSet` — the sole view constructor (existing, RDR 0001)
Constructs the assembled evaluation view that both matchers and
guards read in one resolve. Injects `recognizedTagKey` with
`ProvenanceRecognized` only when `in.Recognized` is non-empty
(the sole injection point — A2's whole-package sweep). Then
merges `in.Observed` and `in.Owned` in D3 precedence order
(`owned` > `observed` > `recognized`), so an owned/observed key
`recognized` would overwrite the recognized entry — the
shadowing this RDR's producer obligation exists to prevent
reaching. Unchanged by this RDR; this RDR guarantees conforming
input never presents that overwrite.

### `<ExportedPredicateName>(in Input) error` — the reserved-key predicate (new)
The single definition behind "one predicate, two call sites."
Traverses `in.Owned` and `in.Observed` keys (byte-exact compare
to `"recognized"`) and each row in `in.Table.Rows`'s
`RequiresOwned` (block 5 widening). Returns non-nil on the first
breach. The check is **unconditional on `in.Recognized`** — a
reserved-keyed owned/observed tag is a breach whether or not the
resolve carries an outcome (block 4's unconditional clause),
because a key whose reservation lapses per-call is not reserved.
`GUESS`: the internal decomposition — whether owned/observed and
RequiresOwned are checked in one loop or two, and which breach
is reported first when both an `Input` key and a `RequiresOwned`
entry offend simultaneously, is not specified. `GUESS`: reports
the `Input`-key breach before the `RequiresOwned` breach, but
the RDR does not say.

### `validateTagDeclarations(model Model) []Failure` — the normalizer name check (new, lands in RDR 0002's normalizer)
Iterates `[tags.<tag>]` declaration keys (post-parse key string,
byte-exact). For each: a recognized-provenance declaration whose
key is not `"recognized"` → `reserved_tag_key` failure; an
owned/observed-provenance declaration whose key is `"recognized"`
→ `reserved_tag_key` failure. Also emits the near-miss advisory
for keys that fold+trim to `"recognized"` but are not byte-equal.
Predicate-position references (`[rule.match.<tag>]`,
`[rule.guard.all.<tag>]`, `[rule.guard.unless.<tag>]`) need no
separate check — they resolve to a declared tag, and an
undeclared `recognized` already fails RDR 0002's `unknown tag`
category. `[rule.write]` is covered by RDR 0002's pre-existing
`write to non-owned tag` category, not by this check.
`GUESS`: the function name and where it sits in 0002's validation
order relative to `unknown tag` / `malformed TOML` / `ambiguous
overlap` — the RDR does not fix the ordering among lint checks.
`GUESS`: the advisory is emitted from this same pass; whether it
is a separate return channel or a field on the result is
unspecified.

## 3. Data model (persisted / passed across the boundary)

### TOML declaration grammar (authored, RDR 0002 surface)

```toml
outcomes = ["round-clean", "round-faulty"]   # root outcome alphabet; closed set

[model]                                       # schema unit; scope of cardinality
...

[tags.<tag>]                                  # <tag> = post-parse key string
provenance = "owned" | "observed" | "recognized"

[rule.match.<tag>]     eq = "..."             # tag-predicate family
[rule.guard.all.<tag>]  ...
[rule.guard.unless.<tag>]  ...
[rule.write]            <tag> = ...            # write position; governed by `write to non-owned tag`
```

- A `[tags.<tag>]` whose key is `recognized` and provenance
  `recognized` is the one conforming recognized-provenance
  declaration. Cardinality ≤ 1 is a consequence of name-keying
  (`[tags.<tag>]` is keyed by name), not a separate check. No
  lower bound — a model declaring no recognized-provenance tag is
  conforming (the affordance is not an obligation).
- TOML quoting is a surface artifact: `[tags."recognized"]` and
  `[tags.recognized]` parse to the identical key string → both
  reserved.

### Kernel `Input` (passed across the producer→kernel boundary)

- `Recognized string` — singleton per resolve (RDR 0001 `Input`
  carries exactly one).
- `Owned`/`Observed` — producer-supplied tag snapshots keyed by
  tag name. Producers MUST NOT supply a key `"recognized"` here;
  the reserved key enters the view only through `Input.Recognized`.
- `Table.Rows[*].RequiresOwned []string` — derived normalizer
  output, no authored TOML form. MUST NOT name `recognized`.

### Failure payload (passed across the normalizer→consumer boundary)

Three machine-readable fields, asserted byte-for-byte
(scenario 7):
1. `OffendingName` — the offending name as authored (post-parse
   key string).
2. `RequiredName` — the literal `recognized`.
3. `RuleID` — the literal `reserved-tag-key/kernel-owned`.

The category discriminator `reserved_tag_key` participates in
RDR 0002's data-level category set; the rule id does not
participate in category dispatch. Where the offending key is the
reserved name, this payload requirement also extends RDR 0002's
pre-existing `unknown tag` failure (authored by this RDR, lands
in 0002's implementation).

### Advisory payload (near-miss)

Carries both spellings (authored + reserved `recognized`); is
non-blocking and MUST NOT change the load/lint verdict. `GUESS`:
carried as a separate advisory list alongside the failures, not
as a failure variant.

## 4. Top-level pseudo-code — `Resolve` with reserved-key enforcement

```
func Resolve(in Input) (Result, error):
    # (1) Entry precondition — this RDR's reserved-key predicate.
    #     One predicate, two call sites: this call + producers call
    #     the same exported function at construction. Unconditional
    #     on in.Recognized.
    if err := <ExportedPredicateName>(in); err != nil:
        return Result{}, err          # zero Result, non-nil error; programmer mistake

    # GUESS: ordering vs RDR 0009's row-shape precondition (A11,
    # reconciled at Stage 7.1). This RDR requires only that its own
    # breach be detected whenever present — never skipped because
    # 0009's also fired. Below: 0008 checked first, but the order
    # is unsettled and a doubly-breaching caller seeing either
    # error is conforming until 7.1 fixes it.
    #   if err := <RDR0009Precondition>(in); err != nil: return Result{}, err

    # (2) Build the assembled view (existing RDR 0001 path).
    view := assemble(in)
    #   - if in.Recognized != "":  view[recognizedTagKey] = {in.Recognized, ProvenanceRecognized}
    #   - merge in.Observed then in.Owned (D3: owned > observed > recognized)

    # (3) For each row: outcome gate, escape/refuse, match, guard, owned deps.
    for row := range in.Table.Rows:
        if row.Outcome != "" and row.Outcome != in.Recognized:
            continue                       # outcome gate (A2: never reads the key)
        if esc := escapeOrRefuse(row, in): # GUESS: existing RDR 0001 helper shape
            return refuse(esc, in), nil

        if !view.matches(row.Match):        # match pattern over the same view
            continue
        if !gate(row.Guard, view):          # guards read the SAME view (A2)
            continue
        if missing := missingOwned(row.RequiresOwned, view); len(missing) > 0:
            return refuse(owned_state_unavailable{MissingOwned: missing}, in), nil

        return Result{Plan: planFrom(row)}, nil   # GUESS: Plan construction

    return refuse(no_match, in), nil        # GUESS: no_match shape; RDR 0001 owns it
```

`GUESS` markers cluster on:
- the exported predicate's **name** and **error arity** (single
  vs list),
- the **ordering** of 0008's vs 0009's entry preconditions (A11,
  explicitly unsettled),
- the **ordering** of reserved-key checks among RDR 0002's
  existing lint categories,
- the concrete **Go types** for `Owned`/`Observed`, `Result`
  field names, `Provenance` representation, the failure/advisory
  **struct shapes**,
- whether `recognizedTagKey` is **exported** for the normalizer,
- the **advisory's data channel** (separate list vs failure
  field),
- the internal **decomposition / first-breach ordering** within
  the predicate, and
- the `escapeOrRefuse` / `planFrom` / `no_match` helper shapes,
  which are RDR 0001's and only sketched here.