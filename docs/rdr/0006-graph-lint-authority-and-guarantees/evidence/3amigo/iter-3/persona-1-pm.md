Model: claude-opus-5[1m]

# Persona 1 - Product Manager

Question: does this RDR actually deliver the user outcome?

Findings:

1. **`Technical Design` — "Each finding carries a stable code…" paragraph, and
   the `Disposition of every model class lint can meet` table (all rows) — no
   report-everything-in-one-pass duty; the disposition table reads as one
   verdict per class.**
   RDR 0003's `Normative Contracts` carry a binding clause this RDR is the
   consumer of: "Lint MUST report every defect it can decide in one pass over a
   row group, not the first one it encounters. Withholding a group's
   exhaustiveness claim MUST NOT suppress overlap findings, coverage findings,
   or further withholding reasons for that same group: each unprovable
   dimension, each refusing row, each overlapping pair, and any coverage gap
   over a provable product is its own finding. A group with two refusing rows
   MUST emit a finding for each." RDR 0006 adopts RDR 0003's grouping,
   escape-row, narrowing, and projection clauses by explicit citation, but
   never cites or restates this one. The disposition table instead presents
   input classes as mutually exclusive rows each with a single "Finding minted"
   cell, and desk-trace step 7 resolves a two-clause collision by picking one
   verdict (`graph-unprovable-coverage`) rather than emitting both. Only
   scenario 8 hints at co-emission ("overlap findings among the group's
   decidable rows are still emitted"), and it is a scenario expectation, not a
   normative duty.
   The user outcome at stake is the actual reason a maintainer runs a lint:
   fix the graph in one edit-run cycle. A lint that reports the first defect
   per group turns a five-defect model into five CI round-trips, and nothing in
   this RDR forbids that implementation.
   **Blocks**: whether the implementer builds a collect-all findings engine or
   a first-failure-per-group engine — and whether the MVV fixture matrix must
   include a multi-defect group asserting every expected code, which it
   currently does not (the illegal matrix is one defect class per fixture).

2. **`Technical Design` — "Non-blocking findings are informational: the
   advisory tier is scoped to redundant rows, unreachable rules, and the
   bare-escape coverage closure" / Finding-codes table (only
   `graph-coverage-closed-by-escape` listed as info) / disposition-table row
   "Redundant row / unreachable rule".**
   Two of the three named advisory classes have no stable code, no invariant
   number, and no MVV assertion. This is not a cosmetic omission: RDR 0002's
   `Normative Contracts` explicitly delegates the duty here — "a rule whose
   predicate set requires `recognized` to be absent is dead, which is RDR
   0006's unreachable-rule finding, not a load failure here." A peer has
   routed a user-visible outcome to this RDR and this RDR named it in prose
   only. The finding-codes table is the artifact an implementer and a test
   author both read as the closed taxonomy; a class present in prose and absent
   from the table will be built as nothing.
   Compare the care given to the third advisory class
   (`graph-coverage-closed-by-escape`): named code, table row, disposition row,
   MVV assertion, scenario 9c. The asymmetry shows the other two were not
   actually designed.
   **Blocks**: whether an author who writes a dead rule gets any signal at all;
   and whether the implementer must mint codes for redundancy/unreachability or
   may ship the advisory tier with exactly one member.

