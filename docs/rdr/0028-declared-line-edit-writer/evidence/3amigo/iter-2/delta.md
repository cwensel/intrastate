Model: claude-opus-5

# 3amigo delta pass (iter-2)

Delta-scoped re-run over the passages the iter-1 fixes touched. This is NOT a
fresh full review: every judgement below is either "a fix opened a new gap" or
"an iter-1 finding is still open". Where a passage merely inherits a
pre-existing wording from an out-of-scope, Verified assumption, it is called out
as inherited rather than counted as fix-introduced.

**Two findings: one major, one minor.** The fix pass is otherwise clean —
the code citations all verify, the five new/rewritten scenarios carry real
criteria, the C1.2 brace-asymmetry is consistent everywhere it is assumed, and
scenario/step numbering is intact.

## Findings

### The `disposition` table is now incomplete — C1.3's new unreadable-target arm has no row

**Anchor:** `0028:§pre-lock-mini-checks` (`disposition` table), `0028:C1` (C1.3
`input:`), `0028:§testing-strategy`.

The iter-1 fix for IMPL-5 added a new refusal arm to C1.3 `input:`:

> A target that cannot be read refuses BEFORE mutation with the OS error in the
> Detail — it does NOT reuse `edit_anchor_unmatched` …

That arm was added to the clause and nowhere else. The `disposition` table is
captioned **"every input class, and what it mints"** and enumerates fourteen
rows including every other pre-mutation refusal (`edit_clear_undeclared`,
`edit_anchor_unmatched`, `edit_anchor_ambiguous`, `edit_anchor_collision`,
`edit_anchor_unstable`, `edit_value_multiline`, the gate arm, the unbound-tag
arm). "Target unreadable / missing / EISDIR / EACCES" is absent, so the table's
own completeness claim is false as of the fix.

The gap is not cosmetic, because this arm is the ONE refusal in C1 that carries
no reason TOKEN. Every other row names a token in the `Error / category` column;
this arm carries "the OS error in the Detail". That makes it the single row that
would have told an implementer what the Detail looks like when there is no
`edit_*` token to put in it, and it is exactly the row missing. It also leaves
the arm untested: a record-wide scan finds the string "cannot be read" at one
line only (C1.3 `input:`) — no scenario in S1..S31 exercises an unreadable
target, and `§failure-modes` F1..F8 does not name it either. The table's closing
sentence — "The two silent arms are the two the contract states as success;
there is no third" — is still true, but the enumeration it summarizes no longer
covers the contract.

**Blocks:** the implementer's Detail shape for the unreadable-target refusal,
and its test. Add the row (input class "target unreadable / absent", outcome
"refused pre-mutation", error `execution_failure` with the OS error in Detail,
artifact none, loud) and either extend `S27`'s error-surface scenario to cover
it or add a scenario; F1/F2's visible list should name it too.

### `§performance-expectations` cites `regexp.Expand` as the empty-group authority that C1.2's fix just disowned

**Anchor:** `0028:§performance-expectations` (Determinism, *Empty / null /
absent*), `0028:C1` (C1.2 `escapes:`).

The iter-1 fix for IMPL-10 rewrote C1.2 `escapes:` to draw a hard line:

> the empty-group rule MATCHES Go `regexp.Expand`'s, but the parse is this
> clause's own and never a call to `Expand`, which re-scans its template at
> expansion time and would break `parse once:`

`§performance-expectations` still reads "a non-participating capture group
expands to the empty string (Go `regexp.Expand`)" — a bare parenthetical
citation with none of the disowning qualifier. Read alone it names `Expand` as
the implementation, which C1.2 now forbids and which S14's bare-`$1` arm
falsifies (`Expand` accepts `$1`; C1.2 refuses it `edit_template_invalid`).

Minor because C1.2 is the normative owner and unambiguously wins; but the
parenthetical is now the only place in the record that points an implementer at
the wrong function, and the record has no other `Expand` citation to correct it.

**Blocks:** nothing normative. Retag the parenthetical as "same rule as Go
`regexp.Expand`'s, not a call to it (C1.2)".

## Verified, not merely restated

