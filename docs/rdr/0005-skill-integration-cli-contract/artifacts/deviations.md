# Deviations — RDR 0005 Skill Integration CLI Contract

Classified per launch.md Phase 2 rules. Every entry carries exactly one Type.

## DEV-1 — `flow next --evaluate-gates` read-accessor narrowing (Phase 0 Q1)

- **Type**: SPEC-UNDER
- **Status**: mechanical translation (derived choice — unattended run)
- **Clause**: REQ-35 fixes the invoked reader set as "exactly those readers
  serving an owned key some candidate row of the requested model requires";
  REQ-37 names `read-state` as the single exception. `flow next` takes no
  `--outcome`, so two readings exist: (a) the union over ALL model rows, or
  (b) only rows still viable under supplied `--tag` facts.
- **Reading taken**: (a) — the union over all model rows.
- **Evidence**: REQ-110's refusal-identity clause does not list caller context
  as an input to `flow-artifact-missing`. Reading (b) would make that refusal
  depend on which `--tag` facts the caller happened to supply, making one
  refusal identity context-sensitive where the identity clause fixes it.
  Reading (a) keeps the refusal identity total and caller-independent.
- **Escalate if**: the MVV fixture asserts a reader must NOT run under `next`.
  Phase 1/2 must surface that as a contradiction rather than silently reconcile.


## DEV-2 — `TestReq71And72_EmitJSONRendersAngleAndAmpersandAsThemselves` asserts the inverse of its own clause

- **Type**: TEST-FIXTURE
- **Status**: RESOLVED in Phase 3c (test corrected; implementation untouched).
  The escape list at `:188` now carries the three `\u00xx` spellings, matching
  the sibling oracle. Verified still discriminating by flipping
  `SetEscapeHTML` to `true`, which the corrected oracle catches.
- **Clause**: REQ-70/REQ-72 — "`<`, `>`, and `&` serialize as themselves and
  never as `<`, `>`, or `&`".
- **What the test does** (`internal/cli/clierr/finding_0005_test.go:188`):
  ```go
  for _, esc := range []string{`<`, `>`, `&`} {
      if strings.Contains(out, esc) { t.Errorf("EmitJSON rendered %s; `<`, `>`, and `&` MUST serialize as THEMSELVES …") }
  }
  ```
  It fails when the output contains a RAW `<`, `>`, or `&` — i.e. exactly when
  the contract is SATISFIED. The sibling oracle at `:252` (and
  `assertNoHTMLEscapes` in `flow_encoder_0005_test.go:395`) loops over
  ``[]string{`<`, `>`, `&`}``, which is the intended list; this
  one appears to have lost the `\u00` prefixes.
- **Evidence**: `clierr.WriteJSONLine` uses `json.Encoder.SetEscapeHTML(false)`
  per REQ-72 and A-8, and emits
  `{"code":"flow-tag-invalid",…,"param":"[\"a<b\",\"x&y\"]"}` — `<` and `&` as
  themselves. The test at `:252` passes against the same bytes.
- **Recommendation**: correct the literal list at `:188` to the three
  `\u00xx` escapes, matching `:252`. No implementation change is warranted —
  changing the encoder to satisfy `:188` would reintroduce the very defect
  REQ-72 removes and would break `:252` and `assertNoHTMLEscapes`.

## DEV-3 — the byte-identity oracles compare a JSON string value against its unencoded form

- **Type**: TEST-FIXTURE
- **Status**: RESOLVED in Phase 3c (tests corrected; implementation untouched).
  Both oracles now DECODE the string carrier before comparing, via a new
  `decodedStringField` helper applying `rawLiteralFor`'s technique, and each
  still sweeps the RAW line for the `\u00xx` spellings so decoding launders
  nothing.
- **Clause**: REQ-71 — "Every site that emits or compares a canonical set
  literal MUST use the same encoder, so plan-to-request copy-through and
  read-back equality are byte equality."
