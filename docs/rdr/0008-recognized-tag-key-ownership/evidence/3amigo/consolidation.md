Model: claude-opus-5[1m]

# 3amigo — consolidation (RDR 0008, iteration 1)

Origin ledger for the resolve pass: PM-1…5, IMP-1…5, QA-1…5 (15 entries).
Three distinct persona lists; no finding was reused across personas.

## Hotspot passages — drew fire from two or more personas

Ranked by how many personas hit the passage and by severity.

### H-1 — Normative block 2's cardinality clause (`a flow MUST carry at most one such declaration`) — 3 personas

**IMP-5(a) · QA-3 · (PM-3 adjacent)**. All three arrive independently at the
same structural defect: the clause is **unfalsifiable in RDR 0002's locked wire
format**. `[tags.<tag>]` (`0002…md:268-274`, normative) is a TOML table keyed by
tag name, so two recognized-provenance declarations must carry two *different*
names — at least one of which already breaches the same sentence's naming clause
— while two declarations both named `recognized` are a TOML duplicate-key error
(`malformed TOML`). No input isolates cardinality. Scenario 5's fixture is
therefore unconstructible and would pass for the wrong reason.

Compounding, "a flow" is an undefined scope unit: 0002's schema unit is
`[model] id`, and `flow` appears there only as authoring prose. Cross-file
cardinality has no defined check site.

**Highest-priority rewrite.** This is a contract defect, not a test-detail gap.

### H-2 — Normative block 5 / Approach item 5 (the `RequiresOwned` reservation) — 2 personas

**IMP-1 · QA-3(adjacent, via scenario 9's lint half)**. `Row.RequiresOwned` has
**no authored TOML source form** — `requires_owned` is absent from 0002's locked
field layout and has zero occurrences repo-wide. It is a *derived* normalizer
output. RDR 0007 (`0007…md:2075`) narrows it to post-guard write-dependency keys,
i.e. derived from `[rule.write]` — in which case a `recognized` entry can only
arise from a write to `recognized`, which 0002 **already** rejects as
`write to non-owned tag` (`0002…md:345`). Block 5's lint half is then dead code
in the wrong category, and scenario 9's lint half has no authorable fixture.

Note this is the clause the **cove lens added last pass** (F-8). The kernel half
(hand-constructed `Row`) remains real and reachable; the *lint locus* is the
defect.

### H-3 — The open Enforcement locus (a)/(b)/(c) — 2 personas

**IMP-2 · QA-1**, with **PM-5** supplying the decision-framework angle. The RDR
defers the locus to Pre-Lock; Pre-Lock is now. IMP-2: three materially different
diffs ((b) adds the first non-nil error path `Resolve` has ever had; (c) invents
an exported name the RDR forbids borrowing from 0009). QA-1: scenario 6's Expected
holds **two mutually exclusive oracles** — (b)/(c) require `err != nil`, (a)
requires `err == nil` plus a refusal — so the test cannot be written in either
direction. PM-5 notes the QOC deciding rows (blast radius, prior art) do not rule
out (a), which A6's own Residual concedes leaves the shadowing "detected by
nothing."

### H-4 — Block 3's failure payload / the `reserved tag key` category surface — 2 personas

**IMP-3 · QA-4**. Block 3 makes payload content normative (offending name,
required name, the rule) but names no carrier: no Go type, package, field names,
or category spelling exists at HEAD (0002 defers even the package name). The
category is written as prose `reserved tag key` where sibling discriminators are
snake_case constants. "The rule" is prose, not a comparable value, so the Phase 3
golden-text test has no oracle. IMP-3 adds an ownership question QA did not:
block 3's second clause reaches *into* 0002's pre-existing `unknown tag` payload
and mutates it — an unowned edit.

## Single-persona findings carried to the resolve pass

Not hotspots, but each exits the ledger with a disposition:

- **PM-1** (High) — the chosen branch denies the naming freedom the Problem
  Statement frames as the user's want, and never says so; "the capability the
  user wanted" shifts meaning between Problem Statement and Alt 2's rejection.
- **PM-2** (High) — "authoring-time diagnosis" is asserted but routed through no
  named user-facing command; RDR 0005's code-string table (`0005…md:599-616`) has
  no entry for a 0002-layer data-level category.
- **PM-3** (High) — the innocent-name residual (`[tags.outcome]`, provenance
  `observed`) is plausibly the *majority* of the Problem Statement's population,
  not a "narrow path": both canonical 0002 fixtures name the recognized tag
  `outcome`.
- **PM-4** (Medium) — no criterion is written from the author's seat; the lint
  message's user-facing text is outside every test.
- **QA-2** (High) — scenario 3's "identical assembled view" has no exported
  observable (`TagSet` exports only `Lookup`/`Len`) and names no seam owner.
- **QA-5** (Medium) — scenario 8's "proves unreachable" is a universal negative
  with no stated method; the overwrite has no observable at the package boundary.
- **IMP-4** (Medium) — phase ordering: six of nine scenarios un-startable at HEAD
  and the MVV is defined end-to-end over a normalizer that does not exist,
  colliding with the Gate's "executed during implementation, not deferred."
- **IMP-5(b)** (Medium) — "post-parse" parse boundary: `[tags.*]` only, or also
  `[rule.match.<tag>]` / `guard.all` / `guard.unless` / `[rule.write]`?

## Cross-cutting note (both PM and QA raised a form of it)

PM's framing question — is the user's outcome *"I get to name this tag"* or
*"my row fires, or I'm told why"*? — and QA's Done-criterion note ("every scenario
below has a green test" cannot go green from this RDR's implementation alone)
are the same boundary problem seen from two ends: **what does this RDR finish, and
for whom?** PM-1/PM-3/PM-5 collapse once the first is answered; IMP-4 and QA's
note collapse once the second is.
