# Recommendation 0014: Repository gate verification oracle

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-28
- **Status**: Draft
  <!--
  - `Deferred` is the parked-with-a-revisit-trigger status for a
    Draft that cannot proceed because **no acceptable mechanism
    exists yet** — every in-our-control path is ruled out and the
    one that would work is outside our control. It is a *pause in
    the lifecycle*, not an exit from it: the RDR stays intact and
    re-enters at the stage it stopped when the trigger fires.
    Carry the condition on the live value:
    `Deferred [revisit when <condition>]`, and say in the same
    field what was ruled out and why (Alternatives Considered
    carries the long form). Distinct from `Abandoned`, which is
    terminal — an Abandoned RDR is closed, owes a post-mortem, and
    never re-enters. A Deferred RDR owes **no** post-mortem
    (nothing was implemented), and its `Priority` records what the
    fix is *worth*, not what is scheduled. Do not defer merely to
    park work that is possible but unfunded — that is a Priority,
    not a Status.
  - `Demoted` is the terminal status for an RDR judged
    *not RDR-shaped* — the decision was never a real
    design fork, so it leaves the RDR lifecycle and is
    refiled as a plain issue. Carry the destination on the
    live value: `Demoted [→ <issue link>]`, and record the
    same link under **Related Issues**. A `Demoted` RDR runs
    no further stages. (Distinct from the 07.1 *demotion*
    below, which is a `Final → Draft` flip that keeps the
    RDR in the lifecycle — that flip never writes
    `Status: Demoted`; see the disambiguation note there.)
  - A Draft demoted from Final by the 07.1 cluster gate
    carries a qualifier on the live value:
    `Draft [revised from Final YYYY-MM-DD; re-verify A2,A4
    — <one-line reason>]`. It is still a `Draft` for every
    binary Draft/Final gate; only Stage 4 (scoped
    re-verify) and Stage 7 (re-lock) parse the qualifier.
    The Stage 7 flip to `Final` overwrites the whole value,
    so the qualifier self-clears at re-lock — no separate
    cleanup. This 07.1 "demotion" is a *verb* describing the
    Final→Draft flip; it is **not** the `Demoted` status
    above (which exits the lifecycle to an issue) — do not
    conflate the two. (`Reverted` above is the unrelated
    terminal "implementation rolled back" status — also do
    not conflate.)
  - A Final tolerated at the 07.1 gate under a JOINT-DECISION
    carries `Final [joint decision → <home §-anchor>: <the
    open question>]`. It is still a `Final` for every binary
    gate. The qualifier is an **open obligation, not a
    coherence claim**: it says the named question is
    unanswered here, not that this RDR agrees with the answer.
    So it does not self-clear. When the home answers, this RDR
    owes a scoped check of that answer against its own
    normative fences before it re-locks or implements —
    consistent → drop the qualifier and record the clearing;
    contradicts fenced text → a 07.1 SPEC-DEFECT. A re-lock
    that comes first carries the qualifier forward unchanged;
    it is never silently dropped.
  -->
- **Type**: Technical Debt
- **Profile**: foundational — provisional: a repository-wide
  verification-oracle rule (the predicate class any RDR's CI-gate
  requirement may be verified by), then applied to RDR 0006's
  REQ-119/120/121.
  <!-- Do not paste the matrix below into the field; it is the
  Stage 5 routing latch, provisional on `Draft`, made
  authoritative by Resolve.
  Sized by BLAST RADIUS — the MAX of two axes, not
  contract count or word count.
  (1) contract axis: small = one contract, no user-facing
  surface (skips Stage 5); mid = one contract + user-facing
  surface OR locks a contract; large = locks an enum/hash/
  format/grammar/destructive-op; foundational = cross-RDR
  producer / spans modules.
  (2) accretion axis (HARD floor): if `Seam Lineage` below
  carries ≥2 closed prior point-fixes at this locus, Profile
  is floored at FOUNDATIONAL regardless of the contract axis
  — a seam with prior point-fixes is never small/mid (it
  spans the prior RDRs/patches = the matrix's cross-RDR
  trigger). The only escape is a written accretion disposition
  in the Seam Lineage field. This floor is what stops a
  "one contract → mid" sizing from under-gating an accreting
  seam.
  Matrix: rdr/stages/README.md. Seed estimates from the design
  shape; Resolve overwrites from the verified count; Stage 8
  Gate locks it at Draft → Final. Never skip lenses off a
  Draft Profile until Resolve has run. -->
- **Priority**: Medium
- **Related Issues**: kata `intrastate#9en6` (1537); roborev job
  6140 findings 1/2/3