- **What the tests do**: `finding_0005_test.go:196` asserts
  `strings.Contains(out, `["a<b","x&y"]`)` and `:245` asserts that literal
  appears exactly twice in `EmitJSON`'s output.
- **Why it cannot hold**: the literal rides in `CLIError.Param` and
  `Finding.Literal`, both Go `string` fields. A JSON string value MUST escape
  the `"` characters its content carries (RFC 8259 §7) — no encoder setting
  changes that, and `SetEscapeHTML(false)` governs only `<`, `>`, `&`. The
  emitted bytes are `"param":"[\"a<b\",\"x&y\"]"`, so the substring
  `["a<b","x&y"]` is absent by construction while the VALUE is carried
  byte-identically. The assertion is only satisfiable if the literal rides as
  a JSON ARRAY rather than inside a string, which contradicts the field types
  the same suite fixes reflectively (`:522`: every named `Finding` field is
  `reflect.String`).
- **Evidence**: the suite's own working comparator is
  `flow_mvv_0005_test.go::rawLiteralFor`, which JSON-DECODES the string
  carrier before comparing (`if strings.HasPrefix(s, "\"") { json.Unmarshal(v,
  &unquoted); return unquoted }`) — i.e. the MVV's three-hop byte-identity
  assertion already unwraps the string, and it is satisfied by this
  implementation.
- **Recommendation**: assert through `rawLiteralFor`-style decoding of the
  carrier, or against the JSON-escaped spelling `[\"a<b\",\"x&y\"]`. The
  substantive obligation — the two carriers agree byte-for-byte and neither
  HTML-escapes — is covered by `:252`, `assertNoHTMLEscapes`, and the MVV's
  three-hop test, all of which pass.

## DEV-4 — A-12's filtered `outcomes[]` is refuted by the MVV fixture

- **Type**: SPEC-UNDER
- **Status**: mechanical translation (derived choice — the MVV settles it)
- **Clause**: REQ-40 — "flow next MUST return the legal recognized-outcome
  alphabet **for the supplied state**"; REQ-74 — "`outcomes[]` (recognized
  outcome tags — nothing else)".
- **The open choice**: ASSUMPTION A-12 read "for the supplied state" as
  FILTERING the declared alphabet down to outcomes carrying at least one
  candidate row that survives the supplied facts.
- **Reading taken**: the alphabet is the model's DECLARED `outcomes` list,
  carried through unfiltered. The state qualifier governs the CANDIDATES.
- **Evidence**: `0005:MVV` requires `flow next` over `flowMVVModel` to report
  `advance`, `hold`, AND `bail` (`flow_mvv_0005_test.go:52-58`). The fixture
  declares `bail` in `outcomes` but authors no rule responding to it, so
  under A-12 it would be filtered out and the MVV would fail. The MVV is
  normative and the assumption is not, so the assumption yields.
- **Why the reading also stands on its own**: `outcomes[]` tells a caller
  what they may legally pass to `flow resolve --outcome`. An outcome with no
  currently-viable row is still legal to request — it refuses `flow-no-match`
  rather than `flow-unmodeled-outcome`, and those are different diagnoses.
  Hiding it would leave a caller unable to discover the outcome exists at
  all, while `candidates[]` already carries the state-dependent picture.

## DEV-5 — `flowReadBackFailModel` requires the same `set-state` to both succeed and fail

- **Type**: TEST-FIXTURE
- **Status**: RESOLVED in Phase 3c (test corrected; implementation untouched).
  `TestReq104And105` now creates the artifact EMPTY via `newFlowArtifact`
  instead of routing the seed through the refusing verb, so the REQ-104/REQ-105
  oracle is reached. The seed was never load-bearing — the refusal follows the
  mutation regardless of prior artifact contents. Verified still discriminating
  by pointing the fixture's writer at a reachable read-back locator, which makes
  the `set-state` succeed and the test fail.
