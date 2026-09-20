Model: claude-opus-5[1m]

# Repeatability resolve — RDR 0030 (full variant, iteration 1)

Origin ledger: `diff.md` — D1–D6, G1–G4. Every entry exits once.

| Finding | Disposition | Lands on |
| --- | --- | --- |
| D1 site-split ownership | fixed — pin | C1, new paragraph |
| D2 `expand`'s totality | fixed — pin | C1, same paragraph; A16 |
| D3 crossing record payload | fixed — pin (same edit as D1/D2) | C1, same paragraph; A16 |
| D4 clear-collision vs C2 bound precedence | fixed — pin | C1 clear paragraph; A14 |
| D5 width check vs enum member guards | dismissed-with-cite | — |
| D6 lint read-back join key | fixed — single-source | C3; A15; `authority` table |
| G1 widened write record shape | fixed — same edit as D3 | C1; A16 |
| G2 admits shim signature | fixed — pin | C1, new paragraph |
| G3 `intWidth` exported spelling | fixed — pin | C1 width paragraph; A12; `authority` table |
| G4 does `internal/resolve` pre-exist | dismissed-with-cite (+ wording fix) | C1 width paragraph |

## Dismissals

**D5** — vacuous, not silent. C1's kind arm is a disjoint switch: the form
is admitted on `int` with both bounds OR `enum` with a non-empty `domain`.
The width check is reachable only on the `int` arm; the `#`-member and
duplicate-member guards only on the `enum` arm. No declaration reaches
both, so there is no precedence to state. The diff itself recorded that no
divergent output was produced.

**G4** — not an RDR silence. `internal/resolve` pre-exists on `main`
(~20 files; `resolve.go::GuardEvaluator`; `guardcontract.go` present), and
the record already says so in A12 and §Phase 2. run-3 marked a GUESS it
could have closed by widening. The phrase "the `internal/resolve` shim"
was nonetheless tightened to name the package as existing, since a
contracts-only reader was the one who tripped.

## Grounding corrections the gate produced

Two diff claims were wrong on `main` and the fixes reflect the code, not
the diff:

- The three-valued verdict type is `resolve.GuardResult` with members
  `GuardTrue` / `GuardFalse` / `GuardUnevaluable` — NOT `GuardUndecided`,
  which run-2 posited and the diff carried.
- `Row` has NO `Span` field; it has `SourceLocator`, and findings set
  `Span: row.SourceLocator`. C3's "the row's `Span` (`model:rule`)"
  misnamed the row-side field. This is what collapsed D6's fork without a
  consult: `graph-overlap` emits a span for the LEFT row only
  (`Element: right.RuleID` has none), so a `Span` join needs a new payload
  field and a repointed `ruleIDsOf`, where recovery-by-truncation needs
  neither.

## Needs (re)verification (Stage 6 closes these)

- **A16 (NEW, Pending)** — `expand` stays total because `renderWrites`
  resolves every refusable step before the crossing record is built.
  Method: Source Search.
- **A12 (already Pending; plan amended)** — now also verifies that the
  moved width core lands as the EXPORTED `resolve.IntWidth` and that both
  cross-package callers compile against that spelling.
- **A14 (already Pending; plan amended)** — now also verifies the clear
  loop still precedes the per-cell walk, so the collision (not C2's bound)
  is what a rule tripping both reports.
- **A15 (already Pending; plan amended)** — the join key is fixed to the
  recovered authored id; verify both slots recover and `ruleIDsOf` is
  unchanged.

No assumption previously Verified was invalidated. No spike is owed.

## Net-new scope

None. Every acted-on finding traces to a `diff.md` entry; nothing was
charted to a successor.

## Determinacy

`Determinacy: fired` written into §Normative Contracts this pass (C1 step
ordering and identity, C3 identity). The fact now reads `fired`.

## Mini-checks

Already fired and present in the draft from earlier lens passes (five
tables: `authority`, `oracle`, `fidelity`, `disposition`, `trace`). This
pass added no new cue; it amended the `authority` table's width row and
its join-key row to match the pinned contracts.
