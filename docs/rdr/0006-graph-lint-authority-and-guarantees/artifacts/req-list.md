# REQ List — RDR 0006 Graph Lint Authority And Guarantees

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0006-graph-lint-authority-and-guarantees.md`. Quotes are verbatim,
taken from the projector (`rdr inspect --select <id>`) for fenced elements and
read from the record for prose outside the fences.

Element ids (`0006:C4`, `0006:MVV`) are carried wherever a REQ derives from a
labelled contract, so a later stage can trace the REQ back to its contract.
The record has 20 `C` elements but far more than 20 REQs: several fences carry
multiple independent obligations, and this record puts the **whole invariant
taxonomy, the finding-code table, the `reason` discriminator, the disposition
table, the reachability relation, and the node ceiling in unfenced prose**.
The fenced set is not the testable surface — it is roughly half of it.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced)
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `AP` = Proposed Solution / Approach
- `INV` = Technical Design, the mandatory invariant set (numbered 1–7)
- `CODES` = Technical Design, the finding-code table
- `REASON` = Technical Design, the `graph-unprovable-coverage` `reason` table
- `DISP` = Technical Design, the disposition table (input class × outcome)
- `LBD` = Technical Design / Load-Bearing Decisions
- `RT` = Technical Design / Round-Trip / Inverse Invariants
- `IC` = Technical Design / Illustrative Code (**input contract is normative**)
- `AUDIT` = Existing Infrastructure Audit
- `CONS` = Trade-offs / Consequences
- `FM` = Trade-offs / Failure Modes
- `TS` = Validation / Testing Strategy (incl. the desk trace)
- `SC-n` = Validation, numbered scenario *n*
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (prerequisites, phases)
- `PE` = Validation / Performance Expectations

---

## Authority and placement

- [REQ-1] "Graph lint MUST be a blocking acceptance gate over the normalized transition model. A model with any blocking lint finding MUST NOT be accepted for CI success, and MUST NOT reach the resolver **through the merge boundary** the gate guards." — (NC, `0006:C1`)

- [REQ-2] "That boundary is the enforcement site and the only one: the kernel does not consult a lint verdict at runtime, so a locally-edited illegal model on a working tree can still be resolved against, and this clause does not claim otherwise." — (NC, `0006:C1`) — negative REQ: no resolver-side lint consultation may be implemented.

- [REQ-3] "Graph lint MUST consume the normalized candidate-row graph from the transition model contract. It MUST NOT define a second sparse-source parser or a parallel transition semantics." — (NC, `0006:C2`)

- [REQ-4] "The authoritative CLI surface for graph acceptance MUST be the root command `intrastate lint` or a same-engine CI invocation of that command. Pre-commit hooks, aliases, and resolver-local validation flags MAY call that engine, but MUST NOT define different acceptance rules." — (NC, `0006:C19`)

- [REQ-5] "Lint command success and failure MUST route through `respond.OK`, `respond.Fail`, and `CLIError`; the command MUST NOT write directly to stdout or stderr." — (NC, `0006:C20`)

- [REQ-6] "Root `lint` is deliberately not under RDR 0005's `flow` group" — the command is registered at root, not under `flow`. — (AP)

- [REQ-7] "every such path must share the same command/request builder and graph-lint engine; none defines acceptance." Assertable-now form: "the lint package exposes exactly **one** exported engine entry point and one request builder, and a test asserts the root command calls them" — (AP, SC-6)

- [REQ-8] "The lint engine boundary is an internal graph-lint package that receives a normalized graph value, not Cobra command state and not sparse TOML." — (TD)

- [REQ-9] "the canonical command and subsystem name is \"lint\"." — (LBD, Naming)

## Command input contract (normative, not illustrative)

- [REQ-10] "The command's **input contract is normative** (the name cannot be authoritative while the way to invoke it is illustrative); only cosmetic flag spelling defers to RDR 0005." — (IC)

- [REQ-11] "`--flow <id>` selects the model through RDR 0005's config discovery; `--model <path>` names one explicitly. They are mutually exclusive — supplying both is a `GroupUserEnv` usage error, not a precedence rule." — (IC)

- [REQ-12] "Exactly one model is linted per invocation: a corpus is linted by invoking the command once per model, so a run's findings never span models." — (IC)

- [REQ-13] "The finding's `model` field is the `[model].id` from the model itself, never the path." — (IC)

- [REQ-14] "Until RDR 0005 lands discovery, only `--model <path>` is instantiable." — A9; `--model <path>` is the form the gate and the MVV use. — (CA A9, IP)

## Minimum input contract the lint engine receives

- [REQ-15] The minimum input contract is: "model identity and version"; "normalized candidate rows with deterministic row identity, source rule id, optional source span, match predicates, guard predicates (atoms in RDR 0007's `Key`/`Operator`/`Literal`/`Block` shape), writes, clears, and — for escape rows — row kind and the declared list of failure classes the row rescues"; declared tags with provenance plus RDR 0003's tag declaration model; "recognized outcome alphabet, declared initial owned state, and declared terminal states"; "accessor references and context references only as normalized identifiers needed for dangling-reference diagnostics." — (TD)

- [REQ-16] "Lint reads every one of these from the declaration and never infers one from a tag's name, value spelling, or a fixture" — (TD)

- [REQ-17] The pipeline has four stages: "load and normalize the transition model through the RDR 0002 table contract"; "derive a graph view from normalized candidate rows"; "run invariant checks over that graph and collect typed findings"; "emit either a success report or one aggregate structured `CLIError` failure through `respond`." — (TD)

- [REQ-18] "The read-set is computed by lint from the row's atoms, never read from `Row.RequiresOwned`." "Lint therefore derives the owned read-set as the owned-provenance keys named by the row's match pattern and its `all`/`unless` guard atoms" — (INV 6)

## Row grouping and selection context

- [REQ-19] "Coverage, overlap, and withholding MUST be decided per scoped row group as RDR 0003 defines it … This RDR MUST NOT define a second grouping predicate; it supplies only which selection contexts are reachable, per the reachability relation in `Load-Bearing Decisions`." — (NC, `0006:C5`)

- [REQ-20] "**Source state is the authored match pattern.** … It is the row's **authored match pattern**, the field `internal/resolve/resolve.go::Resolve` compares (`view.matches(row.Match)`) before the guard gate: rows sharing a match pattern and a recognized outcome are one group." — (TD)

- [REQ-21] "Two rows whose patterns differ are two groups even when some abstract node satisfies both, because no concrete view makes them candidates together." — (TD; the negative control of SC-20)

- [REQ-22] "The reachability relation below decides **which** groups are proven, never **which rows are in one**. Membership is syntactic and per-row; reachability is a filter over contexts." — (TD)

- [REQ-23] "invariants 3 and 4 run once per reachable group." — (TD)

- [REQ-24] "Group *membership* is syntactic (the authored match pattern), so the group count is bounded by the authored rule count regardless of the lattice; only the reachability filter scales with nodes." — (PE)

## Ambiguity, determinism, overlap

- [REQ-25] "Graph lint MUST reject ambiguity instead of relying on source order, rendered-row order, or first-match priority to choose between enabled rows." — (NC, `0006:C4`)

- [REQ-26] "when two ordinary rows qualify for the same selection context, lint rejects the model. It never selects by order; escape rows are modeled graph edges with their own per-class overlap population, not a tie-breaker and not a coverage exemption." — (LBD, Selection / predicate)

- [REQ-27] "**Determinism / overlap** — within a scoped row group, no finite-domain input assignment may enable two ordinary (non-escape) rows." — (INV 3)

- [REQ-28] "Escape rows are checked for overlap in their own populations, one per declared failure class … an escape row overlapping an ordinary row is not a runtime ambiguity and is not reported." — (INV 3; NC `0006:C9`; DISP; SC-9a)

- [REQ-29] "Escape rows MUST participate in the coverage union and MUST be overlap-checked in one population per declared failure class, never against ordinary rows" — (NC, `0006:C9`)

- [REQ-30] "One overlapping pair MUST yield one finding naming both rows; escape rows overlapping across several shared failure classes MUST yield one finding per shared class." — (NC, `0006:C16`; SC-9b, SC-13)

## The mandatory invariant set (7 classes)

- [REQ-31] "Graph lint MUST check at least these blocking invariant classes: dangling edge, dead end, determinism/overlap, guard exhaustiveness/gap, single-valued state, owned-set-before-match, and declared terminal/escape handling." — (NC, `0006:C3`)

- [REQ-32] **Invariant 1, dangling edge** — "every context reference, tag key, tag value, outcome, and accessor reference named by a row must resolve to a declared model element, and the model must declare an initial owned state." — (INV 1)

- [REQ-33] "There is no \"transition target\" term … the edge's destination is checked as tag keys and values against their `[tags.<tag>]` declarations and declared domains. Terminal declarations are checked the same way — as tag predicates, not as state names." — (INV 1)

- [REQ-34] **Invariant 2, dead end** — "every reachable owned-state node that satisfies no declared terminal must be the source of at least one modeled non-escape row (escape self-loops do not count as progress)." — (INV 2)

- [REQ-35] "A terminal is a predicate over owned tags, and a node *satisfies* it when **every** value in each of the node's per-tag value sets meets it. Partial satisfaction is deliberately not enough" — (INV 2)

- [REQ-36] "**The test runs on split nodes, not merged ones.** … Before evaluating this invariant, lint splits each reachable node on the terminal-participating keys: a node whose `status` set is `{done, review}` against a terminal `status=done` is tested as `{done}` … and `{review}` … — independently." — (INV 2; regression guard SC-21)

- [REQ-37] "The split is bounded by the declared domains of terminal-participating keys only, never the whole lattice." — (INV 2)

- [REQ-38] **Invariant 4, exhaustiveness/gap** — "every scoped row group whose participating guard dimensions are all finitely declared claims closed coverage (default-on, never an opt-in annotation), and lint must prove `union(row_i accepted assignments) == scoped product` over RDR 0003's derivation." — (INV 4; NC `0006:C6`)

- [REQ-39] "Escape rows contribute their accepted assignments to the union like any other row; there is no separate \"an escape row exists\" disjunct." — (INV 4)

- [REQ-40] **Invariant 5, single-valued state** — "the model must not produce a view in which a tag declared single-valued … holds two values. This is decided **per row, syntactically**: no row's write block may assign a single-valued tag two values." — (INV 5)

- [REQ-41] "It is deliberately *not* decided over the reachability nodes: a merged node's value set has cardinality > 1 whenever two paths write different single values, which is the legal shape of a two-path merge, not a violation" — (INV 5) — negative REQ: no node-level single-valued check.

- [REQ-42] **Invariant 6, owned-set-before-match** — "a row that reads an owned tag (match key or guard atom) must find it held in every reachable owned-state that satisfies the row's match pattern; the declared initial owned state counts as a write." — (INV 6)

- [REQ-43] **Invariant 7, declared terminal/escape handling** — "terminal states and escape rows are explicit model data; lint must not infer them from missing rows." The minted defect is "a reachable non-terminal node with no outgoing non-escape row and no terminal declaration covering it (an implied terminal), or a group relying on an undeclared escape row to close coverage." — (INV 7)

- [REQ-44] "A model that declares no initial owned state MUST be rejected with a blocking finding; lint MUST NOT treat an absent root as an empty reachable set and report a clean model." — (NC, `0006:C18`; DISP row; SC-15). The disposition table names the code: "`graph-dangling-edge` naming the missing `[model]` declaration".

## Exhaustiveness scope, provability, and withholding

- [REQ-45] "Graph lint MAY claim exhaustiveness only over finite declared domains supplied by RDR 0003's tag declaration model, and MUST read finite domain, optionality, and single-valuedness from that declaration — never inferred from a tag's name, value spelling, or a fixture." — (NC, `0006:C6`)

- [REQ-46] "If a participating dimension is not finite or not projectable, lint MUST emit `graph-unprovable-coverage` for that dimension." — (NC, `0006:C6`)

- [REQ-47] "The claim is default-on for every scoped row group whose participating dimensions are all finitely declared." — (NC, `0006:C6`)

- [REQ-48] "The scoped product MUST include the presence dimension an `exists` atom over an optional key contributes … lint MUST NOT certify a group exhaustive while ignoring that dimension's `{absent}` assignment." — (NC, `0006:C6`; SC-11)

- [REQ-49] "A key declared always-present contributes no presence dimension." — (INV 4; SC-11 control)

- [REQ-50] "Coverage is a universal claim, so it MUST be computed from the group's authored rows and their declared guard domains alone — never from a reachability node." — (NC, `0006:C7`)

- [REQ-51] "The union `union(row_i accepted assignments)` and the scoped product both range over the group's *participating guard dimensions* as RDR 0003 defines participation; match keys are not product dimensions … and contribute no assignment to either side." — (NC, `0006:C7`)

- [REQ-52] "Lint MUST NOT widen a group's row set by pooling rows satisfiable at a shared abstract node" — (NC, `0006:C7`; SC-20's negative control)

- [REQ-53] "Reachability decides only whether a group is proven at all; it never enters the coverage computation." — (NC, `0006:C7`)

- [REQ-54] "lint MUST NOT certify a row group exhaustive when a participating row can refuse `guard_unevaluable`. A withheld claim MUST be emitted as `graph-unprovable-coverage`, naming the participating row and the refusing atom." — (NC, `0006:C8`)

- [REQ-55] "This RDR cites that clause and MUST NOT restate it, mint a second code for it, or carry it in a non-blocking tier." — (NC, `0006:C8`) — negative REQ: no second code, no warning tier for withheld claims.

- [REQ-56] "Lint must not certify a row group exhaustive when any participating row — escape rows included — can refuse `guard_unevaluable`" — the "can refuse" test is "a syntactic test over the declared optionality field, not a reachability query". — (TD, Withheld exhaustiveness claims; TS desk-trace step 5)

- [REQ-57] "Lint MUST publish the model-independent product bound above which it declines to prove coverage, and MUST emit `graph-product-too-large` for a group whose declared finite product exceeds it" — (NC, `0006:C12`)

- [REQ-58] "An atom over a tag not declared single-valued has no projection and MUST take `graph-unprovable-coverage` (`0003::A21`)." — (NC, `0006:C12`)

- [REQ-59] "The bound is an implementation constant, published in the command's help output and asserted by the MVV, not a per-model input." — (TD; SC-18)

- [REQ-60] "Lint MUST therefore publish a **model-independent node ceiling** alongside the product bound, and MUST emit `graph-product-too-large` naming the traversal (not a group) when a model's reachable node set exceeds it, rather than running unboundedly. Both bounds are implementation constants in the command's help output and asserted by the MVV." — (PE; DISP row; SC-22)

## Escape-row class scoping

- [REQ-61] "An escape row closes coverage only for the failure classes it declares." — (NC, `0006:C10`)

- [REQ-62] "Lint MUST therefore compute the coverage union per (scoped row group × declared rescuable class) — `no_match` and `ambiguous_match` only — and MUST NOT let a row declaring one class close the group's other arms." — (NC, `0006:C10`; SC-19)

- [REQ-63] "Lint MUST therefore check the `ambiguous_match` arm's coverage only for a group that carries a `graph-overlap` finding, and MUST treat the arm as vacuously closed for a group whose ordinary population is overlap-free." — (NC, `0006:C10`; DISP "Silent by design" row; SC-19b)

- [REQ-64] "`owned_state_unavailable` is a runtime refusal class no escape row rescues and no lint finding mints; invariant 6's `graph-owned-before-write` is the design-time check whose runtime counterpart it is, and the two MUST NOT be conflated." — (NC, `0006:C10`)

- [REQ-65] "A group whose coverage is closed by a bare escape row MUST emit `graph-coverage-closed-by-escape` naming that row; a bare green MUST NOT satisfy this clause." — (NC, `0006:C9`; SC-9c)

- [REQ-66] "A group whose coverage is closed by a bare escape row (one carrying no guard atoms) passes, but the verdict must say so" — the closure does not fail the run. — (INV 4; DISP)

## Always-present owned keys

- [REQ-67] "An owned key declared always-present MUST be held in every reachable owned-state node; a violation is `graph-always-present-owned`." — (NC, `0006:C11`)

- [REQ-68] "Lint MUST NOT extend this check to observed or recognized keys — those arrive at runtime and `internal/resolve/resolve.go::assemble` reads no declaration — so this RDR discharges only the owned half of RDR 0003's conformance premise." — (NC, `0006:C11`; SC-17)

- [REQ-69] "**The root is included, and that is a requirement on the `initial` declaration, not an exemption.** The declared initial owned state is a reachable node by construction, so an always-present owned key the `initial` table omits is absent at the root and takes the finding there." — (INV 5)

- [REQ-70] "lint reports the omission against the `initial` declaration rather than against an arbitrary downstream node, so the diagnostic names the authored site." — (INV 5)

## Finding shape, codes, and severity

- [REQ-71] "Every blocking finding MUST carry a stable code, model identity, severity, human-readable message, and the source rule/context id or source span when the normalized model can provide one." — (NC, `0006:C13`; SC-5)

- [REQ-72] "A finding attributed to one guard atom MUST carry that atom's `Key`, `Operator`, `Literal`, and `Block`; a finding scoped to an escape population MUST carry the failure class." — (NC, `0006:C13`)

- [REQ-73] The blocking code set is exactly: `graph-dangling-edge`, `graph-dead-end`, `graph-overlap`, `graph-coverage-gap`, `graph-unprovable-coverage`, `graph-single-valued-state`, `graph-always-present-owned`, `graph-owned-before-write`, `graph-terminal-escape`, `graph-product-too-large`. — (CODES)

- [REQ-74] "The advisory tier is closed at `graph-coverage-closed-by-escape`, `graph-redundant-row`, `graph-unreachable-rule`, and `graph-vacuous-atom`." — (NC, `0006:C17`; CODES, severity `info`)

- [REQ-75] "Advisory findings MUST NOT change the success disposition." — (NC, `0006:C17`; DISP)

- [REQ-76] "the advisory tier … must not absorb a withheld claim." — (TD) — negative REQ paired with REQ-55.

- [REQ-77] "A redundant row is one whose accepted assignments are a proper subset of a sibling's in the same group" — code `graph-redundant-row`, `info`. — (NC, `0006:C17`; TD)

- [REQ-78] "an unreachable rule is one no reachable owned-state node satisfies, including RDR 0002's dead-rule case" — a predicate set requiring `recognized` to be absent. Code `graph-unreachable-rule`, `info`. — (NC, `0006:C17`; TD; SC-14)

- [REQ-79] "a vacuous atom is an `exists` atom over an always-present key, which RDR 0003 requires be reported as well-formed-but-vacuous rather than rejected." Code `graph-vacuous-atom`, `info`. — (NC, `0006:C17`; DISP)

- [REQ-80] "**`graph-unprovable-coverage` carries a `reason` discriminator.** … The finding MUST carry a stable `reason` from this closed, append-only set, naming the key, dimension, or row at fault" — the set is `dimension-not-finite`, `tag-not-single-valued`, `row-can-refuse`. — (REASON)

- [REQ-81] "`row-can-refuse` is *not* a model defect — it is an honest withholding — and the message MUST say so rather than reading as an authoring error." — (REASON)

- [REQ-82] "The atom fields are required only when the finding is attributed to one guard atom (the withheld-claim arm, scenario 8); the dimension arm must still name the dimension, so no arm of this code may carry the code alone." — (SC-3)

## Emission completeness and ordering

- [REQ-83] "Graph lint MUST report every defect it can decide in one pass over a row group, not the first it encounters" — (NC, `0006:C16`; SC-13)

- [REQ-84] "Withholding a group's exhaustiveness claim MUST NOT suppress overlap, coverage, or further withholding findings for that group." — (NC, `0006:C16`; SC-8)

- [REQ-85] "Each unprovable dimension, each refusing row, each overlapping pair, and any coverage gap over a provable product is its own finding, so the emitted set never depends on row or dimension iteration order." — (TD)

- [REQ-86] "Graph lint findings MUST be emitted in deterministic order by finding identity: model id, invariant code, source rule/context id or graph element id, then normalized predicate/write fingerprint." — (NC, `0006:C15`; LBD Identity)

- [REQ-87] "The fingerprint MUST be a canonical sortable serialization, never a hash: RDR 0002's canonical atom sort (`(key, block, operator token, literal)`) over the row's predicate atoms and next-state tags, with RDR 0003's canonical set-literal form." — (NC, `0006:C15`)

- [REQ-88] "When one run mixes identity namespaces, a source rule/context id sorts before any graph element id." — (NC, `0006:C15`)

- [REQ-89] "JSON finding order is deterministic by finding identity, asserted as full-list equality against a golden file, not merely as a sorted property." — (SC-4)

## Output envelope

- [REQ-90] "Graph lint failure MUST return one aggregate `CLIError` with code `graph-lint-failed` and `GroupUserEnv`" — exit 2. — (NC, `0006:C14`; DISP)

- [REQ-91] "the individual blocking findings MUST remain machine-readable in JSON mode through an append-only typed `findings` field owned by `clierr`, not through a verb-local wrapper or a text-only `Detail` string." — (NC, `0006:C14`)

- [REQ-92] "That field is a **top-level sibling of `code`** on the marshalled `CLIError` — reached as `.findings` — because `clierr.EmitJSON` marshals the `*CLIError` itself and emits no `{\"type\":\"failed\",\"error\":{…}}` wrapper; lint MUST NOT restructure that envelope to create one." — (NC, `0006:C14`; LBD Wire/byte format)

- [REQ-93] "Lint success MUST carry non-blocking findings under the existing `respond.Success.Data` payload as `data.findings`." — (NC, `0006:C14`)

- [REQ-94] "On both surfaces the key MUST be emitted even when the list is empty — the empty list is the proof's receipt — so the `findings` field MUST NOT be `omitempty`, **and** the success path MUST assign `respond.Success.Data` a non-nil struct value, since `Data` itself is `json:\"data,omitempty\"` and a nil payload would drop the key." — (NC, `0006:C14`; SC-16)

- [REQ-95] "The `Finding` type MUST be defined in `clierr` as a subsystem-agnostic record so `clierr` gains no dependency on the graph-lint package." — (NC, `0006:C14`; TD Type ownership)

- [REQ-96] "`clierr.Finding` carries those as **declared string-typed fields, not enums**: `Key`, `Operator`, `Literal`, `Block`, and `Class` are `string`, and `clierr` ascribes them no meaning. The producing package owns the vocabulary and MUST populate them from its own typed values" — (TD)

- [REQ-97] "Text mode MUST enumerate every finding's code and message, which requires extending `clierr.EmitText` (failure) and the `respond.OK` text branch (success, which renders only `Notes`/`Warnings` today and drops `Data`) — both are booked edits, not reuse." — (NC, `0006:C14`; AUDIT)

- [REQ-98] "Text mode's oracle is concrete: the set of finding codes appearing in text output equals the set in JSON output for the same fixture. Layout and wrapping are unasserted; a text renderer that drops any finding fails." — (SC-4)

- [REQ-99] "both modes return the same exit behavior with no direct stdout/stderr writes." — (SC-4)

- [REQ-100] "Model unreadable / not conforming to RDR 0002's schema | per RDR 0002 | refused before normalization; lint never runs | none of this RDR's codes" — (DISP)

## Reachability relation (this RDR's contribution; peers quantify over it)

- [REQ-101] "the graph lint traverses is the **owned-state graph**: a node is an abstract owned-state (per owned tag: absent, or held with its set of possible declared values; a tag with no finite domain abstracts to held/absent); the root is the declared initial owned state (A6); an edge is a normalized non-escape row whose match pattern over owned tags is satisfiable in the source node, producing the node with that row's writes applied and clears removed." — (LBD, Reachability relation)

- [REQ-102] "**Escape rows are edges too, and they are self-loops** … Lint MUST model them as such rather than omitting them, because a node reachable only through an escape rescue is reachable at runtime" — (LBD)

- [REQ-103] "Escape self-loops are excluded from invariant 2's \"at least one outgoing row\" test — a self-loop is not progress — so a node whose only exit is an escape rescue is still a dead end, and from invariant 3's ordinary-row overlap population" — (LBD)

- [REQ-104] "Match atoms over observed/recognized tags and all guard atoms are **not** pruned — every such edge is taken as traversable." — (LBD) — no guard-aware pruning.

- [REQ-105] "A selection context is *reachable* when some reachable owned-state satisfies its match pattern; a row's *reachable predecessors* are the rows on any root-to-source path." — (LBD)

- [REQ-106] "This is a syntactic dataflow over declarations, writes, and clears — no accessor execution, no runtime trace, no guard evaluation — so it is total, and it over-approximates the runtime" — (LBD)

- [REQ-107] "It is a different predicate from RDR 0003's \"can refuse\" test, which is decided over the optionality field alone and never consults this relation." — (LBD)

- [REQ-108] "the traversal is a **fixpoint over merged nodes, never path-sensitive**: two edges reaching the same successor produce one node whose per-tag value sets are the union of theirs (a lattice widening), and the worklist runs to a fixpoint." — (LBD, Join rule and termination)

- [REQ-109] "The lattice is finite — finitely-declared tags range over their declared domain's subsets, a tag with no finite domain over `{held, absent}` — so the fixpoint terminates on cyclic graphs, including the self-loops RDR 0002 models." — (LBD)

- [REQ-110] "**Existential checks read merged nodes directly.** Overlap …, dangling edge, and owned-set-before-match ask whether *some* witness exists." — (LBD, Soundness direction)

- [REQ-111] "**Universal checks must not.**" Coverage's remedy is that "it never reads a node at all"; dead end's remedy is to "**split** the node on terminal-participating keys first, recovering exactness." — (LBD)

- [REQ-112] "invariant 6's \"every reachable owned-state that satisfies the row's match pattern\" is the normative form, evaluated against merged fixpoint nodes. \"Reachable predecessors\" is descriptive prose for the same relation, not a second per-path enumeration; where the two readings differ the node form governs." — (LBD)

- [REQ-113] "A row *preserves* a tag when it neither writes nor clears it, so the tag's value set passes through the edge unchanged." — (LBD)

- [REQ-114] "a lint finding is identified by `(model id, invariant code, source rule/context id or graph element id, normalized predicate/write fingerprint)`. The fingerprint covers the attributed atom and failure class when present." — (LBD, Identity)

## Precedence when two clauses bear on one verdict

- [REQ-115] Desk-trace step 7: "A group both closable by a bare escape row **and** carrying a row that can refuse resolves to `graph-unprovable-coverage` (blocking, exit 2) — **not** a success with `graph-coverage-closed-by-escape`. Withholding dominates closure" — (TS; SC-12, whose negative half is the assertion that matters)

## No waiver channel

- [REQ-116] "**There is no suppression or waiver mechanism** — no inline ignore comment, no allowlist." — (CONS) — negative REQ.

- [REQ-117] "the cure is a clearer model, never a weaker lint." Accepted false positives (SC-10) "are accepted false positives whose cure is an explicit write, clear, or terminal declaration on the model, never a guard-aware lint." — (CONS; SC-10)

## Round-trip

- [REQ-118] "This RDR introduces no encode/decode, import/export, or inverse operation." — (RT) — negative REQ: no round-trip surface is added.

## Minimum Viable Validation

- [REQ-MVV] "Add a fixture-backed lint invocation that uses the same command shape intended for CI. It must pass one legal transition model and fail one illegal model for each blocking invariant class named in this RDR, asserting stable finding codes and source rule/context identity in JSON mode. The illegal fixture matrix must assert `graph-dangling-edge`, `graph-dead-end`, `graph-overlap` (both an ordinary-row pair and an escape-row pair sharing a failure class), `graph-coverage-gap`, `graph-unprovable-coverage` (both a non-finite dimension and a withheld claim naming row and atom), `graph-single-valued-state`, `graph-always-present-owned`, `graph-owned-before-write`, `graph-terminal-escape`, and `graph-product-too-large`. It must also include one **multi-defect group** asserting every expected code from a single run, so a first-failure engine cannot pass. The legal matrix must include a group closed by a bare escape row, a redundant row, and an unreachable rule, asserting `graph-coverage-closed-by-escape`, `graph-redundant-row`, and `graph-unreachable-rule` on success. Every blocking run must return aggregate `CLIError.Code = graph-lint-failed` with `GroupUserEnv` exit behavior and a machine-readable finding list; every run, clean or not, must emit the `findings` key even when empty." — (MVV, `0006:MVV`)

## Validation scenarios carrying obligations beyond the MVV text

- [REQ-119] SC-7 (CI gate, A5): "`.github/workflows/ci.yml` contains a **`graph-lint`** job — not `lint`, which is the shipped golangci-lint job and must survive untouched — that runs `make build` and then `./bin/intrastate lint --as=json` over the checked-in transition model, with the assertion reading the JSON `code` field, not the exit integer alone." Also: "`make check` is additionally wired to the same command for local parity, which requires adding the `build` edge `check` currently lacks." — (SC-7; CA A5; IP Phase 3)

- [REQ-120] SC-7: "A deliberately illegal commit to the checked-in model must fail that job." — (SC-7)

- [REQ-121] A8: "a transition model is checked into this repo for that gate to lint." Scenario 7 has no subject until one is authored and homed. — (IP Prerequisites; CA A8)

- [REQ-122] SC-23 (false-positive census): "lint the checked-in transition model (A8) in its accepted state. … the blocking finding count is recorded. Zero is the pass condition; a non-zero count on a model its maintainers accept is the mechanical trigger for the guard-aware-pruning successor RDR named in `Consequences`, and MUST be reported rather than waived." — (SC-23; CONS)

- [REQ-123] SC-11 control: "The control's key must carry RDR 0003's **explicit** always-present marker: an unmarked key defaults to optional and takes the ×2 presence row, so an unmarked control would not test the distinction." — (SC-11)

- [REQ-124] SC-20 negative control: "**no run may exit 0 with an empty blocking list**" for the two-different-match-pattern fixture; the paired positive control (same match pattern) "lints clean". — (SC-20)

- [REQ-125] SC-21: "`graph-dead-end` is **not** emitted for the terminated branch (invariant 2 splits before testing); it *is* emitted if the live branch has no outgoing row." — (SC-21)

- [REQ-126] SC-22: "The ceiling is asserted to be the published constant from the command's help output, not a test-local value." — (SC-22)

- [REQ-127] SC-5: "a finding whose normalized row lacks a rule id but has a source span or graph element id … remains actionable and stable by carrying the available source span or graph element id in the identity fields." — (SC-5)

- [REQ-128] SC-16: "the `findings` key is present with value `[]` on success (`data.findings`), asserted against a golden JSON file so an omitted key fails, and present on the failure envelope alongside blocking findings." — (SC-16)

## Deviations linkage (`artifacts/deviations.md`, pre-seeded — not modified here)

D1 and D2 were pre-seeded by the `0002-0009` cluster gate, iteration 3. Their
Stage-8 checks bear on the REQ set as follows.

- [REQ-129] **D1 check 2 mints a code this record's taxonomy does not yet
  carry.** "a root `terminal` entry naming a context whose predicate reads a
  **non-owned** tag → one blocking finding (§D7(i))". JDR 0001 §D7(i) states
  it directly: "A terminal context matching a non-owned tag is a 0006 blocking
  finding." The producer side is already landed and routes to this RDR by name:
  RDR 0002 (`Implemented`) fences that "whether a model *omits* `[initial]` or
  `terminal`, and whether a terminal context matches a non-owned tag, are RDR
  0006's blocking findings, not load failures", and its ownership table row
  reads "missing `[initial]` / `terminal`, terminal context on a non-owned tag
  | lint | RDR 0006 | loud, blocking". Implement as a blocking finding with a
  code row and a scenario. See ASSUMPTION-1 for the code chosen. D1 requires
  escalation only "if 0006's finding taxonomy cannot admit the code without a
  fenced change" — it can: `0006:C3` says "at least these blocking invariant
  classes", and invariant 1 already checks terminal declarations "as tag keys
  and values", so a non-owned terminal key is a dangling-reference defect
  inside an already-fenced class.

- **D1 check 1** is a documentation check over A6/A10's `Pending` status, not a
  code REQ. Its precondition is satisfied upstream: RDR 0002 is `Implemented`
  with root `[initial]` and root `terminal` in the layout
  (`internal/table/model.go::Model.Initial`, `::Model.Terminal`) and the
  write-replaces clause fenced ("A write replaces. A write assigns a tag's
  whole value and supplants whatever was held"). So no SPEC-DEFECT escalation
  against 0002 is warranted. Per the global instruction that RDRs are never
  amended, the status-line edit is not made; the discharge is recorded here and
  in `deviations.md` at Stage 8's disposition step.

- **D1 check 3** ("not a bare identifier" prose vs §D7's context-id list) is a
  prose repair with "no test" — no REQ. It does not change REQ-33: §D7(i)
  makes `terminal` a list of **context ids** that normalization dereferences to
  predicate sets, and RDR 0002 already implements exactly that
  (`internal/table/load.go::loadTerminal` — "A terminal context id is a
  REFERENCE, not a state name"). Lint therefore receives `[][]Atom`, i.e. tag
  predicates, which is what REQ-33 and REQ-35 require. The prose and the
  landed shape agree on substance.

- **D2** is a citation-repair entry over stale "RDR 0002 is `Draft`" claims
  (`0006:14, 141-142, 253-256, 280, 380, 391, 813, 1421, 1425`). No REQ derives
  from it; it constrains no behaviour. Its only bearing on the REQ set is that
  A6, A9, and A10's "Pending" framing must not be read as gating implementation
  — 0002, 0003, and 0007 are all `Implemented`, so REQ-101's root, REQ-35's
  stop set, and REQ-40's write-replaces premise are all live, not deferred.

---

## ASSUMPTIONS

- ASSUMPTION-1: **The D1 terminal-non-owned finding takes the existing
  `graph-dangling-edge` code rather than a new one.** Wording: D1 says "Add the
  code row and a scenario", which reads as possibly minting a new code, while
  `0006:C17` closes the *advisory* tier at four members and the blocking table
  is presented as a taxonomy, not an explicitly closed set. The defensible
  single reading: `0006:C3` requires "at least these blocking invariant
  classes", and invariant 1 (dangling edge) already reads "every context
  reference, tag key, tag value … must resolve to a declared model element" and
  "Terminal declarations are checked the same way — as tag predicates". A
  terminal context whose predicate reads a non-owned tag is exactly an
  unresolvable reference in the terminal position, and REQ-44 already routes
  the sibling defect (absent `[initial]`) to `graph-dangling-edge` by the
  disposition table. Minting a second code would also violate this record's own
  economy (`0006:C8`, "MUST NOT … mint a second code"). Chosen:
  `graph-dangling-edge`, with the message naming the terminal context id and
  the non-owned tag key. If a later phase finds a consumer that must branch on
  a distinct code, that is a Stage-8 deviation, not a REQ change here.

- ASSUMPTION-2: **The published product bound and the published node ceiling are
  two distinct implementation constants that share one code.** `0006:C12`
  publishes "the model-independent product bound"; `PE` requires "a
  model-independent node ceiling **alongside** the product bound". Both emit
  `graph-product-too-large` (CODES row plus the DISP row "Reachable node set
  over the published node ceiling"). The findings are distinguished by what
  they name: the product finding names "group + bound" (DISP), the ceiling
  finding names "the traversal (not a group)" and the ceiling (PE, SC-22). Both
  constants appear in the command's help output (REQ-59, REQ-60).

- ASSUMPTION-3: **`severity` is a two-valued field spelled `blocking` / `info`.**
  The CODES table's severity column uses exactly those two tokens, and the
  illustrative finding JSON carries `"severity": "blocking"`. No third tier
  exists — `0006:C8` and `0006:C17` both forbid one.

- ASSUMPTION-4: **The `reason` discriminator is carried only on
  `graph-unprovable-coverage`.** The REASON table is scoped to that code alone
  and no other code's row mentions a discriminator. `clierr.Finding` carries it
  as another declared string field on REQ-96's terms.

- ASSUMPTION-5: **"Recognized outcome" as a group key is the row's declared
  outcome binding, not a runtime value.** REQ-20 makes a group "rows sharing a
  match pattern and a recognized outcome". Since grouping is syntactic and
  per-row (REQ-22), the outcome half is read from the normalized row's declared
  outcome, matching RDR 0003's scoped-row-group definition.

- ASSUMPTION-6: **`--as=json` is the existing global output-mode flag, not a
  lint-local one.** The IC examples spell `intrastate lint --as=json`, and the
  shipped `respond`/`clierr` gateway already owns mode selection; REQ-5 forbids
  a second output path, so no new flag is introduced.

- ASSUMPTION-7: **The MVV's "same command shape intended for CI" means the
  fixture tests drive the root Cobra command through `ExecuteAndEmit`, not the
  engine API directly.** TS states the matrix must "exercise `intrastate lint`
  through the Cobra path and the same output gateway intended for CI", and
  REQ-119's CI job is a separate, A8-gated artifact. The fixture suite uses
  `--model <path>` per REQ-14.

- ASSUMPTION-8: **`graph-coverage-gap`'s payload follows RDR 0003's fenced
  requirement** — it names "the selection context, every rule id in the group,
  and one uncovered assignment" (`0003`, disposition table). This record fixes
  the code and severity but leaves the payload to the cited peer; REQ-71's
  general carriage clause is the floor, and 0003's row is the specific
  statement.

- ASSUMPTION-9: **A8's checked-in transition model is authored during this
  implementation**, per IP Prerequisites ("the model is authored and homed
  during implementation"). REQ-121, REQ-119, REQ-120, and REQ-122 all depend on
  it; if it cannot be authored within RDR 0006's surface, they become recorded
  deviations rather than silently dropped REQs — the MVV proper (REQ-MVV) is
  fixture-backed and does not consume the CI gate.

---

## QUESTIONS

None blocking.

Five clauses were examined as ambiguity candidates and each resolved against
the record's own text, the deviations file, the landed peers, or predecessor
precedent, so none is recorded as a blocking QUESTION:

1. **Does D1 check 2 mint a new finding code?** Two readings ("Add the code
   row" = a new code, vs. a new *row* documenting an existing code's trigger)
   would produce materially different taxonomies. Resolved against this
   record's own economy clause (`0006:C8` forbids minting a second code for a
   case an existing one names) and against invariant 1's scope, which already
   covers terminal declarations as references. Recorded as REQ-129 plus
   ASSUMPTION-1. Escalate only if a consumer needs to branch on a distinct
   code.

2. **Which node population does invariant 6 read — merged or split?**
   Invariant 6 is existential (REQ-110), so merged nodes are safe, but
   invariant 2's split rule (REQ-36) sits adjacent in the text. Resolved by
   LBD's explicit statement (REQ-112): the node form governs, and splitting is
   scoped to invariant 2's terminal-participating keys only (REQ-37).

3. **Is invariant 5's `graph-always-present-owned` a separate invariant or a
   sibling code?** The invariant set has seven members and this code is
   introduced inside invariant 5's prose as "`graph-single-valued-state`'s
   sibling code". Resolved: it is a distinct blocking code (its own CODES row
   and its own DISP row) fenced separately in `0006:C11`, but not an eighth
   mandatory invariant class — `0006:C3`'s list stays at seven. Recorded as
   REQ-67/REQ-68 without inflating the class count.

4. **Does `graph-terminal-escape` fire on a bare missing terminal
   declaration, or only on reliance-on-inference?** Invariant 7's text narrows
   it: "The input-observable defect this mints `graph-terminal-escape` for is a
   model that *relies* on such inference". Resolved: reliance is the trigger
   (REQ-43); a model omitting `[initial]` takes `graph-dangling-edge` instead
   (REQ-44), and RDR 0002's own table routes "missing `[initial]` / `terminal`"
   to lint without fixing a code, leaving this record's DISP row authoritative.

5. **Do A6/A9/A10's `Pending` statuses gate implementation?** Read literally
   they would defer REQ-101's root, REQ-35's stop set, and REQ-40's
   write-replaces premise. Resolved by D2 plus the landed peers: RDR 0002 and
   0003 are `Implemented`, `internal/table` carries `Model.Initial` and
   `Model.Terminal`, and 0002 fences write-replaces. The statuses are stale
   text, which is exactly what D2 records; no REQ is deferred on them.