- **Predecessors**: 0006-graph-lint-authority-and-guarantees
- **Seam Lineage**: no prior accretion

## Problem Statement

A maintainer relies on the repository's graph-lint gate (RDR 0006
SC-7, REQ-119/120/121) to stop an illegal model from merging. Today
that gate is self-certified by string matching: a CI job that silently
stopped enforcing — ignored failures, validated the wrong output,
linted a different model — would not be caught, because the Go tests
assert on the gate's *text* rather than invoking the gate.

The gate has three carriers with three different oracles: the
`graph-lint` job in `.github/workflows/ci.yml`, whose first step is a
bespoke inline shell block asserting on the JSON `code` field and
whose second step also runs `make graph-lint` — so the checked-in
model is linted twice under two disagreeing oracles; the `Makefile`
`graph-lint` target (asserting on the process exit integer); and
`internal/cli/lint_gate_0006_test.go` (whole-file substring greps over
both files, plus the REQ-120 adversarial calling `runCmd` in-process,
touching neither CI nor Make). The only invocation edge is CI → Make
for the example-models step; the tests invoke neither executable
carrier.

The decision: does a gate stated as a shell/CI artifact have exactly
one executable carrier — a single script or Make target that CI
invokes and the test invokes as a subprocess, making the test's
predicate "running the gate over a defective model exits non-zero" —
or do the carriers stay plural, with the test's predicate legitimately
being artifact shape, in which case REQ-119/120/121 must be restated
to claim shape rather than enforcement? The single-carrier arm also
forces which oracle the carrier carries (the JSON `code` field vs the
exit integer, which disagree today) and prices a subprocess `make`
invocation into the Go test suite. Decide it as a repository-wide
verification-oracle rule — the predicate class any RDR's CI-gate
requirement may be verified by — then apply it to REQ-119/120/121,
rather than settling it for graph-lint alone.

## Critical Assumptions

- **A1 A Go test can invoke the carrier as a subprocess — `make
  graph-lint` from the repo root — and observe a non-zero exit, with
  the `build` prerequisite staying test-tolerable via the Go build
  cache, and concurrent invocation from parallel test packages either
  safe or explicitly serialized (no clobbered-binary flake class).**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: pending — spike: a throwaway test invoking `make
    graph-lint MODEL=<broken>.toml` via `os/exec`, timed cold and
    warm, in a repo-root working directory resolved as
    `lint_gate_0006_test.go::repoRootFor` already does; plus two
    simultaneous invocations to expose the build-output race.
  - **If wrong**: the behavioral predicate cannot live in `go test`
    (or flakes and invites a quarantine skip — the gap reborn);
    enforcement proof falls back to a separate CI step plus the
    in-process adversarial, and C1's application to REQ-120 weakens to
    the status quo.
- **A2 The bespoke ci.yml step passes (exit 0) when `intrastate lint`
  dies emitting no JSON: the `code=$(… | python3 …)` extraction fails,
  `code` is empty, and with no `set -e` the step falls through to
  success.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: pending — spike: run the step's shell verbatim with
    the binary substituted by `false` (and by a stub emitting
    non-JSON); record the step's exit status.
  - **If wrong**: the fail-open rationale for making the exit integer
    the carrier's pass/fail floor drops to duplication-removal alone;
    C3's oracle clause is re-argued at Resolve, not silently kept.
- **A3 Every refusal class the gate must catch surfaces as a non-zero
  exit carrying a JSON `code`: blocking lint aggregates to
  `graph-lint-failed` (`internal/graphlint/taxonomy.go::AggregateCode`)
  and load/env failures carry their `clierr` codes, so the exit-code
  ambiguity (`internal/cli/clierr/clierr.go::ExitCodeFor` maps
  `GroupUserEnv` and `GroupInternal` both to 2) is resolved by the
  reported `code`, whose writer is the failure envelope
  (`internal/cli/clierr/clierr.go::CLIError.Code`, emitted through the
  respond gateway per `docs/cli-output-contract.md`).**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: pending — Resolve walks the lint verb's failure
    paths (loader refusal, schema refusal, blocking findings) to the
    envelope and confirms each emits a `code` on stdout/stderr the
    carrier can extract.
  - **If wrong**: the carrier reports a missing or wrong diagnosis for
    some refusal class and CI logs misattribute the failure; the gate
    still fails closed on the exit integer.
- **A4 `make graph-lint MODEL=<path>` lints the overridden model
  first, and make aborts on its non-zero exit before the
  example-models loop runs.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: pending — spike: run the override against a broken
    model and confirm the exit is non-zero and no example-model lint
    line follows (order/emission claim — never quote-confirmable).
  - **If wrong**: the adversarial subject needs a dedicated entry
    point (a variable or sub-target on the same carrier); only the
    test's invocation shape changes, not the rule.
