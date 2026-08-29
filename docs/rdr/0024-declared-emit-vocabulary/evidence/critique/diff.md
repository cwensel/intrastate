Model: claude-opus-5

# Critique diff — cli/0024, Pass A (opus-5, 18 rows) vs Pass B (sonnet-5, 9 rows)

Diffed by passage anchor, not by ID; the `C-N` numbering is per-file and does
not correspond. Grounded against `intrastate` at `4c545df`, with a live decoder
probe against `pelletier/go-toml/v2 v2.2.4` and direct reads of the consumer
models at the sibling `rdr` working directory.

The headline: the two passes barely overlap. Four merged entries are AGREED;
eleven are A-ONLY; three are B-ONLY. That asymmetry is itself the finding — A
ran the code and B ran the document. Every one of A's six code-grounded rows
verifies CONFIRMED. B's rows are almost entirely reasoning about the record's
scope commitments, and three of its nine land squarely on already-charted
successor scope.

## Merged ledger

| Merged-ID | Passage anchor | A-row | B-row | Class | Grounding | Defect |
|---|---|---|---|---|---|---|
| M-01 | `0024:C1` refusal arms — "a disposition sub-table carrying no keys" vs "an `enum` with no domain"; `0024:S1` fixture-per-arm | C-1 | — | A-ONLY | **CONFIRMED** | Two enumerated `malformed_emit_declaration` arms decode to the identical state (`Domain == nil`); S1 demands a fixture per arm, so one arm is unimplementable as written. |
| M-02 | `0024:C1` "Strictness still catches a typo at the DECLARATION level (`domaim = [...]` refuses)" | C-2 | — | A-ONLY | **CONFIRMED** | The declaration-level typo refuses under `unknown_schema_field`, a category this RDR does not register; C1 presents it inside the emit-grammar paragraph without naming it. |
| M-03 | `0024:F1` "offending file in `locator`"; `0024:MVV` step 2 | C-3 | — | A-ONLY | **CONFIRMED** | `atLine` is wired at exactly one loader site; no emit analogue is specified or scheduled, so every emit refusal renders `file:1`. |
| M-04 | `0024:A8`, `0024:C3` "the property `TestReq146` already pins", `§phase-2` failure signal | C-4 | C-5 | **AGREED** | **CONFIRMED — A8 is FALSE, not merely Pending** | `TestReq146` asserts a hand-maintained list over `table.Row`; `Model.EmitDecls` is model-level and structurally unreachable. Phase 2's stated failure signal cannot fire. |
| M-05 | `0024:S6` category witnesses, "per the convention `dump_test.go`'s REQ-119 witness map" | C-5 | — | A-ONLY | **CONFIRMED** | The convention is REQ-119 *and* REQ-131 — trip the category "and no other". S6 asserts registration only, never exclusivity, under a loader whose order C2 leaves unspecified. |
| M-06 | `0024:C4`, `§phase-3` (ii) "mechanical regeneration is authorized HERE by name"; `0024:A2` | C-6 | C-9 | **CONFLICT** | **CONFIRMED for A; B's framing REFUTED** | A: the 28 goldens are 0011's REQ-118 pre-change capture, and regenerating them voids the guarantee. B: they are ordinary brittle fixtures and the real defect is a compounding fixture tax. The file's own doc settles it for A. |
| M-07 | `0024:C2` "The six steps between `loadTags` and `normalizeRules`" | C-7 | — | A-ONLY | **CONFIRMED** | C2 says six and names five; `load.go:85-89` has five. A normative placement contract that miscounts its own insertion point. |
| M-08 | `0024:A4` "`conformDomain`'s enum arm is exactly the membership test C2 needs"; `0024:C3` throwaway `TagDecl` | C-8 | — | A-ONLY | **CONFIRMED** | The enum arm is guarded by `len(decl.Domain) > 0`; `conformKind` has no `enum`, no `scalar` and no default arm. The reuse is sound only under a C1-fires-first precondition neither A4 nor C3 states. |
| M-09 | `0024:§approach`, `§consequences` "ships capability, not coverage", `0024:A5`, `0024:A3` | C-9 | C-4 | **AGREED → RE-RAISE** | **RE-RAISE (Charted PM-3)** | Merge-day coverage is zero and the only adopter's models live in a sibling repo, unscheduled. Already charted to a consumer-repo successor. |
| M-10 | `0024:§consequences` N-round loop; `0024:C2` fail-fast, order unspecified | C-10 | C-2 | **AGREED → part RE-RAISE** | **RE-RAISE (Charted: multi-defect load reporting) + one A-ONLY residue CONFIRMED** | Both name the adoption grind. A adds a diagnosis B does not reach: this is the first load category that fires N times on a model the author did not change. |
| M-11 | `0024:C1` "`scalar` is the declared-but-unvalidated escape hatch"; `§risks-and-mitigations` scalar hollowing; Meyer condition | C-11 | (C-2 tail) | A-ONLY | **CONFIRMED as fact; defense is partly fair** | `0006:C17` genuinely closes the advisory tier at four codes, so the procedural ground is real — but no lint output distinguishes a fully-`scalar` model from a proven one. |
| M-12 | `0024:C4` "`omitempty` is available on this field and is deliberately not taken" | C-12 | — | A-ONLY | **CONFIRMED** | The stated justification is self-defeating: a pre-field binary emits no `dispositions` key at all (JSON `undefined`), which is exactly the absence case `{}` was meant to disambiguate from. |
| M-13 | `0024:C2` "Proving the AUTHORED form is proving the executed form"; `§risks-and-mitigations` domain drift | C-13 | C-8 | **AGREED → RE-RAISE** | **RE-RAISE (Charted PM-7)** | The model is proved against itself; the seed defect is relocated, not closed, and handed to an unbuilt consumer seam test. Already charted. |
| M-14 | `0024:C1` non-canonical `int` literals (`03`, `+5`, `-0`) | C-14 | — | A-ONLY | **CONFIRMED** | No shipped or consumer model emits an int-shaped value; A5's census finds enums and one `scalar`. An unused kind carries a documented wart into a locked grammar. |
| M-15 | `0024:§phase-1` "Amend three stale comments"; `0010:C3` "undeclared" | C-15 | — | A-ONLY | **PARTIALLY CONFIRMED** | The RDR does address `TestReq27` and does define "undeclared" as "not a tag". But the test's own failure-message prose (`emit_0010_test.go:454,477`) carries the word and is not on the amendment list. |
| M-16 | `0024:S3` "The pass criterion is LOAD-side equivalence, not payload byte-identity" | C-16 | — | A-ONLY | **CONFIRMED** | S3's criterion is measured *after* Phase 3 regenerates the only whole-payload oracle that could falsify it. The criterion is unfalsifiable at the moment it is measured. |
| M-17 | `0024:S7`, `§phase-3` "whichever record lands second"; `0024:C4` JDR §D1 | C-17 | C-6 | **AGREED** | **CONFIRMED** | S7 runs against an oracle that does not exist at HEAD. Diagnoses match; A adds that S7's own transfer clause contradicts its "unsatisfiable while unmet" claim. |
| M-18 | `0024:§decision-rationale` QOC matrix; `0024:ALT1` | C-18 | — | A-ONLY | **CONFIRMED as stated** | The correctness row scores 5 on catching both seed defects, but the value-defect catch is conditional on the domain being right (M-13). 22-vs-18 is a four-point spread on a hand-scored table. |
| M-19 | `0024:A5`, `0024:C1` closed enum domain vs interpolated values; `0010:A6` deferred widening | — | C-1 | B-ONLY | **REFUTED at HEAD; real as forward risk** | B claims `rdr-status.toml` is "one interpolated value away". All 54 blocks / 22 distinct `next` values are fixed literals with zero interpolation. `0010:A6` (Verified) fixes the record number as caller-supplied. |
| M-20 | `0024:A7` (Pending), C1's hand-written arms closure | — | C-3 | B-ONLY | **CONFIRMED as restatement, superseded by M-01** | B re-states A7's own open closure question. A's M-01 does the harder work: it names a specific arm that is not merely unproven but unimplementable. |
| M-21 | `0024:D-naming` `dispositions` plural vs per-key single token | — | C-7 | B-ONLY | **CONFIRMED as fact; low severity** | The rejected-alternatives list compares only sibling names and never interrogates the chosen name's shape-signaling. C4 does document the shape, so this is a docs-clarity nit, not a spec gap. |

