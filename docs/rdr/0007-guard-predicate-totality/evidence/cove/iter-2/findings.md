Model: claude-opus-5[1m]

# Cove (S-04) — RDR 0007, iteration 2 (re-entry, delta-scoped)

Chain-of-Verification, delta-scoped to the Status line's re-verify set:
**A3, A16, A17, A18, A19, A21, A22**. Two isolated sub-agents — Step 0
(grounding sweep) and Steps 1–2 (12 independent verification questions) —
plus main-context grounding of every finding below against source.

Element files: `step-0-grounding.md`, `steps-1-2-questions.md`.

Iteration 1 (loose files in the parent dir) ran pre-Final on the whole draft;
its findings are not re-reported here. This pass covers the post-Final revision.

## Origin ledger

| ID | Origin | Class | One-line |
| --- | --- | --- | --- |
| F-1 | Step 0 #1 + Steps 1–2 Q2 | REFUTED (contradiction with a closed joint decision) | A19 reverses JDR 0001 §D4's CLI-transport ruling without acknowledging it |
| F-2 | Q6 | REFUTED (spike validated a different design) | The A3 spike kept `Refusal.Guard`; the RDR says the payload replaces it |
| F-3 | Q11 | silence (unlisted prerequisite) | Final RDR 0009 asserts on `Refusal.Guard`, which this RDR deletes; no re-lock box |
| F-4 | Step 0 #2 | inconsistency | A16 Verified and A22 Pending rest on the *same* unfulfilled 0002 duty-list state |
| F-5 | Q7 | silence | The nil-seam branch's meaning after the per-atom reshape is unstated |
| F-6 | Step 0 #3 | evidence over-scope | A21 quotes 0009's producer obligation as if it were a kernel-boundary guarantee |
| F-7 | Q9 | under-specified signature | Set-valued tag encoding is unstated, so Phase 3's `contains` contract test cannot be written |
| F-8 | Q12 | silence | A19 is the only Pending assumption with no fallback, yet is marked "Blocks lock" |
| F-9 | Step 0 #4 | anchor hygiene | A3/Testing-Strategy cite proposed symbols with no proposed-vs-shipped marker |
| F-10 | Q3 | understated blast radius | The nested payload is the kernel's first non-flat `Refusal` field, filed as "Reuse + extend" |
| F-11 | Q8 + Step 0 #5 | minor false fact / truncated quote | "closes the set at five"; A17's `exists` quote drops "of an optional tag value" |
| F-12 | main context | silence (behavior delta) | Shipped `gate` selects one `lowest` row for the refusal; the per-row payload retires that, unstated |
| F-13 | main context | silence (doc obligation) | `gate`'s doc comment carries D8's superseded rationale the RDR forbids repeating |

---

## F-1 [REFUTED — highest value] A19 reverses a CLOSED joint decision without acknowledging it

A19 states the structured-envelope-field question is open and blocked on a
§JD-8 amendment, and Prerequisites makes that grant **"Blocks lock."**

JDR 0001 §D4 already ruled on the transport, in the same paragraph that assigns
this RDR its normative home:

> …replacing the single-valued `Refusal.Guard` text, which has no referent once
> guards are atoms; **it reaches the CLI through §JD-8's `Detail`**.
> — `docs/jdr/0001-resolve-kernel-seam.md:206–207`

§JD-12 records §D4 as **Closed** (`docs/jdr/0001-resolve-kernel-seam.md:275`).
§JD-8 is unamended (L252–260) and pre-authorizes `Code` values only, explicitly
"**not new envelope fields or exit groups**."

So the author's 2026-08-21 round chose the option §D4 had already rejected —
structured field over flattened `Detail` — and A19 presents it as an *addition*
to §JD-8 rather than a **reversal inside a closed decision**. The RDR's own
Overrides section claims no supersession. Two isolated sub-agents converged on
this passage independently.

This is a genuine either/or the evidence cannot collapse: reversing §D4 is a
joint-decision reopening, not a drafting fix. See Tiebreaker T-1.

