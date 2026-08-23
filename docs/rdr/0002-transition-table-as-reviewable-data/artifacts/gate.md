Model: claude-opus-5[1m]

# Finalization Gate — RDR cli/0002, Transition Table As Reviewable Data

- **Date**: 2026-08-23
- **Verdict**: **PASS — READY. Locked to Final.**
- **Mechanical pre-sweep**: `evidence/tooling-pass/tooling-pass.md` — BLOCK on
  one MECHANICAL C1 finding (surviving template guidance block in
  `## Critical Assumptions`), fixed in-pass, re-run **PASS**.
- **Pre-lock lenses run** (`Profile: foundational` → cove → 3amigo → critique →
  repeatability): all four complete with dispositions; critique dual-model;
  repeatability ×3 across three distinct model stamps, diff health **Healthy**.

## 1. Contradiction Check

**No contradictions found between research findings, design principles, and
proposed solution.** Four candidates were examined and each resolves:

**Research supports the chosen split, it does not cut against it.** The
Investigation's prior art is uniformly *source-model-plus-rendered-view*:
Sismic's `import_from_dict` / `export_to_dict` (nested source model exported to
another view), `transitions`' positive `conditions` / negative `unless` guard
lists, Stateless's `StateGraph` (symbolic graph built from metadata),
Sismic's `PlantUMLExporter` (rendered output as a view over model data). The
Proposed Solution's sparse-TOML-source + expanded-table-dump is that same shape.
The statechart literature hit (hierarchy/extended state as the answer to
dimensional explosion) is honored by inherited contexts, while the stated
principle "not a runtime FSM engine" is honored by keeping the libraries as
vocabulary/validation/visualization only. No research finding recommends a
representation this RDR rejects.

