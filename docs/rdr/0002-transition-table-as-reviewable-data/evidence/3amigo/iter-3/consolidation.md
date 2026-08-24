Model: claude-opus-5[1m]

# 3amigo Consolidation — RDR 0002, iteration 3

Origin ledger for this lens's resolve pass. Three isolated persona passes
(PM / Implementer / QA), each run as a separate sub-agent seeing only the RDR
and its own persona block. Isolation verified: no persona file references
another persona's finding ids.

**Delta scope**: the JDR 0001 §D7 re-authoring (refine `3c3cb3e`, resolve
`929e931`) and everything the cove iteration-2 pass touched. The iter-1 ledger
(HOT-001…HOT-004) and the iter-2 ledger (HOT2-001) are closed and were not
re-opened; no persona re-raised them.

Consolidation is **mechanical**: hotspots are the passages two or more personas
named independently, computed from the `Passage:` anchors, not re-judged by a
model that read all three. Overlap marks a hotspot passage, not a validated
finding; a single-persona finding is not thereby weaker.

## Hotspots (2+ personas, independently)

- **HOT3-001 the joined-rendering defect is not fully closed — it survives in
  write values, and A13's stated hazard contradicts the trace** (PM3-001,
  IMP3-003, QA3-001) — three personas reached the cove pass's central edit from
  three directions. IMP3-003: the literals clause binds *atom* literals to a
  member sequence but says nothing about set-valued **write** values, which the
  spike comma-joins (`main.go::renderValue`), reproducing the banned collision
  one field over — in a value RDR 0004 compares for equality on read-back.
  PM3-001: A13's `If wrong` says the failure is duplicate atoms "without
  changing match semantics", while the desk trace's `merge atoms` row records
  the opposite (a live expanding row where the author wrote a dead rule), and
  Prerequisites rests the lock on that `If wrong` text. QA3-001: the
  delimiter-bearing control A13 is Pending on has no determinate pass/fail arm —
  its second half is an RDR 0006 verdict the MVV refuses to count, its first
  half is written against "the implementation's own joining delimiter", which a
  black-box test cannot read.
  Passages: Normative Contracts / literals clause; write-replaces clause;
  Critical Assumptions A13; Conditional Mini-Checks `trace` (merge atoms row);
  Implementation Plan / Prerequisites; Testing Strategy scenario 2.

- **HOT3-002 the fail-fast clause is under-specified against the two-pass gate
  and unweighed against the authoring workflow** (IMP3-008, PM3-003, QA3-007) —
  IMP3-008: "the first category a document trips is the refusal" is
  unimplementable across the mandated two-pass gate, because reaching `[model]`
  requires a parse, so a malformed **and** v2 document necessarily refuses
  `malformed TOML` while the gate's precedence is stated absolutely. PM3-003:
  the clause fixes single-category refusal for a hand-authored model carrying 23
  categories without weighing it against the authoring ergonomics the Decision
  Rationale names as the RDR's purpose, and forecloses an accumulating mode in
  RDR 0005 without making that call here. QA3-007: the clause's absorbed-defect
  carve-out names a class and routes its coverage to the merge/literals clauses,
  whose control is the one QA3-001 shows is owed — so zero live assertions cover
  it.
  Passages: Normative Contracts / version-gate clause (two-pass paragraph) and
  the "Load is fail-fast" paragraphs; Decision Rationale.

- **HOT3-003 `[model.metadata]` is fully specified, has no consumer, and no
  oracle past decode** (PM3-004, QA3-002) — PM3-004: the clause mandates a
  normalized-model field that Capability Dependencies names no consumer for,
  Phase 2's deliverable list omits, and no peer RDR reads (grep of
  `docs/rdr/0006-*.md` finds no `metadata` reference; JDR 0001 §D7(v) justifies
  it as "one sanctioned namespace for tooling", naming no tool) — so
  decode-and-discard satisfies every asserted obligation. QA3-002: "carried
  verbatim" has no oracle, because the clause removes it from the Round-Trip
  invariant and the dump field list, the only two comparisons the RDR defines,
  leaving top-level key names as the sole witness — which a loader discarding
  every value would pass.
  Passages: Normative Contracts / layout block, `[model.metadata]` clause;
  Testing Strategy scenario 1; Conditional Mini-Checks `fidelity`;
  Implementation Plan / Phase 2; Capability Dependencies.

