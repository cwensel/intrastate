Model: claude-opus-5[1m]

# Persona 2 - Implementer

Question: if I started coding this Monday, what would I ask in the first hour?

Findings:

1. **"Source state" is never defined, yet three normative surfaces depend on it.**
   Anchor: `Technical Design`, "**Row groups are RDR 0003's.**" — "the same source
   state and the same recognized outcome"; `Load-Bearing Decisions` →
   *Reachability relation* — "a node is an abstract owned-state (per owned tag:
   absent, or held with its set of possible declared values)"; invariant 2
   ("every reachable owned-state that satisfies no declared terminal").
   Clarification needed: RDR 0002's authored shape has no `state` field. Its
   `[[rule]]` carries `use`/`[rule.match.<tag>]`/`[rule.guard...]`/`[rule.write]`
   (`0002::Normative Contracts`, layout clause; the `prelock-flapping-cap`
   example matches `status.eq = "Draft"` under a context). So "source state" must
   be *derived*, and this RDR never says from what. Three incompatible readings
   are live: (a) the full match-pattern over owned tags (an authored predicate);
   (b) the reachability relation's abstract owned-state *node* (a set of
   per-tag value sets); (c) a single distinguished `status`-like tag. Under (a)
   two rows whose match patterns differ syntactically but denote the same owned
   states are in *different* groups and never overlap-checked; under (b) one row
   belongs to many groups (one per node satisfying it), and invariants 3/4 run
   per node, not per row set; under (c) the model would need to name the state
   tag and no clause does. RDR 0003 explicitly declines to close this — it says
   "RDR 0006 supplies the graph traversal that enumerates which selection
   contexts are reachable" (`0003`, row-group paragraph), and its A10 books the
   confirmation. Nobody defines the *identity* of a selection context.
   Blocks: the group-keying function — the single most load-bearing data
   structure in Phase 2, which every invariant 3, 4, 5 and the
   `reachable group` iteration in the desk trace step 3 hangs off. Cannot start
   Phase 2 without it.

2. **RDR 0003's "escape row closes coverage only for the classes it can actually
   rescue" clause has no counterpart in invariant 4, and `owned_state_unavailable`
   appears nowhere in this RDR.**
   Anchor: invariant 4 — "Escape rows contribute their accepted assignments to
   the union like any other row; there is no separate 'an escape row exists'
   disjunct" and "A group whose coverage is closed by a bare escape row (one
   carrying no guard atoms) passes"; the escape-row `normative` block —
   "Escape rows MUST participate in the coverage union … as RDR 0003's
   escape-row clause states"; the disposition table row "Clean model, a group
   closed by a bare escape row | 0 | `respond.OK`".
   Clarification needed: RDR 0003's cited clause is longer than what is adopted.
   It says the kernel's gate returns before `escapeOrRefuse`, so "An escape row
   therefore cannot rescue either class [`guard_unevaluable` or
   `owned_state_unavailable`], however bare its guard," and a bare escape row
   "MUST NOT be read as discharging the narrowing." The desk trace step 7 covers
   half of this (withholding dominates closure), but the *class-scoping* half is
   never stated here: does an escape row declaring only `escape = ["no_match"]`
   close coverage for its group, or only for the `no_match` arm? RDR 0002's
   escape list is a per-class list (`0002::Normative Contracts`, escape-list
   clause: `no_match` and `ambiguous_match`). And `owned_state_unavailable` —
   a class named in RDR 0003's clause — is absent from this entire document, so
   I do not know whether invariant 6's owned-before-write finding and that
   runtime class are two views of one thing or unrelated.
   Blocks: the coverage-union implementation — whether the union is computed once
   per group or once per (group × rescuable class), and whether the disposition
   table's success row is reachable at all for a `no_match`-only escape row.

3. **Invariant 5 ("Single-valued state") is the only invariant with no stated
   algorithm, and its second half is undecidable as written.**
   Anchor: invariant 5 — "no row's writes assign it two values, and no reachable
   path leaves a second value held beside one no write or clear on that path
   removes."
   Clarification needed: the first conjunct is a per-row syntactic check
   (trivial). The second quantifies over *paths*, but the reachability relation
   in `Load-Bearing Decisions` is a fixpoint over *nodes* whose abstraction is
   "held with its set of possible declared values" — it deliberately merges
   paths, so it cannot answer a per-path question. Is the intended check (a) at
   the abstract node level: any reachable node where a single-valued tag's
   possible-value set has cardinality > 1? That is a very different test — the
   node abstraction unions values across joins, so *any* tag written to two
   different values on two different edges into one node trips it, which is
   almost certainly a false positive for an ordinary enum. Or (b) genuinely
   path-sensitive, which needs a different (and unbounded, given loops — RDR 0002
   models self-loops, `0002::A1` evidence) traversal than the one this RDR
   defines? Also unstated: does a write to a single-valued tag implicitly *replace*
   the prior value, or must a `clear` precede it? RDR 0002 carries writes and
   clears as separate lists; if write-replaces-implicitly, the second conjunct is
   vacuous and only the first conjunct exists.
   Blocks: whether Phase 2 implements a second traversal, and whether A7's whole
   "invariant 5 is a sufficient producer for RDR 0003 A18" argument survives — a
   check nobody can implement discharges nothing.