- **Clause**: REQ-104/REQ-105 — a post-mutation read-back that did not
  complete is `flow-write-readback-incomplete` / `-timeout` at exit 3 with a
  "may have been applied" detail.
- **The contradiction**: `TestReq104And105`
  (`flow_setstate_0005_test.go:415-421`) calls
  `seedArtifact(t, model, "status=draft")` — which drives `flow set-state`
  and `t.Fatal`s if it refuses — and then drives the SAME verb over the SAME
  model and artifact with `--write status=final`, requiring THAT to refuse.
  The two invocations differ only in the written value, and the fixture makes
  the read-back unreachable by the writer's declared `path`, which is
  identical for both. REQ-110 fixes determinism ("the same request over the
  same model revision and the same artifact contents must produce the same
  success or refusal"), so no binding can make the first succeed and the
  second fail without making the refusal depend on invocation ORDER — which
  is precisely the non-determinism REQ-110 forbids.
- **Recommendation**: seed this fixture's artifact through a model whose
  read-back DOES complete (the fixture family already has `flowMVVModel`
  shapes for this), or drop the seed entirely — the refusal under test does
  not require pre-existing state, since the write is refused after the
  mutation regardless of what the artifact held. Either change is confined to
  the test.
- **Interim resolution taken** (so the substantive clause still ships
  verified): the binding seals the artifact against read-back only when the
  writer's declared locator is unreachable, which is deterministic per
  (model, request). Under it the fixture's SEED refuses, so this one test
  cannot pass as written; every other REQ-104/REQ-105 obligation — the code,
  the exit-3 group, and the "may have been applied" detail — is exercised and
  green through `accessorFailure`'s `ClassReadBackIncomplete` arm.

## DEV-6 — the artifact on-disk format and the declared-path behaviour vocabulary

- **Type**: IMPL-DECISION
- **Status**: recorded (affects future interpretation)
- **Why recorded**: no REQ constrains either choice — every Phase 1 state
  oracle goes `set-state` -> `read-state` through the CLI — but both fix
  observable behaviour a later RDR would have to honour or deliberately
  change.

**(a) The artifact is a flat JSON object of tag key to tag value, both
strings** (`internal/cli/flowbind`). Two properties are load-bearing:

- Key PRESENCE distinguishes an empty set from a cleared key. `labels`
  holding `[]` is a present key; a cleared `labels` is absent from the
  object. REQ-66 and REQ-107 turn on exactly that, and any format encoding
  absence as an empty string would collapse them.
- Values are stored VERBATIM as the canonical strings that crossed the CLI,
  so read-back equality is byte equality (REQ-71, `0004:C12`). Re-encoding
  members at the storage layer would be a second encoder, which REQ-71
  forbids.

**(b) Gate verdicts and locator unreachability come from the accessor's
DECLARED `path`**, not from the process environment: `flow.gate.deny` denies,
`flow.gate.indeterminate` is indeterminate, a `*.unreachable` or
`*-unreachable` suffix cannot be consulted, and any other gate path allows.

This is forced by `0004:C3` / REQ-31 — nothing may discover authoritative
artifacts or behaviour ambiently — and `path` is RDR 0002's carried
per-accessor locator, the only per-accessor channel a model author controls.
It is what lets the Phase 1 fixture corpus declare a denying gate or an
unverifiable read-back without the CLI consulting anything outside the model
and the caller's `--artifact` bindings.

**Consequence a later RDR must weigh**: a real deployment's gate would
presumably consult something at `path` rather than pattern-match on it. That
is a substitution of the binding, not of the seam — `accessor.GateBinding` is
unchanged — but the `path` vocabulary above becomes reserved.

**(c) A sealed-artifact marker key** (`"\x00flow.readback-unreachable"`)
records that the last write declared an unreachable read-back, so the
verifying re-read refuses `incomplete_read` while the mutation itself stays
applied (`0004:C14`). It lives in the artifact because RDR 0004 selects the
read-back reader by artifact ROLE, not by the writer's locator, making the
artifact the only channel the two share; keeping it in process memory would
make the refusal depend on process history rather than on artifact contents,
which REQ-110 names as part of request identity. The `\x00` prefix keeps it
outside any authorable tag key namespace.

---

## Orchestrator verification of DEV-2, DEV-3, DEV-5 (Phase 2 → Phase 3 boundary)

A TEST-FIXTURE classification is how a real implementation defect can be
laundered into a "bad test", so the orchestrator verified each claim against
the test source directly rather than accepting the implementer's word.

- **DEV-2 — CONFIRMED.** `finding_0005_test.go:188` loops the RAW characters
  `` `<` ``, `` `>` ``, `` `&` `` and errors when they are present, while its own
  failure message reads "MUST serialize as THEMSELVES". The oracle is inverted:
  it fires exactly when REQ-72 is satisfied. The sibling at :250 carries the
  correct `<` / `>` / `&` list and passes against the same code.
- **DEV-3 — CONFIRMED.** `literal` is the bare `["a<b","x&y"]` (:173, :231), and
  the test asserts it appears verbatim inside `param` — a JSON *string* field,
  whose inner quotes MUST escape per RFC 8259 §7. The suite's own
  `rawLiteralFor` (`flow_mvv_0005_test.go:538`) decodes before comparing, which
  is the correct technique; these two oracles skipped that step.
- **DEV-5 — CONFIRMED.** `seedArtifact` (`flow_harness_0005_test.go:45-61`)
  seeds by invoking `flow set-state` against the SAME model the test then
  requires to fail read-back. Under REQ-110 a refusal identity is fixed by its
  inputs, so the seeding call refuses identically and `t.Fatalf` at :58 fires
  before the assertion is reached. Unsatisfiable as written.

Verdict: all three are genuine test defects; the implementation is correct on
the substantive REQs behind them. Per launch.md the implementer records a wrong
test rather than editing it — **Phase 3c is the correct owner of the fix.**

---

## DEV-7 — `revision` renders empty; RDR 0002 gains no `[model].revision` field (ADV-3)

- **Type**: SPEC-UNDER
- **Status**: mechanical translation (derived choice — the evidence settles it)
- **Gap**: Phase 3b's ADV-3 showed `flowRequest.revision()` returns
  `KernelTable().Revision`, which `internal/table/model.go:404` hardcodes to
  `m.ID`. REQ-25 forbids exactly that: "The CLI never derives, hashes, or
  synthesizes it, so a model that declares none renders `revision` empty rather
  than a CLI-invented value." 3b flagged it as possibly needing an RDR 0002
  erratum and asked for an orchestrator decision.
- **Evidence**: `grep -rn 'evision' internal/table/*.go` (non-test) returns
  exactly ONE hit — the hardcoded `Revision: m.ID` at `model.go:404`. RDR 0002's
  `[model]` schema declares no revision field anywhere, which independently
  confirms Phase 0's finding #3 ("`revision` has no carrier in RDR 0002's
  schema").
- **Decision (orchestrator)**: render `revision` EMPTY. No `[model].revision`
  field is added to RDR 0002.
  1. REQ-25's clause is conditional and this is precisely its stated condition:
     no model CAN declare a revision, so every model "declares none", so every
     model renders empty. The clause is satisfied, not strained.
  2. Adding `[model].revision` would be NEW PUBLIC SURFACE that RDR 0002 does
     not name. Under launch.md's "ADDITIVE IS NOT EXEMPT" rule that is a
     SPEC-UNDER requiring author decision, and it would ship with no REQ-N
     behind it — the red-before-green gate is blind to it. A cross-RDR schema
     extension is not this RDR's to make.
  3. This stays inside 0005: only `flowRequest.revision()` changes. RDR 0002's
     `KernelTable()` keeps `Revision: m.ID` for the KERNEL's own use, where it
     is the table-identity the resolver compares — a different consumer with a
     different meaning. 0005 simply stops forwarding it as the model revision.
- **Consequence, recorded honestly**: REQ-110 names the model revision as part
  of request identity. With revision universally empty, that component is
  constant and contributes nothing to distinguishing two revisions of one
  model. That is a real limitation of RDR 0002's current schema, not of this
  fix — and it is a genuine RDR-SEED candidate for a future record that gives
  0002 a revision carrier. It is NOT grounds to invent the field here.

---

## DEV-8 — the narrowed reader set's owned-demand is `RequiresOwned` ∪ the row's owned guard keys; A-4 is superseded (ADV-2)

- **Type**: SPEC-UNDER
- **Status**: mechanical translation (derived choice — the clause settles it)
- **Clause**: REQ-35 — the invoked read-accessor set is "exactly those readers
  serving an owned key some candidate row of the requested model **requires**".
  Reinforced by REQ-30, REQ-44, and REQ-106.
- **The gap**: REQ-35 names the demand as what a candidate row *requires* but
  does not name the FIELD that carries it. ASSUMPTION A-4 read it as
  `table.Row.RequiresOwned` alone. RDR 0002 derives that field at normalization
  from the rule's WRITE BLOCK and CLEAR LIST only — `internal/resolve::Row`'s
  own doc comment says so and adds that "the guard's key set is not required to
  be a subset of this one".
- **Why A-4 under-read the clause**: a row that GUARDS on an owned key it does
  not write cannot be decided without that key. The row REQUIRES it in exactly
  the sense REQ-35 uses. Phase 3b's ADV-2 reproduced the consequence: the
  declared reader serving such a key was never invoked, the kernel received a
  view missing the fact, and `flow resolve` refused `flow-guard-unevaluable` —
  for an invocation where the reader was declared, its artifact role WAS bound,
  and the artifact HELD the value. That is REQ-106's shape inverted: a
  decidable transition turned into a refusal the caller cannot repair, because
  nothing in it names the reader that was skipped. `flow next` reported the
  same key as an unresolved fact.
- **Reading taken**: the owned demand of a row is `Row.RequiresOwned` UNIONED
  with the keys of that row's GUARD-BLOCK atoms (`all` and `unless`) whose
  `[tags.<key>]` declaration is `provenance = "owned"`.
- **Why it stays inside 0005**: `table.Row.Atoms []Atom` already carries each
  atom's `Key` and its authored `Block` verbatim (`internal/table/model.go:170`),
  and `Model.Tags` carries provenance. No RDR 0002 change is needed and none
  was made. The narrowing reads KEYS and PROVENANCE only — never operators,
  never literals — so it evaluates nothing, invents no fact, and does not
  reimplement RDR 0003's typed operator semantics, which stay the kernel's seam
  (REQ-113). Both halves are dumpable normalized model data available before
  any evaluation, which is what CA A7 requires of the narrowing inputs.
- **Blocks: guard only, not `match`**: RDR 0007's row-verdict formula
  (`0007:C5`/`0007:C6`) names `all` and `unless` and nothing else; a `match`
  atom is the kernel's selection pattern, not a guard. Recorded honestly: a
  `match` atom on an owned key a row does not write would demand a reader by
  the same argument, and would produce `flow-no-match` rather than
  `flow-guard-unevaluable`. No REQ, FM row, or failing oracle names that shape,
  and no fixture in the corpus exhibits it — every owned `match` key in the
  corpus is also a write target — so widening to `match` would be an
  unverified additive change. It is a genuine RDR-SEED candidate.
- **Consequence for A-4**: A-4's candidate-row leg (rows are selected from
  normalized model data before reads run, not from a kernel selection) stands
  unchanged. Its owned-demand leg is superseded by this entry. `req-list.md`'s
  A-4 entry carries the pointer.