- **A5 No lint outcome the gate must block emits a JSON `code` with
  exit 0: advisory findings (e.g. `graph-coverage-closed-by-escape`,
  which the `Makefile` comment records as "do not fail") exit 0
  without a failure `code`, and every blocking class exits non-zero —
  so deleting ci.yml's inline JSON-code step narrows nothing.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: pending — Resolve enumerates the lint verb's
    exit/code matrix (blocking, advisory, load-refusal, env-failure
    classes) from the disposition table and envelope emission, backed
    by a fixture run per class.
  - **If wrong**: an exit-0-with-code outcome exists that the old
    inline step blocked; the carrier must then additionally fail on
    any emitted failure `code`, and C3's oracle clause is re-cut
    before the inline step may be deleted.

## Proposed Solution

### Approach

Adopt a repository-wide **single-carrier, behavioral-oracle rule** and
apply it to the graph-lint gate. The rule: a requirement claiming a
repository gate *enforces* something is verified by invoking the
gate's one executable carrier over a defective subject and asserting
refusal — never by grepping the carrier's text; text assertions are
legal only for *wiring* claims (that a named invoker calls the
carrier), and only on the invocation token. Applied here:
`Makefile::graph-lint` becomes the gate's sole executable carrier,
with the JSON `code` diagnosis folded into it; ci.yml's `graph-lint`
job reduces to `make build` + `make graph-lint` (the bespoke inline
shell block is deleted); and `lint_gate_0006_test.go`'s grep
predicates are replaced by (a) a subprocess invocation of the carrier
over a defective model asserting non-zero exit naming
`graph-lint-failed`, and (b) one wiring assertion that ci.yml invokes
`make graph-lint`. This is the repo's own stated convention — "CI
invokes the same targets so local and remote runs share one source of
truth" (`Makefile` header, lines 1–2) ⇒ the bespoke CI step is
already a convention violation, not a design to preserve.

RDR 0006 stays unamended: REQ-119's oracle clause ("the assertion
reading the JSON `code` field, not the exit integer alone", 0006:A5)
is *restated by this RDR* — the pass/fail floor becomes the carrier's
exit integer (fail-closed even when no JSON is emitted), and the
`code` field becomes the mandatory diagnosis reported on failure,
preserving 0006:A5's reason for reading it: `ExitCodeFor` maps
`GroupUserEnv` and `GroupInternal` both to exit 2, so only the `code`
distinguishes refusal classes (`internal/cli/clierr/clierr.go::
ExitCodeFor` ⇒ the carrier must surface `code`, not just fail).

### Technical Design

One executable carrier, three invokers. The carrier
(`Makefile::graph-lint`) owns the gate's logic: build the production
binary (prerequisite edge already present), lint the checked-in model
and the example models, fail closed on any non-zero lint exit, and
report the JSON `code` on failure. CI's `graph-lint` job, local `make
check`, and the Go gate test are pure invokers — none re-implements
any part of the oracle. The gate test is the *verifier*: it invokes
the carrier over a defective model (the same `[initial]`-drop defect
REQ-120 uses today) and asserts refusal; a carrier that silently
stops enforcing turns that test red, which is exactly the
self-certification gap this RDR closes. If the recipe outgrows Make,
the carrier may delegate to a script it owns — the target remains the
single entry point, so no second carrier appears.

#### Normative Contracts

**C1** — the verification-oracle rule (repository-wide)

```normative
A requirement claiming a repository gate ENFORCES a property is
verified behaviorally: invoke the gate's executable carrier over a
subject violating the property and assert refusal — a non-zero
process exit, plus the refusal `code` where the carrier reports one.
A text/shape assertion over a gate artifact may verify only a WIRING
claim — that a named invoker calls the carrier — and is limited to
the invocation token (e.g. `make graph-lint`), never the gate's
logic, subject, or oracle.
```

**C2** — the single-carrier rule (repository-wide)

```normative
A gate stated as a shell/CI artifact has exactly ONE executable
carrier. CI jobs, local entry points, and tests reach the gate only
by invoking that carrier. Gate logic inlined in workflow YAML is a
defect. In this repository the carrier class is a Makefile target; a
carrier may delegate to a script it owns, and that script is not a
second carrier.
```

Application to the graph-lint gate — restates RDR 0006
REQ-119/120/121's verification, without amending 0006:

**C3**

```normative
`Makefile::graph-lint` is the sole carrier of the RDR 0006 gate.
- Oracle: the carrier fails on any non-zero `intrastate lint` exit
  (fail-closed floor) and on failure reports the JSON `code` field
  as diagnosis. This restates 0006:A5's "JSON `code`, never the exit
  integer alone": the code moves from CI's inline assertion into the
  carrier's report, preserving refusal-class discrimination.
