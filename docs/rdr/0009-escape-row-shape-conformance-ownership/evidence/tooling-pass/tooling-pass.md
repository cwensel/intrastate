Model: claude-opus-5[1m]

# Tooling Pass — RDR 0009 escape-row shape conformance ownership

Stage 7 mechanical adherence sweep, run before the Finalization Gate's written
responses. Post-mutation regression check over a draft rewritten by five
pre-lock lenses (cove, 3amigo, critique, repeatability) and the Stage 6
reconcile.

Target: `docs/rdr/0009-escape-row-shape-conformance-ownership.md` (2037 lines).
Checked against `$RDR_HOME/TEMPLATE.md`.

## CHECK 1 — Template section coverage

Every Required (spine) section is **Present-substantive**:

| Section | Verdict |
| --- | --- |
| Metadata | Present-substantive (Date, Status, Type, Profile, Priority, Related Issues, Predecessors, Overrides, Seam Lineage) |
| Problem Statement | Present-substantive |
| Critical Assumptions | Present-substantive (A1–A10, all with the four-field Evidence Record) |
| Proposed Solution / Approach | Present-substantive |
| Technical Design | Present-substantive |
| Normative Contracts | Present-substantive (8 `normative` blocks) |
| Alternatives Considered | Present-substantive (Alt 1, Alt 2, Briefly Rejected ×3) |
| Context / Background / Technical Environment | Present-substantive |
| Research Findings / Investigation / Key Discoveries | Present-substantive |
| Trade-offs / Consequences / Risks and Mitigations | Present-substantive |
| Failure Modes | Present-substantive |
| Implementation Plan / Prerequisites | Present-substantive (2 checked, 5 open forward obligations) |
| Minimum Viable Validation | Present-substantive (3 normative fixtures, stepwise with expected end-state) |
| Phase 1 / Phase 2 / Phase 3 | Present-substantive |
| Validation / Testing Strategy | Present-substantive (scenarios 1–11) |
| Finalization Gate | Template body pre-lock; replaced by the gate.md pointer at lock (per TEMPLATE's own instruction) — not a hollow finding for a locking RDR |
| References | Present-substantive (18 entries, fully populated from citations the body carries) |

Conditional sections present and substantive: Load-Bearing Decisions
(Selection/predicate, Naming), Capability Dependencies (4 rows), Existing
Infrastructure Audit (5 rows), Decision Rationale (scored 4-column matrix +
`Premortem:` and `Joint-check:` verdict lines), Performance Expectations.

Conditional sections cleanly omitted (PASS per the false-positive guard, not
Missing): Cluster (RDR stands alone), Round-Trip / Inverse Invariants (no
encode/decode pair), Illustrative Code, Day 2 Operations (no persistent
resource), New Dependencies (`errors`/`fmt` are stdlib, recorded in Phase 1).

Placeholder scan — zero hits repo-wide on the RDR for `[Conditional`,
`[Resource]`, `[Capability]`, `_Draft placeholder._`, `this is a seed
skeleton`, `TBD`. No surviving template instructional text outside the
Finalization Gate body.

**C1: no findings.**

## CHECK 2 — Method label vocabulary

Ten Evidence Records, ten Method labels, all exactly in the sanctioned set:

- Source Search ×7 — A1, A2, A5, A7, A8, A9, A10
- Spike ×1 — A3
- Peer RDR ×1 — A4
- Prior Art ×1 — A6

No paraphrased, missing, or off-vocabulary label. Status values are all
`Verified` (four carry a qualifying clause — A2 "with a named wiring
obligation", A4 "settled closed", A6 "contrary to the original claim", A7
"structured form adopted", A9 "with a named amendment obligation" — each a
narrowing of a Verified verdict, not a different status). No `Pending`, no
`Unverified`.

**C2: no findings.**

## CHECK 3 — Source Search self-reference

All seven Source Search records resolve to product source, never to
`{RDR_PATH}` or this RDR's artifact directory:

- A1 → `internal/resolve/resolve.go::Resolve`
- A2 → `internal/cli/clierr::ExitCodeFor`, `internal/cli/config/config.go::Load`
- A5 → `internal/resolve/*_test.go`, `cmd/intrastate/main.go`, `internal/cli/root.go::NewRootCmd`
- A7 → `internal/resolve/resolve.go::Row`, `internal/cli/clierr::CLIError`
- A8 → `internal/resolve/resolve.go::compareRefs`, `::rowRefs`, frozen ADV/Fixup tests
- A9 → `internal/cli/clierr/clierr.go::CLIError`, `::CLIError.Cause`, `go list`
- A10 → `internal/resolve/resolve.go::Resolve`, `::Input`

A3 (Spike) cites `evidence/spikes/a3-fixture-conformance.md` plus four captured
`.out` files — a spike artifact path, which is the sanctioned form for Method
Spike, not a self-reference.

**C3: no findings.**

## CHECK 4 — Docs Only on load-bearing claims

Zero `Docs Only` records. The only two occurrences of the string in the file
are inside the Finalization Gate's own template instructional text.

**C4: no findings.**

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

Every cited `path::Symbol` grepped against the working tree. All resolve in
the cited file:

| Anchor | Resolves |
| --- | --- |
| `internal/resolve/resolve.go::Resolve` | `resolve.go:318` — `func Resolve(in Input) (Result, error)` (confirms A10's arity claim) |
| `::Input` | `resolve.go:224` (carries `Table Table`) |
| `::Row` | `resolve.go:169` |
| `::Table` | `resolve.go:204` |
| `::Plan` | `resolve.go:238` |
| `::Result` / `::Result.Refused` | `resolve.go:283` / `:291` |
| `::RowRef` | `resolve.go:276` |
| `::rescues` | `resolve.go:199` — method on `Row` |
| `::planOf` | `resolve.go:502` |
| `::escapeOrRefuse` | `resolve.go:473` |
| `::gate` | `resolve.go:376` |
| `::assemble` | `resolve.go:148` |
| `::refuse` | `resolve.go:515` |
| `::copyTags` | `resolve.go:569` |
| `::compareRefs` | `resolve.go:528` |
| `::rowRefs` | `resolve.go:538` |
| `::TagSet.Lookup` | `resolve.go:113` |
| `internal/cli/clierr::ExitCodeFor` | `clierr.go:110`; `case GroupUserEnv, GroupInternal: return 2` at `:119` |
| `::GroupInternal` / `::GroupUserEnv` | `clierr.go:39` / `:33` |
| `::CLIError` | `clierr.go:48` |
| `::CLIError.Cause` | `clierr.go:66`, `json:"-"` |
| `::CLIError.Unwrap` | `clierr.go:81` |
| `::ErrorCode` | `clierr.go:90` |
| `::EmitJSON` | `clierr.go:131` |
| `internal/cli/root.go::NewRootCmd` | `root.go:36`; registers only `newVersionCmd()` at `:55` |
| `::ExecuteAndEmit` | `root.go:76` |
| `::cobraErrorToCLIError` | `root.go:117` |
| `internal/cli/config/config.go::Load` | `config.go:74`; `Code: "config-read-error"` at `:85` |
| `internal/resolve/fixtures_test.go::escapeRow` | `fixtures_test.go:249` |

Line-referenced fixture sites re-confirmed (these are cited *as* line numbers
because the line is the edit target Phase 2 acts on, the sanctioned exception):

- `fixtures_test.go` builder sets `Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}}` — present, and `NextTags` carries the same pair (matching the RDR's "NextTags stays on the builder" instruction).
- `adversarial_test.go:229` — `escape.Writes = []resolve.Tag{{Key: "status", Value: "Escaped"}}` on `escapeRow("rdr.escape.needsowned", "flows/rdr.toml:90", …)`. Both the override and the RowRef identity used in MVV fixture 1 confirmed verbatim.
- `fixup_test.go:103` — the same override on `escapeRow("rdr.escape.ambiguous", "flows/rdr.toml:99", …)`.

`Table.CheckValid`, `ErrEscapeShapeBreach`, and `EscapeShapeBreachError`
resolve **nowhere** — correct, and not a C5 finding: these are the surfaces
this RDR introduces, pinned by its Normative Contracts, not claims about
shipped code.

Non-finding recorded per the anchor doctrine: the RDR's escape-partition claim
(`len(row.Escape) != 0`) matches `resolve.go:326` and its placement claim
(check precedes `view := assemble(in)`) matches `resolve.go:319`, though
neither line number is asserted by the RDR.

**C5: no findings.**

## CHECK 6 — Status consistency

No assumption is `Pending` or `Unverified`, so no settled-fact prose can lean
on an unsettled assumption — the check is vacuous by construction here.

Checklist-vs-gate: Prerequisites' two checked boxes agree with the assumption
records (all ten terminal; A3's spike output captured in `{SPIKE_DIR}`). The
five unchecked boxes are forward obligations binding *other* RDRs' implement
stages (0002, 0004) and the future `flow` verb — correctly unchecked, and none
is a precondition of this RDR's own MVV.

Cross-checked the four qualified-Verified records against the prose that
depends on them:

- A2's "no RDR 0005 contract change" verdict is relied on in Normative
  Contracts and Consequences — A9's Stage 6 search closed the additive
  question, and the source confirms the `CLIError` doc comment's `omitempty`
  extension allowance verbatim (`clierr.go:44–47`).
- A4's closed verdict ("predicate does not widen to `NextTags`") is relied on
  in the first normative block and Load-Bearing Decisions — both restate the
  closure with the peer-contract dependency carried forward, consistent.
- A6's falsified original claim is correctly reflected in Key Discoveries
  ("Researched and falsified") and in the Consequence note that the approach
  rests on in-repo authority. No prose still asserts the falsified "load-time
  rejection is the standard" reading.
- A9's amendment obligation appears as a Prerequisite and inside the CLI
  normative block — the `Cause` doc comment's now-stale second clause is
  present at `clierr.go:63–65` exactly as quoted.

**C6: no findings.**

## Cluster re-entry note

No `## Refinement Context (cluster re-entry)` block. Grep: 0 hits. This RDR
never entered the 07.1 gate.

## Determinacy requirement (Profile `foundational`, contract names identity/ordering)

Satisfied. The Normative Contracts name multi-breach report ordering, identity
collapse, and `RowRef` identity — squarely in the class requiring the
repeatability lens. `evidence/repeatability/` carries `run-1.md`, `run-2.md`,
`run-3.md` (full foundational variant) and `diff.md`, whose verdict records one
admissible contract silence (D-1, closed as A10), one test-coverage gap (D-2,
closed as scenario 10's mixed-identity requirement), and one already-open
deferred contract (C-1/A9, closed at Stage 6).

## Verdict

**PASS — no findings. Proceed to the Finalization Gate's written responses.**
