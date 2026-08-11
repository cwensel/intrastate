Model: claude-opus-5[1m]
Lens: cove (dispatcher-side grounding, run alongside passes A and B)

Independent Step-0 sweep run in the dispatching context so lens findings could be
adjudicated without a second round-trip. Every entry below was read from source at
`via-claude` HEAD (`20cc49f`).

## Step 0 — Grounding sweep

| RDR claim | Verdict | Source |
| --- | --- | --- |
| `resolve.go::Row` carries `RuleID`, `SourceLocator`, `Escape`, `Writes`, `NextTags` | CONFIRMED | `internal/resolve/resolve.go::Row` — all five fields present; `RuleID`/`SourceLocator` doc'd as "the source identity RDR 0002 requires" |
| `Resolve` returns literal `nil` error on every path | CONFIRMED | `resolve.go::Resolve` — `refuse(...), nil` ×2, `Result{Plan: planOf(...)}, nil`, and the `escapeOrRefuse` arms; no helper returns an error type |
| `planOf` copies `Writes` unconditionally, escape included | CONFIRMED | `resolve.go::planOf` — `Writes: copyTags(row.Writes)`, `Escaped: escaped`; no escape-awareness |
| `copyTags` is nil-preserving | CONFIRMED | `resolve.go::copyTags` — `if in == nil { return nil }` |
| Escape discriminator is `len(row.Escape) != 0` / `rescues` | CONFIRMED | `resolve.go::Resolve` candidate loop; `resolve.go::rescues` (`slices.Contains(r.Escape, kind)`) |
| `Resolve` already walks `Table.Rows` twice + `gate` | CONFIRMED | `Resolve` candidate loop, `escapeOrRefuse` escape loop, `gate(rows …)` |
| `RowRef` fields; `compareRefs`; `rowRefs` sorts | CONFIRMED (with a gap — see F-2) | `resolve.go::RowRef` = `{RuleID, SourceLocator}`; `compareRefs` compares those two; `rowRefs` uses `slices.SortFunc` |
| `adversarial_test.go:229` assigns `escape.Writes = {status Escaped}` | CONFIRMED | exact line 229 in ADV-2b |
| `fixup_test.go:103` assigns the identical override | CONFIRMED | exact line 103, Fixup-1e "missing owned state" |
| **`fixtures_test.go::escapeRow` "populates `NextTags` only — sets no `Writes` at all"** | **REFUTED** | `fixtures_test.go:257` — `Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}},` inside the builder. 16 call sites inherit it. See F-1 |
| A5: no production file imports `internal/resolve` | CONFIRMED | only `internal/resolve/*_test.go` + `docs/rdr/0007-…/evidence/spikes/aggregation-probe_test.go`; `cmd/intrastate/main.go` is `cli.Execute()` only |
| `root.go::NewRootCmd` registers only `version` | CONFIRMED | `internal/cli/root.go::NewRootCmd` |
| `clierr::ExitCodeFor` — `GroupUserEnv, GroupInternal → 2`; non-`CLIError` → 1 | CONFIRMED | `internal/cli/clierr/clierr.go::ExitCodeFor`, incl. doc line "A non-CLIError defaults to exit 1" |
| `config.go::Load` is the shipped wrap precedent | CONFIRMED | `internal/cli/config/config.go::Load` |
| Frozen suite has **zero** nil-sensitive `Writes` assertions and **four** length-based ones | CONFIRMED exactly | `mvv_test.go:110`, `resolve_test.go:668`, `:709`, `:986` — all `len(...)`; no `Writes == nil` / `!= nil` anywhere |
| AGENTS.md "Add new codes as needed" | CONFIRMED | `AGENTS.md:25-27` (phrase wraps across lines 26–27) |
| RDR 0004: "apply **only** planned owned-tag writes … MUST NOT write observed or recognized tags" | CONFIRMED | `docs/rdr/0004-*.md:268-269` verbatim |
| RDR 0004: `NextTags` / "next state" appears **nowhere** | CONFIRMED | case-insensitive grep over RDR 0004 → 0 matches |
| RDR 0002 escape prohibition; `<clear>` write; "malformed escape declaration"; escape-row rendering | CONFIRMED | `0002-*.md:287`, `:338`, `:235`/`:346`, `:294` |
| RDR 0001 reserved-error-path sentence | CONFIRMED | `0001-*.md:213` |
| RDR 0005 "without requiring new exit-code groups" | CONFIRMED | `0005-*.md:114` |
| RDR 0007 rejected `guard_input_missing` | CONFIRMED | `0007-*.md:1488`, `:1667` |
| `rsa.PrivateKey.Validate` doc form — "returns nil if the key is valid, or else an error describing a problem" | CONFIRMED verbatim | `go doc crypto/rsa.PrivateKey.Validate` (go1.26.5) |
| `slogtest.TestHandler` aggregates via `errors.Join` | CONFIRMED | `go doc testing/slogtest.TestHandler` — "combined into a single error with errors.Join" |
| `fstest.TestFS` aggregate-style | CONFIRMED | `go doc testing/fstest.TestFS` |
| `json.UnmarshalTypeError.Field` (location field precedent) | CONFIRMED | `go doc encoding/json.UnmarshalTypeError` |

## Step 0b — Inverse sibling check

The RDR adds a new exported structural-validation predicate + typed error + sentinel.
Searched production code for an existing sibling that already makes this decision:

- **Validation-entry sibling**: none in the kernel — `Resolve` performs no structural
  check on `Table.Rows` before evaluation (read `Resolve` start-to-finish). CONFIRMS the
  RDR's "searched, none exists".