## F-2 [REFUTED] The A3 spike validated a design the RDR does not specify

Normative Contracts: the payload "replaces the single-valued `Refusal.Guard`
text." §D4 agrees. But the spike's own change list reads:

> `Refusal.Guard string` → `Guard []GuardAtom` **plus** `UndecidedAtoms
> []GuardAtom` (sorted, …)
> — `evidence/spikes/a3-reshape-probe.md:156`

The spike **kept** `Refusal.Guard` (retyped) and added the payload beside it.
Its PASS verdict for `fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`
depends on that survival — it re-encodes `!=` to `reflect.DeepEqual` over a
field that, under "replace," does not exist. Under the RDR as written that test
must be **re-decided**, not re-encoded, which is exactly A3's own "If wrong"
branch ("the reshape is a coordinated migration, not a free change").

A3 is currently `Verified — spike executed`. The spike does not verify the
specified design; it verifies a superset. Either the contract says "replaces"
and the spike's PASS is not evidence for it, or the design keeps both fields
and the contract is wrong.

## F-3 [silence] Final RDR 0009 asserts on `Refusal.Guard`; no re-lock prerequisite

RDR 0009 (`Final`) grounds its own non-vacuity on that field:

> …the asserted properties are `Refusal.Kind` / `Refusal.Guard`
> — `docs/rdr/0009-escape-row-shape-conformance-ownership.md:421`

This RDR deletes `Refusal.Guard`, staling a Final peer's evidence. Prerequisites
carries re-lock checkboxes for RDR 0002 and RDR 0005 but **none for 0009**, and
A21 — which cites 0009 for its `Writes`-empty premise — does not mention it.

## F-4 [inconsistency] A16 Verified and A22 Pending rest on the same unfulfilled state

A22 is Pending because the canonicalization duty "is NOT currently on 0002's
Refinement Context Direction list." A16 is `Verified` on the strength of 0002's
joint-check line. But 0002's Direction list carries **neither** duty:

> Name the `RequiresOwned` producer (JD-3) and rename the three canonical
> fixtures (JD-10). … Restate the resolver flow as gate-then-count … Add `row
> kind` and `escape failure classes` …
> — `docs/rdr/0002-transition-table-as-reviewable-data.md:861–864`

The existence-token emission duty (A16) is as absent from that list as
canonicalization (A22). Same evidentiary state, two different verdicts.

## F-5 [silence] The nil-seam branch after the per-atom reshape

`resolve.go::evaluateGuard` returns `GuardUnevaluable` when `seam == nil` for
any non-empty guard. After the reshape the kernel decides existence atoms and
absent-key atoms itself, so a row of all-`exists` atoms — or all absent keys —
would be **fully decidable with no seam at all**. The RDR is silent on whether
the nil seam still forces unevaluable.

This is load-bearing: frozen `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`
and `Input.Guards`' shipped doc both depend on the current answer, and A3's
touch-list does not name the branch.

## F-6 [evidence over-scope] A21 quotes a producer obligation as a kernel guarantee

A21's evidence quotes 0009's clause bare. Source heads it:

> Escape-row shape conformance … **is a PRODUCER obligation on every constructor
> of `resolve.Row` values**: a Row with a non-empty Escape list MUST have an
> empty Writes slice.
> — `docs/rdr/0009-escape-row-shape-conformance-ownership.md:822–826`

The composition A21 derives (`Writes` empty ⇒ `RequiresOwned` empty) therefore
holds only where a *conforming producer* built the row — it is not a kernel-side
invariant. Chained onto §JD-3, which records that no layer is yet obliged to
populate `RequiresOwned` at all, the derivation has no guaranteed implementer.
A21's "Verified (composition)" qualifier is honest but the quote's scope is not.

## F-7 [under-specified signature] Set-valued tags have no stated encoding

`resolve.go::Tag.Value` is a bare `string`. A set-valued tag therefore reaches
the seam as one string with no declared element encoding, yet:

