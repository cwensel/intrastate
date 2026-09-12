Model: claude-opus-5

# Charted to successor — 3amigo, cli/0029

Real findings that are net-new scope for this RDR. Recorded here, dismissed
from this loop.

- **Charted: expose the tier table from the binary** (PM S2, `0029:BR4`).
  `llms.txt` instructs agents to prefer the binary over any file, and this
  RDR's sole user-facing delivery is prose in a file. A machine-readable
  tier surface (`intrastate version --as=json` carrying tiers, or a
  `--help-all` section) would satisfy that precedence rule. Out of scope:
  `0029:BR4` weighed document-vs-document and this is a new emitted surface
  — which would itself need a tier, and would be governed by the very
  contract this RDR is writing. Successor RDR if the prose register proves
  unread. Activation Step 1 now names it as the successor rather than a
  fallback.

- **Charted: a release-note mechanism for C3's disclosure obligation**
  (PM S4). Through `0.x`, C3's disclosure is the only clause that binds, and
  it has no file, template, generator or check anywhere — `0029:F2` books
  the consequence as residual risk. Out of scope: the release *mechanism* is
  explicitly excluded by `0029:§problem-statement`, and this is process
  tooling, not a promise about output. Recorded in §capability-dependencies
  as an unbacked Introduced capability. Successor: a kata, or an RDR if the
  disclosure needs a generated artifact.

- **Charted: a conformance fixture for C1's consumer rules** (Impl F11).
  C1's tolerant-reader MUSTs bind an external consumer and now say so; a
  published fixture an external consumer could run against would make them
  checkable rather than prose. Out of scope: it is a new published artifact
  with its own compatibility surface. Successor RDR if a third-party
  consumer appears.

- **Charted: `docs/cli-output-contract.md` carries no versioning language
  at all** (QA F5, grounded: `frozen`/`append-only`/`growing`/`release`/
  `upgrade`/`breaking`/`compat` all zero). Activation Step 1 now states the
  four acceptance criteria that close this. What stays net-new is a TEST
  that the register agrees with C4 — a docs-vs-contract drift check. Out of
  scope: no such check exists for any doc in this repo, so it is new
  tooling, not an assertion this RDR can add. Successor: a kata against
  `make docs-check`.
