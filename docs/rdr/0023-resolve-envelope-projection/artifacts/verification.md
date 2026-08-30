# Verification — RDR 0023 Resolve-envelope projection opt-out

Stage-8 independent verification artifacts. Each phase appends its own
section; sections are append-only and never rewritten.

## Phase 3a — CoVe

Independent Chain-of-Verification against the committed implementation.
Method: for each normative obligation, an input was chosen that would make a
violation VISIBLE, then executed against the real binary or the real
production symbols. No Phase 1 test file was read; every probe was derived
from `0023:C1`, `0023:C2` and the REQ list alone.

**Tree under verification.** All probes ran against committed `HEAD`
(`597a9d9`), materialized into an isolated clean checkout via `git archive`,
because the shared worktree carried concurrent uncommitted edits from the
parallel Phase 3b agent. Verifying against a mutating tree would have been
unsound; the clean checkout is byte-equal to the committed implementation.

### Result

**No FAIL-N entries. Every normative obligation probed held.**

### Coverage

| Obligation | REQ | Probe | Result |
| --- | --- | --- | --- |
| Flag shape: bool, default false, no shorthand | 1, 15, 16 | pflag introspection over the walked tree | PASS — `type=bool default=false shorthand=""` |
| Non-registration on every other command | 2, 12 | live `--plan-only` on `next`/`read-state`/`set-state`/`flow`/root/`docs`/`lint`/`version`/`help`/`completion` | PASS — `command-error`, exit 2 everywhere |
| Whole-tree structural walk = `{flow resolve}` | 3, 5, 13 | independent total walker (no name skip, no Hidden gate) written from the clause | PASS — 15 commands walked, registrants `[intrastate flow resolve]` |
| Vacuity guard: four flow verbs exist | 4 | assertion over walked set | PASS |
| Forbidden walker is strictly weaker | 6, 14 | contrast run of `help_all.go::walkCommandTree` on the same root | PASS — reaches 9 of 15, drops the whole `completion` subtree |
| Auto-generated commands materialized + asserted | 7, 8, 9 | `InitDefaultHelpCmd` + `InitDefaultCompletionCmd`, then precondition check | PASS — `completion` present in walked set |
| No new refusal code / exit group | 11 | refusal sweep across six classes | PASS — only pre-existing codes |
| Projected key set exactly the plan group | 17, 18, 22, 23 | 8 fixture shapes, key-set comparison | PASS — dropped set is exactly `{model,observed,owned,readers,outcome}` |
| Absent, never null, never placeholder | 19, 113, 123 | grep for `null` + decode of projected wire | PASS — no `null`, no stand-ins |
| Byte-identity of carried values | 20, 26, 34 | per-key `json.dumps` comparison, 8 shapes | PASS — every carried key byte-equal |
| Carried keys keep default relative order | 21, 121 | order-preservation check vs default key list | PASS — `revision` first, order is default-minus-deleted |
| Presence rules unchanged inside plan group | 24, 25 | `emit` `{}` case; `escape_class` present and absent cases | PASS — identical in both widths |
| Default-mode byte-identical to pre-change | 27, 95, 96, 97, 129 | checked-in golden vs live default | PASS — exact match; golden commit `d3369f2` verified ancestor of conversion `c47ab2b` |
| 0011 fences intact | 29, 30 | `flow next --all` works; `resolve --all` rejected | PASS |
| Report-only: identical exit codes | 31, 32, 45 | 6 refusal classes ± flag | PASS |
| Byte-identical refusal envelopes | 33, 100, 131 | unmodeled-outcome, tag-invalid, guard-unevaluable, gate-denied, artifact-missing, missing-outcome | PASS — including full `findings[]` arrays |
| Strict key-subset, values byte-identical | 34, 133 | automated differential, 8 shapes | PASS |
| Identical invoked-reader set, observed at execution | 35, 40, 41, 42, 43 | `readerExecutionHook` seam driven ± flag on a live reader model | PASS — records `[release-state]` (non-empty, so discriminating) in both runs |
| Strictly shorter, full emitted line, both modes | 36, 38, 39, 57 | byte table, 5 shapes × 2 modes | PASS — 46–59% saved, strictly shorter in all 10 |
| Projection site: after last `respond.Fail`, before `respond.OK` | 47, 48 | source position check | PASS — last Fail line 262, projection line 334 |
| Flag read at exactly one lexical site | 49, 53 | grep for every `GetBool` of the flag in production code | PASS — one site, `flow_projection.go:76` |
| No "if refusing, skip projection" guard | 50 | inspection of projection file | PASS — none; refusal blindness is structural |
| Projection separate from assembly, assembly flag-blind | 51, 52 | signature + call-site inspection | PASS — takes assembled payload, returns projected |
| Text: subset, byte-identical per line | 54, 55, 56, 58, 110 | per-line membership check | PASS — 15 → 9 lines, 6 dropped, zero additions, all byte-identical |
| Text determinism | 59, 122 | 10 runs each mode, sha256 | PASS — one hash per mode |
| Partition completeness over the struct | 60, 68, 69 | reflection over `resolvePayload` vs `resolvePayloadGroups` | PASS — 14/14 fields, exactly one entry each, sets equal |
| Declared groups match the census | 61, 62, 63 | entry-by-entry comparison | PASS |
| Field count preserved at 14 | 64, 90 | `reflect.NumField` | PASS |
| Carrier is a standalone table, not a struct marker | 65, 66, 70, 71 | `flow_partition.go` is a separate declaration site | PASS — projection reads `echoFieldNames()`, no list of its own |
| Projection implemented FROM the declaration | 67 | `projectAwayEchoGroup` drives off `echoFieldNames()` | PASS |
| Group membership is the table's, never a Go type's | 72, 89 | `escape_class` is plain `string` + `omitempty`, PLAN, always-keep | PASS — survives projection |
| Always-keep core carried | 73, 74, 76 | `rule`, `escaped`, `escape_class`, `revision` on every shape incl. escaped | PASS |
| Conversion scope: five echo fields only → `*T` | 88, 86 | struct field type/tag dump | PASS — only the five are pointers; all non-nil on success |
| A9 empty-container nil arm | 52, 85 | model with no tags, no owned keys, no readers | PASS — renders `{}`,`{}`,`[]`; no `null` |
| No gateway change | 94 | git log over `internal/cli/respond/` | PASS — untouched |
| MVV step 6 byte table ≥ 40% on checked-in fixtures | 103, 104 | 5 shapes × 2 modes | PASS — min 46%, max 59% |
| Docs regenerated and consistent | 115, 116, 117 | `make docs-check` | PASS — "docs up to date" |
| Help line names the omitted group | 119 | flag usage string | PASS — names `model, observed, owned, readers, outcome` |
| One encoder, HTML escaping disabled | 121, 127 | emit value `a<b&c>d` | PASS — literal, not `<`, in both widths |
| Predecessor suites green, unedited | 29, 135 | `go test ./...` | PASS — whole repo green |

