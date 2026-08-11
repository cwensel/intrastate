Model: claude-opus-5[1m]
Lens: cove — origin ledger (iteration 1)

Reconciled by passage anchor across `findings-main.md` (dispatcher grounding),
`findings-b.md` (adversarial/implementer pass), and `findings-a.md` (claim sweep).
IDs are per-file and do not correspond across passes; the rows below are the merged
origin ledger the resolve pass dispositions against.

| # | Origin | Anchor (grepable snippet) | Type | Disposition |
| --- | --- | --- | --- | --- |
| L-1 | main F-1, B FB-1/FB-2 | "populates **`NextTags` only — it sets no `Writes` at all**" / "The `escapeRow` builder needs no change" | refuted-claim (blocking) | fixed |
| L-2 | main F-2, B FB-7 | "ordered by RowRef identity using the kernel's existing compareRefs ordering" | silence / unsatisfiable contract | fixed |
| L-3 | main F-3, B FB-4 | Infra Audit: "Breach is diagnosable from the error text (A7), not a new code" | internal contradiction | fixed |
| L-4 | main F-4, B FB-3 | "degrades to the single error's own message for one breach" | refuted-claim (precision) | fixed |
| L-5 | B FB-5 | Technical Design: "returns a nil `Result` and a non-nil error" | unexpressible type / contradiction with MVV | fixed |
| L-6 | B FB-6 | "MUST be exported by the kernel package as a construction-time check" | silence (receiver + sentinel export) | fixed |
| L-7 | B FB-8 | "`Resolve` already scans `Table.Rows` twice per call" | refuted-claim (minor, measured) | fixed |
| L-8 | B FB-9 | "`[]Tag{}` provably cannot produce the user-facing symptom" | overstated grounding | fixed |
| L-9 | main Step 0b, A F-4 | `respond.ValidateMode` is the in-repo `Validate*` precedent, uncited | sibling-exists (non-blocking) | dismissed-with-cite |
| L-10 | A F-3 | "the verb MUST wrap it into a *clierr.CLIError … the offending row identity" | silence (no wire carrier) | fixed |
| L-11 | A F-2 | "the shipped config.Load wrapping pattern" cited for `Hint` | refuted-claim | fixed |
| L-12 | A F-6 | "a frozen suite with **zero** nil-sensitive assertions" | refuted-claim (overbroad) | fixed |
| L-13 | A F-7 | "by its own boundary, leaves parse/render fidelity and row-shape validation with RDR 0002" | refuted-claim (misattributed) | fixed |
| L-14 | A F-8 | Failure Modes: `errors.As` presented as the way to read "the row identities" | silence (aggregate traversal) | fixed |
| L-15 | A F-5 (half) | `Result{nil,nil}` vs the frozen "never neither" invariant | refuted finding | dismissed-with-cite |

## L-2 resolution note (tiebreak reversed mid-pass)

The first fix for L-2 pinned `slices.SortStableFunc`, making table position the tiebreak
among equal-identity rows. Checking that against the frozen suite refuted it: RDR 0001's
REQ-2/REQ-10 forbid a diagnostic payload that varies with row order, and two frozen tests
enforce it directly —
`adversarial_test.go::TestAdv3b_MissingOwnedPayloadMustNotDependOnTableRowOrder` and
`fixup_test.go::TestFixup3c_AmbiguousMatchRowsPayloadMustNotDependOnTableRowOrder`. ADV-3b's
comment names the exact trap: "stable order — stable with respect to row order, which is not
the same as stable with respect to the input tuple." A positional tiebreak would have
reintroduced the defect those tests freeze.

Final disposition: equal identities **collapse** to one reported entry, with the breach count
stated alongside. This keeps the payload a function of the input tuple under any permutation.
Recorded as A8 (Verified, Source Search).

## Grounding gate results

All rows above passed the gate: each cites source read at HEAD (`20cc49f`), none
re-litigates an already-decided option, and none conflicts with `{RDR_RESOURCES}`
principles. Two rows were checked and **held** rather than becoming findings:

- A3's non-vacuity argument ("every `Plan.Writes` read on an escape-derived plan sits in
  a `t.Fatalf` format argument, never in a condition") — CONFIRMED at
  `adversarial_test.go:238-241`, `fixup_test.go:54-57`, `fixup_test.go:112-115`. So L-1 is a
  scope correction, **not** a re-opening of the A3 verdict: the 154-PASS conformed run was
  measured with the builder edit applied, and the argument for why stripping is safe covers
  the builder's rows for the same reason.
- The frozen suite's "zero nil-sensitive assertions and four length-based `Writes` checks"
  — CONFIRMED exactly (`mvv_test.go:110`, `resolve_test.go:668`, `:709`, `:986`).

## L-15 dismissal (pass A F-5, second half)

Pass A held that returning `Result{nil, nil}` on a breach violates `Result`'s frozen
"never neither" invariant, asserted at `resolve_test.go:80` ("disposition carries neither a
plan nor a refusal"), and that the RDR narrows it in neither Overrides nor text.

**Refuted by source**: that assertion is reached only through `resolve_test.go::mustResolve`,
which does `if err != nil { t.Fatalf("Resolve returned a Go error for a modeled disposition") }`
*before* returning the `Result`. A breach returns a non-nil error, so the invariant test
fatals on the error path and never observes the zero `Result`. The invariant is scoped to
nil-error dispositions — exactly the "modeled disposition" its own failure message names — so
the entry precondition does not narrow it and owes no Overrides entry.

The *other* half of F-5 (the unexpressible "nil `Result`" wording, and "on the first row"
contradicting the aggregate contract) was valid and is fixed under L-5.

## L-9 dismissal

`internal/cli/respond/respond.go::ValidateMode` is the repo's only existing `Validate*`
function. It returns `*clierr.CLIError`, not `error`, so it neither satisfies nor
contradicts the RDR's naming clause (which governs an `error`-returning kernel predicate and
is argued from `pprof.Profile.CheckValid` / `rsa.PrivateKey.Validate`). Recording it as a
sibling would add a citation without changing the contract. Dismissed as non-blocking;
noted here so a later pass does not re-chase it.
