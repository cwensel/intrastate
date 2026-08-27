Model: claude-opus-5

# RDR 0009 — DX prior-art pass for the `internal/resolve` conformance predicate

Scope: developer experience of the structural-validation API only. The
design choice (an exported conformance predicate over `Table`, plus a call
to that same predicate at `Resolve` entry returning a non-nil Go error on
breach) is settled and is NOT re-litigated here.

Toolchain read: `go env GOROOT` = `/opt/local/lib/go-1.26`, `go1.26.5`.
All GOROOT paths below are rooted at `/opt/local/lib/go-1.26/src`.

Target under discussion, for grounding:

- `./internal/resolve/resolve.go:204`
  — `type Table struct { Revision string; Outcomes []string; Rows []Row }`
- `./internal/resolve/resolve.go:318`
  — `func Resolve(in Input) (Result, error)`
- `resolve.go:293-296` — the current doc sentence: "A modeled refusal
  travels the Result value with a nil error; the error return is reserved
  for programmer mistakes, not for modeled refusals."

---

## Q1 — Error shape: opaque vs sentinel vs typed

### Finding: Go convention is a TYPED error carrying the offending identity as exported fields, optionally wrapping a sentinel for the category.

Go's stdlib does not use formatted-string-only errors when the caller may
need to know *which* element broke. Every structural-error type read here
carries the offending element's identity as exported struct fields, and the
message string is derived from those fields — never the other way round.

**`encoding/json.UnmarshalTypeError`** — `encoding/json/decode.go:125-141`:

```go
// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value  string       // description of JSON value - "bool", "array", "number -5"
	Type   reflect.Type // type of Go value it could not be assigned to
	Offset int64        // error occurred after reading Offset bytes
	Struct string       // name of the struct type containing the field
	Field  string       // the full path from root node to the field, include embedded struct
}
```

`Field` is documented as "the full path from root node to the field" — the
stdlib treats *locating the offending element* as a first-class field, not
as message prose. `Error()` at `:135` is a pure projection of those fields.

**`encoding/json/v2.SemanticError`** — `encoding/json/v2/errors.go:63-93`.
This is the newest stdlib thinking on the same problem (Go 1.26,
`goexperiment.jsonv2`) and it doubles down: the element identity is a
*structured path*, not a string.

```go
type SemanticError struct {
	requireKeyedLiterals
	nonComparable

	action string // either "marshal" or "unmarshal"

	// ByteOffset indicates that an error occurred after this byte offset.
	ByteOffset int64
	// JSONPointer indicates that an error occurred within this JSON value
	// as indicated using the JSON Pointer notation (see RFC 6901).
	JSONPointer jsontext.Pointer
	...
	// Err is the underlying error.
	Err error // may be nil
}
```

Two DX lessons here, both explicit in the source:

1. The **category** is a sentinel in `Err`, the **location** is a field.
   `encoding/json/v2/errors.go:25-40` documents `ErrUnknownName` and the
   exact caller idiom:

   ```go
   err := ...
   serr, ok := errors.AsType[*json.SemanticError](err)
   if ok && serr.Err == json.ErrUnknownName {
       ptr := serr.JSONPointer // JSON pointer to unknown name
       name := ptr.LastToken() // unknown name itself
       ...
   }
   ```

   This idiom is repeated verbatim as a runnable example at
   `encoding/json/v2/example_test.go:374-377`. That is the shape that makes
   an error **assertable without string matching**: `errors.AsType` for the
   type, `== sentinel` for the category, field read for the identity.

2. `requireKeyedLiterals` and `nonComparable` are embedded deliberately, so
   callers cannot construct the type positionally nor compare it with `==`.
   The stdlib forces inspection through `AsType` + fields.

**`go/types.Error`** — `go/types/api.go:51-71` — same pattern for a
*checker over a caller-supplied structure*, which is the nearest analogue to
a table validator:

```go
type Error struct {
	Fset *token.FileSet // file set for interpretation of Pos
	Pos  token.Pos      // error position
	Msg  string         // error message
	Soft bool           // if set, error is "soft"
	...
}
```

Note the split: `Pos` (machine-readable offending location) and `Msg`
(human prose) are separate fields, and `Error()` at `:69` merely formats
`Pos` + `Msg`.

**`time.ParseError`** — `time/format.go:841-847` — carries
`Layout, Value, LayoutElem, ValueElem, Message`: again, the offending
sub-element (`ValueElem`) is a field.

