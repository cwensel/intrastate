Model: claude-opus-5[1m]

# Cove (S-04) — RDR 0007 Guard Predicate Totality

Chain-of-Verification with embedded Step 0 grounding sweep. Iteration 1.
Two independent sub-agents: Step 0 (codebase claim sweep, 24 claim clusters)
and Steps 1–2 (12 verification questions answered independently).

## Step 0 — Grounding sweep result

Every load-bearing codebase claim CONFIRMED against source except as noted in
F-1/F-2/F-3 below. The three claims A3 rests on are exact:

- `resolve.go::gate` L388 calls `missingOwned(survivors, view)`, not `rows` —
  GuardFalse prunes (L381 `continue`) before the owned sweep, which precedes
  the undecidable loop (L396–417).
- `resolve.go::gate` L412 `return nil, &Refusal{Kind: KindGuardUnevaluable, …}`
  discards the named return `selected` — no decided-true sibling escapes past
  an unevaluable candidate.
- `resolve.go::escapeOrRefuse` L485–488 `if blocked != nil { return refuse(in,
  *blocked) }` — the escape set's blocking refusal replaces the candidate-set
  refusal `r`. The RDR's *corrected* aggregation clause matches shipped
  behavior; the pre-correction clause was rightly refuted at Resolve.

The three negative/gap claims are all correct: no `testdata/` exists repo-wide;
`TestAdv3` pairs two *unknown* predicates with `fixtureGuards{}` deciding
neither (order-stability only); no shipped test puts a missing-owned row beside
an unevaluable-guard row in one survivor set.

Peer-RDR claims for 0001 (D5/D8, selection rule, Implemented), 0002 (L292
escape closure, tag-declaration obligation), 0003 (closed vocabulary, `exists`
alone absence-inspecting, silence on absent operands, A3 no-inline-not,
Alternative 3 rejected, `unless` placement-only, row verdict shape,
exhaustiveness permission ceiling, A4 semantic kind, `unknown tag` outside a
normative block) all CONFIRMED verbatim.

## Step 0 — Inverse check (new rule vs. existing sibling)

**Code: searched, none exists.** `GuardResult` is produced only in
`resolve.go` and test fixtures. No existing code decides absent-key⇒unevaluable,
implements three-valued combination, or derives undecidability from an absent
key. `gate` aggregates undecidability *already handed to it*.

**Spec: a sibling DOES exist — see F-1.**

## Findings

### F-1 [Step 0 / REFUTED + sibling-not-cited] — RDR 0004's Normative Contracts already decide the no-collapse question; A6b mischaracterizes the block as silent

A6b Evidence asserts "RDR 0004 Normative Contracts states **only** 'A read
accessor MUST return typed tag values or a typed refusal…'". The block contains
eleven clauses, one directly on point:

> ```normative
> A gate accessor MUST return allow, deny, or indeterminate. Indeterminate MUST
> be a refusal-class result, not a false allow and not a false deny.
> ```
> — `docs/rdr/0004-accessor-execution-safety-model.md` L263–264