Every code citation the fixes rest on was checked against source in
`/Users/cwensel/sandbox/newcoinc/intrastate`. All five hold, and two hold more
precisely than the record needed:

- `internal/cli/flow_input.go::parseTags` — refuses `table.RecognizedTagKey`
  under `flow-tag-reserved` and an owned key under `flow-tag-owned`, both BEFORE
  any accessor. C1.2's "Provenance is NARROW because binding is" is exact. Note
  the third arm `parseTags` deliberately does NOT refuse — an UNDECLARED
  observed key binds fine (the guarded `decl := m.Tags[key]` lookup, commented as
  intentional). C1.2 handles this correctly by requiring the key be "a tag key
  the model declares with `observed` provenance" and making a non-declared key a
  LINT defect (`edit_anchor_invalid`), so the runtime laxity is fenced at load.
  No gap.
- `internal/table/load.go::carrierDefect` — clause 1 is literally
  `case a.Path != nil && a.Command != nil:` / `case !hasPath && !hasCommand:`
  with `hasPath := a.Path != nil && *a.Path != ""`. C1.1's prescribed edit
  (`&& !hasEdit` on the second arm) and its pointer-vs-emptiness discipline
  argument for `edit_carrier_conflict` both land on the real code.
- `internal/accessor/executor.go` deadline arm — `applyCtx, cancelApply :=
  context.WithTimeout(...)`, `err := binding.Apply(...)`, then
  `appliedDeadline := errors.Is(applyCtx.Err(), context.DeadlineExceeded)` and
  `r.applied = true`, with `err` never consulted on that branch. C1.3 `order:`'s
  scoping caveat and A11 describe this arm exactly.
- `flowbind.go::Registry` — `internal/cli/flowbind/registry.go:34`,
  `func Registry(m *table.Model, baseDir string, allowCommands bool)`, building
  `cmdbind.Config{BaseDir, AllowCommands}` and threading it to all three
  capability loops. Sole production construction site, as A10 and C1.3
  `read-back:` claim. `commandBacked(acc) = len(acc.Command) != 0 || acc.Path == ""`
  — S30's quoted form is byte-accurate, including the `acc.Path == ""` residue arm.
- `internal/cli/flow_exec.go::artifactMap` — `out[role] = accessor.Artifact{Role:
  role, Path: path}` built from `r.artifacts` alone, never touching
  `flowRequest.observed []resolveTag` (which exists, `flow_exec.go:27`). A
  repo-wide non-test grep for `accessor.Artifact{` returns this ONE site, and it
  is reached from `runReaders`, `runGates` and `flow_state.go:364` — so A1's
  "sole non-test construction site … shared by `resolve`, `set-state` and
  `read-state`" is exactly right.
- MVV step 2's load-bearing `--tag` claim — `registerTagFlag(cmd)` appears at
  `flow_resolve.go:120` and `flow_state.go:199` (set-state), confirming both
  halves register it. `flow_state.go`'s read-state command (line 49ff) registers
  only `registerSelectionFlags`, so it does NOT take `--tag`; A1 never claims it
  does, only that it shares `artifactMap()`. Accurate as written.
- `§illustrative-code` — executed, not eyeballed. The block parses under
  `tomllib` (the `;` separator and integer `timeout` from IMPL-8 are gone; the
  comment now correctly warns `Timeout` is `*string`). Both anchors compile
  under Go RE2: the Status anchor yields one named group `q`, and
  `ExpandString` with `${q}` round-trips `Draft [joint decision → …]` →
  `Final [joint decision → …]` AND collapses to `- **Status**: Final` on a
  plain-`Draft` line (the non-participating-group rule C1.2 states). The README
  anchor, with `{tag.nnnn}` replaced by a `QuoteMeta` probe, yields two groups
  and `${1}Final${2}` rewrites only the status cell. The illustrative block is
  now self-consistent with C1.2's `${name}` admission.
- `S12`'s fixture — `rdr/tools/rdr/testdata/status/records/0022-cache-metrics-surface.md`
  exists and its Status line reads `- **Status**: Draft [revised from Final
  2026-08-10; re-verify A1,A2 — the` wrapping to a continuation line. Both
  properties S12 requires (wrapped AND `Draft`) are present, and the record it
  displaces (`0021-cache-warmup-order.md`) is the one A3 used. QA-4 is closed on
  a real file.