---

## Adjudications

### M-06 — CONFLICT: what the 28 goldens are

The passes agree the 28-site diff is real and disagree on what it costs.

A (C-6) reads them as `0011`'s REQ-118 discharge — a frozen pre-change
reference whose regeneration destroys the guarantee. B (C-9) reads them as
"deliberately brittle" fixtures and pivots to a different complaint: that the
test architecture will levy this tax again when 0023 lands.

The file settles it in A's favour. `flow_demand_0011_test.go:472-477` heads the
test with `REQ-118 / 0011:S8` and `A-11: discharged against a pre-change build
of the same tree — the oracle may be a golden payload "provided it is captured
from the pre-change behaviour and asserts the FULL payload rather than a
subset — which is the gap S8 names"`. Line 486 repeats it inline
("captured from the PRE-change behaviour and asserts the FULL payload"), and
the sweep's failure text at `:563-566` asserts "every one of these payloads is
byte-identical before and after". Line 623 states it a third time.

These are not fixtures with a brittle style. They are the artifact a peer REQ is
measured against, and `0024:A2` applies its own peer-spec test to `TestReq39`
while exempting these on a distinction the file itself does not make.
B's compounding-tax observation is true and orthogonal; it is not a substitute
for the finding, and taken alone it would license exactly the regeneration
that voids REQ-118. **Resolve on A's diagnosis; B's framing is the safer-sounding
and wronger of the two.**

