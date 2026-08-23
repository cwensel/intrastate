Model: claude-opus-5[1m]

# Critique Dispositions — iter-2 (dual-model, cross-model diff)

Origin ledger: `critique.md` (opus, C-1..C-15) + `critique-modelB.md` (sonnet,
C-1..C-10), reconciled by passage anchor in `diff.md` as M-1..M-22.

## Fixed

- **fixed** — M-3 inheritance merge silently drops constraints: defined the
  combined predicate set as a set over the **full** atom identity
  `(key, block, operator, literal)`, forbade keying the merge on any proper
  prefix, and stated that inheritance accumulates rather than overrides.
  Grounded in the spike: `main.go:440` keys on `block\x00key\x00op`, so a
  differing literal is a last-write-wins loss whose survivor depends on map
  iteration over `use` — and the C-002 identity sort runs after the merge, so
  ordering determinism structurally cannot detect it. Sections: Normative
  Contracts (inheritance); Testing Strategy sc.2.
- **fixed** — M-18 atom in both `all` and `unless` had no set semantics: block is
  part of atom identity, so both survive as distinct atoms; the rule is dead and
  reported by RDR 0006, never silently pruned at load (pruning would convert an
  authoring mistake into a matching rule). Sections: Normative Contracts (guard).
- **fixed** — M-2 "combined predicate set" defined twice with different extents:
  outcome binding now reads **match blocks only**; a `recognized` atom under
  `guard.all`/`guard.unless` is refused as a malformed outcome binding rather
  than lifted. Grounded: `main.go::liftOutcome` matches on `a.Key` with no block
  filter, so an `unless` binding inverts author intent silently. Sections:
  Normative Contracts (outcome binding); load categories; Testing Strategy sc.3.
- **fixed** — M-9 literal identity unspecified while key identity is byte-exact:
  a set literal normalizes to a member **sequence**, never a delimiter-joined
  string. Grounded: `main.go::renderLiteral` space-joins, collapsing
  `["needs work"]` and `["needs","work"]` into one atom. Sections: Normative
  Contracts (identity); Testing Strategy sc.2.
- **fixed** — M-14 identity tuple not total under `#`-bearing rule ids: reserved
  the expansion-suffix separator in both rule ids and alphabet members, with two
  new load categories. Closes the identity-level collision the Round-Trip section
  had conceded only as a rendering caveat. Sections: Normative Contracts
  (identity); load categories; Testing Strategy sc.3.
- **fixed** — M-4 (grounded half) version gate ordered only "before
  normalization": now required **before strict field validation**, via a
  two-pass load, so a v2 file refuses as unsupported version rather than as
  `unknown schema field` on an arbitrary v2-only key. Sections: Normative
  Contracts (version).
- **fixed** — M-20 strict decode stated as library trivia: restated as a
  normative obligation on the format, library-agnostic, so a parser swap cannot
  silently retire the `unknown schema field` category. Sections: Normative
  Contracts (version).
- **fixed** — M-11 (hotspot, both models) escape rescue is per-outcome: stated at
  the escape clause where an author looks, not only inside A8's evidence.
  Grounded on `main`: `::escapeOrRefuse` filters `row.Outcome != in.Recognized`
  before gating. Also resolves diff contradiction X-1 — an escape rule MAY bind
  with `in` and expand, so neither model's "opposite fix" was a design fork; the
  kernel had already decided it. Sections: Normative Contracts (escape).
- **fixed** — M-1 the atom-set → kernel two-field split was unspecified (the
  normalizer's central step): equality atoms populate `Match`; set membership,
  comparison, existence, and every `unless` atom belong to the guard predicate;
  an atom never appears in both. Scoped so RDR 0007's reshape still owns the
  guard's *shape*. Sections: Normative Contracts (row population).
- **fixed** — M-8 `NextTags`/`Writes` contractually forced identical: stated that
  they are equal **under this RDR's authoring surface** and why they stay
  distinct fields, and forbade populating one by aliasing the other so a future
  divergence is a compile-time change, not a silent one. Sections: Normative
  Contracts (row population).
- **fixed** — M-7 round-trip excluded the locator's only varying part, so a
  locator-dropping normalizer passed: presence and the rule-identifying part are
  now compared; only positional detail is exempt. Sections: Round-Trip.
- **fixed** — M-13 `[dump]` render preferences sit in the semantic source file:
  declared presentation-only and excluded from the normalized value, with golden
  tests steered to the value rather than rendered text. Sections: Normative
  Contracts (dump).
- **fixed** — M-19 escape-row enforcement had only kernel-side coverage: added a
  distinct malformed-escape control for an escape rule carrying a **clear list**
  (sc.3 previously listed only the empty-write-block shape) — the normalizer is
  the sole enforcement point, so each shape owes its own control. Sections:
  Testing Strategy sc.3.
- **fixed** — M-10 MVV's two-sibling requirement contradicted the
  promote-verbatim mandate: the MVV requirement wins; promotion is
  extend-then-promote, and the "canonical examples" clause fixes field layout and
  idiom, not row census. Grounded: the RDR fixture binds `round-clean` on one
  ordinary and one escape rule, so no two ordinary candidates contend — matching
  the desk trace's own "no witness possible" row. Sections: MVV.
- **fixed** — M-21 MVV scenario 4 depended on Pending RDR 0006: the overlap
  assertion is explicitly deferred to the Phase 5 lint handshake rather than
  counted satisfied, with this RDR's own obligation (rows carrying enough
  structure for the check) named. Sections: MVV.
