Model: claude-opus-5
Lens: cove (pass A)

RDR: `docs/rdr/0009-escape-row-shape-conformance-ownership.md`
Grounding base: branch `via-claude`, HEAD `20cc49f` (treated as `main`).
Toolchain for stdlib checks: `go env GOROOT` = `/opt/local/lib/go-1.26`.

## Step 0 — Grounding sweep

Every claim below was checked by reading the cited source, not by trusting the RDR.

### Kernel types and functions (`internal/resolve/resolve.go`)

| # | RDR claim | Verdict | Source |
| --- | --- | --- | --- |
| 1 | `Row` carries `RuleID`, `SourceLocator`, `Escape`, `Writes`, `NextTags` | CONFIRMED | `internal/resolve/resolve.go::Row` — all five fields present |
| 2 | `Row` permits the illegal combination (`Escape` + `Writes` both populated) | CONFIRMED | `resolve.go::Row` — single struct, no variants, no constructor guard |
| 3 | `Plan` mirrors `RuleID`/`SourceLocator`; doc separates `NextTags` from `Writes` | CONFIRMED | `resolve.go::Plan` — "NextTags is the next state." vs "Writes are the owned-tag writes for the accessor layer." |
| 4 | A1: `Resolve` — ALL return paths carry a literal `nil` error | CONFIRMED | `resolve.go::Resolve` — 5 `return …, nil` sites (unmodeled_outcome, blocked refusal, exact-one plan, both escapeOrRefuse arms) |
| 5 | Helpers `gate`, `escapeOrRefuse`, `planOf`, `refuse` return no error type | CONFIRMED | signatures in `resolve.go` |
| 6 | `planOf` copies `Writes` unconditionally, escape selections included | CONFIRMED | `resolve.go::planOf:507` `Writes: copyTags(row.Writes)` — no `escaped` branch |
| 7 | `rescues` exists as escape discriminator | CONFIRMED | `resolve.go::Row.rescues` |
| 8 | Candidate loop partitions on `len(row.Escape) != 0` | CONFIRMED | `resolve.go::Resolve` |
| 9 | `RowRef` exists with `RuleID`/`SourceLocator` | CONFIRMED | `resolve.go::RowRef` |
| 10 | `rowRefs` sorts by `compareRefs` for REQ-1/REQ-10 stability | CONFIRMED | `resolve.go::rowRefs` — `slices.SortFunc(out, compareRefs)` |
| 11 | `compareRefs` orders by `RuleID` then `SourceLocator` | CONFIRMED | `resolve.go::compareRefs` |
| 12 | `copyTags` is nil-preserving | CONFIRMED | `resolve.go::copyTags` — `if in == nil { return nil }` |
| 13 | `Resolve` performs no structural check on `Table.Rows` before evaluation | CONFIRMED | `resolve.go::Resolve` — `assemble` then alphabet check; no shape pass |
| 14 | `Resolve` doc: "the error return is reserved for programmer mistakes" | CONFIRMED | `resolve.go::Resolve` doc, verbatim |
| 15 | `RefusalKind` carries the no-string-inspection sentence | CONFIRMED | `resolve.go::RefusalKind` doc line 45 |
| 16 | Kernel already scans `Table.Rows` twice per call + `gate` | CONFIRMED | `Resolve` candidate loop, `escapeOrRefuse` loop, `gate` |

### Test-tree claims — the critical failures