**`net/url.Error`** — `net/url/url.go:32-39` — the wrapper shape:
`Op, URL string; Err error` plus `func (e *Error) Unwrap() error`. This is
the precedent for wrapping a category sentinel inside a context-carrying
typed error.

### `errors.AsType` (Go 1.20 `errors.As` successor, generic form)

`errors/wrap.go:154-167`:

> AsType finds the first error in err's tree that matches the type E, and
> if one is found, returns that error value and true. [...] The tree
> consists of err itself, followed by the errors obtained by repeatedly
> calling its `Unwrap() error` or `Unwrap() []error` method. When err wraps
> multiple errors, AsType examines err followed by a depth-first traversal
> of its children.

Load-bearing for Q1+Q2 jointly: `AsType` traverses `Unwrap() []error`, so a
typed per-row error stays individually assertable *even when aggregated*
via `errors.Join`. Typed-error and aggregate-reporting are not in tension.

### What this rules out

An opaque `errors.New`/`fmt.Errorf` string is the shape stdlib uses only
where the caller has no decision to make (e.g. `net/http.Cookie.Valid()`,
`net/http/cookie.go:328-361`, which returns bare `errors.New("http: invalid
Cookie.Name")` — a boolean-grade verdict for a fixed-shape struct with no
element to name). A `Table` with N rows is not that case: there IS an
offending element (the row) and a test WILL want to assert on it.

A bare sentinel alone is also insufficient — it carries the category but not
the row index. The stdlib answer is *both*: sentinel in a wrapped `Err`
field, identity in exported fields.

---

## Q2 — Aggregate vs first-breach reporting

### Finding: prior art is clear and consistent — validators that check a *document* report ALL breaches; it is the parsers on a byte stream that stop at the first.

The distinction in the stdlib is not arbitrary. It tracks whether
continuing after the first breach is *meaningful*:

**Fail-fast, and why**: `encoding/json.Unmarshal` documents at
`encoding/json/decode.go:88-91` that it returns "an `UnmarshalTypeError`
describing the **earliest** such error. In any case, it's not guaranteed
that all the remaining fields following the problematic one will be
unmarshaled." A byte-stream parser loses positional confidence after a
breach, so continuing is not sound. A table walk over `[]Row` has no such
problem — row *i+1* is independently checkable.

**Aggregate, three stdlib instances, all in the "check a caller-supplied
structure" shape that matches `Table`:**

1. **`testing/fstest.TestFS`** — `testing/fstest/testfs.go:20-37`, doc
   comment:

   > If TestFS finds any misbehaviors, it returns either the first error or
   > a list of errors. Use `errors.Is` or `errors.AsType` to inspect.

   Mechanism at `testfs.go:96-106` and `:92` — accumulate into a slice,
   join at the end:

   ```go
   type fsTester struct {
       fsys   fs.FS
       errors []error
       ...
   }

   // errorf adds an error to the list of errors.
   func (t *fsTester) errorf(format string, args ...any) {
       t.errors = append(t.errors, fmt.Errorf(format, args...))
   }
   ```

   and the return:

   ```go
   return fmt.Errorf("TestFS found errors:\n%w", errors.Join(t.errors...))
   ```

   This is a near-exact structural match for the proposed predicate: an
   exported function taking a caller-supplied structure, walking all of it,
   returning a single `error` that is an `errors.Join` of per-element
   errors, with the doc comment naming `errors.AsType` as the inspection
   route.

2. **`testing/slogtest.TestHandler`** — `testing/slogtest/slogtest.go:248-250`:

   > TestHandler tests a [slog.Handler]. **If TestHandler finds any
   > misbehaviors, it returns an error for each, combined into a single
   > error with [errors.Join].**

   Implementation `slogtest.go:284-292`: accumulate `errs`, `return
   errors.Join(errs...)`. Note `errors.Join` returns nil when every element
   is nil (`errors/join.go:20-30`), so the "no breach ⇒ nil error" contract
   is free.

3. **`go/types.Config.Error`** — `go/types/api.go:159-166`. The stdlib makes
   this an explicit, documented *caller choice*, and the default is
   revealing:

   > If Error != nil, it is called with each error found during type
   > checking; err has dynamic type Error. Secondary errors [...] have
   > error strings that start with a '\t' character. **If Error == nil,
   > type-checking stops with the first error found.**

   The type checker's *considered* behavior when the caller wants a good
   experience is all-errors; first-error is the degraded fallback.

