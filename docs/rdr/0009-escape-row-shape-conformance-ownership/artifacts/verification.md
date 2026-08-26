# Verification — RDR 0009 Escape-row shape conformance ownership

## Phase 3a — CoVe verification (independent, spec-derived)

Method: violating inputs were derived from `0009-escape-row-shape-conformance-ownership.md`
and `artifacts/req-list.md` **alone**. The Phase 1 suites
(`internal/{resolve,table,cli}/escape_shape*_0009_test.go`) were not read, per
the independence constraint on this phase. Every input below was **executed**
against the real implementation through throwaway probe harnesses in
`internal/resolve`, `internal/table`, and `internal/cli` (deleted afterwards),
plus the built `intrastate` binary. Roughly 90 distinct inputs were run.

**Result: no FAIL-N entries. Every executable REQ held under probing.**

### What was probed, and what held

**A/C/D — the producer rule, the entry precondition, the typed error**
(REQ-1..REQ-10, REQ-16..REQ-38)

- Single write-bearing escape row, entered through `Resolve`: returns the zero
  `Result` (`Plan == nil`, `Refusal == nil`, `Refused() == false`) with a
  non-nil error. `errors.Is(err, ErrEscapeShapeBreach)` holds; `errors.As` and
  `errors.AsType[*EscapeShapeBreachError]` both recover
  `RowRef{"rdr.escape.needsowned", "flows/rdr.toml:90"}` with `Count == 1`,
  with no prose parsed. (REQ-18, REQ-23, REQ-24, REQ-25, REQ-72, REQ-84)
- The single-breach case returns the **aggregate**, not the bare per-row error:
  a direct `err.(*EscapeShapeBreachError)` assertion **fails**, and
  `err == ErrEscapeShapeBreach` is false — the uniform shape REQ-33/REQ-34
  demand. `Unwrap() []error` has length 1. (REQ-33, REQ-34, REQ-72)
- `Resolve`'s return exposes **no** `Unwrap() error` at the outermost level and
  is byte-identical in message to the direct `CheckValid()` return — verbatim,
  never `fmt.Errorf`-wrapped. (REQ-36)
- Nil-vs-empty: a row with `Writes: []Tag{}` (non-nil, zero length) and
  `Escape` non-empty produces **no** error. Asserted on the *input row*
  (`row.Writes != nil && len(row.Writes) == 0`), which is the discriminating
  assertion REQ-75 requires. (REQ-5, REQ-74, REQ-75)
- Escape-adjacent non-escapes: a row with `Writes` and `Escape == nil`, and one
  with `Writes` and `Escape == []RefusalKind{}`, both conform. An escape row
  carrying only `NextTags` conforms — the predicate does not extend to
  `NextTags`. (REQ-3, REQ-4, REQ-94)
- A `<clear>`-sentinel value in `Writes` on an escape row breaches, so the one
  Writes-only predicate carries "no clears" at the kernel boundary. (REQ-2)
- The discriminator is `len(Escape) != 0` only: an escape row whose class is
  outside RDR 0002's `no_match`/`ambiguous_match` allowlist (and one carrying a
  garbage kind) still breaches when it carries writes, and still **conforms**
  when write-free — confirming `CheckValid` checks the one property and nothing
  else. (REQ-4, REQ-6, REQ-7, REQ-51)
- No stripping path: a conforming escape row selected via the escape path emits
  `Plan{Escaped: true}` with `len(Plan.Writes) == 0`, from an unconditional
  `planOf`. (REQ-8, REQ-9, REQ-10)

**Whole-table scope and precedence** (REQ-19, REQ-21, REQ-76, REQ-106)

- Dormant breach: an escape row no resolution path reaches still errors — and
  errors for **every** recognized outcome tried (`"a"`, `"b"`, an unmodeled
  value, and `""`). (REQ-19, REQ-73, REQ-106)
