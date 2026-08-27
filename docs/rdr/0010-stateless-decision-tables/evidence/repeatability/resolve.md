Model: claude-opus-5[1m]

# repeatability — resolve, 0010 (stateless decision tables)

Origin ledger = the five admissible findings in `diff.md` §1 (D1–D5). The diff
itself dropped six false GUESSes in its §3 and two decomposition clusters (G3,
G4) in §2; those are not re-adjudicated here.

Diff health: three runs on three distinct models (`claude-sonnet-5`,
`claude-fable-5`, `claude-haiku-4-5-20251001`), GUESS clusters localized to
C1's class surface, disagreements to three contracts. Healthy per the lens's
Expected signal — not the identical-but-confidently-wrong shape that would
call for a model switch.

## Grounding gate

All four code-facing claims were grounded on `main` before any edit.

| Claim | Verdict | Anchor |
| --- | --- | --- |
| Gate step sits after selection, before payload; deny is a refusal | CONFIRMED | `internal/cli/flow_exec.go::runGates`, `::gateVerdictFailure`; `flow-gate-denied` |
| `resolvePayload` field count and order | CONFIRMED — **thirteen**, `Gates` between `Rule` and `Next` | `internal/cli/flow_resolve.go::resolvePayload` |
| `escaped` already a payload field, position 12 of 13 | CONFIRMED | same struct, between `Clear` and `EscapeClass` |
| Escape rescue sets `Plan.RuleID` to the escape row | CONFIRMED | `internal/resolve/resolve.go::planOf` via `::escapeOrRefuse` |
| `Model` has a `Class` field/accessor today | ABSENT (as expected) | `internal/table/model.go::Model` |
| `Model`'s existing convention | **exported fields only**, no field accessor | `internal/table/model.go::Model` (13 exported fields, one projection method `::KernelTable`) |
| Loader step order; `checkAccessorBindings` last | CONFIRMED last; `loadAccessors` 4th, before `normalizeRules` | `internal/table/load.go::run` |

Computed, not argued: the `resolvePayload` field list was enumerated from
source (`Model, Revision, Observed, Owned, Readers, Outcome, Rule, Gates,
Next, Writes, Clear, Escaped, EscapeClass` = 13), which is what refuted C4's
position clause rather than merely questioning it.

## Dispositions

**D1 — `0010:C4`, gate placement / thirteen-field count → fixed (pin + correct a false claim).**
The diff read this as an RDR silence. Grounding shows the *behaviour* is fixed
in code and needs only restating — but it also caught a real defect the diff
did not claim: C4 declared `emit` "immediately after `Rule`", and `Rule`'s
actual successor is `Gates`. C4 now declares `emit` after `Gates`, states the
count as thirteen→**fourteen**, and records that a gate deny is a refusal so
`emit` is never computed on a denied selection. run-3's `gates` payload field
was not an invention; it is already there.

**D2 — `0010:C4`, `escaped` and the escape row's own emit → fixed (single-source).**
Not a new rule: A9 already grounds `escaped` and the escape path, but only in
assumption prose. C4 now states that `escaped` keeps its declared position and
that a rescued plan joins the **escape row's own** `emit` on the same
`Plan.RuleID` path — so C3's per-rule join needs no escape-specific arm. run-1's
omission was the RDR's silence, correctly caught.

**D3 — `0010:C1`, class storage shape → fixed, with the diff's framing inverted.**
The silence is real: the accessor + zero-value clause lived only in the
Technical Design lead-in and the *Authority* cell, outside every contract span.
But the diff called run-3's exported field the defect ("defeats the one-writer
guarantee"); grounding refutes that — `table.Model` is a plain exported-field
record with **no** field accessor anywhere, so run-3 read the type correctly and
an accessor would be its sole exception. C1 now pins the exported `Class string`
field with an empty zero value reading as `state-machine`, and states plainly
that the one-writer/four-reader ownership is a **review obligation, not a
mechanized one** — which *Authority* already said for the fifth-site risk. This
is the pass's strongest fix; it is also where the diff would have led us wrong.

**D4 — `0010:C1`, agreement-check ceiling → fixed (pin both bounds).**
C1 fixed only the `loadTags` floor. run-3's upper bound is correct against
source: `checkAccessorBindings` is `run`'s last step. C1 now fixes the window at
both ends and states the consequence the ceiling buys — the doubly-malformed
decision table (owned tag *and* `[initial]`) refuses as `malformed model
declaration` naming `owned=<n>`, never with the writer-arity diagnostic C2 flags
as "not self-explanatory". The window is fixed; the step's *name* is not, which
dissolves cluster G2.

**D5 — `0010:C3`, `Emit` aliasing through `expand` → fixed (leave-shared, stated).**
Admitted by the diff as its weakest finding, and it holds: C3 reasons about
`TagValue`'s clone machinery as the reason for a new type, then leaves the copy
rule for that type open. Resolved toward **sharing permitted, no clone
required** — requiring a clone would reintroduce exactly the machinery
`EmitValue` exists to avoid. Stated so an implementer does not add a defensive
copy and read it as contract.

## Amendment sweep

C1's storage change is a redefinition of a subject token (`accessor` → `Class`
field), so every site was swept, not just the contract:

- `§technical-design` lead-in ("gains the class as a field with an accessor") → exported field wording.
- `§technical-design` item 1 ("reads the accessor") → "reads the `Class` field (empty = `state-machine`, C1)".
- `§authority` model-class row → exported `Model.Class`; its writer cell also gains the ceiling.
- `A10` ("converting either site to the class accessor") → "the class test" — consistency only; A10's verdict is untouched.
- `S5` checked and **not** edited: it asserts payload shape (`{}` never `null`), never a field count, so the thirteen→fourteen correction does not reach it.

## Needs (re)verification — carried to Stage 6

- **A15 (NEW, Pending)** — exported `Model.Class` with empty-zero-value semantics
  is sufficient for the four class readers, and the C1 window has a step
  available. Method: Source Search. Partially grounded already (field
  convention + loader step slice confirmed); what remains is that each of the
  four readers can reach the field where it runs.
- No previously Verified assumption was invalidated. A9 is reinforced, not
  disturbed — C4 now states normatively what A9 had only in prose.

## Charted to successor

None. Every finding traced to a ledger entry; no net-new scope was absorbed.

## Tiebreakers

None escalated. D3 was the one genuine fork (field vs. accessor) and the
grounding evidence collapsed it — the existing type convention decides it, so
it did not need a §strong-consult or a human call.
