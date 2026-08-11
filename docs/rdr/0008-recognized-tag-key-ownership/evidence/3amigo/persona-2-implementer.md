Model: claude-opus-5[1m]
Persona: Implementer

# 3amigo — Persona 2: Implementer (RDR 0008)

Lens: I am handed this RDR on Monday morning and told to open an editor.
These are the five things I would have to ask before the first line of Go,
ranked by how hard they block me. Grounded against
`internal/resolve/resolve.go` (`recognizedTagKey`, `assemble`, `Resolve`,
`missingOwned`, `TagSet`), `internal/resolve/fixtures_test.go`, and peer
RDRs 0001, 0002, 0007, 0009 at HEAD.

---

## IMP-1 — `Row.RequiresOwned` has no authored source form, so the block-5 lint rule has nothing to inspect

**Severity**: High

**Passage that triggered it** — Technical Design / Normative Contracts,
block 5:

> "A normalized row's `Row.RequiresOwned` MUST NOT name `recognized`. …
> Enforcement is the same load/lint locus and the same `reserved tag key`
> category as the declaration channel."

and Proposed Solution / Approach, item 5:

> "An owned-key *reference* is a third way the name arrives: a row naming
> `recognized` there fails `missingOwned`'s owned-only provenance test and
> refuses `owned_state_unavailable` naming the reserved key. Load/lint
> rejects the name."