- ci.yml's `graph-lint` job reduces to `make build` +
  `make graph-lint`; the bespoke inline JSON step is deleted only
  after the exit/code matrix (A5) confirms deletion narrows nothing.
- The gate test's predicates become: (a) BEHAVIORAL — invoking the
  carrier as a subprocess over a defective checked-in-model variant
  exits non-zero AND names `graph-lint-failed` (both conjuncts: the
  code discriminates a lint refusal from a wrong-reason failure such
  as a broken build or missing file), and the predicate MUST exercise
  the invocation CI performs — the default no-override arm, with the
  defective subject substituted in an isolated copy — not only a
  `MODEL=` override arm; (b) WIRING — ci.yml's `graph-lint` job
  invokes `make graph-lint` as a LIVE, GATING step (not commented,
  no `continue-on-error`, no disqualifying `if:` guard) — the one
  permitted artifact assertion under C1, checking gating-ness, not
  token presence.
```

#### Load-Bearing Decisions

- **Naming** — the carrier and the CI job keep the name `graph-lint`
  (the `lint` name stays golangci-lint's, settled by 0006:A5;
  rejected: `lint-gate`, `model-lint` — renaming would break the
  wiring token every invoker and the wiring assertion share).
- **Selection / predicate** — among the candidate carrier homes
  (inline workflow shell, a `scripts/` script, the Make target), the
  **Make target** is canonical because the repo's stated convention
  routes CI through Make targets (`Makefile` header: "CI invokes the
  same targets so local and remote runs share one source of truth");
  a script, if the recipe outgrows Make, is owned *by* the target and
  invoked only through it.

#### Illustrative Code

Illustrative only — exact recipes and signatures sharpen at
Resolve/Pre-Lock; tests must not assert these literally.

```make
# carrier shape: fail closed on exit, report code on failure
graph-lint: build
	run $(BIN) lint --model <subject> --as=json; \
	on non-zero: extract and echo the JSON "code"; exit non-zero
```

```go
// verifier shape: behavioral predicate over the carrier
out, err := exec.Command("make", "graph-lint",
	"MODEL="+brokenModelPath).CombinedOutput()
// want: err != nil, and out names graph-lint-failed
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Gate execution | `Makefile::graph-lint` | Exit-integer oracle only; no `code` report | Extend | Becomes the sole carrier; gains the `code` diagnosis (C3) |
| CI enforcement | `.github/workflows/ci.yml` `graph-lint` job | Bespoke inline oracle duplicates the target | Replace (step) | Job reduces to `make build` + `make graph-lint` |
| Gate verification | `internal/cli/lint_gate_0006_test.go` | Whole-file substring greps; in-process adversarial only | Replace (predicates) | Behavioral subprocess predicate + one wiring assertion |
| Defective subject | `lint_gate_0006_test.go::dropInitialTable` | none | Reuse | Same minimal defect feeds the subprocess predicate |
| Refusal taxonomy | `internal/graphlint/taxonomy.go::AggregateCode` | none | Reuse | The behavioral predicate's expected code |

### Decision Rationale

Scored matrix (approaches × deciding criteria; A = single carrier +
behavioral oracle, B = plural carriers + shape-restated REQs, C =
workflow-level e2e verification, D = in-process-only verification):

| Criterion | A: single carrier | B: plural + shape REQs | C: workflow e2e (act/canary) | D: in-process only |
| --- | --- | --- | --- | --- |
| Correctness fit (catches a gate that stopped enforcing) | catches carrier neutering and wiring drift | catches nothing new; legitimizes the gap | catches everything incl. runner config | catches engine regressions only |
| Prior-art alignment | matches the repo's Makefile convention and the peer pattern (beads/helm/roborev) | contradicts the repo's stated convention | no precedent in the peer set | partial (REQ-120 already does it) |
| Reversibility | high — greps restorable | high | low — new dependency/infrastructure | high |
| Blast radius | Makefile + ci.yml + one test file | REQ text only | CI infrastructure + new dependency | one test file |
| Cost | subprocess `make`+build priced into `go test` (A1) | zero | act install or canary upkeep; flaky | zero |

