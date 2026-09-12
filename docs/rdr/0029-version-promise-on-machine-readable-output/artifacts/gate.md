# Finalization Gate — cli/0029 version-promise-on-machine-readable-output

Date: 2026-09-12 · Profile: `large` · Verdict: **READY — Gate PASS**

Mechanical pre-sweep: `evidence/tooling-pass/tooling-pass.md` — PASS
(one hollow section fixed in-pass, re-run clean). `rdr lint --locking`
exit 0, `blocking=0 resolution=0`. Item 4, Cross-Cutting Concerns, is
authored in the record and cited as `cli/0029:G-cross-cutting`; it is not
copied here.

## 1. Contradiction Check

No contradictions between Research Findings and the Proposed Solution.
Each contract traces to a discovery rather than to preference:

- The research found `closed` live in two incompatible senses on adjacent
  surfaces (`graphlint::CodeCoverageClosedByEscape`'s "CLOSED at these
  four" against `table::Categories()`'s "closed load-category set" beside
  "The list's size is not a contract"). C2 retires the descriptive use and
  substitutes a three-tier vocabulary — the contradiction the research
  named is the one the contract resolves.
- The research found vocabulary growth is this project's normal release
  event (`table.Categories()` grew under 0002, 0024, 0025, 0028 before any
  tag). C4's `append-only`/`growing` assignments match observed behaviour;
  ALT1 (freeze within a major) is rejected against those same four events,
  not against taste.
- The research found the verdict-safe introduction mechanism already ships
  (`graphlint::Report.Blocking()`/`Advisory()`, `0006:C17`). C3 cites that
  guarantee rather than restating it, and mints no new mechanism —
  consistent with Proportionality's finding that C3 rides `0006:C17`.
- OpenTofu's `json-format.mdx` supplies both halves of C1's consumer rule
  verbatim and the sibling-field placement; C1 adopts both.

Planned features vs stated principles — one tension examined and resolved,
not papered over. C1 imposes a TOLERANT reader on emitted envelopes while
`0002:C3` imposes a STRICT reader on authored model TOML. The RDR states
this inversion is intended and follows the audience (an author learns of a
typo; a consumer survives an addition), and Key Discoveries is explicit
that the literature (Daigneau p.244, DDIA p.121) justifies
tolerant-on-receive/validated-on-send but does NOT state this specific
inversion — so C1's asymmetry paragraph is labelled this project's own
decision rather than received practice. A principle is claimed at the
strength the evidence supports, which is the opposite of a contradiction.

Second tension, also resolved in text: C1 adds a non-`omitempty` field to
the refusal record while `0005:C1` restricts it to "exactly one omitempty
structured field, `findings`". C1 argues `schema_version` is a scalar
version marker carrying no discriminator, so it does not consume the
`JDR 0001 §D10` rule-3 budget whose subject is a structured carrier. The
refusal envelope still has exactly one structured field.

## 2. Assumption Verification

Seven Critical Assumptions, all terminal (`ca=all-terminal`, 7 Verified,
0 Pending, 0 Unverified). Each record internally consistent — Status,
Method and Evidence agree and "If wrong" is non-empty.

| ID | Status | Method | Basis |
| --- | --- | --- | --- |
| A1 | Verified | Peer RDR | `0006:C17`'s closure opened by ruling; `JDR 0001 §D10` precedent |
| A2 | Verified | Source Search | no release exists — `git tag --list` empty, `internal/version::resolve` pinned `"dev"` |
| A3 | Verified | Source Search | `flow_input.go::planEnvelope` is a tolerant four-field decode, no `DisallowUnknownFields` |
| A4 | Verified | Source Search | sixteen emitted vocabularies enumerated; fourteen tiered, two deliberately untiered |
| A5 | Verified | Spike | `cli::UnknownReasons` builds clean; CLIError `code` registry is a hard import cycle — recorded `seam: none (prose-only)` |
| A6 | Verified | Spike | `clierr` is the leaf both terminal records reach; byte-identity golden re-captured |
| A7 | Verified | Spike | snapshot check is a golden-file `go test` on the existing CI `test` job, demonstrated red on promotion |