- **fixed** — M-15 (diff-refuted as a design defect, real as a stale note): the
  trace claimed the single-member `in` case unwitnessed "because the spike
  suffixes on operator, not expansion count." The spike suffixes on
  `len(outcomes) > 1` (`main.go:282`), exactly matching the contract. Corrected
  the note rather than the design — an incorrect self-assessment was being
  carried toward Stage 6. Sections: Conditional Mini-Checks (trace).
- **fixed** — M-5 (hotspot, both models) A8's "sole barrier" overclaimed: the ban
  is a load-time check in a package the kernel does not import, so it covers only
  tables this RDR's loader built. Recorded the bounded reach and named the
  uncovered producers rather than implying table-wide coverage. Sections:
  Critical Assumptions (A8).

## Dismissed with cite

- **dismissed-with-cite** — M-17 cross-model dump ordering unspecified. The draft
  already decides it: "A dump spanning more than one model MUST carry a
  model-unique `model id`, which is what makes the leading tuple field
  discriminating; two models sharing an id is a `duplicate model id` load
  failure. RDR 0006 lints exactly one model per invocation."
- **dismissed-with-cite** — M-22 Prerequisites/RDR-0007 sequencing blocks lock.
  The draft says verbatim: "This item gates implementation sequencing, not lock."
  The interim shim is also named (Testing Strategy sc.6: "the spike mirrors the
  constants locally because the kernel does not yet export them").
- **dismissed-with-cite** — M-4 (migration half) version gate lacks a
  migration/coexistence path. Net-new scope for a v1-only format RDR: the gate's
  job is to refuse unambiguously, and a migration story presupposes a v2 whose
  shape no RDR has defined. The ordering defect inside this row was fixed above.
- **dismissed-with-cite** — B:C-8 ambiguous overlap should be a load failure.
  Re-raise of a Load-Bearing Decision: "Load and lint split on **arity, not
  severity**"; Testing Strategy sc.3 states "Ambiguous overlap is deliberately
  absent from this scenario: it is cross-row and therefore an RDR 0006 lint
  finding (scenario 4), not a load failure."

## Charted to successors

Recorded in `Charted.md`: M-6 (guard-read keys in `RequiresOwned` → RDR 0007),
M-12 (category → CLI code mapping → RDR 0005), M-5 residual (empty
`Input.Recognized` at the caller edge → RDR 0005 / RDR 0001), M-16 (re-readable
dump grammar → a dump-format RDR).

## Needs verification (Stage 6)

New or amended load-bearing claims, all **Pending** — none verified in this pass:

1. **Atom-set cardinality under inheritance** — that two contexts contributing
   the same key+operator with different literals yield two atoms, order-
   independently. Method: MVV test (Testing Strategy sc.2, both `use` orders).
   The spike currently exhibits the *defect*, so this is a claim about the
   production normalizer, not witnessed by existing evidence.
2. **Set literals as member sequences** — that `["needs work"]` and
   `["needs","work"]` remain distinct atoms. Method: MVV test (sc.2). The spike
   space-joins, so existing evidence contradicts the new contract by design.
3. **Outcome binding reads match blocks only** — that a `recognized` atom in a
   guard block refuses as a malformed outcome binding. Method: MVV test (sc.3).
   Unwitnessed: the spike lifts without a block filter.
4. **`#` reserved in rule ids and alphabet members** — two new load categories.
   Method: MVV test (sc.3). Unwitnessed.
5. **Version check precedes strict decode** — that a v2 file refuses as
   unsupported version, not `unknown schema field`. Method: MVV test (sc.3).
   The spike's ordering is unverified against this requirement.
6. **Atom routing across `Match` / `Guard`** — that equality atoms populate
   `Match` and all others the guard. Method: MVV test. Interacts with RDR 0007's
   reshape (Prerequisites), so it may only be assertable after that lands.
7. **Escape rules expand under `in`** — that an escape rule binding with `in`
   yields one rescuing row per outcome. Method: MVV test. Unwitnessed: the
   fixture's escape rule binds with `eq`.
8. **Extended fixture carries two ordinary sibling rows** — required by the MVV
   and absent from the current fixture. Method: Spike/fixture extension.

## Tiebreakers

None. The one genuine fork the diff surfaced (X-1: whether an escape rule may
expand under `in`) was collapsed against the shipped kernel rather than escalated
— `::escapeOrRefuse` filters escapes on `row.Outcome`, which settles it as
documentation, not design.