- the contracts pin `contains` over an **absent** set-valued tag as unevaluable
  (Testing Strategy row 8), and
- Phase 3 exports a `contains` contract test for RDR 0003 to instantiate.

`Evaluate(atom, value) GuardResult` over an unstructured string is insufficient
to write that test, and no document states the encoding. A17 claims the seam
needs nothing beyond atom and value; for `contains` that is unsupported.

## F-8 [silence] A19 has no fallback, yet blocks lock

A22 carries an explicit fallback (narrow to kernel-side only, re-file the duty
on 0002). A19 — marked **"Blocks lock"** — states no path if the grant is
declined. Given F-1, decline is the likelier branch: §D4 already chose `Detail`.

## F-9 [anchor hygiene] Proposed symbols cited without a proposed-vs-shipped marker

`GuardAtom`, `evaluateAtom`, `undecidedAtoms`, `compareAtoms`, `copyAtoms` have
**zero** hits under `internal/` (re-run this pass). They are proposed names, which
is legitimate — but Testing Strategy row 17's evidence column cites
`undecidedAtoms`/`compareAtoms` beside genuinely shipped anchors with no marker
distinguishing them. Shipped analogues: `resolve.go:528 compareRefs`,
`resolve.go:569 copyTags`. A reviewer re-running the sweep reads five NOT-FOUNDs.

## F-10 [understated blast radius] First nested `Refusal` payload

The Existing Infrastructure Audit files the payload as "Reuse + extend."
`Refusal` today is flat scalars plus `[]RowRef` / `[]string`. A two-level
array-of-objects (per row → per atom) is the kernel's **first** nested refusal
payload, which is what makes the A19 envelope question hard in the first place.

## F-11 [minor] Two inexact quotes

(a) A16: "RDR 0003's vocabulary clause closes the set at five" — it closes at
five operator **classes**, which is eight tokens (`lt`/`lte`/`gt`/`gte`).
Immaterial to the design, wrong as written.

(b) A17 quotes 0003 as "existence checks presence"; source reads "existence
checks presence **of an optional tag value**"
(`docs/rdr/0003-guard-predicate-exhaustiveness.md:358–361`). The elision drops
the optionality that is the whole point of the operator.

## F-12 [silence] The `lowest`-row selection is retired but unstated

Shipped `gate` builds the unevaluable refusal from ONE row:

```go
lowest := slices.MinFunc(undecidable, func(a, b Row) int {
    return compareRefs(refOf(a), refOf(b))
})
return nil, &Refusal{Kind: KindGuardUnevaluable, Guard: lowest.Guard, Rows: rowRefs(undecidable)}
```
— `internal/resolve/resolve.go:408–416`

The RDR's payload is per-row over **every** undecidable row, sorted. That
retires the `MinFunc` selection — but the spike carries it forward
(`undecidedAtoms(lowest.Guard, seam, view)`, spike L240), and the draft never
says the single-row choice goes away. F-2's ambiguity is what leaves room for it.

## F-13 [silence] `gate`'s doc comment repeats the superseded D8 rationale

The RDR's normative block: "Authoring docs MUST NOT repeat the superseded
rationale." The shipped kernel comment does exactly that:

> absent owned state is the more precise diagnosis and **is frequently the
> reason the seam could not decide the predicate**
> — `internal/resolve/resolve.go:370–372`

No phase names retiring it. Phase 2 pins the precedence but not the prose.

---

## Verdict

**Healthy pass**, and a materially sharper one than iteration 1 — the delta
scope paid: five of thirteen findings land directly on the seven re-verify
assumptions, and the two highest-value ones (F-1, F-2) both say the same thing
in different registers: **the post-Final revision drifted from the joint
decision and from its own spike.**

A3's re-verification trigger was re-run this pass and is **CLEAN** (no
production importer; 154 PASS / 0 FAIL), so that half of A3 stands.

The pass's centre of gravity is F-1. Everything else is resolvable in-draft;
F-1 is a reversal inside a closed joint decision and is a genuine fork.
