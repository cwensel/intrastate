model: claude-opus-5[1m]
variant: full (profile: foundational)

# Repeatability run-1 — RDR 0008 reconstruction

Reconstructed from `docs/rdr/0008-recognized-tag-key-ownership.md` only.
Guesses where the RDR is silent are marked **GUESS**.

## 1. Public API

The RDR names one required new export and otherwise constrains values, not
structs. Package `internal/resolve` (the RDR 0001 kernel).

```go
package resolve

// Reserved recognized-provenance tag key. The RDR pins the *value* and cites
// the existing unexported constant as its home; it explicitly leaves exporting
// the spelling as implementation latitude ("Whether the normalizer additionally
// references an exported kernel constant for the spelling is implementation
// latitude").
const recognizedTagKey = "recognized" // existing, unexported; RDR 0008 Phase 1 adds a pointer comment

// GUESS — the RDR does NOT require this to be exported, and says the normalizer
// must not depend on a literal-vs-literal comparison. Included only as the
// latitude the RDR permits; a conforming implementation may omit it.
// const RecognizedTagKey = recognizedTagKey

// REQUIRED EXPORT (block 4 + block 5, Enforcement locus).
// "the kernel MUST export a construction-time predicate over `Input` that
// producers may call, and `Resolve` MUST apply that same predicate at entry".
//
// GUESS — name. The RDR makes the exported name implementation latitude and
// explicitly forbids binding it to any RDR 0009 symbol name. `ValidateInput` is
// this run's pick.
// GUESS — signature. The RDR says only that it is a "predicate" over `Input`
// that returns/produces a breach on the Go error path ("returning a non-nil
// error and no `Result` disposition on breach"). A bool-returning "predicate"
// and an error-returning validator are both readings; error-returning is chosen
// because the RDR requires the breach to *travel* the Go error path and to name
// the offending key, which a bool cannot carry.
func ValidateInput(in Input) error

// Existing entry point, unchanged in signature.
// RDR: "`Resolve` MUST apply that same predicate at entry, returning a non-nil
// error and no `Result` disposition on breach." Also: "`internal/resolve/resolve.go::Resolve`
// already returns `(Result, error)` whose error is doc-reserved for programmer
// mistakes with no non-nil path in its body today, so (b) adds a check without
// widening the signature."
func Resolve(in Input) (Result, error)
```

### Error modes

| Mode | Shape | RDR basis |
| --- | --- | --- |
| Reserved key in `Input.Owned` or `Input.Observed` | non-nil `error` from `Resolve`; zero-valued `Result` (no `Plan`, no `Refusal`) | block 4; scenario 6 |
| Reserved key in any row's `Row.RequiresOwned` | same non-nil `error` path, same predicate | block 5; scenario 9 second half |
| Conforming `Input` | nil error; existing dispositions unchanged | block 6; scenario 6's nil-error pin, incl. empty `Input{}` |
| Bypassing caller (package-internal path) reaching `assemble` with a reserved-key owned tag | D3 precedence resolves deterministically; no error | block 6; scenario 8 |
| Bypassing caller with `RequiresOwned: ["recognized"]` | `owned_state_unavailable` refusal, `MissingOwned` contains `recognized` | block 5 residual; scenario 9 first half |

**No new `RefusalKind`.** Stated three times in the RDR (blocks 4, 5, 6).

**GUESS — error type.** The RDR says "travels the Go error path RDR 0001
reserves for programmer mistakes" but never names a type. This run assumes a
plain `fmt.Errorf`/`errors.New`-shaped error naming the offending key and
channel, not a typed sentinel. The RDR does not require callers to discriminate
programmatically, so no `errors.Is` target is specified.

**GUESS — whether `ValidateInput` reports one breach or all.** Block 4 says
"detected whenever present — never skipped because another precondition also
fired", which constrains *this* predicate vs RDR 0009's, not multiple breaches
*within* this predicate. This run returns the first breach found.

### Normalizer-side API (RDR 0002's package; this RDR only constrains values)

The RDR is explicit that "The Go type, package, and field names carrying these
values are RDR 0002's to choose — this RDR constrains the values and their
distinctness, not the struct." So the only fixed API here is the value set:

```
category discriminator (data-level): "reserved_tag_key"
rule identifier (inside the failure payload): "reserved-tag-key/kernel-owned"
required-name field value: "recognized"
offending-name field value: <the key as authored>
```