- Exhaustive over the five refusal kinds: for each of `unmodeled_outcome`,
  `no_match`, `ambiguous_match`, `owned_state_unavailable`, and
  `guard_unevaluable`, the base table was confirmed to produce that refusal,
  then a breaching row was injected. In all five the breach error won with no
  `Result` disposition and no refusal kind. (REQ-20, REQ-21, REQ-76)

**E — multi-breach ordering, collapse, Count** (REQ-39..REQ-46, REQ-85..REQ-89)

- Three distinct identities including a mixed zero-identity row, supplied in
  two different `Table.Rows` permutations: identical reported sequence both
  times, sorted by `compareRefs` (`RuleID` then `SourceLocator`), never by
  position. Confirmed the `SourceLocator` tie-break with `x/l9`, `x/l1`, `w/l5`
  → `w/l5`, `x/l1`, `x/l9`. (REQ-40, REQ-41, REQ-87, REQ-88)
- Same `RuleID`, different `SourceLocator` correctly does **not** collapse
  (two entries), so collapse keys on the full `RowRef`, not the rule id alone.
- Three rows sharing `RowRef{"",""}`: `Unwrap() []error` has length **1**, its
  single element carries `Ref == RowRef{"",""}` and `Count == 3`, read off the
  struct field. `Count == 1` on every single-row breach — uniform, not
  degenerate-only. (REQ-42, REQ-43, REQ-44, REQ-45, REQ-89)
- Traversal is exactly one level: every element type-asserts to
  `*EscapeShapeBreachError`, none carries a nested `Unwrap() []error`, and
  `errors.Is(elem, ErrEscapeShapeBreach)` holds on **each** element via
  `Unwrap() error` returning the sentinel — not only on the aggregate.
  (REQ-30, REQ-35, REQ-39, REQ-86)
- REQ-38's failure-mode claim was itself checked: a single `errors.As` over a
  two-breach aggregate does return only the first element, as the record says.

**F — the exported predicate, two call sites** (REQ-47..REQ-53, REQ-83)

- Agreement run over six tables (single breach, dormant breach, empty-not-nil,
  nil `Rows`, empty `Rows`, and a collapse case): `Table.CheckValid()` called
  directly and the error reached through `Resolve` agreed on every one — both
  nil, or both non-nil with equal extracted `[]RowRef` **and** equal
  per-identity `Count` in the same order. Compared on extracted values, not
  `reflect.DeepEqual` over the error pair. (REQ-47, REQ-53, REQ-83)
- Vacuous conformance: both `Table{}` (nil `Rows`) and `Table{Rows: []Row{}}`
  return a plain untyped `nil` (`err == nil` is literally true — no typed-nil),
  through both call sites. (REQ-52, and Q3's reading (a) is confirmed
  reachable through `Resolve`.)
- Signature and naming confirmed by compilation of the probes against the
  exported surface: `func (t Table) CheckValid() error`, no arguments, value
  receiver — matching the req-list ASSUMPTION. (REQ-31, REQ-37, REQ-48, REQ-49)

**B — the authored-path enforcer** (REQ-11..REQ-15)

Real TOML, loaded through `table.Load`, mutating the shipped
`internal/table/testdata/rdr-fixture.toml` escape rule. All five malformed
variants were rejected, each under category `malformed_escape_declaration`,
each naming the source rule:

| variant | rejected | category |
| --- | --- | --- |
| `[rule.write]` with keys | yes | `malformed_escape_declaration` |
| `[rule.write]` **empty** | yes | `malformed_escape_declaration` |
| `clear = ["prelock_lens"]` | yes | `malformed_escape_declaration` |
| `clear = []` **empty** | yes | `malformed_escape_declaration` |
| `gate = [...]` | yes | `malformed_escape_declaration` |
| unmodified (control) | **no** | — |