### Discriminability controls run

Negative controls were executed to confirm the oracles can fail, not merely
that they pass.

- **Unassigned field (REQ-69, REQ-79).** Added a field to `resolvePayload`
  with no `resolvePayloadGroups` entry. The completeness assertion detects it
  (`fields with NO declaration entry: [Unassigned]`), and the field was
  confirmed to fall through into the *projected* width — exactly the
  "omit-list whose complement is whatever else exists" C2 forbids. The
  declaration-completeness check is therefore load-bearing, not decorative.
  Reverted after the probe.
- **Reader seam is non-vacuous (REQ-40, REQ-43).** The seam recorded a
  non-empty set (`[release-state]`) on a model with a live reader, so a
  projected run that skipped the reader pass would record `[]` and turn the
  assertion red. This distinguishes it from the `invokedReaders(model,
  outcome)` recomputation C1 rejects by name.
- **Total walker vs forbidden walker (REQ-6).** Over the same materialized
  root, the total walk reaches 15 commands and
  `help_all.go::walkCommandTree` reaches 9, dropping the entire `completion`
  subtree — the six commands the clause's scope exists to cover.

### Measured byte table (full emitted line, trailing newline excluded)

| shape | mode | default | projected | saved |
| --- | --- | --- | --- | --- |
| pricing 2×2 | json | 290 | 152 | 47.6% |
| pricing 2×2 | text | 266 | 129 | 51.5% |
| release begin | json | 437 | 238 | 45.5% |
| release begin | text | 396 | 206 | 48.0% |
| release ship-clean | json | 405 | 191 | 52.8% |
| release ship-clean | text | 356 | 153 | 57.0% |
| release escaped | json | 358 | 156 | 56.4% |
| release escaped | text | 331 | 138 | 58.3% |
| review submit | json | 301 | 162 | 46.2% |
| review submit | text | 254 | 127 | 50.0% |

