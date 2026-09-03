# Deviations — RDR 0028 declared-line-edit-writer

Entries pre-seeded by the 7.1 cluster reconcile (0026-0027-0028, 2026-09-03)
are left OPEN for Stage 8: running the named check is the disposition; an
entry escalates only if the check contradicts a contract.

## D1 — `TestReq77`'s tail-slice assertion breaks on C1.4's append

**Type**: TEST-FIXTURE. **Status**: OPEN (pre-seeded, 7.1 pairwise 0027×0028 PW1).

**Situation.** `internal/table/command_carrier_0025_test.go::TestReq77_…`
asserts 0025:C5's six categories as the TAIL of `Categories()`
(`all[len(all)-len(want):]`). `C1.4` appends five after them, so the shipped
test goes red with no behaviour change. C1.4's "the list's size is not a
contract (0025:C5)" holds; the test over-asserts.

**Check.** Rewrite `TestReq77` to relative order — the six contiguous and in
clause order, followed by the five — and run it with this record's S6:
`go test ./internal/table -run 'TestReq77'` green.

## D2 — Description text for the five `edit_*` categories

**Type**: SPEC-UNDER. **Status**: OPEN (pre-seeded, 7.1 pairwise 0027×0028 PW2).

**Situation.** 0027 Phase 3 ships a per-category description surface. C1.4
names wire strings, constants and order only. If 0027's surface is total over
`Categories()` (0027's deviations D2 decides), the five need text on it.

**Check.** After 0027 lands: its accessor's totality test and `make check`'s
docs gate green with the five appended. If 0027 chose per-category opt-in,
close as no-op.

## D3 — Merged clause-3 behaviour of `carrierDefect` after both land

**Type**: TEST-FIXTURE. **Status**: OPEN (pre-seeded, 7.1 pairwise 0027×0028 PW4; critique C-8).

**Situation.** C1.6's declared-tag family check changes clause 3's membership
test and `carrierDefect`'s signature; 0027 rewrites the adjacent exemption
comment and clause 4's predicate. Same check as 0027's deviations D3, owed by
whichever lands second.

**Check.** `["nice","sh","-c","cat {tag.x}"]` with `x` declared, and
`["nice","sh","-c","cat {nope}"]`, both report `command_shell_interpreter`
after both records land, in either order.
