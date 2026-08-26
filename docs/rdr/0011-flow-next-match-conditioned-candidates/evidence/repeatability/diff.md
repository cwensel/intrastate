Model: claude-opus-5[1m]

# Repeatability-LITE diff — cli/0011 vs run-1 (claude-sonnet-5)

Single alternate-model reconstruction. Not a consistency sample, not a
disagreement count. Every finding below stands on its own as a concrete
algorithmic contract the RDR left under-determined.

## Admissible findings

### F1 — `0011:C1`: the gate-id `not-evaluated` source set is unowned

**What run-1 rendered.** Pseudo-code line 146: `for gateID in row.DeclaredGates`
— every gate id *declared on the row*, emitted per candidate, unconditionally
on any reported disposition. Section 1 restates it as "a declared gate id not
run because `--evaluate-gates` was not passed".

**The RDR passage / the silence.** C1's `unknown` paragraph says "(absent
--evaluate-gates) every gate id MUST appear in that candidate's `unknown`
list" and later "A gate id un-run because --evaluate-gates was not passed".
The `disposition` mini-check row reads "gate id, `--evaluate-gates` not
passed → `{gate-id, not-evaluated}`". Neither names the *field or accessor*
the id set is read from, nor whether it is the row's declared gate ids, the
model-level gate declarations reachable from the row, or the gate-accessor
ids `0005:C1`'s gate-opt-in clause already enumerates for a reported
candidate. `0011:MVV` step 7 only asserts the *consequence* ("any declared
gate id lands in `unknown` as `not-evaluated` and the list is non-empty by
construction"), which holds under any of the three readings.

**Why an implementer could land elsewhere.** A row in this model can reference
gates it does not itself declare, and `0005` already ships a distinct notion
of "gate accessors for reported candidates". An implementer reading C1's
"every gate id" as model-scoped emits a superset of run-1's per-row set; one
reading it as the already-shipped gate-accessor list emits ids under a
different vocabulary. All three satisfy every oracle the RDR states, and the
three produce different `unknown` payloads on the same fixture — the payload
`D-identity` and the `(key, reason)` sort are defined over.

### F2 — `0011:C1`: the `not-evaluated` emission is not pinned per disposition

**What run-1 rendered.** Pseudo-code lines 131–147 place the gate-id loop
inside the single `case Plan, GuardUnevaluable, OwnedStateUnavailable` arm, so
`not-evaluated` entries are emitted identically on all three reported
dispositions — including `owned_state_unavailable`.

**The RDR passage / the silence.** C1 states exactly one fidelity limit, and
it is scoped to guard facts: "when the probe refuses
`owned_state_unavailable`, `gate` returns before it collects the guard
payload, so an `uncomparable` guard atom on that row is not reported; every
`absent` fact on the row still is, from the walk". The `disposition` table's
`owned key no reader established` row likewise names only the guard-atom
suppression. Gate ids are a CLI-side class with no kernel payload dependency,
so the stated limit does not reach them — but the RDR never says whether a row
that never reached `gate` at all still reports its gate ids as
`not-evaluated`.

**Why an implementer could land elsewhere.** The reason token's own
justification in `D-undecided-vocabulary` is "a declared gate id not run" —
which is *equally true* of a row the kernel refused `owned_state_unavailable`
before gating and of a plan row the caller declined to gate. An implementer
who reads `not-evaluated` as "the caller did not ask" emits it on all three
arms (run-1's reading); one who reads it as parallel to the guard-payload
limit — report only what the row's own evaluation reached — suppresses it on
`owned_state_unavailable`. Both are consistent with C1 and with `MVV` 7, whose
fixture is explicitly a fully-resolved candidate and so never exercises the
`owned_state_unavailable` arm.

### F3 — `0011:C1`: the guard-payload read's disposition guard is unfixed

**What run-1 rendered.** Pseudo-code line 141: `if result.Disposition ==
GuardUnevaluable` — the `Refusal.Undecided` projection runs *only* on the
`guard_unevaluable` arm.

**The RDR passage / the silence.** C1 says the guard reasons are "read off the
probe's `guard_unevaluable` refusal payload (`Refusal.Undecided`)", and
`0011:A13` confirms `evaluateAtoms` "return[s] the payload only when the row
verdict is `GuardUnevaluable`". That fixes *where the payload lives*. It does
not fix whether `summarize` conditions its read on the disposition or reads
`Refusal.Undecided` unconditionally on whatever `Result` `excluded` hands
back — the mini-check `authority` table says only "now also `summarize` via
`Refusal.Undecided`".

**Why an implementer could land elsewhere.** These diverge observably on the
`no_match` arm. C1 requires an excluded row not be reported, so a row refused
`no_match` carries no `unknown` list — but the kernel's `Resolve` reaches
`gate` only for rows surviving `matches`, and a *guard*-decided-false
exclusion also surfaces as `no_match` (the `disposition` table's "guard
decided false" row). An implementer reading the payload unconditionally and
then filtering by disposition, versus gating the read on
`GuardUnevaluable`, must independently re-derive which refusals carry a
payload at all; the RDR pins the source but not the caller-side guard, and
`0011:S7`'s oracle exercises only the `guard_unevaluable` row.

### F4 — `0011:C1`: `excluded`'s return type is named but its arity/ownership of the presence test is not

**What run-1 rendered.** Two incompatible placements in one document. Section 2
gives `excluded` the presence filter *and* has it "returning the kernel's own
verdict"; section 1's signature is `excluded(row, owned, observed) *resolve.Result`
— i.e. `excluded` rebuilds the view from `owned`/`observed` per row. Section 3
and the pseudo-code instead compute `view := AssembledView(owned, observed)`
once, outside the row loop, and pass presence in.

**The RDR passage / the silence.** C1 requires "The invoked reader set and the
assembled view are fixed ONCE per invocation, before the row loop". The
mini-check `authority` table assigns match-key presence to
"`flow_next.go::excluded` over `assembledView`" with call sites
"`::excluded`, `::summarize`" and the note "one test, two uses — C1 forbids a
second presence vocabulary". So the RDR fixes that there is one presence
*test* and that the *view* is built once, but it never says which function
*owns* the built view, i.e. whether `assembledView` is called once by the
caller and threaded to both consumers, or called by each consumer over the
same inputs.

**Why an implementer could land elsewhere.** `0011:A11`'s verification records
that today `runFlowNext` passes the same `owned` slice *both* to
`assembledView(owned, req.observed)` *and* to `excluded(row, owned,
req.observed)` — the shipped shape rebuilds per row. An implementer preserving
that shape satisfies "same key set" but not "fixed once before the row loop"
in the literal sense; one hoisting the map satisfies both but changes
`excluded`'s signature in a way no contract states. This is not a naming
question: it decides whether the per-key presence decision is recomputed
inside the loop, which is precisely the locus C1 says "forbids a second
presence vocabulary". Run-1 rendered both answers and did not mark either
GUESS.

### F5 — `0011:C2`: the `--all` filter's discriminant is left as "BlockMatch atoms", with no owner for the owned-key exception

**What run-1 rendered.** Pseudo-code line 139:
`unknown := [e for e in unknown if e.key not a MATCH-block atom's key]` —
filtering by *key*, over the merged list. Section 2 (helper 2, clause d) states
the opposite intent: "suppress ONLY the match-atom-derived entries from this
walk while leaving owned-key-absent entries (from the demand-set walk)
intact."

**The RDR passage / the silence.** C2 is emphatic about the *requirement*:
"The filter is scoped to the ATOM WALK's contribution ONLY. It MUST NOT
suppress a `{key, absent}` pair the OWNED-key walk independently produces for
the same key", and C3 pins a fixture whose absent-key row "MUST match on a key
it does NOT write or clear". What the RDR does not fix is the *mechanism* that
makes the two contributions separable at filter time. C1 describes both
sources as producing the same `{key, absent}` pair, deduplicated "by `{key,
reason}` pair, not by key" — which means the two contributions are
*indistinguishable in the merged list* for a matched-and-written key. The
filter therefore cannot be a predicate over the merged list at all; it must be
applied at emission, or the entries must carry a provenance the payload does
not.

**Why an implementer could land elsewhere.** The RDR states the invariant and
even names the failure ("asserted over a matched-and-written key it fails on
the very model MVV 4 runs, and forcing it to pass would delete a shipped 0005
fact") without pinning the data-flow that satisfies it. Run-1 wrote the
filter as a post-merge key predicate — the exact construction C2 forbids —
in its pseudo-code while writing the correct intent in prose. An implementer
who dedups first and filters second silently deletes the shipped 0005 owned-key
fact on every `models/rdr.toml` rule; one who filters at emission and dedups
after does not. Both readings are available because C1's dedup rule and C2's
filter scope are stated in different elements and their *order of application*
is never fixed.

### F6 — `§technical-design` (`0011:C1` / `0011:A15`): the match-owned demand term's row scope for `flow resolve` is stated twice, differently

**What run-1 rendered.** `MatchOwnedKeys(model, rows)` over
`rows := model.NonEscapeRows()` — a union over *all* non-escape rows, with the
same expression reused verbatim for both callers. Section 3 repeats it as
`matchOwnedKeys(model, rows)`. Run-1's `invokedReaders(model, outcome)`
signature retains the `outcome` parameter but its rendered body never reads it.

**The RDR passage / the silence.** C1 binds the term to both callers —
`invokedReaders(req.model, "")` and `invokedReaders(req.model, outcome)` — and
says "The demand set is a union over the outcome's rows, so the reader is
invoked even when the row `resolve` would select does not itself match on the
key". `D-identity` and the `authority` table say the set is "a property of the
MODEL (and, for `resolve`, the requested outcome), not of the mode". So the
*existing* two terms are outcome-filtered for `resolve`, and C1 asserts the
new term is too ("a union over the outcome's rows") — but the term is
introduced as "each row's MATCH-block owned keys" with no restatement of the
outcome filter at the point of definition, and the escape-row treatment of the
new term is never stated (C1 strips the escape list from the *probe*; whether
escape rows contribute to the *demand set* is a different question the
existing `guardOwnedKeys` answer is not quoted for).

**Why an implementer could land elsewhere.** Whether the match-owned term is
outcome-filtered decides the size of the breaking class C1 names and accepts.
Unfiltered, `flow resolve --outcome X` invokes readers demanded only by rows
of outcome Y, and the plan→exit-2/3 regression C1 accepts widens to models
where no row of the requested outcome match-owns anything — strictly more than
"a model with a match-only owned key" for that outcome. `0011:S8`'s oracle
runs one outcome over a purpose-built fixture and cannot separate the two.
Run-1 rendered the unfiltered form for both callers without a GUESS marker,
which is the reading that maximizes the accepted breakage.

## Discarded GUESSes (run error, not RDR silence)

Run-1 carries 7 GUESS markers. All 7 were tested against the contracts; none
is admissible, and two are outright run errors against pinned text:

1. **`excluded` returns `Result` not `bool`** (section 1, marked "RDR states
   behavior, not the exact signature"). The RDR *does* pin it, twice: the
   `authority` mini-check reads "`excluded` returns the `Result` instead of a
   `bool` (A13)", and the Existing Infrastructure Audit's first row records the
   known limit "returns `bool`, discarding the `Result`" against the decision
   "returns the `Result` so `summarize` can read the guard payload". Run error.
2. **Probe omits `Outcomes: []string{row.Outcome}`** (section 3's Probe shape
   binds only `Recognized`). `0011:A1` pins it explicitly — "`excluded` binds
   `Outcomes: []string{row.Outcome}` alongside `Recognized: row.Outcome` … C1's
   probe builder MUST preserve that pairing." Run error; without it
   `unmodeled_outcome` is not in fact unreachable, so run-1's own
   `case UnmodeledOutcome: unreachable` rests on a binding it dropped.
3. **`unknownEntry` Go field casing** — the RDR states it owns "the payload's
   INFORMATION, not its text layout" (`D-undecided-reporting-shape`) and fixes
   the JSON view; Go identifier casing is naming, inadmissible per the lens.
4. **`matchOwnedKeys` helper name** — invented name, inadmissible in itself.
5. **`readers` response field name** — a pre-existing `0005` field this RDR
   does not touch; not 0011's to fix.
6. **"other pre-existing Candidate fields carried forward unchanged"** — C3
   fixes the rename census exhaustively (five test reads, the production field,
   its JSON tag, the initializer, three appends, the `slices.Contains` guards,
   the header comment, the `Long` help, plus two prose docs). Nothing is left
   open.
7. **Section 3's blanket "exact Go struct/field names throughout"** — naming.

## Escalation assessment

**No.**

Criterion (a) does not fire. The run is not suspiciously clean: it carries 7
GUESS markers and, more tellingly, two internal self-contradictions (F4's
`excluded` arity, F5's filter placement) where the prose and the pseudo-code
render different algorithms. That is the signature of genuine
under-determination being papered over per-section, not of shared-model
overfit. Determinacy is readable from this one run.

Criterion (b) is the live question, and it comes closest on **F6**, which lands
on the cross-RDR anchor: `internal/cli/flow_exec.go::invokedReaders` is shared
with 0010, and C1's `JC1` line turns on that anchor. But the cross-RDR weight
is already discharged by argument the diff can check: `0010:A2` puts the
`decision-table` class at zero owned tags, so the added term is *empty* over
0010's class under every reading of F6 — filtered or unfiltered, escape rows in
or out. F6's divergence therefore cannot reach 0010 at all; it is confined to
the `state-machine` class 0010 does not touch, where the only peer reliance is
`0005:C1`'s narrowing clause, already named under Overrides. A second and third
run would re-sample the same silence without adding power over a peer, because
no peer is downstream of it.

The five remaining findings (F1–F5) are all `flow next` payload-internal — they
change what a candidate's `unknown` list contains, not what any peer contract
reads. `0011:C1` is load-bearing *within* 0011, but a one-run sample is
adequate to establish these are open, and they are open concretely enough to
fix by writing three sentences into C1 and C2 (gate-id source set; emission
per disposition; the order of filter-then-dedup) plus one into the demand-set
paragraph (outcome scope of the new term). Escalating to x3 would buy
frequency data on silences already localized to specific clauses, which is not
what x3 is for.

Recommend resolving F1–F6 in place rather than escalating.
