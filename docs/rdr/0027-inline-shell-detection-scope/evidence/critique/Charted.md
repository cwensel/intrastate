# Charted — critique lens, 0027

Net-new scope surfaced by the dual-model critique. Not absorbed into 0027; recorded
here so the disposition is durable and the driver can seed a successor.

- **stdin-appetite successor has no tracking artifact** (`0027:§capability-dependencies`,
  `0027:A5`) — the `stdin = "none" | "envelope"` axis is cited as the owner of the
  admitted stdin-fed forms, but exists only as a line in 0025's `evidence/critique/Charted.md`.
  A dependency with no kata cannot be cited as an owner. Out of scope here: 0027 decides a
  predicate, not the stdin axis. **Suggested successor**: open the kata, then link it from
  0027:A5; or restate C1's out-of-scope routing as an unowned admission.
  Sharpened by critique: the successor as charted would not close the axis anyway —
  `cmdbind.go:207` sets `cmd.Stdin` unconditionally, so it must first introduce a
  withholding point. That is design work the successor's brief does not yet name.

- **`shellInterpreters` has no golden test and no version marker** (`0027:C1`
  `interpreter set:` line) — the map is the OPEN deny-list C1 declares itself the
  amendment path for, but nothing pins its contents; `0027:S6` covers `Categories()`
  order only. Both passes converged here from opposite directions (blast-radius growth
  vs. amendment friction). Out of scope here: adding a test is an implementation
  obligation, and 0027's S-set is already written. **Suggested successor**: a golden
  test over the map naming `0027:C1`, plus a stated amendment path in the clause.

- **`report:` detail-string format is asserted nowhere** (`0027:§normative-contracts`)
  — the format is stable in implementation (`load.go:1180`, `argv[i] + " " + argv[j]`)
  but no contract pins it, while `Categories()` order IS golden-pinned. B-only finding;
  the asymmetry is real. Cheap to close by tightening `0027:S5`'s Expected from
  containment to format — but S5's wording is Stage 6's to amend, not this lens's.
