model: claude-fable-5
variant: full (profile: foundational)

# Repeatability run-2 — RDR 0008 reconstruction

Reconstructed from `docs/rdr/0008-recognized-tag-key-ownership.md` only.
Guesses where the RDR is silent are marked **GUESS**.

## 1. Public API

Two surfaces: the kernel (`internal/resolve`, exists) and the RDR 0002
normalizer (does not exist; obligations carried there).

### Kernel (`internal/resolve`)

```go
package resolve

// Reserved recognized-provenance tag key. Existing unexported constant;
// Phase 1 adds only a pointer comment citing RDR 0008. The RDR explicitly
// does NOT require exporting a spelling constant ("implementation latitude"),
// and the conformance test must be behavioral, never literal-vs-literal.
const recognizedTagKey = "recognized"

// REQUIRED EXPORT (Normative Contracts blocks 4+5; Enforcement locus
// "(b)+(c) paired — one predicate, two call sites"): a construction-time
// predicate over Input that producers may call. It must reject:
//   - any Input.Owned or Input.Observed tag keyed exactly "recognized"
//     (byte-exact; unconditional on Input.Recognized), and
//   - any row in in.Table.Rows whose RequiresOwned names "recognized"
//     (block 5's widening — the predicate reads the table's rows).
// GUESS — the exported name. The RDR leaves it implementation latitude and
// forbids binding it to any RDR 0009 symbol name. Chosen here:
func ValidateInput(in Input) error

// EXISTING SIGNATURE, new behavior: Resolve MUST apply the same predicate
// at entry. On breach it returns a non-nil error and a zero-valued Result
// (no Plan, no Refusal) — the kernel's first non-nil error path, landing on
// the (Result, error) surface RDR 0001 REQ-6 reserves for programmer
// mistakes. Conforming input (including the empty Input{}) keeps a nil
// error. No new RefusalKind; D3 precedence unchanged.
func Resolve(in Input) (Result, error)
```

Error modes, kernel side:

- Predicate breach → non-nil Go `error`, no `Result` disposition. Never a
  `RefusalKind`, never a modeled refusal.
- **GUESS** — the error's type and text. The RDR constrains the channel and
  the trigger, not the error value; a plain `fmt.Errorf`-style error naming
  the offending key satisfies it.
- Doubly-breaching input (this predicate + RDR 0009's row-shape
  precondition): either error is conforming until Stage 7.1 fixes an order
  (A11). Skipping this check because the other fired is non-conforming.
- Bypass residual (package-internal callers not going through `Resolve`):
  a row with `RequiresOwned: ["recognized"]` refuses
  `owned_state_unavailable` with `MissingOwned` containing `recognized`
  (scenario 9); data-channel collisions fall to D3's deterministic
  precedence (scenario 8).

### Normalizer (carried into RDR 0002's implementation)

No Go signatures are this RDR's to fix — "The Go type, package, and field
names carrying these values are RDR 0002's to choose." The RDR binds
behavior and values:

- New data-level validation category, discriminator token
  `reserved_tag_key` (snake_case, matching RDR 0001's RefusalKind grammar;
  `reserved tag key` is prose only), additive under RDR 0002's "including
  at minimum" fence (A1).
- Checked positions: `[tags.<tag>]` declaration keys **only**. Two rules:
  a `provenance = "recognized"` declaration MUST be named `recognized`;
  an `owned`/`observed` declaration MUST NOT be named `recognized`.
  Predicate positions (`[rule.match.<tag>]`, `[rule.guard.all.<tag>]`,
  `[rule.guard.unless.<tag>]`) need no separate check (`unknown tag`
  covers the undeclared case); `[rule.write] recognized = …` is rejected by
  the pre-existing `write to non-owned tag` category, not by this rule.
- Comparison: byte-exact on the post-parse key string — no folding, no
  trimming. `[tags."recognized"]` ≡ `[tags.recognized]` (both reserved);
  `Recognized`, `RECOGNIZED`, `" recognized"` are ordinary unreserved names.
- Near-miss advisory (normative): a declaration whose post-parse key is not
  `recognized` but becomes `recognized` under Unicode-simple case folding
  *and* leading/trailing whitespace trimming raises a **non-blocking**
  advisory naming both spellings. It never changes the load/lint verdict.
- Failure payload (block 3): every `reserved_tag_key` failure — and any
  `unknown tag` failure whose offending key is the reserved name — carries
  three distinct machine-readable fields: offending name as authored,
  required name (literal `recognized`), rule identifier (literal
  `reserved-tag-key/kernel-owned`).

## 2. Three most important internal helpers