### M-19 — B-ONLY, REFUTED at HEAD

B's C-1 is its lead finding and its §3 "assumption that will not survive first
contact". It asserts `rdr-status.toml` is "one interpolated value away" from
breaking A5, and its premortem stages a route "that needs a record id folded
into the command text".

Enumerated at HEAD, `/Users/cwensel/sandbox/newcoinc/rdr/models/rdr-status.toml`
carries 54 `[rule.emit]` blocks and 22 distinct `next` values. Every one is a
fixed literal: 14 `/rdr-*` routes (including multi-token ones like
`/rdr-prelock repeatability diff`, which are fixed strings, not templates), 6
`stopped:` tokens, plus `none` and `resolve:lens`. A5's census reproduces
exactly. There is no interpolated value and no partial one.

`0010:A6` (Verified) already adjudicated the composition B raises: the record
number "is the argument the caller already holds — the user typed it — so it is
never a table value and the table never interpolates it". B cites A6 as
supporting its case; A6 is the passage that forecloses it.

The residual risk — that a future author writes a templated value into an
otherwise-clean key and the grammar's only escape is whole-key `scalar` — is
real and is the one thing B contributes here that A does not. It belongs in
Risks as a named forward condition, **not** as a refutation of A5. Marked
REFUTED as a defect at HEAD, retained as a forward-looking note.

### M-15 — PARTIALLY CONFIRMED, narrower than A states

A's C-15 asserts the RDR "reinterprets `0010:C3`" and omits `TestReq27` from
the amendment set. Half of that is answered in the draft: `§phase-1` names
`TestReq27_AnEmitKeyIsNotATagKey` explicitly, asserts it "keeps passing
unchanged", and defends the reading — "'undeclared' there means **not a tag**
— no provenance, no accessor, no match/guard/write participation — every clause
of which C1 keeps". `[emit]` is a separate namespace from `[tags]`, so the
behavioral claim is correct.

What survives is the prose leg. `emit_0010_test.go:454` prints "an emit key is
undeclared and is not a tag" and `:477` prints "emit keys are undeclared and
uninterpreted". Both restate the word the three amended comments are being
amended *for*, and neither is on the amendment list. A's conclusion (the record
reinterprets a locked contract without an override) overstates; A's observation
(the amendment set is incomplete by its own criterion) holds.

### M-10 — AGREED, but only one pass has the diagnosis

Both passes name the N-round adoption loop, and both quote the RDR's own
`§consequences`. The Charted file already routes multi-defect load reporting to
a successor, so the *remedy* is dismissed-with-cite.

A's residue is not covered by that charting. A observes that every existing load
category in `internal/table/category.go` fires on a single authoring mistake the
author just made, whereas `unknown_emit_key` and `emit_value_out_of_domain` are
the first to fire N times on a model the author did not change — opt-in
strictness applied retroactively over an existing corpus. That is a claim about
whether the tier's cardinality contract *transfers* to this category, not a
request for multi-defect reporting. It compounds with M-03 (`file:1` on every
one of the N runs), and the record cross-references neither section to the
other. Worth a sentence in Consequences; not a re-open.

---

## Where each pass was blind

- **B had the evidence and missed it.** M-04 is the one code claim B checked,
  and it got it right. But B's own §5 scenario for C-3 asks the Finalization
  Gate to verify A7's closure — the exact check A ran with a decoder probe, and
  which returns a specific unimplementable arm (M-01). B stopped at "this is
  Pending" where one probe would have produced the defect. Same for M-08: B
  never opened `load.go`, so the `len(decl.Domain) > 0` guard is unseen.
