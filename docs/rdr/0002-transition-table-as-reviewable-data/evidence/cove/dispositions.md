Model: claude-opus-5[1m]

# Cove Dispositions — RDR 0002, iteration 1

Origin ledger: `findings.md` CV-001 … CV-021. One line per finding.

## Fixed

- **fixed** — CV-001 (A5 parse-code precedent refuted): rewrote A5's Evidence to
  state what source shows — the envelope is verified, but no parse-validation
  code has ever travelled it; `config-invalid` is a doc comment and a `TODO`.
  Named the table categories as the first of their kind, so code names, `Group`,
  and exit mapping are new design work. Status stays `Verified` — the assumption
  (failures *can* be surfaced through the envelope) is unaffected; only the
  precedent clause was false. Sections: Critical Assumptions A5; Existing
  Infrastructure Audit (two rows swept).
- **fixed** — CV-002 (Prerequisites understate the missing kernel surface):
  widened the bullet to the whole atom surface — no `Atom` type, no operator
  token, no block field, `Row.Match` a flat `[]Tag` — and stated that only two
  of Phase 2's four deliverables have fields to land in. Section: Implementation
  Plan / Prerequisites.
- **fixed** — CV-003 (exact-byte key identity has an uncited in-repo sibling):
  added `::TagSet.matches` / `::Lookup` / `::has` and the `ModeOf` counter-example
  as the confirming sibling beside the RDR 0007 A22 cite. Section: Normative
  Contracts / tag-key identity.
- **fixed** — CV-004 (load vs lint routing contradicted across three sites):
  introduced the **arity split** — single-rule checks are load-time and this
  RDR's; cross-row checks are RDR 0006 lint — and applied it to the Technical
  Design validation sentence, the validation-category normative block, and
  Failure Modes. Ambiguous overlap now routes to lint everywhere. Sections:
  Technical Design; Normative Contracts; Failure Modes.
- **fixed** — CV-005 (row kind required here, forbidden at the kernel boundary
  by Final RDR 0009): row kind is now a **derived view property, never a row
  field**, computed from the escape class list; removed it from the normalized
  row's field list and from both dump field lists. Sections: Technical Design;
  Normative Contracts (escape + dump); Round-Trip; Testing Strategy 2.
- **fixed** — CV-006 (`[dump]` field selection vs the "every field" MUST):
  `[dump]` MAY reorder columns, MUST NOT omit a field; a dump missing a field
  is not an expanded table dump. Section: Normative Contracts / dump.
- **fixed** — CV-007 (cyclic context inheritance had no clause and no category):
  added `cyclic context inheritance` to the validation-category block and to the
  Technical Design and Failure Modes lists. The spike already implements the
  check (`main.go::mergeContext`); the contract now names it.
- **fixed** — CV-008 (single-member `in` minted a different identity than `eq`):
  pinned that a suffix attaches **only where a rule expands** — a rule binding
  one outcome normalizes to an unsuffixed row however spelled. Cited RDR 0003's
  repeated-element rejection for the duplicate case. Section: Normative
  Contracts / outcome binding.
- **fixed** — CV-009 (alphabet rules had one undifferentiated category): split
  into missing alphabet and malformed alphabet with the three sub-cases (empty,
  duplicate member, empty-string member) each distinguishable. Section:
  Normative Contracts / validation categories.
- **fixed** — CV-010 (normalization explosion routed two ways, no threshold):
  pinned as a lint **diagnostic** reported per rule, explicitly not a refusal,
  and stated that this RDR sets no threshold. Section: Failure Modes.
- **fixed** — CV-011 (Round-Trip named an inverse never specified): restated the
  invariant over the **normalized value**, not the dump text — `parse ∘
  normalize`, not `parse ∘ normalize ∘ dump` — and said plainly that no dump
  grammar is defined, so no textual round-trip is claimed. Named the three
  lossy rendering sites (`<clear>` unreserved, unescaped separators, unreserved
  suffix separator) as work owed by any later RDR defining a re-readable dump.
  Section: Round-Trip / Inverse Invariants.
