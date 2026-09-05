# Charted to successor — cove lens, 0019

- **Finding 3 (extension)** — extending the empty-store predicate to the
  non-file-backed write carriers (`::NewEditWriter` / 0028 edit-carried,
  `cmdbind.Writer` / 0025 command-backed). Out of scope here: this RDR fixes
  `[initial]` materialization semantics on the file-backed carrier, and a
  carrier-independent emptiness notion is a different design question that
  0025/0028 own the shape of. 0019 refuses those carriers explicitly (C1
  carrier scope) rather than under-defining them. Suggested successor: a new
  RDR seeded once A6's carrier decision lands, since the emptiness surface it
  picks determines whether a carrier-independent form is even expressible.

- **Finding 8 (gate)** — the five Finalization Gate sub-sections are unwritten
  TEMPLATE guidance. Not a cove defect and not resolved here: the gate is
  Stage 7's own work (`/rdr-finalize`), which writes the responses and clears
  the five `placeholder:survived` lint findings. Recorded so the driver does
  not read the pre-lock close as gate-complete.

## Pass 2 (delta) — net-new anchors, dispositioned

The `comm -13` diff names A2, F5, `§testing-strategy` and A6 as net-new
against pass 1's ledger. None is scope expansion:

- `0019:A2`, `0019:F5` — propagation sites of pass-1 Finding 6 (the
  conform retraction in C1). Fixing them IS the amendment sweep, not new
  scope: C1 deleted a step two siblings still relied on.
- `0019:§testing-strategy` — the scenario (S10) owed by pass-1 Finding 3's
  fix, which added the carrier refusal. A new terminal arm owes an oracle.
- `0019:A6` — created by pass-1 Finding 2's fix; pass 2 corrected its
  implementer enumeration (two `ReadBinding` implementers, not three).
