Model: claude-opus-5[1m]

# Grounding Findings — iteration 2

Re-entry sweep. Iteration 1 (`../findings.md`) predates six doc commits that
demoted this RDR from Final to Draft and re-authored it (+430/-181 since
`ff133bb`). No `Ground-sweep:` verdict in Decision Rationale, so scope was
everything, not a propose-diff.

Three sweeps: local Go/build/CI symbols (17 claims), peer-RDR clause citations
(32 claims), sibling prior-art and JDR citations (6 claims).

## REFUTED

- **F1 — `Technical Design` invariant 5 attributed the view-level conformance
  check to RDR 0007, which never accepted it.** The draft read "the view-level
  check at assembly is RDR 0007's". RDR 0007 states no conformance obligation:
  its only `conform` content is **escape-row shape** conformance
  (`0007`, "escape-row shape conformance ... a row that breaches the conformance
  predicate above MUST ..."), a different subject. RDR 0003 A18
  (`Status: Pending`) records the premise as **"currently held by no
  component"** and lists RDR 0007's adoption as a *planned request*, not an
  existing duty. Confirmed against source: `internal/resolve/resolve.go::assemble`
  merges owned/observed/recognized tags into a `TagSet` and reads no
  declaration; `conform` occurs nowhere in the non-test kernel.
  A locked peer was cited for a duty it does not carry.

## NOT-FOUND

- **F2 — "RDR 0001's kernel is stateless and single-step" cited `single-step` as
  peer text.** `stateless` is normative and repeated in RDR 0001
  (`Technical Environment`, `Approach`, `Cross-Cutting`). `single-step` appears
  **nowhere** in RDR 0001 — nor do "one step", "multi-step", or "traversal".
  The inference holds against RDR 0001's text (the kernel "returns either one
  legal transition plan or a typed refusal" per supplied snapshot, and does not
  run subflows), but the citation read stronger than its source.

## Internal inconsistency (not a codebase refutation)

- **F3 — the five tag-declaration fields were partitioned two different ways
  inside this RDR.** A2 listed "value kind, finite domain, optionality,
  single-valued marker, element universe" (five separate fields, matching RDR
  0003's clause); `Technical Design` folded element universe *inside* finite
  domain and reached five by a different split. Both are compatible with RDR
  0003; the RDR disagreed with itself.

## Peer ledger drift