**GUESS — the advisory's carrier.** Load-Bearing Decisions / Identity requires
a "non-blocking advisory naming both spellings", and scenario 4 asserts it, but
the RDR gives it no category token, no rule id, and no field list. This run
assumes a separate advisory/diagnostic channel on the normalizer result
(distinct from the failure list, since it "MUST NOT change the load/lint
verdict"), carrying the authored spelling and the reserved spelling. **This is
the largest under-specified surface in the RDR.**

## 2. Three most important internal helpers

1. **`assemble(in Input) TagSet`** — *existing, unchanged.* Sole constructor of
   the assembled evaluation view; writes exactly three sources (`in.Recognized`,
   `in.Observed`, `in.Owned`, cited at `:151-162`) under D3 precedence
   (`owned` > `observed` > `recognized`). Injects `in.Recognized` under
   `recognizedTagKey` with `ProvenanceRecognized` **only when `in.Recognized` is
   non-empty**. RDR 0008 changes nothing here.

2. **the reserved-key scan behind `ValidateInput`** — walks `in.Owned` keys,
   `in.Observed` keys, and every `in.Table.Rows[i].RequiresOwned` entry, testing
   byte-exact equality against `recognizedTagKey`. Unconditional on
   `in.Recognized` (block 4: "The check is unconditional on `Input.Recognized`").
   **GUESS — decomposition.** The RDR requires "one predicate, two call sites"
   for the *exported* surface; whether the internals split into per-channel
   helpers is unstated. This run keeps one scan.
   **GUESS — the RDR never states that `Input` carries the table.** Block 5 says
   the predicate checks "the rows of the supplied table" and notes this "widens
   the `Input` predicate to read `in.Table.Rows`", so `in.Table.Rows` is taken as
   the accessor path. The field name `Table` is inferred from that phrase.

3. **`missingOwned(row Row, view TagSet) []string`** — *existing, unchanged.*
   Iterates `row.RequiresOwned` and tests `view.has(key, ProvenanceOwned)`; a key
   present under `ProvenanceRecognized` fails the owned-only test and is reported
   missing, yielding `owned_state_unavailable`. RDR 0008 cites this as the
   kernel-side residual it narrows (not removes).

4. *(normalizer side, RDR 0002's code)* **the declaration-key reserved check** —
   over `[tags.<tag>]` post-parse key strings only. Two rules: provenance
   `recognized` ⇒ name MUST be `recognized`; provenance `owned`/`observed` ⇒ name
   MUST NOT be `recognized`. Plus the near-miss advisory (case-fold + trim).
   Explicitly **not** applied to `[rule.match.*]`, `[rule.guard.all.*]`,
   `[rule.guard.unless.*]` (covered by `unknown tag` fall-through) or
   `[rule.write]` (covered by `write to non-owned tag`).

## 3. Data model across the boundary

### Kernel `Input` (RDR 0001's, read by the new predicate)

```go
type Input struct {
    Recognized string            // at most one per resolve (singleton by inheritance)
    Observed   map[string]string // GUESS — shape. RDR says only "tag keys"; map assumed
    Owned      map[string]string // GUESS — shape, same
    Table      Table             // GUESS — field name; inferred from "in.Table.Rows"
}
```

**GUESS — `Owned`/`Observed` element type.** The RDR never states whether a tag
value is a string, a typed value, or a struct. Only the *key* matters to this
RDR, so the value type is invented here.

### `Row` (RDR 0001's / RDR 0002's normalized output)

```go
type Row struct {
    Outcome       string   // outcome gate: `row.Outcome != in.Recognized`
    Match         ...      // tag predicates; GUESS — shape unstated
    RequiresOwned []string // derived; no authored TOML form (no `requires_owned` key)
    Escape, Writes ...     // RDR 0009's subject; named here only as disjoint
}
```

### `TagSet` (the assembled view)

Two exported observables, per scenario 3: `Lookup(key) (value, provenance, ok)`
and `Len()`. Internal `has(key, prov)` compares `tv.provenance == prov`.
Provenance vocabulary: `ProvenanceOwned`, `ProvenanceObserved`,
`ProvenanceRecognized`.

### Validation failure payload (RDR 0002's struct, this RDR's values)

Three distinct machine-readable fields, all asserted byte-for-byte by scenario 7:

| Meaning | Value |
| --- | --- |
| category discriminator | `reserved_tag_key` |
| offending name | as authored |
| required name | `recognized` |
| rule identifier | `reserved-tag-key/kernel-owned` |

The RDR is emphatic that category and rule id "are not alternatives and a
consumer MUST NOT choose between them": the category is what RDR 0005's
exit-code map keys on; the rule id is for golden assertions and remediation
lookup, never dispatch.

The same three-field payload requirement **extends RDR 0002's pre-existing
`unknown tag` failure** where the offending key is the reserved name.

### TOML source shape (RDR 0002's schema, constrained here)

```toml
outcomes = [...]              # root; closed alphabet
[model]
[tags.recognized]             # MUST be this name iff provenance = "recognized"
provenance = "recognized"
[rule.match.<tag>]
[rule.guard.all.<tag>]
[rule.guard.unless.<tag>]
[rule.write]
```

Key comparison is byte-exact on the **post-parse** key string: case-sensitive,
no trimming, no folding. `[tags."recognized"]` == `[tags.recognized]` (quoting
is a surface artifact). `Recognized`, `RECOGNIZED`, `" recognized"` are ordinary
unreserved names that raise the advisory.

Scope unit for cardinality: one `[model]`. No cross-model or cross-file rule.

## 4. Pseudo-code of the main operation

```
func Resolve(in Input) (Result, error):
    # RDR 0008 block 4: the ONLY new kernel step. Entry precondition.
    if err := ValidateInput(in); err != nil:
        return Result{}, err          # zero Result: no Plan, no Refusal, no RefusalKind

    # --- everything below is RDR 0001's existing kernel, unchanged ---
    view := assemble(in)              # sole view constructor; D3 precedence inside
    ...
    for each row in in.Table.Rows:
        if row.Outcome != in.Recognized: continue      # outcome gate; never reads the tag key
        if !in.Table.models(row.Outcome): ...          # alphabet gate, slices.Contains
        if !matches(row.Match, view): continue
        if !guards(row, view): continue                # same view as the matcher (A2)
        if missing := missingOwned(row, view); len(missing) > 0:
            return refusal(owned_state_unavailable, MissingOwned: missing), nil
        return plan(row), nil
    return refusal(no_match), nil

func ValidateInput(in Input) error:
    # Unconditional on in.Recognized — the reservation does not lapse per-call.
    for k := range in.Owned:
        if k == recognizedTagKey:
            return fmt.Errorf(...)    # GUESS: message/type unspecified by the RDR
    for k := range in.Observed:
        if k == recognizedTagKey:
            return fmt.Errorf(...)
    # block 5: same predicate covers the owned-reference channel
    for _, row := range in.Table.Rows:
        for _, k := range row.RequiresOwned:
            if k == recognizedTagKey:
                return fmt.Errorf(...)
    return nil
    # NOTE: RDR 0009 writes a structurally identical precondition over
    # Row.Escape/Row.Writes on this same entry point. Ordering between the two
    # is A11 / Stage 7.1. Interim rule: apply this one and report its breach; if
    # the table also breaches 0009's shape rule, either error is conforming.
    # GUESS — call order relative to 0009's predicate is deliberately unfixed.

# Normalizer side (RDR 0002's code, this RDR's rule)
func validateTagDeclarations(decls) (failures, advisories):
    for name, d := range decls:          # name = post-parse key string
        if d.provenance == "recognized" and name != "recognized":
            failures += reservedTagKey(offending: name, required: "recognized",
                                       rule: "reserved-tag-key/kernel-owned")
        if d.provenance in {"owned","observed"} and name == "recognized":
            failures += reservedTagKey(offending: name, required: "recognized",
                                       rule: "reserved-tag-key/kernel-owned")
        if name != "recognized" and fold(trim(name)) == "recognized":
            advisories += nearMiss(authored: name, reserved: "recognized")   # non-blocking
    return
    # GUESS — the owned/observed-named-`recognized` case reuses the same
    # `required` field value ("recognized"), which reads oddly since the author
    # must RENAME AWAY from it there. The RDR mandates three fields for "every
    # reserved_tag_key failure" without splitting by direction. Flagged.
```

## Determinacy notes — where the RDR left me guessing

Ranked by how load-bearing the silence is:

1. **The advisory's data shape** (Load-Bearing Decisions / Identity, scenario 4)
   — normative and tested, but no category token, no rule id, no field list, no
   carrier. Contrast with the failure payload, which is pinned to the byte.
2. **`ValidateInput`'s signature** — "predicate" (bool) vs. error-returning
   validator. The RDR uses "predicate" throughout but requires a non-nil error
   and an offending-key-naming breach.
3. **`Input`'s table accessor** — `in.Table.Rows` is quoted once in block 5;
   the field is never introduced in the API description.
4. **The `required` field's meaning in the owned/observed direction** — the
   three-field payload is mandated uniformly, but "required name = recognized"
   is wrong advice for a tag that must be renamed away from `recognized`.
5. **`Owned`/`Observed` element types** — irrelevant to the rule, but a
   reconstructor must invent them.
6. **Breach-reporting multiplicity** — first breach vs. all breaches.
7. **Ordering vs RDR 0009's predicate** — deliberately left open (A11), and the
   RDR says so explicitly. Not a defect; recorded so the diff does not count it
   as one.