**Question**: What authored TOML construct does the normalizer inspect to
produce a `RequiresOwned` reserved-key failure? RDR 0002's locked wire
format is `root outcomes`, `[model]`, `[tags.<tag>]`, `[accessors.<id>]`,
`[context.<id>]`, `[[rule]]`, `[rule.match.<tag>]`,
`[rule.guard.all.<tag>]`, `[rule.guard.unless.<tag>]`, `[rule.write]`,
rule-level `clear`, rule-level `escape`, `[dump]` — there is **no**
`requires_owned` key, and a repo-wide sweep of `docs/rdr/` returns zero
occurrences of `requires_owned` / `requires-owned` in any spelling.
`Row.RequiresOwned` is therefore a *derived* normalizer output, not
something an author types. Compounding this, RDR 0007 is concurrently
narrowing the field's meaning to **post-guard write-dependency keys**
(`0007…md:2075` "Narrow `Row.RequiresOwned`'s doc contract to post-guard
…"), which implies derivation from `[rule.write]` tag names. So which is
it:

- (a) the normalizer derives `RequiresOwned` from `[rule.write]` keys, in
  which case a `recognized` entry can only arise from a *write* to
  `recognized` — which RDR 0002 already rejects as `write to non-owned
  tag`, making block 5 dead code and the `reserved tag key` category the
  wrong one; or
- (b) there is some other derivation path (match keys? declared owned
  tags?) I should be reading from, in which case name it; or
- (c) block 5 is really a post-normalization *assertion* over the emitted
  `[]resolve.Row`, not a source-level lint — in which case it lives at a
  different code site than blocks 2–3 and cannot share "the same load/lint
  locus."

**Decision it blocks**: I cannot write the Phase 2 check at all. I do not
know what function signature it has (does it take sparse TOML source
structs, or `[]resolve.Row`?), what it iterates over, or whether the
condition is even reachable. Testing Strategy scenario 9's lint half
("the same name in `RequiresOwned` fails in the `reserved tag key`
category before resolution") is unwritable for the same reason — I cannot
author the fixture that triggers it, because there is no TOML syntax that
puts a name into `RequiresOwned` directly. (Scenario 9's *kernel* half is
fine and runnable today: I can hand-construct
`resolve.Row{RequiresOwned: []string{"recognized"}}` in
`fixtures_test.go` and assert `owned_state_unavailable` — that half needs
no answer.)

---

## IMP-2 — The Enforcement locus is left open, and it is the single largest code-shape fork in the RDR

**Severity**: High

**Passage that triggered it** — Load-Bearing Decisions / Enforcement
locus:

> "**Enforcement locus (open — settle at Resolve)** — where the `Input`
> producer obligation is *checked*. … Three candidates: (a) an unchecked
> documented obligation on producers … (b) a kernel-entry precondition on
> `Resolve` returning the Go error RDR 0001 reserves … (c) an exported
> construction-time predicate over `Input` that producers call … Default
> lean: (b) … Pick the final form at Pre-Lock"

and Normative Contracts block 4:

> "The enforcement locus is open (Load-Bearing Decisions / Enforcement
> locus): it is NOT inherited from RDR 0009"

**Question**: Which of (a)/(b)/(c) — or the (b)+(c) pairing the passage
says "is available here" — am I building? These are three materially
different diffs to `internal/resolve/resolve.go`:

- (a) is a doc-comment edit on `Input` and nothing else — zero test, zero
  new symbol.
- (b) adds a non-nil error return path to `Resolve`'s body, which today
  has **none** (verified: every `return` in `Resolve`, `escapeOrRefuse`,
  and `gate` pairs a `Result` with a literal `nil` error). That is the
  first exercise of that channel in the package and changes what every
  existing test must assert about the second return value.
- (c) adds a new **exported** symbol to the kernel package — and the RDR
  simultaneously forbids me from naming it after 0009's illustrative
  `ValidateTable` ("do not bind to any symbol name from 0009"), so I would
  be inventing the public name with no guidance.

If it is the (b)+(c) pairing, I need to know that up front, because
0009's stated reason for pairing them is "one predicate, two call sites,
so the enforcement points cannot drift" — that is a shared-helper
structure I should build once, not retrofit.

**Decision it blocks**: Everything on the data channel. I cannot write the
Phase 3 producer-obligation breach test (Testing Strategy scenario 6),
whose Expected clause is itself branch-conditional — "caught at whichever
locus A6's Enforcement-locus decision settles … If the locus resolves to
(a) documented-only, this scenario instead pins the D3 backstop." Those
are two different tests asserting two different things (a non-nil `error`
vs. a `KindNoMatch` refusal). I also cannot size the change to the
existing suite: under (b), every test that currently does
`got, err := resolve.Resolve(in)` and ignores `err` becomes a candidate
for revision.

---

## IMP-3 — The `reserved tag key` category has no named Go type, package, or field shape, and block 3's payload is normative

**Severity**: High

**Passage that triggered it** — Normative Contracts, block 3:

> "Every `reserved tag key` failure — and any undeclared-tag failure whose
> offending key is the reserved name — MUST carry, at the data level, the
> offending name, the required name `recognized`, and the rule (the kernel
> owns this key; declarations conform to it). A category consumer may map
> it, but the guidance travels in the failure data, not the renderer."

**Question**: What Go type carries this, and where does it live? The RDR
says "data-level category" throughout, but no such type exists at HEAD:
RDR 0002's Existing Infrastructure Audit records "Resolver/table package
| `internal/` search | Not implemented | Introduce | New internal package
can own sparse source structs and normalized row structs" — i.e. 0002
defers even the *package name*. So I need:

- the package path (a new `internal/table`? `internal/normalize`?
  something under `internal/cli`?);
- whether the category is a string constant, a typed enum like
  `resolve.RefusalKind`, or a struct field;
- the exact spelling of the category value — the RDR writes it as prose
  `reserved tag key` with spaces, but RDR 0001's sibling discriminators
  are snake_case string constants (`no_match`, `owned_state_unavailable`,
  `unmodeled_outcome`), and RDR 0002's peer categories are likewise prose
  in the doc (`unknown tag`, `write to non-owned tag`). Is the on-the-wire
  value `reserved tag key`, `reserved_tag_key`, or `reserved-tag-key`?
- the three payload field names. "the offending name, the required name
  `recognized`, and the rule" is a normative *content* requirement with no
  field names attached, yet Phase 3 demands a "golden refusal-text check
  (offending name, required name, rule)" — a golden test asserts exact
  bytes, which I cannot write against unnamed fields.

Note this is not deferrable to 0002's implementer: block 3's second clause
("**and any undeclared-tag failure whose offending key is the reserved
name**") makes this RDR reach *into* 0002's pre-existing `unknown tag`
category and mutate its payload. Someone has to own that edit.

**Decision it blocks**: The Phase 2 failure construction and the Phase 3
golden-text test (Testing Strategy scenario 7). I also cannot tell whether
scenario 7's "category-enumerating consumer that does not know the new
category" is a real consumer I must find and test through (RDR 0005's exit
code map lives in `internal/cli/clierr`) or a synthetic stub.

---

## IMP-4 — Phase ordering is unrunnable as written: Phase 2 targets code that does not exist, and the Prerequisite gate is not a checkable condition

**Severity**: Medium

**Passage that triggered it** — Implementation Plan / Prerequisites and
Phase 2:

> "- [ ] RDR 0002 implementation underway (the enforcement point is its
> normalizer's load/lint path; this RDR's checks land inside that work,
> not before it)"