The correctness-fit and prior-art rows decide it: B and D cannot
catch a silently-stopped-enforcing gate — the exact defect roborev
job 6140 raised — so choosing them re-labels the problem instead of
closing it; C is the only approach that also verifies the GitHub
runner itself, but at a standing cost and flakiness the wiring
assertion plus a real-runner activation run covers well enough. A
wins with bounded cost (one subprocess invocation, A1) and is the
only arm consistent with the repo's own single-source-of-truth
convention. Rejections, one line each: B — restating REQs to shape
concedes the gate is unverified; C — nektos/act or a canary branch is
a heavy, flaky dependency, and its standing coverage of
workflow-config semantics is instead taken by C3(b)'s gating-aware
wiring assertion plus the activation checklist (required-checks
membership); D — leaves CI/Make free to drift, the raised finding.

Premortem: hardened (hardened) — critic ledger P-1…P-12
(`evidence/propose-premortem/critic.md`) folded: default-arm
predicate and wrong-reason conjunction into C3(a), gating-aware
wiring assertion into C3(b), exit/code matrix as A5, subprocess
concurrency into A1, shell-semantics disposition into Risks and
Load-Bearing Decisions; no finding forced a switch.

Ground-sweep: clean (19 anchors)

Joint-check: clear (12 peers) — open peers 0012, 0013, 0015–0024;
every `graph-lint` hit is the Predecessor slug
`0006-graph-lint-authority-and-guarantees` or an engine mention,
0019's `lint_gate_0006_test.go` citations are source-search evidence
(it neither modifies the file nor the missing-`[initial]` refusal),
and `models/rdr.toml` peers cite it as lint subject, which this RDR
does not modify; absence arm N/A (no refusal converted to
acceptance); bridge sub-check N/A (no Cluster membership).

## Alternatives Considered

### Alternative 1: Plural carriers, REQs restated to shape

**Description**: Keep the three carriers and their three oracles;
restate REQ-119/120/121 as artifact-*shape* claims so the existing
substring greps become honest predicates. The repository rule would
read: CI-gate requirements are verified by artifact shape plus an
in-process behavioral adversarial.

**Pros**:

- Zero code motion; no new cost in the test suite.
- Honest about what the current tests actually prove.

**Cons**:

- Concedes the raised defect permanently: a CI job that silently
  stopped enforcing still merges illegal models undetected.
- Contradicts the repo's own convention (`Makefile` header) and every
  peer-instance read; keeps the double-lint of the checked-in model
  under two disagreeing oracles.

**Reason for rejection**: fails the matrix's correctness-fit row —
it re-labels the self-certification gap instead of closing it.

### Alternative 2: Workflow-level end-to-end verification

**Description**: Verify the CI wiring itself — run the actual
`graph-lint` job locally under nektos/act in a test, or maintain a
scheduled canary (a branch carrying a deliberately broken model)
asserting the job goes red on a real runner.

**Pros**:

- The only arm that also verifies runner configuration and workflow
  syntax — nothing is left to a wiring assertion.

**Cons**:

- Heavy standing dependency (act or canary upkeep), slow and flaky;
  not runnable inside plain `go test`.
- No precedent in the peer set (bounded sweep found none).

**Reason for rejection**: the marginal coverage over A is the runner
itself, which the one-time real-runner activation run (Operational
Activation) plus the wiring assertion buys without a standing cost.

### Briefly Rejected

- **In-process-only verification**: keep REQ-120's `runCmd`
  adversarial as the sole behavioral proof and drop the greps — leaves
  CI and Make free to drift, which is exactly the raised finding.
- **Script as the carrier instead of the Make target**: what is
  rejected is a second *entry point*, not script-housed logic — a
  `set -euo pipefail` script owned and invoked solely by the target
  is the sanctioned home for recipe logic that outgrows Make (see
  Load-Bearing Decisions and the shell-semantics risk); a
  freestanding `scripts/` entry point would invert the repo's
  Make-first convention and re-pluralize the carrier.

## Context

### Background

Raised by roborev job 6140 (findings 1/2/3) against the RDR 0006
implementation; tracked as kata `intrastate#9en6`
(odc-type: test-oracle, odc-trigger: design-conformance). RDR 0006's
deviation D4 discharged the *existence* of the CI job, the Make
target, and the checked-in model — it says nothing about oracle
*fidelity*, so this is beyond D4's disposition. Scope review judged
this RDR-shaped: closing it honestly is a gate-architecture or
REQ-claim decision, not a test edit. (At current HEAD one invocation
edge exists — ci.yml's example-models step runs `make graph-lint`,
added by `39aeb02` — which duplicates rather than replaces the bespoke
step; the test still invokes no executable carrier.) Distinct from
kata `5mhf` (which model is the gate's subject, vs how the gate is
verified).

### Technical Environment