| # | RDR claim | Verdict | Source |
| --- | --- | --- | --- |
| 17 | `adversarial_test.go:229` assigns `escape.Writes = []resolve.Tag{{Key:"status",Value:"Escaped"}}` | CONFIRMED | exact line, in `TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement` |
| 18 | `fixup_test.go:103` assigns the identical override | CONFIRMED | exact line, in the "missing owned state must not rescue an ambiguity" subtest |
| 19 | **A3 Scope correction: `fixtures_test.go::escapeRow` "sets no `Writes` at all"** | **REFUTED** | `internal/resolve/fixtures_test.go:257` — `Writes: []resolve.Tag{{Key:"status",Value:"Blocked"}},` inside the `escapeRow` literal (lines 249-260), one line below `NextTags:` at :256 |
| 20 | "The write-bearing escape rows are two call-site overrides" / "edit all three sites" | **REFUTED** | builder is itself write-bearing; `rg 'escapeRow\('` → 16 call sites. Every escape row in the suite breaches |
| 21 | Phase 2: "The `escapeRow` builder needs no change" | **REFUTED** | same cite as #19. The RDR's own spike (`evidence/spikes/a3-fixture-conformance.md`) lists three edits, first being removal of that exact `Writes` line, and warns "escapeRow is not the only place an escape row gets Writes" |
| 22 | Infra Audit row: "Builder sets `NextTags` only" | **REFUTED** | same cite as #19 |
| 23 | A3: every `Plan.Writes` read on an escape-derived plan sits in a `t.Fatalf` format arg, never a condition | CONFIRMED | `adversarial_test.go:241`, `fixup_test.go:57`, `fixup_test.go:115` |
| 24 | REQ-18/REQ-19/REQ-29/MVV-leg-1 `Writes`-equality assertions run on non-escape fixtures | CONFIRMED | `resolve_test.go:671-673`, `:709`, `:731`, `:980`, `:986`; `mvv_test.go:67-69` |
| 25 | "four length-based `Writes` checks" | CONFIRMED | exactly 4: `resolve_test.go:668`, `:709`, `:986`, `mvv_test.go:110` |
| 26 | "a frozen suite with **zero** nil-sensitive assertions" | **REFUTED as written** | ~26 nil-sensitive assertions on `Plan`/`Refusal` (`resolve_test.go:78-83`, `mvv_test.go:42,45,99,102,140`, …). Narrowed claim (zero nil-sensitive on `Writes`) does hold |
| 27 | `boundary_test.go` uses no escape fixtures | CONFIRMED | AST helpers only; `rg -c 'resolve\.'` → 0 |

### A5 — no production file imports `internal/resolve`

| # | RDR claim | Verdict | Source |
| --- | --- | --- | --- |
| 28 | Only importers are `internal/resolve/*_test.go` + the RDR 0007 spike | CONFIRMED | repo-wide grep → 5 kernel test files + `docs/rdr/0007-…/evidence/spikes/aggregation-probe_test.go`. Zero production importers |
| 29 | `cmd/intrastate/main.go` is `func main() { cli.Execute() }` | CONFIRMED | exact body |
| 30 | A5 fixture-builder list | CONFIRMED | all in `internal/resolve/fixtures_test.go` |
| 31 | A5 file list `{fixtures,resolve,adversarial,mvv,fixup}_test.go` | CONFIRMED | `boundary_test.go` correctly omitted (does not import the package) |

### CLI / clierr

| # | RDR claim | Verdict | Source |
| --- | --- | --- | --- |
| 32 | `case GroupUserEnv, GroupInternal: return 2` exists | CONFIRMED | `clierr.go::ExitCodeFor:119-120`, verbatim |
| 33 | Non-`CLIError` defaults to exit 1 | CONFIRMED | `clierr.go::ExitCodeFor` trailing `return 1` + doc |
| 34 | `CLIError.Unwrap` keeps the cause traversable | CONFIRMED | `clierr.go::CLIError.Unwrap` |
| 35 | `.Detail` "carries hard facts about the failure" | CONFIRMED | field comment, verbatim |
| 36 | `.Hint` and `.Code` exist as separate fields | CONFIRMED | `clierr.go::CLIError` |
| 37 | `ErrorCode`/`ExitCodeFor` branch on `Code`/`Group` only, never text | CONFIRMED | both functions |
| 38 | `config.go::Load` wraps a Go error into `GroupInternal` with `Code`/`Detail`/`Cause` | CONFIRMED (those four fields) | `config.go::Load:84-90` |
| 39 | **A2 + contract: `config.Load` is precedent for `GroupInternal` + stable code **+ `Hint`** | **REFUTED** | `config.go::Load:84-90` carries **no `Hint`**. Every `Hint` in the repo is on a `GroupUserEnv` error (`Discover:64`, `root.go:122`). No shipped `GroupInternal` error carries a `Hint` |
| 40 | `root.go::NewRootCmd` registers only `newVersionCmd()` | CONFIRMED | single `AddCommand`; no `flow` verb |
| 41 | `cobraErrorToCLIError` stamps `GroupUserEnv`/`command-error` | CONFIRMED | `root.go::cobraErrorToCLIError` |
| 42 | AGENTS.md "Add new codes as needed" | CONFIRMED | `AGENTS.md` — "…a `Group` that maps to an exit code. Add new codes as needed; keep envelope fields `omitempty`." |

### Go stdlib precedents (local GOROOT `/opt/local/lib/go-1.26`)

