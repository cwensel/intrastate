Model: claude-opus-5[1m] (diff pass; reconciles claude-opus-5 × claude-fable-5)

# Critique dual-model diff — RDR 0006, iter-2

Reconciled **by passage anchor**, not by ID (`C-N` is per-file). Pass A =
`critique.md` (claude-opus-5, 18 rows). Pass B = `critique-modelB.md`
(claude-fable-5, 12 rows).

## Convergence summary

Both passes independently landed on the same three structural defects, and both
name the same locked-peer contradiction. That is agreement between isolated
contexts on different base models — a hotspot, not a shared blind spot. Twelve
of B's twelve rows have an A counterpart; A carries six rows B did not reach.

**No row in either pass contradicts a row in the other.** The disagreement is
entirely about *scope of consequence* (see Divergences), never about whether a
defect exists.

## Anchor-matched rows

| Passage anchor | A | B | Agreement |
|---|---|---|---|
| §TD "Source state is a tag-set, not a name" / group identity | C-1 | C-1, C-10 | **Both**: per-node satisfiability membership *is* a second grouping predicate, contradicting 0006's own MUST-NOT clause and locked 0003. Both flag A2's "closes 0003::A10" as false. |
| §TD "Lint MUST NOT key groups on the syntactic match pattern" | C-11 | C-10 | **Both**: direct contradiction with locked 0003 "Match keys are not product dimensions … A row's match pattern selects which group the row belongs to". Both note Contradiction Check claims none. |
| §LBD "Join rule and termination" / node count | C-2 | C-2, C-11 | **Both**: product bounded per group, node count unbounded. B adds the merge trigger is circular (value-set-tuple equality ⇒ union is a no-op) and `held ⊔ absent` undefined. |
| §TD envelope `{"type":"failed","error":{…}}` / `error.findings` | C-4 | C-3 | **Both**: shape exists only in respond.go's doc comment; `EmitJSON` marshals `*CLIError` flat. Both note A4 is stamped Verified against a misread symbol. B additionally cites `docs/cli-output-contract.md` saying so explicitly. |
| A5 gate target "a new `lint` job" | C-7 | C-5(a) | **Both**: `jobs.lint` already exists (golangci-lint). |
| A5 "over the checked-in transition model" | C-6 | C-5(b) | **Both**: no checked-in model exists; A5 excludes the fixture corpus by name; nobody books authoring one. |
| Invariant 2 "every value … meets it" + merge | C-8 | C-12 | **Both**: terminal satisfaction is anti-monotone in merging; converging legal flows take structural `graph-dead-end` with no legal cure. |
| §Trade-offs "no suppression or waiver mechanism" / "rare" | C-15 | C-9 | **Both**: rarity asserted, never measured; structural given the other findings. B grounds it in 0003's own desk-trace fixture (8,192 assignments at declared defaults). |
| A7 route-back on locked 0003 A18 | (in C-17 scope) | C-7 | **Both**: verification condition is an amendment to a `Final [locked]` peer; nobody booked. B adds the two Finals disagree *in writing* today. |
| Normative: text mode MUST enumerate / MUST NOT write directly / `data.findings` | C-16 | C-4 | **Related, complementary**: A anchors on `Data`'s own `omitempty` defeating "empty list is the receipt"; B anchors on `respond.OK`'s text branch dropping `Data` and `EmitText` rendering no findings, with the renderer changes unbooked ("Reuse"). Both are real; they are different halves of the same output-surface gap. |
| A6 producer / unlanded 0002 schema clause | (noted in A's §2/§3) | C-6 | **B-stronger**: B rows it — the whole reachability apparatus hangs on an unmade edit to a Draft peer; until it lands every model is blocking-rejected, so lint is unshippable independently. A discusses this in prose but did not row it separately. |
| §Performance "No throughput target is load-bearing" | (folded into C-2) | C-11 | **B-stronger as its own row**: the one dimension declared not load-bearing is the one the design leaves unbounded. |

## A-only rows (B did not reach)

| A row | Passage | Why it stands |
|---|---|---|
| C-3 | §LBD "an edge is a normalized non-escape row" | Escape rescue is a legal runtime transition with no edge in lint's graph. 0002 forbids escape rules carrying write/clear, so the RDR's own prescribed cure is a load failure. Hard deadlock — needs grounding against 0002's escape clause. |
| C-5 | §TD "Type ownership … subsystem-agnostic record" vs "MUST carry `Key`, `Operator`, `Literal`, `Block`" | `Block` and failure class are 0002/0001 vocabulary; either `clierr` grows the dependency the clause exists to prevent, or the fields degrade to strings and "MUST carry" is unenforceable. |
| C-9 | Invariant 5 "A write … **replaces** its prior value" | Write semantics no cited peer owns, used to *delete* the path-accumulation check. |
| C-13 | §Approach "a locally-edited illegal model can still be resolved against … That is accepted" | "MUST NOT be accepted for resolver use" has no enforcement site. |
| C-14 | `graph-unprovable-coverage` "widened trigger" | One code, three triggers with three different remedies. |
| C-17 | A2 "A17 … A19 … A20's lint half … close by citation" | A Draft peer cannot close a locked peer's records. A7 concedes this for A18; A2 does not for the other three. |
| C-18 | Invariant 6 read-set vs `RequiresOwned` | 0007 (Final) narrows `RequiresOwned` to write dependencies; the owned *read*-set has no producer and is absent from the minimum input contract. Risk is a silent false green. |

## B-only rows (A did not reach)

| B row | Passage | Why it stands |
|---|---|---|
| C-8 | Advisory tier "closed at" three codes vs 0003's vacuous-`exists` report | Locked 0003 requires an `exists`-over-always-present atom be reported "well-formed but vacuous"; the closed three-member tier has no code for it, and Scenario 11's control asserts a bare pass. Testable violation of a locked clause. |
| C-5(c) | §Illustrative Code "`--flow <id>` … through RDR 0005's config discovery" | 0005's own audit defers config discovery ("Parser placeholder; no transition-model config yet"), so `--flow` binds to a facility that does not exist. Third leg of B's Scenario-7-uninstantiable finding. |

## Divergences (the dual-model signal)

Substantive disagreement is narrow. Both models found the same defects; they
differ on **which direction the over-approximation fails**:

1. **`[direction split]` — C-1's soundness consequence.** A frames the grouping
   defect as producing **false positives** (rows the runtime never evaluates
   together reported as `graph-overlap`). B frames it as producing a **false
   green** for invariant 4: the coverage union at a merged node is a *superset*
   of the union at each concrete state it abstracts, so lint can certify
   coverage a real reachable state does not have — and B notes the RDR's
   no-false-green derivation covers invariants 2 and 6 only, flipping direction
   for 4. **B's reading is the more dangerous one and is not refuted by A.**
   This is the highest-value disagreement in the diff: the draft's central
   safety claim ("can only produce a false positive, never a false green") may
   be scoped too broadly. Resolve first.

2. **`[merge-rule split]` — is the widening implementable at all?** A assumes
   the fixpoint works and counts its nodes (exponential). B argues the merge
   trigger is *circular* — with no named states, "the same successor" can only
   mean value-set-tuple equality, under which the union is a no-op and the
   widening never fires. If B is right, A's node-count finding is a symptom and
   B's is the cause. Not contradictory; nested.

3. **Scope of the output-surface gap.** A (C-16) and B (C-4) each found a
   different half. Neither is wrong; the union is the finding.

## Gate read

**Healthy, both passes.** Concrete named passages; anchors resolve; origins span
§1/§2/§3/premortem/AT rather than clustering on §1. Ledger rows index defects
the prose actually argues. No generic advice. Neither pass is a re-run
candidate.

Isolated-context convergence on C-1/C-11, the envelope, the CI job, the missing
model, and invariant 2 is agreement, not contamination — the passes ran
concurrently and neither could read the other.

## Resolve scope

The origin ledger for the resolve half is the **union**, reconciled by anchor:
18 A-rows + B's 2 unique rows + the 3 divergences above = **23 distinct
defects**. Rows matched across passes resolve once, against the stronger of the
two framings (noted per row).