4. **The reachability fixpoint has no stated termination or node-identity rule
   over a cyclic graph.**
   Anchor: `Load-Bearing Decisions` → *Reachability relation* — "an edge is a
   normalized non-escape row whose match pattern over owned tags is satisfiable
   in the source node, producing the node with that row's writes applied and
   clears removed."
   Clarification needed: nodes are tuples of per-tag possible-value sets, so the
   node lattice is finite only if every owned tag has a finite declared domain —
   but the same paragraph admits "a tag with no finite domain abstracts to
   held/absent," which is finite, so termination probably holds. What is missing
   is the *join rule*: when two edges reach the same successor with different
   value sets, are those one node (union the sets — a lattice widening) or two
   distinct nodes (path-sensitive, exponential)? Under the first reading invariant
   6's "held in every reachable owned-state that satisfies the row's match
   pattern" is checked against merged states and gets weaker; under the second
   the traversal can blow up on the RDR flow's self-loops. Neither is stated, and
   the "over-approximates the runtime … can only produce a false positive, never
   a false green" guarantee is only true under one of them.
   Blocks: the core worklist algorithm in Phase 2 and the truth of the
   no-false-green claim the `Consequences` section repeats.

5. **A6's producer edit to RDR 0002 is scheduled but unspecified, and Phase 2
   cannot start without its shape.**
   Anchor: A6 `Producer (decided at Stage 4)` — "**RDR 0002's authoring schema**,
   not a sidecar … The edit is additive but touches two of RDR 0002's normative
   clauses together"; Prerequisites — "[ ] A6: RDR 0002's schema declares the
   initial owned state and terminal states."
   Clarification needed: the *arm* is decided; the *shape* is not. Is the initial
   owned state a set of tag=value assignments inside `[model]`, a separate
   `[initial]` table, or a rule id? Are terminals a list of state names (which
   would require the "source state" identity from finding 1 to be a nameable
   thing) or a predicate over owned tags? Invariant 1 requires "terminal state
   named by a row must resolve to a declared model element" — so rows reference
   terminals by name, which implies states are named, which contradicts the
   reachability relation's anonymous abstract-owned-state nodes. Invariant 2's
   "satisfies no declared terminal" says *satisfies*, implying a predicate.
   Blocks: everything rooted at reachability — invariants 2, 6, and the whole
   `reachable group` scoping. Phase 2 is unstartable; Phase 4 fixtures cannot be
   authored.

6. **The `findings` envelope field is required on both `CLIError` and the success
   payload, but the two shipped types support it very differently.**
   Anchor: `Technical Design` — "The implementation must extend
   `internal/cli/clierr.CLIError` with an optional typed `Findings` field
   serialized as `findings`, or an equivalently named respond/clierr-owned typed
   field"; the aggregate-`CLIError` `normative` block — "an append-only optional
   typed `findings` envelope field owned by `clierr`/`respond`, not through a
   verb-local wrapper".
   Clarification needed: `respond.Success` already has `Data any`
   (`internal/cli/respond/respond.go`), so the success half is a payload struct
   under `data`, not a sibling `findings` key — but the clause says "the same
   typed field," and the `Validation` section explicitly names `Success.Data` as
   "a host for the success-side findings list." Those are two different wire
   shapes (`{"type":"ok","data":{"findings":[…]}}` vs
   `{"type":"ok","findings":[…]}`). Which does scenario 1's assertion
   ("`findings` carries only informational entries") test? Second: putting a
   `Finding` type in `clierr` means `clierr` — documented as "the leaf home … so
   other internal packages can construct CLIErrors without importing internal/cli"
   — must now own the graph-lint finding vocabulary, or the lint package must
   define `Finding` and `clierr` must import it (a cycle). Which package owns the
   `Finding` type?
   Blocks: the envelope change in Phase 3 and every JSON assertion in the MVV;
   also whether `clierr` gains a dependency it was explicitly built to avoid.

