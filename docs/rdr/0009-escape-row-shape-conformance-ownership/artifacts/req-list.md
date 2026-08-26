# REQ List — RDR 0009 Ownership of escape-row shape conformance

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0009-escape-row-shape-conformance-ownership.md`. Quotes are verbatim,
copied from the projector (`rdr inspect --select <id>`) or from the record
itself, never transcribed.

Element ids are carried where the REQ derives from a labelled contract
(`0009:C1` … `0009:C8`, `0009:MVV`). Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts
- `LBD` = Proposed Solution / Technical Design / Load-Bearing Decisions
- `TD` = Proposed Solution / Technical Design (prose, outside fences)
- `AP` = Proposed Solution / Approach
- `CA` = Research Findings / Critical Assumptions
- `FM` = Trade-offs / Failure Modes
- `CONS` = Trade-offs / Consequences
- `EIA` = Implementation Plan / Existing Infrastructure Audit
- `PRE` = Implementation Plan / Prerequisites
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (Phases 1–3)
- `TS` = Validation / Testing Strategy
- `PE` = Validation / Performance Expectations

**HEAD-state correction (Phase 0, amended).** The record's A2/A5 evidence
("no `flow` verb exists at HEAD"; RDR 0002 "Final, not yet implemented") is
**stale at this commit** and does not govern executability. Verified in this
worktree:

- The `flow` verb **ships** — `internal/cli/flow.go`, `flow_exec.go`,
  `flow_resolve.go`, `flow_next.go`, `flow_input.go`, `flow_state.go`, with the
  RDR-0005 suites `flow_*_0005_test.go` and a typed-code vocabulary
  (`codeTagInvalid = "flow-tag-invalid"` …, `flow_input.go:34-59`).
- `internal/cli/flow_resolve.go:105-116` **already calls `resolve.Resolve`** and
  already has a live error branch — today it re-codes the kernel error as
  `codeAccessorFailed` and folds it into prose via `err.Error()`. That is the
  exact call site C7 binds, and it is editable now.
- The serialized carrier C7 calls "a NEW omitempty field" **already ships** as
  `clierr.CLIError.Findings []Finding` (`json:"findings,omitempty"`) with the
  `Finding` record (`Code`, `Message`, `Param`, `Locator`, `Hint`, …). It was
  landed by RDR 0008 per JDR 0001 §D10 and is pinned by
  `internal/cli/reserved_key_0008_test.go::TestReq102_TheFindingCarrierHasTheFiveFieldsIncludingHint`.
  This is what `deviations.md` D1 anticipated ("one `omitempty` `findings` list
  on the envelope, which is the carrier for this RDR's row identities").
- **RDR 0002 is `Implemented`** (record Status line; launch `status.md` state
  `COMPLETE`), surface at `internal/table/`.

**Scope note (what is actually deferred, and why).** Scenarios 1–7, 7b, 9, 10,
10b execute. The C7 CLI block (Section H) is **executable and in scope** — the
false "no verb at HEAD" premise is withdrawn. Exactly two deferrals remain, and
both are mandated by the **record's own** TS "Scope of this RDR's Done criteria"
note, not by HEAD's state:

- **REQ-90 (TS scenario 8)** — the record says it "binds RDR 0002's build" and
  places it out of this RDR's Done criteria. RDR 0002 being implemented does
  **not** un-defer it: the record defers it by *ownership*, not by readiness.
  Its authored-path obligation is discharged here as Phase 3's shared fixture
  set (REQ-97/REQ-98), which is what the record asks this RDR to ship.
- **REQ-91 (TS scenario 11)** — the record marks it verbatim "**DEFERRED: not
  executable by this RDR**". Its stated *reason* ("no `flow` verb exists at
  HEAD") is now false, but the deferral is the record's own scope decision and
  is honored. Its substance is not lost: C7's own clauses (REQ-58 … REQ-69) are
  executable and carry the same obligations, so the CLI surface is covered by
  contract REQs even though the record's scenario 11 is not run as written.

No other REQ is deferred.

**Standing obligations from `deviations.md` (D1, D2)** are citation/consistency
repairs, not REQs; they are checked at their own phase and are not restated
below.

---

## A. The producer obligation (the single rule)

- [REQ-1] `0009:C1` "Escape-row shape conformance — an escape row carries no
  owned-state mutation — is a PRODUCER obligation on every
  constructor of resolve.Row values: a Row with a non-empty
  Escape list MUST have an empty Writes slice." — (NC)

- [REQ-2] `0009:C1` "(Authored clears
  normalize to `<clear>` writes per RDR 0002, so the Writes
  predicate carries both "no writes" and "no clears" at the
  kernel boundary.)" — (NC)

- [REQ-3] `0009:C1` "The predicate is Writes-only and does NOT
  extend to NextTags: A4 settled that owned state is reachable
  only through a write accessor, which RDR 0004 scopes to
  "planned owned-tag writes," so an escape row's NextTags
  cannot mutate owned state." — (NC)

- [REQ-4] `0009:D-selection-predicate` "the breach predicate over a
  single row is `len(row.Escape) != 0 && len(row.Writes) != 0`." — (LBD)

- [REQ-5] `0009:D-selection-predicate` "*Nil-vs-empty* — length-based on
  purpose: a non-nil but empty slice is conforming. The predicate tests
  emptiness, never nil-ness." — (LBD)

- [REQ-6] `0009:C8` "The shared resolve.Row type keeps its single shape: escape
  identity remains discriminated solely by a non-empty Escape
  list. No ordinary/escape type split and no row-kind field is
  introduced at the kernel boundary." — (NC)

- [REQ-7] "The escape discriminator stays
  `len(Escape) != 0` / `rescues` — the existing sibling
  signal, reused." — (TD)

- [REQ-8] "`planOf`
  stays unconditional — under the precondition it can no
  longer copy writes onto an escaped plan, so no stripping
  logic is added" — (TD)

- [REQ-9] "Data flow is unchanged for every conforming table:
  `assemble` → candidate partition (`len(row.Escape) != 0`) →
  `gate` → selection → `planOf`." — (TD)

- [REQ-10] `0009:BR2` "Kernel sanitization (strip `Writes` when the selection
  is an escape)" is rejected — recorded as Briefly Rejected; no
  sanitization/stripping path may be introduced. — (Briefly Rejected;
  cross-checked against REQ-8)

## B. The authored-path enforcer (RDR 0002's half)

- [REQ-11] `0009:C2` "On the authored-table path, RDR 0002's normalizer/lint is the
  enforcing implementation: an authored escape rule carrying a
  write block or clear list MUST be rejected at table load
  under RDR 0002's "malformed escape declaration" validation
  class, diagnosable to the source rule." — (NC)

- [REQ-12] `0009:C2` "Rejection keys on the
  PRESENCE of the block, not its contents — RDR 0002 forbids an
  escape rule from containing a write block, so an empty one
  (`writes = []`) is rejected too." — (NC)

- [REQ-13] `0009:C2` "The authored layer is
  therefore stricter than the kernel's length-based predicate,
  deliberately" — (NC)

- [REQ-14] `0009:C2` "Normalized escape rows
  MUST render write-free. This RDR binds that duty with a
  conformance fixture set; RDR 0002's grammar is unchanged." — (NC)

- [REQ-15] "Ownership: **this RDR is the single normative home of the
  escape-row shape rule.**" — (AP)

## C. The kernel entry precondition

- [REQ-16] `0009:C3` "The kernel enforces the same obligation as an entry
  precondition of Resolve." — (NC)

- [REQ-17] `0009:C3` "Resolve's signature is UNCHANGED —
  `func Resolve(in Input) (Result, error)`, one parameter — and
  the checked table is the one already reached through the input,
  `in.Table`; the precondition adds no parameter and takes no
  table argument of its own." — (NC)

- [REQ-18] `0009:C3` "A table containing a row that
  breaches the conformance predicate above MUST cause Resolve
  to return a non-nil Go error identifying the offending row by
  RuleID and SourceLocator, with no Result disposition." — (NC)

- [REQ-19] `0009:C3` "The precondition is evaluated over the WHOLE table before any
  evaluation step: a breach surfaces even when no resolution
  path reaches the offending row, and it precedes every
  modeled disposition (a malformed table is malformed as a
  value, independent of the input tuple)." — (NC)

- [REQ-20] `0009:C3` "The breach travels
  the error path RDR 0001 reserves for programmer mistakes:
  it MUST NOT be a modeled refusal, MUST NOT introduce a
  refusal kind, and MUST NOT alter any disposition of a
  conforming table." — (NC)

- [REQ-21] `0009:D-selection-predicate` "*Precedence* — the breach error precedes
  every modeled disposition, `unmodeled_outcome` included." — (LBD)

- [REQ-22] `0009:D-selection-predicate` "*Placement* — the call is `Resolve`'s
  first statement, `in.Table.CheckValid()`, above the existing
  `view := assemble(in)`, so a breaching
  table costs no view assembly. This is a **cheapness
  preference, not an observable contract**" — (LBD)

- [REQ-23] "When any row carries both a non-empty `Escape` and a
  non-empty `Writes`, `Resolve` returns the zero `Result`
  (`Plan` and `Refusal` both nil — `Result` is a struct, so
  there is no nil `Result` to return) and a non-nil error
  naming every such row." — (TD)

- [REQ-24] "Because the zero `Result` reports
  `Refused() == false`, callers MUST check the error before
  reading the disposition; a caller that branches on
  `Refused()` first would read a success-shaped value with a
  nil `Plan`." — (TD)

## D. The typed error surface

- [REQ-25] `0009:C4` "The breach error MUST be a TYPED error carrying the offending
  row identity as the kernel's existing RowRef value, inspectable
  via errors.As / errors.AsType without parsing message text, and
  MUST wrap a package-level sentinel naming the breach category so
  errors.Is can classify it." — (NC)

- [REQ-26] `0009:C4` "Both the error type and the sentinel
  MUST be EXPORTED" — (NC)

- [REQ-27] `0009:C4` "Row identity MUST NOT be recoverable
  only from formatted prose." — (NC)

- [REQ-28] `0009:C4` "the sentinel is ErrEscapeShapeBreach (package-level, exported);" — (NC)

- [REQ-29] `0009:C4` "the typed error is *EscapeShapeBreachError (exported), carrying
    exactly ONE RowRef field, Ref, plus a Count int field (see the
    multi-breach clause);" — (NC)

- [REQ-30] `0009:C4` "its Unwrap() error returns ErrEscapeShapeBreach, so errors.Is
    classifies EVERY per-row element, not only the aggregate;" — (NC)

- [REQ-31] `0009:C4` "the predicate is Table.CheckValid() error." — (NC)

- [REQ-32] `0009:C4` "These four spellings collide
  with no frozen boundary guard: the kernel's banned exported-name
  lists (REQ-25, REQ-36, REQ-37) are exact-match, and errors/fmt are
  forbidden by no import guard (REQ-26, REQ-28, REQ-37)." — (NC)
  *(the REQ-nn ids in this quote are RDR 0001's, not this list's.)*

- [REQ-33] `0009:C4` "Because errors.Join returns a wrapper even for a single error,
  callers MUST classify and extract with errors.Is / errors.As /
  errors.AsType rather than a direct type assertion or equality
  against the sentinel." — (NC)

- [REQ-34] `0009:C4` "The aggregate is the uniform return shape:
  Resolve MUST NOT return the bare per-row error in the
  one-breach case and the aggregate otherwise, so caller code has
  one shape to handle regardless of breach count." — (NC)

- [REQ-35] `0009:C4` "Unwrap() []error MUST be EXACTLY
  ONE LEVEL deep, every element being a *EscapeShapeBreachError —
  never a nested join." — (NC)

- [REQ-36] `0009:C4` "Resolve MUST return CheckValid's error
  VERBATIM, never fmt.Errorf-wrapped: an added layer would expose
  Unwrap() error at the outermost level and break the flat
  traversal this clause guarantees." — (NC)

- [REQ-37] `0009:D-naming` "Spellings pinned (Pre-Lock, 3amigo T-2): sentinel
  `ErrEscapeShapeBreach`; typed error `*EscapeShapeBreachError`
  with fields `Ref RowRef` and `Count int`; predicate
  `Table.CheckValid() error`." — (LBD)

- [REQ-38] `0009:F1` "A single
  `errors.As`/`AsType` call reports only the FIRST breach in
  the chain, so it is the wrong instrument for reading a
  multi-breach report — a test asserting on all offending
  rows must traverse the aggregate." — (FM)

## E. Multi-breach reporting, ordering, and collapse

- [REQ-39] `0009:C5` "When a table carries MORE THAN ONE breaching row, the
  precondition MUST report every one of them in a single pass —
  not the first alone — combined with errors.Join, so each
  per-row error stays individually inspectable through the
  aggregate's Unwrap() []error." — (NC)

- [REQ-40] `0009:C5` "The reported rows MUST be
  ordered by RowRef identity using the kernel's existing
  compareRefs ordering, never by Table.Rows position, so the
  diagnostic payload is a function of the table value rather
  than of row order" — (NC)

- [REQ-41] `0009:C5` "Table.Rows position MUST NOT be used to break
  the tie." — (NC)

- [REQ-42] `0009:C5` "The report is instead ordered by RowRef identity and made
  total by COLLAPSING equal identities: breaching rows sharing
  one RowRef contribute ONE reported error, not one per row." — (NC)

- [REQ-43] `0009:C5` "Rows carrying no source identity collapse to a single
  RowRef{"",""} entry; that is a degenerate producer, and the
  error MUST remain diagnostic in that case by stating the
  breach count alongside the identities." — (NC)

- [REQ-44] `0009:C5` "The count is carried STRUCTURALLY, as the Count field on each
  per-identity *EscapeShapeBreachError — never only in formatted
  prose" — (NC)

- [REQ-45] `0009:C5` "Count is PER-IDENTITY and counts PRE-COLLAPSE ROWS: three
  breaching rows sharing RowRef{"",""} yield ONE reported error
  with Ref == RowRef{"",""} and Count == 3. A single-row breach
  carries Count == 1, so the field is uniform rather than present
  only in the degenerate case." — (NC)

- [REQ-46] `0009:D-selection-predicate` "*Reporting* — aggregate, not fail-fast: all
  breaching rows in one pass, sorted by `compareRefs`, with
  equal-identity rows collapsed to one entry" — (LBD)

## F. The exported predicate (`Table.CheckValid`)

- [REQ-47] `0009:C6` "The conformance predicate MUST be exported by the kernel
  package as a construction-time check callable by any table
  producer, and Resolve's entry precondition MUST be that same
  function — one predicate, two call sites, so a non-TOML
  producer can fail at build or construction time rather than
  at first production Resolve, and the two enforcement points
  cannot drift." — (NC)

- [REQ-48] `0009:C6` "It MUST be a METHOD ON Table taking no
  arguments" — (NC)

- [REQ-49] `0009:C6` "It MUST
  return error (nil when valid), and its name MUST follow Go's
  error-returning convention — CheckValid or Validate, never
  Valid/IsValid/OK, which Go reserves for bool-returning
  predicates." — (NC)

- [REQ-50] `0009:C6` "copy the latter's doc
  form — "returns nil if <x> is valid, or else an error
  describing a problem."" — (NC)

- [REQ-51] `0009:C6` "The doc comment MUST also name the ONE
  property checked — escape-row shape conformance — and state
  that nothing else is checked (in particular not the
  Escape-class restriction to no_match/ambiguous_match that
  Row's doc records for RDR 0002), so a nil return is never
  read as general table validity." — (NC)

- [REQ-52] `0009:C6` "A table with no rows (empty or nil Rows) conforms VACUOUSLY:
  CheckValid MUST return nil, since no row can breach a predicate
  quantified over rows and errors.Join of nothing is nil." — (NC)

- [REQ-53] `0009:D-selection-predicate` "*Locus* — one exported function, called at
  both sites, so any future change is a single-site edit that cannot
  desynchronize the two enforcement points." — (LBD)

## G. Doc-contract amendments authorized by this RDR

- [REQ-54] `0009:D-selection-predicate` "*Doc-contract amendment* —
  `resolve.go::Resolve`'s numbered evaluation-order list (a REQ-1 contract
  artifact) currently opens at the alphabet check, so the whole-table
  precondition MUST be prepended to it as a new step 0." — (LBD)

- [REQ-55] `0009:D-selection-predicate` "The same authorization covers the two
  disposition-cardinality doc contracts the breach return would otherwise
  contradict — `Result`'s "never both and never neither" and `Resolve`'s
  "returns exactly one disposition" — each amended to scope
  itself to nil-error returns." — (LBD)

- [REQ-56] "state the producer
  obligation on the `Row` doc contract" — (IP, Phase 1)

- [REQ-57] "record the
  kernel's validation surface — exactly which shape property
  it checks vs still assumes — so partial validation cannot be
  read as general kernel ownership." — (IP, Phase 1)

## H. CLI surfacing (C7) — **EXECUTABLE** (the `flow` verb ships)

Every REQ in this section was previously marked DEFERRED on the false premise
that no `flow` verb existed at HEAD. The premise is withdrawn; all of H is
in scope. Numbers are unchanged.

**Grounding — the call site and the carrier both already exist.**
`internal/cli/flow_resolve.go:105` calls `resolve.Resolve`, and its `err != nil`
branch today reads:

```go
return respond.Fail(cmd, internalErr(codeAccessorFailed,
    "the resolution kernel reported a programmer error: "+err.Error()))
