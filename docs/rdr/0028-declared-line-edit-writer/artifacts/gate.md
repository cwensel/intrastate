# Finalization Gate — cli/0028 declared-line-edit-writer

- **Record**: `0028-declared-line-edit-writer.md`
- **Date**: 2026-09-03
- **Verdict**: **PASS** — locked to Final
- **Mechanical pre-sweep**: PASS (`evidence/tooling-pass/tooling-pass.md`;
  `rdr lint --locking 0028` exit 0, blocking=0 resolution=0)

Item 4 (Cross-Cutting Concerns) is authored in the record at
`cli/0028:G-cross-cutting` and is retained there at lock, because peer RDRs
cite it. It is deliberately not copied here.

## 1. Contradiction Check

No contradictions found between research findings, design principles, and the
proposed solution.

Compared: Research Findings (Investigation, Key Discoveries) against Approach,
Normative Contracts C1.1–C1.6, Decision Rationale and Consequences.

The one place the research and the design visibly differ is not a
contradiction but a deliberate, recorded divergence: the surveyed prior art
(Ansible `ansible.builtin.lineinfile`, Puppet `puppetlabs-stdlib file_line`,
kubebuilder `InsertCode`/`ReplaceInFile`) all admit last-match resolution and
insert-on-no-match, and C1.3 `select:` refuses both — exactly one line must
match, 0 is `edit_anchor_unmatched` and ≥2 is `edit_anchor_ambiguous`, with
creation fenced out. The RDR states this as being stricter than its prior art
and gives the reason (an unmatched anchor is a stale model, not a missing
line), so the finding and the design agree about what the finding was.

Principles checked against features rather than assumed:

- *The model carries no path.* C1.3 `target:` states an `edit` entry has
  exactly the authority a `path` entry has over the caller-bound file. The
  two candidate counter-features are A1's context-tag carriage and C1.6's
  `{tag.<key>}` argv placeholder; both carry record IDENTITY for anchor and
  argv selection, never a path — the caller still binds the artifact via
  `--artifact`. No feature reintroduces a model-carried path.
- *The model does not carry the write body.* Preserved: the `edit` block is a
  declaration grammar (anchor + template + closed placeholder vocabulary), not
  a body the model transports, which is the 0025 consequence this RDR names in
  Overrides.
- *Reuse stops where the borrowed mechanism's purpose stops.* Applied
  consistently and stated twice against itself: `flowbind.go::save`'s
  stage-and-rename discipline is reused while its `Writer` unreachable-locator
  seal is explicitly NOT (a markdown target has nowhere to hold it), and
  `flowbind.go::load`'s absent-is-empty precedent is explicitly NOT taken (it
  exists so a first write can create a flow artifact; `edit` fences creation
  out). Both refusals are argued from the borrowed mechanism's purpose, not
  from convenience.

## 2. Assumption Verification

All 12 Critical Assumption Evidence Records (A1–A12) are internally consistent
and terminal. No blocker.

- **Status**: 12/12 `Verified`; 0 `Pending`, 0 `Unverified`, 0 placeholder
  (`rdr status --tags`: `ca=all-terminal`, `ca_pending_ids=[]`).
- **Method vocabulary**: all in-sanctioned-set, 0 off-vocabulary
  (`ca_off_vocabulary_ids=[]`). Distribution: Source Search ×8 (A1, A2, A6,
  A7, A9, A10, A11, A12), Spike ×3 (A3, A4, A5), Peer RDR ×1 (A8). Every
  record carries a Method field.
- **Status/Method/Evidence agreement**: holds for every row. Source Search rows
  cite resolved source sites, not spike output; Spike rows cite run artifacts
  under `evidence/spikes/`; A8 (Peer RDR) cites JDR 0003 §D1 by element and
  quotes its disposition rather than restating the rule.
- **"If wrong" non-empty and substantive** on all 12 — each names a specific
  consequence, not a generic hedge (e.g. A2 names two distinct fallbacks, A8
  states the claim is void per JDR 0003 §D1, A12 names a concrete
  stronger-check fallback).
- **No `Docs Only` records at all**, so the load-bearing-Docs-Only blocker
  cannot arise.
- **No Source Search self-reference** (CHECK 3): every Source Search row's
  Evidence resolves into the consumer source tree; none resolves to the RDR
  itself or under its artifact directory.
- **Symbol resolution** (CHECK 5): 60/60 `source-anchor` edges resolve `true`
  against the repo; none `false`, none absent-unlooked.
- **Status consistency** (CHECK 6): with every assumption terminal there is no
  `Pending`/`Unverified` property relied on as settled fact elsewhere, and no
  checklist-vs-gate disagreement is possible — the gate responses live only
  here and the record carries no second copy.

Two `mentions` edges report unresolved (`0028:MVV → 0022`, `0028:S12 → 0021`).
Both are correct as written and neither is an assumption defect: they are
references to test FIXTURE filenames under
`rdr/tools/rdr/testdata/status/records/` (`0022-cache-metrics-surface.md`,
`0021-cache-warmup-order.md`), which the linter's untyped `mentions` heuristic
pattern-matches on the `NNNN-slug` shape. They are non-blocking, and rewriting
them to point at peer records would make them wrong.