Go module `github.com/cwensel/intrastate`. Carriers:
`.github/workflows/ci.yml` (`graph-lint` job), `Makefile`
(`graph-lint: build` target), `internal/cli/lint_gate_0006_test.go`
(REQ-119/120/121 tests, in-process `runCmd`). Governing record:
RDR 0006 (SC-7, REQ-119/120/121, deviation D4).

## Research Findings

### Investigation

Read before enumerating (evidence:
`evidence/research/prior-art.md`): the three carriers at HEAD
(`.github/workflows/ci.yml` `graph-lint` job, `Makefile::graph-lint`,
`internal/cli/lint_gate_0006_test.go`); RDR 0006's gate decision
(0006:A5, SC-7) and its D4 disposition; the repo's conventions
(`Makefile` header, `docs/cli-output-contract.md`); a peer-instance
pass over the langref checkout set (gh-cli, helm, goreleaser,
roborev, beads) asking how each verifies its own repository gates;
and a bounded DevRef corpus pass (3 queries) on CI-scripts-as-single-
carrier. Sibling-path check (exhibited): the behavioral-oracle
discriminator already exists in adjacent gates —
`Makefile::release-check` runs the built binary and asserts the
stamp width (its comment records that a text-presence oracle was
tried and shown unsound), `Makefile::docs-check` is
regenerate-then-diff, and
`lint_gate_0006_test.go::TestReq120_ADeliberatelyIllegalEditToTheCheckedInModelFailsTheGate`
is already behavioral in-process ⇒ C1 generalizes an existing in-repo
signal rather than inventing a parallel one.

### Key Discoveries

- **Documented** — `Makefile` header: "CI invokes the same targets so
  local and remote runs share one source of truth" ⇒ the bespoke CI
  step is a standing violation of the repo's own convention; the
  single-carrier arm restores it rather than introducing policy.
- **Documented** — peer instance: beads'
  `scripts/ci_capability_selector_test.go` subprocess-invokes the
  exact script its CI runs (`exec.Command("bash", …
  "ci-capability-selector.sh")`); helm and roborev CI run bare `make`
  targets ⇒ the chosen predicate class is the peer norm.
- **Documented** — a bounded sweep found no peer Go test asserting on
  its own CI workflow text ⇒ the current grep oracle has no
  precedent in the comparison set.
- **Documented** — the checked-in model is linted twice in CI today:
  the bespoke step (JSON `code` oracle) and `make graph-lint` (exit
  integer, via the example-models step) ⇒ consolidation deletes a
  duplicate *invocation*; that it deletes no unique *blocking power*
  is exactly A5, proven before the step is removed, not assumed.