7. **A5's exit-code assertion is self-contradicting about what to assert on.**
   Anchor: A5 Evidence — "`clierr::ExitCodeFor` maps `GroupUserEnv` to exit 2 as
   scenario 2 asserts, but shares that code with `GroupInternal`, so the gate
   asserts on `Code` rather than the exit integer alone"; scenario 2 —
   "**Expected**: each run fails … with aggregate `CLIError.Code =
   graph-lint-failed`, exit code 2".
   Clarification needed: A5 says assert on `Code`, scenario 2 asserts both. In
   JSON mode `Group` is `json:"-"` (not serialized), so a CI gate reading stdout
   sees only `code`, and a gate reading `$?` sees only 2 — which cannot
   distinguish a lint failure from an internal panic-path error. If CI is meant
   to be the authority, does it parse JSON (requiring `--as=json` in the gate
   step) or branch on exit code? Related: `Makefile::check` does not depend on
   `build` (A5's own finding), and CI runs `test`/`lint`/`vuln` jobs, none of
   which calls `make check` — so which of the four surfaces gets the new step?
   Blocks: Phase 3's CI wiring, which is the whole content of MVV scenario 7 and
   the only thing that flips A5 to Verified.

8. **"Reachable predecessors" is defined over rows, but invariant 6 is stated
   over states — and the definition admits paths, not a fixpoint.**
   Anchor: `Load-Bearing Decisions` — "a row's *reachable predecessors* are the
   rows on any root-to-source path"; invariant 6 — "must find it held in every
   reachable owned-state that satisfies the row's match pattern; the declared
   initial owned state counts as a write."
   Clarification needed: the two formulations disagree. The "reachable
   predecessors" wording (also the one A3 quotes from RDR 0003: "rejects rows
   that match owned tags unless every reachable predecessor sets or preserves
   them") is a per-path row-set question — enumerating "any root-to-source path"
   is exponential and non-terminating on the self-loops RDR 0002 models.
   Invariant 6's wording is a per-node question over the fixpoint, which is
   cheap. Which is normative? They give different answers: a tag written on
   every path but cleared on a loop back-edge is "held in every reachable node"
   under neither, but "set by every reachable predecessor" under a naive row
   reading. Also: what does "preserves" mean operationally — a row that neither
   writes nor clears the tag?
   Blocks: invariant 6's implementation and its complexity class; scenario 10's
   expected output.

9. **The withholding test's "participating row" population is cited but its
   interaction with `single_valued` provability is not carried over.**
   Anchor: **Withheld exhaustiveness claims** — "Lint must not certify a row
   group exhaustive when any participating row — escape rows included — can
   refuse `guard_unevaluable` … RDR 0003's narrowing clause and its 'can refuse'
   decision procedure (a syntactic test over the declared optionality field, not
   a reachability query) govern".
   Clarification needed: RDR 0003 carries a *second* blocking trigger this RDR
   never names — its operator/kind agreement clause makes the `single_valued`
   marker "a **precondition for provability**: an `eq`/`in`/comparison atom over
   a tag not declared single-valued has no projection and takes the blocking
   outcome" (`0003::A21`). Under that clause the current fixtures produce
   `graph-unprovable-coverage` for essentially every rule, not `graph-coverage-gap`.
   This RDR's code table maps `graph-unprovable-coverage` to "Finite-domain proof
   unavailable, or claim withheld under RDR 0003's narrowing" — is the
   unmarked-single-valued case "finite-domain proof unavailable"? And RDR 0003
   also has a **too-large product** blocking case with a "declared,
   model-independent bound … the implementation MUST publish the bound it
   enforces." This RDR mints no code for it and states no bound.
   Blocks: which code a not-single-valued dimension and an over-large product
   emit, and whether the implementation must publish a cardinality bound (an
   externally-visible contract) that no clause here specifies.

10. **The `--flow`/`--model` illustrative invocation conflicts with how the model
    is actually located, and RDR 0005 owns half of it.**
    Anchor: `Illustrative Code` — "`intrastate lint --flow rdr --model
    ./path/to/rdr-transition-model.toml`" prefaced by "RDR 0005 may still adjust
    flag placement"; scenario 7 — "over the checked-in transition model or
    fixture corpus".
    Clarification needed: RDR 0005 defines `--flow <id>` as selecting the model
    via *config discovery* (`0005`, MVP request grammar; its infra audit says
    "Config discovery … Parser placeholder; no transition-model config yet"),
    while `--model <path>` is an explicit path. Are both accepted, and what
    happens when both are given? Can `lint` take multiple models or a directory
    (scenario 7 says "the checked-in transition model **or** fixture corpus" —
    corpus implies many)? Is the finding's `model` field the `[model].id` or the
    path? Since `lint` is deliberately *not* under RDR 0005's `flow` group, the
    flag vocabulary has no owner.
    Blocks: the Phase 3 command signature and every MVV invocation string;
    also whether one run emits findings across several models (the deterministic
    order clause leads with "model id", implying yes).

11. **The deterministic-order key bottoms out in an unspecified "normalized
    predicate/write fingerprint".**
    Anchor: the deterministic-order `normative` block — "model id, invariant code,
    source rule/context id or graph element id, then normalized predicate/write
    fingerprint"; `Load-Bearing Decisions` → *Identity* — "The fingerprint covers
    the attributed atom and failure class when present."
    Clarification needed: is the fingerprint a hash, or a sortable canonical
    string? Ordering requires a total order, so a hash works only if it is
    compared as bytes — but then the order is stable yet arbitrary, and a
    fixture diff churns whenever an atom is added. RDR 0002 already defines a
    canonical atom sort ("atoms sort by (key, block, operator token, literal)")
    and RDR 0003 canonicalizes set literals — are those the fingerprint's input,
    and is the fingerprint just the canonical serialization? Also, the key
    contains "source rule/context id **or** graph element id" — a disjunction
    across two namespaces with no stated tiebreak when a run mixes both (which
    scenario 5 explicitly produces).
    Blocks: the `sort.Slice` comparator and every golden-output test in Phase 4.

12. **Invariant 1 requires "every transition target … must resolve to a declared
    model element," but RDR 0002 produces no transition target.**
    Anchor: invariant 1 — "every transition target, context reference, tag,
    outcome, accessor reference, and terminal state named by a row must resolve
    to a declared model element."
    Clarification needed: RDR 0002's normalized row carries *next-state tags*
    and *writes* (`0002::Normative Contracts`; scenario 2 asserts "escape rows …
    carry neither writes nor next-state tags"), not a named target. So the
    dangling-edge check over "transition target" is either (a) each next-state
    tag key/value resolving against its `[tags.<tag>]` declaration and declared
    domain — which is already covered by the "tag" term in the same list — or
    (b) a target-state name that does not exist in the data model. If (a), the
    term is redundant and should be dropped so implementers do not go looking for
    a field; if (b), it depends on A6's unwritten schema (finding 5).
    Blocks: what `graph-dangling-edge` actually checks — and the illegal-fixture
    for it in MVV, which must trip a real code path.

13. **The advisory tier names three members but only one has a code.**
    Anchor: `Technical Design` — "the advisory tier is scoped to redundant rows,
    unreachable rules, and the bare-escape coverage closure above"; the finding
    code table lists only `graph-coverage-closed-by-escape` as `info`; the
    disposition table row "Redundant row / unreachable rule | 0 | `respond.OK` |
    informational entry".
    Clarification needed: what stable codes do "redundant row" and "unreachable
    rule" carry? The versioning cross-cutting concern requires codes be "stable
    and append-only," and MVV scenario 1 asserts "`findings` carries only
    informational entries" on the legal fixture — so a test must name them. Also
    undefined: what makes a row *redundant* (accepted-assignment subset of
    another row's? — but that is exactly an overlap, which is blocking), and is
    an *unreachable rule* a row whose selection context the reachability fixpoint
    never reaches (making it a direct consumer of finding 1's group identity)?
    Blocks: the informational half of the finding taxonomy and MVV scenario 1's
    assertion set.

14. **A7 asserts invariant 5 discharges RDR 0003 A18, but A18 also names a
    view-level obligation invariant 5 structurally cannot cover.**
    Anchor: A7 — "A18 names two sufficient producers and requires only one: RDR
    0007's assembly rejecting a non-conforming view, **or** this RDR's lint
    reporting a model-level conformance violation."
    Clarification needed: RDR 0003's conformance clause defines a conforming view
    as one where "every always-present key is present in it and every
    single-valued tag holds at most one of its declared domain values." Invariant
    5 covers only the single-valued half — it says nothing about *always-present
    keys being present*. Since observed and recognized tags come from accessors
    and the caller at runtime (`resolve.go::assemble` merges owned, observed and
    recognized), a model-level lint cannot prove an always-present *observed* key
    will be present. So does invariant 5 need an always-present clause for owned
    tags at least, or does A7 knowingly discharge only half of A18?
    Blocks: invariant 5's scope, and whether A7 flips to Verified honestly.

15. **Scenario 11's control case tests a distinction the invariant text asserts
    but no clause makes decidable for a non-optional-declared key.**
    Anchor: invariant 4 — "A key declared always-present contributes no presence
    dimension"; scenario 11 — "a control group whose `exists` atom is over a key
    declared always-present … the always-present control passes."
    Clarification needed: RDR 0003's declaration model has an *optionality
    marker*, and its assignment-count table (per `0003:1535`) says an **unmarked**
    key "takes the ×2 presence row" — i.e. unmarked defaults to *optional*, not
    always-present. So "declared always-present" is a third state (explicitly
    marked present) distinct from "unmarked." Does invariant 4 mean explicitly
    marked, or merely not-marked-optional? Under the former, the control fixture
    needs a marker RDR 0002 may not yet carry (`0003::A16` covers the
    single-valued marker's authoring location but this is a different field).
    Blocks: the presence-dimension projection code and scenario 11's control
    fixture.