The empty-block cases confirm the loader keys on key **presence**, making the
authored layer deliberately stricter than the kernel's length-based predicate.
(REQ-11, REQ-12, REQ-13) The control load produced 2 normalized escape rows,
both with `len(Writes) == 0`, and its `KernelTable().CheckValid()` returned
nil. Every loadable file in `internal/table/testdata` was swept: all conform.
(REQ-14)

End-to-end through the built binary:
`intrastate lint --as json --model <breaching>` emits
`{"code":"model-invalid", ..., "detail":"malformed_escape_declaration: escape
rule draft-no-match-escape carries a write block"}` at exit 2.

**REQ-98** was checked on the shipped fixture rather than on prose: an authored
`clear = ["prelock_lens"]` normalizes to a kernel write
`{Key: "prelock_lens", Value: "<clear>"}`, so the `<clear>`-sentinel-write
representation the kernel predicate keys on is the same one the normalizer
emits. The two enforcers read one representation.

**H — the C7 CLI surface** (REQ-58..REQ-69, REQ-109)

Driven through `kernelResolveFailure` and serialized with `clierr.EmitJSON`:

- Single breach → `Code: "escape-row-shape-breach"` (unprefixed, per Q4's
  reading of the contract literal), `Group: GroupInternal`, `ExitCodeFor` = 2,
  `Hint: "fix the table producer: an escape row must carry no writes"` matching
  REQ-109 verbatim. (REQ-58, REQ-64, REQ-67, REQ-109)
- The wire envelope carries no `cause` key — `Cause` is `json:"-"` — and the
  identities reach the wire on `findings`, not `detail`. `Detail` is empty on
  the breach envelope. (REQ-59, REQ-60)
- `Finding` fields used are `Rule` (rule id) and `Locator` (source locator),
  both plain strings, all clierr-local. `go list -deps ./internal/cli/clierr`
  reports **zero** dependency on `internal/resolve`: clierr stays a leaf and
  the conversion happens in the verb layer. (REQ-61, REQ-62)
- `Count` serializes on the `Finding` record (Q1 reading (a)): a five-row,
  three-identity table produced three findings and the wire carried
  `"count":3` on the collapsed degenerate entry and `"count":1` on the others.
  The count is never recovered from prose. (REQ-63)
- Discrimination holds in both directions: a generic kernel error and RDR
  0008's live reserved-key breach (produced by a real `Resolve` call) both fall
  through to `flow-accessor-failed`, neither mislabelled as a shape breach; and
  the shape breach never falls through. (REQ-58)
- `CLIError.Cause`'s doc comment at `internal/cli/clierr/clierr.go` no longer
  claims "the wire-visible cause surface is Detail" — it now names both
  `detail` and `findings`, and `docs/cli-output-contract.md` documents both
  `findings` and its `count` column. (REQ-65, REQ-66)
- No refusal kind was added anywhere: `RefusalKinds()` is unchanged and the
  breach carries no `RefusalKind`. (REQ-68, REQ-69)

**I/J — scenarios, fixtures, mutation** (REQ-77..REQ-82, REQ-93..REQ-96)

- **REQ-77/78/79 oracle, regenerated rather than diffed against the spike
  artifact.** The base commit `05607f5` (immediately before Phase 1's
  `f476fbd`) was extracted with `git archive` into a scratch tree and its
  `internal/resolve` suite run; the same suite was run at HEAD with this RDR's
  own `*_0009_test.go` files removed, restricting the comparison to the tests
  existing at the base commit. Both runs produced **493** outcomes, **0 FAIL**,
  and the sorted outcome sets are **identical in both directions** (`comm -23`
  and `comm -13` both empty). Zero disposition change for conforming tables.
- **REQ-81 / Mutant A.** `Table.CheckValid` forced to `return nil`
  unconditionally, this RDR's own tests excluded: the frozen suite still
  **passes**, as the record expects.
