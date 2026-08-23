Model: claude-opus-5[1m]

# Charted to successors — critique iter-2

- **M-6 — `RequiresOwned` excludes guard-read keys, degrading the refusal
  diagnosis.** Verified on `main`: `internal/resolve/resolve.go::gate` evaluates
  guards first, then `missingOwned`, then undecidable. A guard reading an owned
  tag absent from the view returns `GuardUnevaluable`, and because
  `RequiresOwned` names only write/clear keys, `::missingOwned` finds nothing —
  so the author sees `guard_unevaluable` where `owned_state_unavailable` is the
  more precise and actionable diagnosis. Out of scope here: this RDR's text
  already allocates the field's *meaning* to RDR 0007 (post-guard write
  dependencies, JDR 0001 §JD-3) and states it "does not add guard-read keys to
  it." Suggested successor: **RDR 0007** — either widen `RequiresOwned` to
  include guard-read keys, or accept the coarser diagnosis explicitly.

- **M-12 — ~17 stable load-time validation categories have no CLI mapping and no
  Final owner.** This RDR fixes the data-level category set; mapping them to
  `CLIError` codes and exit statuses is RDR 0005's, which is `Pending`. Until it
  lands, every category surfaces as a generic CLI failure and the stability this
  RDR guarantees is invisible to scripts. Suggested successor: **RDR 0005** —
  carry the category set into the CLI output contract as distinct codes.

- **M-5 residual — empty `Input.Recognized` at the caller boundary.** This RDR
  bans the empty string from the alphabet, which breaks the two-conjunct hazard
  at its own boundary, and the A8 record now states the barrier's bounded reach.
  The remaining guard — rejecting an empty `--recognized` before calling
  `Resolve`, and a kernel-side alphabet check covering non-normalizer producers —
  belongs to **RDR 0005** (CLI edge) and **RDR 0001** (kernel). Recorded so the
  residual is not mistaken for table-wide coverage.

- **M-16 — the dump is sold as the review surface but has no grammar.** Real
  tension, correctly scoped out here: this RDR fixes the dump's field list and
  row order and explicitly declines to define a re-readable grammar (three named
  lossy sites in Round-Trip). A re-readable dump format is net-new scope.
  Suggested successor: a **dump-format RDR**, which must close the `<clear>`
  sentinel, separator escaping, and suffix separator (this pass reserved `#` in
  rule ids and alphabet members, which removes the third site's identity-level
  collision but not its rendering caveat).