- **fixed** — CV-012 (the MVV's sharpest assertion had no possible witness):
  the RDR fixture must now carry **two sibling candidate rows binding the same
  outcome**, without which the gate ordering is untestable. Also replaced the
  "one rewind or cluster guard" item, dischargeable by a rule carrying no guard
  at all, with "one guard that actually carries a predicate." Section: MVV.
- **fixed** — CV-013 (`Resolve` handshake unwitnessed; constants asserted
  against absent symbols): scenario 6 now states the assertion references the
  kernel constants directly, compiles only after RDR 0007's reshape, and is
  owed at that point — and notes the spike mirrors them locally. Section:
  Testing Strategy 6.
- **fixed** — CV-014 (category oracle was a substring match): scenario 3 now
  requires one mutated fixture per category, asserted **by category, not message
  text**, tripping one category and no other; added the previously
  uncontrolled classes. Section: Testing Strategy 3.
- **fixed** — CV-015 (determinism oracles pass on a no-op; controls missing):
  added the "every oracle must have a failing control" rail to the MVV —
  category-level assertions, permutation on the RDR fixture rather than only
  kata, `exists = true` as well as `false`, a `Status`/`status` negative
  control, and an unknown-schema-field fixture with the decoder configured to
  reject unmapped keys. Sections: MVV; Testing Strategy 5.
- **fixed** — CV-016 ("comparator is total" satisfied by the defect it
  excludes): replaced with **rule-order permutation** as the assertion that
  actually excludes a positional tiebreak, and said why the adjacent-pair check
  does not. Swept the matching claim in Risks and Mitigations. Sections: MVV;
  Risks; Testing Strategy 5.
- **fixed** — CV-017 (totality proof holds per-model; nothing bounds a dump to
  one model): scoped the totality claim to one model's rows, required a
  model-unique `model id` for a multi-model dump, added a `duplicate model id`
  load category, and recorded that RDR 0006 lints one model per invocation.
  Sections: Normative Contracts / dump + validation categories.
- **fixed** — CV-018 (gate-then-count restated as this RDR's own normative
  block): kept the clause this RDR owns — row order MUST NOT decide selection,
  and normalization emits no positional field — and cited JDR 0001 §D2 / the
  kernel for the procedure with an explicit MUST NOT restate, mirroring the
  discipline already applied to the tag type model. Verified the
  non-escapability property survives as descriptive prose in Technical Design,
  Failure Modes, and Testing Strategy 4. Section: Normative Contracts.
- **fixed** — CV-019 (escape-row emptiness stated as a derived guarantee):
  restated as a **producer obligation, not a kernel guarantee**, citing
  `escapeOrRefuse`'s shared gate and `TestAdv2b_EscapeEdgeMustNotBypassThe`
  `OwnedStateRequirement`, which freezes the opposite as live kernel behavior.
  Section: Normative Contracts / `RequiresOwned`.

## Charted to successor

- **charted-to-successor** — CV-011's dump grammar. Defining a re-readable dump
  syntax (delimiters, escaping, a reserved `<clear>` sentinel, a reserved
  expansion-suffix separator) is a new contract, not a repair to this one. The
  draft now scopes its invariant to the normalized value and names the three
  lossy sites, so nothing is silently assumed. See `Charted:` below.

## Routed to Stage 6 (cross-RDR — not a 0002 edit)

- **needs-verification** — CV-020: RDR 0008's normative block asserts "nothing
  in this RDR or RDR 0002 forbids that alphabet entry" (the empty-string
  outcome). RDR 0002 now forbids it and A8 leans on that clause as the sole
  barrier. The *decisions* are compatible — 0008 scoped itself out ("not this
  RDR's to rule on") — but a Final peer carries a stale factual claim about this
  RDR's contents. We do not amend RDRs; this is a reconcile item.
- **needs-verification** — CV-021: RDR 0002 is `Draft` while 0003, 0006, 0007,
  0008, 0009 are `Final` and 0001 is `Implemented`; RDR 0009 twice cites "RDR
  0002's **Final** escape-rule prohibition" and rests its enforcement argument
  on it. The producer being the cluster's only unlocked document inverts the
  dependency order. Every contract changed in this pass should be checked
  against the Final peers' verified assumptions before lock —
  `/rdr-cluster-reconcile` is the mechanism.

## Dismissed with cite

- **dismissed-with-cite** — "`RequiresOwned` derived from writes+clears
  conflicts with the kernel's `missingOwned` running over guard-unevaluable
  survivors." Re-raise of an adjudicated decision: RDR 0007's Metadata says it
  "*fixes* the meaning of `Row.RequiresOwned` (post-guard write-dependency
  keys)" and names kata `xg7p` — "RequiresOwned conflates guard-input state with
  post-guard transition state" — as the defect it closes. Narrowing away from
  guard-input state is the deliberate call, and 0007 A7 covers the residue
  (every guard-referenced key is declared, so a typo fails at table load).