No `Docs Only` record exists, so the load-bearing-without-a-plan case
cannot arise. No `Source Search` evidence is self-referential: A2/A3/A4
resolve to `internal/version/`, `internal/cli/`, `internal/table/`,
`internal/accessor/`, `internal/graphlint/`, `internal/guard/` and
`internal/resolve/` symbols, none into this record or its artifact dir.
All 58 `source-anchor` edges resolve `true` against `main`; no bare
`file:line` anchor. No `Pending`/`Unverified` property is relied on as
settled fact anywhere in the body, because none remains.

Two assumptions are worth naming as honest negatives rather than clean
passes, and both are recorded that way in the record: A4 verified the
PARTITION but refuted C4's original assignment list as incomplete (seven of
fourteen named), which C4 now corrects; A5 refuted the CLIError `code`
enumeration seam as unbuildable at this layering, so that row carries
`seam: none (prose-only)` and its tier stands explicitly unasserted. A
tier whose seam is absent is a weaker promise, and the table says so.

## 3. Scope Verification

The Minimum Viable Validation is in scope and executed during
implementation, not deferred. It is named in the Implementation Plan's
Phase 1 and walks the whole promise end to end.

**The specific proof**: build at HEAD with `schema_version` implemented;
run `intrastate lint --model <clean-model> --as=json` and assert the
envelope reports `"schema_version":"0.1"` with no findings; add a graph-lint
finding code at severity `info` that fires on that same unmodified model;
re-run and assert exit code still 0, `type` still `ok`, the new finding
present in `data.findings` with `"severity":"info"`, and the schema version
moved `"0.1"` → `"0.2"` on the minor with the major unchanged; then promote
that code to `blocking` and assert exit flips 0 → 2 and the `ok` envelope
is replaced by the bare `CLIError` (observed key set
`code,detail,findings,message,param`).

It is a real oracle, not a smoke test: step 5 is step 4's negative control
(if promotion does not differ, step 4 proved nothing), and step 1's
pre-change baseline MUST fail the `schema_version` assertion. Assertions
are by named key, never over an exact key set, since `notes`/`warnings` are
`omitempty` siblings (S1). One half is explicitly human and said so: no
in-repo definition of a release exists for a test to range over, so whether
a version bump was the right SIZE stays a review judgement, with A7's
snapshot diff supplying the mechanical trigger.

Blast radius: `clustered=false`, `impact_families=none` — this record
declares no Cluster siblings, so no `impact.md` is owed and the fence does
not ask for one. Two cross-record obligations are nonetheless named as
prerequisites of this RDR's own Phase 1, not as follow-ups: the `0006:C17`
amendment with the `want`-literal edit at `findings_0006_test.go:288-293`
(A1), and re-capture of `0023`'s byte-identity golden
`mvv-step1-default-golden.json` (A3), which a new top-level key necessarily
changes. Neither re-cuts a peer's shipped REQ — 0006's amendment is the
ruling A1 already carries, and 0023's golden is a re-capture of the same
assertion, not a change to what it asserts.

## 5. Proportionality

Right-sized; nothing to trim before locking. The record's own
Proportionality sub-section carries the full argument and is retained in
the body above the pointer; the gate's judgement on it:

The split test is seams, not labels. The four labelled contracts are
clauses of ONE durable seam — the `--as=json` terminal envelope — and the
grounding is mechanical: every clause binds the same two terminal records,
`internal/cli/respond::Success` and `internal/cli/clierr::EmitJSON`, and no
third. C1 is the load-bearing contract, C2 the vocabulary its increment
rule ranges over, C4 that vocabulary's census, C3 the release policy over
the one growth event C2's `growing` tier admits. C3 mints no mechanism of
its own — it cites `0006:C17`'s existing partition — and a contract adding
no seam and no machinery is not an independent load-bearing contract.

