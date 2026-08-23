Model: claude-opus-5

# Critique — RDR 0006 Graph Lint Authority And Guarantees (iter-2)

## 6. Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | §Technical Design, "**Source state is a tag-set, not a name.**" — "a row belongs to the group of every reachable node its match pattern over owned tags is satisfiable in. Group identity is therefore `(owned-state node, recognized outcome)`" | This IS a second grouping predicate, contradicting the RDR's own `Coverage, overlap, and withholding MUST be decided per scoped row group as RDR 0003 defines it … This RDR MUST NOT define a second grouping predicate`. RDR 0003 (Final, locked 2026-08-22) defines the group syntactically as rows sharing one *source state* and one recognized outcome; RDR 0006 rebinds "source state" to a merged fixpoint lattice node, making group membership a function of the abstract interpretation rather than of the authored rows. A2's claim to "close RDR 0003 A10" is therefore false: A10 asks whether 0006 reads the *same* division of labour, and it does not. | Same model, two lint implementations (or one implementation before/after a fixpoint bug fix), two different coverage verdicts. A model that passed yesterday fails today with `graph-coverage-gap` on a group the author never authored as a group. | §1, §2, premortem, AT-1 |
| C-2 | §Technical Design, "invariants 3 and 4 run once per reachable group … a row participating in several reachable nodes is checked once per node" | Group count is `|reachable owned-state nodes| × |outcomes|`, not `|authored selection contexts|`. Because the join rule unions per-tag value sets, a model with *k* owned tags of domain size *d* has up to `(2^d)^k` nodes. Every one of those is a group whose scoped product must be materialized against the published `graph-product-too-large` bound. The RDR publishes a bound on the *product per group* and never bounds *the number of groups*. | `intrastate lint` on a real RDR/kata model hangs or OOMs with no output at all. Worst case it emits thousands of `graph-overlap` and `graph-coverage-gap` findings, one per synthetic node, for a model with a dozen authored rules. There is no waiver mechanism (§Trade-offs, "There is no suppression or waiver mechanism"), so the user cannot ship. | §1, §2, premortem, AT-2 |
| C-3 | §Load-Bearing Decisions, "Reachability relation" — "an edge is a normalized non-escape row whose match pattern over owned tags is satisfiable in the source node, producing the node with that row's writes applied" | Escape rows are excluded as edges, but RDR 0002 `Normative Contracts` fix that an escape rule "MUST NOT contain a write block or clear list" and normalizes "to a row with both empty" — so an escape rescue is a legal runtime transition that produces **no** next state. Lint's graph therefore has an entire class of runtime-reachable outcome with no representation. Invariant 2 (dead end) will accuse the node an escape rescue lands on, and invariant 6 will accuse rows downstream of it, because lint believes that path does not exist. | A model whose escape rows are the designed recovery path is rejected with `graph-dead-end` / `graph-owned-before-write` on the recovery arm. The RDR's own cure ("an explicit write, clear, or terminal declaration") is *forbidden* on escape rules by RDR 0002, so the user has no legal edit that makes lint green. Hard deadlock. | §1, §3, premortem, AT-3 |
| C-4 | A4 Evidence, "`internal/cli/respond/respond.go::Fail` emits the structured envelope in text/json modes" (Status: Verified); §Load-Bearing Decisions, "Wire / byte format" — "`error.findings` on the failure envelope"; §Technical Design, "`{\"type\":\"failed\",\"error\":{…,\"findings\":[…]}}`" | The `{"type":"failed","error":{…}}` envelope **does not exist in shipped code**. `respond.Fail` in JSON mode calls `clierr.EmitJSON`, which does `json.Marshal(e)` on the `*CLIError` and prints it — producing `{"code":…,"message":…}` at top level, with no `type` discriminator and no `error` nesting. The shape the RDR quotes appears only in `respond.go`'s *package doc comment*, which the code contradicts. A4 is stamped Verified against a symbol whose behavior it misread. | The JSON contract the RDR specifies cannot be produced without changing `EmitJSON` — a change that breaks every existing consumer of the failure line (including the `version` verb's error path). Either implementation silently ships `findings` at the wrong depth, or CI's `jq '.error.findings'` returns null on every failure and the gate passes a broken model. | §1, premortem, AT-4 |
| C-5 | §Technical Design, "**Type ownership**: … The `Finding` struct is therefore defined **in `clierr`** as a subsystem-agnostic record"; §Normative Contracts, "A finding attributed to one guard atom MUST carry that atom's `Key`, `Operator`, `Literal`, and `Block`; a finding scoped to an escape population MUST carry the failure class." | Self-contradictory. A record carrying `Block` (an RDR 0002 authoring concept: `match`/`all`/`unless`) and a resolver *failure class* (`no_match`/`ambiguous_match`, RDR 0001's `RefusalKind`) is not subsystem-agnostic — it is graph-lint and guard vocabulary wearing a generic name. Either `clierr` grows a dependency on the atom/class types (the cycle the clause exists to prevent) or the fields degrade to bare strings and the "MUST carry" is unenforceable by the type system. | The atom field ships as `map[string]string` or four loose strings, MVV Scenario 8's atom assertion is written against stringly-typed data, and RDR 0003 A20 — which is waiting on exactly this carrier — closes against a field that cannot be validated. First schema change lands within weeks. | §1, §2, AT-5 |
| C-6 | A5 Gate target, "running `make build` then the built `./bin/intrastate lint --as=json` over the **checked-in transition model**"; Validation scenario 7 | **There is no checked-in transition model.** The only `.toml` files in the repo are `.roborev.toml`, `.kata.toml`, and three spike fixtures under `docs/rdr/*/evidence/spikes/`. A5 explicitly rules the fixture corpus out as the subject ("not the fixture corpus"). The RDR's one Pending-by-MVV assumption names an oracle that points at nothing, and never books who authors the model. | Scenario 7 cannot be written. Implementation either invents a transition model nobody specified (and lints an artifact with no owner), or points the gate at a spike fixture and A5's disjunction-resolution is undone. The gate ships green over a model that is not the one maintainers edit. | §1, §3, premortem, AT-6 |
| C-7 | A5 Gate target, "a **new `lint` job in `.github/workflows/ci.yml`**" | `.github/workflows/ci.yml` **already has a job named `lint`** — the golangci-lint/gofmt job. A5's oracle ("the job exists in that file, names that command") collides with a shipped job of the same name. The RDR treats the name as free. | Either the Go lint job is clobbered (formatting and static-analysis regressions stop being caught) or the graph-lint gate lands under an unnamed job and Scenario 7's mechanical check — which asserts on the job's identity in that file — cannot be evaluated. | §1, AT-7 |
| C-8 | Invariant 2 (Dead end), "a node *satisfies* it when **every** value in each of the node's per-tag value sets meets it"; Scenario 10's paired dead-end case | Combined with the union join rule, terminal satisfaction is anti-monotone in merging: the more paths reach a node, the *less* likely it satisfies any terminal. Every terminal node reachable by two distinct paths that write different values becomes a merged node satisfying no terminal, hence `graph-dead-end`. This is not a rare false positive — it is the *normal* shape of a workflow with more than one route to completion. | Every real flow model fails `graph-dead-end` on its own terminal state. The stated cure ("explicit … terminal declaration") does not help: the merged node's set is `{done, review}` and no single terminal predicate over declared values covers it. Users are told to fix the model and cannot. | §1, §3, premortem, AT-8 |
| C-9 | Invariant 5, "A write to a single-valued tag **replaces** its prior value (no preceding `clear` is required), so no path can accumulate a second value and there is no path-accumulation case left to check." | This asserts a write semantics no cited peer owns. RDR 0002's normative clauses cover write blocks, explicit `clear`, and the rule that "Absence from both the write block and the clear list MUST NOT imply deletion" — nothing states replace-vs-accumulate for single-valued tags. RDR 0003 owns the *marker*, not the update rule. This RDR invents the semantics inline and uses it to *delete* a check. | If any peer implements accumulate (or set-valued union) semantics, invariant 5's entire justification for being syntactic-per-row collapses, and the multi-value state it was minted to catch reaches runtime as a `TagSet` whose `assemble` last-write-wins silently discards one value. Lint reported green. | §1, premortem, AT-9 |
| C-10 | A7 / Invariant 5's sibling, `graph-always-present-owned` — "a key declared always-present must be held in every reachable owned-state node" | The initial owned state (A6) is a *table of tag assignments*, necessarily a proper subset of all declared always-present owned keys unless the author declares every one of them at the root. So the root node itself violates the invariant for every always-present key not in the initial table. The RDR never exempts the root, and the reachability relation makes the root a reachable node by construction. | Every model fails `graph-always-present-owned` on its first run, pointing at the root. The only fix is to inflate the `initial` declaration with every always-present key — which makes the initial-state declaration a duplicate of the tag declarations and defeats its purpose. | §1, premortem, AT-10 |
| C-11 | §Technical Design, "Lint MUST NOT key groups on the syntactic match pattern: two rows whose patterns differ in spelling but denote the same owned states are in the same group and MUST be overlap-checked." | Directly contradicts RDR 0003's locked participation clause: "**Match keys are not product dimensions.** A row's match pattern selects which group the row belongs to — the selection context this RDR groups by." RDR 0003 groups *by the match pattern*; RDR 0006 forbids exactly that. Since RDR 0003 is Final and locked, this is not a Draft-side reconciliation — it is a contradiction between a locked and an unlocked document, booked nowhere in this RDR's Contradiction Check (which claims "No CONTRADICTION row"). | The overlap check runs over a row population RDR 0003's derivation was never defined for. Two rows with disjoint match patterns that happen to be satisfiable in the same merged node are reported as `graph-overlap` — a pair the runtime never evaluates together, since `Resolve` filters on `view.matches(row.Match)` before the guard gate. A false positive on a model the kernel resolves correctly. | §1, §3, premortem, AT-11 |
| C-12 | §Technical Design, "for escape rows — row kind and the declared list of failure classes the row rescues"; invariant 4's per-class coverage union; Scenario 19 | The kernel does not compute coverage per rescuable class — `escapeOrRefuse` is consulted only *after* the ordinary population produces zero or many matches, and it rescues on `row.rescues(r.Kind)` for the single kind that actually occurred. Lint's `(group × declared rescuable class)` product is a design-time artifact with no runtime counterpart: there is no runtime state in which the `ambiguous_match` arm's "coverage" is separately consumed. Requiring both arms closed demands escape rows the runtime will never reach. | `graph-coverage-gap` on the `ambiguous_match` arm of a group whose ordinary rows are provably exhaustive — a group that *cannot* produce `ambiguous_match` because overlap is already blocking. Users add dead escape rows to silence lint, which then trip `graph-unreachable-rule`. | §1, premortem, AT-12 |
| C-13 | §Approach, "a locally-edited illegal model can still be resolved against on a working tree. That is accepted" | The RDR's Problem Statement promises "illegal or incomplete transition graphs caught at design time instead of discovered at runtime," and the Normative Contract says a blocking model "MUST NOT be accepted for resolver use." Nothing enforces the second sentence: the kernel takes a `Table` value from whoever built it, and there is no lint verdict on that path. The contract is unenforceable as written — it is a merge-boundary policy dressed as a runtime guarantee. | A maintainer edits the model, runs `intrastate flow next`, gets a confident wrong transition, and believes lint would have stopped it. The failure the RDR exists to prevent happens in exactly the workflow (local iteration) where it hurts most. | §1, §3, premortem, AT-13 |
| C-14 | §Technical Design, "the same blocking code an unprovable dimension takes, with a widened trigger" (`graph-unprovable-coverage`) | One code, three semantically distinct triggers: a non-finite dimension, a narrowing-withheld claim, and (per RDR 0003 A21) an atom over a tag not declared single-valued. Scenario 3 concedes the arms carry different payloads ("atom fields are required only when … the withheld-claim arm"). A machine consumer branching on `code` cannot tell "your tag needs a domain declaration" from "your model is fine but lint declines to promise". | CI failure output says `graph-unprovable-coverage` and the maintainer has no way to know whether to add a declaration, add a marker, or accept that the group is unprovable by design. Diagnosis requires reading the RDR. | §1, §2, AT-14 |
| C-15 | §Trade-offs, "**There is no suppression or waiver mechanism** — no inline ignore comment, no allowlist … a false positive is expected to be rare" | The rate is asserted, never measured. Given C-2 (node explosion), C-8 (merge-anti-monotone terminals), C-10 (root always-present), and C-11 (spelling-blind grouping), false positives are structural, not rare. The RDR pre-commits to having no escape hatch and names a successor RDR as the only response. | On the first real model the gate is red for reasons the author cannot fix, and the only available action is to delete the CI job — which is precisely the "gate bypassed after implementation" outcome A5's `If wrong` names. The authority is lost the week it ships. | §2, §3, premortem, AT-15 |
| C-16 | §Normative Contracts, "Lint success MUST carry non-blocking findings under the existing `respond.Success.Data` payload as `data.findings`. On both surfaces the key MUST be emitted even when the list is empty" | `respond.Success.Data` is declared `Data any \`json:"data,omitempty"\``. The RDR's Testing Strategy claims "every existing optional `CLIError` field is `omitempty`, so appending the typed `findings` field … is wire-compatible" — but never notices that `Data`'s own `omitempty` sits between the verb and the `findings` key. Emitting `[]` inside `Data` works only if the value assigned to `Data` is a non-empty-struct — a fact the RDR does not state and Scenario 16 does not distinguish from the failure case. | Scenario 16's golden-file assertion passes or fails on an implementation detail (whether `Data` holds a struct or a map, and whether that map is nil) that nobody specified. The "empty list is the proof's receipt" guarantee silently evaporates for a verb that assigns `Data` a nil-valued interface. | §1, AT-16 |
| C-17 | A2 Evidence, "**A17** by invariant 3's two-population overlap …, **A18** by invariant 5 …, **A19** by `graph-coverage-closed-by-escape`, and **A20**'s lint half by the atom field on the finding contract" — "close by citation" | RDR 0003 is `Final [locked 2026-08-22]`. A Draft peer cannot close a locked peer's assumption records by asserting it. A7 concedes this for A18 ("that edit is a route-back on a locked peer"), but A2 asserts A17/A19/A20 close "by citation" without the same concession, and the Finalization Gate's Assumption Verification repeats the claim. | RDR 0003 still reads Pending on four records at implementation time. Whoever implements 0003's guard engine and whoever implements 0006's lint each believe the other side closed the question, and the seam ships with two readings of overlap population and verdict shape. | §2, premortem, AT-17 |
| C-18 | §Technical Design, invariant 6 — "the declared initial owned state counts as a write" | RDR 0007 (Final) narrows `Row.RequiresOwned` to *post-guard write dependencies*, explicitly excluding guard-read keys ("guard-input decidability is not RequiresOwned's job"). Invariant 6 checks a *different* set — match keys plus guard atoms — for which the normalized row carries no derived field. RDR 0002's normative clause fixes `RequiresOwned` as "the sorted, duplicate-free set of tag keys named by the rule's write block and clear list" and "does not add guard-read keys to it." Lint must recompute the read-set from raw atoms, which the RDR's minimum input contract never lists as an obligation. | Implementation reads `RequiresOwned` (the field that exists) instead of the atom read-set (the field that doesn't), and invariant 6 checks writes-before-writes instead of reads-before-writes. `graph-owned-before-write` never fires on the defect it was minted for; `owned_state_unavailable` still surfaces at runtime. Silent false green. | §1, premortem, AT-18 |

### Overlap with prior lenses

Read only after this critique was written, to avoid re-raising resolved findings.
Disclosed rather than deleted, because in each case the *resolution* the earlier
lens produced is the thing this critique attacks.

- **C-1 / C-11** sharpen 3amigo **S1** ("source state is undefined; three
  incompatible readings are live"). The draft resolved S1 by *choosing* the
  abstract-node reading. That resolution is the defect: it contradicts RDR 0003's
  locked participation clause, which S1 flagged but did not adjudicate. New: the
  direct quotation of the conflicting locked clause and the false A10 closure.
- **C-2** is new. 3amigo **S3** raised the join rule as unspecified; the draft
  answered it with union-on-merge and a termination proof. Nobody asked what the
  node *count* is, or that group count now scales with it.
- **C-8** sharpens **S9** and **S14** (invariant 2's satisfaction relation; the
  dead-end half has no scenario). The draft added both — the "every value" rule
  and Scenario 10's paired case. New: that the two, combined with the join rule,
  are jointly unsatisfiable for any two-path completion.
- **C-4** and **C-16** are new. 3amigo **H7** raised the two-shapes problem and
  the `clierr` cycle; the draft resolved it by fixing the depths. Nobody checked
  the depths against `EmitJSON`, or `Success.Data`'s own `omitempty`.
- **C-6 / C-7** sharpen **H5** and iter-1 **C3**. The draft resolved the
  disjunction by naming one target. New: that target's artifact does not exist,
  its job name is taken, and Phase 3 reintroduces the disjunction.
- **C-3, C-5, C-9, C-10, C-12, C-14, C-17, C-18** are new.
- **C-13** restates 3amigo **S15** and **C-15** restates **S16**, both unresolved
  in the current draft — carried so the ledger is complete, not as new discovery.

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The grouping predicate forks from RDR 0003, and nobody notices until two engines disagree

**Root cause.** This RDR simultaneously forbids itself from defining a grouping predicate and then defines one.

**The passage that enabled it.** §Normative Contracts:

> Coverage, overlap, and withholding MUST be decided per scoped row group as RDR 0003 defines it … This RDR MUST NOT define a second grouping predicate; it supplies only which selection contexts are reachable.

And, three hundred lines earlier, §Technical Design:

> **Source state is a tag-set, not a name.** … a row belongs to the group of every reachable node its match pattern over owned tags is satisfiable in. Group identity is therefore `(owned-state node, recognized outcome)` … Lint MUST NOT key groups on the syntactic match pattern: two rows whose patterns differ in spelling but denote the same owned states are in the same group and MUST be overlap-checked.

RDR 0003 — `Final [locked 2026-08-22]` — says the opposite in its locked participation clause:

> **Match keys are not product dimensions.** A row's match pattern selects which group the row belongs to — the selection context this RDR groups by.

RDR 0003 groups *by the match pattern*. RDR 0006 *forbids* keying on the match pattern. These are not two phrasings of one idea; they are two different functions from rows to groups. Under 0003, grouping is syntactic, computable from the normalized row alone, and stable. Under 0006, grouping is semantic, computable only after a fixpoint over the owned-state lattice has converged, and every change to any unrelated row's writes can silently re-partition every group in the model.

The RDR knows this is contested — RDR 0003's A10 exists precisely to book "does RDR 0006 read the same row-group division of labour?" — and A2 answers it by *asserting* closure:

> This RDR adopts both by citation (`Technical Design`, `Normative Contracts`) and supplies only the reachability of selection contexts (`Load-Bearing Decisions`), which closes RDR 0003 A10 and A12.

It does not close A10. It answers "no" and records "yes."

**Symptom the user sees.** Two lint runs over the same model, from two implementations (or one implementation before and after a fixpoint fix), produce different `graph-coverage-gap` and `graph-overlap` sets. A group the author never authored — two rules with unrelated match patterns that happen to be jointly satisfiable in a merged node — is reported as overlapping. The user reads the finding, opens both rules, sees match patterns that share no key, and concludes the tool is broken. They are right: `internal/resolve/resolve.go::Resolve` filters candidates on `view.matches(row.Match)` before the guard gate, so those two rows are never in the same candidate set at runtime.

### 1.2 The reachability fixpoint explodes, and the published bound does not bound it

**Root cause.** The RDR publishes a bound on the scoped *product* per group and never bounds the number of *nodes*, while making the number of groups a function of node count.

**The passage that enabled it.** §Load-Bearing Decisions:

> **Join rule and termination** — the traversal is a **fixpoint over merged nodes** … two edges reaching the same successor produce one node whose per-tag value sets are the union of theirs … The lattice is finite — finitely-declared tags range over their declared domain's subsets, a tag with no finite domain over `{held, absent}` — so the fixpoint terminates.

"Finite" is doing heroic work. A node is a point in `∏ₜ (2^domain(t) ∪ {absent})`. RDR 0003's own worked desk trace (its step 3) computes a *single group's* product at 8,192 assignments for three tags. The node lattice for the same three tags is the same order of magnitude — and then §Technical Design multiplies it:

> invariants 3 and 4 run once per reachable group … a row participating in several reachable nodes is checked once per node

So the work is `|reachable nodes| × |outcomes| × (per-group product proof)`. The RDR's only cardinality control is:

> Lint MUST publish the model-independent product bound above which it declines to prove coverage, and MUST emit `graph-product-too-large` for a group whose declared finite product exceeds it

which fires *inside* a group, after the node set is already enumerated. Termination is proved; tractability is never discussed. §Performance Expectations waves it off:

> No throughput target is load-bearing for this RDR. Lint is a single-invocation static check … If fixture runtime becomes material, implementation should profile the invariant engine before changing the contract.

Fixture runtime will not be material. Fixtures are hand-authored to be small. The *real* RDR/kata transition model is what explodes, and by then the contract is locked.

**Symptom the user sees.** `intrastate lint --flow rdr` produces no output and no exit for minutes, then the CI runner kills it. Or it completes and emits four hundred findings, one per synthetic merged node, for a model with fourteen rules. There is no `--max-nodes`, no waiver, no allowlist — §Trade-offs forecloses all three — so the only remediation is deleting the CI job.

### 1.3 The JSON failure envelope the RDR specifies does not exist, and A4 says it was verified

**Root cause.** A4 was verified against a package doc comment rather than against the code that runs.

**The passage that enabled it.** A4, `Status: Verified`, `Method: Source Search`:

> `internal/cli/respond/respond.go::Fail` emits the structured envelope in text/json modes … Graph lint needs new stable codes, not a new envelope or direct output path.

And §Load-Bearing Decisions, Wire / byte format:

> `error.findings` on the failure envelope (a new typed `clierr` field)

with §Technical Design giving the literal shape:

> `{"type":"failed","error":{…,"findings":[…]}}`

The shipped code does not produce that. `respond.Fail` in JSON mode is:

```go
case ModeJSON:
    clierr.EmitJSON(cmd.OutOrStdout(), ce)
```

and `clierr.EmitJSON` is `json.Marshal(e)` on the `*CLIError` followed by `Fprintln`. The emitted line is `{"code":"…","message":"…"}` — top-level, no `type` discriminator, no `error` wrapper. The `{"type":"failed", ...}` form appears exactly once in the repository, in `respond.go`'s package doc comment, which the function below it contradicts. The RDR quoted the comment.

This matters more than a typo, because the RDR builds an argument on the two envelopes having different shapes:

> The two envelopes carry findings at **different depths, deliberately**, because the shipped types differ: failure adds a sibling `findings` key on the error envelope … while success carries them inside the existing `Data any` payload

The asymmetry is real but the *depths* are wrong: today failure is at depth 1 (`.findings`) and success at depth 2 (`.data.findings`). The RDR says failure is at depth 2 (`.error.findings`). Whoever writes Scenario 7's CI assertion picks one.

**Symptom the user sees.** CI's assertion reads `.error.findings` (as A5 requires: "asserts on the JSON `code` field"). Against the shipped emitter that path is `null`, the `jq` test passes vacuously, and an illegal model lands on `main` with a green check. A5's own `If wrong` describes this outcome and the RDR does not see that A4 already caused it.

---

## 2. The one section rewritten within 6 weeks of shipping

**§Technical Design's "Source state is a tag-set, not a name" paragraph, together with the Reachability relation and Join rule bullets in §Load-Bearing Decisions.**

They are one mechanism split across two sections, and they will be rewritten together because every other problem in this RDR routes through them:

- They are the source of the RDR 0003 contradiction (C-1, C-11) — a locked peer's grouping predicate versus this one.
- They make group count unbounded (C-2), so the published product bound does not bound the work.
- They make terminal satisfaction anti-monotone in merging (C-8), so multi-path workflows fail `graph-dead-end` by construction.
- They exclude escape rows from the edge relation (C-3), so the rescue path is invisible to a lint whose stated remedy for the resulting false positive is an edit RDR 0002 forbids.
- They make `graph-always-present-owned` fire on the root (C-10).

The rewrite will not be a tweak. The realistic landing spot is the one the RDR already rejected and pre-emptively named:

> If the rate turns out material in practice, the response is a successor RDR on guard-aware pruning, not a suppression flag.

That successor is not the fix. Guard-aware pruning is orthogonal: the false positives above come from the *join*, not from unpruned guards. The actual rewrite is one of two shapes — either grouping reverts to RDR 0003's syntactic definition and the fixpoint is demoted to a *reachability filter over syntactic groups* (which is what RDR 0003 asked for: "supplies only which selection contexts are reachable"), or the join rule becomes path-sensitive with a depth bound and the "no false green" guarantee is replaced with a bounded one. Either way the paragraph that begins "Source state is a tag-set, not a name" does not survive.

Secondary candidate, for the same six weeks: the **§Technical Design type-ownership paragraph** (C-5). "Subsystem-agnostic record in `clierr`" plus "MUST carry that atom's `Key`, `Operator`, `Literal`, and `Block`" plus "MUST carry the failure class" cannot all hold. The first implementer resolves it by stringly-typing the atom, and the first consumer who needs to programmatically inspect the atom rewrites the clause.

---

## 3. The one assumption that will not survive first contact with a real user

**A5**, and specifically its Gate target:

> a **new `lint` job in `.github/workflows/ci.yml`**, running `make build` then the built `./bin/intrastate lint --as=json` over the **checked-in transition model** — not the fixture corpus, which tests the engine rather than the shipped model.

Three independent ways this fails on contact, all verifiable today:

1. **There is no checked-in transition model.** The repository's entire `.toml` inventory is `.roborev.toml`, `.kata.toml`, and three spike fixtures under `docs/rdr/*/evidence/spikes/`. A5 rules the fixture corpus out by name. So the gate's subject does not exist, and no section of this RDR books who authors it, where it lives, or which RDR owns it. §Implementation Plan Phase 3 restores the disjunction A5 claimed to resolve — "over the checked-in transition model or fixture corpus" — which means the RDR contradicts its own resolution within its own body. (C-6)

2. **A job named `lint` already exists in that file.** It runs `make fmt-check` and `golangci-lint-action@v9`. A5 treats the name as unclaimed. Scenario 7's oracle — "the job exists in that file, names that command" — is satisfiable by *renaming the Go lint job*, which is not what anyone means. (C-7)

3. **`make check` has no `build` edge, and A5 knows it** — it books the cost ("which requires adding the `build` edge `check` currently lacks") but Phase 3 leaves the wiring as "add … once the checked-in transition model or fixture corpus exists." The prerequisite is a checkbox nobody owns.

The user-visible consequence is the one A5's own `If wrong` names and then treats as hypothetical: *"If that gate is bypassed after implementation, the graph may be lintable locally but not enforced at the design-time boundary maintainers actually rely on."* It will not be *bypassed* — it will simply never be wired, because the artifact it is supposed to lint does not exist and no one is assigned to create it. The command ships, the fixtures pass, the RDR closes, and the authority the entire document is about never materializes.

Runner-up: **A6**, which is Pending on RDR 0002 adding `initial` and `terminal` declarations. Its `Behavior before it lands` clause —

> not "vacuous" — a model declaring no initial owned state is rejected with a blocking finding (disposition table), so the unlanded schema cannot be mistaken for a clean model

— means that until RDR 0002 lands the schema change, `intrastate lint` **rejects every model in existence**, including the fixture corpus, with `graph-dangling-edge`. The RDR presents this as safety. It is a hard dependency inversion: 0006 cannot ship a green run until 0002 ships a schema clause that 0002 has not agreed to, and 0002 is a `foundational` document consumed by five other RDRs.

---

## 4. Premortem

*Written as if the failure already happened.*

We shipped `intrastate lint` eleven weeks after locking RDR 0006. It is now disabled in CI and two maintainers have stopped editing the transition model. Here is how it went.

**Week 1–3 — the model that wasn't there.** Phase 3 stalled immediately. Scenario 7 required a "checked-in transition model" and there wasn't one; the three `.toml` files in the repo were spike fixtures under `docs/rdr/*/evidence/spikes/`, which A5 explicitly excluded. We promoted `rdr-fixture.toml` to `models/rdr.toml`, told ourselves it was the shipped model, and wired the gate. Nobody owned that file. It had been authored to exercise RDR 0002's *parser*, not to describe a legal RDR flow, so its first lint run was never going to be clean — but we didn't know that yet, because it couldn't run at all: RDR 0002's `initial`/`terminal` schema clause (A6) had not landed, and every model without a declared root takes `graph-dangling-edge` by the disposition table. We hand-added an `[initial]` block that RDR 0002's parser rejected as an unknown layout key, patched the parser locally, and shipped with a `TODO`.

**Week 4 — the CI job name.** We added `jobs.lint` to `.github/workflows/ci.yml` per A5, and `golangci-lint` stopped running. Nobody noticed for six days, during which two `gofmt` regressions and one `errcheck` violation landed on `main`. We renamed ours to `graph-lint`, which made Scenario 7's mechanical oracle — "a `lint` job in `.github/workflows/ci.yml`" — false. We changed the scenario. That was the first time we edited the RDR to match the implementation.

**Week 5 — `buildOwnedStateGraph` never returned.** `buildOwnedStateGraph` is the fixpoint from §Load-Bearing Decisions: worklist over merged nodes, per-tag value sets unioned on join. On the four-tag fixture it converged in 40ms. On the promoted model — nine owned tags, two of them enums of five values, one bounded int `{0..3}` — it produced 61,000 reachable nodes before the runner's 6-minute step timeout. We had proved termination and never bounded the lattice. `graph-product-too-large` did not help; it fires per group, *after* `enumerateGroups` has already materialized the node set. We added an undocumented `maxNodes = 5000` constant and emitted a bare `graph-lint-failed` when it tripped, which is a code the taxonomy does not contain.

**Week 6 — the first real run, and 1,431 findings.** With the cap in place the run completed. It emitted:

- 812 `graph-dead-end`. Almost all on the *terminal* nodes. `nodeSatisfiesTerminal` implements invariant 2's "every value in each of the node's per-tag value sets meets it," and the join rule had unioned `status: {done}` from the archive path with `status: {review}` from the rework path into one merged node holding `{done, review}`. That node satisfies no terminal, so it needs an outgoing row, so it's a dead end. The RDR anticipated this in Scenario 10 and called it "an accepted false positive whose cure is an explicit write, clear, or terminal declaration." There is no such declaration. You cannot declare a terminal that means "done-or-review" without also making it satisfiable at a node where the flow genuinely has not finished, which invariant 2's own justification forbids ("accepting it would let a merged node close on a path that has not actually terminated"). The over-approximation and the strictness rule are jointly unsatisfiable for any model with two routes to completion.

- 219 `graph-always-present-owned`, every one of them naming the root. `checkAlwaysPresentOwned` walks reachable nodes for keys declared always-present; the `[initial]` table declares four assignments and the model declares eleven always-present owned keys. The root violates for seven of them. The RDR never exempts the root and the reachability relation makes it a node.

- 287 `graph-overlap` between rows whose match patterns share no key. `groupRows` implements "Lint MUST NOT key groups on the syntactic match pattern," so `prelock-flapping-cap` (match `status=Draft`) and `terminal-archive` (match `phase=closed`) landed in the same group at a merged node where both patterns are satisfiable. At runtime `Resolve` never puts them in one candidate set — `view.matches(row.Match)` filters on the concrete view, and no concrete view holds both `status=Draft` and `phase=closed`. Lint accused a pair the kernel cannot produce.

- 71 `graph-owned-before-write` on the recovery arm. `buildOwnedStateGraph` takes edges from non-escape rows only. Our recovery path is an escape row rescuing `no_match`. RDR 0002 forbids escape rules from carrying a write block or clear list, so the escape row has no writes and lint has no edge for it — the node the recovery lands on is unreachable in lint's graph, and every row downstream of it reads owned tags lint believes were never written. The suggested cure ("an explicit write") is a load failure in the parser.

**Week 7 — the maintainer journey that ended it.** A maintainer added one rule to the model: a new `reconcile → propose` rewind edge with a guard on `rewind_target`. `intrastate lint --flow rdr --as=json` returned 1,447 findings — sixteen more than before, none of them about her rule. She could not find her own change in the output. Findings *are* ordered deterministically (`model id, invariant code, source rule/context id, fingerprint`), which sorts by *invariant* first, so her rule's finding sat between two hundred `graph-dead-end`s and three hundred `graph-overlap`s. There is no `--only`, no `--since`, no severity filter, and §Trade-offs forecloses suppression: *"There is no suppression or waiver mechanism — no inline ignore comment, no allowlist."* She reverted her rule and filed the change as a comment in the RDR instead. That is the failure: the design-time gate made design-time editing harder than not editing.

**Week 8 — CI goes green by going away.** We disabled the `graph-lint` job. Nobody objected. The command still exists; `make check` does not call it (the `build` edge was never added); the model file drifted from the flow it describes within two sprints because nothing checked it. §Approach had already conceded the working-tree half — *"a locally-edited illegal model can still be resolved against on a working tree. That is accepted"* — so with CI off, nothing enforces anything, and the Normative Contract *"A model with any blocking lint finding MUST NOT be accepted for resolver use"* became a sentence with no enforcement site anywhere in the system.

**Week 9–11 — the seam split.** RDR 0003's implementation landed. Its `exhaustive.go` groups rows by match pattern, per its locked participation clause. Our `graphlint/group.go` groups by merged node, per §Technical Design. Both cite each other as the authority. `TestCoverage_SharedGroup` passes in 0003's package and fails in ours against the same fixture, and the fix required deciding which locked document was wrong. RDR 0003 A10 — the record that existed *specifically* to catch this — reads Pending in `0003`, and Verified-by-citation in `0006` A2. Nobody looked at both.

**What the postmortem said.** Not "the reachability abstraction was wrong." It said: *we specified an abstract interpretation in a design document, proved the one property that was easy to prove (termination, no false greens), and never once asked what its output looks like on a model with two paths to done.* Every failing invariant above — dead end, always-present, overlap, owned-before-write — is a consequence of the join rule, and the RDR's only sentence about the join rule's cost is *"a false positive is expected to be rare."* That sentence was the whole risk analysis.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time tests: each is answerable by reading the RDR, its peers, and the repository as they stood at draft. Each names the finding it catches.

### AT-1 — Grouping predicate agreement with the locked peer (catches C-1)

```gherkin
Scenario: RDR 0006's grouping predicate is the one RDR 0003 defines
  Given RDR 0003 is Final and its participation clause states
        "A row's match pattern selects which group the row belongs to"
  When  RDR 0006's Technical Design is read for how it assigns rows to groups
  Then  it MUST NOT contain a clause forbidding grouping by the match pattern
  And   two rows with disjoint match patterns MUST be in different groups under
        both documents for the same normalized model
  And   if the two documents disagree, RDR 0006 MUST carry a CONTRADICTION row
        in its Finalization Gate naming RDR 0003's clause by quotation
  And   RDR 0003 A10 MUST NOT be recorded as closed by RDR 0006 while RDR 0003
        reads Pending on it
```

*This fails on the draft.* §Technical Design says "Lint MUST NOT key groups on the syntactic match pattern"; RDR 0003's locked clause says the match pattern *is* the grouping key; and §Contradiction Check asserts "No CONTRADICTION row."

### AT-2 — Node-count bound (catches C-2)

```gherkin
Scenario: the reachability lattice is bounded, not merely finite
  Given the owned-state lattice is the product over owned tags of
        (subsets of the declared domain) ∪ {absent}
  When  a reviewer computes the worst-case node count for RDR 0003's own desk-
        trace fixture (profile: 4 values, prelock_iterations: {0..3}, cluster_eligible: bool)
  Then  RDR 0006 MUST state a model-independent bound on reachable node count
  And   MUST name the finding code emitted when that bound is exceeded
  And   MUST state the total work bound as (nodes × outcomes × per-group product),
        not the per-group product alone
```

*Fails on the draft.* `graph-product-too-large` bounds a group's product; nothing bounds the group count, and §Performance Expectations declines to set any target.

### AT-3 — Escape rows in the edge relation (catches C-3)

```gherkin
Scenario: the escape rescue path is representable in lint's graph
  Given RDR 0002 states an escape rule "MUST NOT contain a write block or clear
        list, even an empty one" and normalizes "to a row with both empty"
  And   RDR 0006's reachability relation defines an edge as "a normalized
        non-escape row"
  When  a model's designed recovery path is an escape row rescuing no_match
  Then  RDR 0006 MUST state which node the rescue transitions to
  And   MUST state whether nodes downstream of a rescue are reachable
  And   the cure it offers for a false positive on that path ("an explicit write,
        clear, or terminal declaration") MUST be legal under RDR 0002
```

*Fails on the draft.* No answer to the first two; the offered cure is a load failure under RDR 0002's escape clause.

### AT-4 — Failure envelope shape against shipped code (catches C-4)

```
Steps:
1. Read internal/cli/respond/respond.go::Fail. Note the JSON branch calls
   clierr.EmitJSON(out, ce).
2. Read internal/cli/clierr/clierr.go::EmitJSON. Note it marshals *CLIError
   directly and prints one line.
3. Determine the JSON path of a CLIError's "code" field on the wire.
   Expected by the RDR: .error.code (from "{\"type\":\"failed\",\"error\":{…}}").
   Actual: .code.
4. PASS only if the RDR's stated failure-envelope shape matches step 3's actual,
   or the RDR books the envelope change as an explicit Existing Infrastructure
   Audit "Extend" row naming EmitJSON and its blast radius on existing verbs.
```

*Fails on the draft.* A4 is `Verified` and cites `Fail`; the shape quoted comes from the package doc comment, which the code contradicts. The audit table's `respond` row says "Reuse."

### AT-5 — `Finding` type ownership is self-consistent (catches C-5)

```gherkin
Scenario: a subsystem-agnostic record does not carry subsystem vocabulary
  Given the normative contract requires Finding be "defined in clierr as a
        subsystem-agnostic record so clierr gains no dependency on the graph-lint package"
  And   the same contract requires a finding to carry the atom's
        "Key, Operator, Literal, and Block" and, for escape findings, "the failure class"
  When  a reviewer types those fields
  Then  Block is RDR 0002's authored-block enum (match/all/unless)
  And   failure class is RDR 0001's RefusalKind
  And   the RDR MUST state whether clierr imports those types or degrades them to strings
  And   if degraded, MUST state how the "MUST carry" obligation is enforced
```

*Fails on the draft.* Neither question is answered.

### AT-6 — A5's oracle names an artifact that exists (catches C-6)

```
Steps:
1. Enumerate .toml files in the repository, excluding _repos/.
   Result: .roborev.toml, .kata.toml, and three files under docs/rdr/*/evidence/spikes/.
2. A5 states the gate subject is "the checked-in transition model — not the
   fixture corpus."
3. PASS only if step 1 contains a file that is a transition model and is not a
   fixture, OR the RDR books its creation with a named owner and location.
4. Also PASS the disjunction check: Implementation Plan Phase 3 must not
   reintroduce "the checked-in transition model or fixture corpus" after A5
   resolved it.
```

*Fails both halves on the draft.*

### AT-7 — CI job name is unclaimed (catches C-7)

```
Steps:
1. Read .github/workflows/ci.yml, list jobs. Result: test, lint, vuln.
2. A5 requires "a new `lint` job in .github/workflows/ci.yml".
3. PASS only if no job of that name exists, or the RDR states what happens to
   the existing one.
```

*Fails on the draft.*

### AT-8 — Terminal satisfaction under the join rule (catches C-8)

```gherkin
Scenario: a two-path workflow's terminal node is not a dead end
  Given a model with two rules both writing a terminal-satisfying owned state,
        with different values (status=done and status=archived)
  And   the join rule unions per-tag value sets on merge
  And   invariant 2 requires "every value in each of the node's per-tag value
        sets" to satisfy the terminal
  When  both edges reach the same successor node
  Then  the merged node holds status: {done, archived}
  And   it satisfies neither declared terminal
  And   invariant 2 emits graph-dead-end
  Therefore the RDR MUST show a legal model edit that clears this finding,
        or state that multi-path completion is unmodelable
```

*Fails on the draft.* Scenario 10 acknowledges this exact case ("a node reaching a declared terminal on one path and a sink on another emits `graph-dead-end`") and calls it an accepted false positive whose cure is "an explicit … terminal declaration" — without exhibiting one that works.

### AT-9 — Write semantics for single-valued tags has a cited owner (catches C-9)

```
Steps:
1. Invariant 5 asserts "A write to a single-valued tag replaces its prior value
   (no preceding clear is required)" and uses it to delete the path-accumulation check.
2. Grep RDR 0002 and RDR 0003 Normative Contracts for a clause stating write
   semantics for single-valued tags.
3. PASS only if such a clause exists and is cited by invariant 5, or the RDR
   books the semantics as a new Critical Assumption with a producer.
```

*Fails on the draft.* RDR 0002's nearest clause is "Absence from both the write block and the clear list MUST NOT imply deletion," which says nothing about replacement; RDR 0003 owns the marker, not the update rule.

### AT-10 — The root satisfies its own always-present check (catches C-10)

```gherkin
Scenario: the declared initial owned state does not violate graph-always-present-owned
  Given the initial declaration is "an `initial` table of owned tag = value
        assignments fixing the root node" (A6)
  And   graph-always-present-owned requires an always-present owned key to be
        "held in every reachable owned-state node"
  And   the root is a reachable node
  When  the model declares more always-present owned keys than the initial table assigns
  Then  the root violates the invariant
  Therefore the RDR MUST either exempt the root, or state that the initial table
        MUST assign every always-present owned key
```

*Fails on the draft.* Neither is stated.

### AT-11 — Overlap population matches the kernel's candidate set (catches C-11)

```gherkin
Scenario: lint never reports an overlap the kernel cannot produce
  Given internal/resolve/resolve.go::Resolve builds candidates by
        "view.matches(row.Match)" before the guard gate
  When  lint groups two rows whose match patterns are jointly satisfiable in a
        merged abstract node but never in one concrete view
  Then  graph-overlap MUST NOT be emitted for that pair
  And   the RDR MUST state the argument that its overlap population is a
        superset of no concrete candidate set, or accept the false positive
        explicitly with a worked example
```

*Fails on the draft.* The over-approximation argument in §Load-Bearing Decisions covers *reachability* (a path lint walks may be infeasible), not *co-membership* (two rows lint co-groups may never be co-candidates). Those are different claims and only the first is argued.

### AT-12 — Per-class coverage has a runtime counterpart (catches C-12)

```gherkin
Scenario: the per-rescuable-class coverage union corresponds to a runtime state
  Given lint computes the union per (scoped row group × declared rescuable class),
        classes being no_match and ambiguous_match
  And   the kernel reaches escapeOrRefuse only after the ordinary population
        yields zero or many, and rescues on the single kind that occurred
  When  a group's ordinary rows are provably exhaustive and non-overlapping
  Then  neither no_match nor ambiguous_match can occur for that group at runtime
  And   lint MUST NOT emit graph-coverage-gap on either arm
  Therefore the RDR MUST state when the ambiguous_match arm is required to be closed
```

*Fails on the draft.* Scenario 19 asserts the uncovered arm "still emits its coverage finding" with no condition on whether that arm is runtime-reachable.

### AT-13 — "MUST NOT be accepted for resolver use" has an enforcement site (catches C-13)

```
Steps:
1. Read the Normative Contract "A model with any blocking lint finding MUST NOT
   be accepted for resolver use or CI success."
2. Read §Approach: "The kernel does not consult a lint verdict at runtime …
   a locally-edited illegal model can still be resolved against on a working tree."
3. Identify the code path that enforces "MUST NOT be accepted for resolver use."
4. PASS only if such a path exists, or the contract is reworded to the merge
   boundary the RDR actually enforces.
```

*Fails on the draft.* No such path; the contract and the approach paragraph state opposite things.

### AT-14 — One code per user action (catches C-14)

```gherkin
Scenario: graph-unprovable-coverage's arms are distinguishable by a machine
  Given the code fires on (a) a non-finite dimension, (b) a narrowing-withheld
        claim, and (c) an atom over a tag not declared single-valued
  And   the remedies are: declare a domain, accept the withholding, add a marker
  When  a CI consumer reads the finding's code field
  Then  it MUST be able to distinguish the three
  And   Scenario 3's concession that "atom fields are required only when the
        finding is attributed to one guard atom" MUST NOT be the only discriminator
```

*Fails on the draft.*

### AT-15 — The false-positive rate claim is grounded (catches C-15)

```
Steps:
1. §Trade-offs asserts a false positive "is expected to be rare and fixable by
   an explicit write, clear, or terminal declaration," and forecloses suppression.
2. Run AT-8, AT-10, AT-11 against the draft. Each produces a structural false
   positive, not an incidental one.
3. PASS only if the RDR exhibits a worked model on which the join rule produces
   zero false positives, or replaces "expected to be rare" with a measured
   number from a spike over RDR 0003's desk-trace fixture.
```

*Fails on the draft.* The claim is unbacked and three structural counterexamples are derivable from the RDR's own text.

### AT-16 — `data.findings` survives `Data`'s omitempty (catches C-16)

```
Steps:
1. Read internal/cli/respond/respond.go::Success. Note `Data any \`json:"data,omitempty"\``.
2. The RDR requires "the key MUST be emitted even when the list is empty" on
   the success surface as data.findings.
3. Determine what respond.OK emits when the verb assigns Data a struct whose
   only field is an empty slice, versus a nil interface, versus a nil map.
4. PASS only if the RDR states which of those the verb assigns, and Scenario 16's
   golden file pins it.
```

*Fails on the draft.* The Testing Strategy analyses `CLIError`'s `omitempty` fields and never looks at `Success.Data`'s.

### AT-17 — A Draft peer does not close a locked peer's records (catches C-17)

```gherkin
Scenario: closure direction respects lock status
  Given RDR 0003 is Final [locked 2026-08-22]
  And   RDR 0006 A2 states A17, A18, A19, A20 "close by citation"
  And   RDR 0006 A7 concedes for A18 that "that edit is a route-back on a locked
        peer, not a Draft amendment"
  When  a reviewer applies A7's concession uniformly
  Then  A17, A19, and A20 MUST also be booked as route-backs, not as closed
  And   the Finalization Gate MUST NOT record them as discharged
```

*Fails on the draft.* A2 claims closure for three of the four; only A18 gets the concession.

### AT-18 — Invariant 6 reads a field that carries what it needs (catches C-18)

```gherkin
Scenario: owned-set-before-match has a producer for its read set
  Given RDR 0007 (Final) narrows Row.RequiresOwned to post-guard write dependencies
        and states "guard-input decidability is not RequiresOwned's job"
  And   RDR 0002 fixes RequiresOwned as the write-block-plus-clear-list key set
        and "does not add guard-read keys to it"
  And   invariant 6 checks "a row that reads an owned tag (match key or guard atom)"
  When  a reviewer looks for the normalized field carrying that read set
  Then  none exists
  Therefore the RDR's minimum input contract MUST list the atom-derived owned
        read set as an explicit obligation, naming its producer
```

*Fails on the draft.* The minimum input contract lists "match predicates, guard predicates … writes, clears" but never states that lint must derive the owned read set itself, and never warns implementers off `RequiresOwned` — the field whose name most closely matches the invariant's prose.