## 3. Scope Verification

The Minimum Viable Validation is in scope and will be executed during
implementation. It is not deferred.

**The specific proof**: a `rdr-write.toml` re-authored as a state-machine over
four declared bindings — a 0025 command reader on role `record`, an `edit`
writer for that record's Status line, a C1.6 command reader over the README row
on role `readme` (with `{tag.nnnn}` in its argv), and an `edit` writer for that
row anchored by `{tag.nnnn}` — driven over three fixture records (a plain
`Draft`; a `Draft [joint decision → …]` on one line; and one whose bracketed
qualifier WRAPS onto a continuation line while still reading `Draft`). Per
record, one invocation of
`flow --allow-commands resolve … --outcome lock | flow --allow-commands
set-state … --plan -` must exit 0 and leave the record's Status reading `Final`
(qualifier preserved by backreference; the wrapped record's continuation line
byte-identical), the README row's status cell reading `Final`, both confirmed
by read-back through the command reader, and every other byte of both files
unchanged — `git diff --stat` showing exactly one changed line per file per
invocation. Step 1 additionally requires `intrastate lint --model
rdr-write.toml` to pass with NO wrapper script anywhere in the fixture, which
is the actual claim of the RDR: the edit is data.

Five negative arms are part of the MVV, not follow-on work: no
`--allow-commands` refuses before mutation; a duplicated README row refuses
`edit_anchor_ambiguous`; a `replace` whose output no longer matches its own
anchor refuses `edit_anchor_unstable`; and — required rather than optional — a
README carrying a row in the live corpus's unlinked `| NNNN |` form refuses
`edit_anchor_unmatched` for that record while the other rows still flip. That
last arm is what keeps the MVV from proving the scenario only against a fixture
curated to exclude the one class of drift the real corpus actually has.

**Prerequisite status, stated rather than glossed**: Phase 4 (the end-to-end
proof, i.e. MVV execution) is BLOCKED on a consumer-side prerequisite the
record names — the row-addressed projector verb the C1.6 reader invokes does
not ship today. Phases 1–3 (grammar and lint; the binding; seam carriage and
gate pre-check) are unblocked and independent of it. This is a sequencing
constraint on implementation, not a deferral of the MVV: the RDR carries the
prerequisite explicitly in Prerequisites and gates Phase 4 on it, which is the
correct shape. Implementation must not report the MVV satisfied until that verb
exists and steps 1–8 have run.

## 5. Proportionality

Right-sized. Nothing to trim before locking.

**Contract count is one** — C1, the `edit` write carrier, whose clauses
C1.1–C1.5 are one grammar over one seam (`internal/table/load.go` carrier
admission plus `internal/accessor` write execution). C1.6 was scored separately
in the Decision Rationale and counted here rather than fenced past: it extends
0025:C2's argv placeholder vocabulary on the READ path, in a different package
(`cmdbind`), which is the shape a split test looks for. It fails the
deployability test for an independent contract — `{tag.<key>}` in a command
reader's argv exists solely so the write carrier's 0004:C12 read-back can
address one row of a shared artifact; it carries no state and no behaviour of
its own, and shipping it without C1.1–C1.5 addresses rows on behalf of nothing.
Not independently justifiable, so not independently locked. The RDR is the sole
author of one contract; no split.

**Profile `large` re-validated against the contracts just counted.** One
contract; user-facing yes (a declared grammar model authors write); locks
format (the `edit` block's declaration grammar and its closed placeholder
vocabulary). The value Resolve wrote still matches what the lenses found. Not
`foundational`: Seam Lineage records no prior accretion on this seam
(`seam_lineage=0`, provenance recorded in the Metadata field). The lens battery
`large` requires did run — grounding, 3amigo, critique (two models, differing
stamps), repeatability-lite (run-1 + diff, `emit.next: none`) — so the Profile
did not route past anything it owed.

**Reviewed at critique**: the reverse reading — C1.6 as a second seam warranting
its own record — was weighed and rejected above rather than deferred. What the
critique lens changed is the honesty of the surrounding claims, not the count:
C1.6's consumer verb is now a blocking MVV prerequisite (Prerequisites) rather
than a Spec Impact cell.

**Length**: the record is long (~1640 lines) for one contract, and the length
was examined rather than accepted. It is carried by material that earns its
place — 12 verified assumptions with real evidence spans, 32 testing scenarios
that ARE the contract's executable form, and an Alternatives section whose
three rejected designs each fix a boundary C1 depends on. The one advisory the
sweep raised is A2's 35-line Evidence field (soft cap 30, `evidence:over-budget`).
Answering the Gate's question directly: the load-bearing anchors are still
findable in it — the field leads with five resolved `path::Symbol` anchors
before its prose — and the balance is genuine verification content the grounding
sweep reads. Not truncated, and not moved to an artifact, because moving it
would put the reasoning a step further from the claim it verifies for a
five-line saving.