3. **`Critical Assumptions` A6 (`Status: Pending`) plus `Implementation Plan >
   Prerequisites` unchecked A6 item, against `Capability Dependencies` row
   "Initial owned state and terminal declarations | RDR 0002 (requested; arm
   decided Stage 4)".**
   The reachability root and terminal stop set are load-bearing for invariants
   2, 6, and 7 and for the `Load-Bearing Decisions` reachability relation — the
   one graph property this RDR claims as its unique contribution ("this RDR
   contributes the one graph property — reachability — those peers cannot state
   for themselves", `Decision Rationale`). A6 correctly states that no peer
   declares either, and correctly decides the RDR 0002 arm. What it does not
   state is what the lint does, or what the user gets, in the window before
   that RDR 0002 clause lands — and A6 explicitly says the gap "does not gate
   this RDR's remaining pre-lock lenses."
   So this RDR can lock, and be implemented, while three of its seven mandatory
   invariants are inert. A6's own `If wrong` concedes the shape of it:
   "invariant 6 is vacuous and invariant 2 cannot distinguish a designed stop
   from a dead end." An inert invariant is worse than an absent one, because
   `respond.OK` with an empty `findings` list is defined in the disposition
   table as "the proof's receipt" — the user is handed a receipt for a proof
   that was never attempted.
   **Blocks**: whether Phase 2 may ship before RDR 0002's schema edit lands, and
   whether a model lacking initial/terminal declarations must be refused
   (loudly, with a code) rather than passing with a vacuous green.

4. **`Critical Assumptions` A5 (`Status: Pending`, Method `MVV Test`) and
   Validation scenario 7 ("`make check` or `.github/workflows/ci.yml` invokes
   the production `intrastate lint` command over the checked-in transition
   model **or fixture corpus**").**
   The whole user outcome of this RDR rests on the gate being real — the
   `Approach` says a defective model "is rejected before it can be used by the
   resolver or accepted by CI," and Alternative 3 is rejected precisely because
   hooks "can miss the gate." But scenario 7's disjunction lets the gate be
   satisfied by linting *the fixture corpus* alone. Linting fixtures proves the
   command works; it does not gate the maintainer's actual transition model.
   A CI job that lints only fixtures is indistinguishable, from the user's
   seat, from having no gate: an illegal production graph merges green.
   A5's `If wrong` sees the risk ("the graph may be lintable locally but not
   enforced at the design-time boundary maintainers actually rely on") but the
   validation that is supposed to close A5 accepts the failing arm.
   **Blocks**: what the CI step's argument is at Phase 3 — the checked-in
   model, the fixture corpus, or both — and therefore whether merging this RDR
   actually stops an illegal graph from landing.

5. **`Problem Statement` / `Approach` — "rejected before it can be used by the
   resolver", against the disposition-table row "Model unreadable / not
   conforming to RDR 0002's schema | per RDR 0002 | refused before
   normalization; lint never runs".**
   The stated outcome is that illegal graphs are caught *before use*, but the
   RDR defines no coupling between lint and the resolver at all: the resolver
   (`flow next`, `flow resolve`, RDR 0005) is never required to have been
   linted, and nothing prevents `flow resolve` from running against a model
   with blocking findings. The only real barrier is CI, and CI's scope is
   itself unsettled (finding 4). "Before use" is therefore an aspiration in the
   Approach, not a contract anywhere in `Normative Contracts` — the closest
   clause says only that a defective model "MUST NOT be accepted for resolver
   use or CI success", which states a prohibition with no enforcer named.
   The user reads "caught at design time instead of discovered at runtime" and
   reasonably expects a local `flow resolve` on a broken model to be refused or
   at minimum warned. It will not be.
   **Blocks**: whether the resolver path gains any lint coupling (a cached
   verdict, a startup check, an explicit "unlinted model" warning), or whether
   "before use" is downgraded in the Problem Statement to "before merge".

6. **`Illustrative Code` — "Illustrative command shape only; RDR 0005 may still
   adjust flag placement: `intrastate lint --flow rdr --model
   ./path/to/rdr-transition-model.toml`" — against `Normative Contracts` "The
   authoritative CLI surface for graph acceptance MUST be the root command
   `intrastate lint`".**
   The RDR locks the command *name* as authority but leaves its *inputs*
   illustrative and cedes flag placement to a peer. For an acceptance gate, the
   input contract is the gate: whether the model is named by `--model` path, by
   `--flow` identifier resolved through RDR 0005's state binding, or by
   discovery of every checked-in model determines whether CI can be written to
   lint *all* models rather than one the author remembered to list. A gate that
   requires enumerating each model by path silently stops covering a model
   added later — the exact regression class this RDR exists to prevent.
   Prior iteration-1 raised the illustrative-invocation problem and iteration-2
   marked it closed on the grounds that the *command* is now authoritative; the
   *argument* half was not closed.
   **Blocks**: whether Phase 3's CI step is one invocation over a discovered
   model set or N hand-listed paths, and therefore whether the gate stays
   complete as models are added.

7. **`Trade-offs > Consequences` — "Reachability over-approximates the runtime,
   so a guard-infeasible path can produce an owned-before-write or dead-end
   finding the runtime would never hit; the cure is a clearer model, never a
   weaker lint" and Validation scenario 10.**
   This is a decided trade-off and the direction is defensible, but the RDR
   never sizes the cost or gives the user any handle. Every over-approximation
   false positive is a blocking exit-2 on a model that is actually correct, and
   the only prescribed remedy is to edit the model to satisfy the lint —
   "explicit write or clear on the model" per scenario 10. There is no
   suppression, no per-rule waiver, no advisory downgrade, and by design no
   escape hatch (the `Briefly Rejected` list kills the warning tier for
   withheld claims specifically, but the reasoning is applied nowhere else).
   For the maintainer this can read as the lint dictating model structure to
   work around its own imprecision. The RDR should either state that the
   expected false-positive rate on the actual RDR/kata graphs is zero or near
   it (which the `Research Findings` "small enough for a declarative table plus
   a short lint" note suggests is plausible and would settle it), or
   acknowledge the adoption risk.
   **Blocks**: whether any suppression or waiver mechanism is in scope at all,
   and whether the fixture corpus must include a realistic-scale model to
   measure how often over-approximation actually bites before this ships.

8. **`Critical Assumptions` A7 (`Status: Pending`) — "This assumption flips to
   `Verified` when RDR 0003 A18 records invariant 5 as its producer."**
   RDR 0003 is `Final [locked 2026-08-22]`. Its A18 remains `Pending` with a
   `Plan` naming both a RDR 0007 arm and a RDR 0006 arm. A7's flip condition is
   therefore an edit to a locked document, and the project rule is that RDRs
   are not amended. As written, A7 can never flip — it is a permanent Pending
   whose stated discharge path does not exist. This is a bookkeeping defect
   with a user-outcome consequence: the single-valued-state guarantee
   (invariant 5, code `graph-single-valued-state`) is presented to the user as
   a blocking invariant while its own assumption records that the guarantee's
   ownership is unresolved.
   Note A7's substance is sound — the `If wrong` correctly bounds the damage to
   "the pair needs reconciling." Only the discharge mechanics are unreachable.
   **Blocks**: whether A7 is restated as a Design Decision this RDR makes
   unilaterally (invariant 5 *is* the producer, full stop) or stays Pending;
   and whether the Finalization Gate's `Assumption Verification` paragraph
   claiming "A7 is Pending on RDR 0003 A18 recording invariant 5 as its
   model-level producer" is accurate at lock.