- **Naming sibling (not cited by the RDR)**: `internal/cli/respond/respond.go::ValidateMode`
  is the repo's one existing `Validate*` function. It returns `*clierr.CLIError`, not
  `error` — so it does not contradict the RDR's `CheckValid`/`Validate` naming clause, but
  it is the in-repo precedent that clause reasons about purely from stdlib. Non-blocking.
- **`errors.Join` / sentinel**: no existing use anywhere in production code
  (`grep errors.Join` → none; `= errors.New` → none). The aggregate shape is genuinely new.
- **Typed-error sibling**: `clierr.CLIError` with `Unwrap`/`Cause` is the only typed error
  surface; it is a CLI-layer type, not a kernel one. No overlap.

## Step 3 — Findings

### F-1 (type a: codebase claim REFUTED) — A3's "Scope correction" inverts its own spike

**RDR anchor** (appears twice — A3 *Scope correction* and Phase 2):
> `internal/resolve/fixtures_test.go::escapeRow` populates **`NextTags` only — it sets no `Writes` at all**.

and

> The `escapeRow` builder needs no change: it sets `NextTags` only

**Source**: `internal/resolve/fixtures_test.go:249-260` — the builder body includes
`Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}},` at line 257. `git log` shows
`fixtures_test.go` unchanged since `2e2112b test(resolve): add RDR 0001 resolution-kernel
spec tests`, so this is not post-spike drift.

**Corroboration from this RDR's own evidence** — `evidence/spikes/a3-fixture-conformance.md`
"What was edited (variant 1: conformance)" lists **three** edits, the first being
"`internal/resolve/fixtures_test.go`, `escapeRow` — removed the line `Writes: …{{Key:
"status", Value: "Blocked"}},`", and states explicitly: "Edits 2 and 3 matter: `escapeRow`
is not the only place an escape row gets `Writes`." The spike's Verdict repeats it: the two
overrides "must be stripped **alongside the `escapeRow` builder**". The RDR text reverses
the spike's finding.

**Blast radius**: `escapeRow` has **16 call sites** across `adversarial_test.go`,
`resolve_test.go`, `fixup_test.go` — every one currently produces a breaching escape row.
Phase 2 as written ("drop `Writes` at the two call-site overrides … the builder needs no
change") would leave 16 breaching rows and the whole frozen suite erroring at `Resolve`
entry the moment Phase 1 lands.

**Note**: the A3 *verdict* is unaffected — the 154-PASS conformed run was measured with all
three edits applied. Only the scope description and Phase 2 inherit the error.

### F-2 (type d: silence revealing a missing requirement) — `compareRefs` is not a total order over distinct rows

**RDR anchor**:
> The reported rows MUST be ordered by RowRef identity using the kernel's existing compareRefs ordering, never by Table.Rows position, so the diagnostic payload is a function of the table value rather than of row order

and Testing Strategy scenario 10:
> The same multi-breach table supplied with its rows in two different orders. **Expected**: Identical reported row sequence

**Source**: `resolve.go::compareRefs` compares only `RuleID` then `SourceLocator`;
`resolve.go::rowRefs` sorts with `slices.SortFunc`, which is **not** stable
(`slices.SortStableFunc` is the stable variant). Two distinct breaching rows sharing a
`RuleID`+`SourceLocator` pair — trivially constructible on hand-built rows, which is exactly
this RDR's target population (A5: the only producers are hand-built fixtures), and the
zero-value case `RowRef{"",""}` for rows built without source identity — compare equal, so
their relative order is unspecified and scenario 10 can fail nondeterministically.

The existing uses are narrower: `rowRefs` today reports *ambiguous-match candidate rows*,
which RDR 0002 gives distinct rule ids. This RDR extends the ordering to a new population
without stating the uniqueness precondition it inherits.

**RDR is silent on**: whether `RowRef` identity is unique per row, what the report does with
duplicate-identity breaching rows, and whether hand-built rows lacking `RuleID`/
`SourceLocator` are in scope. Minimum fix: pin `SortStableFunc` (making table order the
documented tiebreak) or state the uniqueness precondition and what happens when it fails.

### F-3 (type c: internal contradiction) — Infrastructure Audit row contradicts the Normative Contracts

**RDR anchor** — Existing Infrastructure Audit, "Error envelope" row, *Spec Impact* column:
> Breach is diagnosable from the error text (A7), not a new code

**Contradicted by** the Normative Contract:
> the verb MUST wrap it into a *clierr.CLIError carrying a stable Code

and by A7's own updated Status ("the structured form is **adopted**, not merely permitted"),
and by Failure Modes:
> the runbook entry is therefore "branch on the code; the named rows identify the broken producer," **not** "read the error text."

The audit row is stale text from the pre-Resolve draft, when the diagnostic gap was accepted.
It now states the opposite of the contract it summarizes, on both halves ("from the error
text" and "not a new code").

### F-4 (type d: silence, minor) — `errors.Join` single-error parenthetical is imprecise

**RDR anchor**:
> `errors.Join` costs nothing (nil when all-nil; degrades to the single error's own message for one breach)

**Source**: `go doc errors.Join` (go1.26.5) — "A non-nil error returned by Join implements
the `Unwrap() []error` method." A single-error `Join` returns a `*joinError` wrapper, not the
bare error. The *message* does read as the single error's message, so the RDR's claim is true
as stated about message text; but a caller doing `err == ErrX` or a type assertion (rather
than `errors.Is`/`As`) would not see through it. The RDR's contracts consistently require
`errors.Is`/`As`, which do traverse — so this is a precision note, not a defect. Non-blocking.