| # | RDR claim | Verdict | Source |
| --- | --- | --- | --- |
| 43 | `rsa.PrivateKey.Validate` doc form "returns nil if <x> is valid, or else an error describing a problem" | CONFIRMED | `$GOROOT/src/crypto/rsa/rsa.go:233` + doc, verbatim |
| 44 | `rsa.PrivateKey.Validate` has an idempotence short-circuit | CONFIRMED | `rsa.go` Validate — `Precomputed.fips`-consistency early return |
| 45 | pprof `Profile.CheckValid` exported + re-checked at parse boundary | CONFIRMED (existence/export); call-graph half UNVERIFIED | `$GOROOT/src/cmd/vendor/github.com/google/pprof/profile/profile.go:362` |
| 46 | `testing/fstest::TestFS` aggregates via `errors.Join` | CONFIRMED | `$GOROOT/src/testing/fstest/testfs.go:92` |
| 47 | `testing/slogtest::TestHandler` aggregates via `errors.Join` | CONFIRMED | `$GOROOT/src/testing/slogtest/slogtest.go:292`, doc at :250 |
| 48 | `encoding/json/v2::SemanticError` sentinel-in-`Err` + exported location field | CONFIRMED (shape); `Err`-field/`AsType` half PARTIAL | `$GOROOT/src/encoding/json/v2/errors.go:72` — exported `ByteOffset`, `JSONPointer` |
| 49 | `encoding/json::UnmarshalTypeError.Field` | CONFIRMED | `$GOROOT/src/encoding/json/decode.go:132`; v2 variant `v2_decode.go:116` also has `Err error` |
| 50 | `go/types::Error` structured | CONFIRMED | `$GOROOT/src/go/types/api.go:51` |
| 51 | `analysis.Validate` fail-fast for startup wiring | CONFIRMED (existence); fail-fast semantics PARTIAL | `$GOROOT/src/cmd/vendor/golang.org/x/tools/go/analysis/validate.go` |

### Peer-RDR quotes

| # | RDR claim | Verdict | Source |
| --- | --- | --- | --- |
| 52 | RDR 0004: "A write accessor MUST apply only planned owned-tag writes produced by a successful transition." | CONFIRMED | `0004-…md:268` |
| 53 | `NextTags` / "next state" appear **nowhere** in RDR 0004 | CONFIRMED | zero grep hits |
| 54 | RDR 0002: escape rule "MUST NOT contain a write block or clear list" | CONFIRMED | `0002-…md:287-288`, verbatim |
| 55 | RDR 0002 escape-row rendering ("row kind `escape`, …") | CONFIRMED | `0002-…md:294-296`, verbatim |
| 56 | RDR 0002 lists "malformed escape declaration" | CONFIRMED | `0002-…md:235`, `:346`, `:720`, `:770` |
| 57 | RDR 0002: clear "renders as a `<clear>` write" | CONFIRMED | `0002-…md:338` |
| 58 | RDR 0001 reserves the Go error path for "parser bugs, IO failures, or programmer mistakes" | CONFIRMED | `0001-…md:213` |
| 59 | RDR 0001 deferral "Kernel assumes a parsed, reviewable table shape" | CONFIRMED | `0001-…md:296` |
| 60 | RDR 0001 closed five-kind taxonomy / A5 | CONFIRMED | `0001-…md:102`, `:106`, `:214` |
| 61 | RDR 0001 "Each refusal kind must be stable enough … without inspecting error strings" | CONFIRMED | `0001-…md:222-223` |
| 62 | RDR 0001 benchmarking quote | CONFIRMED | `0001-…md:508-509` |
| 63 | RDR 0001 `artifacts/triage.md` "Noted, not filed" entry | CONFIRMED — **and it contradicts RDR 0009** | `triage.md:126-127` — "`escapeRow` (`fixtures_test.go:249`) could drop **`Writes`/`NextTags`**" |
| 64 | RDR 0005 "without requiring new exit-code groups…" | CONFIRMED | `0005-…md:114` |
| 65 | RDR 0005 "Internal parse or invariant failures use the existing internal/user-env mapping" | CONFIRMED | `0005-…md:311` |
| 66 | RDR 0005 `flow set-state` driven by planned owned-tag writes | CONFIRMED | `0005-…md:72`, `:142` |
| 67 | RDR 0007 rejected minting `guard_input_missing` | CONFIRMED | `0007-…md:1488`, `:1667` |
| 68 | **RDR 0006 "by its own boundary, leaves parse/render fidelity and row-shape validation with RDR 0002"** | **PARTIALLY REFUTED** | `0006-…md:407` says only "Parse/render fidelity remains owned by RDR 0002"; `:459` lists boundaries with no row-shape clause. `rg 'row-shape|row shape'` over RDR 0006 → **zero hits** |
| 69 | A6: SCXML IRP test 313 via uSCXML `TESTS.md` | NOT-VERIFIABLE LOCALLY | cited path is outside this repo and not reachable from any working directory here. Recorded honestly; A6 is non-load-bearing by the RDR's own framing |