- **REQ-82 / Mutant B.** `CheckValid` restored and `Writes` re-introduced on
  `fixtures_test.go::escapeRow`, this RDR's own tests excluded: the frozen
  suite **fails**, e.g.
  `TestReq15_EscapeEdgeRescuesOnlyWhenItMatchesExactlyOnce` and
  `TestReq16_ZeroAndMultipleMatchesRefuseUnlessEscapeIsModeled` report
  "Resolve returned a Go error for a modeled disposition". The check is wired
  into the live evaluation path and is not dead code. Both mutants were
  reverted; `git status --porcelain` is clean.
- **REQ-93/94.** The `escapeRow` builder carries no `Writes` and retains
  `NextTags`; neither `adversarial_test.go` nor `fixup_test.go` re-adds a
  `Writes` override on an escape row. **REQ-96**: `internal/resolve/resolve.go`
  imports both `errors` and `fmt`, `fmt` used inside `Error()` for REQ-43's
  diagnostic message (Q2's reading confirmed).

**G — doc-contract amendments** (REQ-54..REQ-57)

- `Resolve`'s numbered evaluation-order list opens at a new step `0` naming the
  whole-table precondition, above the alphabet check. (REQ-54)
- `Result`'s and `Resolve`'s cardinality contracts are both scoped to nil-error
  returns ("WHEN RESOLVE RETURNS A NIL ERROR" / "WHEN IT RETURNS A NIL
  ERROR"). (REQ-55)
- `Row`'s doc carries the PRODUCER OBLIGATION paragraph. (REQ-56)
- `CheckValid`'s doc opens "returns nil if t is valid, or else an error
  describing a problem", names the one property, and states "NOTHING ELSE is
  checked", explicitly excluding the Escape-class restriction. (REQ-50,
  REQ-51, REQ-57)

**L — performance and non-goals** (REQ-99..REQ-108)

- `testing.AllocsPerRun(100, ...)` over a 200-row conforming table of escape
  rows: **0 allocations per run**. Only a breach reaches the recording pass.
  (REQ-99, REQ-100)
- No memoization: mutating `Writes` in place on the rows backing a `Table`
  value flips the verdict on the very next `CheckValid()` call, and flips back
  when restored. The cost is genuinely per-call. (REQ-101)
- **REQ-103 brute-forced.** Twenty input shapes — five write-sets (nil, empty,
  one write, a `<clear>` write, two writes) × two escape classes × two
  ordinary-row configurations driving `no_match` and `ambiguous_match` — were
  run through `Resolve`. Four emitted an escaped plan and twelve errored on the
  precondition; **not one escaped plan carried a non-empty `Writes`**. The
  probe was checked non-vacuous (escaped plans were in fact emitted).
- REQ-102 (no byte-stable hash / canonical serialization introduced),
  REQ-104/REQ-105 (no refusal kind added, taxonomy and mapping unchanged),
  REQ-107 (fixture sites conformed), REQ-108 (`RowRef` + `compareRefs` reused,
  same sorted order) all hold under the above.

**REQ-70 (REQ-MVV)** is discharged by the combination above: the hand-built
write-bearing escape table errors naming `RuleID` and `SourceLocator` with no
`Plan` and no `Refusal`, the same table with writes removed resolves and
escapes unchanged, and the exported validator and the entry check were
exercised as the same predicate on the same tables (REQ-83 agreement run).

### Deferred, not probed

REQ-90 (TS scenario 8) and REQ-91 (TS scenario 11) are DEFERRED by the
record's own scope note and were skipped per this phase's brief. Their
substance is nonetheless covered above: REQ-90's authored-path assertions ran
anyway against `internal/table` (the "strengthening, not a scope breach" the
req-list's ASSUMPTION permits), and REQ-91's CLI assertions are carried by
REQ-58..REQ-69 and REQ-109.

### Observations recorded, not raised as failures

Neither item below is a violation of any REQ in the list; both are recorded so
a later reader does not have to re-derive the judgement.

1. **The C7 breach path is unreachable through authored TOML.** Because the
   RDR 0002 normalizer rejects a write-bearing escape rule at load
   (`malformed_escape_declaration`), no `flow resolve` invocation over an
   authored model can reach `kernelResolveFailure`'s breach branch. This is
   the intended layering — C2 is the authored-path enforcer and C3 is the
   producer/programmer backstop for non-TOML producers (REQ-47's explicit
   rationale) — so the C7 surface is correctly reachable only by a non-TOML
   table producer. It was therefore verified at the function boundary rather
   than through the binary. Not a defect; noted because the end-to-end
   `flow`-verb assertion REQ-109 describes can only be written against a
   hand-built table, never against a fixture file.

2. **`escapeShapeBreaches` degrades under an `%w` wrap.** Probed directly:
   `kernelResolveFailure(fmt.Errorf("ctx: %w", aggregate))` over a three-breach
   aggregate produces one `Finding` instead of three, because the traversal
   reaches the aggregate with a bare `err.(interface{ Unwrap() []error })`
   assertion rather than walking the chain. **This violates no REQ on any
   reachable path**: REQ-36 requires `Resolve` to return `CheckValid`'s error
   VERBATIM and never `fmt.Errorf`-wrapped, and the sole live call site
   (`flow_resolve.go:115`) passes `Resolve`'s return through untouched — which
   was confirmed by probe (bare aggregate → 3 findings). The `Code`, `Group`,
   exit code, and `Hint` all remain correct even under a wrap; only the
   `Findings` cardinality degrades. Recorded as latent brittleness against a
   future intermediate wrapper, not as a FAIL-N. The two other
   `resolve.Resolve` call sites (`flow_next.go:281`, `flow_resolve.go:217`) do
   not surface breaches at all — both strip or filter `Escape` before the call
   and are advisory probes — so neither is a missing-surfacing gap.

### Ambiguities resolved during this pass

- **REQ-62's "`Rule`/`Param` for the rule id"** — the implementation populates
  `Finding.Rule` and leaves `Finding.Param` empty. Read as disjunctive: the
  req-list phrases it as a suggested mapping over two candidate fields, and
  `Rule` is the semantically correct one (`Param` names a CLI parameter
  elsewhere in the envelope). Not treated as a violation.
- **REQ-22's placement** — the RDR calls it "a cheapness preference, not an
  observable contract", and the req-list ASSUMPTION forbids asserting that
  `assemble` was not called. Probed only for its observable consequences
  (REQ-19 whole-table scope, REQ-21 precedence), both of which hold. Noted
  that `CheckValid` runs above RDR 0008's `CheckInput`; a table breaching both
  reports the 0009 breach. JDR 0001 §JD-5 leaves that relative order open, so
  no REQ is engaged either way.
- **REQ-97/98's "named, shared fixture set"** — the fixtures ship inside
  `internal/table/escape_shape_0009_test.go`, which this phase is forbidden to
  read, so the artifact's *form* could not be inspected. Its *substance* was
  verified independently against the shipped TOML corpus (the authored-clear
  case above, plus the testdata sweep). The req-list ASSUMPTION admits "a
  named fixture file plus its expectations" without fixing the file kind, so
  a Go-resident fixture set satisfies the reading on its face.

**Verdict: OK.** No FAIL-N entries.

## Phase 3b — adversarial review (independent)

Written against the RDR's `## Trade-offs > ### Failure Modes` section
only. The Phase 1 test files and Phase 3a's findings were not read.

Added files:

- `internal/cli/escape_shape_adv_0009_test.go`
- `internal/cli/clierr/finding_count_adv_0009_test.go`
- `internal/resolve/escape_shape_adv_0009_test.go`

No existing test was edited or weakened. `golangci-lint run` reports
0 issues; `go vet ./...` is clean; the only failing tests in the tree
are the three adversarial ones recorded below.

### ADV-1 — the per-identity `Count` is dropped by the shared text renderer, so a collapsed breach loses its multiplicity in the CLI's default output mode

**Anchor.** Failure Modes, *Visible break (programmatic producer)* and
*Diagnosis (revised at Resolve)*: the breach must be diagnosable without
reading error prose, and the runbook entry is "branch on the code; the
named rows identify the broken producer." The Normative Contracts make
the count **structural** precisely because collapsing equal identities
discards multiplicity — "never only in formatted prose, which would be
readable solely by the string matching this RDR forbids everywhere else"
— and for the degenerate identity require that "the error MUST remain
diagnostic in that case by stating the breach count alongside the
identities."

**The defect.** `respond.Fail` routes `ModeJSON` to `clierr.EmitJSON` and
**everything else — the default — to `clierr.EmitText`**. `EmitText`
renders each finding through `Finding.identitySuffix()`, whose field list
was **not** widened when `Finding.Count` was added. So `count` reaches the
JSON envelope and is absent from text.

Observed for three breaching rows sharing `RowRef{"dup-rule",
"flow.toml:10"}`:

- JSON: `…"rule":"dup-rule","count":3…`
- text: `escape-row-shape-breach: this escape row carries writes (locator="flow.toml:10" … rule="dup-rule")`

The operator in the default mode cannot tell three rows breached rather
than one. The degenerate case is worse: for three rows built with no
source identity, `identitySuffix()` renders no identity fields at all and
no count, so the text finding line is
`escape-row-shape-breach: this escape row carries writes (hint="…")` —
literally nothing actionable — while the kernel's own `Error()` string on
the same value reads `escape row "" at "" carries writes (3 breaching
rows)`. The one channel that carries the number is the prose the RDR
forbids consumers from parsing.

**Tests.**
`internal/cli/escape_shape_adv_0009_test.go`:
`TestAdv0009_TheCollapsedBreachCountSurvivesTextRendering`,
`TestAdv0009_ADegenerateIdentityBreachIsStillDiagnosticInText`.
Both assert only *presence* of the count in text output — form and layout
are left free.

**Status: FAIL** (both). The remedy is one line: add `Count` to the
`identitySuffix()` field list (it is the record's only non-string field,
so it needs its own rendering arm rather than the `[][2]string` loop).

### ADV-2 — `escapeShapeBreaches` reaches the aggregate with a bare type assertion, so one `%w` wrap silently reduces an N-breach report to its first breach

**Anchor.** Failure Modes, *Visible break (programmatic producer)*: "A
single `errors.As`/`AsType` call reports only the FIRST breach in the
chain, so it is the wrong instrument for reading a multi-breach report —
a test asserting on all offending rows must traverse the aggregate."

**The defect.** `internal/cli/flow_resolve.go:296` reaches the join with

```go
if joined, ok := err.(interface{ Unwrap() []error }); ok {
```

a **bare type assertion**, which matches only the *outermost* error. Its
own comment claims the verb "traverses it rather than calling `errors.As`
once, which would report the first breach alone" — but under any wrap the
assertion fails and control falls through to exactly that single
`errors.As`. Observed on a three-identity breach: unwrapped → 3 findings;
wrapped once with `fmt.Errorf("…: %w", err)` → **1** finding.

The loss is silent. `errors.Is` still classifies through the wrap, so
`Code` stays `escape-row-shape-breach`, `Group` stays `GroupInternal`, and
the exit code stays 2 — every assertion the RDR nominates as the
conformance signal still passes while two offending rows vanish from the
envelope.

Today `Resolve` returns the error verbatim, so the wrap is not reachable
from the shipped call site. That is a property of the *caller*, not a
defense in the verb, and it is exactly the coupling the RDR's own
"returned VERBATIM, never wrapped" comment admits to depending on. The
robust instrument is available and demonstrably works through wraps —
`errors.As(ce, &joined)` finds the join even through `*clierr.CLIError`,
whose `Unwrap` returns a single `error`. Replacing the assertion with
`errors.As` costs nothing and removes the coupling.

**Test.** `internal/cli/escape_shape_adv_0009_test.go`:
`TestAdv0009_MultiBreachTraversalSurvivesAWrappedKernelError`.

**Status: FAIL.**

### ADV-3 — the shared `clierr.Finding.Count` widening opens a text-mode hole every co-owning RDR inherits

**Anchor.** Failure Modes, *Error-channel creep guarded against*, plus
deviations D3 (`Count` is "new public surface on a record RDRs 0005/0006/
0008 co-own"). `0006:C14` fixes text mode as the completeness-preserving
renderer: "a renderer that drops one has dropped a defect the author needs
to see."

The widening carries two obligations that pull against each other:

1. **invisible to co-owners that do not set it** — no shipped envelope,
   golden, or round-trip may change; and
2. **carried by every shared renderer** — a producer that *does* set it
   must not be silently truncated on one output path.

**Obligation 1 holds.** `Count` is `omitempty` with zero value `0`, no
0005/0006/0008 producer sets it, and `graphlint.identityKey` (the
finding-identity tuple) does not include it — so dedup, ordering and
serialization are all untouched. Verified byte-exactly against
representative 0005, 0006 and 0008 findings and against a pre-widening
payload round-trip.

**Obligation 2 fails**, and the failure is on *shared* code:
`identitySuffix()` is 0005/0006 machinery, so the hole is not 0009-local —
any co-owner that adopts `Count` later inherits it silently. This is the
same root cause as ADV-1, recorded separately because the surface, the
owner, and the blast radius differ: ADV-1 is 0009 losing its own
diagnosis, ADV-3 is a shared record gaining a field one of its two
renderers cannot express.

**Tests.** `internal/cli/clierr/finding_count_adv_0009_test.go`:
`TestAdv0009_SharedTextRendererCarriesTheWidenedCountField` (**FAIL**),
`TestAdv0009_TheCountWideningIsInvisibleToProducersThatDoNotSetIt`
(**PASS**), `TestAdv0009_APreWideningPayloadRoundTripsUnchanged`
(**PASS**).

**Status: FAIL** (obligation 2); obligation 1 pinned as passing guards.

### ADV-4 — properties attacked and found sound (kept as labelled regression guards)

Each of these was written to break the implementation and did not.

- **Collapse, ordering, and determinism.** The collapsed report is a
  function of the input tuple: four permutations of a four-row table with
  a duplicated identity produce byte-identical JSON envelopes. Equal
  `RowRef`s collapse to one entry carrying the pre-collapse row count;
  `RuleID` and `SourceLocator` are both part of the identity, so
  `("a","f:1")` and `("a","f:3")` stay distinct.
  `TestAdv0009_TheCLIBreachEnvelopeIsAFunctionOfTheTableNotRowOrder` —
  **PASS**.
- **Aggregate shape and `errors.Is`/`As` fidelity.** The join is exactly
  one level deep, no element is itself an aggregate, every element
  classifies on its own via `Unwrap() → ErrEscapeShapeBreach` (not only
  the aggregate), and a single `errors.As` over the aggregate returns the
  first element — the RDR's stated hazard, pinned so a future change
  cannot quietly make the documented traversal wrong.
  `TestAdv0009_TheAggregateIsFlatAndEveryElementClassifies` — **PASS**.
- **Zero disposition change for conforming tables.** A conforming table
  reaches the same plan on the ordinary arm *and* on the escape-rescue arm
  (`Escaped == true`), and an empty/nil-`Rows` table conforms vacuously
  and still refuses `no_match` rather than erroring.
  `TestAdv0009_AConformingTableKeepsItsDisposition` — **PASS**.
- **Breach precedence and the zero `Result`.** A breaching table with a
  recognized outcome outside the declared alphabet errors rather than
  refusing `unmodeled_outcome`, and returns neither `Plan` nor `Refusal`.
  `TestAdv0009_ABreachOutranksEveryModeledDispositionAndCarriesZeroResult`
  — **PASS**.
- **The discrimination in `kernelResolveFailure`.** RDR 0008's
  reserved-key breach still maps to `flow-accessor-failed`, not to
  `escape-row-shape-breach`; exit stays 2. The change that stopped
  collapsing kernel errors into the generic accessor code did not
  over-recode.
  `TestAdv0009_ANonBreachKernelErrorKeepsTheGenericInternalCode` —
  **PASS**.
- **`Count` round-trips the wire.** Decoding the emitted envelope and
  re-reading `Count` per identity recovers 2 and 1 without touching prose.
  `TestAdv0009_ThePerIdentityCountRoundTripsThroughJSON` — **PASS**.
- **Precondition ordering vs. RDR 0008.** A table breaching both
  preconditions deterministically reports the 0009 breach (0009 runs
  first); a *conforming* table with an 0008 reserved-key breach still
  reports 0008. The documented "order is not an observable contract"
  claim holds at HEAD. Probed; not pinned, because pinning it would
  freeze an ordering the RDR deliberately leaves open.
- **The authored path.** `table.Load` already rejects
  `neg-escape-with-write.toml`, `neg-escape-with-clear.toml` and
  `neg-escape-with-empty-write.toml` under
  `malformed_escape_declaration`, including the presence-based rejection
  of an *empty* write block that the RDR flags as deliberately stricter
  than the kernel's length-based predicate. Verified by direct `Load`
  probe; no new test added, since the existing `rules_test.go` /
  `roundtrip_test.go` coverage already pins it.

### ADV-5 — recorded observations (no test added)

- **`neg-escape-with-empty-write.toml` is absent from the
  `roundtrip_test.go` load-category map** (`internal/table/roundtrip_test.go:352-354`
  lists the `-with-clear`, `-with-gate` and `-with-write` fixtures but not
  `-with-empty-write`), even though `Load` does reject it with
  `malformed_escape_declaration`. It is covered by `rules_test.go` instead,
  so this is a coverage-map gap rather than a behavior gap. Recorded, not
  edited — the map is an existing test.
- **REQ-109 asks for an end-to-end `flow` test** ("in the style of the
  existing `flow_*_0005_test.go` suites") asserting the `Code`, the
  `Findings` identities, and the `Hint`. Such a test appears
  unsatisfiable at HEAD: the only route from a loaded model to
  `resolve.Resolve` is `Model.KernelTable()`, and `table.Load` rejects
  every escape rule carrying a write block, so no *loadable* model can
  produce a breaching kernel table. The obligation is therefore reachable
  only at the `kernelResolveFailure` unit boundary, which is where the
  present tests exercise it. Recorded as an ambiguity resolved in favor of
  the unit boundary; if REQ-109 is meant literally it needs a
  construction seam that does not exist.
- **`CLIError.Detail` is left empty** on the breach envelope; the human
  sentence rides `Message` instead. Normative text at RDR line 1862 says
  "`Detail` carries the human sentence, the new field carries the
  identities," but the binding clause describes the `config.Load` pattern
  without mandating a non-empty `Detail`, and `Message` serves the same
  role in the rendered output. Defensible as-is; recorded because a
  literal reading of line 1862 would call it a miss.

### Verdict

**BLOCK.** Three added tests fail against the implementation
(ADV-1 ×2, ADV-2, ADV-3 ×1 — four failing test functions across three
failure modes). All three share the CLI/`clierr` seam; the resolution
kernel's half of RDR 0009 withstood every attack made on it. The
remedies are small and local: widen `identitySuffix()` to render `Count`
(closes ADV-1 and ADV-3), and swap the bare type assertion in
`escapeShapeBreaches` for `errors.As` (closes ADV-2).
