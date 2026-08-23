Model: claude-opus-5

# Reconciliation — RDR 0003 iteration-3 hostile critiques (A: opus-5, 17 rows; B: sonnet-5, 10 rows)

Reconciled by passage anchor, then by prose. `C-N` ids do not correspond across
the two source files; the `A row` / `B row` columns carry each pass's own id.

## 1. Reconciled ledger

| R-NN | Passage anchor | Defect | A row | B row | Agreement | Severity |
| --- | --- | --- | --- | --- | --- | --- |
| R-01 | `trace` step 6 `:1141` + step 3 `:1138` against `evidence/spikes/guard-fixture.toml` | The document's only end-to-end worked example stamps coverage OK for the `status eq "Draft"` group, whose two rows accept `{mid,large}` and `{foundational}` over a 4-value `profile` domain — `small` is uncovered. Step 3 also omits the `{present,absent}` presence dimension that `cluster_eligible` (no `optional` marker → optional by `:775-777`) requires by `:922-934`, while step 5 discusses the `exists` atom without adding the dimension to the product. The one desk-check that would catch a coverage-arithmetic bug was run without doing the arithmetic. | C-1 | — | A-ONLY | **BLOCKING** — false-green in the RDR's own witness |
| R-02 | Cardinality clause `:1053-1055` ("the product of every participating dimension's **declared domain size**") vs. kind table `:841` (`set` → "its element universe") | For a `set` dimension the assignment space is `2^N` over an N-element universe, not N. The bound clause specifies the wrong quantity for exactly the kind whose blowup the too-large bound exists to catch, and the clause's own "a bound discovered by exhausting memory is not a conforming bound" is violated by its own arithmetic. | C-2 | — | A-ONLY | **BLOCKING** — wrong computation in a normative clause |
| R-03 | `resolve.go::gate` vs. escape-coverage clause `:936-941` + A17 `:522-546` + two-population clause `:943-958` | `gate()` returns `blocked` for `KindGuardUnevaluable` and `KindOwnedStateUnavailable` **before** the exact-one count, so `Resolve` returns `refuse(in, *blocked)` and never reaches `escapeOrRefuse`. An escape row cannot rescue those two refusal classes at all, yet the RDR treats escape rows as a coverage-closing population and cites this same kernel as its authority. Bare escape row closes lint coverage; runtime still refuses. | C-6 | — | A-ONLY | **BLOCKING** — false-green; P5 breached by the clause citing the kernel |
| R-04 | Disposition table `:1114` ("Full `unless` block decides true \| Row disabled \| Excluded intersection subtracted") + coverage derivation `:698-701` vs. RDR 0007 `:1435-1440` / `:1412-1421` | RDR 0003 models `unless` as two-valued set subtraction; RDR 0007 (the kernel this RDR delegates verdict combination to) computes `all_result ∧ ¬(unless_conj)` in strong Kleene with `¬U = U`, `T ∧ U = U`. An unevaluable atom inside `unless` makes the row verdict `U` at runtime while lint subtracts a decided set. The disposition table has **no row** for "`unless` atom unevaluable". | C-4 | — | A-ONLY | **BLOCKING** — false-green through the one block the narrowing never examines |
| R-05 | Optionality default `:775-777` × single-valued default `:810-818` | Two independently-defaulted-to-worst-case rules compose and no passage multiplies them. An unmarked 5-value enum is `2^5 × 2 = 64` assignments where the author meant 5; four such tags ≈ `1.7 × 10^7`. B raises the same optionality default from the runtime/authoring side: it silently licenses `guard_unevaluable` and doubles the product, invisible unless the author reads this RDR's prose. | C-3 | C-8 | **BOTH** | **BLOCKING** — the composed default makes the central guarantee unreachable in practice (A11's own stated "If wrong") |
| R-06 | Withheld-claim clause `:968-975` ("MUST name the participating row and the atom that can refuse") + `Consequent duty` note `:977-984`; `resolve.go:270-273` `Refusal.Guard string`; `gate()` `slices.MinFunc` | A normative MUST whose producer field does not exist. A frames it against RDR 0006's finding contract (`0006:363-367` has no atom-level field, conceded in the RDR's own block quote, then declared non-gating). B frames it against the shipped kernel: `Refusal.Guard` is a single opaque string and `gate()` reports only the lexicographically-lowest undecidable row's guard, contradicting "one finding per refusing row, never just the first". Same defect — the MUST has no satisfiable producer on either surface. | C-11 | C-4 | **BOTH** (different surfaces) | **BLOCKING** — a normative MUST no consumer can satisfy |
| R-07 | Coverage clause `:936-941` + overlap-exclusion `:952-954` vs. anti-opt-out clause `:685-698` | A bare escape row closes coverage by itself **and** is excluded from the ordinary overlap check, so it generates no finding. One line, always available, blocking-free — precisely the opt-out the RDR spends two paragraphs forbidding for declarations. B reaches the same construct from the other end (premortem incident 4: authors route around the bound with "an escape row that swallows everything", named as A7's own anti-pattern). | C-7 | C-7 (premortem) | **BOTH** (A normative, B via consequence) | SUBSTANTIVE — near-blocking; the guarantee is voidable by one authored line |
| R-08 | Participation clause `:910-919` ("Match keys are not product dimensions") vs. fixture rules `profile-to-grounding` (`[rule.match.status]`) and `reconcile-rewind-legality` (`[rule.guard.unless.status]`) | The clause partitions keys by *authoring block*, not by key, and the RDR's own fixture uses `status` both ways. No rule is given for a key that is both; `resolve.go::Row` carries `Match` and `Guard` as separate fields so the ambiguity survives into the shipped shape. Divergent dimension counts across implementations. | C-5 | — | A-ONLY | SUBSTANTIVE |
| R-09 | Participation clause `:910-919` ("A key every row constrains identically still bounds the product … MUST NOT be dropped") | A harmless, universally-shared, zero-discriminating-power guard atom newly requires its key's domain to be finite-declared, converting a previously-clean group into a blocking lint failure. Distinct from R-08: same clause, different defect (R-08 is *which* keys are dimensions; R-09 is the cost of a non-discriminating one). | — | C-6 | B-ONLY | SUBSTANTIVE |
| R-10 | RDR 0006 `disposition` rider — "default-on for every scoped row group whose participating dimensions are all finitely declared" (`0006:913-915`) | Default-on means declaring one more domain for an unrelated reason (documentation) silently opts an existing working group into a blocking gate. A docs-only change breaks CI; the author's cheapest fix is to revert the declaration, making the model less precise — the outcome A11's rehoming was meant to prevent. No migration path or grace tier is stated. | — | C-5 | B-ONLY | SUBSTANTIVE |
| R-11 | A17 `:522-546` + `authority` table `:1100` + Contradiction Check `:1897-1907` vs. JDR 0001 §JD-14 (`jdr/0001:321-337`) and RDR 0006 (`0006:906-914`) | This RDR asserts a two-population overlap reading while JDR 0001 §JD-14 still records "**Decided 2026-08-22: 0003's reading governs** — escape rows are ordinary participants in **both** the union and the overlap check", and RDR 0006 independently records the same. Neither peer has been amended. Two implementers reading different cluster documents ship contradictory blocking lint. | — | C-1 | B-ONLY | SUBSTANTIVE (partially adjudicated — see §4, R-11) |
| R-12 | A14 `:434-464` ("a gap to close, not a contradiction to resolve") vs. RDR 0002 `:309-311` ("Normalization MUST combine both into one candidate-row predicate set") + identity tuple `:1162-1170` | The identity tuple requires `block` to survive as a per-atom field; RDR 0002's normative clause says *combine*. A calls this a contradiction between two normative clauses downgraded to a request by asserting compatibility without citing a reading under which both hold. B grounds the same gap in RDR 0002's own refine backlog (`0002:897-899`, listed as an open item, not a contract) and adds the semantic consequence: if `Block` does not survive, guard exclusion silently becomes inclusion — a semantic inversion with no compile- or lint-time error. | C-8 | C-3 | **BOTH** | SUBSTANTIVE |
| R-13 | A15 `:465-499` + MVV Scenario 3 `:1784-1800` | The bound clause is provisionally accepted on a test that has not run, and A15's own evidence concedes "proof cost … may depend on dimension count or per-dimension width, not on cardinality alone". A adds the sharper defect: the test reads the implementation's published bound `B` at fixture setup and sizes its products relative to it, so it derives its oracle from the system under test and cannot falsify the claim it is named to discharge — it never measures provability at all. B adds the user-facing consequence: the diagnostic reports cardinality, so an author shrinks their domain when the actual cost driver is dimension count. | C-10 | C-7 | **BOTH** (A: unfalsifiable test; B: clause may not survive validation) | SUBSTANTIVE |
| R-14 | A10 `:...` + A12 `:...` (both `Pending`, homed at "RDR 0006's refine") + owned-tag clause `:1075-1078` + phase table `:1699-1704` | Two load-bearing derivations — the row-group predicate every product is scoped by, and the reachability relation making the owned-tag clause decidable — are deferred to a `Draft` peer's unstarted refine. A12's own text concedes RDR 0006 frames the check as *graph reachability* while this RDR requires a *syntactic decision over declarations*: two decision procedures over one normative predicate, neither chosen. A10's evidence concedes "a cycle with no floor", broken by fiat on one side only. | C-12 | C-2 | **BOTH** | SUBSTANTIVE |
| R-15 | Metadata `:22-24` + Prerequisites `:1567-1659` + phase table `:1699-1704` | Six of seventeen assumptions are Pending, four homed at other documents' refine passes, each individually argued "not lock-blocking" with no passage assessing the aggregate. Phases 1–3 are "Partially" startable and Phase 4 is "No" — the RDR locks a contract whose integration phase is admittedly not startable. B frames the same fact as churn rate (five re-verification passes, a rehomed model, A16/A17 opened the same day) predicting a rewrite. | C-9 | C-9 | **BOTH** | SUBSTANTIVE |
| R-16 | `Illustrative Code` `:1219-1245` — sole declaration example, `single_valued` spelling disclaimed at `:1243-1245` | The RDR owns a five-field type system and gives exactly one example whose spelling it disowns. A adds: no example anywhere shows a `set` declaration with `contains`, a computed cardinality, or a withheld claim, and `guard-fixture.toml` declares none of `optional`, `single_valued`, `elements` — so three spellings coexist. B adds the first-contact symptom: a developer copying the block gets an unknown-field error from RDR 0002's loader with no way to know from RDR 0003 alone that the fix is pending elsewhere. | C-13 | C-10 | **BOTH** | SUBSTANTIVE |
| R-17 | Conformance clause `:787-798` + `disposition` table `:1110-1126` | Conformance is the premise every lint claim is conditional on, and the clause routes the check to "the kernel's view assembly" without citing a clause in RDR 0007 or 0006 that requires it. `resolve.go::assemble` performs no conformance check. There is no disposition row for a value atom whose key is declared always-present but absent in a non-conforming view. A premise with no enforcer. | C-16 | — | A-ONLY | SUBSTANTIVE |
| R-18 | Scope Verification `:2010-2025` + Proportionality `:2051-2064` | The tag declaration model was absorbed at Stage 6, after grounding, 3amigo, and critique had run against a document that did not contain it, and only "that row" was re-run. Six alternatives are analyzed for the operator grammar; zero for the type system. A11's only rejected alternative is a *placement* alternative — a Design Decision whose sole rejected alternative is where to put the decision has not evaluated the decision. | C-14 | — | A-ONLY | SUBSTANTIVE |
| R-19 | A1 `:139-173` + Testing Strategy `:1758-1767` + Phase 3 `:1726-1728` | A1 is `Verified` on a harness that "carries its own TOML-to-atom encoding … and consumes no declared domain", against a fixture declaring none of the five fields this RDR now owns. `contains` is the one operator with no evaluated-plus-declared evidence, and Phase 3 is required to add a `contains` predicate over a declared set-valued tag "before the full operator vocabulary is accepted" — i.e. the vocabulary is accepted at lock on a condition scheduled for after lock. | C-17 | — | A-ONLY | SUBSTANTIVE |
| R-20 | Metadata Status block `:9-25` + assumption records generally (A8 `:276-306`, A17 `:522-546`) | The document has become a ledger of its own governance: a reader must parse a four-document dependency graph, two JDR decisions, one re-opened JDR decision, and a cluster-gate report before reaching guard evaluation. The implementer reads Normative Contracts and skips the assumptions, where every unresolved arithmetic question lives. | C-15 | (C-9 overlap) | A-ONLY | NOISE — real but a readability complaint, not a defect in the contract |

## 2. Independent verification against source

Verified directly, not taken from either pass. Note one correction to both
passes' citation: the fixture lives at
`docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/guard-fixture.toml`,
not at a repo-root `evidence/spikes/`.

### V-1 — A's C-1 (R-01), trace step 6 coverage claim — **VERIFIED**

Fixture, `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/guard-fixture.toml`:

- `[tags.profile] domain = ["small", "mid", "large", "foundational"]` — four values.
- `profile-to-grounding`: `[rule.match.status] eq = "Draft"`, `[rule.guard.all.profile] in = ["mid", "large"]`.
- `foundational-to-cove`: `[rule.match.status] eq = "Draft"`, `[rule.guard.all.profile] eq = "foundational"`.

Both rows share `match status eq "Draft"`, so by trace step 2 they form one
group. Union of accepted `profile` assignments = `{mid, large} ∪ {foundational}`
= three of four. **`profile = small` is accepted by no row in the group.** The
group has a coverage gap; trace step 6 (`:1141`) stamps `OK` and the trace's
closing line states "No CONTRADICTION row and no GAP row."

Optionality marker: **VERIFIED absent.** `[tags.cluster_eligible]` declares
`provenance = "observed"`, `kind = "bool"`, `domain = [true, false]` — no
`optional` field. Under `:775-777` ("A declaration carrying **no optionality
marker declares the key optional**") it is optional, and under `:922-934` an
optional key contributes a `{present, absent}` presence dimension. Trace step 3
builds the product as `profile` (4) × `prelock_iterations` (`{0..3}`) = 16 and
does not include the presence dimension; step 5 discusses the `exists` atom in
detail and never adds the dimension to step 3's product.

One nuance in the RDR's favour, which does not rescue the row: step 6's cell is
worded as which product each *kind* of group computes over, rather than as an
arithmetic result. That makes the omission an unperformed check rather than a
stated-wrong number — but the step is still stamped `OK`, and the group it
governs demonstrably has a hole.

**Verdict: VERIFIED.** Anchor: `guard-fixture.toml` rules `profile-to-grounding`
/ `foundational-to-cove`; RDR `:1138`, `:1140`, `:1141`, `:1147`.

### V-2 — A's C-2 (R-02), `set` cardinality — **VERIFIED**

The cardinality clause reads, at `:1053-1055`:

> The quantity both figures report is the **cardinality of the scoped product** —
> the number of assignments in it, the product of every participating dimension's
> declared domain size — not a bitset width, byte size, or row count.

"declared domain size" confirmed verbatim. The kind table (`:841`) row for `set`
reads: `| set | its element universe | yes (required to claim exhaustiveness) |
**no** — values co-occur by construction |`. So a `set`'s declared finite domain
*is* its element universe, of size N, and the cardinality clause multiplies N.

A `set`-valued tag's assignment space over an N-element universe is the power
set: `2^N`. The RDR itself supplies the reasoning that makes this unavoidable —
the single-valued clause (`:810-818`) says a tag whose values can co-occur
contributes "one **independent boolean dimension per value** (`2^|domain|`)",
and the kind table's own `set` row says values "co-occur by construction" and
forbids the single-valued marker. The two clauses therefore *require* `2^N` for
`set` while the cardinality clause counts N.

**Verdict: VERIFIED.** The bound clause specifies the wrong quantity, by an
exponential factor, for the one kind whose blowup the bound exists to catch.
Anchors: RDR `:1053-1055`, `:841`, `:810-818`.

### V-3 — A's C-4 (R-04), `unless` two-valued vs. three-valued — **VERIFIED**

RDR 0007 `docs/rdr/0007-guard-predicate-totality.md:1412-1414`:

> Atom verdicts combine in the KERNEL under strong-Kleene … negation is
> `¬T = F`, `¬F = T`, `¬U = U`

and `:1435-1436`:

> The row verdict is `all_result ∧ ¬(unless_conj)`, where `unless_conj` is the
> conjunction of the `unless` block's atoms

RDR 0007 additionally states the exact case at `:1999`: "absent key is `¬U = U`
and vetoes — correct by the rule". So with `all_result = T` and one `unless`
atom `U` (none `F`), `unless_conj = U`, `¬(unless_conj) = U`, row verdict `U` →
`guard_unevaluable`, a blocking refusal.

RDR 0003's model of the same object is two-valued. Coverage derivation
(`:698-701`): a row's accepted assignments are "the intersection of all positive
`all` atom domains **minus** the single conjunctive assignment set matched by
the row's full `unless` block". Disposition table (`:1114`): `| Full unless block
decides true | Row disabled | Excluded intersection subtracted | none | Silent by
design |`.

Disposition row for "`unless` atom unevaluable": **confirmed absent.** I read the
full table at `:1110-1126`. It carries "Value atom over an absent key" and
"Existence atom over an absent key" rows, and the `unless`-decides-true row
above, and no row whose input class is an undecidable `unless` block. The
narrowing clause (`:961-965`) and the withheld-claim clause do formally cover a
value atom "over a key declared optional" without naming a block, so the case is
not formally unreachable — but the subtraction language actively contradicts it
(you cannot both subtract a decided set and withhold the claim), and the
disposition table is where an implementer goes for input-class-to-outcome
mapping.

**Verdict: VERIFIED** on all three sub-claims (0007 computes three-valued; 0003
models two-valued subtraction; no disposition row exists). Anchors:
`0007:1412-1414`, `0007:1435-1436`, `0007:1999`; RDR 0003 `:698-701`, `:1114`,
table `:1110-1126`.

### V-4 — A's C-6 (R-03), `gate` blocks before the exact-one count — **VERIFIED**

`internal/resolve/resolve.go::Resolve` (line 318) and `::gate` (line 376).

`Resolve`, lines 341-354:

```go
	selected, blocked := gate(candidates, in.Guards, view)
	if blocked != nil {
		return refuse(in, *blocked), nil
	}

	switch len(selected) {
	case 1:
		return Result{Plan: planOf(in, selected[0], false)}, nil
	case 0:
		return escapeOrRefuse(in, view, Refusal{Kind: KindNoMatch}), nil
	default:
		return escapeOrRefuse(in, view, Refusal{
			Kind: KindAmbiguousMatch,
			Rows: rowRefs(selected),
		}), nil
	}
```

`escapeOrRefuse` is reachable **only** from the `case 0` and `default` arms of
the switch, and the switch is reached only when `blocked == nil`.

Inside `gate`, both blocking returns precede the count:

```go
	if missing := missingOwned(survivors, view); len(missing) > 0 {
		return nil, &Refusal{Kind: KindOwnedStateUnavailable, ...}
	}
	...
	if len(undecidable) > 0 {
		...
		return nil, &Refusal{Kind: KindGuardUnevaluable, ...}
	}

	return selected, nil
```

So for `KindOwnedStateUnavailable` and `KindGuardUnevaluable`, `gate` returns
`(nil, blocked)`, `Resolve` takes `return refuse(in, *blocked)`, and
`escapeOrRefuse` is **never called**. An escape row cannot rescue either class.

**Verdict: VERIFIED.** A17 (`:522-546`) reads `Resolve` step 2 correctly — escape
rows are excluded from ordinary candidates — and cites that reading as its whole
evidence base, but neither A17 nor the two-population clause (`:943-958`) nor the
escape-coverage clause (`:936-941`) states that two of the four refusal classes
are decided before escape rows are consulted at all. Anchors:
`internal/resolve/resolve.go::Resolve`, `internal/resolve/resolve.go::gate`,
`internal/resolve/resolve.go::escapeOrRefuse`.

### V-5 — A's C-5 (R-08), `status` as both match key and guard key — **VERIFIED**

In `guard-fixture.toml`:

- `[[rule]] id = "profile-to-grounding"` → `[rule.match.status] eq = "Draft"` — `status` as a **match** key.
- `[[rule]] id = "reconcile-rewind-legality"` → `[rule.guard.unless.status] eq = "Implemented"` — `status` as a **guard** key.

Also `foundational-to-cove` uses `[rule.match.status]`. The participation clause
(`:910-919`) says "Match keys are not product dimensions" and justifies the
split by authoring block and by `resolve.go::Row` carrying `Match` and `Guard`
distinctly — which I confirmed (`Row.Match []Tag`, `Row.Guard string`). No clause
states the rule for a key appearing in both roles.

Mitigating: the two rules are in *different* groups (`match status eq "Draft"`
vs. `match stage eq "reconcile"`), so the fixture does not exhibit the
same-group collision, and trace step 3 does explicitly reason about `status` for
the Draft group ("`status` does **not** enter the product: its `status eq
"Draft"` atom is authored under `[rule.match.status]`"). The clause is
well-defined per-group as written.

**Verdict: PARTIAL.** The factual claim — the RDR's own fixture uses `status` as
a match key in one rule and a guard key in another — is **VERIFIED**. The
severity claim (that this produces divergent dimension counts) is weaker than A
states: the clause partitions per-atom-per-group, and no group in the fixture
contains `status` in both roles. The real gap is that no clause *says* the rule
is per-atom-per-group rather than per-key-per-model. Anchors:
`guard-fixture.toml` rules `profile-to-grounding` / `reconcile-rewind-legality`;
RDR `:910-919`; `internal/resolve/resolve.go::Row`.

### V-6 — A's C-16 (R-17), conformance enforcement — **VERIFIED**

`internal/resolve/resolve.go::assemble` (line 148) in full is a three-way merge
of `Recognized`, `Observed`, and `Owned` into a `TagSet` map, with owned values
overwriting observed. It reads no declaration model, checks no always-present
key, and validates no single-valuedness. There is no other conformance check in
`Resolve`'s path — `matches`, `missingOwned`, and `evaluateGuard` all test
individual row requirements, never the view against a declared model.

The conformance clause (`:787-798`) states conformance is "the premise every
lint claim in this RDR is conditional on" and routes the check to "the kernel's
view assembly", adding that whether it is refused at assembly or reported as a
lint finding "is RDR 0007's and RDR 0006's respectively" — without citing a
clause in either document that requires the check.

**Verdict: VERIFIED.** Anchor: `internal/resolve/resolve.go::assemble`; RDR
`:787-798`.

### V-7 — B's C-4 (R-06), does `gate()` report one row or all? — **REFUTED as stated; the underlying defect PARTIAL**

B claims `gate()` "picks the row with the lexicographically lowest `RowRef` among
all undecidable rows … as 'the' refusing guard, discarding the others' guard
strings entirely", and that the CLI "names only one row".

The code (`internal/resolve/resolve.go::gate`):

```go
	if len(undecidable) > 0 {
		lowest := slices.MinFunc(undecidable, func(a, b Row) int {
			return compareRefs(refOf(a), refOf(b))
		})
		return nil, &Refusal{
			Kind:  KindGuardUnevaluable,
			Guard: lowest.Guard,
			Rows:  rowRefs(undecidable),
		}
	}
```

`Rows: rowRefs(undecidable)` carries **every** undecidable row, not one. The
`lowest` variable feeds only the single `Guard` string field. So the refusal
payload does name every refusing row; what it cannot carry is more than one
row's *guard text*, because `Refusal.Guard` is a single `string`.

B's stated symptom ("the CLI error names only one row … not 'every refusing
row'") is therefore **REFUTED** — `Rows` is the full set. The residual defect is
real but narrower: `Refusal.Guard` is one opaque string, so the *atom-level* half
of the RDR's withheld-claim MUST ("MUST name … the atom that can refuse") has no
carrier in the shipped shape, and among N undecidable rows only the lowest row's
guard text survives.

Runtime vs. lint: this is the **runtime** path. `gate` is the kernel's resolution
gate; the RDR clause B cites ("one finding per refusing row, never just the
first", disposition `:1125`) is in the **Lint outcome** column and governs RDR
0006's lint findings, not `Refusal`. B conflates the two surfaces. The RDR's own
disposition row splits them explicitly: "RDR 0007 payload at runtime **and** a
blocking lint finding naming the refusing atom".

**Verdict: REFUTED as framed (one-row reporting); PARTIAL on the merged R-06
defect** — the atom-naming MUST genuinely has no producer, which A's C-11
establishes more accurately against RDR 0006's finding contract
(`0006:363-367`, no atom-level field), a gap the RDR concedes in its own block
quote at `:977-984`. Anchors: `internal/resolve/resolve.go::gate`;
`internal/resolve/resolve.go` `Refusal.Guard`; RDR `:968-975`, `:977-984`,
`:1125`.

## 3. Where the two passes contradict

### CT-1 — Whether the kernel reports one refusing row or all (R-06)

- **B (C-4)** asserts `gate()` reports only one row, discarding the others.
- **A (C-11)** never makes this claim; it locates the same MUST's failure in RDR
  0006's finding contract having no atom-level field.

**Evidence favours A.** `Rows: rowRefs(undecidable)` demonstrably carries every
undecidable row. B misread `slices.MinFunc`'s scope: `lowest` populates only the
scalar `Guard` field, not `Rows`. A's framing — the *atom* has no producer field,
on the lint surface where the "never just the first" clause actually lives — is
the accurate statement of the defect. B additionally crosses the runtime/lint
boundary the RDR's disposition table keeps separate.

### CT-2 — Whether the escape-row overlap position is a unilateral reversal (R-11)

- **B (C-1)** ranks this its top finding: the RDR "unilaterally reverses a
  JDR-recorded joint decision … without amending the JDR or RDR 0006".
- **A** does not raise the JDR-reversal framing at all. A's escape-row findings
  (C-6, C-7) attack the *substance* — that escape rows cannot rescue
  `guard_unevaluable` in the shipped kernel, and that a bare escape row is a
  free coverage exemption.

**Evidence is split, and favours A's framing as the more load-bearing.** B's
facts check out: JDR 0001 §JD-14 (`docs/jdr/0001-resolve-kernel-seam.md:333-336`)
still reads "**Decided 2026-08-22: 0003's reading governs** — escape rows are
ordinary participants in both the union and the overlap check", and RDR 0006
(`0006:912`) independently records "Decided: RDR 0003's reading governs; repair
invariants 3 and 4 to match the Load-Bearing Decision". Neither has been amended.
But B's charge of *unilateral* reversal is directly adjudicated by the RDR's own
text (see §4, R-11), which books the correction as a Pending assumption with a
named plan rather than asserting it as settled. What B does *not* address, and A
does, is that the substantive premise A17 rests on — the kernel's escape/ordinary
split — is only half the story: the same kernel makes escape rows unable to
rescue two of four refusal classes, which cuts against the coverage clause A17
leaves untouched. A's C-6 is the finding that actually decides the resolve here;
B's C-1 is a real but procedurally-handled process defect.

### CT-3 — Which section gets rewritten within six weeks

- **A** names the declaration-model clauses (`:741-861`) — newest text, absorbed
  at Stage 6, defaults unusable as composed, kind table missing real cases.
- **B** names RDR 0006's row-group and predecessor-reachability contract, and
  this RDR's clauses that assume A10/A12 close cleanly.

**Not strictly contradictory** — these are different sections, and both
predictions can hold. **Evidence favours A** on likelihood: A's argument is
grounded in an arithmetic defect verifiable today (V-2, R-05) that changes
verdicts for every previously-green model, whereas B's rests on predicting how a
peer document's unstarted refine will land. B's is the better-argued *process*
risk; A's is the better-evidenced *contract* risk.

### CT-4 — Which assumption will not survive first contact

- **A** names A11 (finite domain ⇒ every claim has a producer).
- **B** names A15 (a scalar cardinality bound predicts provability).

**Not contradictory, and both are well-argued.** They are also connected in a way
neither pass notices: R-02's `set` miscount means the bound does not measure the
assignment space at all, so A15's "does cardinality predict provability?"
question is asked of a number that is itself the wrong quantity. A's A11 attack
is the stronger of the two as written — A15 is at least stamped `Pending` with a
named discharge, whereas A11 is stamped `Verified` on evidence that is entirely
an ownership argument (see R-18).

## 4. Re-raises of already-adjudicated material

Rows below are ones the RDR's own text explicitly addresses. Adjudication does
not automatically make a row wrong — but each must be argued *against* the
adjudicating text, and neither pass does so.

### R-11 (B's C-1) — escape-row overlap "unilateral reversal" — **ADJUDICATED**

The Contradiction Check (`:1897-1907`) addresses this head-on:

> **One gate decision is partially re-opened, deliberately and on evidence.**
> JDR 0001 §JD-14 decided this RDR's escape-row reading governs on both the
> coverage union and the overlap check. The coverage half stands. The overlap
> half reasoned from "RDR 0001's runtime refusal of ambiguity", and the shipped
> kernel refutes that premise … This RDR now checks overlap in two populations
> and books the cross-document confirmation as **A17**; it does not re-decide
> §JD-14 unilaterally.

A17's Plan line (`:543-546`) further states: "carry to the RDR 0006 refine that
already owes §JD-14's invariant 3/4 repair, and record the correction against
§JD-14 at the next JDR 0001 touch", with an "If wrong" reverting the clause to a
single population. The `authority` table (`:1100`) labels the row
"**Contested — A17**".

B's row asserts "nothing in JDR 0001 or RDR 0006 has actually been amended",
which is factually true, and characterises the RDR as "describing itself
*disagreeing* with a joint decision". The RDR's answer is that it books the
disagreement as an open Pending assumption with a named forcing plan, which is
the mechanism this process provides for exactly this case. B quotes the
"partially re-opened" line but dismisses it as "doing a lot of unearned work"
without engaging the A17 record, the Plan, the If-wrong, or the `Contested`
label. **Downgraded to SUBSTANTIVE on that basis** — the residual defect is real
(the peers still read the other way today, and two implementers can still
diverge until the refine lands) but it is a scheduling risk the RDR has
recorded, not an unnoticed contradiction.

### R-14 (A's C-12, B's C-2) — A10/A12 deferred to RDR 0006's refine — **PARTIALLY ADJUDICATED**

Both records are stamped `Pending` with explicit Plans and are covered by the
"standing tolerance" mechanism, and RDR 0006 acknowledges the debt at its own
`0006:919`: "Still owed to RDR 0003 as standing tolerances (its A10, A12)". The
RDR discloses A12's two-decision-procedure problem in A12's own evidence line —
both passes quote it *from* the RDR, which is the record working as intended.

Not fully adjudicated, and the rows survive: what the RDR adjudicates is that the
questions are *open and tracked*; neither the RDR nor either peer states which
procedure wins if graph reachability and syntactic decision diverge, while the
owned-tag clause (`:1075-1078`) remains a normative MUST quantifying over "every
reachable predecessor". A normative MUST with two candidate decision procedures
and no arbiter is a live defect regardless of the tolerance. **Severity held at
SUBSTANTIVE.**

### R-06 (A's C-11) — withheld-claim MUST has no producer — **ACKNOWLEDGED, NOT ADJUDICATED**

The RDR concedes the gap in its own block quote at `:977-984`:

> RDR 0006's finding contract … has **no atom-level field**
> (`docs/rdr/0006-graph-lint-authority-and-guarantees.md:363-367`). §JD-4 records
> the atom-level extension as a decided duty on RDR 0006, discharged at its
> refine. Until that lands the atom-naming half of the clause above has no
> producer; the row-naming half is satisfiable today. This does not gate this
> RDR's lock — a missing producer field surfaces as an RDR 0006 refine item.

This is disclosure, not adjudication: the RDR states the MUST is unsatisfiable
and then declares that fact non-gating, without qualifying the MUST itself. A's
row is precisely that a normative MUST no consumer can satisfy is not a contract.
**Severity held at BLOCKING** — the clause remains unqualified normative text
while its own document says it has no producer. MVV Scenario 8, the RDR's
flagship false-green defense, asserts the finding "names that row and atom".

### R-16 (A's C-13, B's C-10) — illustrative code disclaims its own spelling — **DISCLOSED**

The RDR marks the block "illustrative only … the field's authoring location is
RDR 0002's to fix (A16)" at `:1243-1245`, and A16 is a live Pending record. Both
passes quote the disclaimer, so neither is unaware of it. The rows survive
because a disclaimer is not a substitute for a normative example when the
document owns a five-field type system and this is its only worked declaration —
but the severity is capped by the disclosure. **SUBSTANTIVE, not blocking.**

### R-01, R-02, R-03, R-04, R-05 — **NOT adjudicated anywhere**

I checked §Contradiction Check (`:1874-1917`), the `trace` closing notes
(`:1147-1158`), the `authority` table (`:1081-1104`), and every assumption record
for any text addressing: the fixture's `profile = small` gap; the `set`
cardinality quantity; `escapeOrRefuse` unreachability for `guard_unevaluable`;
an undecidable `unless` block; or the composition of the two conservative
defaults. **None is mentioned.**

The trace's closing paragraph does name two clauses as "unexercised by this
fixture rather than gapped" — the escape-row overlap population and the
single-valued marker — with a named test for each. Neither is the coverage gap,
and the paragraph's existence shows the pass that wrote it was looking for
unexercised clauses while missing an uncovered *assignment* in the group it had
just stamped OK.

The five blocking rows are all first-raised, all verifiable against artifacts
already in the repo, and none has a standing disposition.