- **HOT3-004 the source locator has no provenance, and its Round-Trip
  comparison has no projection** (IMP3-002, QA3-008) — IMP3-002: the normative
  floor ("at least the model id and rule id") is contradicted by the only
  implementation, which stores the authored `source` string; both `kata-fixture`
  rules author `source = "kata:review"`, so two distinct rows carry one locator
  and the "rule-identifying part" comparison is unwritable. QA3-008: comparing
  presence and the rule-identifying part while excluding line/column requires a
  structured locator, and the RDR specifies only a floor while the spike's is an
  opaque token.
  Passages: Technical Design / normalized candidate row (locator floor);
  Round-Trip / Inverse Invariants (locator paragraph).

- **HOT3-005 the load-time category floor does not cover the shapes scenario 3
  asserts on** (IMP3-004, IMP3-005, IMP3-006, IMP3-007, QA3-003, QA3-004) — six
  findings across two personas against one contract surface. Guard-block
  operators are unconstrained while `unknown operator` is a floor category
  (IMP3-004); `[model]` well-formedness — absent `id`, absent `version`
  (silently `unsupported version 0`), absent `[tags.recognized]` — has no
  category (IMP3-005); the accessor-entry clause states a rule only for
  `timeout`, and that rule has no oracle and is unimplemented (IMP3-006);
  `[dump]`'s column vocabulary and malformed forms are unspecified with no
  category (IMP3-007). On the test side, the `[initial]`-undeclared-key control
  is assigned a category the spike refuses under a different one (QA3-003), and
  three escape-shape controls share one category so they cannot discriminate
  under the category-not-message oracle the draft mandates (QA3-004).
  Passages: Normative Contracts / load-time category floor, match-block operator
  restriction, version-gate clause, Accessor tables clause, dump clause;
  Testing Strategy scenario 3.

## Single-persona findings (not weaker — no overlap by construction)

- **IMP3-001 the closed layout omits `[[rule]].source` and `[model].description`,
  which every rule in the normative fixture authors** — the closed-layout clause
  enumerates admitted keys exhaustively and closes "No other root key or table is
  admitted", but both fixtures author `source` (and `[model].description`), and
  scenario 1 requires those fixtures to decode under strict decoding. An
  implementer building the struct from the normative text gets a decoder that
  refuses the normative fixture — the same failure the RDR celebrates for the
  pre-§D7 fixtures. Highest-severity implementer finding; the only one that makes
  the promoted fixtures fail their own contract.
  Passage: Normative Contracts / closed-layout clause; Load-Bearing Decisions /
  Wire / byte format.

- **PM3-002 the reviewable artifact is admitted ambiguous where reviewers read
  it, charted to no named successor** — the Round-Trip section admits the
  rendered dump makes a set-valued atom literal non-recoverable, answers "the
  normalized value is the contract", and charts the inverse to "a successor
  dump-format RDR" named nowhere in Capability Dependencies, References, or the
  phase plan. The same collision is treated as fatal one layer up and tolerated
  one layer down where a human disambiguates.
  Passage: Round-Trip / Inverse Invariants (lossy-site paragraph); Conditional
  Mini-Checks `fidelity` (dump row, read-back row).

- **QA3-005 witness counts stated three incompatible ways** — scenario 3 says
  35; A1 says 32 negative + 2 positive; the transcript has 37 refusals and 5
  positive controls. Prevents the promotion guard the not-narrowed rule calls
  for.
  Passage: Testing Strategy scenario 3 Expected vs. Critical Assumptions A1
  Evidence.

- **QA3-006 probe fixtures cited by non-existent paths** — cited as
  `evidence/spikes/iter-2/probe-a.toml` / `probe-b.toml`; actual paths are
  `neg/probe-a-initial-observed-no-writer.toml` and
  `neg/probe-b-initial-observed-with-writer.toml`. Prevents re-verification of
  the enumeration that justifies omitting the observed-key control.
  Passage: Testing Strategy scenario 3, probe evidence citation.

## Recorded non-findings (checked, not defects)

- The `6ccfe901…` digest **does** reproduce over the row block of `output.txt`;
  the full-file SHA differs only because the file carries a header line.
- The 18-of-23 category count in scenario 3 checks out against the transcript.
- Both `output.txt` and `negative-cases.txt` open with a stray
  `Model: claude-opus-5[1m]` stamp — an evidence-file convention leaking into
  spike artifacts. Harmless to the digest as computed, but it travels into the
  production test tree on promotion. Out of scope for this draft's contracts;
  a hygiene item for the resolve pass to route.
