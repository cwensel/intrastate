# RDR 0017 — Stage 2 prior-art record

Date: 2026-08-28. Stage: 2-propose. Budget: ≤3 corpus queries per claim,
≤5 opened hits per claim.

## Accepted citations (in-repo domain priors — load-bearing)

1. `docs/cli-output-contract.md` §Structured findings, field table:
   "`code` | the finding's own discriminator — for a model-load failure,
   the load category slug". ⇒ The shipped contract doc already defines
   per-finding `code` as the finding's OWN discriminator; an envelope-code
   echo on an allowing gate contradicts the documented field meaning.
2. RDR 0005 `artifacts/req-list.md` [REQ-24]: "`code` = category slug,
   `locator` = file:line" for `flow-model-invalid` findings. ⇒ In-repo
   precedent that a per-finding `code` draws from a producer-local
   underscore-slug vocabulary carrying no exit group and no refusal-table
   row (`internal/cli/flow_input.go::loadFindings` populates
   `Code: string(category)`).
3. RDR 0005 `artifacts/req-list.md` [REQ-51]: "A denied or indeterminate
   gate MUST surface as flow-gate-denied or flow-gate-indeterminate with
   one findings[] entry per gate; it is never a plan and never an escape
   class." ⇒ Fixes the ENVELOPE code and the findings cardinality; it
   does not assign per-finding `code` values — the vocabulary question is
   genuinely open.
4. RDR 0005 `artifacts/req-list.md` [REQ-18]: field partition — `class`
   (and `severity`) are RDR 0006's to populate. Plus
   `internal/cli/clierr/clierr.go::Finding.Class` doc: "the declared
   failure class an escape-scoped finding carries". ⇒ Refutes the
   carry-disposition-in-`class` arm on the record.
5. `internal/cli/flow_exec.go::gateResult` (`ID`/`Result`/`Reason`,
   `Result = string(result.Verdict)`) and
   `internal/accessor/model.go::VerdictAllow|VerdictDeny|VerdictIndeterminate`
   = `"allow"|"deny"|"indeterminate"`. ⇒ The sibling verb (`flow next
   --evaluate-gates` `gates[]`) already reports an allowing gate's
   disposition structurally; the disposition vocabulary exists.

## External pass (bounded; foundational profile)

Queries run (arc, semantic):

- Q1 `StateMachineRes` limit 5: "per-finding error code vocabulary in
  structured CLI diagnostics — does each finding carry its own code
  distinct from the top-level error code" → 5 hits, none on-point
  (CLI-internals/OpenAPI index pages). Rejected: no passage answers the
  identity question.
- Q2 `DevRef` limit 5: "error report design: each item in a list of
  failures carries its own code or category, separate from the overall
  status code" → best hit `RESTful Web APIs.pdf` p.328 (problem-detail
  documents — per-problem machine-readable type distinct from the HTTP
  status code). Frame-level support only; not opened deep enough to
  quote a per-item rule, so treated as supporting context, not
  load-bearing.

⚠ no prior-art coverage for the INSTANCE question (how peer CLIs code a
NON-FAILING subject reported inside a failure list) in the local corpora;
the test-runner analogy (TAP `ok`/`not ok`; per-subject disposition
markers with a separate process exit status) is from the model prior and
is demoted to a Resolve assumption rather than leaned on.

## Rejected branches

- `class` as the disposition carrier — refuted by accepted citations 3+4
  (triage had already refuted the seed's own recommendation the same way).
- `severity` as the carrier — same REQ-18 partition objection
  (`severity` is 0006's), plus severity ranks, it does not classify.
- New `flow-gate-allowed` refusal-table row — the `flow-*` table is
  exit-group-bearing refusal surface (REQ-97/98 shape); an allow is not
  a refusal, so a table row would mint an exit group for a non-refusal
  and trip the Stage-8 "additive is not exempt" table rule for no
  consumer gain.
