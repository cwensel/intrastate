model: claude-opus-5[1m]
variant: full (profile: foundational)

## Reconstruction basis

Read via projector only. Elements read by id: `C1`, `C2`, `MVV`, `S1`–`S5`.

**Widened past the contract spans** (recorded per the lens instruction — each
widening was needed because the contract left a signature, type, site, or step
order under-determined):

- `§technical-design` (lines 485–909) — C1/C2 name no Go symbols for the
  projection function itself. Widening supplied `§source-authority-census`
  (the 14-field table with writers), `§disposition` (exit/envelope table), and
  `§oracle-discriminability` (the negative controls). Without this the entire
  data model in item 3 would have been a GUESS.
- `A2` — C1's `D-wire-byte-format` explicitly *defers* the in-code mechanism.
  A2 fixes it: pointer-valued echo fields with `omitempty`, populated in
  default mode, nilled at projection. Reading only C1 would have manufactured
  a "projected struct" GUESS the record actually settles.
- `A9` — the nil-container invariant that A2's mechanism rests on; it names
  the four container writers (`tagMap`, `observedTagMap`, `emitMap`,
  `readerIDs`) and their `make(...)` opens.
- `A5`, `A6` — the tree walker's shape and the refusal-path flag-blindness.
- `§approach`, `§existing-infrastructure-audit`, `§technical-environment` —
  supplied the concrete registration site (`newFlowResolveCmd`), the help
  symbol (`flowResolveExtendedDesc`), and the confirmation that
  `registerSelectionFlags` is the shared registrar the flag must never enter.
- `§phase-1`/`§phase-2`/`§phase-3`, `§prerequisites`, `§desk-trace` — fixed
  step order for item 4 (golden capture *before* the struct edit; projection
  site after the last `respond.Fail`) and the pinned `NumField() == 14`.

Spans deliberately *not* widened: `§alternatives-considered`,
`§briefly-rejected`, `§research-findings`, `§finalization-gate` — none bear on
the API surface.

---

## 1. Public API

The module is `internal/cli` in Go module `github.com/cwensel/intrastate`.
The RDR contracts a CLI surface, not an exported Go package; the "public API"
is therefore (a) the command-line surface, which C1 fixes normatively, and
(b) the in-package Go surface the projection introduces, which C1 leaves to
the mechanism A2 verified.

### 1a. Command-line surface (normative, C1)

```
intrastate flow resolve [existing flags] [--plan-only]
```

| Property | Value | Source |
| --- | --- | --- |
| Name | `--plan-only` | C1, `D-naming` |
| Type | boolean | C1 |
| Default | `false` | C1 |
| Shorthand | none (long form only) | C1 |
| Registration site | `newFlowResolveCmd` (verb-local) only | C1, infra audit |
| Forbidden sites | `registerSelectionFlags`, the `flow` group's own or persistent flag set, the root's persistent set | C1 |
| Accepted by | exactly `{flow resolve}` over the whole command tree | C1, S4 |

Error modes on the CLI surface — the flag mints **no new refusal code and no
new exit group** (C1 states this twice):

| Input | Exit | Envelope |
| --- | --- | --- |
| `--plan-only` on `flow resolve`, success | 0 | projected success payload |
| `--plan-only` on `flow resolve`, refusal | unchanged from default mode | refusal, byte-identical (A6) |
| `--plan-only` on `next` / `read-state` / `set-state` | 2 | `command-error` shared bucket (cobra parse failure) |
| `--plan-only` on a pre-RDR binary | 2 | `command-error` |

### 1b. In-package Go surface

GUESS (signatures): C1 and C2 fix the *behaviour* and the *carrier
constraint* but name no function. Grounded on A2's mechanism (pointer echo
fields, nilled at projection) and C2's "standalone table keyed by field name,
in its own declaration site":

