# Charted — cli/0023 critique lens

Net-new scope surfaced by the critique lens, recorded here rather than
absorbed into the RDR (rdr-common: never edit the current RDR to absorb
net-new scope).

- **Adoption / migration change (M-3, from `intrastate#srz2` P1).** This
  RDR ships the projection capability; the motivating token saving is
  realized only when a caller passes `--plan-only`. Out of scope because
  it is a consumer-side call-site change with a different blast radius
  from the CLI surface this RDR owns. Suggested successor: a migration
  change (kata, not an RDR — no design fork) that adds the flag at the
  navigator/status call site and records the realized reduction against
  A1's table. §Consequences now states that a green build unblocks
  rather than closes `srz2`.

- **Traffic-mix measurement (M-18).** The economic case assumes success
  calls dominate the chained-call pattern; the refusal/success mix of the
  motivating consumer is explicitly unsized and this RDR does not measure
  it. Out of scope because the flag is strictly non-negative on every
  call — the mix changes how much is saved, never whether the change is
  worth landing (the RDR's own decided text). Suggested successor: the
  same migration change measures the realized mix when it measures the
  realized reduction.

- **Always-keep oracle quantified over modes (M-12).** C2 already records
  this as an obligation on the successor that adds a second projection
  mode, and explains why minting it now would assert exactly what S2
  asserts. No new charting needed — noted here so the critique row is not
  read as unresolved.

- **`revision` activation (M-5 residual).** A projected plan is
  unattributable while `revision` is constant-empty and `model` is ECHO.
  Out of scope: `revision`'s emptiness is owned by RDR 0002's `[model]`
  block admitting no revision key, not by this RDR. Suggested successor:
  whichever RDR gives `revision` a live value closes the residual without
  moving a field between groups.
