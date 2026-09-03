# Deviations — RDR 0027 inline-shell-detection-scope

Entries pre-seeded by the 7.1 cluster reconcile (0026-0027-0028, 2026-09-03)
are left OPEN for Stage 8: running the named check is the disposition; an
entry escalates only if the check contradicts a contract.

## D1 — S6 pins `Categories()` order relatively, not as a tail golden

**Type**: TEST-FIXTURE. **Status**: OPEN (pre-seeded, 7.1 pairwise 0027×0028 PW1).

**Situation.** `0028:C1.4` appends five categories after 0025:C5's six.
Shipped `internal/table/command_carrier_0025_test.go::TestReq77_…` asserts
the six as the TAIL slice of `Categories()` and breaks on that append.
`0027:S6`'s "membership and order unchanged" is this record's own delta claim
and holds in either landing order only if authored as relative order (the
`command_shell_interpreter` constant's position relative to its neighbours),
never `len`- or tail-equality — 0025:C5 makes the list's size a non-contract.

**Check.** S6 green against a `Categories()` with five extra members appended
at the tail.

## D2 — Is the Phase 3 description surface total over `Categories()`?

**Type**: IMPL-DECISION. **Status**: OPEN (pre-seeded, 7.1 pairwise 0027×0028 PW2).

**Situation.** Phase 3 adds an accessor pairing a category with its text,
rendered on `--help-all` and mirrored into `docs/cli-reference.md`. If that
accessor or the docs staleness gate asserts totality over `Categories()`,
0028's five `edit_*` categories fail it or ship undescribed, and neither
record names the obligation.

**Check.** Decide total vs per-category opt-in and record it here; the
accessor's test states which. If total, 0028's deviations D2 owes the five
texts. `make check`'s docs gate green after 0028's append.

## D3 — Merged clause-3 behaviour of `carrierDefect` after both land

**Type**: TEST-FIXTURE. **Status**: OPEN (pre-seeded, 7.1 pairwise 0027×0028 PW4; critique C-8).

**Situation.** 0028 turns clause 3's placeholder membership test into a
declared-tag family check (a signature change) and both records rewrite the
clause-3 exemption comment and the function's doc comment; `0027:A1`'s line
anchors shift. The JCs disposed the function as disjoint clauses; clause 3 is
touched by both.

**Check** (whichever lands second): `["nice","sh","-c","cat {tag.x}"]` with
`x` declared, and `["nice","sh","-c","cat {nope}"]`, both report
`command_shell_interpreter`; the clause-3 exemption comment no longer reads
"argv0".