- **Documented** — `internal/cli/clierr/clierr.go::ExitCodeFor` maps
  `GroupUserEnv` and `GroupInternal` both to exit 2; the
  discriminating `code` is written by
  `internal/cli/clierr/clierr.go::CLIError.Code` through the respond
  gateway ⇒ the carrier must report `code` as diagnosis even with an
  exit-integer pass/fail floor (0006:A5's rationale preserved).
- **Assumed** (A2, spike at Resolve) — the bespoke step's `code`
  extraction fails open when the binary crashes emitting no JSON ⇒ if
  confirmed, the exit integer is the only fail-closed floor.
- **Assumed** (A1/A4, spikes at Resolve) — subprocess `make
  graph-lint MODEL=<broken>` is test-tolerable and aborts before the
  example-models loop.
- ⚠ no prior-art coverage found for the *rule-as-policy* form (a
  repository-wide written verification-oracle rule); the rule's
  scope, beyond its graph-lint application, is from the model prior
  grounded in the in-repo convention.

## Trade-offs

### Consequences

- Positive: one oracle, one carrier — CI/local/test drift at this
  gate becomes structurally impossible, and the gate test verifies
  enforcement instead of text.
- Positive: the duplicate lint of the checked-in model in CI
  disappears, and with it the bespoke step's python3 dependency.
- Negative: the Go test suite gains a subprocess dependency on `make`
  and a `go build` via the carrier's `build` edge (priced by A1);
  environments without `make` cannot run the gate test.
- Negative: the oracle's logic moves out of ci.yml into the Makefile
  — workflow reviewers see only an invocation line.
- C1/C2 bind future RDRs' CI-gate requirements: their tests must
  invoke carriers, which is a (deliberate) constraint on how cheap a
  future gate test can be.

### Risks and Mitigations

- **Risk**: the subprocess predicate is slow or flaky in CI (cold
  build, parallel `go test` packages both invoking make).
  **Mitigation**: A1's spike times it before lock; the build cache
  makes rebuilds incremental, and the predicate can serialize behind
  a package-level lock if needed (Resolve decides).
- **Risk**: restating REQ-119's oracle is read as amending Final RDR
  0006.
  **Mitigation**: 0006 is untouched; C3 records the restatement in
  this RDR with 0006:A5 cited, the successor-record pattern already
  used across 0001–0011.
- **Risk**: the carrier's `code` extraction re-introduces a fail-open
  path inside Make.
  **Mitigation**: C3 fixes the pass/fail floor to the exit integer;
  the `code` is diagnosis only, so a broken extraction can lose the
  message but never the failure.
- **Risk**: Make/shell semantics diverge across the carrier's two
  runtimes (macOS GNU Make 3.81 + bash-as-sh locally, Make 4.x + dash
  on ubuntu runners) — per-line recipe shells and loop status-masking
  make the "single carrier" environmentally plural (critic P-2/P-6).
  **Mitigation**: the recipe stays single-command-simple, or its
  logic moves into a `set -euo pipefail` script the target owns and
  invokes (the Load-Bearing Decisions escape hatch) — carrier
  identity stays the Make target either way; Resolve picks when the
  recipe crosses that line.

### Failure Modes

- Carrier neutered (recipe edited to a no-op): visible — the
  behavioral gate test invokes the carrier over a defective model and
  goes red when refusal doesn't happen; this was the silent case
  under the grep oracle.
- Carrier neutered in one arm only (a conditional on `MODEL` that
  enforces under the test's override but no-ops on CI's default
  invocation): visible — C3(a) requires the predicate to exercise the
  default arm CI actually runs (critic P-3).
- Wrong-reason non-zero exit (broken build, missing file) read as
  refusal: prevented — C3(a)'s predicate is the conjunction of
  non-zero exit AND `graph-lint-failed`, so a build breakage turns
  the gate test red for the stated wrong reason instead of silently
  green (critic P-7/P-11).
- Vacuous pass (the subject path drifts and the carrier lints
  nothing, exiting 0): the REQ-121 model-conformance test keeps
  `models/rdr.toml` pinned and loadable, and the carrier names each
  subject it lints; the vacuity-kill assertion's exact shape is
  settled at Resolve (critic P-7).
- CI invocation present but non-gating (`continue-on-error:`, a
  disqualifying `if:`, commented out): C3(b)'s wiring assertion
  checks gating-ness, not token presence; required-status-checks
  membership is confirmed at Operational Activation (critic P-4/P-9).
- `make` or the Go toolchain missing where the gate test runs, or
  repo-root discovery failing after a package move: must fail loudly,
  never skip silently — a skipped enforcement test is the
  self-certification gap reborn; the skip/fail policy per environment
  is settled at Resolve (critic P-12).
- Diagnosis path: the test relays the carrier's combined output, so a
  red run shows the lint envelope including the `code`; locally,
  `make graph-lint MODEL=<path>` reproduces the exact gate.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1–A4)

### Minimum Viable Validation

1. Copy `models/rdr.toml`, drop its `[initial]` table (the REQ-120
   defect); run `make graph-lint MODEL=<copy>` → non-zero exit,
   output names `graph-lint-failed`.
2. Run the converted gate test on HEAD → green.
3. Neuter the carrier (stub the `graph-lint` recipe to a no-op) → the
   same test goes red. This step is the RDR's point: under the old
   grep oracle it stayed green.
4. Restore the carrier; confirm `ci.yml`'s `graph-lint` job contains
   only `make build` + `make graph-lint` (bespoke step gone) and
   `make check` still reaches the gate through the same target.

### Phase 1: Carrier Consolidation

Fold the JSON-`code` diagnosis into `Makefile::graph-lint` (C3
oracle) and reduce ci.yml's `graph-lint` job to `make build` +
`make graph-lint`.

### Phase 2: Test Conversion

Replace `lint_gate_0006_test.go`'s substring-grep predicates with the
behavioral subprocess predicate over the carrier — covering CI's
default no-override arm as well as the override arm, plus the
wrong-reason discrimination case — and the single gating-aware wiring
assertion; keep the defective-subject helper and the REQ-121
model-conformance check.

### Phase 3: Convention Codification

Point future gate authors at C1/C2 from the contributor-facing
conventions (one line in `CONTRIBUTING.md`), so the rule outlives
this application.

### Operational Activation

Push and watch the `graph-lint` job run green on a real runner —
never declare the CI change done from local verification alone;
exercise one scratch-branch run with a deliberately broken model to
watch the job fail for the stated reason; and confirm the job's
membership in the branch-protection required status checks (the one
wiring fact no in-repo assertion can see — critic P-9).

## Validation

### Testing Strategy

