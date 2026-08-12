Model: claude-opus-5

# Hostile Set Critique — RDR 0002–0009

Assignment: hostile critique of the Final-and-unimplemented cluster (0002, 0003,
0004, 0005, 0006, 0007, 0008, 0009). RDR 0001 is Implemented and is context, not
a member. No hedging, no balance.

**Verdict up front.** This set will not ship. Not because any single RDR is
wrong — most are individually careful, several are excellent — but because the
set has been optimized for a property no user wants (internal consistency of a
document graph) at the cost of the only property that matters (a running
program). The cluster has produced 29,018 lines of RDR and evidence against
3,971 lines of Go, of which the entire product surface is a `version`
subcommand. The last commit touching `internal/` was 2026-08-09. On 2026-08-11
the project produced 62 commits, all specification. Three of the eight members
(0007, 0008, 0009) were seeded, proposed, refined, resolved, four-lens
pre-locked, reconciled, and finalized in a single day, and together they add
345KB of normative text that changes, by their own accounting, almost no
behavior. RDR 0007 states this outright: "this RDR's own probe still reproduces
on the shipped kernel after Phases 1–3 land, because no code path changes."

That is the shape of a project that has substituted specification for delivery,
and the failure below is the mechanical consequence.

---

## Findings Ledger

| ID | RDR | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-----|-------------|--------------|-------------------|--------|
| C-1 | 0007 | A10 "Obligation destination (NAMED): RDR 0003's `### Phase 1: Predicate Model` implement-stage step" | Obligation routed into a 2-sentence Phase in a Final peer that never accepted it; no mechanism forces acceptance | Guard evaluator ships with an ad-hoc guard encoding; absence silently folds to false; a plan is emitted where the honest answer was "I could not tell" | §1 |
| C-2 | 0007 | A12 "Verified — INDETERMINATE-BY-SILENCE, resolved as an inherited obligation, NOT a refutation" | "Verified" status assigned to an unanswered question; Gate PASS obtained on a hole | Author writes the blessed two-row absence pattern; graph lint rejects the table as `ambiguous overlap` with no stated fix | §1, §3 |
| C-3 | 0007 | A15 "Verified — as a NAMED OBLIGATION on RDR 0003's implement stage (silence, not prohibition)" | The disjunction pattern's value row has no grammar; the only sanctioned absence-disjunction route is unauthorable | Author cannot express "absent OR equals v"; falls back to sentinel-stamping, the anti-pattern the RDR forbids | §1, §3 |
| C-4 | 0007 | Risks: "Residual status: UNMITIGATED, not partially mitigated. Nothing in this repo fails if RDR 0003 never instantiates the harness" | Self-declared unenforced obligation shipped as Final | The conformance harness exists and is never run; drift is undetected until a wrong plan reaches an artifact | §1, premortem, AT-4 |
| C-5 | 0003 | `### Phase 1: Predicate Model` (2 sentences, "normalized predicate atom shape used by resolver and lint") | Receiving Phase is far too small to hold A10+A12+A15; Prerequisites are `[x]` checked and predate the obligations | 0003 implementer reads a checked-off Prerequisites list, builds the evaluator, and never learns three contracts bind them | §1, §2, AT-1 |
| C-6 | 0003 | Prerequisites "- [x] All Critical Assumptions verified" | Checked prerequisite list is stale relative to 0007/0008/0009's routed obligations; no invalidation mechanism | Implementation starts from a green checklist that is materially false | §2, AT-1 |
| C-7 | set | 0007 A10; 0007 A12; 0007 A15; 0007 A6b; 0009 Prereq "RDR 0002's implement stage binds to this RDR's conformance fixtures" | Cross-RDR obligations accumulate into an unowned queue; no RDR holds the integration contract | Nothing integrates; each RDR is individually Implemented and the product does not run | §1, premortem |
| C-8 | 0007 | Normative Contracts: "the only out-of-band surface that exists is a PANIC … the kernel MUST NOT recover the panic into a verdict or refusal" | A library seam is normatively required to crash the process on a producer defect | CLI exits with a Go panic and stack trace instead of a structured envelope; violates 0005's never-silent structured-error contract | §1, premortem, AT-6 |
| C-9 | 0005 | Normative Contracts: verbs MUST route failure through `respond.Fail(cmd, *clierr.CLIError)`; root.go "never-silent contract" | 0007's mandated panic bypasses the entire 0005 error envelope; no RDR reconciles the two | Operator sees a raw panic where every other failure is a typed JSON envelope; skill integration breaks on unparseable output | §1, AT-6 |
| C-10 | 0007 | A6b "Status: Pending — DOWNGRADED at Stage 6 by explicit decision … ROUTED to RDR 0004's implement stage" | Read-completeness gap parked on a peer that has not agreed; negative existence atoms are unsound until closed | A truncated accessor read makes `exists = false` decide TRUE; a plan is produced against artifact state that exists | §1, §3, premortem, AT-5 |
| C-11 | 0004 | Normative: "A read accessor MUST return typed tag values or a typed refusal" | Disjunction with no completeness requirement on the values branch; partial reads are unmodeled | Silent wrong plan from a partially-parsed artifact; no refusal, no diagnostic | §3, premortem, AT-5 |
| C-12 | 0007 | A13 "Verified — as an ACCEPTED EXPOSURE … an operator can unstick a `guard_unevaluable` refusal by supplying the missing tag on the command line" | The set's headline guarantee (missing state is unmaskable) is knowingly bypassable from the CLI | Operator hits a refusal, adds `--tag`, gets a plan; writes computed from typed-in state, nothing in output marks it | §1, §3, premortem, AT-3 |
| C-13 | 0005 | "`--tag name=value` … remains context already known to the caller"; no allowlist or accessor-origin requirement | The `--tag` channel is normatively unconstrained; 0007's follow-up seed is unchecked `- [ ]` | Same as C-12; the CLI is the delivery vehicle for the bypass | §3, AT-3 |
| C-14 | 0008 | Approach item 4: "producers of kernel `Input` … carry a producer obligation not to supply a tag keyed `recognized`" | Producer obligation with the CLI as an unconstrained producer; 0005 never learns of the reservation | Operator passes `--tag recognized=x`; either a reserved-key breach error from a legal-looking flag, or silent shadowing | §1, premortem, AT-7 |
| C-15 | 0008 | Normative: "A tag declaration with provenance `recognized` MUST be named `recognized`" | Withdraws naming freedom that all three committed fixtures exercised (`outcome`, `outcome`, `rewind_target`) | Every existing authored table fails load with `reserved_tag_key`; no migration tooling exists | §2, premortem, AT-8 |
| C-16 | 0008 | Normative: "a model whose rows were meant to match the recognized tag but which declares no recognized-provenance tag lints clean and refuses `no_match` at resolve time" | The RDR's own Problem Statement failure (silent no-match) is left open by the chosen design | Author's row never fires; lint is green; refusal says `no_match` and not why — the exact seed symptom | §2, premortem, AT-8 |
| C-17 | 0009 | Normative: exported `ErrEscapeShapeBreach`, `*EscapeShapeBreachError`, `Ref`, `Count`, `Table.CheckValid()`, `Unwrap() []error` EXACTLY ONE LEVEL, "Resolve MUST return CheckValid's error VERBATIM" | Implementation-level API minutiae frozen as normative contract for a defect with zero production reachability | Nothing user-visible; the cost is spec surface that must be honored and cannot be refactored | §2, AT-9 |
| C-18 | 0009 | Normative: "Resolve MUST return the zero `Result` … callers MUST check the error before reading the disposition; a caller that branches on `Refused()` first would read a success-shaped value with a nil `Plan`" | New nil-plan/success-shaped trap introduced by a hardening RDR; obligation pushed onto every caller | A future `flow` verb branches on `Refused()`, dereferences a nil `Plan`, and panics | §2, premortem, AT-10 |
| C-19 | 0009 | Normative multi-breach clause: "breaching rows sharing one RowRef contribute ONE reported error, not one per row" plus `Count` | Report collapses distinct rows because `compareRefs` cannot totally order them; hand-built rows all carry `RowRef{"",""}` | Producer with three broken rows gets one error naming `("","")` and `Count == 3`; cannot locate any of them | §2, AT-11 |
| C-20 | 0006 | Group membership never defined beyond "the same state/outcome" (cited by 0007 A12) | Row-group membership is undefined, so overlap/exhaustiveness cannot be computed for the pattern 0007 blesses | Lint is either unimplementable or silently non-deterministic in what it compares | §1, §3, AT-2 |
| C-21 | 0002 | "ambiguous overlap" retained among categories validation MUST reject, with no downgrade valve | Mandatory reject collides with 0007's blessed absence pattern under an undefined projection rule | Table author's legal pattern is hard-rejected with no relief and no documented fix | §1, AT-2 |
| C-22 | set | 0007 Consequences: "this RDR's own probe still reproduces on the shipped kernel after Phases 1–3 land"; kata `xg7p` "MUST be closed against that distinction explicitly rather than against a green test run" | RDRs reach Implemented while the defect they exist to fix stays open | Tracker shows the work done; the bug is still there; it is re-found later against a kernel that never changed | §1, §2, premortem, AT-12 |
| C-23 | set | `cluster-reconcile/0001-0006/report.md` — all seven pairwise rows "NO CONFLICT", verdict "RECONCILED" | The reconcile gate that should have caught cross-RDR gaps returned all-clear on a set that 0007 later proved was silent on three load-bearing questions | Gate passes; the set ships; the gaps are discovered at implementation | §1, §5 |
| C-24 | set | 0007/0008/0009 Metadata `Status: Final`; README index rows | Eight Final RDRs, one implemented; ratio of spec to code is 7:1 and rising with no delivery in between | No product. The user has a `version` command | §1, §2, premortem |
| C-25 | 0007 | Minimum Viable Validation: "A 'green MVV' therefore MUST NOT be read as validating the domain rule" | The RDR's own validation gate is declared non-validating; there is no gate that does validate before lock | Lock is granted on evidence the RDR itself says proves nothing about its content | §5, AT-13 |
| C-26 | 0007 | Existing Infrastructure Audit: "View-reading evaluator for the domain-rule scenarios — None … Phase 2 must build it alongside the vectors; it is not free reuse" | The conformance vectors that enforce the whole contract require an evaluator that does not exist and is blocked on C-1 | Phase 2 cannot complete; the enforcement mechanism never lands | §1, AT-4 |
| C-27 | 0008 | Technical Design: "the conformance test MUST run the real kernel and assert the assembled view binds the recognized outcome under the normalizer's spelling" | Test requires a normalizer (0002, unimplemented) and thus cannot run at 0008's own implement stage | 0008 ships "Implemented" with its only behavioral test unrunnable | §2, AT-14 |
| C-28 | 0005 | "Commands accept a model or flow identifier plus state tags supplied as repeated `--tag name=value` flags in the MVP" | 0005 is the sole user-facing member and sits downstream of every unresolved obligation in the set | The user's first contact with the product is the last thing built, after five spec-only dependencies | §2, §3, premortem |
| C-29 | 0004 | Prerequisites: all four items unchecked — "- [ ] All Critical Assumptions verified" — while Status is Final and the Gate asserts "Critical Assumptions A1-A7 are verified" | RDR reached Final with its own Prerequisites list contradicting its Finalization Gate | Implementer cannot tell whether 0004's preconditions hold; the gate record is unreliable as an input | §1, AT-16 |
| C-30 | 0004 | Gate/Assumption Verification: "None of the verified evidence cites this RDR or its artifact directory as self-proof" | False on its face — A1/A2/A6/A7 cite `docs/rdr/0004-.../evidence/spikes/`, and A5 cites 0004's own sections | The self-reference guard that is supposed to stop circular verification reports clean on a circular verification | §1, AT-16 |
| C-31 | 0006 | A5 "Status: Pending"; Prerequisites "- [ ] A5 production gate proof pending MVV scenario 7" | RDR whose entire value is *blocking* CI authority locked Final with the CI wiring unproven | Lint exists and CI does not run it; models are lintable locally and unenforced at the boundary maintainers rely on | §1, AT-17 |
| C-32 | 0005 / 0006 | 0005 A5 "not new envelope fields or exit groups" + Audit "no new envelope fields or exit groups required" vs 0006 N7 "MUST remain machine-readable … through an append-only optional typed `findings` envelope field" | Two Final RDRs give opposite answers on whether the CLI envelope gains fields; the 0001–0006 reconcile pass dispositioned this "NO CONFLICT" | Whichever ships second either breaks the other's contract or silently wins | §1, AT-18 |
| C-33 | 0002 / 0003 / 0006 | 0002 Failure Modes "read-before-write condition … is a lint failure"; 0002 Phase 5 hands it to 0006; 0003 Normative "A row that matches an owned tag MUST be rejected unless every reachable predecessor sets or preserves that tag"; 0006 invariant 6 `graph-owned-before-write` | One check claimed by three RDRs with no arbiter | Either implemented twice with divergent semantics, or each implementer assumes a peer owns it and it ships uncovered | §1, AT-18 |
| C-34 | 0001 / 0005 | 0001 refusal kinds `owned_state_unavailable`; 0005 error taxonomy table (twelve `flow-*` codes) | 0001's five kernel refusal kinds do not map 1:1 onto 0005's codes — `owned_state_unavailable` has no counterpart | Skill branching on a documented refusal class finds no stable code for it; falls back to string matching, which 0001 A5 exists to prevent | §1, AT-19 |