```go
// internal/cli/flow_resolve.go

// resolvePayload is the shipped success payload. Under this RDR the five
// ECHO fields change TYPE only (concrete -> pointer); field COUNT stays 14
// (pinned: decision_table_0010_test.go NumField()==14, :435).
type resolvePayload struct { /* see §3 */ }

// projectPlanOnly returns p with the ECHO group nilled, so encoding/json
// omits those keys under `omitempty`. Plan-group fields are untouched;
// carried values are the same values, so their bytes are identical.
// Implemented FROM the declaration table (C2), never as a bare omit-list.
func projectPlanOnly(p resolvePayload) resolvePayload    // GUESS: name, value receiver, non-pointer return

// internal/cli/flow_partition.go  (its own declaration site — C2 forbids a
// marker on resolvePayload)
type payloadSide string                                   // GUESS: type name

const (
    sideEcho payloadSide = "echo"                         // GUESS: identifier names
    sidePlan payloadSide = "plan"
)

// resolvePayloadSides maps every resolvePayload JSON field name to exactly
// one side. One entry per field, 14 entries. C2: this is the assignment
// SOURCE; the projection reads it and the S3 oracle asserts against it.
var resolvePayloadSides = map[string]payloadSide{...}     // GUESS: name and map type
```

Error modes of the Go surface: **none**. `projectPlanOnly` returns no
`error`. GUESS, but a well-grounded one — the projection is total key
deletion over a fixed partition with no I/O, no parsing, and no
caller-supplied field list (C2: "a fixed partition, not a caller-supplied
field list"), so there is nothing to fail on. Consistent with C1's "mints no
new refusal code".

GUESS: the flag value is read once, into a field on the existing resolve
request struct (`req`, per the census's `req.modelRef` / `req.observed` /
`req.runGates` / `req.revision()` accessors) — call it `req.planOnly` — and
read at exactly one site. Phase 1 makes this normative in spirit: "The flag is
read once, at the projection site, and nowhere in gate evaluation, rule
selection, or refusal construction."

---

## 2. Three most important internal helpers

**(1) `projectPlanOnly(resolvePayload) resolvePayload`** — GUESS on the name;
the responsibility is contracted. Applies the ECHO/PLAN partition by nilling
the five pointer-valued echo fields (`model`, `observed`, `owned`, `readers`,
`outcome`), leaving all nine plan fields (`revision`, `rule`, `gates`,
`emit`, `next`, `writes`, `clear`, `escaped`, `escape_class`) untouched.
Three constraints it must satisfy, all from C1/C2:

- It is a *key deletion*, never a re-encoding. Surviving values keep their
  default-mode bytes exactly; surviving keys keep default-mode relative
  (declaration) order.
- It is driven **from** `resolvePayloadSides`, not from an inline literal
  list. C2: "the projection MUST NOT be implemented as a bare omit-list whose
  complement is 'whatever else exists'."
- It is called from the SUCCESS path only, **after the last `respond.Fail`
  return** and **before `respond.OK`** — never inside the respond gateway and
  never per output mode. Phase 1 forbids a defensive `if refusing, skip
  projection` guard: refusal flag-blindness must be structural (A6), and a
  guard would prove the site is misplaced.

**(2) The declared assignment table (`resolvePayloadSides`)** — a helper by
role even though it is data. C2 constrains its carrier *hard*: it MUST be a
standalone table keyed by field name in its own declaration site, and MUST
NOT be a per-field struct tag or marker on `resolvePayload`. The reasoning is
load-bearing and worth restating: A2's mechanism rewrites the echo fields'
types and tags on that struct, so a marker living there would be edited in
the very commit that changes the projection — the implementer moves the
assertion and its subject together, and S3 degenerates into asserting the
implementation agrees with itself. Independence has to be *structural*.

**(3) The total command-tree walker (test-side)** — GUESS on the name; call
it `walkCommandTreeTotal(root *cobra.Command) []*cobra.Command`. Descends
every child of the root with **no name-based skip and no `Hidden` gate**,
including the auto-generated `help` and `completion` commands. It MUST NOT
reuse `internal/cli/help_all.go::walkCommandTree`, which skips children named
`help` or `completion` before recursing (and whose `docs.go` callers further
gate on `Hidden`) — that walker asserts a strictly weaker negative. Its
caller must first materialize both auto-generated commands via cobra's own
initializers (`InitDefaultHelpCmd`, `InitDefaultCompletionCmd` reachable from
`ExecuteC`) — never a hand-constructed stand-in — and assert both are IN the
walked set as a precondition. `completion` is the discriminating member: a
bare `NewRootCmd()` tree has children `docs`, `flow`, `help`, `lint`,
`version` and no `completion` at all, so a `help`-only control cannot tell a
total walker from an untested tree. This is A5, still Pending: its
verification IS Phase 2 build work.

*Honourable mention, not in the top three:* the four container writers
`tagMap`, `observedTagMap` (`flow_exec.go`), `emitMap` (`flow_resolve.go`)
and `readerIDs` (`flow_next.go`). They are not new, but A9 makes their
unconditional `make(...)` load-bearing under the pointer conversion: a nil
container reached by any of them renders `"readers":null` / `"owned":null`
under `omitempty`, silently changing DEFAULT-mode bytes. A9 is Pending and
its Phase 1 oracle (empty-container request must render `{}`/`[]`, never
`null`) is a precondition on the conversion landing.

---

## 3. Data model across the boundary

The boundary is the NDJSON wire record written to stdout. Nothing is
persisted to disk by this RDR.

### Envelope

```
{"type":"ok","data":{ …payload… }}\n
```

The **full emitted line**, envelope included and trailing newline excluded, is
the normative measurand for C1's STRICTLY SHORTER clause — not the `.data`
payload alone. Asserted in both output modes.

### `resolvePayload` — 14 fields, one assembly site

All fourteen are written at the single `resolvePayload{…}` literal in
`internal/cli/flow_resolve.go`. No fallback arm, no second writer — which is
what makes the partition a single-site edit.

| # | JSON key | Writer | Go type after conversion | C2 side |
| --- | --- | --- | --- | --- |
| 1 | `model` | `req.modelRef` | `*string` (GUESS: pointer) | ECHO |
| 2 | `observed` | `observedTagMap(req.observed)` | `*map[string]string` (GUESS) | ECHO |
| 3 | `owned` | `tagMap(owned)` | `*map[string]string` (GUESS) | ECHO |
| 4 | `readers` | `readerIDs(readers)` | `*[]string` (GUESS) | ECHO |
| 5 | `outcome` | `outcome` | `*string` (GUESS) | ECHO |
| 6 | `revision` | `req.revision()` | `string` | PLAN (identity slot) |
| 7 | `rule` | `plan.RuleID` | `string` | PLAN |
| 8 | `gates` | `req.runGates` | `[]…` (GUESS: gate-result slice) | PLAN |
| 9 | `emit` | `emitMap(row.Emit)` | `map[string]string` (GUESS) | PLAN |
| 10 | `next` | `tagMap(plan.NextTags)` | `map[string]string` (GUESS) | PLAN |
| 11 | `writes` | `plan.Writes` split on `ClearSentinel` | `map[string]string` (GUESS) | PLAN |
| 12 | `clear` | `plan.Writes` split on `ClearSentinel` | `[]string` (GUESS) | PLAN |
| 13 | `escaped` | `plan.Escaped` | `bool` | PLAN |
| 14 | `escape_class` | `escapeClassOf(...)` | `string,omitempty` (GUESS: presence via `omitempty`) | PLAN |

The record itself notes the census is thirteen ROWS for fourteen fields
because `writes` and `clear` share the one `plan.Writes` writer.

Field count is pinned at 14 by `decision_table_0010_test.go` (`:435`), with
the 13-key wire list at `:422-426`. A2's conversion changes field **types
only** — adding a field is out of scope for this RDR by construction, since
it would falsify A3.

### Two widths

- **Default** — all 14 fields. Byte-identical to today's output. Empty
  containers render `{}` / `[]` (never `null`) — the A9 invariant.
- **Projected (`--plan-only`)** — exactly `revision`, `rule`, `gates`,
  `emit`, `next`, `writes`, `clear`, `escaped`, plus `escape_class` **iff
  the same request's default output carried it** (the producer's own presence
  rule, `0005:A-3` / `escapeClassOf` — not an unconditional "when escaped").
  The five echo keys are **ABSENT** — never `null`, never `{}`, never `""`.

