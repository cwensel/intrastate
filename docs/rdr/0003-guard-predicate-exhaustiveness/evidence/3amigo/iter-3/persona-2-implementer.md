Model: claude-opus-5[1m]

# Persona 2 — Implementer, iteration 3 (re-entry, delta-scoped)

Scope: the tag declaration model rehomed here 2026-08-21 (value kind, finite
domain, optionality, single-valuedness, element universe) and the 0002 split;
A8/§JD-4; §JD-13's single-valued marker; §JD-14's escape-row reading; A5
re-verification; the open A10/A12/A14/A15. Question asked of each passage: *if I
started coding Monday, what would I ask in the first hour?*

Grounded against `docs/rdr/0002-transition-table-as-reviewable-data.md`,
`docs/rdr/0006-graph-lint-authority-and-guarantees.md`,
`docs/rdr/0007-guard-predicate-totality.md`,
`docs/jdr/0001-resolve-kernel-seam.md` (§JD-4, §JD-13, §JD-14),
`internal/resolve/resolve.go`, and
`docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/guard-fixture.toml`.

---

## Severity 1 — blocks writing the first line of Phase 1/2 code

### P2-1 — `Normative Contracts`, the optionality clause, vs. the "can refuse" clause

**Anchor:** `Normative Contracts` → "A tag declaration MUST be able to state
whether the key may be absent…" (0003:704-712) read against `Normative
Contracts` → "A row 'can refuse' when the row carries a value atom whose key is
not declared **present-for-every-reachable-predecessor**…" (0003:859-865).

**Request.** The optionality clause gives a declaration exactly two states —
*always-present* and *optional*. The "can refuse" clause tests a **third,
differently-named** property: `present-for-every-reachable-predecessor`. Which
is it? Three readings are live and they give different MVV Scenario 6 verdicts:

1. `present-for-every-reachable-predecessor` is a synonym for the optionality
   clause's *always-present* — in which case one of the two spellings should go,
   because an implementer grepping for a field name finds two.
2. It is a strictly stronger property that additionally quantifies over the
   predecessor relation — in which case it is **not** decidable from the tag
   declaration alone, it needs A12's (still `Pending`) predecessor contract, and
   the clause's own promise that lint "MUST decide this syntactically over
   declarations so the test is total" is unmet. A12's Stage 6 disposition
   asserts the opposite ("no clause in this RDR now claims a decision procedure
   it cannot cite"), so this reading contradicts A12's survivability argument.
3. It is a third declarable field nobody has declared — it appears in no
   declaration clause, in the `authority` census's "Per-tag optionality" row, or
   in RDR 0002's authoring list.

**Decision blocked.** Whether the declaration model needs one presence field or
two, and whether the "can refuse" predicate is a pure declaration lookup
(codeable today, Phase 2 "startable: yes for the product/proof core") or a graph
query gated on A12 (not startable). This is the exact predicate MVV Scenario 6
asserts positively and its negative control asserts green on, so the test cannot
be written either.

### P2-2 — `Normative Contracts`, the optionality clause: no stated default

**Anchor:** `Normative Contracts` → the optionality clause (0003:704-712); the
declaration-model clause says a declaration "**MAY** carry … an optionality
marker" (0003:686-687).

**Request.** What does a tag that carries *no* optionality marker mean —
always-present, optional, or a declaration error? The RDR states no default
anywhere (grep for `default` in 0003 returns only "not by default" at :1477,
about phase sequencing). The fixture at `evidence/spikes/guard-fixture.toml`
declares seven tags and **not one** carries an optionality marker, so under
"default = optional" every one of its groups is withheld and under "default =
always-present" every one certifies green — opposite verdicts on this RDR's own
representative fixture.

The choice is not symmetric and the RDR has an argued position on one side of
it: the "Leaving a dimension undeclared is not an opt-out" paragraph
(`Technical Design`, 0003:637-641) establishes that silence about a *domain*
takes the loud blocking outcome, never a silent pass. Silence about
*presence* is the same shape of question and gets no answer.

**Decision blocked.** The presence dimension A7's projection clause builds for
every key in every scoped product; whether the MVV's "negative control" group
(0003 MVV: "a row group whose keys are all declared always-present") requires
new fixture syntax or is the fixture's existing default; and whether the
declaration parser rejects, defaults, or infers.

### P2-3 — `Normative Contracts`, the participation clause, vs. the `trace` table step 3

**Anchor:** `Normative Contracts` → "A dimension **participates** in a row group
when any row in that group carries an atom over that key, **in either `all` or
`unless`**" (0003:777-782), read against `trace` table step 3 (0003:975): "the
group's shared `status eq "Draft"` atom does not discriminate but still requires
a declared domain".

**Request.** Do `match`-block predicates enter the scoped product, or only
`all`/`unless` guard atoms? The participation clause says guard atoms only. The
trace's own witness contradicts it: in `evidence/spikes/guard-fixture.toml`,
`status eq "Draft"` is authored under `[rule.match.status]`, **not** under
`rule.guard.all` — and the `trace` table's step 2 names it as `match status eq
Draft` while step 3 counts it as a participating product dimension. RDR 0002
(`:228-231`) and the shipped kernel (`internal/resolve/resolve.go::Row.Match`,
a `[]Tag` separate from `Row.Guard`) both treat match and guard as distinct
fields, so this is not a naming slip an implementer can shrug off.

The stakes are concrete: step 3 says the participation clause "is what puts
`status` in the product — under a discriminates-only reading it would silently
drop out". If match predicates are outside the product, `status` drops out
anyway and the trace's justification for the participation clause evaporates
along with the clause's worked example.

**Decision blocked.** What the scoped product is built from — the single most
load-bearing input to A2's coverage/overlap derivation and to Phase 2's product
construction. It also decides whether match-block keys need the declaration
model's finite domain (and therefore whether an undeclared match key triggers
the blocking inability-to-prove outcome).

### P2-4 — `Normative Contracts`, the escape-row clause (§JD-14), vs. RDR 0002's exact-one clause and the shipped kernel

**Anchor:** `Normative Contracts` → "A declared escape row participates in the
coverage identity like any other row… An escape row carrying no guard atoms
denotes the whole scoped product and therefore closes coverage by itself…
MUST NOT exclude escape rows from overlap checks" (0003:799-807); Metadata
Status "§JD-14 confirmed this RDR's escape-row reading governs, no edit owed"
(0003:15-16).

**Request.** §JD-14 resolved this RDR against **RDR 0006**. It did not touch
RDR 0002 or the kernel, and both of those partition escape rows *out* of the
candidate set:

- RDR 0002 `Normative Contracts` (`:328-334`): "Ordinary transition success
  requires exactly one matching **non-escape** normalized candidate row. If zero
  or multiple non-escape rows match, the resolver MAY return a modeled escape
  disposition only when exactly one escape row matches…"
- The shipped kernel, `internal/resolve/resolve.go` `Resolve` doc comment step 2
  (`:304`): "candidate rows are the **non-escape** rows for that outcome"; escape
  rows are reached only via `escapeOrRefuse` (`:460-500`).

So at runtime an escape row that matches the same assignment as a guarded row is
**not** an ambiguity — the guarded row wins outright and the escape row is never
consulted. This RDR's clause calls that same pair "the ambiguity RDR 0001
refuses at runtime" and requires a blocking overlap finding. Two follow-ons:

(a) An escape row with no guard atoms "closes coverage by itself" — so under
this clause **any** group containing a bare escape row is exhaustive, and also
overlaps every guarded row in the group, so the same fixture is simultaneously
green on coverage and blocking on overlap. Is that the intended pairing, or does
the escape row's own `match` block scope it (which loops back to P2-3)?

(b) Does the row group defined in `Technical Design` (0003:617-620 — "rows
sharing one selection context") include escape rows at all? The definition
grounds itself on "exactly the set RDR 0001 resolves exact-one over" — and per
the kernel, that set is the **non-escape** rows.

**Decision blocked.** Whether the lint's group construction filters escape rows,
and whether a coverage-gap finding can be discharged by an escape row.
Directly gates MVV Scenario 2 ("one complete partition, one intentional gap, one
intentional overlap") — a fixture built under either reading fails the other.

---

## Severity 2 — blocks a specific field or clause, workaroundable for a sprint

### P2-5 — `Normative Contracts`, the single-valued clause (§JD-13), vs. RDR 0006 invariant 5

**Anchor:** `Normative Contracts` → "A tag declaration MUST be able to state
that the tag is **single-valued**: at most one of its declared domain values
holds in any conforming evaluation view… RDR 0006's grouping-dependent lint
findings read this field" (0003:714-724).

**Request.** RDR 0006's invariant 5 — the consumer §JD-13 names — reads
"**writes** must not produce two values for a **tag class** that the model
declares single-valued" (`docs/rdr/0006-...md:295-296`). Two mismatches an
implementer hits immediately:

- **Quantity.** This RDR's marker is a property of *one tag* ("at most one of
  **its** declared domain values holds"). RDR 0006's invariant is over a *tag
  class* — plural tags sharing a class, "such as one lifecycle/status value". A
  per-tag boolean does not express "these three tags are one class"; a class
  needs a grouping key. Which shape does the declaration model actually
  provide?
- **Site.** This clause constrains the *evaluation view* (what holds at match
  time). Invariant 5 constrains *writes* (what a row produces). A per-tag
  marker over views does not by itself decide whether a row's write block
  violates the invariant.

**Decision blocked.** The declaration field's type and the parser that reads it
(`bool` on the tag vs. a class identifier), and whether
`graph-single-valued-state` — a **mandatory blocking** code in RDR 0006's table
(`:318`) and asserted by its MVV (`:643`) — has the producer §JD-13 believed it
gained. MVV Scenario 4 requires rejecting a single-valued marker "on a `set`
kind", which is unwritable until the field's shape is fixed.

### P2-6 — `Normative Contracts`, the declaration-model clause: `opaque scalar` has no spelling and no rules

**Anchor:** `Normative Contracts` → "The value kinds are `enum`, `bool`, `int`,
`set`, and opaque scalar" (0003:688-689).

**Request.** Four kinds are back-ticked tokens; the fifth is unquoted English
prose. Three questions:

1. What is its authored token — `scalar`? `opaque`? `string`? The operator
   matrix in `Approach` spells the same thing a *third* way: "**string-like
   scalar**" (0003:573-574).
2. The finite-domain clause (0003:694-701) enumerates a finite domain for
   `enum`, `bool`, `int`, and `set` — and says nothing about opaque scalar. Can
   an opaque scalar declare a finite domain? If it can, it is an `enum` by
   another name; if it cannot, `eq`/`in` over it are permanently
   inability-to-prove.
3. The domain/kind agreement clause (0003:727-734) rejects three specific
   mismatches. None of them names opaque scalar, so the rejection table has a
   hole exactly where the least-typed kind sits.

**Decision blocked.** The value-kind enum in the declaration parser, and
whether the `eq`/`in` matrix rows accepting "string-like scalar" name a fifth
kind or are a stale spelling of `enum`. MVV Scenario 4's declaration-error cases
need the closed kind set to enumerate against.

### P2-7 — `Capability Dependencies`, the tag-declaration-model row, vs. RDR 0002's authoring clause

**Anchor:** `Capability Dependencies` → "**Tag declaration model** (value kind,
finite domain, optionality, set-element universe) | **This RDR** |
**Introduced**" (0003:1055) and "Authoring location + normalization carriage for
tag declarations | RDR 0002 | Pending" (0003:1058).

**Request.** The single-valued marker §JD-13 homed here on 2026-08-22 has **no
authoring location**. RDR 0002's normative clause that receives this model
(`docs/rdr/0002-...md:342-350`) enumerates exactly four fields — "value kind,
and optionally a finite domain, an optionality marker, and a set-element
universe" — and `single-valued` occurs **zero** times in RDR 0002 (verified by
grep). This RDR's own `Capability Dependencies` row (:1055) also lists only the
four, while the `authority` census row (:937) lists five. So the field this RDR
declares normatively has no place an author can write it, and no carriage
guarantee, and the RDR's two internal tables disagree on whether it is in the
model.

**Decision blocked.** Whether the fixture/TOML schema gains a `single_valued`
key under `[tags.<tag>]` now or whether this is a fifth request booked against
RDR 0002 alongside A14. §JD-13 says "RDR 0006 cites it at its refine" but names
no producer for the authoring surface, so today the field is declarable in
theory and unauthorable in practice — which is the precise defect the
2026-08-21 rehoming was performed to close.

### P2-8 — `Normative Contracts`, the finite-domain clause: no `int` bound spelling, and the fixture uses a different one

**Anchor:** `Normative Contracts` → "an `int` declares a `{min..max}` bound"
(0003:696-697); the same `{min..max}` token in the domain/kind agreement clause
(0003:728).

**Request.** `{min..max}` is used as if it were a wire spelling, but the RDR's
own `evidence/spikes/guard-fixture.toml` declares the bounded integer as two
separate keys:

```toml
[tags.prelock_iterations]
kind = "int"
min = 0
max = 3
```

Is `{min..max}` a notation for "a bounded range, however RDR 0002 spells it", or
a literal authored form? A13 set the precedent that a *literal's* canonical
spelling is this RDR's to fix (and this RDR discharged it for set literals),
which makes the silence on the domain spelling read as deliberate — but the
domain/kind agreement clause requires *rejecting* a `{min..max}` bound on an
`enum`, and a rejection rule needs to know what shape it is detecting.

Related: are the bounds inclusive at both ends? A2's derivation computes
cardinality from "every participating dimension's declared domain size", and
`{0..3}` is 4 values inclusive, 3 exclusive — a factor the too-large-bound
diagnostic (0003:885-901) reports as a number.

**Decision blocked.** The declaration parser's int-domain branch, the
domain/kind rejection rule, and the cardinality arithmetic A15's MVV Scenario 3
compares two products with.

---

## Severity 3 — worth an answer before the phase starts, not before the sprint

### P2-9 — `Normative Contracts`, the domain/kind agreement clause: literal-outside-domain timing

**Anchor:** `Normative Contracts` → "A value appearing in a guard literal that
lies outside its tag's declared domain MUST likewise be rejected **before
resolution**" (0003:731-734), in the same clause whose first sentence rejects
declaration errors "**before normalization completes**".

**Request.** Two different phase boundaries in one clause. Is the
literal-outside-domain check a normalization-time error (same gate as the
declaration errors) or a later pre-resolution gate? It matters because
normalization is RDR 0002's stage and resolution is the kernel's, so the two
readings put the check in two different packages with two different owners.
MVV Scenario 4 asserts declaration errors "are rejected **before normalization
completes**" but says only "before resolution" for the atom errors, preserving
the ambiguity into the test.

**Decision blocked.** Which package the domain-membership check lives in and
which error surface reports it (this RDR's predicate semantic kinds vs. RDR
0002's structural validation categories, which A11's rationale notes has "no
type or value-domain category").

### P2-10 — `Minimum Viable Validation` / `Phase gating` table, Phase 2 row

**Anchor:** `Phase gating` table → "2 Finite-Domain Lint Semantics | A10, A12
(RDR 0006 agreements) — A11, A2 and A8 are closed… | **Yes** for the
product/proof core; the peer agreements gate integration, not construction"
(0003:1470).

**Request.** The row says the product/proof core is startable today, but A12's
"can refuse" predicate (P2-1) and the row-group escape-row question (P2-4) are
both *construction* inputs to the product/proof core, not integration concerns.
Concretely: can Phase 2 build the group-construction function and the withheld-
claim decision, or only the coverage/overlap arithmetic over an
already-constructed product? The RDR's own A12 record says the reachability
quantifier "was removed from assertion and booked here" — but the "can refuse"
clause (0003:859-865) still spells `present-for-every-reachable-predecessor`
in a MUST.

**Decision blocked.** Sprint planning for Phase 2 — specifically whether the
group-construction and withholding logic are in scope before RDR 0006's refine
lands, or whether Phase 2 ships only the arithmetic and stubs its inputs.

### P2-11 — `Technical Design`, the "claims closed coverage" default-on reading, vs. A10 `Pending`

**Anchor:** `Technical Design` → "this RDR reads it as **default-on for every
scoped row group whose participating dimensions are all finitely declared**…
If RDR 0006 intends an explicit per-group annotation instead, that is a
divergence for its refine pass to settle" (0003:628-634); §JD-14 (0003:15-16)
confirms default-on.

**Request.** §JD-14 settled *default-on*, and the Status line records it as
"no edit owed". But A10 — whether RDR 0006 reads the same **row-group division
of labour** — is still `Pending`, and RDR 0006's Status line (`:13-15`) lists
"RDR 0003's A10/A12 agreements (row-group division…)" as owed at its refine. So
default-on is decided while *what it defaults on over* is not. If RDR 0006's
refine returns a different grouping predicate, every group's default-on verdict
changes with it.

Question: is the group-construction function I write in Phase 2 expected to be
stable across RDR 0006's refine, or should it be behind a seam I expect to
change?

**Decision blocked.** Whether group construction is a fixed internal contract or
a provisional one — which decides how much MVV fixture investment (Scenario 2's
gap and overlap both live inside one group) is safe before RDR 0006 refines.