---

## 1. The inter-RDR failure mode

### The failure

**The set fails as a set because it has invented a mechanism — the "obligation
routed to a peer's implement stage" — that lets an RDR reach `Final` with a Gate
PASS while its load-bearing questions remain unanswered, by relocating those
questions into documents that never agreed to receive them and that no process
step will re-open.** The obligations accumulate in a queue nobody owns. At
implementation time the queue is invisible, because the receiving RDRs' own
Prerequisites are already checked off. The result is that every RDR ships
individually conformant and the composed system does not work.

This is not a hypothetical. It is the dominant structural feature of the newest
three members, and it is at its worst in RDR 0007.

### Root cause across the RDRs

The project holds a rule — stated in the user's own standing instruction and
enforced by the RDR engine — that **RDRs are never amended**. That rule is
correct in isolation: it preserves history and prevents retroactive rewriting.
But combined with a lifecycle in which RDRs are locked to `Final` *before*
implementation, it produces a trap. When RDR 0007 discovers, at Stage 6, that
RDR 0003 is silent on three questions 0007's contract depends on, 0007 has
exactly three legal moves:

1. Return RDR 0003 to Draft (the cluster gate's demotion path) — expensive,
   and it re-opens a document with a Gate PASS.
2. Decide the questions itself — but they are 0003's domain, and deciding them
   in 0007 would be the cross-RDR contradiction the reconcile pass exists to
   prevent.
3. Declare the assumption **Verified** with a status qualifier and route the
   obligation to 0003's *implement stage*, which is not a document and therefore
   not an amendment.

The set takes option 3 every time. RDR 0007's Prerequisites are a monument to
it — six assumptions closed at Stage 6, of which **four** are closed by
relocation rather than by answer:

- A10 → "RDR 0003's `Phase 1: Predicate Model` implement-stage step"
- A12 → "RDR 0003's implement stage must define the projection, and RDR 0006's
  must define row-group membership"
- A15 → "the same RDR 0003 `Phase 1: Predicate Model` step that owns A10's
  mapping"
- A6b → "RDR 0004's implement stage"

RDR 0009 does the same, more quietly, in its own unchecked Prerequisites:
"RDR 0002's implement stage binds to this RDR's conformance fixtures",
"RDR 0004's implement stage keeps the write accessor scoped".

The root cause, stated plainly: **the RDR set has a mechanism for discovering
cross-document gaps and no mechanism for closing them.** The reconcile stage
finds them; the no-amendment rule forbids fixing them at the source; the
"implement stage" is used as a landfill because it is the only destination that
does not require re-opening a Final document. And the landfill has no bottom,
because an implement stage is not a document that can be checked, reviewed, or
gated. Nothing anywhere fails if the obligation is never discharged. RDR 0007
says this itself, about its own central risk, in the plainest sentence in the
entire corpus:

> **Residual status: UNMITIGATED, not partially mitigated.** Nothing in this
> repo fails if RDR 0003 never instantiates the harness — that is the
> definition of an unenforced obligation, and it is the same hole as the risk
> it answers (P-10).

That is a Final RDR, Gate PASS, stating that its primary enforcement mechanism
is unenforced. It was locked anyway.

### The specific passages that enabled it

**Enabler 1 — the "Verified" status is redefined to mean "asked and not
answered."** RDR 0007 A12:

> **Status**: Verified — INDETERMINATE-BY-SILENCE, resolved as an inherited
> obligation, NOT a refutation. RDR 0003 does not say the pattern is rejected;
> it says nothing about the case at all.

An assumption whose evidence is "the peer document is silent" is not verified.
It is open. Calling it Verified is what allows the Finalization Gate's
"Assumption Verification" check to pass. A15 uses the identical construction
("Verified — as a NAMED OBLIGATION … (silence, not prohibition)"), as does A10
("Verified — as a NAMED OBLIGATION … The mapping is specified nowhere; it is
unowned, not impossible"). Three of RDR 0007's fifteen assumptions are marked
Verified on the strength of a demonstrated absence of evidence. The gate cannot
distinguish this from a real verification, so it passes.

**Enabler 2 — the receiving Phase is two sentences long and predates the
obligations.** RDR 0003 `### Phase 1: Predicate Model`, in its entirety:

> Define the tag-kind/operator compatibility matrix and normalized predicate
> atom shape used by resolver and lint.

RDR 0007 claims this Phase "already charters" its obligation and therefore "the
obligation lands inside 0003's existing scope with no amendment to a Final
peer." That claim is doing enormous work on a very thin sentence. Phase 1 as
written is a modeling task. What 0007 routes into it is: (a) a total, lossless
encoding for `Row.Guard string` that no RDR has ever specified (A10); (b) the
projection rule for existence atoms onto the declared-domain product, on which
graph lint's overlap check depends (A12); (c) whether one tag-table may carry
multiple conjoined operator keys and whether a tag may be guarded in both `all`
and `unless` (A15). Each is a design fork of the kind that, by this project's
own standards, would justify its own RDR. Three of them have been posted into a
two-sentence phase of a document whose Prerequisites read:

> - [x] All Critical Assumptions verified
> - [x] RDR 0002's sparse table/container contract is stable enough to host
>   guard atoms.

Both boxes checked. Nothing in RDR 0003 mentions 0007, A10, A12, or A15. The
implementer of 0003 opens a Final RDR with a clean checklist and no reason to
read 0007 at all.

**Enabler 3 — the cluster reconcile gate returned all-clear on this exact
seam.** `cluster-reconcile/0001-0006/report.md` compared four pairs and found
seven "NO CONFLICT" rows, verdict "RECONCILED", with the note:

> No open SPEC-DEFECT remains.

Among the pairs compared was 0002-0003, dispositioned "RDR 0002 owns the
recognized alphabet and tag declarations; RDR 0003 consumes recognized
provenance for predicate/lint reasoning." Six weeks later RDR 0007 established
by source search that RDR 0003 is silent on absent-operand behavior, on
existence-atom projection, on conjoined atom shape, and on the `Row.Guard`
encoding — and that RDR 0006 never defines row-group membership. The reconcile
pass looked for *contradictions* and found none, because there were none: the
documents do not contradict each other, they simply do not cover the ground
between them. **A gate that only detects contradiction cannot detect silence,
and silence is this set's actual failure mode.** The gate's clean bill of health
is what let the set proceed.

**Enabler 4 — the Finalization Gate is prose nobody diffs against the body.**
The gate is supposed to be the backstop that catches exactly this. It does not,
because its written responses are never checked against the document they
describe. Two mechanical greps prove it, both against RDR 0004:

*First.* RDR 0004 is `Status: Final`. Its Prerequisites list is:

> - [ ] All Critical Assumptions verified
> - [ ] RDR 0001 keeps the resolver stateless …
> - [ ] RDR 0002 carries accessor references …
> - [ ] RDR 0003 consumes accessor-produced tag values …

Four boxes, none checked. Its Finalization Gate, in the same file, asserts:
"Critical Assumptions A1-A7 are verified." The document contradicts itself on
its own lock condition and locked anyway. (0002's equivalent list is `[x]`, so
this is not a house convention — it is an oversight the gate did not catch.)

*Second.* The same gate section closes with:

> None of the verified evidence cites this RDR or its artifact directory as
> self-proof.

That sentence is false against its own assumptions. A1, A2, A6, and A7 all cite
`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes/` — the RDR's own
artifact directory — and A5's evidence cites 0004's own Cross-Cutting Concerns,
Normative Contracts, and Alternatives sections. The self-reference guard exists
precisely to stop an RDR from proving itself with its own artifacts. It reported
clean on an RDR that does exactly that.

If the gate cannot catch a false claim about the paragraph directly above it, it
cannot catch an obligation routed into a peer's implement stage. These are the
same failure at two magnitudes.

**Enabler 5 — the set contains outright contradictions the reconcile pass
scored "NO CONFLICT."** Two, both verified:

*The envelope.* RDR 0005 A5: "Resolver/accessor-specific values only require new
`Code` constants or literals, **not new envelope fields or exit groups**," and
its Infrastructure Audit repeats it: "no new envelope fields or exit groups
required." RDR 0006 N7: blocking findings "MUST remain machine-readable in JSON
mode through an append-only optional typed `findings` envelope field owned by
`clierr`/`respond`." One RDR says the envelope gains no fields; the other says
it MUST gain one. The 0001–0006 reconcile pass compared exactly this pair and
recorded: "RDR 0005's 'no new envelope fields' claim is scoped to
resolver/accessor failures; RDR 0006 owns the graph-lint findings exception."
That is a rationalization, not a resolution — 0005's text carries no such
scoping, and whichever RDR ships second either violates the other or silently
wins.

*Read-before-write.* Three RDRs claim one check. RDR 0002's Failure Modes lists
"read-before-write condition" as a lint failure it owns, then Phase 5 hands
"read-before-write checks" to RDR 0006. RDR 0003 states it normatively: "A row
that matches an owned tag MUST be rejected unless every reachable predecessor
sets or preserves that tag before the match." RDR 0006 makes it invariant 6 with
its own code, `graph-owned-before-write`. Three owners, no arbiter. It will be
implemented twice with divergent semantics, or zero times because each
implementer assumes a peer has it.

**Enabler 6 — a normative panic.** RDR 0007's Normative Contracts:

> An evaluator MUST surface it outside the verdict channel — never by answering
> GuardUnevaluable — and at this seam the only out-of-band surface that exists
> is a PANIC. … A crash is the honest surface precisely because the case is
> unreachable for lint-passed guards: it is a programmer defect, and the kernel
> MUST NOT recover the panic into a verdict or refusal.

The reasoning is internally sound and the conclusion is still wrong at the set
level. RDR 0005 requires every user-facing failure to route through
`respond.Fail(cmd, *clierr.CLIError)`, and `internal/cli/root.go` documents a
"never-silent contract" in which "ExecuteAndEmit converts any cobra-level error
into a CLIError." A panic from inside the guard evaluator satisfies neither. It
produces a Go stack trace on stderr, no envelope, no exit-code group, no
`Code`, and — for the skill-integration consumer that 0005 exists to serve —
unparseable output. The unreachability argument that justifies the panic
depends on A7, which depends on lint rejecting undeclared tags, which depends on
RDR 0003 and RDR 0006 being implemented and correct. The panic is reachable
precisely in the window where the set is half-built, which is the entire
delivery period. No RDR in the set reconciles the mandated panic with 0005's
envelope contract, because 0005 was locked before 0007 existed and 0007 cannot
amend it.

### The symptom the user sees

The user is a skill author driving `intrastate flow next` / `flow resolve`
against a transition table they wrote.

They see: **nothing, for a very long time.** The set's user-facing member (0005)
is the last in the dependency order and is blocked on 0002, 0003, 0004, and now
on obligations routed from 0007, 0008, and 0009. Today the binary answers
`version`.

Then, when 0003's evaluator finally lands, they see the failure this whole
cluster was written to prevent. The 0003 implementer, working from a checked
Prerequisites list, invents a guard encoding for `Row.Guard string` (nobody
specified one — A10), and for the value-comparison operators they take the
obvious reading that a missing tag means the comparison fails. Their evaluator
never runs 0007's conformance harness, because nothing makes it (C-4, the
self-declared unmitigated risk). The user writes a table with a modeled
`no_match` escape. An accessor read comes back partial (C-10/C-11 — RDR 0004
never required the values branch to be complete). The guard over the missing tag
evaluates FALSE, the row prunes, the escape fires, and the user gets a **plan**.

They act on it. The plan's writes are computed against artifact state that was
never read. Nothing in the output distinguishes this from a resolution that read
everything. That is verbatim the outcome RDR 0007's Problem Statement promises
to prevent:

> At that point the table looks like it worked, and nothing in the output
> distinguishes "this row does not apply" from "I could not read what I needed
> to find out."

The set will have spent 345KB of specification on that sentence and shipped the
behavior anyway, because the contract that forbids it was written in a document
the implementer had no reason to open.

---

## 2. The RDR that will be rewritten within six weeks of shipping

**RDR 0009 — Ownership of escape-row shape conformance.**

RDR 0007 is the most consequential member and 0003 is the most under-specified,
but 0009 is the one that will be *rewritten*, and quickly, because it is the
only member that commits fine-grained implementation surface to normative text
for a defect with zero production reachability. Everything about it is
calibrated to be revised the moment real code touches it.

**It freezes an API before there is a caller.** Its Normative Contracts fix, as
MUSTs:

> - the sentinel is ErrEscapeShapeBreach (package-level, exported);
> - the typed error is *EscapeShapeBreachError (exported), carrying exactly ONE
>   RowRef field, Ref, plus a Count int field (see the multi-breach clause);
> - its Unwrap() error returns ErrEscapeShapeBreach …
> - the predicate is Table.CheckValid() error.

Plus: `Unwrap() []error` MUST be "EXACTLY ONE LEVEL deep", and "Resolve MUST
return CheckValid's error VERBATIM, never fmt.Errorf-wrapped." These are code
review comments promoted to constitutional law. They are the kind of decision
that is correct in the abstract and wrong in three specific ways as soon as a
real consumer exists — and by A1's own evidence, **there is no consumer**: "no
production file imports the package." The first `flow` verb that has to render
this to an operator will want the aggregate wrapped with command context, which
the "VERBATIM, never fmt.Errorf-wrapped" clause forbids. That is a rewrite.

**It creates a new footgun while closing an old one.** The RDR mandates:

> Resolve returns the zero `Result` (`Plan` and `Refusal` both nil …). Because
> the zero `Result` reports `Refused() == false`, callers MUST check the error
> before reading the disposition; a caller that branches on `Refused()` first
> would read a success-shaped value with a nil `Plan`.

The RDR identifies the trap, states it clearly, and then ships it, pushing the
obligation onto every future caller. RDR 0001's entire five-return-path design
holds `error` at literal `nil` so that disposition is total; 0009 makes it
partial and mitigates with a comment. The first `flow resolve` implementation
that branches on `Refused()` first — the natural order, since that is what every
existing test does — dereferences a nil `Plan`. That will be a bug report, and
the fix will be a `Result` shape change, which is a rewrite of this RDR's
central clause.

**Its diagnostic collapses on its own target population.** The multi-breach
clause is forced into a corner by RDR 0001's frozen row-order-independence
tests (ADV-3b, Fixup-3c) and resolves it by collapsing:

> breaching rows sharing one RowRef contribute ONE reported error, not one per
> row. … Rows carrying no source identity collapse to a single RowRef{"",""}
> entry

and concedes in the same block that hand-built rows "are this precondition's
whole target population (A5) and collide trivially: rows built with no source
identity all carry RowRef{"",""}." So for the *entire* population the check
exists to serve, the report is one error naming `("","")` with a count. A
producer with three broken rows learns that three rows are broken and cannot
locate any of them. The first person who hits this files a bug, and the fix —
carrying a locator, or an index, or anything that discriminates — collides with
the REQ-2/REQ-10 determinism rule the clause was contorted to satisfy. That is
not a tweak; it reopens the RDR's hardest paragraph.

**And the trigger is cheap.** 0009's whole subject is a class of bug that,
per A3's spike, requires editing exactly three lines of test fixtures to make
impossible today. The cost/benefit will look indefensible in retrospect: 1,969
lines of RDR, a locked exported API, a new nil-plan trap, and a collapsing
diagnostic, to prevent a defect that no production code can currently reach.
Six weeks after a real `flow` verb exists, someone will open this file to make
the error renderable and discover that four separate MUSTs forbid the natural
implementation. It gets rewritten.

*(Runner-up, for the record: RDR 0008. Its own normative text concedes that a
model which meant to match the recognized tag but declared nothing "lints clean
and refuses `no_match` at resolve time" — the exact silent failure from its
Problem Statement, left open. And all three committed fixtures name that tag
`outcome` or `rewind_target`, so the reserved-key rule breaks every existing
example with no migration path. But 0008 will be *patched*; 0009 will be
*rewritten*.)*

---

## 3. The cross-cutting assumption that will not survive first contact

**The assumption: that the tags reaching the kernel are a faithful, complete
picture of artifact state — so that "key absent from the view" means "the
artifact does not have this state," and every contract built on presence and
absence is sound.**

This is not written as an assumption anywhere, which is precisely why it will
break. It is the load-bearing floor under the entire cluster, distributed across
the set in pieces so that no single RDR owns it and no single verification
covers it:

- **RDR 0007** makes presence/absence the *entire semantics* of guard
  evaluation. The domain rule, `exists` totality, strong-Kleene combination, the
  aggregation veto — all of it quantifies over whether a key is in the view.
- **RDR 0007 A6a** verifies only the narrow half: `assemble` does not drop keys
  *from the input tuple*. Faithfulness to the tuple, not to the artifact.
- **RDR 0007 A6b** identifies the real gap and leaves it open: "whether a
  *failed* accessor read is distinguishable from *genuine* absence at the
  accessor→kernel boundary is not stated by any current contract." Status:
  Pending, routed to RDR 0004's implement stage.
- **RDR 0004** — the RDR that owns the boundary — says only: "A read accessor
  MUST return typed tag values or a typed refusal." 0007's own analysis of that
  clause is devastating and correct: "This is a disjunction with **no
  completeness requirement** on the 'typed tag values' branch: nothing says the
  returned values are the complete tag set, and nothing requires a
  partially-read artifact to take the refusal branch."
- **RDR 0005** then opens a second, wider hole in the same floor by making the
  caller a first-class producer of view contents: `--tag name=value` is
  "context already known to the caller," with, per 0007 A13, "no allowlist,
  validation, or accessor-origin requirement."
- **RDR 0008** adds a third producer channel and defends it only by producer
  obligation: "producers of kernel `Input` … carry a producer obligation not to
  supply a tag keyed `recognized`."

Each RDR has verified its own edge. Nobody verified the floor.

### How it breaks

**First contact, path one — the partial read.** A real artifact is a file on
disk written by another tool mid-operation. It is truncated, or it is a YAML
document whose third key failed to parse, or the accessor timed out after
reading two of four tags. Under RDR 0004 as written, returning the two tags it
got is *conformant*. `Input.Owned []Tag` has no error channel. At the kernel
boundary, per 0007 A6b, "a truncated snapshot and a genuinely-absent tag are
byte-identical."

Now run 0007's contract over it. A value-comparing atom on a dropped key is
`GuardUnevaluable` — annoying but honest. But a **negative existence atom**
(`exists = false`) — the one construct 0007 A14 blesses as the sanctioned
single-atom absence test, the one it points authors toward to avoid
sentinel-stamping — decides **TRUE**. The row fires. A plan is produced against
state that exists and was not read. 0007 names this outcome exactly ("a
dropped-but-real tag makes a negative existence atom decide TRUE, yielding a
plan against state that exists") and ships anyway, because the fix belongs to a
peer that has not agreed to make it.

**First contact, path two — the operator.** An operator hits `guard_unevaluable`.
The refusal, per 0007's own Failure Modes, does *not* name the absent tag: it
carries `Refusal.Guard` (one row's guard text) and `Refusal.Rows`. The operator
reads the guard text, sees it mentions `cluster_eligible`, and does the obvious
thing:

```
intrastate flow resolve --flow x --tag cluster_eligible=true
```

Because presence is provenance-blind (0007's normative contract, forced by
`TagSet.Lookup` being the only exported accessor), the tag is now present. The
guard decides. A plan is produced whose writes derive from a value the operator
typed. 0007 A13 documents this in full and classifies it as an **accepted
exposure**:

> "the refusal is non-escapable" is a weaker guarantee than it reads: it is
> unroutable *around within the table*, not unbypassable by the caller.

The mitigation is a follow-up seed against RDR 0005, recorded as `- [ ]`,
unchecked, not a lock blocker. The kernel already refuses this exact
substitution on the owned-state path — `TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot`,
`TestAdv5_ObservedTagCannotSatisfyAnOwnedStateRequirement` — which means the
project identified this attack, wrote adversarial tests against it, and then
left the identical hole one seam over with no test contending it.

The set's single headline guarantee — *missing artifact state can never be
masked* — is defeated by one command-line flag, and the set knows it.

**First contact, path three — the collision.** The operator, or a skill, passes
`--tag recognized=approved`. RDR 0008 reserves that key by producer obligation
and enforces it with a kernel-exported predicate over `Input`. RDR 0005, which
defines the `--tag` flag, was locked before 0008 existed and knows nothing about
it. So the operator either gets a `reserved_tag_key` breach from a flag the CLI
documentation presents as free-form, or — if the wiring misses, since 0008's
enforcement is a predicate producers "may call" — silent shadowing resolved by
D3 precedence *against* the declared owner, which is premortem P-1/P-6 in 0008's
own text.

Three producers, three unenforced obligations, one floor that nobody owns.
Every RDR verified its own contribution; the composition is unverified and
unsound.

---

## 4. Premortem

*Written from twelve months after the cluster shipped.*

We shipped all eight. Every RDR reached `Implemented`. The index is green. The
product does not work, and the reason is that we spent a year proving documents
consistent with each other instead of proving code consistent with reality.

**Month 1–3. The queue nobody owned.** Implementation started in dependency
order: 0002, then 0003. The 0002 normalizer took eight weeks — longer than
planned, because `normalizeRow` had to satisfy four masters (0002's own grammar,
0008's `reserved_tag_key` category, 0009's conformance fixtures, 0007's
`<clear>`-renders-as-write claim) and no document listed all four in one place.
The 0009 fixture binding was found late; it was in 0009's Prerequisites as
`- [ ] RDR 0002's implement stage binds to this RDR's conformance fixtures`,
which nobody read because they were reading 0002.

Then 0003. The implementer opened `0003-guard-predicate-exhaustiveness.md`,
read Prerequisites — `- [x] All Critical Assumptions verified` — and started on
`Phase 1: Predicate Model`, two sentences long. They defined the operator matrix
and the atom shape. They did not know that 0007's A10, A12, and A15 had been
routed into that phase, because 0003 does not mention 0007. They needed a
representation for `Row.Guard string` and, finding none, chose an infix
expression string, matching the `"reviews >= quorum"` convention already sitting
in `fixtures_test.go`. They wrote `evaluateAtom` to compare a tag value against
a literal, and for a missing key they returned false — which is what every
programmer writes, and what SCXML §5.9.1 does, and what the 0002 fixtures'
authors clearly assumed. `go test ./internal/resolve/` was green, because
0007's conformance vectors were parameterized over a `GuardEvaluator` and ran
against `fixtureGuards`, the stub. Nothing bound the new evaluator to the
harness. 0007 had said so in writing: "Nothing in this repo fails if RDR 0003
never instantiates the harness."

**Month 4. The absence pattern died in lint.** The first real table author
needed a rule for "no reviewer assigned yet." They read the authoring guide,
which taught 0007 A5's blessed two-row pattern for "absent OR equals v": one row
guarded on `exists = false`, one guarded on the value. They wrote it. `intrastate
lint` rejected the table: **ambiguous overlap**. RDR 0006's row-group membership
had been implemented as "same source state and recognized outcome" — the only
definition available — so both rows landed in one compared product; RDR 0003's
overlap check computed accepted assignments over the declared domain
(`reviewer` had no absent member, because RDR 0002's tag-declaration schema has
no optionality field); and RDR 0002 lists `ambiguous overlap` among the
categories validation MUST reject, with no downgrade valve. Exactly A12.

The author asked how to fix it. There was no answer, because the projection rule
had never been written. They did what the guide called an anti-pattern: they
stamped `reviewer = "__none__"` upstream so the key was never absent. Six other
tables copied the idiom. Partiality was dead in the field within a month of the
first table, and RDR 0007's Risks section had predicted it word for word: "an
anti-pattern warning with no working alternative is a prohibition, not a
mitigation."

**Month 6. The incident.** `flow resolve` on the release-gating flow produced a
plan that marked a build promotable. It was not. The postmortem chain:

`readAccessor.Read` on `release.yaml` hit its timeout after parsing two of five
keys and returned the two it had, which RDR 0004 permits — "typed tag values or
a typed refusal," no completeness requirement on the values branch (A6b, routed
to 0004's implement stage, never discharged; the implementer had no reason to
know the obligation existed). `Input.Owned` carried the truncated set with no
error channel. `assemble` faithfully built a view missing `signoff_received` —
A6a holds, and is irrelevant. The row guarded `signoff_received exists = false`,
the sanctioned absence idiom from A14. It decided **TRUE**. `gate` kept it,
`missingOwned` found nothing (0007 had narrowed `RequiresOwned` to post-guard
write dependencies, so the guard-read key was not listed — by design), no
sibling was unevaluable, `Resolve` returned exactly one match, and `planOf`
built a plan. `flow set-state` applied the writes.

Nobody had asked for a plan. The system had been asked whether signoff existed,
could not read the file, and answered "no." Every RDR behaved as specified.

**Month 7. The masking bug we wrote 2,400 lines to kill, still alive.** The
incident review reopened kata `xg7p`. It had been closed against RDR 0007 —
correctly, per 0007's Consequences: "The tracking kata (`xg7p`) is resolved by
the contract, not by the behavior, and it MUST be closed against that
distinction explicitly rather than against a green test run — otherwise the
same finding is re-derived later against a kernel that never changed." It was
re-derived later, against a kernel that never changed. The distinction had been
recorded in an RDR nobody read at closing time.

**Month 8. The operator bypass.** With `flow` verbs live, on-call learned the
trick. A `guard_unevaluable` refusal names the guard text but not the absent
tag (0007's Diagnostic gap). The runbook grew an entry: *"read the tag out of
the guard, pass it with `--tag`, re-run."* It worked every time, because
presence is provenance-blind and `--tag` is normatively unconstrained "context
already known to the caller" (0005), and RDR 0001 normatively *affirms*
"caller-supplied observed tags." The non-escapable refusal — the entire product
outcome of RDR 0007 — became a speed bump with a documented workaround. A13 had
called it: an accepted exposure whose follow-up seed against 0005 was left
`- [ ]`.

**Month 9. The panic.** A `flow next` invocation hit an evaluator mapping
failure on a guard that had passed an earlier lint version. Per 0007's normative
contract the evaluator panicked, because "the only out-of-band surface that
exists is a PANIC" and "the kernel MUST NOT recover the panic into a verdict or
refusal." The skill calling `intrastate flow next --as json` got a Go stack
trace on stderr and a non-envelope exit. The skill's JSON parse failed, and it
retried in a loop. RDR 0005's never-silent structured-error contract had never
been reconciled with 0007's mandated crash — 0005 was locked first and could not
be amended.

**Month 11. The rewrite.** `EscapeShapeBreachError` was rewritten. The `flow`
verb needed to wrap the aggregate with command context; "MUST return CheckValid's
error VERBATIM, never fmt.Errorf-wrapped" forbade it. A producer with four
broken rows got one error reading `RowRef{"",""}` with `Count == 4` and could not
find any of them, because the report collapses equal identities and hand-built
rows are the whole target population. And `flow resolve` had shipped a nil-plan
panic in week two, from branching on `Refused()` before checking `err` — the
trap 0009 documented and shipped anyway.

**What actually killed us.** Not one bad decision — the individual decisions
were mostly right, and were argued better than most projects ever argue
anything. What killed us was that between 2026-06-17 and the first user, we
wrote 29,018 lines of specification and 3,971 lines of Go, and the only
executable thing anyone could run for over a year was `intrastate version`. We
built an elaborate apparatus for proving documents consistent with one another
and never once got a table author in front of a working `flow next`. Every gap
that mattered — the absence pattern dying in lint, the partial read, the `--tag`
bypass, the panic in the envelope — would have been found in an afternoon by
one real user with one real table. We had the mechanism to find them; it was
called Stage 8, and we kept deferring it in favor of Stage 7.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are RDR-review-time gates, not code tests. Each is executable against the
document set as it stood, before lock. Each names the finding it catches.

### AT-1 — A routed obligation must be accepted by its destination
*Catches C-1, C-5, C-6*

```gherkin
Given RDR 0007 states an "Obligation destination (NAMED)" pointing at
  RDR 0003's "Phase 1: Predicate Model"
When the Finalization Gate runs for RDR 0007
Then RDR 0003's own text MUST contain a reference to RDR 0007 and to the
  obligation by identifier (A10)
And RDR 0003's Prerequisites MUST carry an unchecked item naming it
And if RDR 0003 is Final and cannot be amended, the gate MUST FAIL and route
  the obligation to a demotion of RDR 0003 to Draft, or to a new RDR
```
An obligation is not routed until the receiver's document shows it. "The
implement stage will read it" is not a destination; an implement stage has no
text and cannot be reviewed. Run today against 0007's A10/A12/A15 → three
failures.

### AT-2 — Every blessed authoring pattern must survive load-time lint
*Catches C-2, C-20, C-21*

```gherkin
Given an RDR blesses an authoring pattern (0007 A5's two-row absence pattern)
When the Gate runs
Then the pattern MUST be traced through every load-time validation the set
  defines: RDR 0002 validation categories, RDR 0003 overlap/exhaustiveness,
  RDR 0006 row-group membership
And for each, the RDR MUST cite the passage that ADMITS the pattern
And a citation of the form "the peer is silent" MUST FAIL the gate
```
0007 A12 says in terms that RDR 0003 "never states how an existence atom
projects onto the declared-domain product" and that RDR 0006 "never defines
group membership beyond 'the same state/outcome'". Under AT-2 that is a FAIL,
not a Verified assumption.

### AT-3 — Every guarantee must be tested against the CLI producer channel
*Catches C-12, C-13*

```gherkin
Given RDR 0007 claims missing artifact state "can never be masked behind an
  escapable refusal class"
When the Gate evaluates the guarantee
Then it MUST be evaluated against every producer of Input named in the set,
  including RDR 0005's "--tag name=value"
And if any producer can supply the state that converts the refusal into a
  decided verdict, the guarantee MUST be restated in the Problem Statement
  with its actual scope, or the producer channel MUST be constrained before
  lock
```
0007 A13 already performs this analysis and reaches the right conclusion —
"unroutable *around within the table*, not unbypassable by the caller" — and
then leaves the fix as an unchecked follow-up seed. AT-3 makes the restatement
or the constraint a lock condition rather than a footnote.

### AT-4 — A conformance harness must have a failing consumer
*Catches C-4, C-26*

```gherkin
Given RDR 0007 Phase 2 exports a conformance harness taking a GuardEvaluator
When the Gate asks how drift is prevented
Then there MUST exist a check that FAILS if a GuardEvaluator implementation
  exists in the tree and does not instantiate the harness
And "RDR 0003's implement stage MUST instantiate it" MUST NOT satisfy this
And an RDR that records "Residual status: UNMITIGATED" against its primary
  enforcement mechanism MUST NOT reach Final
```
An interface-satisfaction sweep in CI, or a build-tagged test that enumerates
`GuardEvaluator` implementers, would discharge this in a dozen lines. 0007
already diagnoses the hole precisely; it just declines to close it.

### AT-5 — The read-completeness floor must be owned before presence semantics lock
*Catches C-10, C-11*

```gherkin
Given RDR 0007's contract makes presence/absence load-bearing
And RDR 0007 A6b records that a failed read is indistinguishable from
  genuine absence
When the Gate runs for RDR 0007
Then RDR 0004's read-accessor clause MUST state whether "typed tag values"
  carries a completeness guarantee
And until it does, any TOTAL operator over absence (exists = false) MUST be
  marked unsound-for-production in the authoring guidance
And the RDR MUST NOT lock with the assumption in "Pending"
```
0007 locks with A6b Pending and simultaneously blesses `exists = false` as the
sanctioned absence idiom (A14). Those two facts are incompatible and AT-5 is
what surfaces the incompatibility.

### AT-6 — No normative contract may bypass the CLI error envelope
*Catches C-8, C-9*

```gherkin
Given RDR 0005 requires every user-facing failure to route through
  respond.Fail(cmd, *clierr.CLIError)
When any RDR mandates a failure surface (0007's "the only out-of-band surface
  that exists is a PANIC")
Then the Gate MUST trace that surface to a clierr.CLIError with a stable Code
  and an exit group
And a surface that produces no envelope MUST FAIL the gate
And "the case is unreachable" MUST NOT satisfy it, because unreachability
  depends on peer RDRs not yet implemented
```

### AT-7 — Reserved names must be enforced at every producer the set defines
*Catches C-14*

```gherkin
Given RDR 0008 reserves the tag key "recognized"
When the Gate runs
Then every producer of kernel Input named anywhere in the set MUST be
  enumerated, and for each, the enforcement point named
And RDR 0005's --tag flag MUST be shown either to reject the reserved key or
  to be covered by the exported Input predicate at a named call site
And a "producer obligation" with no enumerated enforcement point per producer
  MUST FAIL
```

### AT-8 — A naming rule that breaks existing artifacts needs a migration
*Catches C-15, C-16*

```gherkin
Given RDR 0008 requires recognized-provenance declarations be named
  "recognized"
And A4 records that all three committed recognized-provenance declarations
  chose an author name (outcome, outcome, rewind_target)
When the Gate runs
Then the RDR MUST carry a migration item covering every non-conforming
  artifact in the tree
And the RDR MUST state how an author who declares NO recognized tag but whose
  rows meant to match it is diagnosed
And "lints clean and refuses no_match at resolve time" MUST FAIL, being the
  Problem Statement's own symptom left open
```

### AT-9 — Exported API surface may not be frozen before a consumer exists
*Catches C-17*

```gherkin
Given RDR 0009 fixes ErrEscapeShapeBreach, *EscapeShapeBreachError, Ref,
  Count, Table.CheckValid, and the exact Unwrap() shape as normative
And A1 records that no production file imports internal/resolve
When the Gate evaluates proportionality
Then normative clauses naming exported identifiers MUST cite at least one
  consumer that requires that exact spelling
And absent a consumer, the identifiers MUST be recorded as implementation
  latitude, not as MUSTs
```
0008 gets this right for the spelling constant ("implementation latitude the
implementer may resolve"); 0009 gets it wrong for four identifiers plus a
wrapping prohibition.

### AT-10 — No RDR may introduce a documented footgun
*Catches C-18*

```gherkin
Given RDR 0009 states "a caller that branches on Refused() first would read a
  success-shaped value with a nil Plan"
When the Gate reads a hazard the RDR states about its own design
Then the RDR MUST either eliminate the hazard or carry a test obligation that
  fails when a caller commits it
And documenting the hazard in prose MUST NOT close it
```

### AT-11 — A diagnostic must discriminate on its own target population
*Catches C-19*

```gherkin
Given RDR 0009's precondition targets hand-constructed rows (A5)
And such rows carry RowRef{"",""}
When the Gate evaluates the multi-breach report
Then the report MUST distinguish N breaching rows on the target population
And a report that collapses the entire target population to one entry MUST
  FAIL, count field notwithstanding
```

### AT-12 — Implemented must mean the defect is closed
*Catches C-22*

```gherkin
Given RDR 0007 states its own probe "still reproduces on the shipped kernel
  after Phases 1-3 land"
When the RDR would be marked Implemented
Then the status MUST record the tracking defect (kata xg7p) as OPEN, with the
  named condition that closes it (RDR 0003's evaluator passing the Phase 2
  harness)
And the index MUST NOT show a green row for an RDR whose defect is open
```

### AT-13 — The MVV must validate the contract
*Catches C-25*

```gherkin
Given RDR 0007 states "A 'green MVV' therefore MUST NOT be read as validating
  the domain rule"
When the Gate evaluates Minimum Viable Validation
Then the MVV MUST exercise at least one clause of the RDR's Normative
  Contracts end to end
And an RDR whose MVV validates none of its own normative content MUST NOT
  reach Final; it is a design note, not a locked contract
```

### AT-14 — Every conformance test must be runnable at its own implement stage
*Catches C-27*

```gherkin
Given RDR 0008 requires a conformance test that "MUST run the real kernel and
  assert the assembled view binds the recognized outcome under the
  normalizer's spelling"
And no normalizer exists (RDR 0002 unimplemented)
When the Gate runs
Then every test an RDR mandates MUST be runnable using only components
  available at that RDR's own implement stage
And a test blocked on an unimplemented peer MUST be recorded as a deferred
  obligation on the index, not as this RDR's validation
```

### AT-15 — The Gate must be checked against the RDR's own body
*Catches C-29, C-30*

```gherkin
Given RDR 0004 Status is Final
When the Finalization Gate runs
Then no Prerequisites checkbox may be unchecked while the Gate asserts the
  same condition holds
  (0004: "- [ ] All Critical Assumptions verified" vs Gate: "Critical
  Assumptions A1-A7 are verified")
And the self-reference claim MUST be verified mechanically against every
  assumption's Evidence path
  (0004 Gate: "None of the verified evidence cites this RDR or its artifact
  directory as self-proof" — false for A1, A2, A6, A7, and A5)
And a Gate whose written responses contradict the body MUST FAIL
```
This is two mechanical greps. Both would have failed 0004 at lock. The Gate is
currently a prose assertion nobody diffs against the document it describes.

### AT-16 — An authority RDR may not lock with its authority unproven
*Catches C-31*

```gherkin
Given RDR 0006 exists to make graph lint a BLOCKING acceptance gate
And A5 ("CI can run intrastate lint as the blocking authority") is Pending
And Prerequisites carry "- [ ] A5 production gate proof pending MVV scenario 7"
When the Gate runs
Then an RDR whose central claim is enforcement MUST NOT reach Final with the
  enforcement path unproven
And the Pending assumption MUST block lock rather than be deferred to the MVV
```
0006 is the only member that could have caught the others' defects at design
time, and it locked with its own enforcement wiring unproven. An unenforced
blocking gate is an advisory gate with a stronger adjective.

### AT-17 — Every shared surface must have exactly one named owner
*Catches C-32, C-33, C-20*

```gherkin
Given two or more RDRs make claims about the same surface
When the cluster reconcile pass runs
Then for each surface it MUST record exactly one owning RDR
And contradictory claims MUST FAIL rather than be dispositioned "NO CONFLICT"

Examples of surfaces this set leaves unowned or multiply-owned:
  | surface                        | claimants            |
  | clierr envelope field set      | 0005 (no new fields) vs 0006 (findings field MUST) |
  | owned-tag read-before-write    | 0002 -> 0003, 0002 -> 0006, 0003 claims, 0006 claims |
  | row-group membership           | 0003 defers to 0006; 0006 never defines it |
  | guard structure at Row.Guard   | nobody (0007 A10)    |
```
The 0001–0006 pass compared four pairs. This set has at least four
multiply-claimed surfaces, none of which is a *contradiction* in the narrow
sense the pass tested for — 0005 and 0006 do not mention each other's clause —
and all of which break at integration.

### AT-18 — Every kernel refusal kind must have a CLI code
*Catches C-34*

```gherkin
Given RDR 0001 pins exactly five refusal kinds and requires each be "stable
  enough for RDR 0005 to map to a CLI error code without inspecting error
  strings"
When the Gate runs for RDR 0005
Then the error taxonomy table MUST carry one code per kernel refusal kind
And owned_state_unavailable MUST have a named counterpart
And a missing mapping MUST FAIL, because the fallback is the error-string
  branching RDR 0001 A5 exists to prevent
```

### AT-19 — The set-level delivery gate
*Catches C-7, C-23, C-24, C-28*

```gherkin
Given a cluster of Final, unimplemented RDRs
When any member would be locked to Final
Then the cluster reconcile pass MUST scan for SILENCE, not only for
  contradiction: for each cross-RDR obligation, the destination document must
  contain matching text
And the ratio of specification lines to shipped product lines MUST be
  reported
And no more than N members may be Final-and-unimplemented at once
  (recommend N=2)
And at least one member MUST have reached Implemented, with a user-facing
  surface, before a new member may be seeded
```

This is the one that actually matters. The 0001–0006 reconcile pass compared
four pairs, found seven "NO CONFLICT" rows, and declared "No open SPEC-DEFECT
remains" — on a set that RDR 0007 later proved was silent on the guard encoding,
existence-atom projection, conjoined atom shape, row-group membership, and read
completeness, and that carries at least two live contradictions (the `findings`
envelope field, read-before-write ownership) the pass scored clean. **The gate
looked for contradiction and the failure mode was silence — and where there
*was* contradiction, the pass rationalized it.** AT-19 inverts the question: not
"do these documents disagree?" but "for every obligation one document places on
another, does the receiver's text show it?" Run against this cluster today, it
fails on at least seven obligations and on a 7:1 spec-to-code ratio with one
Implemented member and no user-facing surface.

---

*Ledger rows: 34. Acceptance tests: 19. Sections 1–5 above are the argument;
the ledger is the routable record. Any defect not carried in the ledger is
unreported by construction.*

---

## Coda: the one-line diagnosis

Eight Final RDRs. One implemented. 29,018 lines of specification, 3,971 lines
of Go, and a binary that answers `version`. Four obligations routed into a
two-sentence phase of a document that never agreed to receive them. One RDR
locked with all four Prerequisites unchecked and a Gate asserting the opposite
(0004). One RDR whose entire purpose is blocking enforcement, locked with the
enforcement unproven (0006). Two RDRs giving opposite answers about the CLI
envelope, dispositioned "NO CONFLICT" (0005/0006). One check claimed by three
owners (read-before-write). And a headline guarantee — missing artifact state
can never be masked — that the set itself documents as defeatable by one
command-line flag (0007 A13).

None of this is a competence problem. The analysis in these documents is
better than most shipped systems ever receive; 0007's A2 truth-table
derivation and 0009's A3 mutation spike are genuinely first-rate. That is
exactly the diagnosis. **The project has become extremely good at the activity
it is measured on, and the activity it is measured on is not shipping.** The
gates all pass because the gates test documents against documents. The one
test nobody runs is a real table author in front of a working `flow next`, and
every defect above dies on first contact with that test.

Stop writing RDRs. Implement 0002, 0003, and 0005 — the minimum path to a
user-visible `flow next` — and let the four routed obligations be answered by
the code that has to honor them. The set will teach you more in one week of
implementation than in another 29,000 lines of specification.