1. **Reserved-key declaration check** (normalizer load/lint pass; GUESS at
   name, e.g. `checkReservedTagKey(decls) []Failure`). Iterates the parsed
   `[tags.<tag>]` map; for each declaration applies the two naming rules
   byte-exactly on the post-parse key and emits `reserved_tag_key` failures
   with the three-field payload. Runs before any resolution, in the same
   pre-acceptance layer as RDR 0002's `unknown tag` / `write to non-owned
   tag` (upstream of RDR 0006's graph lint). Cardinality falls out free:
   a second recognized-provenance declaration necessarily has a different
   name and fails the naming rule (two failures for two wrong-named
   declarations); two `[tags.recognized]` are a `malformed TOML`
   duplicate-key error before this rule is reached.

2. **Near-miss advisory detector** (GUESS at name, e.g.
   `nearMissAdvisory(key string) (Advisory, bool)`). Computes
   `fold(trim(key)) == "recognized"` for keys that are not the reserved
   word; on match emits the non-blocking advisory carrying the authored and
   reserved spellings. Targeted, not blanket: an ordinary name like
   `result` raises nothing (scenario 4). Kept separate from the check in
   (1) because it must not affect the verdict.

3. **The shared predicate core** (kernel; the single definition behind the
   two call sites — the exported predicate *is* this function or wraps it).
   Scans `in.Owned` and `in.Observed` for a key equal to
   `recognizedTagKey`, then walks `in.Table.Rows` scanning each
   `RequiresOwned` for the same key; first hit returns the breach error.
   Unconditional on `in.Recognized`. `assemble` itself is deliberately
   **unchanged** — it stays the sole view constructor, injecting
   `in.Recognized` under the reserved key with `ProvenanceRecognized` only
   when non-empty.

## 3. Data model

Persisted / boundary-crossing values this RDR constrains:

- **Category discriminator**: `reserved_tag_key` — the stable data-level
  value category consumers (RDR 0005's exit-code map) key on. Participates
  in RDR 0002's data-level category set.
- **Rule identifier**: `reserved-tag-key/kernel-owned` — carried *inside*
  the failure payload; for golden byte-for-byte assertions and remediation
  lookup, never for category dispatch. The two tokens are not alternatives;
  both travel, in their respective positions.
- **Failure payload** (shape GUESS — struct/field names are RDR 0002's):
  ```
  { category:  "reserved_tag_key",          // or the unknown-tag category
                                            // when the offending key is the
                                            // reserved name
    offending: "<name as authored>",
    required:  "recognized",
    rule:      "reserved-tag-key/kernel-owned" }
  ```
- **Advisory record** (shape GUESS): non-blocking; carries the authored
  spelling and the reserved spelling. Whether it shares the failure
  payload's type or has its own is unstated — GUESS: its own type, since
  it is not a validation failure and must not alter the verdict.
- **Kernel types**: none added, none changed. `Input`, `Result`, `TagSet`,
  `Row` keep their shapes; no new `RefusalKind`; the assembled view still
  binds the recognized outcome at key `recognized` with
  `ProvenanceRecognized` (empty `Input.Recognized` → no key, and such a
  resolve is outside the binding obligation, not in breach).
- **Migration data** (Phase 2, gated on A9): three fixture declarations
  (`rdr-fixture.toml` / `kata-fixture.toml` `[tags.outcome]`,
  `guard-fixture.toml` `[tags.rewind_target]`) plus all six tag-predicate
  reference sites rename to `recognized`.

## 4. Top-level pseudo-code — the reserved-key enforcement path

```
# --- load/lint time (RDR 0002 normalizer; carried obligation) ---
load_table(toml):
    doc = parse_toml(toml)                       # dup keys → malformed TOML,
    failures, advisories = [], []                #   before any naming rule
    for key, decl in doc.tags:                   # post-parse key strings
        if decl.provenance == "recognized" and key != "recognized":
            failures += reserved_tag_key(offending=key,
                                         required="recognized",
                                         rule="reserved-tag-key/kernel-owned")
        if decl.provenance in {"owned","observed"} and key == "recognized":
            failures += reserved_tag_key(offending=key, ...)
        if key != "recognized" and fold(trim(key)) == "recognized":
            advisories += near_miss(authored=key, reserved="recognized")
    # predicate positions: no reserved-key check — unknown tag / write to
    # non-owned tag (pre-existing categories) cover them; where the
    # offending key is "recognized", unknown-tag failures also carry the
    # three-field payload
    if failures: reject(failures, advisories)    # typed, pre-resolution
    return normalize(doc), advisories            # advisories never block

# --- resolve time (kernel; this RDR's own implementation) ---
ValidateInput(in):                               # GUESS at exported name
    for key in keys(in.Owned) + keys(in.Observed):
        if key == "recognized": return err(producer breach: reserved key)
    for row in in.Table.Rows:
        if "recognized" in row.RequiresOwned:
            return err(producer breach: RequiresOwned names reserved key)
    return nil                                   # unconditional on in.Recognized

Resolve(in):
    if err = ValidateInput(in); err != nil:      # same predicate, 2nd call site
        return Result{}, err                     # zero Result, no disposition
    # (RDR 0009's row-shape precondition also lands here; order open — A11)
    view = assemble(in)                          # unchanged: injects
    ...                                          #   in.Recognized under
    return result, nil                           #   "recognized" iff non-empty
```

GUESS summary: the exported predicate's name; the breach error's type and
text; the failure-payload and advisory struct shapes (RDR 0002's to
choose); whether the advisory reuses the failure type (guessed: no); the
exact fold used for the advisory ("Unicode-simple case folding" is named
but not pinned to a Go API — guessed `strings.ToLower`-equivalent
simple fold).
