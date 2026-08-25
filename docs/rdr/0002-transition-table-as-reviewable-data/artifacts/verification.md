# Verification — RDR 0002

## Phase 3b — Adversarial review

Independent adversarial pass. Anchored in the record's
`## Trade-offs > ### Failure Modes` section. Tests were ADDED only; no
existing test was weakened or modified. Tests live in
`internal/table/adversarial_0002_test.go`.

Baseline: before these tests, `go test ./...` passed clean across
`internal/cli`, `internal/resolve`, and `internal/table`. Every failure now
in the tree is one of the tests below.

---

### ADV-1 — Expansion escapes the match blocks

**Failure mode.** The Failure Modes section splits load-time refusals
("malformed predicate atoms … each is decidable from one rule plus the
declarations") from cross-row lint ("ambiguous expansion … is therefore an
RDR 0006 lint failure"). Expansion is the record's match-block-only
mechanism: REQ-72 (`0002:C13`) scopes it to "a rule's MATCH BLOCKS — local
`match` and inherited context `match`"; REQ-73 has the expanded atom "still
carrying block `match`"; REQ-51 (`0002:C7`) forbids folding `unless` atoms
into `all` or match atoms into either guard block.

`internal/table/normalize.go::expand` decides expansion on
`c.atom.Operator == "in"` with **no block test**, so a `guard.all` or
`guard.unless` `in` atom expands too.

**Why it is quietly wrong.** A guard `in` is one predicate over a member
set: `unless profile in {small, mid}` means "not small AND not mid" — one
row. Expanding it mints one row per member, each excluding only its own
member, converting the author's conjunction into a disjunction. Every minted
row now MATCHES a view the author excluded. Nothing downstream can recover
the intent: the extra rows are well-formed, carry a plausible expansion
suffix, and the dump certifies them as canonical; RDR 0006 sees rows, not
the authored guard. The over-reach also bleeds `#` into the rendered row
identity, because REQ-60/REQ-90 deliberately scope the `#` ban to
match-block `in` members ("since only match blocks expand"), so a `#`-bearing
guard member is legally authorable — and, once wrongly expanded, forges the
ambiguous identity `0002:C11` exists to close.

**Tests.**
- `TestAdv1_GuardBlockInMustNotExpand` — **FAILS.** Rule expands to 4 rows,
  want 2; each row carries a 2-element suffix, want 1; the `unless` atom is
  rendered `eq` on one member, want `in` over `[mid small]`; and each row
  excludes only one of the two members the author excluded.
- `TestAdv1b_GuardInMustNotBleedHashIntoRowIdentity` — **FAILS.** Identity
  renders as `rdr.continue-prelock#a#b#large` — 3 `#` separators where 1 is
  legal, so splitting the rendered identity on `#` is no longer unambiguous.

---

### ADV-2 — A rule writing a non-owned tag is refused by the wrong authority

**Failure mode.** The Failure Modes section names "**a write to an
undeclared or non-owned tag**" as its own stable load refusal, and closes
the list with "each is decidable from one rule plus the declarations".
`0002:C14` (REQ-77) then relies on it: `RequiresOwned` is derived from the
write block and clear list, and "**by the write-to-non-owned-tag rule** every
such key is an owned tag".

`internal/table/load.go` checks provenance only on `[write.<id>].keys`
(`loadAccessors`). A rule whose write block or clear list names an observed
tag is caught downstream by `checkAccessorBindings`' writer-arity count and
reported as `malformed_accessor_binding`.

**Why it is quietly wrong.** Two things. (1) REQ-110 makes the category
identifiers an API surface — the caller is told the model's accessor wiring
is wrong when the defect is the rule, and the existing
`TestReq27_ReaderArityAndWriterProvenanceCarryDistinctCategories` freezes
that these two obligations must not collapse into one category; they collapse
here. (2) The refusal is not decidable from one rule plus the declarations —
it needs the whole writer table and every rule's write set — so the RDR's
own decidability claim for this category is false, and `RequiresOwned`'s
all-owned guarantee rests on a check that never inspects the rule.
`RequiresOwned` is the kernel's owned-state gate
(`resolve.go::missingOwned`), so a non-owned key reaching it is a runtime
refusal the load was supposed to pre-empt.

**Test.** `TestAdv2_RuleWriteToNonOwnedTagIsRefusedAsSuch` — **FAILS**, all
three arms:
- `write block names an observed tag` → got `malformed_accessor_binding`,
  want `write_to_non_owned_tag`.
- `clear list names an observed tag` → same.
- `rule provenance does not collapse into accessor arity` → both the arity
  defect (`neg/neg-owned-no-reader.toml`) and the rule-provenance defect
  refuse with `malformed_accessor_binding`.

*Honest note:* an earlier third arm — an observed write with a writer
declared to serve it — PASSED, but for a coincidental reason (the refusal
fires on the writer entry, not the rule), so it was not a discriminating
case and was replaced by the category-non-collapse oracle above.

---

### ADV-3 — A repeated `in` member breaks the identity tuple's totality

**Failure mode.** The Failure Modes section requires malformed predicate
atoms to be refused at load. `0002:C13` (REQ-72/REQ-76) states the
dependency outright: "**Because RDR 0003 rejects a repeated `in` element at
parse, no two rows of one rule share a suffix, which is what keeps the
identity tuple total.**" REQ-64 (`0002:C22`) lands the enforcement in this
package, which owns "the two load categories that carry RDR 0003's rejection
rules" — `malformed predicate atom` being one. RDR 0003's parser is not what
loads this document; this package is. The premise is unenforced here.

**Why it is quietly wrong.** `in = ["round-clean", "round-clean"]` loads
clean and mints two byte-identical rows under one identity. REQ-98's total
order over `(model id, rule id, expansion suffix)` is void for that model,
so the dump's determinism guarantee lapses silently; the per-rule expansion
count the Failure Modes section says is "derived from the dump" is inflated
against a rule the author wrote once. On an **escape** rule it is worst:
`0002:C5` (REQ-41) has each escape row rescue the outcome it binds and the
kernel counts survivors, so two identical escape rows binding one outcome are
two survivors — a self-inflicted `ambiguous_match` on the very row that
exists to rescue one.

**Tests.**
- `TestAdv3_RepeatedInMemberIsRefused` — **FAILS**, both arms (outcome
  binding and inherited context match atom): the document loads clean where
  `malformed_predicate_atom` is required.
- `TestAdv3b_RowIdentityIsTotalAcrossTheModel` — **FAILS.** Consequence
  oracle, satisfied by either refusal or unique identities. Observed:
  `rdr.draft-no-match-escape#round-clean` carried by 2 rows;
  `rdr.continue-prelock#large` and `rdr.continue-prelock-cluster#large`
  each carried by 2 rows.

---

### Phase 3b summary

| Entry | Test | Currently |
| --- | --- | --- |
| ADV-1 | `TestAdv1_GuardBlockInMustNotExpand` | FAILS |
| ADV-1 | `TestAdv1b_GuardInMustNotBleedHashIntoRowIdentity` | FAILS |
| ADV-2 | `TestAdv2_RuleWriteToNonOwnedTagIsRefusedAsSuch` (3/3 arms) | FAILS |
| ADV-3 | `TestAdv3_RepeatedInMemberIsRefused` (2/2 arms) | FAILS |
| ADV-3 | `TestAdv3b_RowIdentityIsTotalAcrossTheModel` (2/2 arms) | FAILS |

5 added test functions, all failing against the current implementation. No
pre-existing test was weakened, and no pre-existing test fails.