`rdr-write --outcome profile` emits `stopped:split-signal` here: it counts
four `**Cn**` labels against `contracts_durable=2+`. The count is right and
the inference does not transfer — the tool counts labels, the split test is
seams. The record states this rather than silently overriding it, and notes
that the remedy the stop's `surface` names (fold the clauses under one
`**C1**`) is presentation only, changing no normative word. Accepted as
recorded; not a lock blocker.

Profile `large` re-validated at this gate: `--outcome floor` returns
`floor: none` (`rule: floor-below-two` — fewer than two prior point-fixes
at the locus), and the lens row that actually ran (grounding → 3amigo →
critique, plus repeatability-lite) is `large`'s row, so the field and the
lenses agree. `--outcome repeatability` returns `emit.next: none`
(`rule: repeatability-lite-complete`); `--outcome critique` returns `none`
(`critique-large-diffed`, two passes on differing models).

One advisory accepted rather than acted on: `evidence:over-budget` at A4
(46 lines, soft cap 30). Profile is `large`, not `foundational`, so it does
not block. The judgement the check asks for: the load-bearing anchors
(`graphlint::AggregateCode`, `resolve::RefusalKinds`, `accessor::Verdicts`
and six more) are still findable in the field itself, and the balance is
the per-surface reasoning that justifies each of fourteen tier assignments
— the content the grounding sweep reads. Relocating it to `{ARTIFACT_DIR}`
would put the justification one hop from the assignment it justifies. Kept
in place; not truncated.

## Joint-decision fence

Ran `--outcome fence` over `overlap_uncited,rulings_open,clustered,
impact_families`. First run: `stopped:overlap-uncited` — three in-flight
peers (0012, 0014, 0021) shared an anchor or contract literal with no
cross-citation. Each was FIRED as a joint-decision question against the
peer's own contracts, not synced:

- **0021 — genuine joint decision, settled here.** `0021:C2` mints a
  required `schema` marker (`intrastate.graph/1`) versioning its exported
  document, and `0021:C5` embeds that document in this envelope's `data` —
  two version markers on one wire; and 0021 emits vocabularies (`--emit`,
  `graph-export-too-large`) that C4's census does not enumerate while C4
  declares an unassigned machine-readable surface a defect. The
  propose-time `clear` ("shares `--as=json` but adds no envelope field")
  was right about the envelope and wrong about the consequence — a stale
  clear, corrected in Decision Rationale. Home: `0029:C4`, which now states
  the boundary (the document versions itself, `schema_version` versions the
  envelope, later surfaces tier in the record that adds them). 0021 is an
  early Draft with every assumption Pending; it takes the reciprocal
  citation in its own flow. No peer's locked text was edited.
- **0014 — coincidental co-mention.** `0014:C1/C3` decide gate-oracle
  mechanics (the carrier reports JSON `code`; the test asserts
  `graph-lint-failed` discriminates a lint refusal); C4 here decides
  stability tiers for the same tokens. Complementary — C4's `frozen` tier
  is what guarantees the literal `0014:C3` asserts on will not move.
  `Makefile::docs` is non-normative prior art in both. Acknowledged.
- **0012 — coincidental co-mention.** `internal/resolve` and
  `internal/table` appear as package paths only; 0012's contracts add no
  refusal kind, no `Block` member, no table category, so no vocabulary C4
  tiers is touched.

Re-run after the citations: 0 uncited pairs on both intersects, fence
resolves `op = none` (`rule: fence-clear`). `rulings_open=0`.
`joint_check_home=homed` — JC1's home edge
`0029:JC1 → cli/0029:§normative-contracts` resolves.

## Verdict

**READY.** Implementation may begin. The record specifies what to build
(four labelled contracts over one seam), how to prove it (an MVV with an
oracle table and negative controls), and what to do when a step fails
(Failure Modes names the recovery path per class). Every assumption is
terminal, including the two that came back refuted and are recorded as
such rather than smoothed over — A4's incomplete census, now corrected,
and A5's unbuildable seam, now legible as `seam: none (prose-only)`. The
two cross-record obligations (0006's C17 amendment, 0023's golden
re-capture) are prerequisites inside Phase 1 with visible failure modes,
not deferred coordination.
