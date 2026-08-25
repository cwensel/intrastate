Model: claude-opus-5

# Finalization Gate — RDR cli/0005 skill-integration-cli-contract

- **RDR**: `0005-skill-integration-cli-contract`
- **Date**: 2026-08-24
- **Verdict**: **READY — Gate PASS**

Mechanical pre-sweep: `evidence/tooling-pass/tooling-pass.md` (BLOCK, 3 findings,
all MECHANICAL, fixed in-pass) → `evidence/tooling-pass/iter-2/tooling-pass.md`
(PASS).

## 1. Contradiction Check

No conflict between Research Findings and Proposed Solution.

The Investigation establishes three constraints and the solution honors all
three. (a) *No resolver CLI exists* — the plan is greenfield under
`internal/cli`, and the Existing Infrastructure Audit / reconcile reuse row both
record "nothing to fold in", consistent with `internal/cli` registering only
`newVersionCmd`. (b) *The kernel already exists and is pure* — the Approach
keeps selection in `internal/resolve` and binds reads/gates in the verb, which
is what JDR 0001 §D8 resolved; nothing in the design makes the kernel stateful.
(c) *The output gateway is the mandated path* — every verb routes through
`respond.ValidateMode` / `respond.OK` / `respond.Fail`, and the Investigation's
one drift note (`newVersionCmd` still using `cmd.Println`) is carried forward as
an explicit instruction to extend `respond.OK` rather than copy the shortcut.
That is the finding being *acted on*, not contradicted.

Planned features vs stated principles: the "navigated, not orchestrated"
principle taken from the sibling prior is honored — the four verbs answer
questions and apply one planned mutation each; no verb drives a loop. The peer
split is respected rather than absorbed: 0002 keeps the model format, 0003 the
predicate grammar, 0004 the accessor safety model, 0006 the lint invariants, and
the four cross-seam decisions live at JDR 0001 §D8–§D11 and are cited, not
restated.

One tension was checked and is not a contradiction: the RDR claims request-level
*semantic* determinism while also making one *byte*-level claim (the canonical
set literal). Cross-Cutting Concerns states exactly this split, and the byte
claim is scoped to the single encoder rule A6 pins. The two coexist by
construction.

## 2. Assumption Verification

All seven records are internally consistent; Status / Method / Evidence agree
and every "If wrong" is non-empty and concrete.

| ID | Status | Method | Consistent? |
| --- | --- | --- | --- |
| A1 | Verified | Source Search | Yes — `respond`/`clierr`/`ExecuteAndEmit` symbols all resolve on the current tree. |
| A2 | Verified | MVV Test | Yes — the named scenario is in the Validation section and covers all four verbs. |
| A3 | Verified | Peer RDR | Yes — cites JDR 0001 §D2/§D4 + `0002::Normative Contracts` + `0003::Approach`, each an element. |
| A4 | Verified | Peer RDR | Yes — JDR 0001 §D8/§D9 + `0004::Normative Contracts`; kernel-purity leg additionally confirmed by two named tests. |
| A5 | Verified | Source Search | Yes — exit mapping verified in `clierr::ExitCodeFor`; the not-yet-existing `Findings []Finding` is stated as absent, not claimed present. |
| A6 | Verified | Design Decision | Yes — a scoping choice with the rejected alternative named (leaving grammar/codes to implementation-time invention). |
| A7 | Verified | Peer RDR | Yes — both legs anchored on `0002::Normative Contracts` *Accessor tables* / *Gate references* and `Row.RequiresOwned`. |

- **No load-bearing `Docs Only`**: no `Docs Only` record exists in this RDR.
- **No `Source Search` self-reference**: A1 and A5 both cite `internal/cli/**`
  production source; neither resolves to this RDR or its artifact directory.
- **Method vocabulary**: all seven labels are among the sanctioned eight.

Two records (A5, A6) carry Evidence fields over the advisory 30-line budget
(59 and 45 lines). Both are settled-at-round content — the flat `omitempty`
`Finding` shape with its prior-art tiebreak, and the HTML-escaping-disabled pin
with its normative fixtures. The load-bearing anchors stay findable in-field, so
the fields are kept as written; C9 is advisory and does not block.