- **dismissed-with-cite** — "`next-state tags` and `writes` being value-identical
  undermines RDR 0009 A4." Inverts A4, which settled that the breach predicate
  stays `Writes`-only and explicitly does **not** widen to `NextTags`, because
  only `Writes` reaches the accessor layer. Two fields populated from one
  rendered set with different downstream consumers is coherent and matches
  `resolve.go::planOf`, which copies both.
- **dismissed-with-cite** — "The normative spike evidence is missing from the
  repo." False against disk: `evidence/spikes/` holds `output.txt` (1899 B),
  `negative-cases.txt` (1711 B), both fixtures, and `main.go`; all read directly
  this pass, `git status` clean.

## Needs (re)verification — carried to Stage 6

No previously `Verified` assumption was invalidated. A5 stays `Verified` (its
load-bearing half was confirmed; only a precedent clause was corrected). New
load-bearing claims introduced by this pass, none silently absorbed:

1. **The load/lint arity split** (CV-004) is a new normative boundary asserting
   that every cross-row check is RDR 0006's and every single-rule check is this
   RDR's. Needs checking against RDR 0006's actual finding set — it must have a
   category for ambiguous overlap, gap, dead row, and read-before-write, and
   must not expect this RDR to refuse any of them at load.
2. **`duplicate model id`** (CV-017) is a new load category with no witness and
   no peer binding. Needs an RDR 0005 code and a fixture.
3. **`cyclic context inheritance`** (CV-007) is a new load category. The spike
   implements the check but `negative-cases.txt` has no fixture for it.
4. **A suffix attaches only where a rule expands** (CV-008) is a new
   normalization rule with no witness — the spike's `liftOutcome` currently
   suffixes on the `in` operator, not on the expansion count, so a single-member
   `in` would take a suffix under the spike as written.
5. **Row kind as a derived property** (CV-005) needs confirming against RDR
   0006, which consumes the dump and may expect a kind field.
6. **The two-sibling fixture requirement** (CV-012) changes the canonical
   fixture, so `output.txt` and its recorded SHA-256 will change when the
   fixture is regenerated. The Performance Expectations digests are now stale
   with respect to the fixture the MVV requires.

## Charted

- **Charted:** A re-readable expanded-table dump grammar — delimiters, escaping,
  a reserved `<clear>` sentinel distinct from the tag-value space, and a
  reserved expansion-suffix separator — is out of scope for this RDR, which
  defines the dump's field list and ordering but treats rendering as a review
  surface. Suggested successor: a dump-format RDR, required before any consumer
  parses a dump rather than the normalized value. Raised by CV-011.

## Tiebreakers

None. Every fork collapsed on evidence: the load/lint split resolved on arity
(single-rule vs cross-row, which also matches where the spike already refuses);
row kind resolved on RDR 0009's normative prohibition plus the kernel's actual
`Row`; the round-trip resolved by scoping the invariant to what this RDR
actually specifies rather than inventing a grammar.
