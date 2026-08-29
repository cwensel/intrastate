# RDR 0020 — Stage 2 prior-art research trail

Stage: Propose (2026-08-28). Budget: ≤3 corpus queries per claim, ≤5
opened hits per claim.

## Accepted citations (in-repo, quotable)

- `internal/cli/flow_input.go::parseTags` — the guarded-lookup comment:
  "a missing declaration yields the zero `TagDecl`, which conforms
  everything and leaves the shape-only behaviour a caller already
  relies on intact. Refusing an undeclared key is a different decision,
  under a different code, and this is not the place that makes it."
- `internal/table/load.go::ConformValue` — doc: "A zero TagDecl
  conforms everything: an undeclared kind has no domain to violate,
  which is what keeps an undeclared `--tag` key shape-only."
- `internal/cli/flow_input.go::canonicalValue` — the contradicting arm:
  `isSet := decl.Kind == "set"`; `!isSet && looksArray` refuses
  "the tag `<key>` is not set-valued".
- `internal/table/normalize.go` (`CatUnknownTag` sites) — "references
  the undeclared tag", "writes the undeclared tag", "clears the
  undeclared tag"; `internal/table/load.go::accessorTable` — "names the
  undeclared tag"; `[initial]` arm — "assigns the undeclared tag".
  ⇒ every model-side reference to an undeclared tag refuses at load;
  an undeclared `--tag` key is structurally unreadable by any rule,
  accessor, write, clear, or initial assignment.
- `internal/table/model.go::EmitValue` — doc: "An emit key is not a tag
  key: it is undeclared, uninterpreted, and compared by exact byte
  equality" — the in-repo pure-carrier precedent for an undeclared
  vocabulary (RDR 0010 `0010:C3`).
- `internal/cli/flow_state.go::parseWrites` — "A `--write` key always
  has a declaration — `writerFor` just proved a writer names it" ⇒ the
  zero-decl arm of `canonicalValue` is reachable only from `parseTags`.
- Peer proposals read via projector: `0012:§proposed-solution` C2
  ("Undeclared-key ADMISSION policy … is upstream and deliberately not
  decided here"), `0012:§decision-rationale` (Joint-check names 0020 as
  nearest coupling, pins no contract at the admission locus);
  `0024:§proposed-solution` (opt-in declarations constrain authoring,
  never evaluation; zero declarations = today byte-for-byte);
  `0018:§proposed-solution` (kernel sentinel `ErrReservedTagKey`;
  CheckInput reserved-key surface untouched by this RDR).

## Corpus queries (external class/instance claims)

Corpus: `StateMachineRes` (arc, semantic).

1. "undeclared context key admission open world caller-supplied
   extended state validation" — top hits header-thin (fizz api.md
   "Context And Runtime"); nothing quotable. Rejected.
2. "unknown or undeclared input keys ignored or rejected schema
   validation of event payload" — hits on OpenAPI error schemas,
   scxmlcc attribute docs; not this operator. Rejected.
3. "event payload data passed through opaque not validated by the
   engine arbitrary fields" — scxmlcc/inngest/awf-cli header chunks;
   nothing quotable. Rejected.
4. (claim 2, instance) "SCXML datamodel assign undeclared location
   error.execution event data unvalidated" — state-machine-cat and
   scxmlcc doc headers; the W3C-spec split (declared datamodel
   closed-world vs `_event.data` open) not quotable from this corpus.
   Rejected.

⚠ no prior-art coverage (quotable) for the external problem class
"undeclared caller-context key admission in peer engines"; the external
analogies (SCXML's declared-datamodel/open-event-data split; protobuf
unknown-field retention; Kubernetes label open-world) are model-prior
and were DEMOTED to a Resolve assumption (0020:A4) rather than leaned
on. The choice rests on the in-repo anchors above.

## Grounding micro-sweep (step 7.5)

One fresh-context sub-agent, factored brief (17 anchors, no justifying
prose). Verdict: PASS — 17 CONFIRMED, 0 REFUTED, 0 NOT-FOUND,
including the negative anchor (`flow-tag-undeclared` absent from
internal/ and docs/cli-output-contract.md). Recorded in the record as
`Ground-sweep: clean (17 anchors)`.

## Joint-decision check (step 8)

Open peers: the 12 Draft batch siblings (0001–0011 are `Implemented`,
closed — not open peers; no `Final` peer exists). Greps: `parseTags`,
`canonicalValue`, `flow_input.go`, `flow-tag-invalid`,
`flow-tag-undeclared`, `flow-tag-reserved`, `flow-tag-owned`,
`flow-tag-duplicate`, `not set-valued`, `unknown_tag`, `TagDecl`,
`model-schema.md`, `ConformValue`. Hits and triage: 0023
(`flow_input.go`, `flow-tag-invalid` — A7 Source-Read evidence, echo
group homed at JDR 0002 §D1; not a modify-anchor), 0017
(`flow_input.go::loadFindings` — different symbol, citation), 0012 /
0024 (`TagDecl` — vocabulary citation / explicit "NOT TagDecl"
contrast). Absence arm: due-diligence sweep of closed records — 0005's
kind-mismatch clauses are declared-key-scoped; `not set-valued`
appears in no other record. Verdict: clear (12 peers).
