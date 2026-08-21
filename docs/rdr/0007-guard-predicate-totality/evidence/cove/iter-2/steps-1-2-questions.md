Model: claude-opus-5[1m]

# Cove Steps 1–2 — verification questions (iteration 2, delta-scoped)

Delta scope: A3, A16, A17, A18, A19, A21, A22 and the contracts/design/plan/testing
passages that depend on them. Iteration 1 findings read only to avoid re-reporting;
no answer below inherits iteration 1's conclusions.

Question kinds: **[SRC]** answerable only by reading source; **[RDR]** RDR-internal;
**[PEER]** peer-document.

---

## Q1 [SRC] — Does `internal/cli/clierr.CLIError` carry any structured (non-string) serialized field today, and is `Detail` a string?

**A. No structured field; `Detail` is a `string`.** Read
`internal/cli/clierr/clierr.go::CLIError` (L48–67). Every serialized field is a
plain `string`:

```go
type CLIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Param string `json:"param,omitempty"`
	// Detail carries hard facts about the failure (the underlying
	// syscall reason, a parser diagnostic). May be multi-line.
	Detail string `json:"detail,omitempty"`
	Hint string `json:"hint,omitempty"`
	Group ErrorGroup `json:"-"`
	Cause error `json:"-"`
}
```

`Group` and `Cause` carry `json:"-"`. `EmitJSON` (L131) marshals the struct
directly, and `EmitText` (L142) renders `Detail` as a single `detail: %s` line.
A19's refutation claim ("carries only flat strings — `Code`, `Message`, `Param`,
`Detail`, `Hint`") is **confirmed against source**, including the doc-comment
invitation it quotes ("Extend with new optional fields as needed — keep them
`omitempty`", L46–47).

---

## Q2 [PEER] — Does JDR 0001 §D4 actually leave the payload transport open, as A19 and the Prerequisites assume, or does §D4 already decide it?

**A. §D4 already decides it — and decides it the OTHER way. This is a live
contradiction with A19 and with Prerequisite 3.** `docs/jdr/0001-resolve-kernel-seam.md`
§D4's resolution paragraph (L203–208):

> The `guard_unevaluable` refusal names what blocked the verdict per row and per
> atom — key, block, reason `absent` | `uncomparable` — replacing the
> single-valued `Refusal.Guard` text, which has no referent once guards are
> atoms; **it reaches the CLI through §JD-8's `Detail`.**

§JD-8 (L253–257) says exactly what A19 quotes: resolver-specific values "only
require new `Code` constants or literals, **not new envelope fields or exit
groups**", and "`Detail` and `Hint` ship and render in text and JSON."

So the JDR — the very document A19 names as the place the field "cannot be
stated by this RDR" until amended — has already routed this payload through
`Detail`, i.e. flattened text. A19 states the opposite ("2026-08-21: structured
field, not flattened text") and asserts §JD-8 "does NOT cover this payload."
Under the JDR as written, §D4 + §JD-8 jointly DO cover it, via `Detail`. A19
therefore does not merely await a decision; it contradicts a decision already
recorded, and the RDR nowhere acknowledges that §D4 says `Detail`.

---

## Q3 [SRC] — Does the shipped `resolve.go` contain a `RowRef`-keyed or key-keyed structured refusal payload precedent that a new per-atom payload would parallel, or is `Refusal` entirely flat-plus-slices today?

**A. `Refusal` is flat scalars plus two flat slices; there is no nested-object
precedent in the kernel.** `internal/resolve/resolve.go::Refusal` (L255–273):
`Kind RefusalKind`, `Revision string`, `Flow string`, `Recognized string`,
`Rows []RowRef`, `MissingOwned []string`, `Guard string`. `RowRef` (L276–279) is
`{RuleID, SourceLocator string}` — a flat two-string struct, which is the
closest existing analogue and is only one level deep. The RDR's payload contract
("a two-level array-of-objects (per undecidable row its `(RuleID, SourceLocator)`;
per unevaluable atom its key, block, and reason)") has **no** shipped precedent
in `Refusal`. Nothing refutes the design; but the RDR's Existing Infrastructure
Audit row "Undecidability verdict + refusal mapping … Reuse + extend" understates
that the extension introduces the kernel's first two-level refusal payload.

---

## Q4 [SRC] — Would adding new exported kernel symbols (`GuardAtom`, existence-token constants) break any shipped boundary test that pins the kernel's exported surface?

**A. No.** `internal/resolve/boundary_test.go::exportedKernelSymbols` (L92) returns
every exported top-level name, but the three call sites all use **ban-lists, not
allow-lists**: `resolve_test.go:895` (banned symbols implying work/persistence),
`resolve_test.go:1186` (banned encode/decode names), and `resolve_test.go:1208`
(`banned := []string{"Hash", "Checksum", "Canonicalize", "Fingerprint", "Digest"}`).
Adding `GuardAtom`, `OpExists`-style constants, or boolean-literal constants
trips none of these — **unless** the exported existence-literal constants were
ever named with a banned token. Note the near-miss: REQ-37 bans an exported
symbol literally named `Canonicalize`, and A22's normative clause discusses
canonicalization (assigned to 0002, not the kernel), so no collision arises. A3's
"bounded" claim survives this check.

---

## Q5 [SRC] — Does `resolve.go::gate` read guard TEXT for any control-flow decision, such that replacing `Row.Guard string` with an atom slice could perturb the prune → owned → undecidable partition?

**A. No — `gate` touches `row.Guard` at exactly two sites, neither of which is a
branch condition.** `resolve.go::gate` (L376–420): L380
`verdict := evaluateGuard(seam, row.Guard, view)` — pass-through; L414
`Guard: lowest.Guard` — payload copy. Every branch is on the `GuardResult`
(`if verdict == GuardFalse`, L381; `switch verdicts[i]`, L398). `missingOwned`
(L444) reads only `row.RequiresOwned`; `escapeOrRefuse` (L473–498) never mentions
`Guard` and delegates entirely to `gate` (L485). A3's "partition unchanged /
`escapeOrRefuse` zero changes" is **structurally confirmed**, independent of the
spike's empirical claim.

---

## Q6 [SRC] — Do the two shipped `Refusal.Guard` reader assertions really avoid guard-text identity, as A3's "If wrong" condition requires?

**A. Yes, both.** `fixup_test.go:64` is
`if got.Refusal.Guard != escape.Guard` — a comparison against the value the test
itself set on the escape row, i.e. a round-trip identity, not a text assertion.
`resolve_test.go:1091` is `if r.Guard == ""` with message "refusal names no
undecidable guard for diagnosis" — a non-emptiness check. Neither depends on the
guard's *representation*. A3's "If wrong" branch does not fire.

**But note a gap A3 does not name:** `resolve_test.go:1091`'s `r.Guard == ""`
re-encodes to `len(r.Guard) == 0` **only if `Refusal.Guard` survives as a slice**.
The Technical Design says "`Refusal.Guard` is replaced by the per-row, per-atom
payload" and the normative block says the payload "replaces the single-valued
`Refusal.Guard` text" — i.e. the field is *deleted*, not retyped. The spike
(§4 "The reshape applied") kept BOTH: `Refusal.Guard []GuardAtom` **plus**
`UndecidedAtoms []GuardAtom`. So the spike that verifies A3 did not implement the
RDR's own stated design; it implemented a retype-plus-add. Under the RDR's actual
"replace" wording, `fixup_test.go:64`'s `escape.Guard` comparison has no field to
compare against and must be re-decided, not re-encoded.

---

## Q7 [SRC] — Does `resolve.go::evaluateGuard`'s empty-guard branch short-circuit before the seam, so an unguarded row is outside the domain rule?

**A. Yes.** `resolve.go::evaluateGuard` (L425–433):

```go
func evaluateGuard(seam GuardEvaluator, guard string, view TagSet) GuardResult {
	if guard == "" {
		return GuardTrue
	}
	if seam == nil {
		return GuardUnevaluable
	}
	return seam.Evaluate(guard, view)
}
```

The `guard == ""` branch precedes the nil-seam check, so an unguarded row is
`GuardTrue` without any view consultation. The RDR's normative block "A row with
no atoms is decided TRUE without consulting the view (shipped: `evaluateGuard`'s
empty-guard branch)" is **accurate as to shipped behavior**.

Under-specified, though: the RDR nowhere states what happens to the **nil-seam**
branch after the reshape. Once the seam is per-atom and is never called for
absent keys or `exists` atoms, a nil seam no longer blocks a row whose atoms are
all `exists` or all absent-key — the kernel can decide it alone. The shipped
contract (`Input.Guards` doc, L230–232: "A nil seam means the table's guarded
rows cannot be decided") and the frozen test
`fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` both rest
on nil-seam ⇒ unevaluable for *guarded* rows. The RDR's contracts are silent on
whether a nil seam still makes an all-`exists` row unevaluable, or whether the
kernel now decides it. That is a shipped frozen contract the reshape can flip,
and A3 does not name it.

---

## Q8 [PEER] — Does RDR 0003's operator/kind matrix give the kernel exactly one presence token, as A16 claims, with no second absence-inspecting operator?

**A. Yes, one token.** `docs/rdr/0003-guard-predicate-exhaustiveness.md` L234–240,
the full matrix:

| Operator | Accepted tag value kinds | Literal shape | Lint proof role |
| `eq` | enum, boolean, integer, string-like scalar | one typed scalar | Narrows the tag domain to one value. |
| `in` | … | non-empty typed scalar set | Narrows the tag domain to the listed values. |
| `lt`, `lte`, `gt`, `gte` | integer | one typed integer | Narrows a bounded integer domain by comparison … |
| `exists` | optional scalar or optional set-valued tag | boolean | Tests presence or absence, not value equality. |
| `contains` | set-valued tag with a declared element universe | non-empty typed element set | Narrows the set-valued domain … |

`exists` alone is stated in presence terms; the vocabulary clause (L296–298)
closes the set: "The initial operator vocabulary MUST be closed and typed …
Unknown operators MUST be rejected during parse or lint before resolution."
A16's premise holds.

**Caveat A16 overstates:** the matrix row `lt, lte, gt, gte` is **one row
carrying four tokens**. A16 says "The other four rows (`eq`, `in`,
`lt`/`lte`/`gt`/`gte`, `contains`) are stated purely as domain-narrowing value
comparisons, and RDR 0003's vocabulary clause closes the set at five." The
vocabulary clause closes the set at five *operator classes*, which is **eight
tokens**. That does not refute A16 (the kernel needs only to recognize one token
and treat all others as value atoms), but "closes the set at five" is a
mischaracterization: the kernel's fail-closed rule (Contract 3: "an existence
atom carrying a foreign token is a value atom to the kernel") means the kernel
never needs a token whitelist — so the count is irrelevant, and stating it as
five is simply wrong.

---

## Q9 [PEER] — Does RDR 0003 actually establish that declared kind/bounds are lint-only inputs and never evaluation inputs, as A17's "seam needs neither view, provenance, siblings, nor declarations" requires?

**A. Partly — the bounds half is quoted correctly, the `contains` half is not
supported.** The bounds sentence A17 cites is verbatim at 0003 L278–281: "A guard
may still compare an unbounded integer at runtime, but lint must report that it
cannot prove exhaustive coverage for that dimension" ⇒ bounds are indeed a lint
input. And the semantic-identity tuple is `(tag, operator, literal)` with no kind
field (0003 L351–352).

But A17 also asserts "`contains`'s declared element universe is likewise a
lint-proof input." 0003's matrix (L240) says `contains` accepts a "set-valued tag
**with a declared element universe**" and 0003 L275 says "set-element universes
are finite by declaration; set-valued tags must be [declared]". More decisively:
the kernel's `Tag.Value` is a bare `string` (`resolve.go` L20–23:
`type Tag struct { Key string; Value string }`), so a set-valued tag reaches the
seam as ONE string. To decide `contains`, the seam must know how to **decompose
that string into elements** — a representation fact 0003 never states, and which
A17's "needs neither … nor tag declarations at evaluation time" would forbid it
from taking from the declaration. A17 is **under-evidenced on `contains`**: the
seam signature `Evaluate(atom, value)` is only sufficient if set encoding is
carried in-band in the value string, which no document states.

This propagates: Testing Strategy Scenario 8 (`contains` over an ABSENT
set-valued tag ⇒ unevaluable) is decidable kernel-side, so it survives; but
Phase 3's exported contract test ("present value × literal × operator") cannot be
written for `contains` without the encoding fact.

---

## Q10 [PEER] — Does RDR 0003's only parse-rejection clause really cover the AUTHORED LITERAL only, leaving unparseable RUNTIME values open, as A18 claims?

**A. Yes — A18's characterization is verbatim-accurate.** 0003 L302–303:
"Each operator MUST declare which tag value kinds it accepts. A predicate whose
literal cannot be parsed as the declared tag kind MUST be rejected before
resolution." The subject is the **literal**, not the runtime value. 0003's Testing
Strategy scenario 4 (L748–749) lists "unknown tag, unknown operator, unsupported
operator/tag-kind pair, and **literal parse mismatch**" — all parse-time,
literal-side. A grep for a runtime-value-parse clause in 0003 returns nothing.
The gap A18 anticipates is real, and A18's supporting source cite is also
accurate: `resolve.go` L56–58 documents `KindGuardUnevaluable` as "the guard seam
reported a predicate it could not decide." A18 verifies clean.

---

## Q11 [PEER] — Does RDR 0009 (Final) support A21's composition, and does it cite `Refusal.Guard` in a way A21 or the RDR must account for?

**A. The `Writes`-empty premise is verbatim; the `Refusal.Guard` citation is real
and A21 does not mention it.** 0009's Normative Contracts (L820–828):

> Escape-row shape conformance — an escape row carries no owned-state mutation —
> is a PRODUCER obligation on every constructor of resolve.Row values: a Row with
> a non-empty Escape list MUST have an empty Writes slice.

Verbatim as A21 quotes, and 0009 is `Status: Final [joint decision → JDR 0001
§JD-5]` (L9). A21's kernel-side half checks out too: `resolve.go::missingOwned`
(L444–458) ranges `for _, key := range row.RequiresOwned`, so an empty slice
never enters the loop and `missing` stays nil.

**What A21 omits:** 0009 L421 cites `Refusal.Guard` as a load-bearing property of
its own non-vacuity argument — "the asserted properties are `Refusal.Kind` /
`Refusal.Guard`". This RDR deletes `Refusal.Guard`. The RDR's Joint-check
paragraph does record this ("the per-atom payload replacing `Refusal.Guard`
(which Final RDR 0009 cites; stale, rides to its re-lock per §JD-12)"), and
JDR §JD-12 confirms it — so the RDR is not silent overall. But A21, which is the
assumption that reads 0009 for this RDR, does not carry it, and the Prerequisites
list has **no checkbox** for 0009's re-lock, unlike the checkboxes it does carry
for 0002 (§JD-3, A22) and 0005 (A19). A Final peer whose evidence cites a field
this RDR deletes is an unlisted prerequisite.

---

## Q12 [RDR] — Does the RDR state a consistent status for A19, or do the Metadata Status line, A19's own Status, and the Prerequisites disagree?

**A. They disagree with each other and with the peer document, three ways.**

1. The Metadata Status line (L9–10) reads "Draft [revised from Final 2026-08-12 —
   re-verify A3, A16, A17, A18, A19, A21, A22]" — seven assumptions.
2. Prerequisites (L1406–1408) reads: "A3, A16–A19, A21, A22 verified (Resolve).
   A3, A16, A17, A18, A21 verified 2026-08-21; **A19 and A22 remain Pending**."
   "A16–A19" is a range that includes A17 and A18 — consistent — but the Status
   line's set and this line's set agree, so (1)/(2) are compatible.
3. A19's own Status (L669–678) says "Pending — blocked on a peer-document change,
   not on evidence" and names the plan "raise the field at JDR 0001 as a §JD-8
   amendment (or at RDR 0005's re-lock)". Prerequisite 3 (L1399–1405) escalates
   this to "**Blocks lock**, not just implementation."

The contradiction is against the peer document, not internal (see Q2): JDR 0001
§D4 already states the payload "reaches the CLI through §JD-8's `Detail`", which
is precisely the flattened-text option A19 records the author as having rejected.
So the RDR asserts a *blocking* prerequisite whose named home has already decided
against the RDR's chosen shape, and the RDR never says so. Either the JDR line is
stale and should be named as such, or A19's "the original 'without a new envelope
field' reading is refuted" is refuting the JDR's own live decision without citing
it.

Separately, the RDR nowhere states what happens if the amendment is declined —
A22 carries an explicit fallback ("narrow A22 to kernel-side only"), A19 carries
none, yet A19 is the one marked "Blocks lock."

---

## Answers that reveal a problem

- **Q2 — contradiction (peer document).** JDR 0001 §D4 already routes the
  per-atom payload to the CLI "through §JD-8's `Detail`" — flattened text. A19
  asserts the opposite direction (structured field) and asserts §JD-8 "does NOT
  cover this payload," without acknowledging that §D4 already answered. The
  prerequisite A19 says "blocks lock" is aimed at a document that has already
  decided the other way.
- **Q6 — false claim / spec-vs-spike divergence.** The RDR says `Refusal.Guard`
  is *replaced*; the A3 spike that verifies A3 kept `Refusal.Guard` (retyped to
  `[]GuardAtom`) **and** added `UndecidedAtoms`. Under the RDR's own "replace"
  wording, `fixup_test.go:64`'s `got.Refusal.Guard != escape.Guard` has no field
  and must be re-decided, not "re-encoded" — the exact condition A3's "If wrong"
  branch describes. A3's Verified status rests on a spike that implemented a
  different design.
- **Q7 — silence / unnamed frozen-contract risk.** The RDR is silent on the
  nil-seam branch after the reshape. With a per-atom value-only seam, a nil seam
  no longer makes an all-`exists` or all-absent-key row unevaluable — yet
  `Input.Guards`' shipped doc ("A nil seam means the table's guarded rows cannot
  be decided") and frozen `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`
  both depend on it. A3's touch-list does not name it.
- **Q8 — false claim (minor).** A16 says RDR 0003's "vocabulary clause closes the
  set at five." It closes at five operator *classes* = eight tokens
  (`lt`/`lte`/`gt`/`gte`). Immaterial to the design (the kernel needs no
  whitelist) but the stated fact is wrong.
- **Q9 — under-specified / under-evidenced.** A17's claim that the seam needs
  only `(atom, value)` is unsupported for `contains`: `resolve.go::Tag.Value` is
  a bare `string`, so a set-valued tag arrives as one string with no stated
  element encoding. Phase 3's exported `contains` contract test cannot be written
  from what any document states.
- **Q11 — unlisted prerequisite.** RDR 0009 (Final) cites `Refusal.Guard` as a
  load-bearing asserted property (0009 L421); this RDR deletes the field. The
  Joint-check paragraph notes the staleness, but A21 (the assumption that reads
  0009) omits it and Prerequisites carries no 0009 re-lock checkbox, unlike the
  0002 and 0005 boxes it does carry.
- **Q12 — internal inconsistency around A19.** A19 is the only Pending
  assumption marked "Blocks lock" and the only one with no stated fallback if the
  named peer-document change is declined (contrast A22's explicit narrow-to-kernel
  fallback).
- **Q3 — understated blast radius (minor).** The Existing Infrastructure Audit
  files the payload as "Reuse + extend"; `Refusal` today is flat scalars plus
  `[]RowRef`/`[]string`, so the two-level array-of-objects is the kernel's first
  nested refusal payload, with no shipped precedent.

Clean (no problem surfaced): **Q1** (clierr flat-strings claim confirmed),
**Q4** (exported-symbol tests are ban-lists; A3 bounded), **Q5** (`gate` never
branches on guard text; partition-unchanged is structural), **Q10** (A18's read
of 0003's literal-only parse clause is verbatim-accurate).
