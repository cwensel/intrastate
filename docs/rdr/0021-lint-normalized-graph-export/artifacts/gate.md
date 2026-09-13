# Finalization Gate — RDR 0021, lint's normalized-graph export

- **Record**: `cli/0021` (`0021-lint-normalized-graph-export`)
- **Date**: 2026-09-12
- **Verdict**: READY — Gate PASS, locked to Final (re-lock after the
  `0021-0029` cluster reconcile).

Mechanical pre-sweep: `evidence/tooling-pass/iter-3/tooling-pass.md` — PASS
(`rdr lint --locking` exit 0, `blocking=0 resolution=0 placeholder=0`). The
7.1 re-entry note's six items all verify closed in live text (that report
carries the item-by-item table); the note is deleted at this lock. Item 4,
Cross-Cutting Concerns, is authored in the record at `0021:G-cross-cutting`
and is deliberately not copied here.

## 1. Contradiction Check

No contradictions between Research Findings and the Proposed Solution, and
none between planned features and stated principles.

- Findings read `0005:C1` as carving `lint`/`dump`/`parse`-class command
  groups out of the flow contract, "owned by the RDR that names them". C1
  claims a root `graph` verb under exactly that carve-out, and A4 (Method:
  Peer RDR) independently verifies no 0005 envelope amendment is needed —
  now re-verified against the post-0029 envelope.
- Findings read 0002 as fixing the dump's field list and row order while
  declining to define a dump grammar. C2 owns the JSON schema on that seeded
  ground and reuses 0002's field vocabulary and row order rather than minting
  a second row contract, amending no 0002 text, and closes
  `0002:§round-trip-inverse-invariants`'s lossy set-literal rendering for
  this document only, claiming no export→load inverse (RT3).
- Findings record a negative corpus result: no opened peer wraps DOT in a
  JSON envelope. C5 claims no peer precedent for that cell; it defines it
  normatively as `data.dot`, one required string member, unwrapped
  `jq -r .data.dot`. The negative result is carried as a constraint, not
  silently overridden.