**Exact-one vs. first-match is a deliberate divergence, stated as such.** The
Investigation names first-match FSM engines; Decision Rationale rejects them
explicitly ("priority order is convenient in code, but it makes review harder
and lets source or dump reordering change behavior"), and A2 verifies at source
that the shipped kernel already implements gate-then-count. A named, argued
divergence from surveyed prior art is not a contradiction.

**The "hand-authored, not generated" principle and the dump are compatible.**
Background says the table "must be hand-authored from the legal graph audits,
not generated." The RDR generates the *expanded table*, never the *source
model* — Approach states "the rendered table is review support, not the source
authors maintain," and the Consequences bullet says the same. The generated
artifact is the derived view; the authored artifact stays hand-written.

**The producer-locks-after-consumers inversion is booked, not hidden.** RDR 0009
cites "RDR 0002's **Final** escape-rule prohibition" and RDR 0008 asserts
"nothing in this RDR or RDR 0002 forbids that alphabet entry" about the
empty-string outcome, which this RDR now *does* forbid. The Risks section names
both by hand, classifies them correctly as stale peer **fact** rather than
contract conflict (RDR 0008 explicitly scoped itself out: "not this RDR's to
rule on"), and routes the drift to `/rdr-cluster-reconcile`, which owns
cross-RDR reconciliation. Per project doctrine RDRs are never amended, so
carrying the obligation forward is the correct disposition. Not a lock blocker.

## 2. Assumption Verification

**14 Evidence Records (A1–A14). Every record is internally consistent:** Status,
Method, and Evidence agree, and no `If wrong` is empty.

**No `Docs Only` records** — the vocabulary distribution is `Spike` ×4
(A1, A3, A6, A7), `Source Search` ×3 (A2, A5, A8), `Peer RDR` ×1 (A4),
`MVV Test` ×6 (A9–A14). Nothing blocks on the Docs-Only rule.

**No self-referential `Verified` stamp.** The three `Source Search` records cite
`internal/resolve/` and `internal/cli/` paths; none resolves to this RDR file or
its artifact directory. Each cited symbol resolves in the working tree — `Row`
(`resolve.go:169`), `Resolve` (`:318`), `assemble` (`:148`), `gate` (`:376`),
`escapeOrRefuse` (`:473`), `missingOwned` (`:444`), `Table.models` (`:216`),
`recognizedTagKey` (`:110`), `CLIError` (`clierr.go:48`), `Fail`
(`respond.go:133`), `Load` (`config.go:74`), and
`TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`
(`adversarial_test.go:266`). No bare `file:line` anchor appears in any Evidence
field.

**Each `Verified` stamp proves its own claim, not an adjacent one.** Three worth
naming:

- **A2** does not merely assert order-independence — it cites the kernel's own
  control flow (`gate` returns on `blocked != nil` *before* `switch
  len(selected)`), the sorted `missingOwned` output, the source comment
  delegating row order to this RDR, and the adversarial test that freezes the
  property. This discharges JDR 0001 §D2's "Lands in 0002" clause.
- **A5** is unusually honest and correspondingly stronger: it verifies the
  envelope exists *and* states that no parse-validation code has yet travelled
  it, that `config-invalid` is a doc comment plus a `TODO`, and that this RDR's
  categories would be the first of their kind — so code names, `Group`, and exit
  mapping are new design work rather than a pattern to copy. The stamp covers
  exactly the verified half.
- **A8** verifies the outcome filter at source and then states its own
  **residual** rather than overclaiming: the alphabet ban is a load-time check
  in the normalizer, the kernel imports no normalizer and validates no alphabet
  well-formedness, so the barrier covers exactly one producer. Closing it for
  all producers is RDR 0001's to make. This is the correct shape for a bounded
  guarantee.

**Six records are `Pending` (A9–A14), and every one is a deliberate Stage 6
DOWNGRADE with a named MVV assertion — not an unexamined gap.** Prerequisites
enumerates them individually and Testing Strategy carries each oracle:

| ID | Why Pending | Where the oracle lives |
| --- | --- | --- |
| A9 | Blocked on RDR 0007's kernel reshape — unverifiable before it lands | Scenario 2 normalization assertions + scenario 4 `Resolve` run against one normalized value |
| A10 | Decidability half source-verified (`Model` is a singular struct field); paired-document behavior owed | Scenario 3, a **pair** of documents sharing a `model id` |
| A11 | Critique lens introduced the clause; spike fixture binds with `eq`, so escape expansion has no captured witness | Scenario 2, escape rule binding `in` over two alphabet members → two rows, distinct suffixes, empty write sets |
| A12 | Blocked on RDR 0007's reshape, like A9 | Scenario 2, union of `Match` and guard atoms equals the normalized set and intersection is empty; `status.eq=closed@unless` is the discriminating case |
| A13 | Mirror of the distinct-literal assertion; landed in no assertion until now | Scenario 2, byte-identical atom from rule and inherited context collapses to count one |
| A14 | The existing `unsupported version` fixture is v1-shaped and trips on either ordering, witnessing nothing about precedence | Scenario 3, a **v2-shaped** document asserted to refuse `unsupported version` and not `unknown schema field` |

None is load-bearing for a *pre-lock* MVV assertion; each is survivable per its
`If wrong`. A9 and A12 are additionally sequenced behind RDR 0007's reshape,
which the Prerequisites verify is genuinely unimplemented on `main` (no `Atom`
type, no `OpExists` / `LiteralTrue` / `LiteralFalse`, `Row.Guard` still a
`string`, `Row.Match` still a flat `[]Tag`). That is an implementation-sequencing
gate, not a lock gate, and the RDR says so.

**Status consistency holds.** No settled-fact prose leans on an unsettled
assumption. A11's escape-expansion clause, A12's routing totality, A13's merge
idempotence, and A14's version-gate precedence appear in Normative Contracts as
obligations this RDR is *making* (MUST/MAY), never as verified statements about
existing code, and each is paired with its owed oracle. The RDR's `Status: Draft`
line carries no 07.1 qualifier while the README row still reads
`Draft [revised from Final 2026-08-12; re-verify A2, A7]`; per TEMPLATE.md that
qualifier self-clears at this flip, and the re-verification it named is
discharged in Prerequisites ("A2, A7 re-verified; A8 resolved").

## 3. Scope Verification

**The Minimum Viable Validation is in scope and executed during implementation,
not deferred — with one item explicitly and correctly excluded.**

The specific proof: parse the RDR and kata sparse TOML fixtures into typed
source structs; normalize them into candidate rows; dump the expanded table;
validate the full load-time category set; then prove by unit test that one
sample tag-set resolves through `internal/resolve::Resolve` to exactly one
ordinary row, one unmatched tag-set resolves to exactly one modeled escape row,
and one unsupported-version variant is refused before normalization. Testing
Strategy scenarios 1–6 carry the assertions, and nine load-time categories are
already witnessed one-refusal-per-mutated-fixture in
`evidence/spikes/negative-cases.txt`.

**The one excluded item is properly excluded, not quietly dropped.** The
ambiguous-overlap assertion is cross-row and therefore RDR 0006's by the arity
split, and RDR 0006 is unimplemented. The RDR states it is "deferred to the
Phase 5 lint handshake rather than counted as satisfied here," names what this
RDR still owes before that handshake (overlapping rows distinguishable by row
identity and comparable by predicate set, which the identity tuple and atom-set
contract already fix), and says plainly that "marking the overlap assertion
green before RDR 0006 lands would be asserting on a stub." That is a scoped
handoff with a named obligation, not a deferral of this RDR's own validation.

**The MVV's oracle discipline is unusually strong and worth recording.** The
"every oracle must have a failing control" rail is not boilerplate — it names
three assertions that *previously passed against a wrong implementation* and
replaces each: the category oracle (was asserting message text, now asserts
category), the determinism oracles (a normalizer emitting a constant passed;
now paired with the value assertion and run on the RDR fixture that carries the
escape row, the `in`-expansion, `<clear>`, and inherited contexts), and "the
comparator is total" (was satisfied by the positional tiebreak it existed to
exclude; now rule-order permutation). The existence-atom test exercises both
`exists = true` and `false`; byte-equality gets a negative control (`Status` vs
`status` must fail `unknown tag`); unknown schema fields get a fixture, since
strict decoding is not the TOML library's default and is otherwise a silent
no-op.

**Fixture promotion is `extend-then-promote`, and the reasoning is sound.** The
spike fixtures lack two sibling candidate rows binding the same outcome, without
which the gate's ordering is untestable — the desk trace honestly records "no
witness possible in the current fixture" for that row. The RDR resolves the
conflict in favor of the MVV over the promote-verbatim mandate, and constrains
the escape: a fixture may be extended, never narrowed.

## 4. Cross-Cutting Concerns

**Versioning.** Owned here. `[model].version = 1` is the only accepted version,
and the version gate is normatively ordered **before** strict field decoding —
the only ordering this RDR fixes — so a future v2 document refuses as
`unsupported version` rather than as `unknown schema field` on whichever v2-only
key the decoder happens to reach first. A14 owes the v2-shaped precedence
fixture.

**Canonical form / determinism.** The core concern, and the reason this RDR is
`foundational`. It claims byte-identical dumps, so the sub-clause applies item
by item:

- *Hash function + library* — **deliberately not a contract.** "No hash is part
  of the contract: the SHA is evidence for this spike output only, and
  production golden tests must assert the normalized expanded-table value the
  normative contract defines." The SHA-256
  `c4be7447a241fb724632c53a9e5f39c7a9b5c7cc78e0a1e74ae4132271432a1b` is spike
  evidence, not an identity commitment. Correct scoping — nothing downstream
  content-addresses this.
- *Pre-image byte layout / primitive encodings* — the dump's field list is fixed
  and every field is present; the invariant is stated over the **normalized
  value, not the dump text**, because this RDR fixes the field list and row
  ordering but defines no dump grammar.
- *Map iteration order* — explicitly neutralized: "Go's map iteration is never
  the emission order (rows, atoms, next tags, writes, required-owned keys, and
  escape classes are all sorted slices before rendering)." Witnessed by three
  consecutive byte-identical runs despite Go's randomized iteration.
- *Ordering* — the total row order is the identity tuple
  `(model id, rule id, expansion suffix)` compared field by field,
  byte-lexicographically, with the source locator excluded so unrelated source
  edits cannot churn golden tests. Byte-lexicographic comparison is a Go spec
  guarantee, so the ordering is locale- and platform-independent. Source-order,
  map-iteration, and renderer-specific ordering are each explicitly rejected.
- *Whitespace policy / case folding* — "no case folding, trimming, or namespace
  rewriting is applied at any stage," matching RDR 0008's byte-exact reserved-key
  comparison, with a negative control (`Status` against a declared `status` must
  fail `unknown tag`).
- *Empty / null / absent distinguishability* — the expansion suffix is absent or
  an alphabet member, never the empty string, so the identity tuple is total;
  the empty string as an alphabet member is a load-time refusal; an alphabet
  member or rule id containing the expansion-suffix separator is its own
  distinguishable category, which is what keeps `(rule id, suffix)` recoverable.
- *Version marker for future evolution* — `[model].version`, gated first.

The three named **lossy rendering sites** (`<clear>` unreserved in the tag-value
space, unescaped separators in multi-entry fields, unreserved suffix separator)
are disclosed rather than hidden, and the obligation to close them is charted to
a successor dump-format RDR. Since the contract is the normalized value and the
dump is a review surface, this is a bounded, stated limit — not an unresolved
determinism gap.

**Character encoding.** Tag keys and literals are exact byte strings with no
folding, trimming, or namespace rewriting; set-valued literals render with
members sorted so two authored orderings are one literal. Byte identity is the
rule throughout.

**Incremental adoption.** The table is new data with no existing model to
migrate; `internal/resolve::Row` is reused rather than reshaped by this RDR, and
Phases 2–3 sequence behind RDR 0007's reshape rather than forcing it.

**Build tool compatibility.** `github.com/pelletier/go-toml/v2` (v2.3.1, MIT) is
a candidate only; no production dependency is added until implementation.
Strict decoding is stated as an obligation on the *format*, not a property of
the library, so a parser swap cannot silently retire the `unknown schema field`
category.

Peer-owned policies this RDR conforms to rather than authors: predicate operator
grammar and tag type model (RDR 0003), accessor execution and read-back safety
(RDR 0004), CLI surface and error mapping (RDR 0005), cross-row graph lint
(RDR 0006), guard seam / atom shape / existence constants (RDR 0007), reserved
`recognized` key name (RDR 0008), escape-row shape conformance (RDR 0009).

Concurrency model, secret/credential lifecycle, memory management, licensing,
deployment model, and IDE compatibility do not apply — omitted per the template
rather than N/A-bulleted.

## 5. Proportionality

**Right-sized. Nothing flagged for trimming; lock as written.**

**Contract count — one seam, not several.** This RDR is the sole author of
exactly one independent load-bearing contract: *the sparse transition-model
format and its normalization to candidate rows.* The wire schema, the version
gate, the category taxonomy, the identity tuple, and the expanded-table dump
order are all facets of that single seam — none is separable, since each is
defined in terms of the normalized candidate row. Every adjacent seam is
explicitly delegated: operator grammar → 0003, accessor execution → 0004, CLI
mapping → 0005, cross-row lint → 0006, atom shape and existence constants →
0007, reserved key name → 0008, escape-row shape → 0009. The RDR does not lock
seams it does not own, and the load/lint **arity split** (single-rule → load
here, cross-row → lint at 0006, evaluation-time → kernel refusal at 0001) is the
principle that keeps the boundary decidable rather than negotiated case by case.
No split warranted.

**Profile re-validated against the contracts just counted.** `foundational` is
correct and not inflated: this RDR is a cross-RDR **producer** — RDRs 0003,
0006, 0007, 0008, and 0009 all consume the wire format, the normalization
semantics, the dump's total ordering, or the validation category taxonomy, and
0008 and 0009 both name this RDR's normalizer as the enforcement point their own
checks land inside. That is the matrix's cross-RDR-producer trigger on the
contract axis, independent of the accretion axis (`Seam Lineage`: no prior
accretion, so no floor applies). The Profile field carries the value plus one
clause naming the contracts, with no matrix or provenance prose left from the
template — correct form.

**Lenses match the Profile.** `foundational` owes cove → 3amigo → critique →
repeatability, and all four ran with complete evidence and dispositions
(critique dual-model; repeatability ×3 across `claude-opus-5[1m]`,
`Claude Sonnet 5`, and `claude-fable-5`, all `variant: full`). No lens was
skipped against a wrong Profile, so the latch's backstop finds nothing to
correct. No `Transient`-marked contract exists to discount.

**Length is load-bearing, not bloat.** The document is long, but the mass sits
where a `foundational` producer must be precise — 23 normative blocks, a
complete category taxonomy five consumers depend on, and per-assertion failing
controls. Three specific stretches earn their length by recording a *rejected*
reading that a reader would otherwise reconstruct wrongly: the "combined
predicate set" exclusion in outcome binding (R-6, the wider term MUST NOT be
read into the narrow lifting rule), the row-kind clause (R-9, a render-time view
MAY materialize the derived column but MUST NOT feed it back), and the fail-fast
ordering clause (R-3, order among independent defects is *deliberately*
unspecified, so no later implementation is bound to whichever order the first
one happened to use). Each was a repeatability-diff finding where independent
model runs diverged — removing the prose would restore the ambiguity that
produced the divergence.

The one advisory note: A8's Evidence field runs 32 lines against a 30-line
budget (CHECK 9, `1 field over budget, 32 lines`). Reviewed per assumption as
the check asks — the mass is the residual-scope analysis, the load-bearing
anchors (`::Table.models`, `::escapeOrRefuse`, `::assemble`,
`recognizedTagKey`) remain findable within it, and relocating the balance to
`{ARTIFACT_DIR}` would separate a bounded guarantee from the reasoning that
bounds it. Keep as written. Advisory, non-blocking, no action.

---

## Verdict

**READY.** No open blockers. The mechanical sweep passes after one in-pass
conformance fix, all four `foundational` lenses ran with dispositions, the six
`Pending` assumptions are deliberate Stage 6 downgrades each carrying a named
MVV oracle and a survivable `If wrong`, the MVV is in scope with its one
cross-row assertion correctly handed to RDR 0006, and the RDR owns exactly one
contract at the `foundational` profile its consumer set requires.

Status → **Final**.
