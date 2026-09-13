# Finalization Gate — RDR 0021, lint's normalized-graph export

- **Record**: `cli/0021` (`0021-lint-normalized-graph-export`)
- **Date**: 2026-09-12
- **Verdict**: READY — Gate PASS, locked to Final.

Mechanical pre-sweep: `evidence/tooling-pass/iter-2/tooling-pass.md` —
PASS (`rdr lint --locking` exit 0, `blocking=0 resolution=0`). Item 4,
Cross-Cutting Concerns, is authored in the record at
`0021:G-cross-cutting` and is deliberately not copied here.

## 1. Contradiction Check

No contradictions found between research findings, design principles,
and proposed solution. The three findings that bear on the design each
license the clause that relies on them:

- Findings read `0005:C1` as carving `lint`/`dump`/`parse`-class
  command groups out of the flow contract, "owned by the RDR that
  names them". C1 claims a root `graph` verb under exactly that
  carve-out, and A4 (Method: Peer RDR) independently verifies no 0005
  envelope amendment is needed. Consistent.
- Findings read 0002 as fixing the dump's field list and row order
  while explicitly declining to define a dump grammar, naming a
  re-readable form as follow-up. C2 owns the JSON schema on that
  seeded ground and reuses 0002's field vocabulary and row order
  rather than minting a second row contract — and amends no 0002
  text. Consistent, and C2 closes
  `0002:§round-trip-inverse-invariants`'s lossy set-literal rendering
  for this document only, claiming no export→load inverse (RT3).
- Findings record a negative corpus result: no opened peer wraps DOT
  in a JSON envelope. C5 does not claim peer precedent for that cell;
  it defines it normatively as `data.dot`, a single required string
  member, with the documented unwrap `jq -r .data.dot`. The negative
  result is carried as a constraint, not overridden silently.

Principles vs planned features: the stated neutrality principle ("an
export is never a lint pass") is not merely asserted — C4 makes it
normative and the MVV step 5 gives it a mechanism-independent oracle,
byte-comparing lint's refusal with and without the export present.
C2's "no verdict or finding field, and no per-node terminal marking"
is the same principle applied to the document surface, with the
dead-end quantifier left to RDR 0015 (JDR 0001 §JD-23) rather than
re-decided here.

## 2. Assumption Verification

All eight Critical Assumptions are internally consistent and terminal.

- **Status**: 8/8 `Verified`; zero `Pending`, `Unverified`, or
  placeholder. No settled-fact prose anywhere in the record leans on
  an unverified property.
- **Method**: every label is drawn from the sanctioned set — `Spike`
  (A1, A2, A5, A6), `Source Search` (A3, A7, A8), `Peer RDR` (A4).
  Zero off-vocabulary members. No `Docs Only` record exists, so the
  load-bearing Docs-Only bar is vacuous rather than waived.
- **Evidence**: Status, Method, and Evidence agree on every row, and
  each "If wrong" is non-empty.
- **Self-reference**: none. The three `Source Search` rows resolve
  into the product tree (`graphlint/reach.go::Reach`, `::Node`,
  `::reach`, `guard/declaration.go::AssignmentCount`,
  `guard/product.go::Groups`, `graphlint/analysis.go::newAnalysis`,
  `graphlint/engine.go::Run`, `graphlint/taxonomy.go::CodeProductTooLarge`),
  never to this record or its artifact directory.
- **Symbol resolution**: all 22 `source-anchor` edges report
  `resolved: true` against `main` — none false, none absent, so the
  lookups genuinely ran. A5's Evidence is written
  `clierr.go:174::WriteJSONLine`; the symbol resolves, and the stale
  line component is a documented non-finding.

## 3. Scope Verification

The Minimum Viable Validation is in scope and executed during
implementation, not deferred. It is Phase-1 work, gated on one
authored fixture pair rather than on any later phase.

The specific proof: author a two-owned-state / one-terminal /
one-escape-row state-machine fixture plus a decision-table fixture,
then assert four invocations — (a) `intrastate graph --model
<fixture>` run twice, stdouts BYTE-identical, parsing as JSON and
carrying every C2 field including the `reach` abstraction marker;
(b) `--emit dot | dot -Tsvg` renders and its node/edge id set equals
the JSON `reach` block; (c) `--as=json | jq .data` equals the
document value-for-value; (d) `intrastate lint` over a fixture with
blocking findings refuses byte-identically against a pre-change
capture, while `graph` over that same model succeeds with an asserted
document. That fourth invocation is the neutrality oracle, and C4
binds it mechanism-independently.

Blast radius: not applicable. This record is unclustered
(`clustered=false`, `cluster=[]`, no `impact_families`), so it
retires and renames no peer's literals and no `impact.md` projection
is owed. The joint-decision fence resolved `op = none`
(`fence-clear`): every open-peer overlap touching this record is
cross-cited, `joint_check_home=clear`, and no author ruling is
unabsorbed. JC1 records the check against the twelve open peers
0012-0020 and 0022-0024.

## 5. Proportionality

Right-sized; nothing flagged to trim before locking.

**Contract count, not word count.** C1-C5 are five labelled clauses
of ONE independent load-bearing contract: the export of the
normalized graph. C1 fixes its surface (verb, arm set, flag order),
C2 the document it emits, C3 that document's byte stability, C4 the
neutrality boundary against lint, C5 how the two `--as` modes carry
it. None is separately adoptable — a consumer cannot take the
document without the surface that emits it — so this is one seam
stated five ways, not five seams locked together. No split.
`contracts_transient=0`, so no lifespan disposition distorts the
count.

**Profile re-validated.** The Metadata field reads `large`, and that
still matches the contracts just counted: one contract, user-facing
yes (a new root verb and a documented output document), locking the
`intrastate.graph/1` field list and marker. The form is correct —
value plus one clause naming the contract, with no matrix or
provenance prose left from the template. The lens battery `large`
demands did run: grounding, 3amigo, critique (two models, differing
stamps, diff written) and repeatability-lite (run 1 stamped
`variant: lite (profile: large)`, diff written, lens complete). No
lens is stale against a re-entry — there is none.

Size: 1199 lines across 45 sections for a five-clause seam that
locks a document format other records will cite. The density sits in
C2, where each field spelling is normative and therefore load-bearing
at implementation; the Pre-Lock Mini-Checks and Decision Rationale
carry the reasoning that keeps those spellings from being re-litigated
in Phase 1. Proportionate.