Worked projected record (illustrative in the RDR; the normative bytes live
in `evidence/spikes/a2-encoder-mechanism.md` §Reference output):

```json
{"revision":"","rule":"free-eu","gates":[],
 "emit":{"dpa":"required","plan":"basic"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

### Presence and ordering rules

- `emit` stays present as `{}` even when empty (`0010:C4`).
- `escape_class` follows its producer's rule unchanged; oracles compare it
  against the SAME run's default output, never an unconditional literal.
- **Wire** keys keep struct declaration order. **Text** is alphabetically
  sorted by the shipped `flatten` (`respond/text.go`, `sort.Strings`). These
  are two different orderings and no clause requires them to agree; the text
  subset property is set membership, not subsequence.

### Always-keep core

`rule`, `escaped`, `escape_class`, `revision` — projection-invariant, not
unconditionally present. The laundering guard actually rests on `escaped`
(unconditionally present); `escape_class` keeps its presence rule.
`revision` is contracted but **presently vacant**: `0002`'s `[model]` block
admits no revision key, so the accessor returns `""` on every payload this
CLI can emit. Its always-keep clause is forward-binding.

This RDR mints **no** separate always-keep oracle — with one projection mode,
always-keep and S2's key-set literal have the same extension. The obligation
falls on the successor RDR that adds a second mode, which owes an always-keep
oracle quantified over modes.

---

## 4. Top-level pseudo-code

Two operations are load-bearing. The runtime one first, then the ordering of
the change itself (Phase 1's golden-capture step is a real step-order
constraint the record fixes, so it belongs here).

```
# ---- runtime: runFlowResolve ----
runFlowResolve(cmd, args):
  1  req        = parseRequest(cmd, args)          # --model, --tag, --outcome, ...
  2  planOnly   = cmd.Flags().GetBool("plan-only") # read ONCE, here only
  3
  4  model, err = loadModel(req.modelRef)
  5  if err: return respond.Fail(...)              # flag-blind (A6)
  6
  7  owned, readers, err = runReaderPass(model, req.outcome)   # ALWAYS runs.
  8  if err: return respond.Fail(...)              # C1 forbids skipping work
  9                                                # whose only consumer is a
 10                                                # projected-away field.
 11  plan, row, err = resolve.Resolve(model, Input{Owned: owned, Observed: req.observed, ...})
 12  if err: return respond.Fail(...)              # flag-blind
 13
 14  gates = req.runGates(row)
 15  if gatesRefuse(gates): return respond.Fail(...)  # flag-blind
 16
 17  # ---- LAST respond.Fail is above this line. Success path only below. ----
 18
 19  writes, clear = splitOnClearSentinel(plan.Writes)
 20  payload = resolvePayload{                      # all 14 fields, one site
 21      model: &req.modelRef, observed: ptr(observedTagMap(req.observed)),
 22      owned: ptr(tagMap(owned)), readers: ptr(readerIDs(readers)),
 23      outcome: &outcome,                                       # ECHO x5
 24      revision: req.revision(), rule: plan.RuleID, gates: gates,
 25      emit: emitMap(row.Emit), next: tagMap(plan.NextTags),
 26      writes: writes, clear: clear, escaped: plan.Escaped,
 27      escapeClass: escapeClassOf(model, plan),                 # PLAN x9
 28  }
 29
 30  if planOnly:
 31      payload = projectPlanOnly(payload)         # nil the ECHO side, FROM
 32                                                 # resolvePayloadSides (C2).
 33                                                 # No guard for refusals —
 34                                                 # unreachable by placement.
 35  return respond.OK(cmd, payload)                # ONE gateway, both modes:
 36                                                 # json marshals directly,
 37                                                 # text goes through flatten.

