Model: claude-opus-5[1m]

# Tooling Pass — RDR 0006 (mechanical adherence sweep)

Run as the Stage 7 mechanical pre-step, after the Stage 6 reconcile
(`755ba52` / `9a0d894`). Post-mutation regression sweep: six pre-lock rounds and
the reconcile rewrote this draft, so every check below is asking whether one of
them hollowed a section or disturbed an evidence record.

Subject: `docs/rdr/0006-graph-lint-authority-and-guarantees.md` (1858 lines).
Iteration 1 (no prior report in this folder — no loop-breaker comparison owed).

## CHECK 1 — Template section coverage

Every **Required** spine section in `TEMPLATE.md` is Present-substantive:
Metadata, Problem Statement, Critical Assumptions, Proposed Solution
(Approach / Technical Design / Capability Dependencies / Existing
Infrastructure Audit / Decision Rationale), Alternatives Considered, Context,
Research Findings, Trade-offs, Implementation Plan, Validation, Finalization
Gate, References. No section is Missing-bodied; every heading carries body text
before the next heading. `## References` is filled (11 entries) and all eleven
resolve — see CHECK 5.

No surviving verbatim template bracket: `[Required`, `[Conditional`,
`[Resource]`, `[Capability]`, `[NUMBER]`, `[TITLE]`, `[Name]`,
`_Draft placeholder._`, seed-skeleton header, and bare `TBD` all return zero
hits.

Two blocks of **older-template instructional text** survive. Both are
MECHANICAL, both fixed in this pass:

- **C1a — `## Finalization Gate` section body (lines 1718-1729).** The section's
  own body is the template blockquote ("Complete each item with a written
  response before marking this RDR as **Final** … First run the mechanical
  pre-sweep"). Its five `###` children (Contradiction Check, Assumption
  Verification, Scope Verification, Cross-Cutting Concerns, Proportionality)
  are all Present-substantive, so this is a hollow *parent* over substantive
  children, not a hollow leaf. **Fixed by the lock itself**: the READY action
  replaces the whole section body — instruction blockquote and all five
  sub-sections — with the single pointer line to `artifacts/gate.md`. No
  separate edit owed.

- **C1b — "Method vocabulary" block (lines 395-436).** Legacy template
  instructional text inside `## Critical Assumptions`: the eight-label glossary
  ("pick exactly one per assumption"), the self-reference paragraph, and the
  exactness-claim paragraph. The current `TEMPLATE.md` no longer ships it.
  **Not fixed, and deliberately not**: it is present in RDRs 0001, 0002, 0003,
  and 0005 as well, two of which (0003, 0005) are `Final [locked]`. Stripping it
  from 0006 alone would make this document the sole cluster member without the
  glossary its peers cite, and we never amend locked RDRs. Recorded as
  cluster-wide template drift, not an 0006 defect. Does not block: it is inert
  reference prose in an otherwise fully-authored section, and CHECK 2 confirms
  every actual Method label conforms to the vocabulary it describes.

**FALSE-POSITIVE GUARD** — no Conditional section fires. `#### Round-Trip /
Inverse Invariants` (1131-1135) is a substantive-by-negation declaration that
names the owner (RDR 0002) and routes determinism to the finding-identity
decision, not an omission.

## CHECK 2 — Method-label vocabulary

Ten Evidence Records, A1-A10, sequential with no gaps. Every Method is exactly
one sanctioned label:

| Label | Count | Records |
| --- | --- | --- |
| Peer RDR | 7 | A1, A2, A3, A6, A7, A9, A10 |
| Source Search | 2 | A4, A8 |
| MVV Test | 1 | A5 |

None missing, none paraphrased, none off-vocabulary. Records added during the
rounds (A8, A9, A10 — booked at critique iter-2) each carry a conforming label.
**PASS.**

## CHECK 3 — Source Search self-reference

Two `Source Search` records. Neither resolves to `{RDR_PATH}` or to any path
under this RDR's artifact directory:

- **A4** → `internal/cli/clierr/clierr.go::CLIError`, `::ExitCodeFor`,
  `internal/cli/respond/respond.go::Fail`, `internal/cli/root.go::ExecuteAndEmit`.
- **A8** → the repository `.toml` census (`.roborev.toml`, `.kata.toml`, spike
  fixtures under `docs/rdr/*/evidence/spikes/`) — other RDRs' evidence dirs, not
  this one's, and cited as the *absence* proof A8 asserts.

**PASS.**

## CHECK 4 — Docs Only on load-bearing claims

Zero `Docs Only` records. **PASS, vacuously.**

## CHECK 5 — Symbol resolution of anchors

Zero bare `file:line` anchors; zero peer-RDR `~line N` anchors. Every
`path::Symbol` resolves in the working tree:

| Anchor | Result |
| --- | --- |
| `internal/cli/clierr/clierr.go::CLIError` | resolves |
| `internal/cli/clierr/clierr.go::ExitCodeFor` | resolves |
| `internal/cli/respond/respond.go::Fail` | resolves |
| `internal/cli/respond/respond.go::OK` | resolves |
| `internal/cli/respond/respond.go::Success` | resolves |
| `internal/cli/root.go::ExecuteAndEmit` | resolves |
| `internal/cli/root.go::NewRootCmd` | resolves |
| `internal/cli/version.go::newVersionCmd` | resolves |
| `internal/resolve/resolve.go::Resolve` | resolves |
| `internal/resolve/resolve.go::assemble` | resolves |
| `Makefile::build`, `Makefile::check` | resolve |
| `.github/workflows/ci.yml::jobs` | resolves; `jobs.lint` is `name: Lint` (ci.yml:26-27), confirming A5's "key is taken, MUST NOT be reused" |

Peer-RDR anchors `0003::A10`, `::A12`, `::A17`, `::A18`, `::A19`, `::A20`,
`::A21` all resolve as records in
`0003-guard-predicate-exhaustiveness.md`. All eleven `## References` entries
resolve as files/dirs, and the JDR anchors §JD-4, §JD-13, §JD-14 resolve in
`docs/jdr/0001-resolve-kernel-seam.md`.

**PASS.** (The `path::Symbol` occurrence at line 434 is the glossary's own
notation example, not an anchor — see C1b.)

## CHECK 6 — Status consistency

Six `Pending` records (A5-A10), each carrying a Stage 6 DOWNGRADED disposition.
Delegated scan for settled-fact prose depending on them, then verified at
source. **No finding survives verification.** The candidates and why each is a
non-finding:

- **Normative/contract prose specifying a requirement on a pending artifact is
  not a settled-fact claim.** Invariant 5's root passage (604-616) reads "A6's
  *requested shape* … must therefore cover every always-present owned key";
  invariant 6's "the declared initial owned state counts as a write" (622) and
  the reachability relation's "the root is the declared initial owned state
  (A6)" (1037) both state what the relation requires, with the `(A6)`
  cross-reference carrying the pendency. The `Illustrative Code` input contract
  (1139-1146) is explicitly labelled normative and A9 records `--flow` as the
  deferred arm with `--model <path>` instantiable. Specifying a contract over a
  pending declaration is the RDR's job; asserting the declaration already
  exists would be the defect, and no passage does.
- **"The verified assumptions imply a production-command test matrix" (1494)**
  then enumerates A1-A7 by what each *supplies* to the matrix. Read against
  1505-1512, which states the tested code paths "none exists yet", the sentence
  is naming obligations the matrix inherits, not claiming A5-A7 are verified.
  Thin wording, not a status contradiction.
- **Scenario 23 / Consequences on the checked-in model (1357-1359, 1690-1691)**
  read present-tense, but scenario 23 is an MVV scenario — a description of what
  the test does when it runs — and A8 plus `Prerequisites` (1409-1411) both
  record that its subject does not exist yet. The MVV is the named plan A8's
  disposition points to.
- **Checklist-vs-record (b):** the two `A5` boxes (1402 checked, 1404 unchecked)
  are two different claims — "gate *surface* verified" (the Stage 4 source
  search, genuinely done) versus "production gate *proof* pending scenario 7".
  Not a contradiction. The checked box at 1422 is on *RDR 0005 command
  placement*, which is settled, and its own body then states A9's narrowing.
  No box asserts a Pending record satisfied.

**PASS.**

## CHECK 9 — Evidence-field budget (advisory)

Longest Evidence fields, in lines: A7 21, A2 16, A5 15, A6 12, A3 9, A10 8,
A8 7, A4 7, A1 6, A9 4.

**0 fields over the 30-line budget, 0 lines.** No report question owed at the
Gate.

## Determinacy trigger (Stage 7 prompt, not a tooling-pass CHECK)

`Profile: large` and the Normative Contracts name deterministic finding order
and a canonical (never hashed) sortable fingerprint (975-984), so the trigger
**fires** — the Stage 6 reconcile read it as not firing on a keyword scan for
"determinacy", which is narrower than the prompt's topic list.

Disposition: **satisfied by written MVV disposition, no repeatability run
owed.** Determinacy here is a property of the *specified output*, not of the
RDR-authoring process the repeatability lens re-runs. It carries a mechanical
oracle in Validation scenario 4 — "JSON finding order is deterministic by
finding identity, asserted as full-list equality against a golden file, not
merely as a sorted property" — plus the identity tuple at 1029-1033 and the
`Cross-Cutting Concerns` scoping that deterministic claims mean semantic finding
identity, never byte-identical output or content-addressed hashes. Recorded in
`artifacts/gate.md` §3.

## Verdict

**PASS** — 0 blocking findings. One mechanical C1 item (the Finalization Gate
instruction body) is cleared by the lock action itself; one (the legacy Method
vocabulary block) is cluster-wide template drift held deliberately, not an 0006
defect. Proceed to the Gate's written responses.