> "**Phase 2: Normalizer name validation** — Add the reserved-key checks
> to RDR 0002's load/lint path"

and Testing Strategy:

> "Scenarios 2–5 and 7 land inside RDR 0002's normalizer work (no
> normalizer exists at HEAD — reuse audit); 1, 3, and 6 are kernel-side
> and runnable against `internal/resolve` as it stands; 9 has a half on
> each side."

**Question**: What exactly am I expected to land on Monday, and what is
the definition of done for *this* RDR's implementation? Six of nine test
scenarios (2, 4, 5, 7, and the lint halves of 2 and 9) are explicitly
un-startable at HEAD. Concretely I need to know which of these is the
intent:

- (i) This RDR ships **Phase 1 + the kernel-side scenarios only** (1, 3,
  6, kernel half of 9), and Phases 2–3's lint half become a written
  obligation carried forward into 0002's implementation prompt — in which
  case this RDR's Status can flip to implemented with the MVV *not*
  executed, which contradicts the Finalization Gate's Scope Verification
  item ("will be executed during implementation, not deferred").
- (ii) This RDR **blocks** until 0002 is underway, in which case "underway"
  needs an operational definition (a merged package? a passing parse test?)
  — an unchecked checkbox is not something I can evaluate.
- (iii) This RDR implements a *slice* of 0002's normalizer itself (enough
  parse + tag-declaration handling to host the name check), which would be
  a large unscoped expansion nobody has costed.

The Minimum Viable Validation compounds this — it is defined as "One
end-to-end pair over the normalizer + kernel," so the MVV itself cannot
run at HEAD under (i).

**Decision it blocks**: My branch plan and my commit scope. I cannot tell
whether I am writing a ~30-line kernel-side change plus tests, or standing
up a TOML parser. It also determines whether I write Phase 3's fixtures in
`internal/resolve/fixtures_test.go` (where
`recognizedTagSensitiveTable` already lives) or in a new package's test
tree.

---

## IMP-5 — "at most one such declaration" and "post-parse key string" need their scoping unit and parse boundary named

**Severity**: Medium

**Passage that triggered it** — Normative Contracts, block 2:

> "A tag declaration with provenance `recognized` MUST be named
> `recognized`, and **a flow MUST carry at most one such declaration**. …
> The reserved-key comparison is byte-exact on the **post-parse** key
> string: case-sensitive, no trimming, no folding … TOML quoting is a
> surface artifact, not a name variant: `[tags."recognized"]` and
> `[tags.recognized]` parse to the identical key string and are therefore
> both reserved."

**Question**: Two sub-questions, both about where in the pipeline the
check sits.

