# Finalization Gate — cli/0029 version-promise-on-machine-readable-output

Date: 2026-09-12 · Profile: `large` · Verdict: **READY — Gate PASS**
Re-lock after the 0021-0029 cluster reconcile (Stage 7.1, iteration 1).

Mechanical pre-sweep: `evidence/tooling-pass/iter-2/tooling-pass.md` — PASS
(the surviving cluster re-entry note deleted after both its obligations were
discharged in live text; re-run clean). `rdr lint --locking` exit 0,
`blocking=0 resolution=0 placeholder=0 advisory=1`. Item 4, Cross-Cutting
Concerns, is authored in the record and cited as `cli/0029:G-cross-cutting`;
it is not copied here.

**What changed since the previous lock.** Two items, both from the cluster
gate's re-entry note, both now closed:

- Decision Rationale described peer 0021 as "an early Draft with every
  assumption still Pending" and justified a prose delegation with "nothing
  is edited in a peer's locked text". 0021 was Final from 2026-08-28 —
  two weeks before this record's date — so the premise was false and the
  justification inverted: because 0021's text WAS locked, the delegation
  could not be discharged in its own flow, and nothing happened. Both
  sites now state 0021's real status and cite the tracked obligation
  (0021 demoted at the cluster gate with `re-verify A4 @refine`) instead.
- The standing joint decision on two-marker read order is answered in C4.

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
  guarantee rather than restating it, and mints no new mechanism.
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
decision rather than received practice.

Second tension, also resolved in text: C1 adds a non-`omitempty` field to
the refusal record while `0005:C1` restricts it to "exactly one omitempty
structured field, `findings`". C1 argues `schema_version` is a scalar
version marker carrying no discriminator, so it does not consume the
`JDR 0001 §D10` rule-3 budget whose subject is a structured carrier.

**Re-checked at this lock.** The C4 addition introduces no contradiction:
its two mixed-case rules are compositions of clauses already on the record
(C1's unsupported-major rejection and its never-projected-into-`data`
rule) with `0021:C2`'s additive-within-`/1` promise, and it contradicts
neither. The corrected Decision Rationale paragraphs now agree with the
corpus about 0021's status, which is what the previous lock got wrong.

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
All 78 edges resolve `true`; zero `resolved: false`; no bare `file:line`
anchor. No `Pending`/`Unverified` property is relied on as settled fact.

Two assumptions are worth naming as honest negatives rather than clean
passes, and both are recorded that way in the record: A4 verified the
PARTITION but refuted C4's original assignment list as incomplete (seven of
fourteen named), which C4 now corrects; A5 refuted the CLIError `code`
enumeration seam as unbuildable at this layering, so that row carries
`seam: none (prose-only)` and its tier stands explicitly unasserted.

**Re-entry scope honoured.** The qualifier read `re-verify none` — no
assumption was disturbed by the cluster finding, which was a prose defect
about a peer's status, not about anything an assumption claims. A4's
Verified status still rests on the same census against `main`; the note's
own C-4 item records that it re-verifies for free once 0021 assigns its
tiers, and no separate repair is owed here.

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
declares no Cluster siblings in its Metadata, so no `impact.md` is owed and
the fence does not ask for one. (The 0021-0029 cluster reconcile was a
pairwise 7.1 run over two records that each declare no Cluster field; it
demoted both, and its report lives under `cluster-reconcile/0021-0029/`.)

Three cross-record obligations are named as prerequisites of this RDR's own
Phase 1, not as follow-ups: the `0006:C17` amendment with the `want`-literal
edit at `findings_0006_test.go:288-293` (A1); re-capture of `0023`'s
byte-identity golden `mvv-step1-default-golden.json` (A3), which a new
top-level key necessarily changes; and — surfaced by the cluster gate and
now tracked rather than delegated in prose — 0021's assignment of tiers to
`--emit`, `graph-export-too-large` and C2's field spellings, per `0029:C4`.

**No predicted re-cut of a peer's shipped REQ.** 0006's amendment is the
ruling A1 already carries; 0023's golden is a re-capture of the same
assertion, not a change to what it asserts; and 0021 is Draft, not shipped —
its obligation is discharged in its own open refine pass, which is exactly
what the correction to Decision Rationale now records. No Overrides entry
is owed.

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
the one growth event C2's `growing` tier admits.

`rdr-write --outcome profile` emits `stopped:split-signal` here: it counts
four `**Cn**` labels against `contracts_durable=2+`. The count is right and
the inference does not transfer — the tool counts labels, the split test is
seams. The record states this rather than silently overriding it.

Profile `large` re-validated at this gate: `--outcome floor` returns
`floor: none` (`rule: floor-below-two`), and the lens row that ran
(grounding → 3amigo → critique, plus repeatability-lite) is `large`'s row.
`--outcome repeatability` returns `emit.next: none`
(`rule: repeatability-lite-complete`).

One advisory accepted rather than acted on: `evidence:over-budget` at A4
(46 lines, soft cap 30). Profile is `large`, not `foundational`, so it does
not block. The judgement the check asks for: the load-bearing anchors
(`graphlint::AggregateCode`, `resolve::RefusalKinds`, `accessor::Verdicts`
and six more) are still findable in the field itself, and the balance is
the per-surface reasoning that justifies each of fourteen tier assignments.
Relocating it would put the justification one hop from the assignment it
justifies. Kept in place; not truncated.

**The C4 addition is proportionate.** One paragraph, no new contract label,
no new mechanism, no new seam — it states the composition of two rules the
record and its peer already carry. The alternative (a fifth contract for
marker precedence) would have split one seam's clause into its own contract
against the same split test applied above.

## Joint-decision fence

`--outcome fence` over `overlap_uncited,rulings_open,clustered,impact_families`
resolves `op = none` (`rule: fence-clear`), re-run after this pass's edits.
`rulings_open=0`. `joint_check_home=homed` — JC1's home edge
`0029:JC1 → cli/0029:§normative-contracts` resolves.

The two fired joint decisions stand where the previous lock homed them,
with one correction:

- **0022** — C3 here is the single normative home for the
  introduce-at-`info` / disclose-on-promotion rule; 0022 drops its
  restatement of the C17 closure and cites `0029:C3`. Unchanged.
- **0021** — C4 here is the single normative home for the
  envelope-versus-document boundary, now including the read order and its
  two mixed cases. What changed is the REMEDY, not the home: the previous
  lock recorded that 0021 "takes the reciprocal citation in its own flow"
  on the false premise that it was an early Draft. 0021 was Final, so it
  had no open pass in which to do that, and the delegation went untracked.
  It is now a tracked obligation on 0021's re-entry
  (`re-verify A4 @refine`), recorded as such in Decision Rationale.

`unhomed` does not apply: both fires resolve to elements in this record.

## Verdict

**READY.** Implementation may begin. The record specifies what to build
(four labelled contracts over one seam), how to prove it (an MVV with an
oracle table and negative controls), and what to do when a step fails
(Failure Modes names the recovery path per class). Every assumption is
terminal, including the two that came back refuted and are recorded as
such rather than smoothed over.

The defect that demoted this record is closed at its root, not papered
over: the delegation to 0021 now rests on 0021's actual status and on a
tracked re-entry obligation, so the thing the previous lock assumed would
happen by itself is now something a stage owes. Implementation of THIS
record does not wait on 0021 — C4's census is complete as of this record's
implementation by construction, and 0021's surfaces tier in 0021.