The pricing 2×2 JSON row reproduces A1's normative fixture exactly
(290 → 152 B). Text reduction exceeds JSON reduction on every shape, as A8
predicted. Every shape clears the MVV 40% bar in both modes.

### Notes

- The escaped shape (`escape_class` present, 9 projected keys) and the
  unescaped shape (8 projected keys) were both exercised, confirming A-2's
  conditional key count and REQ-25's "compare against the same run's default
  output" rule.
- `llms.txt` carries no flag names for any command, so `--plan-only`'s
  absence there is correct rather than a docs gap; `make docs-check` passes.

## Phase 3b — Adversarial

Three failure modes attacked, anchored in the RDR's own `Failure Modes`
list. Tests added in `internal/cli/flow_adversarial_0023_test.go`; every
claim below was established by mutation against the shipped implementation,
not by inspection.

Baseline: `go test ./... -count=1` green before and after the additions.

### ADV-1 — `0023:F2`: the shipped declaration table is never asserted total

**Test**: `TestAdversarialF2_TheShippedDeclarationTableIsTotalOverThePayloadStruct`
**Status**: PASSES today; **sole discriminator** of a real coverage hole.

REQ-68 requires the C2 oracle assert, against the declaration, that "every
struct field has exactly one entry; the entry set and the field set are
equal (neither a field without an entry nor an entry without a field)".
The shipped oracle
(`TestReq60And61And62And63And68And69_EveryPayloadFieldIsAssignedToExactlyOneC2Group`)
asserts those two clauses against a map **authored in the test file**, and
its third clause end-to-end against emitted keys. That is correct for C2's
independence requirement, but it leaves the production carrier
`internal/cli/flow_partition.go::resolvePayloadGroups` unverified *as a
carrier*: nothing asserts the shipped table is itself total over
`resolvePayload`.

The hole is reachable because the two sides are not symmetric.
`projectAwayEchoGroup` drives off `echoFieldNames()`, a map lookup by Go
field name that silently yields the zero `payloadGroup` for an absent name.
So an ECHO field the table loses is caught (an echo key survives the flag
and the end-to-end assertions fire), while a **PLAN field the table loses
changes nothing observable** — the projection never consults plan entries,
so the field reaches the projected width by falling out of the map. That is
"defaulting into either width" (REQ-69) and the "bare omit-list whose
complement is whatever else exists" (REQ-60), reached through the carrier
rather than through the projection.

Mutation M-A, run to establish this: rename one PLAN key in
`resolvePayloadGroups` (`"Escaped"` → `"EscapedXX"`), leaving
`resolvePayload` and every JSON tag untouched. Result — **the entire shipped
suite stays green**, C2 completeness oracle included; only ADV-1 turns red.
The same shape arises naturally from a Go-identifier rename, where the JSON
tag survives so every wire-keyed oracle still sees the same key set.