```

That branch is the C7 defect in the concrete: one generic code for every kernel
error, and the row identities reachable only by re-parsing `err.Error()` prose —
the re-parse this RDR forbids everywhere else. C7's obligation is to
discriminate the breach out of that branch. The established pattern to model on
is `flow_input.go`'s typed-code vocabulary (`codeTagInvalid`,
`codeWriteInvalid`, …) plus the `Findings` carrier RDR 0008 already uses.

- [REQ-58] `0009:C7` "When a breach reaches the CLI, the verb MUST wrap it into a
  *clierr.CLIError carrying a stable Code, Group GroupInternal
  (exit 2), the offending row identity, and a Hint stating the
  remedy" — (NC)
  *Site*: `internal/cli/flow_resolve.go`'s post-`Resolve` error branch, which
  must classify the breach (`errors.Is(err, resolve.ErrEscapeShapeBreach)`)
  before falling through to `codeAccessorFailed`.

- [REQ-59] `0009:C7` "Because clierr.CLIError.Cause is
  json:"-", the Go error chain that carries the RowRef values
  is NOT wire-visible: the offending row identities MUST reach
  the envelope through a serialized field." — (NC)
  *Confirmed at HEAD*: `CLIError.Cause` is `json:"-"`; `EmitJSON` marshals the
  `*CLIError` with no custom `MarshalJSON`.

- [REQ-60] `0009:C7` "The chosen carrier is a
  NEW omitempty field added under the type's own "Extend with
  new optional fields as needed" allowance — NOT the existing
  Detail" — (NC)
  **Already satisfied at HEAD, do not re-add.** The carrier ships as
  `clierr.CLIError.Findings []Finding` (`json:"findings,omitempty"`), landed by
  RDR 0008 per JDR 0001 §D10 — which `deviations.md` D1 names as "the carrier
  for this RDR's row identities". The obligation this REQ still binds is to
  **use** `Findings`, not `Detail`, and not to mint a second parallel field.

- [REQ-61] `0009:C7` "The clause binds the carrier CLASS only; the
  field's NAME stays the implementer's." — (NC)
  Resolved by HEAD: the name is `Findings`, already fixed by the shipped type.

- [REQ-62] `0009:C7` "the field's type MUST be a clierr-local
  representation (plain strings or a small row-identity struct
  declared in clierr) and MUST NOT be resolve.RowRef, keeping
  clierr a leaf that does not import internal/resolve — the CLI
  verb layer, which already imports both, does the conversion" — (NC)
  *Satisfied in shape by HEAD*: `clierr.Finding` is clierr-local and its fields
  are plain strings. The live obligation is the mapping — the verb layer
  converts each `*EscapeShapeBreachError`'s `RowRef{RuleID, SourceLocator}` into
  a `Finding` (`Rule`/`Param` for the rule id, `Locator` for the source
  locator), and `clierr` must gain no `internal/resolve` import.

- [REQ-63] `0009:C7` "the per-identity Count MUST serialize alongside the
  identities" — (NC)
  `Finding` carries no `Count` field at HEAD. See **Q1** — the count must reach
  the wire on the `Finding` record without a prose re-parse.

- [REQ-64] `0009:C7` "The stable
  Code is "escape-row-shape-breach"." — (NC)
  Note the spelling is fixed by the contract and is **not** `flow`-prefixed,
  unlike `flow_input.go`'s existing vocabulary. The contract's literal wins.

- [REQ-65] `0009:C7` "the Cause comment's now-stale second clause MUST be
  amended in the same change (a Prerequisite; amending a stale
  code comment is not reopening RDR 0005's envelope contract)." — (NC)
  **Pre-existing deviation, not this RDR's to introduce**: the `Cause` comment
  at `internal/cli/clierr/clierr.go` still reads "the wire-visible cause surface
  is Detail" while `Findings` already ships beside it, so the sentence is
  *already* false at HEAD. This RDR does not create the staleness; it inherits
  it. Repair is in scope per REQ-66. If a phase judges the repair to belong to
  RDR 0008's landing rather than here, that is a recorded deviation, not an
  edit dropped silently.

- [REQ-66] "The same change that adds the `omitempty` identity
  field to `clierr.CLIError` amends the now-stale second
  clause of `internal/cli/clierr/clierr.go::CLIError.Cause`'s
  doc comment ("Not serialized — the wire-visible cause
  surface is Detail"), which the addition makes false, and
  updates `docs/cli-output-contract.md`'s error-envelope
  field list alongside." — (PRE)
  `docs/cli-output-contract.md` already documents `findings` (its "Structured
  findings" section), so the doc half is largely discharged; the `Cause` comment
  half is not.

- [REQ-67] `0009:C7` "conformance MUST be
  asserted on the Code, never on the exit code alone." — (NC)
  Directly testable now, in the manner of the RDR-0005 `flow_*` suites and
  `internal/cli/reserved_key_0008_test.go`.

- [REQ-68] `0009:C7` "Adding this
  Code does NOT open RDR 0001's refusal taxonomy: the five
  refusal kinds stay closed and gain no member" — (NC)

- [REQ-69] "This is an additive envelope change, not a change to RDR
  0005's refusal-to-exit-code mapping" — `0009:C7` (NC)
  Strengthened by HEAD: the field is not even added — it already ships — so
  nothing in the envelope changes at all.

- [REQ-109] "A breach surfaced through the CLI. **Expected**: exit 2 via
  `CLIError{Group: GroupInternal}` carrying `Code: "escape-row-shape-breach"`,
  the offending row identity, and the remedy `Hint`" — (TS 11, restated as an
  executable C7 obligation)
  The record defers *scenario 11* by name (REQ-91), but C7's clauses are
  normative and executable, so the CLI assertion lands here instead: an
  end-to-end `flow` test asserting the `Code`, the `Findings` identities, and
  the `Hint`, in the style of the existing `flow_*_0005_test.go` suites. This is
  the only REQ minted by this amendment.

## I. Executable test scenarios

- [REQ-70] `0009:MVV` "A hand-constructed table containing an escape row with
  populated `Writes`: `Resolve` returns a non-nil error naming
  that row's `RuleID` and `SourceLocator`, with no `Plan` and
  no `Refusal` — while the same table with the writes removed
  resolves, escapes, and refuses exactly as the frozen ADV
  suite proves today (zero disposition change for conforming
  tables). The exported validator and the entry check are
  exercised as the same predicate." — (MVV) — **REQ-MVV**

- [REQ-71] `0009:MVV` "The three scenarios below are **normative fixtures**
  (approved at Stage 4). Their values are read from the frozen
  suite and the A3 spike runs, not invented" — (MVV)

- [REQ-72] `0009:MVV` "**breach-yields-error-not-plan** — an escape row
  (`RuleID` `rdr.escape.needsowned`, `SourceLocator`
  `flows/rdr.toml:90`, `Escape` `[no_match]`) carrying
  `Writes: []Tag{{Key: "status", Value: "Escaped"}}`, in a
  table where it is selected via the escape path.
  **Expected**: `Resolve` returns the zero `Result`
  (`Plan: nil, Refusal: nil`) and a non-nil error; the error
  yields that row's `RowRef{"rdr.escape.needsowned",
  "flows/rdr.toml:90"}` through `errors.As`/`AsType`
  without parsing text, and `errors.Is` classifies it as
  `ErrEscapeShapeBreach` — both through the `errors.Join`
  aggregate, which wraps even this single-breach case. The
  aggregate's `Unwrap() []error` has length **1**; its single
  `*EscapeShapeBreachError` carries `Count == 1`." — (MVV / TS 1)

- [REQ-73] `0009:MVV` "**dormant-row-still-errors** — a table that resolves
  cleanly on its own, plus an escape row no resolution path
  reaches (an outcome outside the tuple's reach) carrying
  the same `Writes`.
  **Expected**: the same error, naming the dormant row —
  pinning the whole-table decision against an on-path-only
  reading." — (MVV / TS 2)

- [REQ-74] `0009:MVV` "**empty-not-nil-conforms** — an escape row whose `Writes`
  is a non-nil, zero-length slice (`[]Tag{}`).
  **Expected**: no error; the table resolves, escapes, and
  refuses identically to the A3 baseline, and the emitted
  plan carries `len(Plan.Writes) == 0`. Asserted on length,
  never on nil-ness" — (MVV / TS 3)

- [REQ-75] `0009:S3` "The **discriminating** assertion is on the
  *input row* (`row.Writes != nil && len(row.Writes) == 0`
  yields no error): asserting only on the plan would not
  separate the two states this scenario exists to separate" — (TS 3)

- [REQ-76] `0009:S4` "Breach precedence over each modeled
  disposition, run once **per refusal kind** — a breaching
  row coexisting with a table state that would otherwise
  yield `unmodeled_outcome`, `no_match`, `ambiguous_match`,
  `owned_state_unavailable`, and `guard_unevaluable`
  respectively.
  **Expected**: The breach error wins in every case, with no
  `Result` disposition. Exhaustive over the five kinds" — (TS 4)

- [REQ-77] `0009:S5` "The full frozen `internal/resolve` suite
  (ADV/MVV/boundary) run against conformed escape fixtures.
  **Expected**: Every assertion still passes." — (TS 5)

- [REQ-78] `0009:S5` "The comparison baseline is the A3 spike's **conformed** run
  (`a3-conformed.out`: 154 PASS, 0 FAIL), restricted to the
  tests existing at that revision — **not** `a3-baseline.out`
  (the pre-conformance tree) and not the raw post-Phase-2
  count, which necessarily exceeds 154 because this RDR adds
  tests. The oracle is the sorted outcome set over that
  pre-existing test set, which must be identical." — (TS 5)

- [REQ-79] `0009:S5` "The baseline is regenerable,
  not frozen to the spike artifact: re-run the frozen suite
  at the implementation's base commit immediately before
  Phase 1 lands and diff sorted outcome sets over the tests
  existing at that commit" — (TS 5)

- [REQ-80] `0009:S6` "Mutation check — the check's non-vacuity.
  Two named mutants, run separately, with this RDR's **own**
  new tests (scenarios 1–4, 7, 9, 10, 10b) **excluded** from
  the oracle" — (TS 6)

- [REQ-81] `0009:S6` "*Mutant A* — `Table.CheckValid` returns `nil`
  unconditionally. **Expected**: the frozen suite still
  passes." — (TS 6)

- [REQ-82] `0009:S6` "*Mutant B* — re-introduce the breach into the fixtures
  (restore `Writes` on `fixtures_test.go::escapeRow`) with
  the entry check intact. **Expected**: the suite fails,
  proving the check is wired into the real evaluation path
  and not dead code." — (TS 6)

- [REQ-83] `0009:S7` "The exported predicate called directly by a
  producer at construction time, on the same tables as
  scenarios 1–3.
  **Expected**: Verdicts identical to `Resolve`'s entry check
  — one predicate, two call sites, no drift. "Identical" is
  defined as: both nil, or both non-nil with equal extracted
  `[]RowRef` **and** equal per-identity `Count` values, in
  the same order. Not `reflect.DeepEqual` over the two error
  values" — (TS 7)

- [REQ-84] "**Scenario**: The zero-`Result` caller trap on a breach.
   **Expected**: `Resolve` returns `Result{Plan: nil, Refusal:
   nil}` alongside the non-nil error, so `Refused() == false`
   on a call that did not succeed. Asserted directly" — (TS 7b)

- [REQ-85] `0009:S9` "A table with two or more breaching rows.
  **Expected**: One `Resolve` call reports **every** breaching
  identity, not just the first; each per-row error is
  individually recoverable by traversing the aggregate's
  `Unwrap() []error`" — (TS 9)

- [REQ-86] `0009:S9` "The traversal is **exactly one level deep** and every element
  type-asserts to `*EscapeShapeBreachError`; `errors.Is(elem,
  ErrEscapeShapeBreach)` holds on **each element**, not only
  on the aggregate. The assertion compares the extracted
  `[]RowRef` — never `reflect.DeepEqual` over two
  `errors.Join` values" — (TS 9)

- [REQ-87] `0009:S10` "The same multi-breach table, rows carrying
  *distinct* `RowRef` identities, supplied in two different
  orders.
  **Expected**: Identical reported row sequence — ordered by
  `RowRef` identity (`compareRefs`), never by `Table.Rows`
  position." — (TS 10)

- [REQ-88] `0009:S10` "The table MUST be **mixed** — at least one row carrying
  source identity and at least one built without it
  (`RowRef{"",""}`) — so the report's ordering is pinned
  across the two identity classes rather than only within
  one." — (TS 10)

- [REQ-89] `0009:S10` "A multi-breach table whose breaching rows
  share one `RowRef` identity (including the zero-value
  `RowRef{"",""}` of rows built without source identity),
  supplied in two different orders.
  **Expected**: Identical report under both permutations —
  equal identities collapse to one reported entry. For three
  rows sharing `RowRef{"",""}`: `Unwrap() []error` has
  length **1**, its single element has `Ref ==
  RowRef{"",""}` and `Count == 3` (pre-collapse rows, read
  off the struct field — no assertion parses message text)." — (TS 10b)

- [REQ-90] `0009:S8` "(binds RDR 0002's build): authored escape rule
  carrying a write block, one carrying a clear list, and one
  carrying an *empty* write block (`writes = []`).
  **Expected**: All three rejected at load under "malformed
  escape declaration," naming the source rule; normalized
  escape rows render write-free." — (TS 8)
  *(DEFERRED by the record's own TS scope note, which places scenario 8 outside
  this RDR's Done criteria as binding "RDR 0002's build". Note RDR 0002 is
  `Implemented` at HEAD (`internal/table/`), so the deferral is one of
  OWNERSHIP, not readiness — the record assigns the assertion to 0002's build,
  and this RDR discharges its side as the Phase 3 shared fixtures, REQ-97/98.)*

- [REQ-91] `0009:S11` "— **DEFERRED: not executable by this RDR**
  (binds the future verb that first calls `Resolve`; no
  `flow` verb exists at HEAD, per A2/A5). A breach surfaced
  through the CLI.
  **Expected**: exit 2 via `CLIError{Group: GroupInternal}`
  carrying `Code: "escape-row-shape-breach"`, the offending
  row identity, and the remedy `Hint` "fix the table
  producer: an escape row must carry no writes"" — (TS 11)
  *(DEFERRED because the record itself marks this scenario "not executable by
  this RDR". Its stated reason — "no `flow` verb exists at HEAD" — is FALSE at
  this commit, but the deferral is the record's scope decision and RDRs are
  never amended. The substance is not lost: C7's clauses REQ-58 … REQ-69 are
  executable and REQ-109 carries the CLI assertion.)*

- [REQ-92] "Done means, in user terms: **an escape row can no longer carry
  an owned-tag write to the accessor layer, so no tag value can
  appear in owned state that no authored rule set**" — (TS preamble)

## J. Fixture conformance (Phase 2)

- [REQ-93] "Bring the write-bearing escape rows into conformance at all
  three sites: drop the `Writes` field from the
  `internal/resolve/fixtures_test.go::escapeRow` builder (line
  257, inherited by its sixteen call sites), and drop the two
  call-site overrides —
  `internal/resolve/adversarial_test.go:229` (ADV-2b) and
  `internal/resolve/fixup_test.go:103` (Fixup-1e "missing
  owned state")." — (IP, Phase 2)

- [REQ-94] "`NextTags` stays on the builder: A4 did not
  widen the predicate." — (IP, Phase 2)

- [REQ-95] "Phases 1 and 2 land as ONE change: the moment the entry check
  exists, every write-bearing escape fixture errors at `Resolve`
  entry and `mustResolve` fatals across the frozen suite" — (IP, Phase 1)

- [REQ-96] "The package gains `errors` (and
  `fmt` for the message); both clear the frozen import guards." — (IP, Phase 1)

## K. Phase 3 — shared normalizer conformance fixtures

- [REQ-97] "Encode the authored-path half as a named, shared fixture set
  binding the future RDR 0002 build" — (IP, Phase 3)

- [REQ-98] "one canonical authored-clear case pins the
  `<clear>`-sentinel-write representation identically for the
  kernel suite and the normalizer suite, so the two enforcers
  of one invariant cannot drift." — (IP, Phase 3)

## L. Performance and non-goals

- [REQ-99] "The precondition adds one O(rows) pass over `Table.Rows` per
  `Resolve` call." — (PE)

- [REQ-100] "The scan is allocation-free in the conforming case
  — it reads `len(row.Escape)` and `len(row.Writes)` and
  allocates only when a breach is found and a `RowRef` is
  recorded" — (PE)

- [REQ-101] "The cost is per-call because the kernel is stateless by
  contract; no caching or memoization is introduced" — (PE)

- [REQ-102] "No byte-stable hash, canonical serialization, or ordering
  guarantee is introduced, so the determinism checklist does
  not apply." — (PE)

- [REQ-103] `0009:F3` "**Silent failure guarded against**: an escape plan
  carrying writes no rule sanctioned … Under the
  precondition it cannot be emitted; the conformance
  fixtures are the tripwire against regression." — (FM)

- [REQ-104] `0009:F5` "**Error-channel creep guarded against**: this RDR's
  normative clause states what qualifies for the error path
  (a producer/programmer contract breach, never a modeled
  condition)" — (FM)

- [REQ-105] "Positive: the five-kind taxonomy, RDR 0005's mapping, RDR
  0002's grammar, and the kernel's type vocabulary all stand
  unchanged; the check reuses a shipped discriminator and a
  reserved channel." — (CONS)

- [REQ-106] "Negative (behavior change, deliberate): a hand-built table
  carrying a *dormant* malformed row — one no resolution
  path reaches — previously resolved fine and now errors on
  every `Resolve` (whole-table precondition)." — (CONS)

- [REQ-107] `0009:EIA` "Escape fixtures … | Extend (conform all three sites) |
  Fixtures must satisfy the precondition; supersedes the "Noted, not filed"
  cleanup" — (EIA)

- [REQ-108] `0009:EIA` "Row-identity diagnostics | `internal/resolve/resolve.go::RowRef`
  + `rowRefs`/`compareRefs` | … | Reuse | The typed breach error carries
  `RowRef`s in the same sorted order" — (EIA)

---

## ASSUMPTIONS

- **ASSUMPTION (REQ-4 / REQ-31 / REQ-52):** `CheckValid` is a method on
  `Table` with a **value** receiver (`func (t Table) CheckValid() error`),
  not a pointer receiver. C6 says "a METHOD ON Table taking no arguments"
  and REQ-22's call site is `in.Table.CheckValid()` where `in.Table` is a
  `Table` value field of `Input`; both precedents cited
  (`pprof.Profile.CheckValid`, `rsa.PrivateKey.Validate`) are pointer
  receivers on pointer-held types, but `Table` is held by value here. Either
  receiver compiles at the call site; the value receiver is chosen as it
  matches the package's value-oriented types (`TagSet.Lookup`,
  `Result.Refused`).

- **ASSUMPTION (REQ-29 / REQ-45):** `Count` counts breaching **rows** sharing
  the identity, not breaching *writes*. C5 states "counts PRE-COLLAPSE ROWS"
  and gives the worked example (three rows sharing `RowRef{"",""}` →
  `Count == 3`), so a row carrying two illegal writes still contributes 1.

- **ASSUMPTION (REQ-35 / REQ-39 / REQ-42):** the collapse happens **before**
  `errors.Join`, i.e. one `*EscapeShapeBreachError` is constructed per
  distinct `RowRef` and the joined slice holds exactly those. C5's "breaching
  rows sharing one RowRef contribute ONE reported error" plus C4's "Unwrap()
  []error MUST be EXACTLY ONE LEVEL deep, every element being a
  *EscapeShapeBreachError" admit no other construction.

- **ASSUMPTION (REQ-36 / REQ-52):** `Resolve` returns `CheckValid`'s error
  verbatim, so on a conforming table `CheckValid` returns a plain untyped
  `nil` (never a typed-nil `*EscapeShapeBreachError`), and `Resolve`'s
  existing nil-error return paths are unchanged.

- **ASSUMPTION (REQ-30 / REQ-33):** `*EscapeShapeBreachError` implements
  `Unwrap() error` (singular) returning the sentinel, and does **not** also
  implement `Unwrap() []error` — a type may not usefully have both, and C4
  assigns the plural form to the `errors.Join` aggregate only.

- **ASSUMPTION (REQ-96):** the error message text is unspecified by this RDR
  beyond "MUST remain diagnostic … by stating the breach count alongside the
  identities" (REQ-43). No test may assert on it (REQ-27, REQ-44, REQ-89), so
  the wording is the implementer's.

- **ASSUMPTION (REQ-77 / REQ-78):** "the frozen `internal/resolve` suite"
  means the test files existing at the implementation's base commit, and the
  oracle is regenerated at that commit per REQ-79 rather than diffed against
  the literal `a3-conformed.out` bytes.

- **ASSUMPTION (REQ-80 … REQ-82):** the mutation check is executed as a
  manual/recorded procedure during the build (revert, run, restore) and its
  outcome is recorded in the phase artifact — not shipped as a permanent test
  that mutates the source.

- **ASSUMPTION (Section H) — CORRECTED:** the C7 REQs **are built** in this
  RDR. The earlier reading (not built, no verb at HEAD) rested on the record's
  stale A2/A5 evidence and is withdrawn. `internal/cli/flow_resolve.go` calls
  `resolve.Resolve` today and its error branch is editable, so every C7 clause
  has a live site. The deviations file's D1 check
  (`grep -n '§D10' docs/rdr/0009-*.md`) remains a citation obligation on the
  record, run at its own phase; its substantive half — "row identities are
  carried in `findings`" — is REQ-60/REQ-62 here and is satisfiable, since
  `Findings` ships.

- **ASSUMPTION (REQ-22):** because placement is explicitly "a cheapness
  preference, not an observable contract", no test asserts that `assemble`
  was not called; the observable obligation is REQ-19/REQ-21 (whole-table
  scope and precedence), which is what scenarios 2 and 4 pin.

- **ASSUMPTION (REQ-90 / REQ-97 / REQ-98):** Phase 3's fixture set is a
  data/spec artifact under this RDR's tree (a named fixture file plus its
  expectations). REQ-14 calls it "a conformance fixture set" and the
  Consequences say "Phase 3 ships the binding fixtures, not the enforcement."
  **Note the record's stated reason is stale**: RDR 0002 is `Implemented` at
  HEAD (`internal/table/`), so a fixture set *could* be executed against the
  real normalizer. The record nonetheless scopes Phase 3 to shipping the
  fixtures and defers the assertion to 0002's build (REQ-90). Proceeding on the
  record's scoping; if a later phase finds executing them against
  `internal/table/` is free, running them is a strengthening, not a scope
  breach.

---

## QUESTIONS

- **Q1 — RE-EVALUATED (was: "does `Count` serialize only through the deferred
  CLI field?"). How does the per-identity `Count` reach the wire, given
  `clierr.Finding` carries no `Count` field at HEAD?**
  The original reading deferred this to a successor build. That is withdrawn:
  the `flow` verb ships, so REQ-63 ("the per-identity Count MUST serialize
  alongside the identities") is executable **now** and must be answered here.
  The kernel side is settled — C4 pins `Count int` as an exported field on
  `*EscapeShapeBreachError` and REQ-89 reads it off the struct. The wire side
  is genuinely open, because `Finding` (`internal/cli/clierr/clierr.go`) has
  `Code, Message, Model, Severity, Param, Locator, Hint, Rule, Span, Element,
  Reason, Dimension, Key, Operator, Literal, Block, Class, Fingerprint` — and
  no count. Three readings:
  (a) add a `Count int json:"count,omitempty"` field to `clierr.Finding`;
  (b) emit `Count` per-row by repeating the `Finding` N times;
  (c) render the count into the `Finding`'s `Message`/`Hint` prose.
  **Proceeding under (a).** (c) is excluded outright — REQ-44 says the count is
  carried "never only in formatted prose", and REQ-27/REQ-86 forbid recovering
  structured data from message text. (b) is excluded by REQ-42: equal
  identities collapse to ONE reported entry, and repeating the record
  reintroduces the per-row multiplicity collapsing exists to remove. (a) is
  additive under the same `CLIError`/`Finding` "omitempty" allowance A9 cites
  and that `Findings` itself landed under, and `Count == 1` for the ordinary
  single-breach case means the key elides on every non-degenerate path.
  **Residual risk**: `Finding` is a shared record owned jointly by RDRs
  0005/0006/0008; adding a field touches a surface those RDRs pin. The
  0008 test that pins it
  (`reserved_key_0008_test.go::TestReq102_TheFindingCarrierHasTheFiveFieldsIncludingHint`)
  asserts the *presence* of five named fields, not an exact field set, so an
  addition does not break it — but a phase that finds an exact-set assertion
  elsewhere should record a deviation rather than weaken that test.

- **Q4 — Does the breach get its own `Code` constant in `flow_input.go`'s
  vocabulary, and does the `flow`-prefix convention apply?**
  Every shipped code is `flow`-prefixed (`flow-tag-invalid`,
  `flow-write-invalid`, …), but REQ-64 fixes this one literally as
  `"escape-row-shape-breach"`, unprefixed. **Proceeding on the contract's
  literal** — C7 pins the exact string and REQ-67 makes the `Code` the
  conformance oracle, so a prefixed variant would fail the contract. Recorded
  because the inconsistency with the neighbouring vocabulary is visible and a
  reviewer may read it as a typo; it is not.

- **Q2 — Are `errors` and `fmt` both required imports, given no test may read
  message prose?** REQ-96 says the package gains "`errors` (and `fmt` for the
  message)". `errors` is load-bearing (`errors.Join`, sentinel via
  `errors.New`). `fmt` is needed only if `Error()` formats the identity and
  count into prose — which REQ-43 requires ("MUST remain diagnostic … by
  stating the breach count alongside the identities"). **Proceeding**: both
  imports land, `fmt` used solely inside `Error()`. Recorded because REQ-32's
  import-guard clearance is quoted for both and a build that omits `fmt`
  would leave REQ-43 unsatisfiable.

- **Q3 — Is REQ-52's vacuous-conformance clause reachable through `Resolve`,
  or only through the direct `CheckValid` call?** RDR 0001's REQ-20
  empty-input region is cited as the reason the shape is admitted. Two
  readings: (a) a nil/empty-`Rows` table reaches `Resolve` and must produce
  its existing disposition unchanged (nil error from the precondition, then
  the normal path); (b) the clause is only about the exported predicate.
  **Proceeding under (a)** — REQ-47 makes the two call sites the same
  function and REQ-20 forbids altering "any disposition of a conforming
  table", and a rowless table is conforming. No predecessor contradicts this;
  recorded because a nil return there must never be read as "validation
  skipped" (REQ-52's own words).