A5 and A6 each rest partly on a **Design Decision made here rather than
verified upstream** — the flat widened `Finding` record, and the set-literal
encoder. Both are correctly labeled, both name the rejected alternative, and
both are covered by the MVV (the `<`/`&` byte-identity assertion). This is the
right disposition for a decision no peer owns; it is recorded, not hidden.

## 3. Scope Verification

The Minimum Viable Validation is **in scope for implementation, not deferred**.

The named proof: one fixture-backed flow exercised through the production Cobra
path, asserting all four verbs — `flow next` (legal outcome alphabet with gate
ids as unresolved facts), `flow resolve` (owned state read from a fixture
artifact, one outcome → one plan, one gate on the selected row), `flow
read-state` (fixture artifact tags per reader), and `flow set-state` (one scalar
write, one set write as a JSON array, one `--clear`, each read-back verified).

It is genuinely load-bearing rather than a smoke test, because it carries the
assertions that discriminate this RDR's contested decisions:

- **The encoder pin (A6)** — the set write's members must include one containing
  `<` and one containing `&`, asserted byte-identical through
  plan → request → read-back. This is what moves the HTML-escaping rule from
  review opinion into a test.
- **The narrowing claim (A7)** — a declared reader that no candidate row needs,
  with its role unbound: `next`/`resolve` must succeed without invoking it,
  while `read-state` must invoke it and refuse. The RDR correctly notes this
  pair is the *only* assertion distinguishing a narrowed invoked set from "run
  every declared reader".
- **The gate-scope claim** — `next --evaluate-gates` over one gated candidate
  plus one guard-excluded row: only the reported candidate's gate runs, and a
  deny there still exits 0.
- **The error taxonomy (A5)** — one escaped plan, one gate deny, one kernel
  refusal, and one exit-3 accessor failure, in both `--as=text` and `--as=json`,
  asserting `findings[]` where the table names it.

Phases 1–4 sequence the implementation, and the MVV is not pushed past them.

## 4. Cross-Cutting Concerns

- **Versioning** — verb-specific JSON `data` payloads and the `findings` field
  are append-only under the existing CLI output envelope. This is the same
  additive convention the shipped `clierr.CLIError` doc comment already fixes
  ("keep them `omitempty` so the envelope stays append-only and stable for
  tools"), so the RDR inherits a policy rather than inventing one.
- **Incremental adoption** — `--model <path>` with fixture-backed loading ships
  before config-resolved `--flow <id>`, so the surface is testable before the
  config seam lands.
- **Secret/credential lifecycle** — not introduced here. Accessor execution and
  external availability are owned by RDR 0004; this RDR only maps their refusals
  to exit codes.
- **Canonical form / determinism** — owned here and scoped explicitly: the
  general claim is request-level semantic determinism over the same model
  revision and artifact contents, *not* byte-identical output. The single
  byte-level claim is the JDR 0001 §D13 canonical set literal surviving
  plan → request → read-back, which holds only because one encoder — HTML
  escaping disabled — renders it at every emit site. The RDR records the
  consequent obligation (both shipped emit sites, `clierr::EmitJSON` and
  `respond::writeJSONLine`, move off bare `json.Marshal` onto one shared
  helper) and records the JCS non-conformance as a deliberate accepted
  deviation rather than an oversight.
- **Error taxonomy / exit codes** — owned here by assignment: RDR 0004's
  Disposition Table explicitly defers ("`exit code` is out of scope here — RDR
  0005 owns the CLI mapping"), and JDR 0001 §D10 supplies the rule. No second
  leg competes.

Concerns that do not apply (no data migration, no persistence schema of its own,
no auth surface, no network protocol) are omitted rather than N/A-bulleted.

## 5. Proportionality

Right-sized. `Profile: mid` remains correct: the contract is user-facing but not
foundational, and Seam Lineage records no prior accretion.

The RDR is long, but the length sits where the decisions are — A5's `Finding`
shape and A6's encoder pin are the two places where this RDR had to *decide*
rather than cite, and both carry their tiebreak evidence and normative fixtures.
The peer-owned material (model format, predicate grammar, accessor safety, lint
invariants, the four JDR decisions) is cited in a line or two each rather than
restated, which is the discipline that keeps a `mid` RDR from becoming a
foundational one.

Nothing is flagged for trimming before lock. One item is deliberately *kept*
against the advisory budget: the A5/A6 Evidence mass, for the reason given in
item 2 — cutting it would blind the grounding sweep that reads those anchors.