For contrast, two mutations the shipped suite *does* catch, recorded so the
gap is not overstated: deleting a table entry outright (caught, by the
carrier probe's source-text grep for the key literal) and flipping
`"Outcome"` to `groupPlan` (caught end-to-end, six oracles red).

### ADV-2 — `0023:F2` via C1's reader clause: the mandated oracle was never written

**Test**: `TestAdversarialF2_ObservedReaderExecutionIsEqualPlusOrMinusTheFlag`
**Status**: PASSES; recorded as a **regression guard**, not a caught defect.

C1 fixes the measurand as OBSERVED READER EXECUTION and rejects the
recomputed form by name. The build carries the seam it mandates —
`internal/cli/flow_exec.go::readerExecutionHook`, correctly positioned at
the per-reader `exec.Read` call inside `runReaders`, correctly flag-blind —
but **no shipped test installs it**: `rg readerExecutionHook internal/`
returns the declaration, the call site, and one *string literal* in
`flow_partition_0023_test.go:743` where a source-text probe checks the
seam's position. The seam is dead code in the shipped build.

The Phase 1 reader oracle asserts a proxy instead: ±-flag equality of
`rule`/`next`/`writes`/`clear`. The proxy only discriminates a skipped
reader pass when the owned state that pass establishes changes the
selection; a skip on a request whose rule does not turn on that reader's
keys leaves every plan-group field identical, and `readers`/`owned` — the
two fields the pass yields together — are exactly what the flag projects
away. That is REQ-37's breach with nothing left to see it.

Verified to discriminate by mutation M-B (a flag-conditional skip of the
reader pass in `runFlowResolve`): the `reader-backed/owned-decides` arm goes
red on the observed set. The `decision-table/owned-does-not-decide` arm
carries no reader by construction and is the non-vacuity control, not a
discriminator — recorded as such rather than dressed up.

The implementation is right on this axis. What the fixup phase owes is the
oracle's existence, not a code change.

### ADV-3 — `0023:F3`: the refusal differential stops short of the site boundary

**Tests**: `TestAdversarialF3_RefusalsAtTheSiteBoundaryAndInTheExit3GroupAreFlagBlind`,
`TestAdversarialF3_NoRefusalReturnFollowsTheProjectionSite`
**Status**: both PASS; recorded as **regression guards**.

`TestReq33And100And131_RefusalEnvelopesAreByteIdenticalPlusOrMinusTheFlag`
runs three refusals — `flow-tag-invalid`, `flow-unmodeled-outcome`,
`flow-ambiguous-match`. All three are exit-2 and all three return from
`runFlowResolve` well above the gate site. Two classes are uncovered, and
both are what C1's site rule is written about:

1. **`flow-gate-denied` / `flow-gate-indeterminate`** — the last
   `respond.Fail` returns in the function, immediately preceding the
   projection branch. REQ-48 places the projection "AFTER the last
   respond.Fail return", so this is the exact boundary the clause draws, and
   the single refusal band a site that drifted one statement upward would
   newly reach. A projection applied above the three shipped arms would
   still be above their own returns, so they stay green through that drift.

2. **The exit-3 group.** REQ-32 asserts "identical exit codes" and the
   shipped arms exercise exit 2 only, pinning the assertion on one group. An
   environment failure refuses from inside the reader pass — the one pass
   REQ-37 names as the reachable surface for a report-width shortcut — where
   a flag-conditional skip surfaces as a changed exit, not as changed bytes.

Verified to discriminate by mutation M-B: all three arms turn red, the
exit-3 arm on the exit-code assertion the shipped exit-2-only differential
cannot make.

The companion test asserts the property the behavioural arms cannot express
— that no `respond.Fail` return follows the projection call in
`runFlowResolve`, and that no defensive "if refusing, skip projection" guard
exists (REQ-50). Today's refusals being flag-blind is a fact about the
fixtures; this is the structural reading `0023:A6` actually rests on.

### Axes attacked and found sound

Recorded so the absence of a finding is a result rather than a silence.

- **Absent-vs-null on the pointer conversion.** `--plan-only=false` is
  byte-identical to no flag; `--plan-only` twice equals once; the five echo
  pointers are allocated unconditionally at the one construction site, and
  every writer (`observedTagMap`, `tagMap`, `readerIDs`) is branch-free and
  `make`-allocating, so no nil arm is reachable on a success payload.
- **Aliasing / shared backing arrays.** `projectResolvePayload` takes and
  returns a value; the three container writers each allocate fresh, so
  nilling on the copy cannot reach caller state. The projected run leaves
  the artifact byte-identical.
- **Key order stability.** Projected wire order is default order minus the
  deleted keys (`revision` first), preserved by the declaration-order
  deletion loop; text mode sorts in the shipped `flatten`.
- **Echo-side table drift.** A Go-identifier rename of an ECHO field
  (`Observed` → `ObservedTags`, JSON tag unchanged) is caught by seven
  shipped oracles. Only the PLAN side is blind, which is ADV-1.
- **Plan-side field addition.** Adding an undeclared field to
  `resolvePayload` is caught, largely by 0010's pinned count and key list.
