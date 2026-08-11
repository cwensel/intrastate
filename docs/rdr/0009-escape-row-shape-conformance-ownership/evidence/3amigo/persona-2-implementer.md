Model: claude-opus-5[1m]
Lens: 3amigo (persona 2 — implementer)

# Persona 2 — Implementer: top 5 clarification-requests

Grounded against `intrastate` @ HEAD (branch `via-claude`, d82f1b5), Go 1.26.5,
`go.mod` `go 1.26.3`.

Scope note: signatures, symbol names/placement, interaction with shipped code,
ordering/precedence, implementability-as-written. Problem framing (P1) and
test pass/fail criteria (P3) deliberately excluded.

---

## IMP-1 — The exported error type and sentinel have no spelling, and the RDR explicitly defers it to a stage that already ran

**Severity**: high

**Passage** — *Technical Design → Normative Contracts* (typed-error clause):

> The breach error MUST be a TYPED error carrying the offending row identity
> as the kernel's existing RowRef value, inspectable via errors.As /
> errors.AsType without parsing message text, and MUST wrap a package-level
> sentinel naming the breach category so errors.Is can classify it. Both the
> error type and the sentinel MUST be EXPORTED

and *Load-Bearing Decisions → Naming*:

> The predicate is `CheckValid`/`Validate`-shaped (returns `error`); the exact
> type/sentinel spelling is sharpened at Pre-Lock.

**The question**: What are the exported identifiers, exactly? The RDR binds four
distinct public symbols into `internal/resolve` but names none of them:

1. the error **type** (`EscapeShapeBreachError`? `ShapeBreachError`? `RowShapeError`?)
2. its **fields** — the clause says it carries "the offending row identity as
   the kernel's existing RowRef value" (singular). But the A7 DX finding says
   "a typed breach error carrying `[]RowRef`" (plural), and Phase 1 says "the
   offending `RowRef`s" (plural). Simultaneously the aggregate design gives
   **one error per distinct identity**, which implies each typed error carries
   exactly **one** `RowRef` and the plurality lives in the `errors.Join`
   aggregate. Three passages, two incompatible shapes.