## New scenarios S28–S31 — each earns its clause

All four carry a falsifiable criterion and each genuinely reaches a clause no
other scenario could:

- **S28 / C1.4 `precedence:`** — a multi-defect entry is the only input that can
  observe fail-fast order, and the scenario says so in its own justification
  ("with one defect per fixture every evaluation order passes"). Expected names
  the winning category AND "only it", so a implementation reporting both fails.
  Correctly pairs `edit_carrier_conflict` (1st) against an uncompilable anchor
  (3rd) and `edit_key_mismatch` (2nd) against a malformed `replace` (4th) —
  both orderings match C1.4's registration list.
- **S29 / C1.5 `one-way:`** — genuinely two-invocation, which no other scenario
  is; S17 covers `clear` within one invocation and cannot reach the
  irreversibility. Expected is two distinct assertions (key ABSENT, then
  `edit_anchor_unmatched`), and the second is derived correctly from C1.3
  `select:`'s no-append rule.
- **S30 / C1.1 `in-memory:`** — reaches the residue arm at RUNTIME, which no
  load-time scenario can (C1.4 makes it unreachable through the loader, as
  `registry.go`'s own comment says). The stated risk is real and specific: the
  discriminator's `acc.Path == ""` arm IS the residue path, so a three-way
  rewrite could silently reroute residue to a `Path: ""` file binding that fails
  open. The failure it would produce ("every declared key reads absent … would
  confirm an unapplied write") is the exact hazard `registry.go` documents.
  Pass/fail is structural rather than token-shaped, which is the weakest
  criterion of the four, but the failing behaviour is observable (keys read
  absent instead of a refusal) so it is not un-testable.
- **S31 / F4+F5 `read_back_mismatch`** — the only scenario asserting a
  post-mutation failure with `Applied()` TRUE, and it says why no pre-mutation
  scenario can reach it. Consistent with A8/JDR 0003 §D1's rule 3 (applied sense
  rides `Applied()`) and with F4/F5's text.

`S20`, `S22` and `S23` also now carry the witnesses iter-1 said they lacked:
S20's unchanged INODE (with the correct argument that content cannot separate
"did not write" from "wrote identical bytes", and A5's confirmation that
stage-and-rename always replaces the inode), S22's explicit demotion of the
`os.Rename` behaviour to "stated as the reason, NOT as an assertion this
scenario runs", and S23's gate-OFF file-backed-reader run as the separator
between "not consulted" and "consulted and permitted" — with the S23/S25
contrast made explicit in both directions. S8's fix is likewise sound: the
in-scope predicate ("carries a `- **Status**:` bullet, which excludes
`*-postmortem.md` and `BUILD-ORDER.md`") is now part of the assertion, and the
cardinalities are explicitly demoted to descriptive with the reason given
(a corpus that grows with every seeded RDR).

## C1.2 brace asymmetry — no passage still assumes symmetry

A record-wide scan for `{{`, `}}`, "escape", "Expand", "symmetric" and "both
templates" returns four sites. C1.2 `vocabulary:` states the asymmetry
explicitly and names each side's rule. S14 tests both sides and closes with "The
two templates are NOT symmetric, which is the point of the scenario". C1.2
`parse once:` says "both templates are parsed at LOAD into segments", which is a
statement about the parse PHASE, not the vocabulary, and does not conflict. The
`§decision-rationale` mention is a premortem ledger line ("escapes … folded into
C1.2/C1.3") that asserts nothing about vocabulary. The only residue is the
`§performance-expectations` parenthetical filed as the minor finding above.
`edit_anchor_invalid`'s predicate in C1.4 was also updated in step ("carries any
other `{…}` form" is the LINT category text; C1.2's narrower "never for a brace
RE2 itself accepts" governs) — these two want a second look at lock, but C1.2 is
the stated owner and C1.4's category line is a summary, so it is not a
contradiction as written.

## Numbering and reference integrity — intact

- Scenarios are contiguous `1..31` with no gap and no duplicate (checked by
  extracting every `N. **Scenario**` heading in order).
- No reference anywhere in the record points above S31.
- The MVV has seven numbered steps. Every step reference resolves: "MVV 3" and
  "MVV 4" in the `fidelity`/`disposition` tables, "MVV step 5" in S20 and S25,
  "MVV step 7" in `§risks-and-mitigations`, and "steps 5–7" in both the
  mini-check preamble and `§testing-strategy`'s Acceptance line. All seven trace
  rows (1, 2, 3a–3f, 4, 5, 6, 7) map onto real MVV steps.
- S20's "No MVV step covers this — step 5's files are untouched because the
  write REFUSED" agrees with MVV step 5 and with the `disposition` table's
  separation of the no-op row from the gate-refusal row. The iter-1 wrong
  cross-references (QA-3) are gone.
- The `trace` table's step 2 row does not restate MVV step 2's new "bound on
  BOTH halves" emphasis, but it does not contradict it either — its assertion
  ("`--tag nnnn=NNNN` bound on the invocation") and witness (`registerTagFlag`,
  `parseTags`) are both true of each half. Not filed as a finding; a one-clause
  echo would strengthen it.
- The `disposition` table's gate row ("command-backed reader, gate off →
  refused pre-mutation … never `read_back_incomplete`") agrees with C1.3
  `read-back:`, S25 and trace step 5. The `fidelity` table's S20 and S21/S22
  rows agree with the rewritten scenarios.

## Iter-1 findings: all closed except as noted

Checked each of the twenty-two consolidated findings against the current text.
Twenty are closed. #11 (unreadable target) is **partially** closed — the clause
was added, the disposition row was not, which is the major finding above. #21
(`regexp.Expand` citation) is closed in C1.2 and left open in
`§performance-expectations`, which is the minor finding above. The two that were
hardest to verify — #1 (gate pre-check site) and #3 (two-hop context plumb) —
are closed with source-accurate SITE clauses: C1.3 `read-back:` now names
`def.Accessor.Command` and the `accessor.Registry` field set at
`flowbind.go::Registry`, states the no-cycle argument correctly (`cmdbind`
imports `accessor`), and A10 keeps Pending status with an honest "to check" list
rather than overclaiming. #14 (C1.6 unscored) is closed by
`§decision-rationale`'s dedicated paragraph, which weighs the per-record-role
alternative A1 named and gives the rejection reason on the matrix's own
drift-removal row. #7 ("zero caller edits") is closed by
`§problem-statement`'s new "Both words are narrow" paragraph, which names the
consumer's remaining surface (`--tag nnnn=NNNN`, the projector verb) and points
at Capability Dependencies — where both rows are in fact present and correctly
marked. #15 is closed by `§prerequisites`' new consumer-side bullet, which
scopes the two unlinked rows as adoption-blocking rather than
implementation-blocking and states the partial-adoption consequence plainly.

## Checked and CLEAN

`0028:C1` C1.1 (`carrier:`, `in-memory:`), `0028:C1` C1.2 (`anchor admits:`,
`escapes:`, `vocabulary:`), `0028:C1` C1.3 (`target:`, `select:`, `order:`,
`write:`, `read-back:`), `0028:C1` C1.4, `0028:C1` C1.5, `0028:C1` C1.6,
`0028:MVV`, `0028:A1`, `0028:A10`, `0028:A11`, `0028:S8`, `0028:S12`,
`0028:S14`, `0028:S20`, `0028:S22`, `0028:S23`, `0028:S28`, `0028:S29`,
`0028:S30`, `0028:S31`, `0028:§problem-statement`, `0028:§decision-rationale`,
`0028:§illustrative-code`, `0028:§prerequisites`, and the `fidelity` and `trace`
tables in `0028:§pre-lock-mini-checks`.

Not clean: `0028:C1` C1.3 `input:` ⨯ the `disposition` table in
`0028:§pre-lock-mini-checks` (major, above), and `0028:§performance-expectations`
(minor, above — outside the nominated scope but reached by the C1.2 fix).
