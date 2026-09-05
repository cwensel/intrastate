# Recommendation 0020: Undeclared --tag key admission — what the zero TagDecl means

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
- **Type**: Bug Fix
- **Profile**: mid — the meaning of the zero `TagDecl` at `--tag` admission; user-facing yes; locks contract
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
- **Related Issues**: kata `intrastate#3exy` (1601); kata `p1p6`
  (companion doc-gap findings)
- **Predecessors**: 0005-skill-integration-cli-contract,
  0008-recognized-tag-key-ownership
- **Seam Lineage**: no prior accretion

## Problem Statement

A caller composing a producer's whole fact vector into `flow resolve
--tag` flags expects undeclared keys to pass through — only guard
atoms make dimensions, and the guarded-lookup comment in
`internal/cli/flow_input.go` promises exactly that ("the zero
`TagDecl` conforms everything"). But the zero decl does not conform
everything: its empty `Kind` makes `isSet` false in `canonicalValue`,
so an undeclared key carrying a canonical JSON array is refused
`flow-tag-invalid` with "the tag `extra` is not set-valued" — a
message pointing at a declaration that does not exist. Reproduced at
HEAD: an undeclared scalar passes, an undeclared array refuses. In the
motivating composition, two set-valued facts in the producer's output
caused every record of a 157-record corpus to refuse. The workaround —
declare the set-valued keys with no guard atom — is sound and shipped,
but undiscoverable from the error, and it forces a model to enumerate
keys it has no interest in.

The decision, made once — what the zero `TagDecl` means at the
admission site: (1) an undeclared key is a **pure carrier** — the zero
decl conforms everything, arrays included; the value passes through
uninterpreted since no guard atom can read it; (2) declaration is an
**identity precondition of admission** — an undeclared key is refused
as such, under a code naming the absence (e.g.
`flow-tag-undeclared`), with its own migration story, since callers
rely on undeclared scalars passing today; or (3) **ratify today's
split** — scalars pass, array literals refuse — with the message
rewritten to name the real cause and the remedy ("declare it, even
with no guard atom"). The chosen arm fixes the refusal code and
message as a consequence, and resolves the live contradiction between
the admission comment and `canonicalValue`. The asymmetry is
unrecorded: REQ-61 governs `--write`/`--clear` keys, RDR 0008
adjudicates only the reserved key `recognized`, and no RDR or the CLI
output contract speaks to the undeclared-key/array case.

## Critical Assumptions

- **A1 No caller, test, or fixture load-bears on the undeclared-array
  refusal** — nothing branches on `flow-tag-invalid` fired for a key
  the model does not declare (the "is not set-valued" arm over the
  zero decl).
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: the message literal has one producer,
    `internal/cli/flow_input.go::canonicalValue`, and NO asserting
    test — `grep "is not set-valued" --include="*_test.go"` returns
    zero repo-wide; the string's only occurrence anywhere is that
    producer. One test reaches the arm without asserting its text:
    `internal/cli/flow_input_0005_test.go`
    `TestReq29And80_WrongKindAndEmptyTagValuesAreFlowTagInvalid`,
    subtest `array-for-scalar-key`, which uses `--tag profile=[…]`
    against a key DECLARED at `internal/cli/flow_fixtures_0005_test.go`
    `[tags.profile]` (`kind = "enum"`, `single_valued = true` — a
    declared non-set kind, so the `!isSet` arm fires) — i.e. the
    declared arm, which C1 keeps; it asserts code and `param` only,
    via `requireRefusal`, never `ce.Message`. Every
    other `flow-tag-invalid` assertion is a different arm (empty
    `--outcome`, malformed no-`=` tag, declared-enum domain
    violation). The repo asserts the converse directly:
    `internal/cli/flow_input_0005_test.go`'s `an undeclared observed
    key still passes` requires success for
    `--tag sidechannel=whatever-the-caller-likes`.
    negative: no assertion, caller branch, or fixture requires
    `flow-tag-invalid` over an UNDECLARED key.
  - **If wrong**: the acceptance-conversion silently changes a
    consumer's error-handling path; surfaces as a caller branch on
    `flow-tag-invalid` that never fires again.
- **A2 An undeclared `--tag` value cannot influence resolution** —
  dimensions come only from rule atoms, and every model-side reference
  to an undeclared tag refuses at load, so the carried value is
  structurally unreadable.
  - **Status**: Verified
  - **Method**: Source Search + Spike
  - **Evidence**: all five `CatUnknownTag` refusal sites confirmed —
    `internal/table/normalize.go::atomsFromBlock` (match/guard/unless
    atoms in every block), the write-block and clear-list arms of
    `internal/table/normalize.go`, `internal/table/load.go::accessorTable`,
    and `internal/table/load.go::loadInitial`. Closed-world reader
    enumeration over every non-test `m.Tags` read found two disjoint
    shapes, neither able to observe an undeclared key's value:
    iterations over the model's OWN key set (vacuously excluding
    undeclared keys), and lookups keyed by a rule-atom/accessor/initial
    key (each already past one of the five refusals). Runtime: the
    observed tag enters the kernel view at
    `internal/resolve/resolve.go::assemble`, but the view has exactly
    four readers, none able to surface an undeclared key's value:
    `matches` (iterating `row.Match` — rule atoms only), `has` (keyed
    by `row.RequiresOwned`), guard `Lookup` (keyed by `atom.Key`), and
    `internal/resolve/resolve.go::conflicting`, which is value-blind —
    it returns `s.tags[key].conflicted`, a bool, and guard
    `evaluateAtom` reads it only to answer `ReasonUncomparable`;
    `internal/resolve/precondition.go::CheckInput`
    reads `Observed` only for the reserved key NAME, never a value.
    Runtime leg — spike `evidence/spikes/selection-diff-a2-vs-b.txt`
    (empty): selection is byte-identical with and without an undeclared
    scalar carrier; every projected field (rule, emit, dispositions,
    gates, next, writes, clear, escaped, outcome) diffs empty and the
    only payload change is `observed` gaining the carried key. The
    array-carrier leg is not runtime-verifiable until the change lands
    (it refuses at HEAD) and is covered by MVV 2/3.
  - **If wrong**: pass-through becomes semantically live and the
    carrier arm's safety argument collapses — routing could depend on
    unvalidated bytes; surfaces as a selection diff in the MVV pair.
- **A3 No byte-equality obligation reads an undeclared observed
  value** — no read-back, plan copy-through, or comparison path
  consumes it, so verbatim (non-canonicalised) carry is safe.
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: tag-value byte-equality lives in exactly two places,
    both confined to one domain — the write plan, the accessor-read
    baseline, and the accessor RE-READ values — and neither reads the
    `--tag` observed slice: `internal/accessor/executor.go::verifyReadBack`,
    and its CLI-side re-rendering
    `internal/cli/flow_exec.go::readBackFindings` (REQ-103, one finding
    per divergent key), whose inputs are `accessor.Refusal`'s
    `Expected`/`Observed` — the same plan-vs-re-read pair, already
    compared upstream. The plan comes
    from `internal/cli/flow_state.go::parseWrites`, which proves a
    declaration via `writerFor` BEFORE calling `canonicalValue` (its
    lookup "is total, unlike `parseTags`'s"), so every byte-compared
    value is canonicalised under a real declaration. All other
    consumers of a `--tag` observed value are non-comparing verbatim
    copies: `internal/cli/flow_exec.go`'s `observedTagMap` /
    `kernelTags`, and the payload `observed` echo
    (`internal/cli/flow_resolve.go`, routed to `groupEcho` by
    `internal/cli/flow_partition.go`). `internal/resolve/resolve.go`'s
    `merge` does compare a repeated key's value, but
    `internal/cli/flow_input.go` refuses a repeated `--tag` key under
    `codeTagDuplicate` (REQ-28) before the kernel sees it.
    negative: no read-back, plan copy-through, hash, sort, or dedupe
    path consumes an UNDECLARED observed value.
  - **If wrong**: a verbatim (unsorted/duplicated) array breaks an
    equality check somewhere downstream; surfaces as a spurious
    mismatch refusal.
- **A4 External prior art aligns with the carrier arm** — SCXML's
  declared-datamodel/open-event-data split, protobuf unknown-field
  retention, Kubernetes-style open label vocabularies.
  - **Status**: Verified
  - **Method**: Prior Art
  - **Evidence**: three opened citations, trail at
    `evidence/research/resolve-prior-art.md`. SCXML — W3C SCXML
    Recommendation §5.4 (`<assign>`): an illegal or non-existent
    datamodel location MUST raise `error.execution` (the closed,
    declared arm), against §5.10.1 (`_event`): the processor "SHOULD
    reformat this data to match its data model, but MUST NOT otherwise
    modify it" (the open, carried arm). Protobuf — Language Guide
    (proto3) §Unknown Fields: "Proto3 messages preserve unknown fields
    and include them during parsing and in the serialized output."
    Kubernetes — *Mastering Kubernetes* (Packt) §Label/§Annotation,
    pp.6–7: key/value SYNTAX is validated while the value's meaning is
    not ("Kubernetes just stores the annotations"); analogous rather
    than identical, since the API server has no declared-key table.
    Corpus-first per the budget: `StateMachineRes` / `StateMachineLit`
    / `DevRef` / `PapersFast` were negative for the SCXML and protobuf
    instance claims (the Propose-stage negative reproduced); the web
    fallback fetched the two primary specs without a 403.
  - **If wrong**: the alignment claim weakens; the decision still
    stands on the in-repo anchors (non-fatal — record and move on).
- **A5 RDR 0005's stable code table tolerates an arm of
  `flow-tag-invalid` becoming unreachable for undeclared keys** — the
  table fixes spellings and exit groups, not per-arm reachability, so
  no amendment of 0005 is required.
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: `0005:§failure-modes`'s table header fixes the scope
    exactly — "Stable code strings (the spellings are normative; the
    group fixes the exit)": spellings and exit group, nothing about
    which inputs reach which arm. `0005:§normative-contracts` C1's only
    `--tag` clauses cover the owned/reserved refusal and "`--tag`
    values MUST enter as observed context"; no clause addresses an
    undeclared, non-owned, non-reserved key's kind-checking. The
    nearest text, `0005:§technical-design` "State input" ("a bare
    scalar for a set key, or an array for a scalar key, is
    `flow-tag-invalid`"), is conditioned on a KNOWN declared kind, so
    it does not reach the undeclared case. Message text is nowhere
    contract in 0005 — C1 requires only that `Finding.message` be
    self-sufficient, not verbatim-pinned. Spellings mirror at
    `internal/cli/flow_input.go`'s constant block
    (`codeTagInvalid`/`codeTagDuplicate`/`codeTagReserved`/`codeTagOwned`),
    unchanged by C1. 0005 Status: Implemented.
    negative: no 0005 clause pins the undeclared-array refusal or ties
    the "is not set-valued" message to any input class.
  - **If wrong**: the change needs a successor-RDR amendment path for
    0005's table before it can land; surfaces at Stage 7's
    contradiction check.

## Proposed Solution

### Approach

Adopt arm (1): an undeclared `--tag` key is a **pure carrier**. The
zero `TagDecl` means what two of the three in-code authorities already
say it means — `internal/cli/flow_input.go::parseTags`'s guarded-lookup
comment ("conforms everything … shape-only behaviour a caller already
relies on") and `internal/table/load.go::ConformValue`'s doc ("A zero
TagDecl conforms everything … keeps an undeclared `--tag` key
shape-only") — and the one dissenting arm,
`internal/cli/flow_input.go::canonicalValue`'s `!isSet && looksArray`
refusal, is the bug. Admission distinguishes *undeclared* (key absent
from the model's tag table) from *declared scalar* (zero-valued `Kind`
never occurs for a declared key — the loader requires a kind), and for
an undeclared key admits the value verbatim, array literals included:
no canonicalisation, no kind or domain check, no comparison. The value
is provably uninterpreted — every model-side reference to an undeclared
tag (rule atom, accessor key, write, clear, `[initial]`) already
refuses at load under `unknown_tag` ⇒ nothing that routes can read the
carried bytes — so passing it through cannot corrupt selection. This is
the same disposition the emit namespace shipped under `0010:C3`
(`internal/table/model.go::EmitValue`: "undeclared, uninterpreted, and
compared by exact byte equality" — its writer is the rule author, its
carrier the payload). The provenance guards are untouched and still
precede admission: reserved key, owned key, duplicate (REQ-26/27/28).
No new refusal code is minted; the "is not set-valued" message becomes
truthful because it can only fire for a key whose declaration actually
says so.

### Technical Design

One seam moves: the admission path in `internal/cli/flow_input.go`.
`parseTags` switches its guarded lookup to the two-value form and
routes an undeclared key around `canonicalValue` (or passes the
declaredness bit into it — implementation latitude, bounded by C1's
refusal list: whichever shape is taken, the empty-value arm must still
fire for a carrier and must still leave a declared set on its
conformance message), so the kind/shape/domain arms run only under a
real declaration. The carried
value flows exactly where an undeclared scalar already flows today:
into the kernel's Observed view (`resolve.Input.Observed`, where only
atoms — all declaration-checked at load — can read keys) and out
through the resolve payload's `observed` echo, byte-for-byte as given.
`--write`/`--clear` are out of scope: `internal/cli/flow_state.go::
parseWrites` proves a writer binding (and therefore a declaration)
before its `canonicalValue` call, so the zero-decl arm is reachable
from `parseTags` alone.

#### Normative Contracts

One contract — the meaning of the zero `TagDecl` at `--tag` admission.

Determinacy: fired — C1 (identity: presence in `m.Tags`, byte-exact, no
folding or sentinel; and step order: the admission refusal precedence the
carrier branch splices into).

**C1**

```normative
PURE CARRIER. At `--tag` admission, a key ABSENT from the loaded
model's normalized tag table (the two-value `m.Tags[key]` lookup) is
a pure carrier: the CLI admits its value VERBATIM — byte-preserved,
JSON array literals included — and never canonicalises, conforms,
kind-checks, or compares it. The only refusals reachable for an
undeclared key are, unchanged and still preceding any accessor
(REQ-27):