Principles vs planned features: the neutrality principle ("an export is never
a lint pass") is not merely asserted — C4 makes it normative and MVV step 5
gives it a mechanism-independent oracle, byte-comparing lint's refusal with
and without the export present. C2's "no verdict or finding field, and no
per-node terminal marking" applies the same principle to the document
surface, leaving the dead-end quantifier to RDR 0015 (JDR 0001 §JD-23).

Post-reconcile check: the tier declarations added to C1 and C2 agree with the
additive-within-`/1` rule C2 already stated — the STABILITY paragraph names
that as the same promise in `0029:C4`'s tier vocabulary, not a second rule.
No clause contradicts another.

## 2. Assumption Verification

All eight Critical Assumptions are internally consistent and terminal.

- **Status**: 8/8 `Verified`; zero `Pending`, `Unverified`, or placeholder
  (`ca=all-terminal`). No settled-fact prose leans on an unverified property.
- **Method**: every label sanctioned — `Spike` (A1, A2, A5, A6),
  `Source Search` (A3, A7, A8), `Peer RDR` (A4). Zero off-vocabulary members.
  No `Docs Only` record exists, so the load-bearing Docs-Only bar is vacuous
  rather than waived.
- **Evidence**: Status, Method and Evidence agree on every row; each
  "If wrong" is non-empty.
- **A4, the reconcile's `re-verify` target**: discharged in place. Its premise
  is `0005:C1`'s scope sentence, not envelope immutability, so `0029:C1`'s
  non-`omitempty` `schema_version` leaves the conclusion standing; C-11's
  second ask is answered with a negative — this record asserts no exact
  envelope key set, C2 predicating its field list of the DOCUMENT.
- **Self-reference**: none. The three `Source Search` rows resolve into the
  product tree, never to this record or its artifact directory.
- **Symbol resolution**: every `source-anchor` edge reports `resolved: true`
  — none false, none absent, so the lookups genuinely ran. A5's Evidence is
  written `clierr.go:174::WriteJSONLine`; the symbol resolves and the stale
  line component is a documented non-finding.

## 3. Scope Verification

The Minimum Viable Validation is in scope and executed during implementation,
not deferred. It is Phase-1 work, gated on one authored fixture pair rather
than on any later phase.

The specific proof: author a two-owned-state / one-terminal / one-escape-row
state-machine fixture plus a decision-table fixture, then assert four
invocations — (a) `intrastate graph --model <fixture>` run twice, stdouts
BYTE-identical, parsing as JSON and carrying every C2 field including the
`reach` abstraction marker; (b) `--emit dot | dot -Tsvg` renders and its
node/edge id set equals the JSON `reach` block; (c) `--as=json | jq .data`
equals the document value-for-value; (d) `intrastate lint` over a fixture with
blocking findings refuses byte-identically against a pre-change capture, while
`graph` over that same model succeeds with an asserted document. That fourth
invocation is the neutrality oracle, and C4 binds it mechanism-independently.

Blast radius: not applicable. This record is unclustered as it locks
(`clustered=false`, `cluster=[]`, no `impact_families`) — the `0021-0029`
pairing rests on reciprocal `cross-cutting-owner` edges, not a declared
`Cluster:` field — so it retires and renames no peer's literals and no
`impact.md` projection is owed.

Joint-decision fence: `op = none` (`fence-clear`) after one real blocker was
fired and cleared in this pass. `index --literal-intersect` reported
`0017 0021 UNCITED 1 shared: code`. Grounded per §ground-before-ask
(`--outcome ground` → `apply`, settled in source), the check fires CLEAR:
`0017:C1` decides what `findings[i].code` names for a multi-subject producer;
`0021:C4` forbids this verb to emit findings at all, its one new code
`graph-export-too-large` is a scalar ENVELOPE code, and the `model-invalid`
arm C1 mirrors reaches the shared
`internal/cli/flow_input.go::loadFindings`, which already populates each
entry's `code` with the REQ-24 load-category slug (`Code: string(category)`)
rather than the envelope code. 0017's rule is satisfied by the shipped loader
this record reuses; 0021 decides nothing in 0017's domain. Recorded as a
citation on the Joint-check line, not as a synced copy. JC1's own fire
(→ 0029, home `cli/0029 §Normative Contracts` C4) is `homed`, and C-7's
two-marker read-order question stays at that home, cited here rather than
restated. `rulings_open=0` — all nine author rulings are marked absorbed.

## 5. Proportionality

Right-sized; nothing flagged to trim before locking.

**Contract count, not word count.** C1–C5 are five labelled clauses of ONE
independent load-bearing contract: the export of the normalized graph. C1
fixes its surface (verb, arm set, flag order), C2 the document it emits, C3
that document's byte stability, C4 the neutrality boundary against lint, C5
how the two `--as` modes carry it. None is separately adoptable — a consumer
cannot take the document without the surface that emits it — so this is one
seam stated five ways, not five seams locked together. No split, and the
author ruling Q1 records why collapsing to a single `**C1**` is blocked:
`0029` (Final) holds a resolved `cross-cutting-owner` edge into `0021:C2` and
cites `0021:C5` by label, so relabelling would break inbound citations from a
no-amend record. `contracts_transient=0`, so no lifespan disposition distorts
the count.

**Profile re-validated.** The Metadata field reads `large`, and that still
matches the contracts just counted: one contract, user-facing yes (a new root
verb and a documented output document), locking the `intrastate.graph/1`
field list and marker, which `0029` consumes by a resolved edge. The form is
correct — value plus one clause naming the contract, no matrix or provenance
prose left from the template. The lens battery `large` demands did run:
grounding, 3amigo, critique (two models, differing stamps, diff written) and
repeatability-lite (`--outcome repeatability` → `none`, `rule =
repeatability-lite-complete`). `lens_stale=none` — the qualifier is cleared
and no lens folder predates it.

**Growth from the reconcile.** The re-entry added two tier declarations and
two peer citations to text this record already carried, and this lock removes
a 102-line re-entry note. The density sits in C2, where each field spelling is
normative and therefore load-bearing at implementation; the Pre-Lock
Mini-Checks and Decision Rationale carry the reasoning that keeps those
spellings from being re-litigated in Phase 1. Proportionate.
