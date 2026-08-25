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

---

## Phase 3a — CoVe verification

Independent Chain-of-Verification pass over the 147 REQs in
[`req-list.md`](req-list.md). Method: for each REQ, name an input that would
make a correct implementation visibly violate it, then run that input against
the shipped package through its public API (`table.Load`, `table.Dump`,
`table.DumpAll`, `table.CheckModelIDs`, `Row.KernelRow`, `Model.KernelTable`,
and `resolve.Resolve`). Probes were throwaway `go run` programs over
hand-authored TOML, deleted after the run.

**Independence.** No `internal/table/*_test.go` file and no
`internal/table/testdata/` fixture was read. Every fixture below was authored
fresh for this pass.

### Coverage

Meaningfully probed: **131 of 147** REQs. The 16 not probed are REQ-116
(source rewrite, explicitly out of scope), REQ-104 (a MUST-NOT-restate
clause with no runtime surface), REQ-111/REQ-112/REQ-117/REQ-144 (RDR 0006's
lint, unimplemented and deferred by REQ-117 itself), REQ-118/REQ-119/REQ-124/
REQ-131/REQ-132/REQ-133/REQ-134/REQ-139 (obligations on the *test suite*,
which this phase is barred from reading), REQ-145 (a derived, non-blocking
dump property), and REQ-7 (parser choice, confirmed statically in `go.mod`).

REQ-MVV was probed end-to-end: a hand-authored model normalizes, dumps, and
resolves through `internal/resolve::Resolve` to exactly one ordinary row, to
exactly one modeled `no_match` escape row, and refuses `guard_unevaluable` on
the unevaluable-sibling tag-set — matching REQ-137's expectation. The overlap
leg is deferred per REQ-117.

Conforming under probe, among others: the version gate's precedence over
strict decoding (REQ-14/15/135), the full load-category taxonomy (REQ-106/110,
25 identifiers), block-agnostic atom validation across all three blocks
(REQ-88), per-operator conformance including the unchecked comparison bound
(REQ-89), the merge over the full atom identity under all five forgery
delimiters (REQ-46/58/126/128), set-write member sequences under the same
(REQ-38/127), the dump's total ordering and 200-run byte-identity
(REQ-98/99/102/122/146), `[dump]` never reaching the normalized value
(REQ-97), the round-trip invariant over key order, rule order, and
`eq`-vs-single-member-`in` (REQ-113), the no-alias mutation control
(REQ-82/125), byte-exact tag-key identity at all five sites (REQ-54/121), and
the §D13 canonical JSON array at the kernel seam for declared `set` tags,
including the one-member and `<clear>` cases (REQ-93).

### FAIL-1 — A multi-member literal on a single-value operator is silently truncated at the kernel seam

**REQ violated.** REQ-93 (§D13 landing, `deviations.md` D5): *"`Tag.Value`
stays `string`; a set crosses as its canonical JSON array — members sorted,
duplicate-free, compact encoding — and 0002 declares it."* The req-list adds
that this covers *"any set-valued atom literal handed to the guard seam"* and
that RDR 0007 *"already wrote its `contains` leg against it"*, so this build
must **match** that form. Also engages REQ-84/REQ-85 (the handoff routes each
atom to exactly one destination) and REQ-56 (a set-valued literal MUST NOT be
rendered into a single string as its normalized value — here the seam value
is not merely joined, it is *dropped*).

**Failing input.** Any multi-member literal on a tag whose declared `kind` is
not `set`. Minimal case, on a `kind = "scalar"` tag `status`:

```toml
[tags.status]
provenance = "owned"
kind = "scalar"

[[rule]]
id = "advance"
[rule.match.recognized]
eq = "approve"
[rule.guard.all.status]
contains = ["done", "open"]
[rule.write]
status = "done"
```

**Observed.** Loads clean. The normalized atom is
`{status all contains [done open]}` — correct. `Row.KernelRow()` then yields:

```json
[{"Key":"status","Operator":"contains","Literal":"done","Block":"all"}]
```

`"open"` is gone, with no refusal and no diagnostic.

**Spec-required.** `internal/resolve/guardcontract.go` — the cross-RDR seam
contract RDR 0007 shipped — fixes the byte form for exactly these operators:
`{"contains/superset", "contains", ["alpha"], ["alpha","beta"], GuardTrue}`
and `{"in/member", "in", ["alpha","beta"], "alpha", GuardTrue}`. The literal
must cross as `["done","open"]`.

**Mechanism.** `model.go::seamValue` branches on `isSet`, which is
`Row.setKeys` membership — the *declared kind*. For any non-`set` tag it
returns `members[0]` unconditionally, discarding members 1..n. It is the
only writer of `resolve.Tag.Value` and `resolve.GuardAtom.Literal`.

**Reproduced on four independent paths**, all on a non-`set` tag:

| Site | Normalized value | Seam value |
| --- | --- | --- |
| `guard.all` `contains = ["done","open"]` | `[done open]` | `"done"` |
| `guard.all` `in = ["alpha","beta"]` | `[alpha beta]` | `"alpha"` |
| `guard.all` `eq = ["done","open"]` | `[done open]` | `"done"` |
| `[rule.write] status = ["done","open"]` | `[done open]` | `"done"` |

The `guard.all` `in` arm is the sharpest: the kernel hands that literal to
RDR 0003's evaluator, whose contract parses it as a JSON array. Given the
bare string `alpha` the seam must answer `GuardUnevaluable` — so the row
deadlocks the flow under RDR 0007's veto rather than refusing at load.