3. the **sentinel** (`ErrEscapeShapeBreach`?)
4. the **method** — `CheckValid` or `Validate`? The clause offers both
   ("its name MUST follow Go's error-returning convention — CheckValid or
   Validate") and never picks. It then cites *both* precedents
   (`pprof.Profile.CheckValid`, `rsa.PrivateKey.Validate`) without choosing.

**Grounding**: `internal/resolve/resolve.go` has **zero** exported error
symbols and does not import `errors` or `fmt` today (imports are `slices`,
`strings` only). Repo-wide there is no `var Err… = errors.New(…)` sentinel
anywhere and no existing `CheckValid`/`Validate` method — the only near-match
is `internal/cli/respond/respond.go::ValidateMode`, which returns
`*clierr.CLIError`, not `error`. So there is no in-repo house style to copy
from; every one of these is a genuinely free choice the implementer must make
and then permanently expose on RDR 0001's locked package surface.

**Decision it blocks**: Cannot write the type declaration, the sentinel `var`,
the `Error()`/`Unwrap()` methods, or the method name — i.e. cannot write
Phase 1 at all. Every downstream test in Validation scenarios 1, 9, 10, 10b
asserts through these identifiers, so they cannot be written either. Also
blocks the singular-vs-plural field decision, which changes whether
`errors.AsType[*T]` returns one row or all rows, and therefore whether the
Failure-Modes warning ("a single `errors.As`/`AsType` call reports only the
FIRST breach") is even true of the chosen shape.

---

## IMP-2 — "Collapse equal identities" + "state the breach count" is under-specified as a data shape, and interacts badly with the one-error-per-identity aggregate

**Severity**: high

**Passage** — *Technical Design → Normative Contracts* (multi-breach clause):

> The report is instead ordered by RowRef identity and made total by
> COLLAPSING equal identities: breaching rows sharing one RowRef contribute
> ONE reported error, not one per row. … Rows carrying no source identity
> collapse to a single RowRef{"",""} entry; that is a degenerate producer, and
> the error MUST remain diagnostic in that case by stating the breach count
> alongside the identities.

**The question**: Where does the count live, and is it per-identity or
whole-table? Three readings are all consistent with the text:

- **(a)** each per-row error carries `Count int` (rows collapsed into *this*
  identity), so `RowRef{"",""} ×3` renders "3 rows";
- **(b)** the aggregate carries a single whole-table count — but `errors.Join`
  returns an opaque `*joinError` with no room for a field, so this requires
  *not* using `errors.Join` as the outer wrapper, contradicting the clause
  that mandates it;
- **(c)** the count is only rendered into `Error()` prose — but A7/the DX
  finding forbid identity being "recoverable only from formatted prose", and a
  prose-only count is exactly the non-assertable diagnostic the DX finding
  rejects.

A8 ("If wrong") gestures at this and leaves it open: *"If a consumer needs a
per-row breach count rather than per-identity, the report must carry a count
field instead of collapsing"* — which reads as though collapsing and a count
field are **alternatives**, while the normative clause requires collapsing
**and** a count simultaneously.

**Grounding**: verified experimentally on Go 1.26.5 — `errors.Join(a, b)`
yields `"breach a\nbreach b"`, exposes `Unwrap() []error` of len 2, and
`errors.AsType[*E]` returns only the first. `errors.Join(nil)` returns nil.
There is no field on the join wrapper; a whole-table count genuinely cannot
ride the aggregate without a custom wrapper type — which the RDR never
authorizes and which would break the mandated `Unwrap() []error` traversal
shape if done naively.

**Decision it blocks**: Cannot define the struct fields of the typed error, and
cannot decide whether `Table.CheckValid` builds `[]error` then `errors.Join`s
(reading a/c) or needs a bespoke aggregate type (reading b). This is the core
data-shape decision of Phase 1 — everything else is mechanical once it's fixed.

---

## IMP-3 — `Resolve` must call the check "at entry", but the shipped `Resolve` computes `assemble(in)` first, and the whole-table precondition inverts a frozen precedence

**Severity**: medium

**Passage** — *Technical Design*:

> The one addition is a structural conformance predicate over `Table.Rows`,
> exported as a `Table` method so any producer can call it at construction
> time, and called by `resolve.go::Resolve` at entry before `assemble`

and *Load-Bearing Decisions → Precedence*:

> the breach error precedes every modeled disposition, `unmodeled_outcome`
> included.

**The question**: Two sub-questions the text does not resolve.

1. **Placement is stated against code that does something else first.**
   `internal/resolve/resolve.go::Resolve` opens with `view := assemble(in)`
   (line 319) and only then checks `in.Table.models(in.Recognized)` (line 321).
   The RDR says "at entry before `assemble`" — implementable, but it means the
   check must be inserted *above* the existing first statement, making the
   breach path skip `assemble` entirely. Is skipping `assemble` intended
   (i.e. is the check required to be cheap/side-effect-free-first), or is
   "before `assemble`" merely descriptive prose that would be equally satisfied
   by placing it after? `assemble` is pure and allocation-only, so both work —
   but the RDR pins one and the Performance Expectations section reasons about
   pass counts, so I cannot tell if the ordering is load-bearing or incidental.

2. **The precedence pin silently reorders a shipped behavior.** Today
   `unmodeled_outcome` is the *first* disposition and "therefore outranks a
   mere zero-match" per `Resolve`'s own doc comment (lines 300-302). Inserting
   a whole-table breach check above it makes the breach outrank
   `unmodeled_outcome`. The RDR asserts this ("precedes every modeled
   disposition, `unmodeled_outcome` included") but does **not** say the
   `Resolve` doc comment's numbered evaluation-order list must be amended.
   That list is a locked RDR 0001 contract artifact ("Evaluation order is fixed
   so the same tuple always replays the same disposition (REQ-1)").

**Grounding**: confirmed at `internal/resolve/resolve.go:318-323`. The doc
comment enumerates steps 1-4 and step 1 is currently the alphabet check. The
RDR's Phase 1 says to "record the kernel's validation surface" but never says
to renumber or prepend to this doc contract.

**Decision it blocks**: Whether to edit `Resolve`'s doc comment (a frozen RDR
0001 contract surface) as part of this change, and whether the check goes above
or below `assemble`. Editing a locked peer's doc contract without being told to
is exactly the kind of thing that gets reverted in review; not editing it leaves
the shipped doc actively wrong about evaluation order.

---

## IMP-4 — The CLI-wrap clause is normative but has no verb to attach to, and mandates a `clierr` envelope change with no field name

**Severity**: medium

**Passage** — *Technical Design → Normative Contracts* (CLI clause):

> When a breach reaches the CLI, the verb MUST wrap it into a *clierr.CLIError
> carrying a stable Code, Group GroupInternal (exit 2), the offending row
> identity, and a Hint stating the remedy … the offending row identities MUST
> reach the envelope through a serialized field, either the existing Detail
> (rendered from the row identities, never re-parsed by any consumer) or a new
> omitempty field added under the type's own "Extend with new optional fields
> as needed" allowance.

**The question**: Is any of this in scope for *this* RDR's implementation, and
if so what is the field and the code string? Concretely:

- The Implementation Plan has three Phases (predicate+precondition, fixtures,
  normalizer fixtures) — **none** of them mentions `clierr` or the verb. The
  CLI obligation appears only under *Prerequisites* as an unchecked box
  (`[ ] The flow verb that first calls Resolve wraps a breach…`), which reads
  as future work bound to a future RDR.
- But Validation **scenario 11** ("binds the verb that first calls `Resolve`")
  asserts exit 2 and identities readable off the **serialized** JSON envelope
  — a test that cannot exist without the verb.
- The clause mandates a "stable Code" but never spells it. The repo's
  precedent codes are kebab-case (`config-not-found`, `config-read-error`,
  `config-invalid`, `command-error`). Is it `escape-row-shape-breach`?
- The `Detail`-vs-new-field fork is left to the implementer, but a new field is
  a permanent addition to the wire envelope of an **Implemented** RDR 0005
  surface.

**Grounding**: confirmed — `internal/cli/root.go::NewRootCmd` registers only
`newVersionCmd()`; there is no `flow` verb at HEAD (the RDR states this
correctly in A2/A5). `internal/cli/clierr/clierr.go:48-67` shows `Cause error
\`json:"-"\`` and `Detail string \`json:"detail,omitempty"\`` — the RDR's
`json:"-"` claim is accurate. The `config.Load` precedent at
`internal/cli/config/config.go:84-90` is real and uses
`Code/Message/Detail/Group/Cause` — and, as the RDR says, carries **no** Hint
(the Hint precedent is one frame up at `config.go:60-65`, the
`config-not-found` branch, and at `root.go:122`).

**Decision it blocks**: Whether Phase 1 ends at the kernel boundary or must also
touch `internal/cli/clierr` (adding an `omitempty` field) — i.e. whether this
change is one package or three. Also blocks writing scenario 11 at all: if no
verb exists, that scenario is unimplementable this pass and should be marked as
binding a future RDR, the way scenario 8 explicitly is ("binds RDR 0002's
build").

---

## IMP-5 — "One predicate, two call sites" is stated as drift-proof, but the exported method and the entry check have different return obligations

**Severity**: low

**Passage** — *Technical Design → Normative Contracts* (exported-predicate
clause):

> The conformance predicate MUST be exported by the kernel package as a
> construction-time check callable by any table producer, and Resolve's entry
> precondition MUST be that same function — one predicate, two call sites, so
> a non-TOML producer can fail at build or construction time rather than at
> first production Resolve, and the two enforcement points cannot drift. It
> MUST be a METHOD ON Table taking no arguments … It MUST return error (nil
> when valid)

and *Technical Design*:

> When any row carries both a non-empty `Escape` and a non-empty `Writes`,
> `Resolve` returns the zero `Result` … and a non-nil error naming every such
> row.

**The question**: Does `Resolve` return `Table.CheckValid()`'s error verbatim,
or wrap it? The clause says the two call sites share *the same function*, but
the two call sites have different callers with different needs:

- called directly by a producer, the error is about a `Table` the producer
  holds — self-evidently scoped;
- returned from `Resolve`, the same error is one of several possible failure
  classes on a function whose error return is otherwise unused, and the CLI
  clause wants it wrapped into a `CLIError` with a code.

If `Resolve` wraps (e.g. `fmt.Errorf("resolve: %w", err)`), the aggregate gains
an extra layer — which is harmless for `errors.Is`/`AsType` but changes what
`Unwrap() []error` returns at the top level, and Validation scenario 9
explicitly asserts on "traversing the aggregate's `Unwrap() []error`". A
wrapped aggregate has `Unwrap() error`, not `Unwrap() []error`, at its
outermost layer — so the test as written would fail against a wrapping
implementation.

**Grounding**: verified on Go 1.26.5 — `errors.Join(a)` returns a wrapper whose
`Error()` degrades to the single error's own message (`"breach a"`), exposing
`Unwrap() []error`. The RDR's claim that Join "wraps even this single-breach
case" is correct. But nothing in the repo establishes whether kernel errors get
a package prefix: there are no existing errors in `internal/resolve` to imitate
(A1 verified all five return paths carry literal `nil`).

**Decision it blocks**: Whether `Resolve`'s breach return is
`return Result{}, err` or `return Result{}, fmt.Errorf(…: %w, err)`. Cheap to
change later, but it determines the exact traversal shape every multi-breach
test asserts against, so it should be pinned before those tests are written.

---

## Grounding summary — RDR claims checked against HEAD

| RDR claim | Verdict |
| --- | --- |
| `fixtures_test.go::escapeRow` sets `Writes: {status Blocked}` at line 257 | **Confirmed** — line 257 exactly |
| Builder inherited by **16** call sites | **Confirmed** — 16 `escapeRow(` call sites across adversarial/resolve/fixup tests |
| Two overrides at `adversarial_test.go:229`, `fixup_test.go:103` | **Confirmed** — both are `escape.Writes = []resolve.Tag{{Key: "status", Value: "Escaped"}}` |
| A1: all `Resolve` return paths carry literal `nil` error | **Confirmed** — `resolve.go:322,344,348,350,352` |
| `resolve.go::Resolve` opens with `assemble`, then the alphabet check | **Confirmed** — lines 319, 321 |
| `Row` carries `RuleID`/`SourceLocator`; mirrored on `Plan` and `RowRef` | **Confirmed** — lines 173-174, 240-241, 276-279 |
| `compareRefs` returns 0 on equal (RuleID, SourceLocator); `rowRefs` uses non-stable `slices.SortFunc` | **Confirmed** — lines 528-533, 546 |
| Three `fixup_test.go` sites build identical `escapeRow("rdr.escape.ambiguous", "flows/rdr.toml:99", …)` | **Confirmed** — lines 79, 101, 125 |
| `clierr.CLIError.Cause` is `json:"-"`; `Detail` is serialized | **Confirmed** — clierr.go:56, 66 |
| `ExitCodeFor`: `GroupUserEnv, GroupInternal → 2`; non-`CLIError` → 1 | **Confirmed** — clierr.go:121, 127 |
| `config.Load` wraps a Go error into `GroupInternal` with `Code`+`Detail`+`Cause`, **no** Hint | **Confirmed** — config.go:84-90 |
| No `flow` verb at HEAD; `NewRootCmd` registers only `version` | **Confirmed** |
| `errors.Join` wraps even one error; `Error()` degrades to the single message; `AsType` finds only the first | **Confirmed experimentally**, Go 1.26.5 |
| `errors.AsType[E]` exists (RDR cites it as an assertion instrument) | **Confirmed** — `go doc errors.AsType` resolves on Go 1.26.5 |
| No existing `CheckValid`/`Validate` returning `error` in-repo to imitate | **Confirmed** — only `respond.ValidateMode`, which returns `*clierr.CLIError` |
| `internal/resolve` imports only `slices`, `strings` | **Confirmed** — will need `errors` and probably `fmt` |

No RDR claim about the code was found to be **false**. The gaps above are
under-specification, not misstatement.
