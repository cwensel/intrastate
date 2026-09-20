Model: claude-opus-5[1m]

# 3amigo — Persona 1: Product Manager — cli/0012

## Widened

The starting set (§problem-statement, §approach, §decision-rationale, MVV)
did not carry enough to answer "does the user get the outcome". Three
silences sent me out:

- MVV step 2 names a refusal payload but not who experiences it, so I read
  §failure-modes, §consequences and §pre-lock-mini-checks (`disposition`
  table) for the user-visible break inventory.
- The problem statement scopes the defect to the OWNED reader ingress, but
  C5 lands a refusal on the CLI and on authored models — a second and third
  user population the stated problem does not name. I read C1, C2, C4, C5,
  A1, A3, A4, S5, S6 and §existing-infrastructure-audit to see whose
  outcome each clause changes.
- Nothing in the starting set addresses existing persisted artifacts on
  upgrade, so I grepped a scratch copy of the record for
  persisted/migrate/upgrade and then grounded the owned path in source
  (`internal/accessor/model.go::OwnedSnapshot` at :419/:424,
  `internal/resolve/resolve.go::assemble` merging `ProvenanceOwned` at
  :230, `internal/cli/flow_next.go::probeRow` at :469).

Grounding confirms the record's factual spine: the owned ingress is real
and unconformed, `probeRow` does route the same seam, and
`resolve.ReasonUncomparable` is the existing vocabulary the fix reuses. No
finding below disputes a fact; they are all about stated user outcome.

## Findings

### PM-1 — The stated problem is one ingress; the shipped change breaks two more — severity: medium
- anchor: 0012:C5 (vs 0012:§problem-statement)
- finding: §problem-statement is emphatic that "the defect's home is the
  OWNED ingress, and only it", and explicitly clears the CLI path
  ("`--tag iter=many` refuses `flow-tag-invalid` at input and never reaches
  the seam"). C5 then lands a NEW user-visible refusal on exactly that
  cleared path (`--tag n=07`, `--write n=+1` — accepted today, refused
  after) and on authored model text. §failure-modes F1 and the
  `disposition` mini-check both own this honestly ("a second visible
  break"), and D-c5-is-not-alternative-3 argues the venue is legitimate.
  What is missing is the user-outcome framing: the problem the record
  states is a silent misroute for owned readers; the problem C5 solves is
  lint/runtime disagreement (A4), which no user reported and which S5
  measures at zero occurrences in a 123-model corpus. A reader of
  §problem-statement alone cannot tell that shipping this RDR takes away a
  CLI input form that works today.
- blocks: whether C5 ships in this RDR at all, and whether it ships in the
  same release as C1/C2. A stakeholder weighing "fix a silent misroute"
  against "break `--tag n=07` for every existing script" cannot make that
  call from the stated problem, because the second half of the trade is
  only discoverable in F1 and the mini-check table. Phase 4 is separable
  (its own phase, its own fixture S6) — nothing in the record says whether
  it MUST be atomic with Phase 1, and the F1 argument ("leaving the CLI lax
  would let `--tag iter=07` flip a guard verdict") argues it must, which is
  a release-scoping decision the reader should not have to reconstruct.

### PM-2 — No account of existing persisted artifacts at upgrade — severity: medium
- anchor: 0012:MVV (and 0012:F1)
- finding: The MVV's user is a freshly-authored `mvv-int` model with a
  reader returning a chosen value. F1's recovery line is "fixing the
  reader's value or the persisted artifact" — the only place the installed
  base appears. But §problem-statement establishes that "every persisted
  owned value and every hand edit to the artifact enters through" the owned
  door, i.e. the affected population is every artifact already on disk. The
  record measures the authored-model blast radius precisely (S5: 123
  models, zero `eq`/`in` over an `int` tag) and measures nothing on the
  held side, where the flip is both directions: malformed → Unevaluable
  (loud, the fix) and non-canonical `"07"` → False-to-True (a SILENT verdict
  change — C5 explicitly does not reach reader values, per the `trace`
  mini-check step 3). Consequences names this flip; no one sized it.
- blocks: the rollout decision. A reader cannot decide whether this is a
  patch-level behavior fix or a change needing a release note, a scan, or a
  `--fix` affordance, because no one knows how many artifacts in use hold a
  value that will newly refuse (loud) or newly match (silent). The silent
  leg is the sharper gap: the RDR's whole premise is that silent verdict
  changes are the harm, and it ships one unmeasured.

### PM-3 — The named beneficiary appears once and is never carried to an outcome — severity: low
- anchor: 0012:§problem-statement
- finding: "For a consumer like the rdr navigator that is a silent
  misroute" is the only appearance of an actual user of this system. The
  MVV's end-state reduces to "step 2's invocation is the seed defect
  (kata `intrastate#cq5p`) refusing loudly instead of misrouting" — a
  synthetic model, not the navigator. Nothing states whether the navigator
  currently misroutes in practice, or what the navigator operator sees
  after the fix (a refusal where a plan used to appear — is that better for
  them, or a new failure to handle?).
- blocks: acceptance. Whoever signs off cannot confirm the reported defect
  is actually gone for the reporting consumer, only that a constructed
  fixture behaves. Also blocks the downstream question of whether the
  navigator (and any other `resolve.Resolve` library caller) needs to
  handle a newly-possible `flow-guard-unevaluable` exit it never saw — the
  record asserts "every `GuardEvaluator` consumer inherits the fix"
  (§consequences) without asking whether every consumer is ready to receive
  a refusal instead of a plan.

### PM-4 — Profile claims "user-facing yes" but no gate response backs it — severity: low
- anchor: 0012:G-cross-cutting (and 0012:G-proportionality)
- finding: §metadata Profile reads "foundational — ... user-facing yes;
  locks cross-rdr", and the record does change three user-visible surfaces
  (`flow resolve` exit, `flow next` reason vocabulary, CLI `--tag`/`--write`
  acceptance). But G-cross-cutting and G-proportionality are still
  unfilled TEMPLATE guidance blocks — the record's own lint confirms
  (`placeholder:survived` at 1205-1209 and 1230-1255, five total). Status
  is Draft so this is expected sequencing, not a defect; I flag it because
  the two gates that would force a written answer to "incremental
  adoption?" and "is the Profile right given the user-facing surface?" are
  exactly the two that would have caught PM-1 and PM-2.
- blocks: nothing yet — but it means the lock gate is the last chance to
  decide C5's release scoping and the persisted-artifact blast radius. If
  those gates are answered mechanically, both findings above ship unstated.

Note on what is NOT a finding: the record is unusually strong on outcome
mechanics. The `disposition` mini-check is a complete input-class → what
the user sees table; the `oracle` mini-check states what each MVV step
fails on plus a negative control (the `4`-prunes leg); F1 walks the
second verb (`flow next`) rather than stopping at `flow resolve`; and
§decision-rationale's O4 row rejects the alternative on the user-outcome
ground (a non-CLI caller keeps the misroute) rather than on architecture
taste. I did not manufacture findings against those.
