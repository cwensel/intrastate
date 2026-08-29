Model: claude-opus-5[1m]

# Charted to successor — 3amigo, RDR 0024

Findings that are real but net-new scope for this record. Recorded here and
dismissed from this loop; none was folded into the draft.

- **PM-3 (adoption leg) — schedule an adopter for the declared vocabulary.**
  The motivating consumer's models (`rdr-status.toml`, `rdr-write.toml`, 54
  `[rule.emit]` blocks) live in the sibling `rdr` engine repo, not in
  `intrastate`. Declaring them is real work against a different repo on a
  different model and cannot be a phase here. This record now states the
  merge-day scope honestly (Consequences) rather than implying coverage it does
  not ship. **Successor**: a kata or RDR in the consumer repo that declares
  those models' emit vocabulary and adds the cross-file seam test below.

- **PM-7 / domain-drift and disposition-token-drift — the consumer seam test.**
  "Every declared emit domain member is a real command or stop token" is the
  outcome-closing check, and it is twice named in Risks as out of intrastate's
  scope by design (generic core). It cannot be written here: intrastate has no
  command surface to check against. It is currently unscheduled, which is the
  live edge PM-7 correctly identifies. **Successor**: same consumer-repo
  successor as PM-3 — the seam test is what converts this RDR's declaration
  into the consumer's guarantee.

- **Q-8 (second leg) — `lint.go`'s help text names a code it does not emit.**
  `internal/cli/lint.go` describes the load-refusal branch as
  `codeModelInvalid` (`flow-model-invalid`) while that branch emits
  `model-invalid`. Pre-existing repo drift, not introduced or worsened by this
  RDR; C2 now names both surfaces so no implementer "fixes" the emitted value
  to match the prose. **Successor**: a kata against the help text.

- **Multi-defect load reporting.** Raised in effect by PM-4 (N-round adoption
  loop). Accumulating findings at the load tier requires inventing the total
  order `internal/table/load.go` deliberately withholds — a pipeline-wide
  change the draft already scopes out by name. Unchanged by this pass; the
  adoption cost is now priced in Consequences instead.
