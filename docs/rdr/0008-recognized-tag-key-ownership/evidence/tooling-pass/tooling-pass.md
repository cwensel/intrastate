Model: claude-opus-5[1m]

# Tooling Pass — RDR 0008 (recognized-tag-key-ownership)

Mechanical adherence sweep, run as the Stage 7 pre-step. Post-mutation
regression check over the draft the four pre-lock lenses and the Stage 6
reconcile rewrote. Findings only; no RDR edits made by this sweep.

Repo HEAD at sweep: `1caf456`.

## CHECK 1 — Template section coverage

Every **Required** (spine) section is Present-substantive:

| TEMPLATE section | Verdict |
|---|---|
| H1 / Metadata | Present-substantive (no `[NUMBER]`/`[TITLE]` residue) |
| Problem Statement | Present-substantive |
| Context › Background, Technical Environment | Present-substantive |
| Research Findings › Investigation, Key Discoveries | Present-substantive |
| Critical Assumptions | Present-substantive (13 records, A1–A13) |
| Proposed Solution › Approach, Technical Design | Present-substantive |
| Normative Contracts | Present-substantive (6 blocks + closure argument) |
| Load-Bearing Decisions | Present-substantive (Identity, Naming, Enforcement locus) |
| Capability Dependencies | Present-substantive (5 rows) |
| Existing Infrastructure Audit | Present-substantive (3 rows) |
| Decision Rationale | Present-substantive (QOC matrix + sibling-path check) |
| Alternatives Considered (2 + Briefly Rejected) | Present-substantive |
| Trade-offs › Consequences, Risks and Mitigations, Failure Modes | Present-substantive |
| Implementation Plan › Prerequisites, MVV, Phases 1–3 | Present-substantive |
| Validation › Testing Strategy | Present-substantive (9 scenarios) |
| Finalization Gate | Pointer to `artifacts/gate.md` at lock (this pass) |
| References | Present-substantive (5 real citations, no brackets) |

**Conditional sections cleanly deleted — PASS, not Missing** (false-positive
guard applied): `Round-Trip / Inverse Invariants` (this RDR defines no
encode/decode inverse pair), `Illustrative Code` (shape-only, optional),
`Day 2 Operations` (creates no persistent resource), `New Dependencies`
(adds none — zero `go.mod`/third-party references), `Performance
Expectations` (compares no alternatives on empirical performance).

Phase headings are renamed from the template's illustrative
`Phase 1: Code Implementation` / `Phase 2: Operational Activation` to
`Phase 1: Contract ratification` / `Phase 2: Normalizer name validation` /
`Phase 3: Conformance pinning` — the template names those as titles to
replace, so this is conformance, not drift.

No surviving `_Draft placeholder._`, no `seed skeleton` header, no
old-template instructional block. The only bracketed text in the document is
the Finalization Gate's own five response prompts, which this pass replaces
with the `gate.md` pointer at lock — correct state entering Stage 7, not a
finding.

**C1: no findings.**

## CHECK 2 — Method label vocabulary

All 13 Critical Assumption Evidence Records carry a Method drawn from the
sanctioned eight (compound labels are two sanctioned labels joined, not
paraphrase):

| Record | Method | In vocabulary |
|---|---|---|
| A1 | Peer RDR | yes |
| A2 | Source Search | yes |
| A3 | Prior Art | yes |
| A4 | Source Search | yes |
| A5 | Derivation | yes |
| A6 | Peer RDR | yes |
| A7 | Source Search + Peer RDR | yes (both) |
| A8 | Source Search + Peer RDR | yes (both) |
| A9 | Spike | yes |
| A10 | Source Search + Peer RDR | yes (both) |
| A11 | Peer RDR | yes |
| A12 | Peer RDR + Source Search | yes (both) |
| A13 | Source Search | yes |

No missing, paraphrased, or off-vocabulary Method. **C2: no findings.**

## CHECK 3 — Source Search self-reference

Every `Source Search` record (A2, A4, A7, A8, A10, A12, A13) resolves its
Evidence to `internal/resolve/*.go`, `internal/cli/respond/respond.go`, peer
RDR bodies, peer spike fixtures, or the RDR engine's own prompt tree. **No
record cites `0008-recognized-tag-key-ownership.md` or any path under this
RDR's `artifacts/` as its own evidence.**

