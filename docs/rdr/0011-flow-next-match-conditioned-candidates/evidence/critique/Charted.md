# Charted to successor — critique, RDR cli/0011

Findings that were real but net-new scope. Each is dismissed from this loop and
recorded here rather than absorbed into the RDR (rdr-common: never edit the
current RDR to absorb net-new scope).

## 1. `registerSelectionFlags` pins no negative for verb-scoped flags

**Origin**: A-13 (pass A). **Why out of scope**: C2 keeps `--all` off
`flow resolve` / `read-state` / `set-state` by NON-REGISTRATION, asserting on
cobra's shipped unknown-flag error. A-13 is right that a negative held only by
"nobody added it to the shared registrar" is not a contract — but pinning
verb-scoped flag surfaces is a CLI-framework convention spanning every verb,
not this record's contract. RDR 0011 already routes the assertion to a real
error class and exit code, which is this record's obligation.

**Suggested successor**: a kata or small RDR on flag-surface negatives — whether
`registerSelectionFlags` should carry an explicit deny-list or a test that
enumerates each verb's accepted flag set, so a future symmetry refactor cannot
silently widen a selection site.

## 2. Zero-owned decision tables are unloadable, so 0010's coupling is verified on a proxy

**Origin**: A-14 (pass A), corroborating a limit A5 already records. **Why out of
scope**: A5 states the limit honestly — a genuine zero-owned 0010 decision table
is refused `malformed_rule_shape` / `write_to_non_owned_tag` on this build, so
the fixture carries a dummy owned tag. That the real shape cannot yet exist is a
property of what 0010 has shipped, not a defect in 0011's contract, and
`0010:A11` already pre-commits to re-verifying against 0011's predicate once
0011 locks.

**Suggested successor**: carry into RDR 0010's own validation — when the
zero-owned table becomes loadable, 0010 owes the `flow next` oracle over it.
Worth naming in 7.1 cluster-reconcile rather than seeding separately.
