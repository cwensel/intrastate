Model: claude-opus-5[1m]

# Cove Dispositions — RDR 0002, iteration 2

Origin ledger: `iter-2/findings.md` CV2-001 … CV2-004. One line per finding.
Delta-scoped to the cluster-reconcile iter-3 re-verify set (A1, A9, A12) and the
§D7 re-authoring; iteration 1's ledger (CV-001…CV-021) is closed and was not
re-opened.

## Fixed

- **fixed** — CV2-001 (the merge key joins set-literal members, collapsing
  distinct atoms). The literals clause was correct but bound only the *stored*
  normalized value; the spike keeps `Literal []string` as required and still
  loses the atom by keying its merge map on `strings.Join(a.Literal, ",")`.
  Extended the clause to bind **every comparison of an atom — merge key,
  dedup, sort, any identity** — to the member sequence, and to say why any
  delimiter is wrong here rather than only a delimiter a member might contain
  (no delimiter is banned from a tag value; `#` is refused only in rule ids,
  alphabet members, and match-block `in` members). Sections: Normative
  Contracts / literals.
  **Amendment sweep** (§amendment-sweep) over the three dependent sites:
  - **A13** demoted `Verified` → `Pending` — its two witnesses use scalar
    literals, under which a joined rendering and a member sequence are the same
    string, so they establish "not keyed on `(block, key, operator)`" but not
    "keyed on the full identity." Recorded what the spike actually does and what
    witness closes it. Section: Critical Assumptions A13.
  - **Prerequisites** — the verified-except list and the "A7, A11, A13, A14
    promoted" sentence were stale; A13 moved to the `Pending` set with its
    reason. Section: Implementation Plan / Prerequisites.
  - **Testing Strategy 2** — added the delimiter-bearing distinct-literal
    control (`["a,b","c"]` vs `["a","b,c"]` → two atoms, rule dead for lint) and
    stated why the existing space-bearing control does not discriminate; the
    idempotence-mirror paragraph's "promote a captured result" claim was
    narrowed to the `(block, key, operator)` question it actually settles.
  - **Round-Trip / fidelity** — the lossy-site note claimed the `#`-ban made the
    renderings unambiguous; the ban never reaches atom literals. Named the
    set-valued atom literal as non-recoverable in the rendered form and scoped
    the `#` argument to row identity. Sections: Round-Trip; `fidelity` table.
  - **desk trace** — added the missing `merge atoms` row carrying the
    contradiction, and corrected the closing line from "No CONTRADICTION row
    survives" to name the one that does.
  - **`oracle` mini-check** — recorded the merge control as the fourth
    non-discriminating oracle, with the shared pattern named.

- **fixed** — CV2-002 (three enumerated categories have no fixture). Scenario 3
  already disclosed the gap and named all three among "five still owed," so the
  contract needed no change — but two of the three are owed only a *fixture*,
  not an implementation, which the draft did not say. Confirmed both fire
  against mutated fixtures this pass (`cyclic context inheritance at "prelock"`;
  `malformed tag declaration: "cluster_ready" has no kind`) and recorded that
  `gen-cases.py` should mint them rather than leaving them to implementation
  discovery. Section: Testing Strategy 3.

- **fixed** — CV2-003 (`[model.metadata]` exemption had undefined edges). Stated
  the three unstated consequences: it is model-level and reaches no candidate
  row, so it is outside the dump field list and outside the Round-Trip
  invariant (two documents differing only in it normalize to identical row sets
  — intended, the invariant is over rows); a dump MAY render it as presentation;
  its internal shape is deliberately unconstrained, because constraining it
  would make it a schema and defeat the point of an extension namespace.
  Section: Normative Contracts / layout.

- **fixed** — CV2-004 (fail-fast order-unobservability argument was unscoped).
  Bounded it to defects that reach a check, and named the absorbed-defect case
  (trips zero categories, so its one-defect fixture asserts nothing) as a
  normalization defect the merge and literals clauses exist to prevent, not a
  case the fail-fast freedom licenses. CV2-001 is the concrete instance.
  Section: Normative Contracts / version gate.

## Charted to successor

None. Every finding landed inside this RDR's existing contract surface — the
literals clause, the category list, the layout clause, and the fail-fast clause
are all already this RDR's. No net-new scope was absorbed.

## Dismissed with cite

None. All four findings passed the grounding gate: CV2-001 and CV2-002 are
source-grounded with reproductions against the committed spike binary; CV2-003
and CV2-004 cite exact draft passages. None re-litigates an adjudicated
decision — in particular CV2-001 is **not** a re-raise of the contexts clause's
prefix-key prohibition, which forbids dropping a *field* from the key; this is a
lossy rendering of the last field, a route that clause does not close.

## Needs (re)verification — carried to Stage 6

1. **A13 is now `Pending`** (was `Verified`). Verification: the delimiter-bearing
   distinct-literal control in Testing Strategy 2 — two contexts contributing
   `in = ["a,b","c"]` and `in = ["a","b,c"]` on one key and block must yield two
   atoms, and the rule must be a dead rule for lint. The spike as committed fails
   this; the fixture and a spike fix are owed. This is the only assumption this
   pass moved.
2. **The extended literals clause is a new load-bearing normative claim** — it
   binds every atom comparison, not just the stored value. It needs checking
   against RDR 0003 (atom identity tuple is theirs downstream) and RDR 0006
   (which serializes atoms and may compare them) to confirm neither depends on a
   joined rendering.
3. **The `[model.metadata]` round-trip exclusion** (CV2-003) is a new statement
   about the Round-Trip invariant's extent. It needs checking against RDR 0006,
   which consumes the normalized model and may expect metadata to participate.
4. **`gen-cases.py` must mint two new fixtures** (CV2-002) — `cyclic context
   inheritance` and `malformed tag declaration`. Mechanical; the checks already
   exist and were confirmed firing.
5. **The spike's `output.txt` digest is unaffected by this pass** — no fixture
   was changed, only the draft. The digest in Testing Strategy 2 and the
   Performance Expectations digests remain valid. Recorded so Stage 6 does not
   re-check them needlessly.

No previously `Verified` assumption other than A13 was invalidated. A1, A9, and
A12 — the cluster-reconcile re-verify set — are all confirmed at source this
pass: A1's §D7 layout claims hold against the re-authored fixture under strict
decoding, and A9's and A12's `Pending` status and blocking reason (`resolve.Row`
carries no atom shape) are CONFIRMED verbatim in `internal/resolve/resolve.go`.

## Charted

None this iteration.

## Tiebreakers

None. The one fork that could have needed adjudication — whether CV2-001 is a
spike bug outside the draft or a contract gap — collapsed on evidence: the spike
*conforms to the letter* of the literals clause (it does keep `Literal` a
sequence in the normalized value) and still commits the defect, so a conforming
implementation can commit it and the clause is what needs tightening. No
`§strong-consult` was needed.

## Mini-checks

Fired cues re-read this pass, per the Stage 5 mini-check obligation: `fidelity`,
`disposition`, `oracle`, `trace`. Three of the four were edited by this pass
(`fidelity`, `oracle`, `trace`); `disposition` was re-read and needed no change —
the arity split and its `duplicate model id` exception are unaffected by these
findings. Tables persist in the draft; no new cue fired.