**Consequence.** Two rows whose normalized atoms differ only past member 0
are distinguishable in the dump and indistinguishable at the seam.
Demonstrated: rules `advance` (`status.eq=[AAA open]`) and `bravo`
(`status.eq=[ZZZ open]`) normalize to distinct atoms and cross as
`{"status":"AAA"}` and `{"status":"ZZZ"}` — both dropping `open`, the member
the author shared between them.

**Note on scope.** REQ-89 says `eq`, `in`, and `contains` *"take members"*, so
the RDR does not obviously forbid authoring a multi-member `eq`. Whichever way
that reads, the seam must not silently truncate: either the load refuses the
literal as a `malformed predicate atom`, or `seamValue` encodes the full
member sequence. It currently does neither.

### FAIL-2 — A nested one-element array literal is silently flattened

**REQ violated.** REQ-6 (`0002:C3`): *"**Strict decoding is an obligation on
this format, not a property of a library.** The decoder MUST reject unmapped
keys so an unknown schema field is a stable refusal rather than a silent
no-op."* The governing principle — a malformed authoring is a stable refusal,
never a silent reinterpretation — is defeated here at the value layer.
Engages REQ-64/REQ-106's *"literal ill-formed for its operator"* arm of
`malformed predicate atom`, and REQ-30's *"every value MUST be well-formed for
that tag's declared kind"* for the `[initial]` site.

**Failing input.** Any array literal containing a one-element nested array:

```toml
[rule.match.status]
in = [["a"], "b"]
```

**Observed.** Loads clean and normalizes as if the author had written
`in = ["a", "b"]` — expanding into rows `probe.advance#a` and
`probe.advance#b`. Arbitrary nesting depth flattens the same way:
`labels = [[[["a"]]], "b"]` normalizes to the write value `[a b]`.

**Spec-required.** `malformed_predicate_atom` (or
`malformed_initial_declaration` at the `[initial]` site). A nested array is
not a member sequence and the format admits no such spelling.

**Mechanism.** `load.go::valueMembers` recurses on `[]any` and rejects an
element only when `len(members) != 1`. A one-element nested array satisfies
that test, so it is unwrapped rather than refused. A *multi*-element nested
array *is* caught — `labels = [["a","b"], "c"]` refuses — which makes the
refusal depend on the inner array's arity rather than on its shape.

**Reproduced at four sites**, all accepted where a refusal is required:
a match-block `in`, a `guard.all` `contains`, a `[rule.write]` value, and an
`[initial]` assignment. Two controls confirm the gap is specific to
`valueMembers`' arity test rather than general: an inline table
(`eq = {a = 1}`) and a TOML datetime both refuse `malformed_predicate_atom`
correctly, and the `<clear>` ban still fires post-flattening.

**Severity.** Lower than FAIL-1 — no value is lost, and the flattened reading
is the one the author most likely meant. It is filed because REQ-6 makes
"stable refusal, never a silent no-op" the format's obligation, and because
the arity-dependent behaviour means the same authoring error refuses or
succeeds depending on how many elements the nested array happens to hold.

### Cross-checks against Phase 3b

Phase 3b's three findings were re-derived independently here before its
section was read, and all three reproduce:

- **ADV-1** — `[rule.guard.unless.status] in = ["small", "mid"]` expands to
  **2 rows** carrying suffixes `#small` / `#mid`, each with the atom rewritten
  to `eq` in the `unless` block. Confirms the REQ-72 (match-blocks-only) and
  REQ-51 (no rewriting across blocks) violation Phase 3b filed. Not double-
  counted here.
- **ADV-2** — a rule write naming an observed tag refuses
  `malformed_accessor_binding`, not `write_to_non_owned_tag`. Confirmed.
- **ADV-3** — `in = ["approve", "approve"]` loads clean and mints **2 rows
  sharing the identity** `probe.advance#approve`. Confirmed.

### Unconfirmed — suspected but not demonstrated

- **REQ-70 / REQ-131, outcome-literal category.** A `recognized` literal
  outside the alphabet refuses `malformed_outcome_binding` only when the
  `recognized` declaration's `domain` does not already exclude it. When
  `[tags.recognized]` is `kind = "enum"` with `domain` equal to the alphabet
  — the natural authoring, and the shape a reviewer would expect — the same
  mutation refuses `malformed_predicate_atom` instead. REQ-131 requires each
  variant to refuse *"the **one** category its mutation targets"*, so a
  negative control written against the natural declaration would observe a
  different category than one written against `kind = "scalar"`. Marked
  unconfirmed rather than a FAIL because the RDR fixes the *failure* but not
  the *category* for this arm (`0002:1094` says only "is a load failure"),
  and REQ-16 leaves the order of independent checks deliberately unspecified.

- **REQ-30, `[initial]` provenance arm.** An observed key in `[initial]`
  refuses `malformed_accessor_binding`, not `malformed_initial_declaration`.
  **Not a defect**: `0002:725-734` states the provenance arm is *"unreachable
  by construction"* and that a non-owned key *"refuses as `malformed accessor
  binding` … before the owned-tag predicate here is consulted"*. The
  implementation matches the record exactly. Recorded so a later reader does
  not re-open it.

### Phase 3a summary

| Entry | REQ | Status |
| --- | --- | --- |
| FAIL-1 | REQ-93 (also REQ-84/85, REQ-56) | Confirmed, 4 paths |
| FAIL-2 | REQ-6 (also REQ-64/106, REQ-30) | Confirmed, 4 sites |

2 confirmed violations, 131/147 REQs probed, 3 Phase 3b findings independently
reproduced.