- the flag grammar (`--tag` takes `name=value`), `flow-tag-invalid`;
- `flow-tag-reserved` (REQ-26), `flow-tag-owned` (REQ-27),
  `flow-tag-duplicate` (REQ-28);
- the empty-value arm, `flow-tag-invalid` — an empty observed value is
  indistinguishable from unset, which is a grammar-level fact about the
  value that holds with or without a declaration. It therefore binds the
  carrier too, and MOVES to the admission path ahead of the carrier
  branch; today it sits inside `canonicalValue`'s `!isSet` arm, which
  the carrier no longer enters. The hoisted arm carries the
  SCALAR-shaped message for a carrier — ``the tag `X` was given an
  empty value`` — which is byte-identical to what an undeclared key
  already receives at HEAD (normative fixture **F4**, Testing Strategy
  scenario 4, from `evidence/spikes/d-undeclared-empty.txt`), so the
  hoist is behaviour-preserving. The set-specific empty-value message ("is set-valued
  and takes a JSON array literal; got ") is a CONFORMANCE message and
  stays declared-only, inside the arms below: it presupposes a
  declared kind, which a carrier by definition has none of.
  The hoisted arm is therefore NOT unconditional on `value == ""`: it
  fires for a key that is undeclared, or declared with a NON-set kind,
  and it must not intercept a declared SET key, whose empty value keeps
  the set-specific conformance message (normative fixture **G**,
  `evidence/spikes/g-declared-set-empty.txt`: `--tag labels=` ⇒ ``the
  tag `labels` is set-valued and takes a JSON array literal; got ``).
  A declared scalar's empty value keeps the scalar message it has at
  HEAD (`evidence/spikes/e-declared-empty.txt`), which is the same
  string the hoisted arm emits, so routing it through either site is
  byte-identical. Hoisting on the bare value test alone would move a
  declared set key off its conformance message — a declared-key change
  this RDR does not make.

`flow-tag-invalid`'s kind, shape, and domain arms — including "is not
set-valued" — are reachable only for a DECLARED key, whose loaded
declaration (kind required at load, one of RDR 0003's five tokens)
is what the message then truthfully reports. No new refusal code is
minted; `flow-tag-undeclared` does not exist. The carried value rides
the kernel's Observed view and echoes in the resolve payload's
`observed` field byte-for-byte as given; declared set values keep the
canonical-array form of `docs/cli-output-contract.md` §Set values on
the wire. The model-side closed world is untouched: a rule atom,
accessor key, write target, clear target, or `[initial]` assignment
naming an undeclared tag still refuses at load (`unknown_tag`).
Declaring a key later is the opt-in tightening: admission then
enforces that declaration's kind and domain.
```

#### Load-Bearing Decisions

- **Identity** — "declared" means the key is present, byte-exact, in
  the normalized model's tag table (`m.Tags`, two-value lookup); the
  same presence signal `internal/table/load.go::accessorTable` already
  branches on (`_, ok := l.model.Tags[key]`). No parallel registry, no
  case folding, no zero-decl sentinel.
- **Wire / byte format** — an undeclared value crosses verbatim, never
  re-canonicalised: canonicalising would interpret as a set a value
  whose declaration never asserted set-ness. Declared set values keep
  the §Set-values-on-the-wire canonical array unchanged.
- **Naming** — no new code; the rejected alternative's
  `flow-tag-undeclared` is deliberately not minted (Alternatives,
  arm 2).
- **Selection / predicate** — refusal precedence at admission is
  unchanged in order: grammar → reserved → owned → duplicate →
  empty-value (undeclared or non-set-declared keys) → (declared keys
  only) kind/shape/domain conformance;
  the carrier branch sits where the conformance arms would have run.
  The empty-value arm is the one that moves — out of `canonicalValue`
  and up to the shared path — because it must bind the carrier, which
  no longer reaches that function. It lands AFTER the duplicate arm and
  the `seen[key]` mark, not before: at HEAD the duplicate check already
  precedes the empty-value site (`flow_input.go:663` vs the `:723` arm
  inside `canonicalValue`), so a repeated key refuses
  `flow-tag-duplicate` today and must still do so — hoisting the arm
  ahead of the duplicate case would invert that precedence and change
  behaviour this RDR does not touch. The order above is the whole
  constraint on the splice point; where the arm sits within it is
  implementation latitude.

#### Authority census

Fired: the draft names three in-code authorities on the zero `TagDecl`'s
meaning and adjudicates between them (two promise carrier, one dissents).

| Input / decision | Writer | Readers | Call sites | Sibling arms | Canonical |
| --- | --- | --- | --- | --- | --- |
| "is this key declared?" at `--tag` admission | model author (`[tags.*]`) | `parseTags` | `flow_input.go::parseTags` (one-value today; two-value under C1) | `accessorTable`, `loadInitial`, `atomsFromBlock`, `renderWrites` write/clear arms — all `_, ok := Tags[key]` | **the two-value `m.Tags[key]` presence test** — C1 adopts the sibling form; no new signal |
| what a zero `TagDecl` MEANS | — | `parseTags` comment (carrier); `ConformValue` doc (carrier); `canonicalValue`'s `!isSet && looksArray` (refuse) | `flow_input.go:669-675`, `load.go:1811`, `flow_input.go:718-722` | the two comments agree; the third dissents | **C1** — the comments are right, the `looksArray` arm is the bug |
| the carried value's meaning | caller | payload `observed` echo only | `groupEcho` via `flow_partition.go` | model-side readers all refuse at load (`unknown_tag`, 5 sites) | **uninterpreted** — no reader can resolve it (A2) |

#### Disposition

Fired: C1 assigns exit outcomes across input classes at one admission seam.

| Input class | Exit | Code / message | Artifact | Silent or loud |
| --- | --- | --- | --- | --- |
| undeclared key, scalar value | 0 | — | `observed` echo, verbatim | loud (echoed) |
| undeclared key, array literal | 0 (**changed**; refuses at HEAD) | — | `observed` echo, verbatim | loud (echoed) |
| undeclared key, empty value | 2 | `flow-tag-invalid` "was given an empty value" | none | loud |
| undeclared key that is a MISSPELLING of a declared one | 0 | — | `observed` echo, verbatim | **silent** — accepted open-world cost; diagnosis is the echo (Failure Modes) |
| declared scalar/enum key, array literal | 2 | `flow-tag-invalid` "is not set-valued" — now truthful | none | loud |
| declared set key, empty value | 2 | `flow-tag-invalid` "is set-valued and takes a JSON array literal; got " | none | loud |
| reserved / owned / duplicate key, any value | 2 | `flow-tag-reserved` / `-owned` / `-duplicate` | none | loud (unchanged, still first) |

#### Desk trace

Fired: C1, the Testing Strategy Expected lines, and fixtures F1–F4 and G all bear
on one output surface (the admission refusal path and the `observed` echo).
Walk of the MVV, each step against the assertions in force, with a witness
from the named spike.

| MVV step | Assertions in force | Witness | Verdict |
| --- | --- | --- | --- |
| 1 — fixture model: scalar `tier`, set `labels`, nothing named `extra`/`extras` | C1's identity rule (absence from `m.Tags`) | `evidence/spikes/carrier-table.toml` | consistent |
| 2 — `--tag extra=plain --tag extras=["a","b"]` | C1 carrier admits verbatim; Testing 1 (exit 0, byte-for-byte echo); F1 | F1 scalar leg at HEAD: `"observed":{"extra":"plain","labels":"[\"security\"]","tier":"free"}`; array leg is post-change, red test adds it | consistent — the array leg is the only unwitnessed cell, correctly marked post-change |
| 3 — same resolve without carriers | C1 (nothing can read it); A2 runtime leg; Testing 2; F2 | `selection-diff-a2-vs-b.txt` empty; `"rule":"free"`, `"emit":{"plan":"basic"}` both runs | consistent |
| 4 — declared scalar given array | C1's kind arm, declared-only; Testing 3; F3 | `c-declared-scalar-array.txt`: exit 2, "the tag `tier` is not set-valued; got the array literal [\"a\",\"b\"]" | consistent — and the message is now truthful, since `tier` IS declared |
| 5 — undeclared key, empty value | C1's hoisted empty-value arm (scalar-shaped message, binds carriers); Testing 4; F4 | `d-undeclared-empty.txt`: exit 2, "the tag `extra` was given an empty value" — byte-identical to HEAD | consistent — hoist is behaviour-preserving on the witness |
| 5b — declared SET key, empty value (not an MVV step; the hoist's blast radius) | C1's guard: the hoisted arm skips a declared set, which keeps the conformance message; Testing 5; fixture G | `g-declared-set-empty.txt`: exit 2, "the tag `labels` is set-valued and takes a JSON array literal; got " | consistent ONLY with the guarded arm — an unconditional `value == ""` test ahead of the declared lookup would emit the scalar message here and regress a declared key |
| end state | `D-selection-predicate` order: grammar → reserved → owned → duplicate → empty-value (undeclared or non-set-declared) → (declared only) conformance | source order confirmed at grounding: duplicate (`flow_input.go:663`) precedes `canonicalValue` (:677); empty-value today sits at :723 inside `!isSet` | consistent — the hoist moves one arm forward across no other arm, and the guard keeps declared sets on their existing arm |

No CONTRADICTION row. Row 5b was a contradiction at the repeatability
lens (the Illustrative Code's unconditional `value == ""` test against
this table's declared-set row) and is resolved by C1's guard.

#### Illustrative Code

Illustrative — shape only, not load-bearing:

```go
switch { // unchanged, and still first
case key == table.RecognizedTagKey: // reserved
case slices.Contains(owned, key): // owned
case seen[key]: // duplicate — still precedes the empty-value arm
}
seen[key] = true

decl, declared := m.Tags[key]
if value == "" && (!declared || !decl.IsSet()) {
    // Empty is indistinguishable from unset. Binds carriers and declared
    // scalars; a declared SET key keeps its conformance message (fixture G).
    return nil, userErr(codeTagInvalid, key, "…was given an empty value")
}
if !declared {
    // Pure carrier (0020:C1): admit verbatim; nothing can read it.
    out = append(out, resolveTag{Key: key, Value: value})
    continue
}
canonical, ce := canonicalValue(key, value, decl, "tag")
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Undeclared-key discrimination at admission | `internal/cli/flow_input.go::parseTags` + `canonicalValue` | zero-decl lookup conflates "undeclared" with "declared non-set" | Extend | two-value lookup; carrier branch precedes the conformance arms |
| Declaration-presence signal | `internal/table/load.go::accessorTable` (`_, ok := l.model.Tags[key]`) | none | Reuse | same discriminator; no new signal invented |
| Uninterpreted-vocabulary precedent | `internal/table/model.go::EmitValue` (`0010:C3`) | different namespace (emit, not tags) | Reuse (as precedent) | carrier semantics mirror an adjudicated in-repo disposition |

### Decision Rationale

Key factors. (a) Two of the three in-code authorities already promise
carrier semantics — the `parseTags` guarded-lookup comment and
`ConformValue`'s zero-decl doc — and only `canonicalValue`'s
first-byte sniff dissents; the fix restores the documented contract
rather than inventing one. (b) Safety is structural, not disciplinary:
the loader's `unknown_tag` refusals make an undeclared key unreadable
by any rule, accessor, write, clear, or initial assignment ⇒ the
carried bytes cannot influence selection. (c) The house already
adjudicated this shape once: emit keys are "undeclared, uninterpreted,
byte-compared" (`0010:C3`), with opt-in declarations arriving later
(0024) — carrier-by-default, declare-to-tighten is the established
pattern. (d) The user outcome: a producer's whole fact vector composes
into `--tag` flags without the model enumerating alien keys — the
exact failure that motivated this RDR. (e) `looksArray` on an
undeclared value is an interpretation of bytes the system pledges not
to interpret, and it misclassifies a legitimate scalar that merely
begins with `[`.

Rejections: arm (2) (`flow-tag-undeclared`) breaks the shipped
undeclared-scalar pass-through, forces models to enumerate keys they
have no interest in (the motivating pain, made mandatory), and buys a
typo guard the payload's `observed` echo already surfaces; arm (3)
ratifies a first-byte heuristic with no principled defense, leaves the
motivating composition broken, and still refuses legitimate
`[`-prefixed scalars; per-model opt-in strictness (the 0024 shape) is
compatible later work, not the default's meaning — see Briefly
Rejected.

Premortem: survived (paragraph) — the shipped-and-failed narrative is
the silent-typo swallow: a caller misspells a declared set-valued key,
the misspelling passes as a carrier, the intended rule silently fails
to match, and resolution refuses no-match or routes to an escape row —
a debugging session with no error naming the key. The approach
answers it: this is today's shipped behavior for scalars (the comment
calls it "behaviour a caller already relies on"), so arm (1) extends
an accepted open-world cost rather than creating one; the resolve
payload's `observed` echo lists the stray key beside the declared
ones, which is where the pinned MVV test points a diagnostician; and
the remedy — declare the key — is exactly the opt-in tightening C1
names. The failure the approach could not answer would be a routing
path that reads undeclared bytes; A2 pins that no such path exists,
and if Resolve refutes A2 the choice reopens. Recommendation stands.
Ground-sweep: clean (17 anchors).
Joint-check: fired → 0025 (home: cli/0020 §Normative Contracts C1 /
cli/0025 §Normative Contracts C1 — mutual tolerance: disjoint rule
families at `accessorTable`; tag-admission identity homed here,
accessor-entry shape homed in 0025; composes). Remaining peers clear
(12) — grep hits triaged, reported as
context: 0023 carries `flow_input.go` and `flow-tag-invalid` in its
A7 assumption's Source-Read evidence (the echo group is an echo), not
as a modify-anchor or a fenced pin — carrier admission strengthens
the verbatim-echo reading it rests on, and the echo-group content
question is already homed at JDR 0002 §D1's partition doctrine (the
disposed 0023↔0024 coupling), not a new fire; 0017 cites
`flow_input.go::loadFindings` — a different symbol, precedent
citation only; 0012 answers declared-key VALUES at the guard seam one
seam downstream and its C2 states "Undeclared-key ADMISSION policy …
is upstream and deliberately not decided here" (its own Joint-check
names 0020 as its nearest coupling, clear — compositional); 0018
exports `ErrReservedTagKey` around the kernel's reserved-key channel,
which precedes admission and does not move. No peer carries
`parseTags`, `canonicalValue`, `not set-valued`, or
`flow-tag-undeclared`. Absence arm (this proposal converts a refusal
into an acceptance): no `Final`-status peer exists — 0001–0011 are
`Implemented` (closed; any coupling rides to 7.1) — and a
due-diligence sweep of them found the kind-mismatch language scoped
to DECLARED keys only ("a bare scalar for a set key, or an array for
a scalar key", 0005), no reliance on the undeclared-array refusal;
A5 pins that reading of 0005's code table at Resolve. Bridge
sub-check: n/a — no surface here is scheduled for deletion by a
sibling plan and this plan retires none.

## Alternatives Considered

### Alternative 1: Declaration as identity precondition (`flow-tag-undeclared`)

**Description**: Arm (2) — an undeclared `--tag` key is refused as
such, under a new code naming the absence (`flow-tag-undeclared`),
with a migration window for callers relying on undeclared scalars
passing today. Symmetric with the model-side closed world: the
caller's vocabulary becomes as declared as the author's.

**Pros**:

- Catches a misspelled key loudly instead of silently carrying it.
- One uniform rule — no declared/undeclared branch at admission.
- Symmetric with the loader's `unknown_tag` closed world.

**Cons**:

- Breaking: undeclared scalars pass today and the guarded-lookup
  comment calls that "behaviour a caller already relies on" — a
  migration story and a new code (0005's table grows) for no consumer
  demand.
- Mandates the exact DX failure that motivated this RDR: every model
  must enumerate a producer's whole fact vocabulary, including keys it
  has no interest in.
- Diverges from the in-repo precedent for uninterpreted vocabulary
  (`0010:C3`: emit keys undeclared and uninterpreted; 0024 makes
  declaration opt-in, not mandatory).

**Reason for rejection**: converts the motivating bug into a mandate
for the workaround's worst property (enumerate-everything), at the
cost of a breaking change plus a new refusal code, to buy a typo
guard the payload's `observed` echo already provides.

### Alternative 2: Ratify today's split, rewrite the message

**Description**: Arm (3) — keep scalars-pass/arrays-refuse, and
rewrite the refusal to name the real cause and remedy ("`extra` is
not declared; declare it — even with no guard atom — to pass a set
value").

**Pros**:

- Zero behavior change; cheapest; the message defect (blaming a
  nonexistent declaration) is fixed.

**Cons**:

- Enshrines a first-byte heuristic (`looksArray`) as contract:
  admission of an undeclared key depends on whether its value begins
  with `[` — an interpretation of bytes the system pledges not to
  interpret, with no principled line to defend at Finalization.
- Refuses a legitimate undeclared scalar that merely begins with `[`.
- Leaves the motivating composition broken: set-valued facts still
  force declarations, now merely with a better error.

**Reason for rejection**: ratifies an accident as a contract; the
better message treats the symptom while the asymmetry it apologizes
for remains indefensible.

### Briefly Rejected

- **BR1 Per-model opt-in strictness (the 0024 shape — a model switch
  making undeclared `--tag` keys refuse)**: compatible later work
  layered on top of the carrier default, but it answers "may a model
  opt out of the default?" while this RDR must first fix what the
  default *means* — and nothing motivates it yet.
- **BR2 Canonicalise undeclared array-looking values as sets**:
  interprets undeclared bytes (the same sin as `looksArray`) and
  invents a byte-equality obligation no consumer holds (A3).

## Context

### Background

Found while authoring the motivating consumer model for RDR 0010's
decision-table class; tracked as kata `intrastate#3exy`, with the
companion documentation findings on kata `p1p6`. No false green — the
site refuses with exit 2, so nothing routes on bad input; the cost is
DX and composability, and every producer/consumer pair with set-valued
outputs meets it. Each recurrence costs a debugging session rather
than a line precisely because the error is undiscoverable from its own
message. Whichever arm wins: land a red test asserting the chosen
disposition for an undeclared array literal *with the scalar case
beside it*, so the asymmetry is pinned rather than incidental; and
under arm (1) or (3), a `docs/model-schema.md` line for authors
composing a producer's whole tag output.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/cli/flow_input.go` (the guarded-lookup comment,
`decl := m.Tags[key]`, and `canonicalValue`'s `!isSet` +
`looksArray` arm), `docs/cli-output-contract.md` (§Set values on the
wire — canonical JSON arrays, which is what makes an undeclared set
indistinguishable from a malformed scalar at this site). Governing
records: RDR 0005 (refusal-code table, REQ-61), RDR 0008 (reserved
key `recognized`).

## Research Findings

### Investigation

Read at Propose (trail: `evidence/research/propose-research.md`): the
admission seam (`internal/cli/flow_input.go::parseTags` /
`canonicalValue`), the runtime conformance surface
(`internal/table/load.go::ConformValue` and its zero-decl doc), the
loader's undeclared-tag refusal sites (`internal/table/normalize.go`,
`load.go::accessorTable`, the `[initial]` arm — all `unknown_tag`),
the `--write` path (`internal/cli/flow_state.go::parseWrites`), the
emit-namespace precedent (`internal/table/model.go::EmitValue`,
`0010:C3`), `docs/cli-output-contract.md` §Set values on the wire, and
the sibling proposals 0012 / 0018 / 0024 via the projector. External
corpus reads (4 arc queries over `StateMachineRes`) surfaced no
quotable passage for the class "undeclared caller-context key
admission in peer engines"; the external analogies were demoted to A4
rather than leaned on, and the choice rests on the in-repo anchors.
Resolve reproduced that corpus negative instance-by-instance (SCXML,
protobuf) and closed the gap on the web fallback: the primary specs
fetched without a 403 and A4 now carries three opened citations. The
choice still rests on the in-repo anchors — prior art corroborates it,
it does not carry it.

### Key Discoveries

- **Documented** — the contradiction is internal to the code: the
  `parseTags` guarded-lookup comment and `ConformValue`'s doc both
  promise "the zero `TagDecl` conforms everything / shape-only", while
  `canonicalValue`'s `isSet := decl.Kind == "set"` plus `looksArray`
  refuses an undeclared array ⇒ the dissenting arm, not the comments,
  is the odd one out.
- **Documented** — an undeclared key is structurally unreadable: every
  model-side reference (rule atom, accessor `keys` member, write,
  clear, `[initial]` target) refuses at load under `unknown_tag`
  (`internal/table/normalize.go`, `internal/table/load.go`) ⇒ a
  carried value can never reach selection; its writer is the caller
  and its only reader the payload's `observed` echo.
- **Documented** — a declared key never carries an empty `Kind`: the
  loader refuses any kind outside RDR 0003's five tokens ⇒ zero-`Kind`
  at the admission seam identifies "undeclared" exactly, and the
  two-value lookup makes that explicit.
- **Documented** — the zero-decl arm of `canonicalValue` is reachable
  only from `parseTags`: `parseWrites` proves a writer binding (and
  thus a declaration) before its call ⇒ the decision is `--tag`-scoped
  and REQ-61's `--write`/`--clear` surface does not move.
- **Documented** — the in-repo precedent for undeclared vocabulary is
  pure carrier: `EmitValue` is "undeclared, uninterpreted, and
  compared by exact byte equality" (`0010:C3`), with declarations
  arriving later as opt-in (0024) ⇒ carrier-by-default,
  declare-to-tighten is the established house pattern.
- **Documented** (Resolve, corrected at grounding) — no caller
  load-bears on the undeclared-array refusal: the "is not set-valued"
  literal has one producer and NO asserting test (the one test that
  reaches the arm is over a DECLARED key and asserts code and `param`
  only), and the repo already asserts the converse ("an undeclared
  observed key still passes") ⇒ A1 verified, more strongly than the
  Resolve reading. No byte-equality path reads an undeclared observed
  value: the two comparison sites
  (`internal/accessor/executor.go::verifyReadBack` and its CLI
  re-rendering `internal/cli/flow_exec.go::readBackFindings`) share
  one domain, the declaration-proved write plan ⇒ A3 verified.
- **Documented** (Resolve) — external prior art aligns after all: the
  Propose-stage corpus negative reproduced, but the primary specs
  opened cleanly on the web fallback — SCXML §5.4 vs §5.10.1,
  protobuf proto3 §Unknown Fields, and Kubernetes labels/annotations
  ⇒ A4 verified, trail at `evidence/research/resolve-prior-art.md`.

## Trade-offs

### Consequences

- Positive: a producer's whole fact vector composes into `--tag` flags
  with no alien-key declarations; the 157-record motivating corpus
  resolves.
- Positive: the comment/code contradiction closes on the side both
  comments already document; the false "is not set-valued" message can
  no longer fire without a declaration to blame.
- Positive: the open-caller/closed-model split becomes a recorded
  contract (C1) instead of an accident.
- Negative: the silent-typo swallow extends from scalars to arrays — a
  misspelled declared set key now passes as a carrier instead of
  refusing with a wrong message; accepted as the open-world cost
  already shipped for scalars, mitigated by the `observed` echo and by
  declaring the key.
- Negative: declaring a previously-carried key later tightens
  admission (kind/domain refusals appear) — intended, and the same
  opt-in semantics 0024 ships for emit.

### Risks and Mitigations

- **Risk**: a consumer branches on the undeclared-array refusal today
  and the acceptance-conversion breaks its error path.
  **Mitigation**: A1's Resolve sweep before implementation; the pinned
  asymmetry test documents the new disposition.
- **Risk**: some downstream path compares observed values byte-wise
  and a verbatim (unsorted) array breaks it.
  **Mitigation**: A3's Resolve trace; if refuted, the carrier arm
  revisits verbatim-vs-canonical carry before lock.

### Failure Modes

Visible: every refusal that survives is unchanged and still fires at
exit 2 before any accessor — reserved, owned, duplicate and grammar for
ANY key, and empty-value for any key that is not a declared set (a
declared set's empty value refuses on its conformance message instead);
wrong kind, wrong shape and out-of-domain for
a DECLARED key only (per C1's list). Silent: a misspelled key (declared or not)
passes as a carrier and the intended rule fails to match; resolution
then refuses no-match or routes to an escape row. Diagnosis: the
resolve payload's `observed` field echoes every carried key
byte-for-byte — the stray spelling sits beside the declared keys in
the same envelope the refusal rides. Recovery: fix the spelling, or
declare the key (no guard atom needed) to put it under conformance.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1–A3, A5; A4 is non-fatal)

### Minimum Viable Validation

1. Author a fixture model declaring one scalar tag and one set tag,
   with no declaration for `extra` or `extras`.
2. Run `flow resolve` with `--tag extra=plain` and
   `--tag extras=["a","b"]` beside the declared tags: the invocation
   exits 0 and the payload's `observed` field carries both values
   byte-for-byte as given.
3. Run the same invocation without the two carrier flags: the selected
   rule and outcome are identical — the carried keys influenced
   nothing (A2's runtime leg).
4. Hand the DECLARED scalar tag an array literal: still refused
   `flow-tag-invalid` "is not set-valued" — now truthfully, and the
   red test pins this beside step 2 so the asymmetry is authored, not
   incidental.
5. Run `--tag extra=` (undeclared, empty): still refused
   `flow-tag-invalid` "was given an empty value" — the one arm the
   carrier does not escape, pinned so the hoist out of
   `canonicalValue` cannot silently drop it.

### Phase 1: Code Implementation

#### Step 1: Carrier branch at admission

Switch `parseTags`'s guarded lookup to the two-value form and route an
undeclared key past the conformance arms per C1. Grammar, reserved,
owned and duplicate are unchanged and in order. The empty-value arm
lifts out of `canonicalValue`'s `!isSet` branch onto the shared
admission path, ahead of the carrier branch, so it still binds a key
that no longer reaches that function — guarded to fire only for an
undeclared or non-set-declared key, since a declared set's empty value
must keep the set-specific conformance message (fixture G);
`canonicalValue` keeps its copy
for the `--write` carrier, which enters by its own path.

#### Step 2: Truthful comments

Rewrite the guarded-lookup comment to cite this RDR's carrier contract
instead of promising a different decision elsewhere;
`ConformValue`'s zero-decl doc stays true as written.

#### Step 3: Pinned tests and docs

Land the red test of MVV steps 2–5 (undeclared scalar AND array pass,
declared-scalar-given-array still refuses, undeclared-empty still
refuses), and the `docs/model-schema.md` line for authors composing a
producer's whole tag output (Background's obligation).

## Validation

### Testing Strategy

The scenarios are the Minimum Viable Validation's five steps, landed as
the red test of Phase 1 Step 3; they are authored there and not
restated here. Coverage goal — done is all six green, with the
undeclared/declared pair adjacent in one test so the asymmetry C1 fixes
is pinned by construction rather than by two tests that could drift:

1. **Scenario**: undeclared scalar and undeclared array literal
   admitted (MVV 2).
   **Expected**: exit 0; `observed` echoes both byte-for-byte.
   Normative fixture **F1** (scalar leg, read at HEAD from
   `evidence/spikes/a2-undeclared-scalar-only.txt`), for a model
   declaring scalar `tier` and set `labels` and nothing named `extra`:
   `--tag tier=free --tag labels=["security"] --tag extra=plain` ⇒
   `"observed":{"extra":"plain","labels":"[\"security\"]","tier":"free"}`.
   The array leg (`extras=["a","b"]` echoing verbatim) is the
   post-change extension of F1 and is what the red test adds; it has no
   HEAD witness (the array literal refuses there), so the asserted wire
   form follows from `resolve.Input.Observed`'s `map[string]string`
   type — the raw argument text as a string value,
   `"extras":"[\"a\",\"b\"]"`, the same shape F1 already witnesses
   for the declared set key `labels`.
2. **Scenario**: the same resolve with and without the carrier flags
   (MVV 3).
   **Expected**: identical selected rule and outcome. Normative
   fixture **F2** (`evidence/spikes/b-baseline-no-carriers.txt`,
   diffed at `evidence/spikes/selection-diff-a2-vs-b.txt`, empty):
   `"rule":"free"`, `"emit":{"plan":"basic"}`, with `gates`, `next`,
   `writes`, `clear` empty and `escaped:false` in both runs — every
   projected field diffs empty and the sole payload delta is
   `observed` gaining the carried key.
3. **Scenario**: declared scalar given an array literal (MVV 4).
   **Expected**: refused `flow-tag-invalid`, "is not set-valued".
   Normative fixture **F3**
   (`evidence/spikes/c-declared-scalar-array.txt`): `--tag
   tier=["a","b"]` ⇒ exit 2, ``{"code":"flow-tag-invalid","message":"the
   tag `tier` is not set-valued; got the array literal
   [\"a\",\"b\"]","param":"tier"}``.
4. **Scenario**: undeclared key given an empty value (MVV 5).
   **Expected**: refused `flow-tag-invalid`, "was given an empty
   value" — the arm the carrier does not escape. Normative fixture
   **F4** (`evidence/spikes/d-undeclared-empty.txt`): `--tag extra=`
   ⇒ exit 2, ``{"code":"flow-tag-invalid","message":"the tag `extra`
   was given an empty value","param":"extra"}`` — byte-identical to
   HEAD, so C1's hoist of this arm is behaviour-preserving. The
   declared-SET empty-value message differs
   (`evidence/spikes/g-declared-set-empty.txt`: "is set-valued and
   takes a JSON array literal; got ") and stays inside the
   conformance arms, declared-only.
5. **Scenario**: declared SET key given an empty value — the hoist's
   one regression risk, pinned beside scenario 4 so the guard on the
   hoisted arm cannot be dropped silently.
   **Expected**: refused `flow-tag-invalid` on the SET-specific
   conformance message, not the scalar-shaped one. Normative fixture
   **G** (`evidence/spikes/g-declared-set-empty.txt`): `--tag labels=`
   ⇒ exit 2, ``{"code":"flow-tag-invalid","message":"the tag `labels`
   is set-valued and takes a JSON array literal; got ","param":"labels"}``
   — byte-identical to HEAD. An unconditional `value == ""` test ahead
   of the declared lookup fails this scenario, which is what makes it
   the guard's witness.


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

- `internal/cli/flow_input.go` — `parseTags` (guarded lookup, REQ-26/27/28
  ordering) and `canonicalValue` (the `!isSet` + `looksArray` arm, the
  empty-value arm).
- `internal/table/load.go` — `ConformValue`'s zero-`TagDecl` doc;
  `accessorTable`'s `_, ok := l.model.Tags[key]` presence signal.
- `internal/table/normalize.go` — the `unknown_tag` refusal sites.
- `internal/cli/flow_state.go` — `parseWrites`, which proves a writer
  binding before its `canonicalValue` call (why `--write` is out of scope).
- `internal/table/model.go` — `EmitValue`, the uninterpreted-vocabulary
  precedent (`0010:C3`).
- `docs/cli-output-contract.md` §Set values on the wire.
- RDR 0005 (refusal-code table, REQ-61), RDR 0008 (reserved key
  `recognized`), RDR 0003 (the five kind tokens), RDR 0010 / 0024 (emit
  namespace: carrier-by-default, declare-to-tighten), RDR 0012 (declared-key
  values one seam downstream).
- Kata `intrastate#3exy` (1601); companion doc-gap kata `p1p6`.