**C3: no findings** (the regression this check guards did not occur).

## CHECK 4 — Docs Only on load-bearing claims

Zero records carry Method `Docs Only`. **C4: no findings.**

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

Every cited `path::Symbol` resolves in the file named, verified at HEAD
`1caf456` (delegated read; per-anchor confirmations below). Line numbers were
not checked, per the check's own instruction.

`internal/resolve/resolve.go` — `recognizedTagKey` (`const … = "recognized"`,
L110); `assemble` (L148, merges Recognized/Observed/Owned into one `TagSet`);
`Resolve` (L318, `func Resolve(in Input) (Result, error)`); `escapeOrRefuse`
(L473); `refuse` (L515); `missingOwned` (L444); `gate` (L376);
`evaluateGuard` (L425); `Input` (L224, carrying `Owned []Tag`,
`Observed []Tag`, `Recognized string`); `Tag` (L20, `{Key, Value string}`);
`TagSet` (L98) with `has` (L125), `Lookup` (L113), `Len` (L122); `Row` (L168)
with `Outcome`, `Match`, `RequiresOwned`, `Escape`, `Writes`;
`ProvenanceOwned` (L34); `ProvenanceRecognized` (L39); `Table.models` (L216).

`internal/resolve/fixtures_test.go` — `recognizedTagSensitiveTable` (L226);
`fixtureGuards.Evaluate` (L20) with signature
`Evaluate(guard string, _ resolve.TagSet) resolve.GuardResult` — the view
**is** discarded, confirming A2's test-gap claim and scenario 3's net-new
status.

`internal/resolve/resolve_test.go` —
`TestReq17_OwnedObservedAndRecognizedTagsAllReachSelection` (L614);
`mustResolve` (L14) `t.Fatalf`s on non-nil error (L17-19), confirming A8's
"red test, never silently ignored"; `Resolve(resolve.Input{})` nil-error pin
at L748 inside `TestReq20_EmptyInputTupleStillYieldsAValueDisposition`.

`internal/cli/respond/respond.go` — the reserved-but-unenforced
`"ok"`/`"failed"` terminal type names (doc comment L22-25); `OK` (L112) sets
`s.Type = "ok"` unconditionally. The Decision Rationale's sibling-path
evidence for how Enforcement-locus candidate (a) ages is accurate.

**A8/A6 property re-confirmed independently**: `Resolve`'s body has four
`return` statements and every one returns a literal `nil` error; the file
contains no `errors.` or `fmt.Errorf` usage. The RDR's claim that the
reserved error return is unexercised surface holds at HEAD.

Peer-RDR and fixture anchors also resolve: RDR 0002's validation fence
("including at minimum … ambiguous overlap") and its canonical-examples
sentence; the three recognized-provenance declarations
(`rdr-fixture.toml:25`, `kata-fixture.toml:22`, `guard-fixture.toml:36`) and
**all six** reference sites (`rdr-fixture.toml:59,74,86`,
`kata-fixture.toml:47,57`, `guard-fixture.toml:72`) — A4's corrected
inventory verifies exactly. RDR 0009 carries **zero** `Input` occurrences
(case-sensitive, whole file), confirming A6 and A11's disjointness leg.

No bare `file:line` anchor stands in place of a symbol where a symbol exists.
**C5: no findings.**

## CHECK 6 — Status consistency

No assumption is `Pending` or `Unverified`. Twelve of thirteen are
`Verified`; A12 is `Refuted as stated` and — critically — no settled-fact
prose leans on it: the Prerequisites gate, the Testing Strategy Done-split,
and A1's scope limit all treat the handoff as an **explicit obligation**
rather than an assumed traversal, which is the shape a refuted
discoverability assumption requires.

A11's `Verified (independence); report order remains a Stage 7.1 item by
design` is internally consistent with block 4, which licenses the interim
"either error is conforming" rule rather than asserting a settled order —
the critique iter-2 OPEN item on this exact point is closed in live text.

No checklist box contradicts the gate: the Prerequisites checklist carries
three unchecked items, all correctly unchecked (they gate *implementation*,
not lock), and the gate records them as prerequisites rather than as
satisfied.

**C6: no findings.**

## Verdict

**PASS** — no findings across C1–C6. Proceed to the Finalization Gate's
written responses.