**`errors.Join` semantics** — `errors/join.go:11-19`:

> Join returns an error that wraps the given errors. Any nil error values
> are discarded. Join returns nil if every value in errs is nil. The error
> formats as the concatenation of the strings obtained by calling the Error
> method of each element of errs, with a newline between each string. A
> non-nil error returned by Join implements the `Unwrap() []error` method.
> The errors may be inspected with [Is] and [As].

So: single `error` return type preserved, one-error case degrades to that
error's own message verbatim (`join.go:46-49` — `if len(e.errs) == 1 {
return e.errs[0].Error() }`), and every joined element stays assertable.

### Peer state-machine engines (sibling repo)

`../state-machines/repos/awf-cli` — a Go
workflow/state-machine CLI, the closest domain peer found. It is
**aggregate, and says so in its package doc**:

- `repos/awf-cli/pkg/validation/doc.go:13` — feature list includes
  "Aggregated error reporting (not fail-fast)".
- `repos/awf-cli/pkg/validation/doc.go:49` — "All errors are collected and
  returned together (not fail-fast)."
- `repos/awf-cli/internal/domain/workflow/validation_errors.go:68-73` — the
  per-issue type, which independently corroborates Q1:

  ```go
  // ValidationError represents a single validation issue.
  type ValidationError struct {
      Level   ValidationLevel
      Code    ValidationCode
      Message string
      Path    string // e.g., "states.validate.on_success"
  }
  ```

  A stable machine-readable `Code` (`ErrCycleDetected`,
  `ErrInvalidTransition`, ... at `validation_errors.go:30-66`) plus a
  `Path` locating the offending element — i.e. sentinel-category plus
  identity-field, arrived at independently of the stdlib.
- `validation_errors.go:89-92` — `ValidationResult{Errors, Warnings
  []ValidationError}` — a result *collection*, not a single error.

`repos/statewright/crates/agent/src/validator.rs:180` and
`crates/engine/src/types.rs:374` — both `pub errors: Vec<String>`. Aggregate
in shape, but stringly-typed; corroborates Q2 only, and is a negative
example for Q1.

### Rejected / not found

- JSON Schema "all errors vs fail-fast" convention and protovalidate: **no
  local source available.** Neither is vendored in GOROOT nor present in
  the sibling repo checkouts. I did not cite them from memory. This leaves
  Q2's cross-ecosystem breadth thinner than the stdlib evidence, but the
  stdlib evidence is itself unambiguous and directly on-shape.
- Kubernetes `field.ErrorList` / `utilerrors.NewAggregate`: **no local
  checkout.** Not cited.
- `go vet` diagnostics: `cmd/vendor/golang.org/x/tools/go/analysis/unitchecker/unitchecker.go:474`
  (`act.diagnostics = append(act.diagnostics, d)`) accumulates all
  diagnostics per analyzer rather than stopping at the first. Corroborating
  but weaker than the three stdlib cases above (diagnostics are not
  `error`s), so used as support, not as the decisive citation.

---

## Q3 — Where validation belongs: constructor vs entry-check vs both

### Finding: "both" has direct stdlib precedent, and the standard optimization is a cheap idempotence short-circuit — not the removal of the entry check.

**Decisive precedent: `github.com/google/pprof/profile.Profile.CheckValid`**
(vendored in the toolchain, `cmd/vendor/github.com/google/pprof/profile/profile.go`).
This is exactly the proposed arrangement — one exported predicate over a
caller-supplied structure, called at more than one boundary:

- Defined exported at `profile.go:362` — `func (p *Profile) CheckValid() error`.
- Called defensively at the parse boundary, `profile.go:202-204`:

  ```go
  if err := p.CheckValid(); err != nil {
      return nil, fmt.Errorf("malformed profile: %v", err)
  }
  ```

- Called again at the end of a mutating operation, `profile.go:487`:
  `return p.CheckValid()`.

So the same exported predicate serves both the producer (who can call it
directly) and the library's own entry points. This is the pattern under
discussion, in vendored production code, with no apologetics in the
comments.

**Second precedent, with an explicit cost note:
`crypto/rsa.PrivateKey.Validate`** — `crypto/rsa/rsa.go:229-233`:

```go
// Validate performs basic sanity checks on the key.
// It returns nil if the key is valid, or else an error describing a problem.
//
// It runs faster on valid keys if run after [PrivateKey.Precompute].
func (priv *PrivateKey) Validate() error {
```

and inside, `rsa.go:240-244`:

```go
// If Precomputed.fips is set and consistent, then the key has been
// validated by [rsa.NewPrivateKey] or [rsa.NewPrivateKeyWithoutCRT].
if priv.precomputedIsConsistent() {
    return nil
}
```

Two DX lessons: (a) the constructor itself calls the exported predicate —
`rsa.go:524-527`, `priv.Precompute(); if err := priv.Validate(); err != nil
{ return nil, err }` — so construction-time and defensive use share one code
path; (b) the *cost* of re-validation is addressed by **short-circuiting on
a recorded already-validated marker**, not by deleting the check. That is
the reusable technique if per-call `O(rows)` ever measures.

**Third: `x/tools/go/analysis.Validate`** —
`cmd/vendor/golang.org/x/tools/go/analysis/validate.go:14-24`:

```go
// Validate reports an error if any of the analyzers are misconfigured.
// Checks include:
// that the name is a valid identifier;
// that the Doc is not empty;
// that the Run is non-nil;
// that the Requires graph is acyclic;
// ...
func Validate(analyzers []*Analyzer) error {
```

Exported, over a caller-supplied `[]*Analyzer` (structurally the same shape
as `[]Row`), and called by the driver at its entry point —
`unitchecker/unitchecker.go:99`: `if err := analysis.Validate(analyzers);
err != nil {`. Same "producer may call it, entry point does call it" split.
Note this one is **fail-fast** (each check `return`s immediately), which is
a genuine split in the prior art on Q2 — see the honesty note below.

**On the cost of defensive revalidation — Meyer, OOSC 2e §11.13**
(`References/DevRefGit/Books/Object Oriented
Software Construction-Meyer.pdf`, printed pp. 396-397 = PDF pages 418-419):

- p. 396 quotes Hoare (1973): "It is absurd to make elaborate security
  checks on debugging runs, when no trust is put in the results, and then
  remove them in production runs, when an erroneous result could be
  expensive or disastrous."
- p. 396 gives the empirical figure: "In ISE's experience the cost for
  monitoring **preconditions** (the default option, including of course
  array bounds checking) is on the order of **50%**. What is frustrating is
  that more than 75% of that cost is due not to precondition checking per
  se but to the supporting machinery of monitoring calls — recording every
  routine entry and every routine exit."
- p. 397 contrasts: "As to **postcondition and invariant** checking, they
  can bring the penalty to **100% to 200%**. (Although circumstances vary,
  preconditions are often relatively simple consistency conditions [...]
  whereas many postconditions and invariants express more advanced semantic
  properties.)"
- p. 397, the governing principle: "There is never any justification for
  compromising on correctness for the sake of other concerns, such as
  efficiency."

Applied here: an entry-point structural check on `Table` is a
**precondition**, Meyer's cheap class — and the 50% figure is dominated
(>75% of it) by *tracing machinery* that a plain Go `if err := ...` does not
have. The literature does not support removing the entry check on cost
grounds absent a measurement; it supports keeping it and short-circuiting if
measurement ever demands (the `rsa` technique above).

---

## Q4 — Naming

### Convention for return type

Surveyed all exported `Validate|Check|Verify|IsValid|Valid|OK` declarations
in GOROOT (`grep -rhn -E "^func (\([^)]*\) )?(Validate|Check|...)[A-Za-z]*\("`,
`_test` excluded). The split is clean and is about **return type**, not
about the verb:

**`Is…`/bare `Valid` ⇒ `bool`, cheap, no reason available:**
- `netip.Addr.IsValid() bool` — `net/netip/netip.go:390`
- `netip.Prefix.IsValid() bool` — `net/netip/netip.go:1325`
- `netip.AddrPort.IsValid() bool` — `net/netip/netip.go:1151`
- `reflect.Value.IsValid() bool` — `reflect/value.go:1657`
- `token.Position.IsValid() bool` — `go/token/position.go:34`
- `obj.FuncInfo.Valid() bool` — `cmd/internal/obj/link.go:2068`

**`Validate`/`Check…`/`Verify…` ⇒ `error`, when the *reason* matters:**
- `rsa.PrivateKey.Validate() error` — `crypto/rsa/rsa.go:233`
- `analysis.Validate([]*Analyzer) error` — `x/tools/go/analysis/validate.go:24`
- `pprof.Profile.CheckValid() error` — `pprof/profile/profile.go:362`
- `x509.Certificate.CheckSignature(...) error` — `crypto/x509/x509.go`
- `x509.Certificate.VerifyHostname(h string) error` — `crypto/x509/verify.go`
- `http.CrossOriginProtection.Check(req *Request) error` — `net/http/csrf.go`
- `types.Config.Check(path, fset, files, info) (*Package, error)` — `go/types/api.go`

**The one instructive outlier**: `net/http.Cookie.Valid() error`
(`net/http/cookie.go:327-328`), whose doc comment reads "Valid reports
whether the cookie is valid" — *`reports whether`* is Go's documented
phrasing convention for a `bool`-returning predicate, yet the function
returns `error`. The name and the doc comment both mislead about the return
type. This is a known wart, not a model: it argues against naming an
`error`-returning function `Valid`.

**`Verify`** is not neutral in Go — it is overwhelmingly cryptographic
(`Certificate.Verify`, `VerifyHostname`, `MAC.Verify`, `verifier.Verify`).
Using it for structural conformance would import the wrong connotation.

**`Check`** alone is used in Go for operations that *do work and return a
product* (`types.Config.Check` returns `(*Package, error)`), and for
request-time policy checks. `CheckValid` disambiguates and is the exact
precedent shape (`pprof.Profile.CheckValid`).

**`OK`** — no exported instance found in GOROOT for this role. Not a Go
convention for this.

### Doc-comment convention

Go's `reports whether` phrasing is reserved for `bool` returns (visible
across `netip`, `reflect`, `errors.Is` — `errors/wrap.go:27` "Is reports
whether any error in err's tree matches target"). An `error`-returning
validator should use the `rsa.PrivateKey.Validate` phrasing instead:
"returns nil if [...] is valid, or else an error describing a problem"
(`crypto/rsa/rsa.go:230`).

---

## Query / source ledger

### Accepted (read in full, cited above)

| # | Source | Used for |
|---|---|---|
| 1 | `encoding/json/decode.go:88-91,125-141,156-166` — `UnmarshalTypeError`, `InvalidUnmarshalError`, `Unmarshal` doc | Q1 typed error w/ `Field`; Q2 *why* stream parsers fail fast |
| 2 | `encoding/json/v2/errors.go:25-40,63-93` — `ErrUnknownName`, `SemanticError` | Q1 sentinel-in-`Err` + identity-in-field; `requireKeyedLiterals`/`nonComparable` |
| 3 | `encoding/json/v2/example_test.go:360-380` | Q1 runnable caller idiom w/ `errors.AsType` |
| 4 | `errors/wrap.go:27,154-167` — `Is`, `AsType` | Q1 assertability; Q2 `Unwrap() []error` traversal; Q4 `reports whether` |
| 5 | `errors/join.go:11-62` — `Join`, `joinError` | Q2 aggregate mechanics, nil-when-all-nil, 1-element degradation |
| 6 | `testing/fstest/testfs.go:20-37,92,96-106` — `TestFS`, `fsTester.errorf` | Q2 decisive: exported structure-validator returning `errors.Join`, doc names `AsType` |
| 7 | `testing/slogtest/slogtest.go:248-250,284-292` — `TestHandler` | Q2 decisive: "returns an error for each, combined into a single error with errors.Join" |
| 8 | `go/types/api.go:51-71,159-166` — `Error`, `Config.Error` | Q1 `Pos`/`Msg` split; Q2 all-errors-vs-first as documented caller choice |
| 9 | `cmd/vendor/github.com/google/pprof/profile/profile.go:202-204,362,487` — `CheckValid` | Q3 decisive: exported predicate called at 2 internal boundaries; Q4 name |
| 10 | `crypto/rsa/rsa.go:229-244,524-527` — `PrivateKey.Validate` | Q3 constructor calls same predicate + idempotence short-circuit; Q4 name & doc phrasing |
| 11 | `cmd/vendor/golang.org/x/tools/go/analysis/validate.go:14-24` + `unitchecker/unitchecker.go:99` | Q3 exported validator over caller slice, called at driver entry; Q2 counter-example (fail-fast) |
| 12 | `net/http/cookie.go:327-361` — `Cookie.Valid() error` | Q4 instructive outlier: `Valid`+`reports whether` misnames an `error` return |
| 13 | `net/url/url.go:32-53` — `url.Error` | Q1 wrapper shape `Op/URL/Err` + `Unwrap` |
| 14 | `time/format.go:841-855` — `time.ParseError` | Q1 offending sub-element as field (`ValueElem`) |
| 15 | GOROOT sweep: `grep -rhn -E "^func (\([^)]*\) )?(Validate\|Check\|Verify\|IsValid\|Valid\|OK)[A-Za-z]*\("`, `_test` excluded | Q4 the `bool` vs `error` naming split |
| 16 | `state-machines/repos/awf-cli/pkg/validation/doc.go:13,49` | Q2 peer engine: "Aggregated error reporting (not fail-fast)" |
| 17 | `state-machines/repos/awf-cli/internal/domain/workflow/validation_errors.go:16-73,89-92` | Q1+Q2 peer: `Code`+`Path` typed issue, `ValidationResult` collection |
| 18 | Meyer, *OOSC 2e* §11.13, pp. 396-397 (PDF pp. 418-419), incl. Hoare 1973/1981 quotes | Q3 cost of defensive revalidation: precondition ≈50% (>75% of it tracing machinery); postcondition/invariant 100-200% |

### Rejected

| Query / source | Why rejected |
|---|---|
| `arc search semantic --corpus DevRef --limit 6 "API design validation errors report all errors versus fail fast usability"` | 6 hits, top 0.603, all generic microservices/HTTP-API prose (`advanced-microservices.pdf:54`, `microservices-from-day-one.pdf:55-63`, `api-driven-devops.pdf:61`). No coverage of validator error-shape or aggregation. No citation taken. |
| `arc search semantic --corpus DevRef --limit 6 "parse don't validate make illegal states unrepresentable constructor invariant"` | 6 hits, top 0.539. Meyer OOSC pp. 388-401 (class invariants) and `tdd-ebook-sample.pdf:261` (factory-method validation) are adjacent but generic OO; neither addresses the exported-predicate-plus-entry-check question. **Not cited** — the concrete Go precedents (pprof, rsa) are strictly stronger. "Parse, don't validate" itself returned no corpus hit. |
| JSON Schema validators' all-errors-vs-fail-fast convention | No local source. Not vendored in GOROOT, absent from sibling checkouts. **Not cited from memory.** Q2's cross-ecosystem breadth is therefore thinner than its stdlib evidence. |
| protobuf / protovalidate | Same — no local checkout. Not cited. |
| Kubernetes `field.ErrorList` / `utilerrors.NewAggregate` | Same — no local checkout. Not cited. |
| `state-machines/repos/statewright` (`crates/agent/src/validator.rs:180`, `crates/engine/src/types.rs:374`, both `pub errors: Vec<String>`) | Aggregate in shape → weak Q2 support only. Stringly-typed → **negative** example for Q1. Not load-bearing. |
| `state-machines/repos/StateSmith` `Validate()` hits | All in `src/StateSmithTest/*.cs` test files; no production validator API surfaced. Rejected. |
| `state-machines/repos/xstate`, `study/` sweep for "all errors"/"fail-fast" | No matches outside `node_modules`. No evidence either way. |
| `go vet` / `unitchecker.go:474` diagnostic accumulation | Corroborating but weak — `analysis.Diagnostic` is not an `error`, so it does not speak to error *shape*. Used as support, not as a decisive Q2 citation. |
| Go blog / wiki error-handling posts in the local distribution | Searched; the Go 1.26 distribution ships no `doc/` blog or wiki prose on error handling. Guidance was taken from normative doc comments in `errors/`, which is stronger evidence anyway. |

---

## Honesty note: where prior art is split

**Q2 is not unanimous.** `analysis.Validate`
(`x/tools/go/analysis/validate.go:24`) is fail-fast — every check `return`s
on first breach — and it validates a caller-supplied slice, the same shape
as `[]Row`. The split is explicable: `Validate` there guards a *developer
wiring mistake at process startup* where one message suffices, whereas
`fstest.TestFS` and `slogtest.TestHandler` validate an *artifact the user is
iterating on*, and both chose aggregate. The `Table` case is the latter — a
table author iterating on rows — so the aggregate side of the split
applies. But the split is real and should not be presented as consensus.

**Q3 has no direct precedent for a *hot-path* re-check.** `pprof.CheckValid`
runs at parse and after aggregation, `rsa.Validate` at key construction —
neither is a per-call hot path. The `rsa` idempotence short-circuit is the
closest thing to guidance, and it is a technique, not an endorsement. No
source read here measures the cost of per-call `O(rows)` structural
revalidation on a stateless entry point. Meyer's figures are the nearest
quantitative anchor and they concern instrumented contract monitoring, not
a hand-written Go loop, so they bound the question loosely at best.