- **F4 — RDR 0003 routes six open records to RDR 0006; this RDR named two.**
  RDR 0003's `Status` line assigns **A10, A12, A17, A19** to RDR 0006's refine
  and **A18, A20** across RDR 0006 and RDR 0007. This draft claimed only "closes
  RDR 0003 A10 and A12". Checked each remaining done-condition against the
  current text — all four are substantively discharged but unnamed, so the
  peer's ledger cannot close by citation:
  - **A17** ("cluster confirms escape-row overlap is checked among escape rows
    for one failure class") — invariant 3 adopts exactly the two-population
    reading. A17's own text says what it "still owes is the cluster's
    confirmation ... not its mechanics".
  - **A18** ("one named component owns the check") — invariant 5 is the
    model-level producer; A18 names two sufficient producers and requires one.
  - **A19** ("a code or field distinguishes that verdict") — the draft mints
    `graph-coverage-closed-by-escape`. A19's cited gap ("today has neither a
    code nor a field") is stale.
  - **A20** lint half ("the refusing atom is addressable") — the finding
    contract now carries `Key`/`Operator`/`Literal`/`Block`. A20's cited gap
    (`0006:363-367` has no atom-level field) is stale.

## Inverse sweep — new rule vs. existing sibling

Searched `internal/`, `cmd/` (case-insensitive) for `graph`, `lint`,
`invariant`, `reachab`, `traverse`, `dataflow`, `worklist`, `BFS`, `DFS`,
`visit`, `predicate`, `Finding`, `Diagnostic`. **Searched, none exists** — every
hit is a prose comment, never code. No existing path computes reachability,
graph traversal, or a stable finding identity.

Nearest siblings, recorded because they are precedent rather than duplication:
`internal/resolve/resolve.go::RowRef` (`{RuleID, SourceLocator}`) with
`::compareRefs` is a *row* identity that makes refusal payloads order-independent
— not a finding identity (no code, severity, or invariant id), but the same
ordering discipline this RDR's finding-identity tuple needs.
`::RefusalKind`/`::RefusalKinds` is the existing stable-code enumeration pattern.

## Confirmed, with recorded imprecision

- **Kernel escape consultation.** "consults escape rows only for `no_match` /
  `ambiguous_match`" is true of the **call sites** — `escapeOrRefuse` is invoked
  only from `Resolve`'s `case 0:` and `default:` arms — but is **not enforced by
  the `Row` type**: `Escape []RefusalKind` may name any of the five
  `RefusalKinds()`. Gating inside `escapeOrRefuse` is by `row.rescues(r.Kind)`.
  The RDR's per-class partition is correctly general and unaffected; RDR 0003
  states the restriction normatively ("An `escape` list MUST contain only ...
  `no_match` and `ambiguous_match`", `0002::Normative Contracts`), so the
  alphabet is closed upstream rather than by the kernel type. No edit.
- **`{"type":"failed"}` envelope.** `respond.go`'s package comment claims a
  `type`-discriminated failure record that `respond.Fail` does not emit
  (`clierr.EmitJSON` marshals the bare `CLIError`); `docs/cli-output-contract.md`
  is correct and the code comment is stale. This RDR never asserts the
  discriminator — it says findings ride "both the error envelope and the success
  payload" — so the design is correctly grounded. No edit; noted because an
  implementer reading `respond.go` would be misled.
- **CI trigger.** `.github/workflows/ci.yml` restricts push to `branches: [main]`
  and adds `workflow_dispatch`; its lint job runs `golangci-lint-action`, not
  `make lint`. A5's claim (discrete jobs, none invoking `make check`) is
  unaffected.
- **`CLIError` field list.** The RDR's enumeration omits `Cause error json:"-"`.
  Non-serialized, so the `omitempty` wire-compatibility claim is unaffected.
- **Aggregation veto subject.** RDR 0007's veto subject is the *resolution*, not
  the row; "a row can refuse" is RDR 0003's defined row-level predicate, which
  this RDR cites explicitly. Faithful inheritance — no edit.
- **Default-on provenance.** "default-on" originated as RDR 0003's *reading of*
  RDR 0006 and an invitation to settle divergence; RDR 0006 now states it, so
  the loop closes in the right direction.

## Confirmed without exception

All remaining claims resolved on `main` or in the cited peer section, most
verbatim. Load-bearing highlights:

- `clierr.go::CLIError` (all six listed fields), `::ExitCodeFor`,
  `GroupUserEnv`+`GroupInternal` sharing exit 2 in one case arm — the RDR's
  "assert on `Code`, not the exit integer" instruction is well-founded.
- Every optional serialized `CLIError` field carries `omitempty` (exhaustive tag
  audit); the package doc states the append-only invariant as intent. The typed
  `findings` field is wire-compatible.
- `respond.go::OK`/`::Fail`/`::ValidateMode`/`::Success` (with `Data any`);
  `root.go::NewRootCmd`/`::ExecuteAndEmit` (conversion verified at the call
  flow, not just the symbol).
- `Makefile::check` = `fmt-check vet lint test`, does **not** depend on `build`;
  `Makefile::build` → `./bin/intrastate` from `./cmd/intrastate`; CI jobs are
  exactly `test`/`lint`/`vuln`, none invoking `make check`.
- Exactly **eight** non-test Go files under `internal/`+`cmd/`; sole registered
  verb is `internal/cli/version.go::newVersionCmd` (one `AddCommand` call).
- The three verbatim RDR 0003 quotes this RDR puts in normative blocks all
  resolve: the narrowing clause, "A declared escape row participates in the
  coverage identity", and "Lint MUST NOT drop `exists` atoms from the product"
  (including the always-present carve-out).
- A6's evidence is exact: `initial` occurs nowhere in RDR 0002, and `terminal`
  only as the fixture rule id `terminal-archive` in its Validation scenario 2.
  RDR 0002's layout enumeration and `[model]` clause match byte-for-byte; RDR
  0002 is `Draft` with A1–A8 all `Verified`; both spike `[model]` blocks carry
  only `id`, `version`, `description`.
- RDR 0005 scopes its normative surface to "one command group" (`flow …`) and
  lists graph lint authority as out of scope, so root `intrastate lint` does not
  collide.
- Prior art: Seed 6 names all six invariants verbatim; `MODEL-transition.md` L1/
  L2/L3 and the provenance→determinism link; `RESOLVER-DESIGN.md` "Run once at
  design time over the **table**"; `repos/README.md` supersedes model-checker
  adoption while explicitly preserving the design-time validator carve-out.
  Alternative 4's tool list (Quint/TLA+/SCXML/sismic) is drawn from
  `attic/TOOL-EVAL-validator.md`, not invented.
- JDR 0001 §JD-4/§JD-13/§JD-14 exist; none contradicts root `intrastate lint`,
  RDR 0006 owning reachability, or the tag model living in RDR 0003.
  **§JD-14's overlap half is itself stale** — RDR 0003 A17 records that it
  reasoned from a premise "the shipped kernel refutes" and that RDR 0006's
  reading "was closer to right than the gate credited". The JDR is the stale
  document here, not this RDR.