[Required — never omit. Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

### Performance Expectations

[Conditional — omit (don't N/A-bullet) this section unless
comparing alternatives on empirical performance grounds.
Do not include effort estimates or speculative
throughput targets. Rough performance metrics are
appropriate only when comparing alternatives — note
empirical data or obvious gains that support the
chosen approach over a rejected one.]

## Finalization Gate

> Complete each item with a written response in
> `{ARTIFACT_DIR}/gate.md` before marking this RDR as
> **Final**. Written responses prevent rubber-stamping
> and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses.
>
> At lock, replace Contradiction Check, Assumption
> Verification, Scope Verification and Proportionality
> with the one-line pointer to gate.md — those four
> judge THIS record at THIS lock and no peer cites
> them. **Cross-Cutting Concerns stays here**, below
> the pointer: it names the project-wide policy other
> RDRs conform to, so it must stay projected and
> citable as `cli/NNNN:G-cross-cutting`. Cite it that
> way, not by section name.

### Contradiction Check

[Gate key: contradiction — a gate response is cited as
`cli/NNNN:G-<key>`, so the key is a stable id and is
not derived from this heading, which may be reworded.]

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

[Gate key: assumptions]

[Confirm every Critical Assumption Evidence Record
is internally consistent: Status, Method, and
Evidence agree, and "If wrong" is non-empty. List
any record whose Method is `Docs Only` (these block
lock unless paired with a Spike or Source Search
plan) and any that remain `Pending` or `Unverified`
with a plan to verify before implementation begins.
Confirm no `Verified` stamp is self-referential or
proves only an adjacent claim, and that each cited
`path::Symbol` resolves on `main`. **Status
consistency:** no assumption marked `Pending` or
`Unverified` may have settled-fact prose elsewhere in
the RDR depending on it.]

### Scope Verification

[Gate key: scope]

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[Gate key: cross-cutting]

[Retained at lock — this sub-section stays in the RDR
when the other gate responses move to gate.md, because
peer RDRs cite it as `cli/NNNN:G-cross-cutting` and an
element that is not projected cannot be cited.]

[List only concerns that apply to this RDR. For each,
state either how this RDR addresses it, or which peer
RDR owns the project-wide policy this RDR conforms
to. Omit (rather than N/A-bullet) anything that does
not apply.]

Candidate concerns (include only those that apply):
versioning · build tool compatibility · licensing ·
deployment model · IDE compatibility · incremental
adoption · secret/credential lifecycle · memory
management · concurrency model · character encoding ·
canonical-form / determinism (see note below).

If this RDR claims byte-identical output,
content-addressed identity, or replay-stable hashes,
also confirm: hash function + library, pre-image
byte layout, primitive encodings, map iteration order,
whitespace policy, case folding, empty/null/absent
distinguishability, and a version marker for future
evolution.

### Proportionality

[Gate key: proportionality]

[Is the document right-sized for the change? Flag
any sections that should be trimmed before locking.
The split test is **contract count, not word count**:
confirm this RDR is the sole author of at most one
independent load-bearing contract (per the Normative
Contracts split signal). If it owns more than one
seam, flag it for splitting rather than locking the
seams together.

Re-validate the **Profile** Metadata field against the
contracts you just counted: confirm the value Resolve
wrote still matches (one contract + no user-facing
surface → `small`; etc. per the applicability matrix).
If the lenses that actually ran disagree with the
Profile (e.g. Profile says `small` but the change locks
a contract that warranted `mid`+ lenses, or the lenses
were skipped on a wrong `small`), correct the field and
do not lock until the missing lenses have run. This is
the latch's backstop — a wrong Profile cannot route
past the lens battery undetected. A `Transient`-marked
contract with a named deleting sibling and schedule is a
recorded lifespan disposition, not an under-sized
Profile — do not count it when re-deriving. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- RDR 0006 (`0006:A5`, SC-7/S7; implementation artifacts
  `req-list.md` REQ-119/120/121, `deviations.md` D4)
- `.github/workflows/ci.yml` (`graph-lint` job), `Makefile`
  (`graph-lint`, `release-check`, `docs-check`),
  `internal/cli/lint_gate_0006_test.go`,
  `internal/cli/clierr/clierr.go`, `internal/graphlint/taxonomy.go`,
  `docs/cli-output-contract.md`
- Peer-instance reads: langref checkouts (beads
  `scripts/ci_capability_selector_test.go`; helm
  `build-test.yml`; roborev `ci.yml`) — see
  `evidence/research/prior-art.md`
- Continuous Delivery (Humble & Farley) p.187; The DevOps 2.0
  Toolkit p.226 (DevRef corpus)
- kata `intrastate#9en6`; roborev job 6140 findings 1/2/3
