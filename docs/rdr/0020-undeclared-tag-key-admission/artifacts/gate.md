# Finalization Gate — cli/0020 undeclared-tag-key-admission

Date: 2026-09-04 · Verdict: **READY** — locked Final.

Mechanical pre-sweep: PASS (`evidence/tooling-pass/tooling-pass.md`;
`lint --locking` exit 0, blocking=0). Item 4 (Cross-Cutting Concerns) is
authored in the record and retained there for citation as
`cli/0020:G-cross-cutting`; it is deliberately not copied here.

## Contradiction Check

No contradictions between Research Findings and the Proposed Solution.

The two pull the same direction, which is the unusual feature of this record:
the research did not weigh an open question so much as identify which of three
in-code authorities is wrong. Key Discoveries establish that `parseTags`'s
guarded-lookup comment and `ConformValue`'s doc both promise the zero
`TagDecl` conforms everything, while `canonicalValue`'s `isSet := decl.Kind ==
"set"` plus `looksArray` refuses an undeclared array. The Proposed Solution
adopts the documented promise and names the dissenting arm the bug. A solution
adopting arm (2) or (3) would have contradicted those findings; arm (1) is the
one that does not.

The remaining findings corroborate rather than qualify: an undeclared key is
structurally unreadable (every model-side reference refuses at load under
`unknown_tag`), so carrying its bytes cannot influence selection; a declared
key never carries an empty `Kind`, so the zero-decl test identifies
"undeclared" exactly rather than approximately; and the house already
adjudicated this shape once for the emit namespace (`0010:C3` — "undeclared,
uninterpreted, and compared by exact byte equality", with declarations arriving
later as opt-in in 0024). Carrier-by-default / declare-to-tighten is the stated
principle and C1 is an instance of it, not an exception.

One planned behaviour deserves explicit clearing against a stated principle:
C1 *moves* the empty-value refusal arm ahead of the carrier branch, so the
carrier does not escape every refusal. That is not a contradiction of "pure
carrier" — an empty observed value is indistinguishable from unset, a
grammar-level fact about the value that holds with or without a declaration,
which is why it binds the carrier too. The set-specific empty-value message
stays declared-only because it presupposes a declared kind. The hoist is
behaviour-preserving at the byte level for both affected classes (fixtures F4,
E, G).

## Assumption Verification

Five Critical Assumptions, all **Verified**, all internally consistent.

| ID | Status | Method | Consistency |
|---|---|---|---|
| A1 | Verified | Source Search | Status/Method/Evidence agree; "If wrong" non-empty |
| A2 | Verified | Source Search + Spike | agree; "If wrong" non-empty |
| A3 | Verified | Source Search | agree; "If wrong" non-empty |
| A4 | Verified | Prior Art | agree; "If wrong" non-empty |
| A5 | Verified | Peer RDR | agree; "If wrong" non-empty |

- **No `Docs Only` record.** Nothing blocks on the Docs-Only rule.
- **Nothing `Pending` or `Unverified`**, so no settled-fact prose anywhere in
  the RDR depends on an unverified property. Status consistency holds trivially.
- **No self-referential `Verified` stamp.** The three `Source Search` records
  (A1, A2, A3) anchor exclusively at `internal/…` symbols; none cites the
  record or its artifact directory. A2's spike outputs and A4's research trail
  live under `evidence/`, which is where Spike and Prior Art evidence belongs.
- **Every cited `path::Symbol` resolves on `main`** — 28 source-anchor edges,
  all `resolved: true` (C5 of the sweep).
- **No record proves only an adjacent claim.** Two were checked closely because
  their claims are the load-bearing ones. A1 claims nothing load-bears on the
  undeclared-array refusal, and its evidence is stronger than a bare absence:
  the message literal has one producer and zero asserting tests, the single
  test reaching the arm does so over a *declared* key and asserts code and
  `param` only, and the repo asserts the converse directly ("an undeclared
  observed key still passes"). A2 claims an undeclared value cannot influence
  resolution, and closes the closed-world both ways — five `CatUnknownTag`
  load refusals, plus a reader enumeration over every non-test `m.Tags` read
  finding two shapes neither able to observe an undeclared key's value, plus a
  runtime spike showing a byte-identical selection diff. A2's array-carrier leg
  is honestly scoped as not runtime-verifiable until the change lands (it
  refuses at HEAD) and is routed to MVV steps 2/3 — a named plan, not a gap.
- A4 is the one non-decisive record and says so: if the prior-art alignment
  weakens, the decision still stands on the in-repo anchors. Correctly
  non-load-bearing.

## Scope Verification

The Minimum Viable Validation is **in scope and executed during
implementation**, not deferred. Implementation Plan Step 3 ("pinned tests and
docs") carries it, and the Testing Strategy authors the fixtures it names.

The specific proof is a five-step test over a fixture model declaring one
scalar tag and one set tag with no declaration for `extra`/`extras`:

1. the fixture model itself;
2. `--tag extra=plain --tag extras=["a","b"]` exits 0 and the payload's
   `observed` field carries both values byte-for-byte (the carrier claim);
3. the same invocation *without* the carrier flags selects an identical rule
   and outcome (A2's runtime leg — the carried keys influenced nothing);
4. a **declared** scalar handed an array literal is still refused
   `flow-tag-invalid` "is not set-valued", now truthfully — a red test pinning
   the asymmetry so it is authored rather than incidental;
5. `--tag extra=` (undeclared, empty) is still refused `flow-tag-invalid` "was
   given an empty value" — the one arm the carrier does not escape, pinned so
   the hoist out of `canonicalValue` cannot silently drop it.

Steps 4 and 5 are what make this an MVV rather than a happy-path demo: they
pin the two boundaries the change could plausibly erode, and step 5 is the
direct guard on C1's hoist.

## Proportionality

Right-sized. No section flagged for trimming before lock.

**Contract count — the split test.** This RDR is the sole author of exactly
one independent load-bearing contract: **C1**, the meaning of the zero
`TagDecl` at `--tag` admission. One seam, so nothing to split. The refusal-code
and message consequences are downstream of that single decision, not separate
seams: no new code is minted, `flow-tag-undeclared` does not exist, and the
"is not set-valued" message becomes truthful as an effect of the arm becoming
declared-only. The contract is durable (`contracts_durable=1`,
`contracts_transient=0`), so there is no lifespan disposition to discount.

**Profile re-validation.** The recorded `Profile: mid — the meaning of the zero
`TagDecl` at `--tag` admission; user-facing yes; locks contract` still matches
what the contracts just counted: one contract, user-facing yes (a CLI admission
behaviour callers invoke directly and whose refusal they can branch on), locks
a contract. That is `mid` under the applicability matrix, not `small` (which
requires no user-facing surface). The form is correct — value plus one clause
naming the contract, no matrix or provenance prose left from the template. The
lenses that actually ran agree with the field: grounding, 3amigo, and
repeatability-lite, which is the `mid` row; `--outcome repeatability` returns
`emit.next: none` (variant lite, run-1 and diff written). No lens is owed and
none was skipped on a wrong `small`.

**Length.** 1058 lines is at the upper end for a one-contract Bug Fix, and it
earns it in the two places where the record does real work: C1 is long because
the empty-value hoist has a genuine three-class boundary (undeclared, declared
scalar, declared set) that must be written out or an implementer will
collapse it, and each class is pinned to a named fixture. A2's evidence is
long because the safety argument is a closed-world enumeration, which is only
convincing if the enumeration is shown. Neither is padding. The Desk Trace,
Authority Census, and Illustrative Code sections are compact and each carries
a distinct load.

**Joint-decision fence.** `--outcome fence` emitted `stopped:overlap-uncited`
naming four uncited in-flight pairs. Each was assessed against the peers'
contracts before this lock, and all four are incidental co-mentions, not
shared decisions:

- **0019** (`canonicalValue`, `parseWrites`) — 0019:C1 decides init-state
  semantics and routes *around* `canonicalValue` (seeds come from
  loader-normalized `Model.Initial` via `canonicalSet`; "C1 re-conforms
  nothing"); it cites `parseWrites` as evidence only. 0020 is the sole record
  modifying `canonicalValue`. The nearest contact, an empty scalar
  `note = ""`, is explicitly out of scope in 0019:F5 and routed to RDR 0002,
  and 0020's hoist is byte-preserving for declared scalars. No double-edit.
- **0017** (`loadFindings`) — 0017 decides per-finding `code` identity on
  `clierr.Finding` and cites `loadFindings` as a producer-local slug
  precedent. 0020 mentions it once in Decision Rationale as a precedent
  citation. Neither modifies it.
- **0018** (`CheckInput`) — 0018:C1/C2 modify the kernel seam
  (`CheckValid`-before-`CheckInput` precedence, exported `ErrReservedTagKey`);
  0020:A2 cites `CheckInput` only as evidence that it reads `Observed` for the
  reserved key *name*, never a value. The reserved-key channel precedes
  `--tag` admission and 0020 does not move it.
- **`[initial]` literal (0019, 0021)** — 0019:C1 makes the only normative
  claim about `[initial]` behaviour and 0021:C2 exports the declared
  assignments as a read-only graph field; 0020 mentions `[initial]` twice
  solely as one of five `unknown_tag` load-refusal sites proving an undeclared
  tag is unreadable. No competing claim.

The one genuine joint decision this RDR carries was already fired and homed:
**JC1 → 0025**, home `cli/0020 §Normative Contracts C1 / cli/0025 §Normative
Contracts C1` — mutual tolerance over disjoint rule families at
`accessorTable`, with tag-admission identity homed here and accessor-entry
shape homed in 0025. `joint_check_home` reads `unhomed` because the home is a
paired-clause register the tool cannot resolve to a single element; this
written response is what vouches for it. `rulings_open` is 0.
