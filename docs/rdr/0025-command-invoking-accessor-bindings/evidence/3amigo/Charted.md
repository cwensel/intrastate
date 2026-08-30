Model: claude-opus-5[1m]

# Charted to successor — 3amigo (0025)

- **User-scope configuration surface for `allow_commands`** (from the A10/C6
  hotspot, all three personas). Real, but net-new scope: the repo has no
  `internal/cli/config/`, no `intrastate.toml` reader and no user-scope
  configuration of any kind, so building one means deciding a search order, a
  precedence rule against the flag, and an absent/malformed-file behaviour
  (fail-closed) — a subsystem, not a clause. 0025 v1 gates on
  `--allow-commands` alone (A10 Verified); the successor inherits an
  already-gated seam and owes only the file's discovery + precedence +
  fail-closed-on-malformed rules. Suggested successor: a config-surface RDR;
  worth an `/rdr-seed`.

- **Automatic resolve→apply wiring** (from the PM's MVV finding). Grounded:
  `internal/cli/flow_state.go` is the sole production caller of the executor's
  `Write` path and is driven by `--write` flags, so no code path takes a
  resolved decision and invokes the declared write. 0025 fixes the
  *declaration* half only; the caller still carries the decision across by
  hand. Recorded in the Problem Statement as surviving residue rather than
  absorbed. Suggested successor: a resolve→apply RDR.

- **Retirement of `flowbind`'s magic-path suffix vocabulary** (Selection LBD,
  already recorded there as a successor's decision — restated here only so the
  3amigo ledger closes on it rather than leaving it implicit).