# ---- change ordering: Phase 1 -> 2 -> 3 ----
 38  P1.0  capture the pre-change golden from the CURRENT binary and check it
 39        in, on a stash-clean tree, BEFORE any struct edit. It is the only
 40        pre-change side in the whole battery; S1-S5 all compare the new
 41        build against itself and cannot see a default-mode regression that
 42        moves both sides together.
 43  P1.1  land A9's oracle: default-mode run with empty containers renders
 44        {} / [] and no null. Must be green BEFORE the conversion.
 45  P1.2  convert the five echo fields to pointer + omitempty; NumField()
 46        stays 14. Add resolvePayloadSides in its own file. Add
 47        projectPlanOnly. Register --plan-only on newFlowResolveCmd only.
 48  P1.3  regenerate docs/cli-reference.md and llms.txt IN THIS COMMIT —
 49        make docs-check fails at a commit that registers the flag with
 50        stale generated artifacts.
 51  P2    write S1-S5. S4's total walker is A5's verification (Pending):
 52        materialize help + completion via cobra's initializers, assert
 53        both are in the walked set, then assert registrants == {flow
 54        resolve}. If it cannot reach that scope, C1's structural negative
 55        needs its own form -- this gates LOCK, not Phase 1.
 56  P3    prose: cli-output-contract.md gains the projected worked payload;
 57        flowResolveExtendedDesc gains the flag under "Reading a
 58        successful plan"; C2's partition statement lands with the field docs.