### Step 0 tally

**REFUTED: 6** (#19, #20, #21, #22, #26, #39) — plus **1 partially refuted** (#68).
**NOT-FOUND: 0.**
**Unverifiable / partially verified: 4** (#45, #48, #51, #69).

#19-#22 are four surface-manifestations of one underlying false codebase fact (`escapeRow` does set `Writes`), collapsed into F-1.

## Step 0b — Inverse sibling check

1. **Exported structural conformance predicate returning `error`.** Searched, none in `internal/resolve`. Repo-wide the only exported validator is `internal/cli/respond/respond.go::ValidateMode` — a real near-miss sibling the RDR never cites (returns `*clierr.CLIError`, not `error`). See F-4.
2. **A typed domain error in `internal/resolve`.** Searched, none exists. `resolve.go` imports only `slices` and `strings` — no `errors` import at all.
3. **A package-level sentinel (`var ErrXxx`).** Searched, none exists repo-wide (`rg 'var Err[A-Z]'` → 0). The repo classifies errors solely via `clierr.CLIError.Code` strings.
4. **An `errors.Join` aggregate.** Searched, none exists repo-wide (0 hits). Would be the repo's first; `clierr.CLIError` has no field able to carry N structured identities. See F-3.
5. **An existing table-shape validation entry point in the kernel.** Searched, none exists — matches the RDR's own claim (#13, CONFIRMED).

## Step 1 — 12 verification questions

Codebase-answerable (must read source): Q1-Q7, Q11 — 8 of 12.

1. Does `fixtures_test.go::escapeRow` set `Writes`?
2. How many escape rows in the frozen suite would breach the proposed predicate?
3. Does `config.go::Load`'s `GroupInternal` branch carry a `Hint`?
4. Does the frozen suite contain zero nil-sensitive assertions?
5. Does any file outside the kernel tests + RDR 0007 spike import `internal/resolve`?
6. Does `clierr.ExitCodeFor` return 2 for `GroupInternal` and 1 for a non-`CLIError`?
7. Does `clierr.CLIError` have any field able to carry N structured row identities?
8. Does RDR 0006 anywhere assign "row-shape validation" to RDR 0002?
9. Does RDR 0009 specify what `Result` `Resolve` returns on breach, consistent with `Result`'s doc invariant?
10. Does the RDR specify the exported predicate's receiver/granularity?
11. Do the cited Go stdlib precedents exist in the local toolchain?
12. Does the RDR reconcile `errors.Join` aggregation with `errors.As` row recovery?

## Step 2 — Independent answers

**A1 (Q1).** Yes. `internal/resolve/fixtures_test.go:257` — `Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}},` inside the `escapeRow` return literal (lines 249-260), immediately after `NextTags:` at :256. Corroborated by `docs/rdr/0001-resolution-kernel/artifacts/triage.md:126` and by RDR 0009's own spike edit list.

**A2 (Q2).** `rg -n 'escapeRow\(' internal/resolve/` → 17 lines: definition at `fixtures_test.go:249` plus 16 call sites (`adversarial_test.go:128,171,193,227`; `resolve_test.go:539,553,554,566,598`; `fixup_test.go:43,79,101,125,215,216,284`). Since the builder sets `Writes`, **all 16 constructed escape rows breach** `len(Escape)!=0 && len(Writes)!=0`. The two `escape.Writes = …` overrides merely overwrite an already-present value.

**A3 (Q3).** No. `config.go::Load:84-90` constructs `Code`, `Message`, `Detail`, `Group`, `Cause` — no `Hint`. Repo-wide every `Hint` is on a `GroupUserEnv` error (`config.go::Discover:64`, `root.go::cobraErrorToCLIError:122`).

**A4 (Q4).** No. ~26 nil-sensitive assertions, e.g. `resolve_test.go:78-83` (`case got.Plan != nil && got.Refusal != nil:`, `case got.Plan == nil && got.Refusal == nil:`, `if (got.Refusal != nil) != got.Refused()`), plus `mvv_test.go:42,45,99,102,140`. Not nil-sensitive about `Writes`: the only presence test is `mvv_test.go:110` `len(got.Plan.Writes) > 0`, length-based.

**A5 (Q5).** No. Six importers total, all tests/spike. `cmd/intrastate/main.go` is `func main() { cli.Execute() }`; `NewRootCmd` adds only `newVersionCmd()`.

**A6 (Q6).** Yes to both. `clierr.go::ExitCodeFor:119-120` `case GroupUserEnv, GroupInternal: return 2`; final `return 1` documented as "A non-CLIError defaults to exit 1".

**A7 (Q7).** No. `CLIError` fields: `Code`, `Message`, `Param`, `Detail`, `Hint` (all string), `Group`, `Cause error` (single). No `Rows`, no `Errors []error`. `Cause` is `json:"-"` — "Not serialized — the wire-visible cause surface is Detail." N identities can reach the wire only as `Detail` prose.

**A8 (Q8).** No. Zero grep hits for "row-shape"/"row shape" in RDR 0006. Its boundary statements are `:407` "Parse/render fidelity remains owned by RDR 0002" and `:459` "table source and normalization stay in RDR 0002, …".

**A9 (Q9).** Partially silent and internally conflicting. Technical Design says "`Resolve` returns a nil `Result`" — but the signature is `func Resolve(in Input) (Result, error)`, by value; no nil `Result` exists. MVV fixture 1 says `Result{Plan: nil, Refusal: nil}`. Meanwhile `resolve.go::Result` documents "never both and never neither", asserted at `resolve_test.go:80`. The RDR never states it is narrowing that invariant.

**A10 (Q10).** Silent. The contract pins name shape, return type, doc form, scope, ordering, aggregation — never the receiver or signature. Both cited precedents are value-receiver methods, implying `func (t Table) CheckValid() error`, but `Table` is an RDR 0001 type and Overrides claims "Neither peer's normative surface is reopened". Load-Bearing Decisions state a per-row predicate while the scope pin is whole-table; neither is designated the exported surface. Validation scenario 7 covers drift behaviorally.

**A11 (Q11).** Yes, all checkable ones exist (see #43-#51). Two sub-claims not asserted: pprof's "re-checked at the parse boundary" call graph, and `analysis.Validate`'s fail-fast semantics.

**A12 (Q12).** Silent on mechanics. A single-target `errors.As` on a multi-branch `errors.Join` resolves to the first matching branch and discards the rest. Failure Modes ("read the row identities off the error value — `errors.As`/`AsType`") reads as if `As` recovers all; scenario 9 more carefully requires iterating `Unwrap() []error`. Also unpinned: whether a single breach returns the bare typed error or a one-element Join (changes `As`/`Unwrap` behavior at N=1 vs N≥2), and whether the sentinel is wrapped once at the aggregate or once per row-error.

## Step 3 — Findings

### F-1 — (a) codebase claim REFUTED — A3's scope correction is false, and it inverts Phase 2

**RDR anchor** (A3, "Scope correction (found by this spike)"):
> `internal/resolve/fixtures_test.go::escapeRow` populates **`NextTags` only — it sets no `Writes` at all**. The write-bearing escape rows are two *call-site overrides*

and Phase 2: "The `escapeRow` builder needs no change: it sets `NextTags` only"; and the Infra Audit row "Builder sets `NextTags` only".

**What's wrong.** `escapeRow` **does** set `Writes`. `internal/resolve/fixtures_test.go:257`:
```go
Writes:        []resolve.Tag{{Key: "status", Value: "Blocked"}},
```
inside the `escapeRow` return literal (lines 249-260), one line below `NextTags:` at :256. The builder has **16 call sites**, so **every escape row in the frozen suite breaches the proposed predicate** — not two. The two `escape.Writes = …` overrides merely overwrite an already-present value.

Three consequences, increasing in severity:
- **Phase 2 as written would not produce a green suite.** Dropping `Writes` only at the two overrides leaves 14 breaching escape rows; under the new entry precondition every test using them errors at `Resolve` — the opposite of A3's "changes no frozen assertion outcome."
- **The RDR contradicts its own evidence.** `evidence/spikes/a3-fixture-conformance.md` lists three edits, the first being removal of that exact `Writes` line from `escapeRow`, and warns "escapeRow is not the only place an escape row gets `Writes`. Stripping only the builder would have left two escape rows still carrying writes." The spike is correct; the RDR body inverted the sentence.
- **It contradicts the RDR 0001 disposition it claims to supersede.** `docs/rdr/0001-resolution-kernel/artifacts/triage.md:126-127`: "`escapeRow` (`fixtures_test.go:249`) could drop **`Writes`/`NextTags`** for fixture realism."

The A3 *verdict* (154 PASS, identical outcome sets) is unaffected — the spike edited all three sites. What is broken is the RDR's prose scope description and the Phase 2 instruction derived from it.

**Source cite.** `internal/resolve/fixtures_test.go:249-260` (`Writes` at :257); `evidence/spikes/a3-fixture-conformance.md` §"What was edited"; `docs/rdr/0001-resolution-kernel/artifacts/triage.md:126`.

### F-2 — (a) codebase claim REFUTED — `config.Load` is not a precedent for `GroupInternal` + Hint

**RDR anchor** (Normative Contracts, CLI wrap clause):
> the verb MUST wrap it into a `*clierr.CLIError` carrying a stable `Code`, `Group GroupInternal` (exit 2), the offending row identity, and a `Hint` stating the remedy — the shipped `config.Load` wrapping pattern

**What's wrong.** `internal/cli/config/config.go::Load:84-90` constructs the `GroupInternal` error with `Code`, `Message`, `Detail`, `Group`, `Cause` — and **no `Hint`**. Repo-wide, every `Hint` sits on a `GroupUserEnv` error (`config.go::Discover:64`; `root.go::cobraErrorToCLIError:122`). **No shipped `GroupInternal` error in this repo carries a `Hint`.** The `Hint` requirement may still be good, but it is new, not precedented, and the RDR should not lean on `config.Load` to authorize it — which also weakens the implicit claim that the wrap is a mechanical copy of an existing site.

**Source cite.** `internal/cli/config/config.go::Load:84-90`; `internal/cli/clierr/clierr.go::CLIError`.

### F-3 — (d) RDR is silent — the aggregate breach report has no carrier through the CLI envelope

**RDR anchor** (Normative Contracts, aggregate clause):
> When a table carries MORE THAN ONE breaching row, the precondition MUST report every one of them in a single pass … combined with `errors.Join`

crossed with the CLI clause requiring a `*clierr.CLIError` "carrying … the offending row identity", and Validation scenario 11.

**What's wrong / missing.** `clierr.CLIError` has no field able to carry N structured row identities: `Code`, `Message`, `Param`, `Detail`, `Hint` (all string), `Group`, `Cause error` (single). `Cause` can hold the `errors.Join` aggregate for in-process `errors.As`, but it is `json:"-"` ("Not serialized — the wire-visible cause surface is Detail"). On the wire, N identities can only be flattened into `Detail` prose — exactly what the RDR's typed-error contract forbids ("Row identity MUST NOT be recoverable only from formatted prose") and what Failure Modes promises against ("Operators branch on the code and tests assert on `RowRef` values — neither reads error prose").

The RDR is silent on the resolution. Three unreconciled options: (i) `CLIError` gains a structured field — an RDR 0005-owned envelope change, contradicting A2's "no RDR 0005 contract change"; (ii) the wire form degrades to `Detail` prose while only in-process callers get structure; (iii) the aggregate collapses to one row at the CLI. The clierr package comment invites (i) ("Extend with new optional fields as needed"), but that surface is RDR 0005's.

**Source cite.** `internal/cli/clierr/clierr.go::CLIError` (fields; `Cause` `json:"-"`); `clierr.go::EmitJSON` (serializes the struct as-is).

### F-4 — (b) new rule with an existing sibling — `respond.ValidateMode` is the repo's only exported validator and contradicts the imported convention

**RDR anchor** (Normative Contracts, exported-predicate clause):
> It MUST return `error` (nil when valid), and its name MUST follow Go's error-returning convention — `CheckValid` or `Validate` … Precedent: … `pprof.Profile.CheckValid` … and `rsa.PrivateKey.Validate`

**What's wrong.** The RDR reaches for stdlib precedent while the repo already has exactly one exported validator, uncited: `internal/cli/respond/respond.go::ValidateMode` — `func ValidateMode(cmd *cobra.Command) *clierr.CLIError`, doc "returns a CLIError when `--as` carries an unrecognized value. Verbs call this at the top of `RunE` to fail fast." It matches the naming rule but returns a **`*clierr.CLIError`, not a bare `error`**, and AGENTS.md reinforces that direction. The RDR's `MUST return error` therefore diverges from the only in-repo precedent without weighing it. The divergence is probably correct (the kernel must not import `clierr` — `internal/resolve/boundary_test.go` polices the kernel's dependency surface), but that argument appears nowhere, so a reader cannot tell whether the sibling was rejected or missed.

Also from the same sweep, as unacknowledged scope: `internal/resolve/resolve.go` imports only `slices` and `strings` — **no error type, no sentinel, no `errors` import**; and repo-wide there is **zero** `errors.Join` usage and **zero** `var Err[A-Z]` sentinels. The typed error, category sentinel, and `errors.Join` aggregate would each be the repo's first, in a codebase standardized on `clierr.CLIError.Code` strings. "Searched, none exists" for all three.

**Source cite.** `internal/cli/respond/respond.go::ValidateMode:89-108`; `internal/resolve/resolve.go` imports (lines 13-16); `AGENTS.md`.

### F-5 — (c) internal contradiction — "nil `Result`" vs `Result`'s frozen never-neither invariant

**RDR anchor** (Technical Design):
> On the first row with non-empty `Escape` and non-empty `Writes`, `Resolve` returns a nil `Result` and a non-nil error naming the row.

versus MVV fixture 1: "`Resolve` returns `Result{Plan: nil, Refusal: nil}` and a non-nil error".

**What's wrong.** Three problems in one sentence.

1. **"a nil `Result`" is not constructible.** `resolve.go::Resolve` is `func Resolve(in Input) (Result, error)` — `Result` is returned **by value**. The MVV's form is the only realizable one.
2. **"On the first row" contradicts the aggregate contract.** The same document requires reporting "every one of them … not the first alone" and "aggregate, not fail-fast." This is pre-aggregate text left standing.
3. **The zero-disposition return contradicts a frozen invariant the RDR never narrows.** `resolve.go::Result` documents "exactly one transition plan or exactly one typed refusal, **never both and never neither**", asserted directly at `internal/resolve/resolve_test.go:80` (`case got.Plan == nil && got.Refusal == nil:` → failure). `Result{nil, nil}` is the "neither" state. Either the `Result` doc contract must be narrowed (an RDR 0001 surface edit not recorded in Overrides) or the RDR must scope the invariant to nil-error returns. It does neither.

**Source cite.** `internal/resolve/resolve.go::Resolve:318` and `resolve.go::Result:281-288`; `internal/resolve/resolve_test.go:78-83`.

### F-6 — (a) codebase claim REFUTED — "a frozen suite with zero nil-sensitive assertions"

**RDR anchor** (Load-Bearing Decisions, *Nil-vs-empty*):
> a frozen suite with **zero** nil-sensitive assertions and four length-based `Writes` checks

**What's wrong.** The suite contains ~26 nil-sensitive assertions on `Plan`/`Refusal`, e.g. `internal/resolve/resolve_test.go:78-83`:
```go
case got.Plan != nil && got.Refusal != nil:
case got.Plan == nil && got.Refusal == nil:
if (got.Refusal != nil) != got.Refused() {
```
plus `mvv_test.go:42,45,99,102,140` and ~18 more. False as written.

The *defensible* narrowed claim — zero nil-sensitive assertions **on `Writes`** — does hold: the only `Writes` presence test is `mvv_test.go:110` `len(got.Plan.Writes) > 0`, length-based; and "four length-based `Writes` checks" is exactly right (`resolve_test.go:668`, `:709`, `:986`, `mvv_test.go:110`). The nil-vs-empty decision survives on the narrowed claim; only the overbroad wording fails — and it is one of three named grounds for the length-based predicate, and the one a reviewer would check first.

**Source cite.** `internal/resolve/resolve_test.go:78-83`; `mvv_test.go:42,45,99,102,110,140`; `rg 'len\([a-zA-Z.]*Writes\)' internal/resolve/` → exactly 4.

### F-7 — (a) peer-quote PARTIALLY REFUTED — RDR 0006 never says "row-shape validation" stays with RDR 0002

**RDR anchor** (Technical Environment, RDR 0006 bullet):
> by its own boundary, leaves parse/render fidelity and row-shape validation with RDR 0002 — so this RDR's load-time enforcement assignment does not overlap its authority

**What's wrong.** Half the attributed phrase is absent from RDR 0006. `rg 'row-shape|row shape'` over that file → **zero hits**. Its actual boundary statements are `:407` "Parse/render fidelity remains owned by RDR 0002" (quoted correctly) and `:459` "table source and normalization stay in RDR 0002, predicate semantics stay in RDR 0003, accessor safety stays in RDR 0004, and CLI presentation stays in RDR 0005."

Why this matters rather than being a wording nit: this bullet is the RDR's sole argument that its load-time enforcement does not overlap RDR 0006's authority, and the non-overlap rests on the half that does not exist. RDR 0006 claims blocking lint authority over the normalized model (`:11`), and its A1 (`:98-108`) records that normalized rows expose **writes** to lint ("RDR 0002 `Technical Design` defines normalized candidate rows with source locator, predicates, and writes"). A write-bearing escape row is therefore visible to RDR 0006's lint. The overlap question is real and currently answered by a misquote.

**Source cite.** `docs/rdr/0006-graph-lint-authority-and-guarantees.md:407`, `:459`, `:11`, `:98-108`; zero-hit grep for "row-shape"/"row shape".

### F-8 — (d) RDR is silent — `errors.As` on an `errors.Join` aggregate recovers one breach, not all

**RDR anchor** (Failure Modes):
> Diagnosis: read the row identities off the error value — `errors.As` / `errors.AsType` for callers and tests

crossed with: "combined with `errors.Join`, so each per-row error stays individually inspectable through the aggregate's `Unwrap() []error`".

**What's wrong / missing.** A single-target `errors.As` against a multi-branch `errors.Join` tree resolves to the **first matching branch** and silently discards the rest. For the multi-breach case the promised diagnosis idiom returns one `RowRef` while the contract promises every breaching row. Validation scenario 9 is more careful ("via `errors.As`/`AsType` **over `Unwrap() []error`**"), implying the caller must type-assert to `interface{ Unwrap() []error }` and iterate — but the RDR never states that a plain top-level `errors.As` is insufficient, and Failure Modes reads as though it suffices.

Two further unpinned points: (i) whether a **single** breach returns the bare typed error or a one-element `errors.Join` — this changes `errors.As`/`Unwrap()` behavior at N=1 vs N≥2, and the parenthetical "degrades to the single error's own message for one breach" addresses the *message*, not the *type*; (ii) whether the category sentinel is wrapped once at the aggregate level or once per row-error. Given F-4 (zero prior `errors.Join` usage and zero sentinels in the repo, so implementers have no local pattern), leaving these unpinned invites exactly the drift this RDR exists to prevent.

**Source cite.** RDR text as quoted (Failure Modes vs Normative Contracts vs scenario 9); `rg 'errors\.Join' --glob '*.go' internal/ cmd/` → zero hits; `internal/cli/clierr/clierr.go::ErrorCode` (repo's only `errors.As` idiom, single-target).

### F-9 — (d) RDR is silent — the exported predicate's receiver and granularity are unspecified

**RDR anchor** (Normative Contracts, exported-predicate clause):
> The conformance predicate MUST be exported by the kernel package as a construction-time check callable by any table producer, and `Resolve`'s entry precondition MUST be that same function — one predicate, two call sites

**What's wrong / missing.** The RDR pins the name shape, return type, doc form, scope, ordering, and aggregation — but never the **receiver or signature**. Both cited precedents are methods on the validated value (`pprof.Profile.CheckValid`, `rsa.PrivateKey.Validate`), implying `func (t Table) CheckValid() error`; but `Table` is an RDR 0001 type and adding an exported method is a surface change Overrides does not record (it claims "Neither peer's normative surface is reopened"). Compounding this, Load-Bearing Decisions state a **per-row** predicate (`len(row.Escape) != 0 && len(row.Writes) != 0`) while the scope pin is **whole-table** — two granularities, neither designated the exported surface, and no word on whether a per-`Row` form is also exported for producers building rows incrementally.

This is a real gap rather than an implementation detail because the clause's whole purpose is that "the two enforcement points cannot drift" — a future producer can only call the same function if the RDR says what that function's signature is. Validation scenario 7 covers drift behaviorally, which mitigates but does not close it.

**Source cite.** `internal/resolve/resolve.go::Table:203-213` (currently one unexported method, `models`); `$GOROOT/src/crypto/rsa/rsa.go:233` and `$GOROOT/src/cmd/vendor/github.com/google/pprof/profile/profile.go:362` (both value-receiver methods — the shape the precedents imply but the RDR never states).