(a) **Scoping unit.** What is "a flow" as a value I can iterate? RDR 0002
says "Each flow has a named model with **flow-level** tag declarations"
and its Identity decision is `(model id, rule id)` — so declarations live
under one `[model]`, and cardinality is per-model. But `[tags.<tag>]` is a
TOML *table* keyed by tag name, which already makes duplicate keys a TOML
parse error, not a semantic one. So how does "two recognized-provenance
declarations" (Testing Strategy scenario 5) even arise? It can only be two
*differently named* tags both carrying `provenance = "recognized"` — e.g.
`[tags.outcome]` and `[tags.result]` — in which case the failure is
already covered by the block-2 naming rule (both are wrong-named) and the
cardinality clause is redundant. Or does it mean something across multiple
`[model]` blocks / multiple files, which would need a cross-file check
site? I need the scenario-5 fixture spelled out, because the two readings
produce different failure counts (two failures vs. one).

(b) **Parse boundary.** "post-parse key string" — post *whose* parse? For
`[tags.<tag>]` the answer is clear (the TOML decoder's map key). But the
same reserved word must also be checked in `[rule.match.<tag>]`,
`[rule.guard.all.<tag>]`, `[rule.guard.unless.<tag>]`, and `[rule.write]`
per RDR 0002's field layout, since Technical Design says the normalizer
"emits normalized rows whose `Match`/`Guard` references to the recognized
outcome use the reserved key." Are the reserved-key checks applied at
every one of those sites, or only at `[tags.*]` with the others covered
transitively by 0002's existing `unknown tag` rule? Note the two produce
different diagnostics: a `Match` on `recognized` where no
`recognized`-provenance tag is declared currently yields `unknown tag`,
not `reserved tag key`, and block 3's carve-out ("any undeclared-tag
failure whose offending key is the reserved name") suggests the RDR
already knows this but does not say which check fires first.

**Decision it blocks**: The check's call site(s) and its cardinality
bookkeeping — whether I write one pass over the `[tags.*]` map or a
walker over every tag-key position in the normalized rule set. It also
blocks Testing Strategy scenarios 4 and 5: scenario 4's fixtures
(`Recognized`, `RECOGNIZED`, `[tags."recognized"]`, `[tags." recognized"]`)
need a decided call site to assert against, and scenario 5 has no
constructible fixture until (a) is answered.

---

## Not raised (in scope, unambiguous, buildable today)

For the record, so the list above reads as targeted rather than
exhaustive — these I can write on Monday with no further input:

- **Phase 1** in full: the pointer comment on
  `internal/resolve/resolve.go::recognizedTagKey` citing RDR 0008. The
  constant is already `"recognized"` and already conforming; no behavior
  change.
- **Testing Strategy scenario 3** (A2, guard/matcher same-view): I need a
  guard seam that inspects its `TagSet`. `fixtureGuards.Evaluate` at
  `internal/resolve/fixtures_test.go:20` currently has signature
  `Evaluate(guard string, _ resolve.TagSet)` and discards the view; I add
  a new fixture type (rather than change that one, which the existing
  suite depends on) that reads `view.Lookup("recognized")` and asserts the
  matcher and guard agree. `TagSet.Lookup` is already exported and
  sufficient.
- **Testing Strategy scenario 9, kernel half**: hand-construct
  `resolve.Row{RequiresOwned: []string{"recognized"}}` and assert
  `KindOwnedStateUnavailable` with `MissingOwned: ["recognized"]`.
  Confirmed reachable by reading — `missingOwned` (`resolve.go:444`) tests
  `view.has(key, ProvenanceOwned)` and `TagSet.has` (`:125`) compares
  `tv.provenance == prov`, so a key present under `ProvenanceRecognized`
  fails the owned-only test.
- **Testing Strategy scenario 8, backstop half**: the non-conforming
  hand-constructed case is directly expressible against `assemble`'s
  existing `Observed`-then-`Owned` loop order (`resolve.go:157-162`).
