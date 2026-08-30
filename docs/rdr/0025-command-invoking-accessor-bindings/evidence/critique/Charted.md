# Charted to successor — critique lens, RDR 0025

One line per finding that was real but net-new scope for this record.

- **Declared stdin appetite (`stdin = "none" | "envelope"`) on a command
  entry** (diff D-13; pass A C-13, pass B C-6). The spike-proven `tee
  {artifact}` hazard destroys the user's artifact and is not statically
  detectable from argv, so no C5 arm can catch it. Withholding the envelope
  from a write command that never declared it would make the footgun
  unauthorable. Out of scope here: it changes the carrier shape C1 fixes, and
  C1's field set is this record's locked surface. Suggested successor: a
  carrier-shape RDR that adds the field and the C5 arm together.

- **User-scope configuration for `allow_commands`** (diff D-11; pass B C-4).
  Every peer in this RDR's own reversal ledger ended at a persistent gate; the
  per-invocation flag is v1's shape only because no config surface exists
  (A10, Verified). Out of scope here: a config file means building discovery,
  precedence and malformed-file handling from nothing. Already charted by the
  3amigo lens (`evidence/3amigo/Charted:`) — recorded again here because the
  critique reached it independently, which raises its priority rather than
  duplicating it. Suggested successor: the config-surface RDR, which composes
  disjunctively with the flag (either grants, neither revokes).
