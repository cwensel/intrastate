# Charted to successor — critique lens, RDR 0019

- **In-band discovery pointer from `flow next` to `init-state`**
  (critique U-16, raised by pass B C-8). Real: BR4/BR5 close both
  automatic-invocation routes, so `init-state` is typed by hand, and
  `flow next`'s `unknown[].reason: absent` is where the operator meets
  the wall this RDR removes. Out of scope HERE because carrying the
  pointer in-band means adding a field or reason-string to `next`'s
  payload — a second override of `0005:C1`'s I/O, where this record's
  override is the verb enumeration alone. RDR 0019 discharges the
  obligation as documentation (Activation Step 1) and discloses the
  residual in Consequences.
  Suggested successor: an RDR extending `flow next`'s unknown-report
  payload with a remediation pointer, which would also serve the
  `[initial]`-key-added-after-seeding residual (F4).