The RDR's elided failure taxonomy (`"artifact unavailable, timeout, execution
failure…"`) drops exactly the third-value terms: the full line L527 reads
"…execution failure, **gate denied, gate indeterminate**, write attempted for a
non-owned tag, and read-back mismatch." Carried forward normatively by RDR 0005
("MUST NOT … coerce gate allow, deny, or indeterminate results into tag
values") and shipped as CLI code `flow-gate-indeterminate`.

A6b's *substantive* claim survives — the **read-accessor** clause genuinely
carries no completeness requirement, so read-failed vs. genuinely-absent remains
unowned. But the evidence sentence asserts a negative about the whole block that
the block refutes, and the omitted clause is the strongest in-repo precedent for
0007's own rule: a peer Final RDR already ruled a third value MUST NOT collapse
into either boolean, phrased almost identically to 0007's "never to false and
never to true." The RDR sources that principle only from SCXML and PostgreSQL
while an adjacent seam in this repo already decided it.

Distinction that keeps 0007's contract intact: 0004's gate accessor *reports*
indeterminate as an execution outcome; 0007 *derives* unevaluability from an
absent key in the view. 0004 says nothing about absent-key⇒unevaluable, so it
does not pre-empt this contract.

### F-2 [Step 0 / REFUTED] — the `RequiresOwned` conflation lives in a Go doc comment, not in RDR 0001

Problem Statement: "`Row.RequiresOwned` is documented as naming 'owned tag keys
the row's evaluation needs'". Overrides: this RDR "*narrows* RDR **0001's
documented** `Row.RequiresOwned` meaning".

The quoted string appears exactly once repo-wide, in code:
`internal/resolve/resolve.go` L180 `// RequiresOwned names owned tag keys the
row's evaluation needs.` Grep of `docs/rdr/0001-resolution-kernel.md` and
`0001-resolution-kernel/artifacts/*.md` returns zero hits. RDR 0001 never
defines the field's meaning.

Favorable to the RDR — the Phase 1 edit is unambiguously in-scope and no peer
RDR needs reopening — but Overrides currently claims to narrow peer text that
does not exist. Attribution belongs on the doc comment.

### F-3 [Step 0 / completeness gap] — A6a's reader enumeration is incomplete

A6a names `TagSet.Lookup` / `TagSet.has` as the readers. `resolve.go::TagSet`
has two more: `matches` (L132–140, the only reader on the candidate-selection
path) and `Len` (L122). Not a refutation — `matches` is a match-pattern test
that cannot remove a key, so presence-monotonicity holds — but a reviewer
re-running the Source Search finds unnamed readers.

### F-4 [Step 2 Q8 / codebase-refuted claim] — `fixtureGuards` cannot exercise the domain rule; the Infrastructure Audit overclaims and Scenarios 4/6/7/8 have no evaluator

Existing Infrastructure Audit: "`fixtures_test.go::fixtureGuards` … Table-driven;
answers unknown predicates `GuardUnevaluable` … Already domain-rule-conforming
— MVV rows 1–3 build on it rather than adding a stub."

Source: `fixtures_test.go::fixtureGuards.Evaluate` L20 —
`func (g fixtureGuards) Evaluate(guard string, _ resolve.TagSet) resolve.GuardResult`.
The `TagSet` is **discarded**. Verdicts come solely from `decided map[string]bool`
keyed on guard *text*. It is structurally incapable of exercising tag
presence/absence, operator semantics, `exists` totality, `contains`-over-absent-set,
or any strong-Kleene combination.

It can back MVV rows 1–3 (which only need a verdict handed in). It **cannot**
back Testing Strategy Scenarios 4 (strong-Kleene matrix), 6 (`exists` totality),
7 (`contains` over absent set-valued tag), or 8 (two-row absence pattern) — all
four require a view-reading evaluator that does not exist and is budgeted in no
phase. Highest-value finding: it is the RDR's own conformance mechanism.

### F-5 [Step 2 Q10 + Q1 / silence] — the RDR is silent on the empty-atom-block identity and the unguarded row

(a) **Empty `unless` block.** Under `all ∧ ¬(unless_conj)` the conventional
empty-conjunction identity `unless_conj = T` would disable **every** row.
Neither 0007 nor 0003 states the empty-block identity. Genuine footgun.

(b) **Unguarded row.** `resolve.go::evaluateGuard` L425–433 returns `GuardTrue`
for `guard == ""` *before* the nil-seam check — the seam is never called, so the
domain rule has no purchase. The kernel documents this (`resolve.go` L183–184:
"Guard is the predicate handed to the guard seam. Empty means no guard.") but
the RDR's Normative Contracts never state that unguarded rows sit outside the
partial-function domain.

### F-6 [Step 2 Q11 / silence] — Phase 2's conformance obligation has no owner and no enforcement mechanism

Risks: "the conformance fixture set (Phase 2) is normative and lives in the
kernel's test tree; 0003's implementation must run it, and the Prerequisites
bind 0003's implement stage to this RDR being Final." The Prerequisites checkbox
binds **sequencing** (0007 Final precedes 0003 implement), not **execution**.
The vectors live in `internal/resolve/`; 0003's evaluator is a separate future
package, so `go test ./internal/resolve/` passes green against the stub
regardless of what 0003 builds — precisely the P-10 failure ("a diverging
evaluator otherwise ships green") the phase exists to prevent. The mitigation
has the same hole as the risk.

### F-7 [Step 2 Q9 / evidence overstatement] — A2 presents a lint-domain formulation as RDR 0003 "fixing" the runtime row verdict

RDR 0003's formulation is set-theoretic over *complete* assignments in a finite
declared domain (Technical Design L266–270: "a row's accepted assignments are
the intersection of all positive `all` atom domains **minus** the single
conjunctive assignment set matched by the row's full `unless` block"). Set
subtraction over total assignments equals `all ∧ ¬(unless_conj)` only in a
two-valued world. The lift to three-valued runtime K3 is **0007's own
extension**, which is legitimate and is exactly what A1 says this RDR is for —
but A2's Evidence presents it as RDR 0003 already fixing the row verdict.

### F-8 [Step 2 Q6 / non-finding, recorded] — `recognized`-key provenance flip

An owned/observed tag keyed `"recognized"` overwrites the recognized entry's
value *and* provenance (`resolve.go::assemble` L148–164; `recognizedTagKey` is
the bare string `"recognized"`, L110, with no reservation). Does not refute A6a
(presence is still monotone; this is a value/provenance change). Key-binding is
RDR 0008's declared territory.

### F-9 [Step 2 Q7 / non-finding, recorded] — no shipped test changes verdict under the narrowing

Fixtures list `RequiresOwned: ["status"]` on rows that write `status` — legal
under the narrowed reading. `adversarial_test.go`'s `"never-present"` poison key
models nothing real post-narrowing, but it is an adversarial probe, not a
contract example. "No kernel behavior change" survives.

## Verdict

Healthy pass. The RDR's central claims hold up under source. Both self-declared
test gaps are real and correctly characterized; the escape-half self-correction
is consistent across the normative block, A3, the spike, and Scenario 5; D8
prune-first ordering and its test sensitivity verify exactly as described.

Findings cluster where the RDR meets surfaces it does not own: the conformance
mechanism (F-4), the peer-RDR precedent it overlooked (F-1), and three silences
(F-5, F-6).