```

---

## Guess register

Every GUESS above, collected. These mark where the record is silent, not
where it is wrong.

| # | Guess | Why it was needed |
| --- | --- | --- |
| G1 | `projectPlanOnly` as the projection function's name, its value receiver, and its non-pointer return | C1/C2 fix behaviour and carrier; no symbol is named |
| G2 | `resolvePayloadSides` / `payloadSide` / `sideEcho` / `sidePlan` as identifiers, and `map[string]payloadSide` as the table's type | C2 constrains the carrier's *shape and site* ("standalone table keyed by field name") but names nothing |
| G3 | `internal/cli/flow_partition.go` as the table's file | C2 says "its own declaration site", not which file |
| G4 | The projection returns no `error` | Follows from "mints no new refusal code" + fixed partition, but is not stated |
| G5 | `req.planOnly` as the flag's landing field | The census shows `req.*` accessors; the flag's storage is unstated |
| G6 | Concrete Go types for all 14 fields (`*string`, `*map[string]string`, `*[]string`, etc.) | The census gives writers and JSON keys, never Go types; A2 fixes only "pointer-valued echo fields with `omitempty`" |
| G7 | `escape_class` realises its presence rule via `,omitempty` | The rule is contracted (`0005:A-3`); the encoding mechanism is not |
| G8 | `gates` element type | `req.runGates` is the writer; the element type is never given |
| G9 | `walkCommandTreeTotal` as the test walker's name and its `[]*cobra.Command` return | C1 and A5 fix its semantics exhaustively; no name or signature |
| G10 | The exact runtime step order in lines 4–15 (load, reader pass, resolve, gates) | The record fixes only *where the projection sits* — success path, after the last `respond.Fail`, before `respond.OK` — and that the reader pass must not be skipped; the surrounding sequence is inferred from the census's writers |
| G11 | `ptr(...)` helper wrapping the container writers at the assembly site | A2 requires non-nil pointers to possibly-empty containers; the idiom is unstated |

## Where the record is notably NOT silent

Recorded because it is the counterweight to the guess register — this RDR
over-determines several things a typical record leaves open, and a diff
should not read these as reconstruction luck:

- The **measurand** of both soft assertions: STRICTLY SHORTER is the full
  emitted line including the `{"type":"ok","data":{…}}` envelope, and the
  invoked-reader set is OBSERVED READER EXECUTION at the invocation site —
  with the two inadmissible forms (read from the projected payload; recompute
  `invokedReaders(model, outcome)`) named and refused explicitly.
- The **carrier** of the assignment table, and why a struct marker is
  forbidden.
- The **totality** of the walk, down to `completion` being the discriminating
  member and cobra's own initializers being the required materialization.
- The **absence** of a defensive refusal guard at the projection site.
- The **negative control** for each of S1–S5.

Two clauses ride on Pending assumptions and would move if those refute:
A5 (the total walker's implementability — gates lock, not Phase 1) and A8
(text-mode strict width — if refuted, C1's width clause narrows to JSON
*before* lock). A9 is Pending with a Phase 1 oracle that gates the conversion.