- **A never looked where B looked.** A does not engage the interpolation
  headroom question (M-19) at all, and does not raise the `dispositions`
  plural/shape naming (M-21). Both are document-level readings A's
  code-first method skipped. Neither is load-bearing, but M-19's forward-risk
  residue is a genuine addition.
- **Both missed** any check of whether `dispositions` reaches text mode's
  renderer correctly, and neither tested the `[emit]` grammar against
  `rdr-write.toml`'s `op` key, which is a second, differently-shaped real
  corpus (10 distinct values, and the `edit` key's multi-line shell script that
  A5 correctly flags as forcing `scalar`).

---

## Resolve these first

Ranked by whether the draft is wrong as written, then by blast radius.

1. **M-01 — C1's unimplementable refusal arm.** The only finding that makes a
   normative contract undischargeable. `[emit.next.domain]` with no keys and an
   absent `domain` both decode to `nil` under the repo's decoder; C1 lists them
   as two arms and S1 demands a fixture for each. Merge the arms explicitly, or
   license the source-text scan that separates them. Probe reproduced at HEAD.

2. **M-04 — A8 is false, not Pending.** Both passes agree; the code confirms.
   `TestReq146`'s sweep iterates a literal list of `table.Row` field names and
   its determinism leg compares `got.Rows`. `Model.EmitDecls` is unreachable
   from both by construction. Flip A8 to Refuted, drop the corroboration from
   C3's normative text, and let scenario 4 carry the determinism proof alone —
   which C3 already half-concedes in a parenthesis.

3. **M-08 — the `ConformValue` reuse precondition.** `conformDomain`'s enum arm
   is guarded by `len(decl.Domain) > 0` and `conformKind` has no `enum`,
   `scalar` or default arm. A4 quotes the enum arm as "exactly the membership
   test C2 needs" and omits the guard. State in C3 that the reuse *depends* on
   C1's empty-domain and unknown-kind arms having already refused. Without that
   sentence, an implementer who relaxes a C1 arm (see M-01) ships a loader that
   admits everything under a declaration that reads strict.

4. **M-06 — the golden regeneration voids REQ-118.** Phase 3 authorizes
   destroying a peer RDR's discharge artifact while asserting the guarantee is
   untouched. Either preserve a genuine pre-change capture, or demote REQ-118
   explicitly and say so. Resolving this also resolves **M-16**, whose oracle is
   the same artifact.

5. **M-07 — C2 counts six steps and names five.** Trivial to fix, and it is a
   placement contract. Five: `loadAccessors`, `loadDump`, `loadContexts`,
   `loadInitial`, `loadTerminal` (`load.go:85-89`).

6. **M-03 — every emit refusal renders `file:1`.** `atLine` has exactly one
   call site. F1 and the MVV both promise a locator; neither says where it comes
   from. Either name an `emitHeaderLine` analogue in Phase 1, or strike the
   locator promise and say refusals are file-scoped. Compounds with M-10.

7. **M-02 — the declaration-level typo is someone else's category.** One
   clause: name `unknown_schema_field` where C1 claims strictness catches
   `domaim`, so S1's fixture asserts the right category.

8. **M-05 — the witness "and no other" half.** REQ-119/REQ-131 require each
   witness to trip its category and no other; S6 asserts registration only. Add
   the exclusivity leg or state why it does not apply under C2's unspecified
   order.

9. **M-12 — `omitempty`'s stated justification is self-defeating.** A pre-field
   binary emits no key at all, so `{}` cannot be the signal C4 says it is. Keep
   the decision if it is wanted for wire-shape stability, but give it a reason
   that holds.

10. **M-17 — S7 is discharged by a paragraph.** Both passes agree. Either give
    Phase 3 a local assertion that `dispositions` is an accounted-for payload
    field, or stop calling S7 a scenario.

**Dismiss with cite** (already charted, `evidence/3amigo/Charted.md`): M-09
(PM-3, adoption scheduling), M-13 (PM-7, consumer seam test), M-10's remedy
(multi-defect load reporting). Fold M-10's *residue* — that this is the first
load category to fire N times on an unchanged model — into Consequences as one
sentence rather than reopening the tier decision.

**Note, do not resolve**: M-19 (interpolation headroom) as a forward condition
in Risks, refuted as a present defect; M-21 (`dispositions` plural) as a docs
clarification if the output contract is touched; M-14 and M-18 as accepted
costs the record can state plainly.
