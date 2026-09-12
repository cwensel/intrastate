# 3amigo consolidation — cli/0029

Personas: PM (persona-1-pm.md), Implementer (persona-2-implementer.md), QA
(persona-3-qa.md). Three isolated sub-agent passes, no cross-persona visibility;
isolation verified (no file references another's output). 27 findings total.

**Hotspots** (ids anchored by >=2 independent personas; overlap marks a hotspot
passage, not a validated finding — a single-persona finding is not thereby weaker):

    3  0029:MVV   0029:C4   0029:C2   0029:C1
    2  0029:S9 S8 S7 S6 S3 S2 S1 · 0029:D-naming · 0029:A4 A3 A2
       0029:§capability-dependencies

All four contracts and the MVV were reached independently by all three personas.

---

## Hotspot 1 — `0029:C4`'s tiered vocabularies have no enumeration seam (3 personas)

- **QA F1 (BLOCKING)** — 3 of 5 newly-`frozen` vocabularies have no enumerable
  accessor, so `0029:S6`'s by-value assertion is unwritable: exit-code classes
  (`clierr::ExitCodeFor` is a `switch`, enumerates nothing), the stderr advisory
  `level` set (bare literals in `respond.go`), and the `graph-lint-failed`
  aggregate code (a single `const`, not a set). Writable for 2 of 5
  (`Operators()`, `Verdicts()`). Blocks scoping Phase 1 Step 3 — three test
  files, or three test files plus three new production accessors?
- **QA F2 (BLOCKING)** — same defect on the `append-only` side and at larger
  scale: the CLIError `code` vocabulary has no enumeration anywhere (codes are
  minted as literals at raise sites), so `0029:S7`'s membership-and-uniqueness
  assertion has no subject and the append-only promise has no mechanical guard.
  Same for the `graph-unprovable-coverage` `reason` set and the `flow next`
  unknown-`reason` set.
- **Impl F4 (SEV-1)** — `findings[].operator` is anchored ambiguously: two
  `Operators()` exist (`internal/guard::Operators` and
  `internal/table/model.go:96`). C4 tiers one, is silent on the other; S6 has
  two candidate subjects.
- **Impl F5 (SEV-2)** — `exit-code classes` is not a vocabulary at C4's stated
  granularity: `ExitCodeFor` has two enumerations behind it (six `ErrorGroup`
  constants, documented as growing and `json:"-"` so not machine-readable at
  all; versus the integers `0,2,3,130` plus an unlisted default `1`).
- **Impl F12 (SEV-3)** — `graphlint.BlockingCodes()` is cited in
  §capability-dependencies but C4's row carries no symbol anchor while its
  neighbours all do.
- **PM S5** — C4 is a point-in-time snapshot ("at this RDR's implementation")
  while `0029:C2` says the tier lives in `docs/cli-output-contract.md` beside
  the vocabulary: two artifacts, two shapes, no statement of which the consumer
  reads or that they must agree. `0029:F3` concedes detection is a manual
  audit, and `0029:A4`'s own history is that the gap recurred twice inside the
  fix.
- **QA F5 (MAJOR)** — `docs/cli-output-contract.md` today contains zero
  occurrences of `frozen`/`append-only`/`growing`, no scenario covers the doc,
  and the register obligation lands only in Phase 2. The one test that would
  stop A4's recurrence is absent.

## Hotspot 2 — `0029:C2`'s `closed` retirement is under-scoped (2 personas)

- **QA F3 (BLOCKING)** — `0029:S9` scopes the grep to two files; `closed` is
  live describing tiered vocabularies in at least four more
  (`guard/grammar.go` x6 on the `frozen` `Operators()`, `accessor/model.go:308`,
  `resolve/resolve.go:45,66`, and `cli/flow_input.go:349`, which states the
  *opposite* of the `append-only` tier C4 now assigns the CLIError codes). S9 as
  written ships a false pass. Also undecided: whether a quotation of a peer
  contract counts as an occurrence (C4 itself quotes `0003:C7`'s "MUST be
  closed and typed").
- **Impl F3 (SEV-1)** — same clause from the other side: one occurrence is a
  shipped wire-visible identifier, `CodeCoverageClosedByEscape =
  "graph-coverage-closed-by-escape"`, a member of the `growing` advisory set,
  so a bare grep can never go to zero without a breaking rename. C2's scope is
  the whole repo and docs; §step-2 names two files. Needs an assertion form
  (allowlist grep? scoped grep? review item?).

## Hotspot 3 — `0029:C1`'s `schema_version` literal has no stated home or bump mechanism (2 personas)

- **Impl F2 (SEV-1)** — no contract says where the literal lives (`respond`?
  `clierr`? a third package — they are deliberately separate to avoid an import
  cycle) or what enforces the bump. A duplicated literal can drift, which is the
  defect `0029:D-naming` cites as this RDR's reason to exist.
- **QA F6 (MAJOR)** — `0029:MVV` step 4's oracle reads the movement by value,
  and its negative control ("a no-op release MUST move neither") is unwritable:
  nothing defines a release in-repo for a test to range over, and there is no
  mechanical link between "a vocabulary gained a member" and "the literal
  changed".

## Hotspot 4 — `0029:C1`'s `planEnvelope` anchor is misdescribed (2 personas)

- **Impl F1 (SEV-1)** — `planEnvelope` is a struct that decodes `--plan` INPUT
  (`flow_input.go:338`), not the output envelope, and there is no
  `func planEnvelope`. Its own doc comment says it discriminates on presence of
  a top-level `type` key and *refuses* a refusal envelope — so it is not the
  tolerant-reader precedent C1, `0029:A3` and `0029:S3` hold it up as. Open: does
  `schema_version` ride the top level only, or is it ever projected into `data`?
- **QA F6 (second half)** — the repo's one in-repo consumer has no
  `schema_version` field at all, and no scenario asks whether it implements C1's
  two consumer MUSTs. In scope (a Phase 1 obligation) or out (consumer rules
  advisory to third parties only)?

## Hotspot 5 — the envelope key sets in `0029:S1`/`0029:S2` ignore omitempty siblings (QA, single persona)

- **QA F4 (MAJOR)** — `respond.Success` marshals four fields (`Type`, `Notes`,
  `Warnings`, `Data`); the captured `data,type` baseline is one clean run's key
  set, not the envelope's shape. S1 read as an exact key-set assertion fails the
  moment `lint` emits a Note. S2 handles the refusal record's optional fields
  carefully — showing the RDR knows the distinction and did not apply it to S1 —
  but omits `Hint` (`clierr.go:60`).

## Hotspot 6 — cardinality assertions on a growing set (Impl, single persona)

- **Impl F8 (SEV-2)** — `0029:C2`'s `append-only` forbids cardinality
  assertions and `growing` is defined as "`append-only`, and additionally…",
  yet `0029:S8` keeps `TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers`
  alive "updated" for a vocabulary C4 moves to `growing`. Its name and its
  `want` literal both encode "exactly four"; updating it reproduces the defect
  at five.

## PM-only — does the record deliver its own stated outcome

- **PM S1 (top severity)** — `0029:§problem-statement` names a *pre-upgrade*
  discovery outcome ("needs to know, before upgrading"); every mechanism in
  `0029:§approach` is observable only post-upgrade (`C1` requires running the
  new binary; `C2`/`C4` live in a markdown doc; `C3` binds release notes).
  `0029:MVV` proves survival after the fact. Blocks: does this RDR owe a
  pre-upgrade affordance, or does §problem-statement re-scope to the governance
  problem the record actually solves?
- **PM S2** — Activation Step 1 is the entire user-facing delivery and is
  unfalsifiable as written. `llms.txt` opens by instructing agents to prefer the
  binary over any file, and describes the target doc as "worked JSON payloads
  and the rationale behind the envelope shape" — no term an agent scanning for
  compatibility matches on. The RDR requires the doc's body to change and never
  the pointer. `0029:BR4` weighed document-vs-document only; exposing the tier
  table from the binary was never considered.
- **PM S3** — the RDR names its own success test ("no code change at 1.0.0 —
  that is the test of whether this RDR worked") and places it past an event
  `0029:A2` establishes has not happened. §scope-verification is unanswered.
- **PM S4** — the `0.x` reframing (`§approach`, `C2`'s last paragraph: tiers are
  intent, not guarantees) leaves `C3`'s release-note disclosure as the only
  binding deliverable of the period in which the RDR ships — and release notes
  have no file, template, generator or check anywhere in the Implementation
  Plan. §problem-statement was never reconciled to the reframing.
- **PM S6** — the consumer discriminator rule ("`code`'s presence
  discriminates", the two terminal records being structurally asymmetric) exists
  only in `0029:MVV`'s trace table, not in any published contract clause, and is
  not among what Activation Step 1 publishes.
- **PM S7** — §scope-verification, §proportionality and §cross-cutting-concerns
  are unfilled template blocks. Stage-7 gate mechanics (matching this lens's
  baseline lint), raised only where it gates the S3 scope claim and the one-seam
  claim on a `Profile: large` record.

## Impl-only — remaining clarifications

- **F6 (SEV-2)** — `respond`'s package doc (`respond.go:13,24,49`) and
  `readPlan`'s comment (`flow_input.go:361`) still document the refusal envelope
  as `{"type":"failed"}`. Behaviour matches C1; the machine-read documentation
  does not. A reviewer reads `0029:S2` as a regression. Does a step correct them,
  and is `docs/cli-output-contract.md` stale the same way?
- **F7 (SEV-2)** — Step 1 re-captures another RDR's checked-in artifact
  (0023's golden). §prerequisites rules on the `0006:C17` amendment but says
  nothing about 0023. Mechanical refresh, or cross-RDR coordination?
- **F9 (SEV-3)** — `internal/version`'s `--as=json` payload takes no tier and is
  not among C4's named exemptions, though `0029:D-naming` reasons about it
  directly. Under C4's own closing sentence it must be tiered or declared exempt.
- **F10 (SEV-3)** — `0029:D-identity` ("same shape iff same major") read alone
  is unimplementable during `0.x`, where every envelope shares major `0` while
  `C1` says they may change incompatibly. C1 carries the resolution; D-identity
  does not carry the qualifier.
- **F11 (SEV-3)** — C1's tolerant-reader consumer MUSTs have no CLI-side
  observable; `0029:S3` tests the one in-repo decoder as a proxy.

## QA-only — remaining

- **F7 (MINOR)** — `0029:S5` and MVV step 5 are the same assertion; promotion
  means mutating a package-level var, and neither says whether a test does that
  (no seam exists) or whether S5 is manual-only. `0029:S8` is explicit about
  being a source edit, which is why S8 is actionable and S5 is not.
- **F8 (MINOR)** — `0029:S3`'s "differs from its predecessor by exactly one key"
  is a review criterion over a git diff, not a test; the predecessor leaves the
  tree on re-capture. Stated in the Testing Strategy as though it were a test row.

## Checked and clean (all three personas reported explicit non-findings)

- Every test symbol and artifact path the RDR names exists where it says (QA).
- `0029:S4` writable as specified; `0029:S7`'s refusal to add a byte-golden is
  correct and well-argued; `§performance-expectations` correctly scopes the
  byte-stability non-claim (QA).
- `0029:A1`'s claim that the four-element `want` literal is the only blocker on a
  fifth advisory code holds — `severityFor` derives, nothing counts (QA).
- C4's two deliberately tier-less namespaces are consistently handled and A4
  backs them (QA).
- `§decision-rationale`'s reframing of the seed's fork, and `ALT1`'s rejection
  against four shipped growth events, are grounded rather than preference (PM).
- `D-naming`'s rejection of `version` verified against `internal/cli/version.go`
  — the collision would have been real (PM).
- `0029:C1`'s tolerant-reader/strict-input asymmetry is argued from audience
  rather than convenience (PM).
