Model: claude-opus-5

# Critique — RDR 0011 (flow next: match-conditioned candidates)

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0011:C1` — "An `indeterminate` row therefore SURVIVES into `gate` as a non-selectable row, is carried alongside the guard-undecidable rows" | `0007:C10`'s resolution-level veto is never cited anywhere in 0011 (0 occurrences of `0007:C10`; `0007:C11` cited 5×). `gate` is table-scoped: `resolve.go:558` `if len(undecidable) > 0 { return nil, &Refusal{...} }` discards the built-up `selected`. One indeterminate row therefore refuses the ENTIRE resolve, destroying a co-resident `GuardTrue` row's plan. On `models/rdr.toml` the `Recognized` filter leaves 7–8 co-resident rows per resolve, every one a chance to nuke a valid plan. | `flow resolve --outcome advance` refuses `guard_unevaluable` naming a rule the user never asked about, on a state where it previously produced a correct plan. Every advance/abandon/revise on the whole model dies the moment one `stage` key goes unbound. | §1, §2, premortem, AT-1 |
| C-2 | `0011:C1` — "an ordinary row whose match verdict is `indeterminate` MUST NOT be selected as a plan, and MUST NOT count toward the `len(selected)`" | `gate` has no representation for "survives but is non-selectable". Its partition is exactly `{GuardFalse → pruned, GuardTrue → selected, GuardUnevaluable → veto}` (`resolve.go:521-556`). `gateMatchRow()` (`match_conflicted_test.go:45`) is **guard-free**, so `evaluateAtoms` returns `GuardTrue` under `0007:C5` regardless of match. The RDR names no third channel, no new field on `Row`, and no signature change to `gate` — Phase 1 says only "survives that filter into `gate` as a non-selectable row". The implementer must invent the mechanism, and every available invention (fold into `GuardUnevaluable`) triggers C-1. | Implementer ships whichever reading compiles first; `flow resolve` becomes non-deterministic across builds, and `S5`'s five-oracle enumeration cannot be checked because the disposition depends on an unspecified choice. | §1, §2, AT-2 |
| C-3 | `0011:JC1` — "the coupling stays one-way … cannot turn a 0010 acceptance into a refusal" | `0010:A9`: the "otherwise" row of a decision table is an **escape row rescuing `no_match`**. `0011:C1` converts an undecidable match into `guard_unevaluable`, which `0002` excludes from escape lists and which `gate` returns BEFORE `escapeOrRefuse` is ever reached (`resolve.go:479-481` short-circuits above `resolve.go:483`). So an unsupplied `--tag` on a 0010 table stops reaching the catch-all. JC1 asserts the exact opposite and calls the coupling "strictly widening". | `flow resolve` on a decision table with a partially supplied tag set refuses `guard_unevaluable` instead of firing the "otherwise" row. 0010's headline feature — the default row — silently stops working, and 0010 has already shipped its MVV against the old behaviour. | §1, §3, premortem, AT-3 |
| C-4 | `0011:A7` — "the five identified by inspection … are the claim's floor, not its verified extent"; Status **Pending** | A7 censuses six files. Nine files in `internal/resolve` build rows with `Match` blocks: the census omits `adversarial_test.go` (7 `Match:` sites), `fixup_test.go` (2), `reserved_key_sameview_0008_test.go` (2, incl. a two-key match at `:68`), `escape_shape_adv_0009_test.go` (1), `resolve_test.go`. A7 is Pending while `S5`, `MVV` step 8 and Consequences all state "and no others" as settled fact — the exact Status-consistency breach `§Assumption Verification` forbids. | `make check` goes red across five unnamed files on the implementation branch. The implementer cannot tell an expected inversion from a regression, and re-expects them all to get green. | §1, premortem, AT-4 |
| C-5 | `0011:S5` vs `0011:MVV` step 8 / `0011:§mini-checks` trace step 8 | S5 says "**Five** shipped oracles change disposition" and lists five. MVV step 8 and the trace table both say "the **four** oracles S5 names (three in `match_conflicted_test.go`, one in `reserved_key_0008_test.go`)" — excluding `TestReq78`, which S5 lists as retired. Two normative counts of the same set inside one locked record. | The build asserts the wrong arity in its own gate check; whichever count the implementer keys off, the other clause is violated at lock. | §1, AT-4 |
| C-6 | `0011:S5` — "`TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues` — **re-expected**" | That test is by its own docstring "the **positive control** for the test above" (`match_conflicted_test.go:190-194`). Re-expecting it to a `guard_unevaluable` refusal deletes the only oracle proving that a clean escape row still rescues when a conflicted key is present — i.e. that the fix pruned ONLY conflicted matches. 0011 destroys the control while keeping the thing it controls for. | A later regression that over-prunes escape rows ships undetected: users lose escape rescue on models that have nothing to do with conflicted keys, with no test to catch it. | §1, §2, AT-5 |
| C-7 | `0011:A8` — "widening `guard.go::isGuardBlock`'s filter … admits `BlockMatch` entries while leaving every existing guard payload byte-identical"; Status **Pending** | `isGuardBlock` is not a payload filter. In `evaluateAtoms` (`guard.go:120-155`) its `continue` is a single early-exit gating THREE things: the seam call in `evaluateAtom`, the payload append, and the `kleeneAnd` operand contribution. Its own comment forbids the edit A8 proposes — sweeping such an atom in "would let it be DECIDED, and a decided-FALSE one would prune its row along with the row's owned-state obligation (D8) — reopening the masking path `0007:C11` ratifies pruning against" — and names `BlockMatch` as the reachable vector. A8 cites `0007:C11` in support of it. Widening also crosses RDR 0003's typed-operator boundary: match atoms would reach `seam.Evaluate` at `guard.go:219`. | Match atoms reach the guard seam and the owned-state masking defect 0007 closed reopens. `guard_adversarial_0007_test.go::TestAdv0007_1` / `::TestAdv0007_2` (whose `"§D12 match block"` arm asserts a `BlockMatch` atom must NOT prune) and `guard_fixup_0007_test.go::TestFix0007Fail1` (whose 4th subtest asserts "a match atom is never handed to the seam", `len(calls) != 0` ⇒ fail) go red — none named in S5. | §1, §2, AT-6 |
| C-8 | `0011:C3` — "The build MUST **mechanically** re-home every shipped read of the `unresolved` key to `unknown` … six TEST reads" | The rename is not mechanical. All six reads go through `stringsAt` (`flow_harness_0005_test.go:426`), which does `e.(string)` per element and returns `nil, false` on a `{key,reason}` object. Two of the six call sites **discard the ok bool**: `flow_adversarial_0005_test.go:446` and `flow_mvv_0005_test.go:66` are `unresolved, _ := stringsAt(c, "unresolved")`. The adversarial one is a NEGATIVE assertion (`if containsString(unresolved, "flag")`) and passes **vacuously** on the empty slice. | The suite goes green while the oracle is dead. The specific assertion that `next` must not report a fact a bound reader holds stops testing anything — and C3 has already pre-declared a green suite as the pass criterion for this class of change. | §1, §2, AT-7 |
| C-9 | `0011:C1` — "the CLI reads the match facts off it exactly as it reads guard facts today, and no new field on `Plan` is required" | `summarize(row, view, gatesRan)` (`flow_next.go:168`) receives no kernel output at all, and `excluded(row, owned, observed) bool` (`flow_next.go:276`) throws the `Result` away. To carry match facts, `excluded` must return the `Refusal` and `summarize` must take it — a signature change on both, plus the `runFlowNext` loop at `flow_next.go:114-129`. The Existing Infrastructure Audit marks both rows "Extend" but names neither signature; Phase 2 says only "carry the kernel's undecided match atoms into the candidate's `unknown`". | Implementer keeps `excluded`'s bool, re-derives match facts CLI-side from `assembledView`, and thereby reintroduces Alternative 3 — the design the RDR spent three sections rejecting. Conflicted keys report `absent`, which is the wrong remedy. | §1, §2, AT-8 |
| C-10 | `0011:C1` compound clause — "a CONFLICTED match key … MUST be reported with reason `absent` rather than omitted" | The RDR calls this "the ONLY case in which `unknown` does not distinguish `uncomparable` from `absent`". It is not the only case, given C-9: whenever the implementer takes the natural CLI-side path the collapse is universal. Worse, `absent` and `uncomparable` are stated (D-undecided-vocabulary) to have *different remedies* — "bind a reader or supply the tag" vs "fix the overlapping readers". Reporting `absent` for a conflicted key tells the user to supply a tag that is already supplied twice. | User binds the reader / supplies the tag as instructed, re-runs, gets the identical `unknown` entry, and concludes the tool is broken. The actual remedy (two readers writing one key) is never surfaced. | §1, §3, premortem, AT-9 |
| C-11 | `0011:A3` (Status **Verified**) — "The conflicted exception is real and user-reachable … two readers writing the same key with differing values conflicts it"; `0011:C3` — the fixture "MUST be built from two invoked readers writing the same owned key with differing values" | **Refuted at model load.** `internal/table/load.go:360-375` `checkAccessorBindings` refuses any model where an owned tag is served by ≠1 reader: `if readerCount[key] != 1 { return fail(CatMalformedAccessorBinding, "owned tag %s is served by %d readers; want exactly one") }`. A3 argues reachability from the absence of dedup in `runReaders` — true but irrelevant, because no model that could exercise it will load. `internal/accessor/executor.go:87-105` closes the second route: a reader's output is classified against `def.RequestedKeys()`, so it cannot emit a key it did not declare. With `--tag` refused by `parseTags` (`flow-tag-duplicate`, per the RDR itself), `tv.conflicted` is **unreachable from the CLI on this build**. | The fixture C3 makes mandatory cannot be constructed. `uncomparable` never ships an end-to-end oracle. Worse, the load-bearing justification collapses: Alternative 3 is rejected because "the failure is real, not theoretical", and `D-placement` rejects the CLI-side option because "conflictedness is not visible to the CLI" — for a condition no CLI user can produce. | §1, §3, premortem, AT-10 |
| C-12 | `0011:§consequences` — "Three shipped `internal/resolve` oracles change disposition"; `0011:A7` says five; `0011:S5` says five; `MVV` 8 says four | A third count, in the section that describes blast radius to a reviewer. The same paragraph calls this "the largest behaviour change in the RDR". | Reviewer sizing the risk reads "three"; implementer discovers five-plus-five-unnamed. The Profile (`large`) was justified on contract count, not on this number, so nothing else catches the undersizing. | §1, AT-4 |
| C-13 | `0011:C2` — "`--all` MUST NOT be accepted by flow resolve, flow read-state, or flow set-state — satisfied by NON-REGISTRATION" | Non-registration is not a contract; it is the absence of one. `--all` is registered "on `next` only, not on the shared `registerSelectionFlags`" — but `registerSelectionFlags` is a shared helper any future verb-add will reach for, and nothing pins the negative. The RDR even routes the assertion to "cobra's shipped unknown-flag path", asserting on a third-party error class. | The day someone moves `--all` to the shared registrar for symmetry, `flow resolve --all` silently starts ignoring match at a SELECTION site — a write-performing plan on an unchecked precondition, the exact masking C1 forbids. | §2, AT-11 |
| C-14 | `0011:A5` — "a genuine zero-owned decision table is **unloadable on this build**"; `0011:JC1` — "the 'empty `required`' half is downstream of 0010 shipping, not of this contract" | The RDR's own key coupling to 0010 is verified against a fixture the RDR admits is not the real thing (it "carries one dummy owned tag and `required` is `['answer']`, never empty"). The one case where C1's absent-key provision is supposed to earn its keep — the no-tag decision table — is validated on a shape that cannot exist and a shape that will exist is never tested. | 0010 ships, the real zero-owned table appears, and `flow next` over it behaves in a way no oracle in either RDR covers. The user's first decision table is the experiment. | §3, AT-12 |
| C-15 | `0011:§problem-statement` — "Success is therefore 'no row is listed that the supplied state excludes', not a target cardinality" | The success criterion is unfalsifiable by the user. The RDR's own second-order admission — "that second case is the old symptom's shape with a diagnosis attached, and it is the intended behaviour, not a residual defect" — means the reported defect (21 of 22 rows) is still a legal output. Meanwhile MVV 3 pins the number `3`, and the whole justification rests on `models/rdr.toml`, where every rule matches on `stage` and `stage` is `required = true` and reader-served. The design is validated exclusively on the one model where it cannot fail. | On any model whose match keys are observed rather than owned, the user gets the identical wall of candidates that motivated the RDR, plus a `unknown` list per row. They report the bug again. It is closed as working-as-designed. | §2, §3, premortem, AT-13 |
| C-16 | `0011:§consequences` — "a candidate can still be listed that `resolve` would refuse (an undecidable match key)" | Buried as a one-line Negative, this is now the *primary* relationship between the two verbs given C-1: `next` reports the row as a candidate, and `resolve` on that same state refuses the whole table. The RDR frames this as "the same relationship `next` already has with `guard_unevaluable`" — but under `0007:C10` that relationship is a table-wide veto, not a per-row one, which the RDR never states. | A skill reads a candidate from `flow next`, invokes `flow resolve` with it, and gets a refusal naming a different rule. The skill has no way to distinguish "my candidate was wrong" from "some other row poisoned the table". | §1, premortem, AT-1 |
| C-17 | `0011:A3` "If wrong" — "the kernel already carries an undecided match verdict, and C1's Phase 1 is a reporting change rather than a semantics change" | A3 is stamped **Verified** by `Source Search`, but its user-reachability claim is refuted by a file the search never opened (`internal/table/load.go`). A3 is the sole evidentiary basis for placing the verdict in the kernel: it is cited in the Approach ("Three facts force that placement"), in `D-placement`, in the Alternative 3 rejection, and in `§Decision Rationale`. With the conflicted case unreachable from the CLI, the only remaining discriminator between the chosen design and Alternative 3 is the ABSENT-key case — which `assembledView` **can** see (A2 says so outright: "a match key the view lacks lands in the candidate's undecided list today"). | The entire kernel-seam change — the `large` Profile, the `flow resolve` blast radius, C-1 through C-7 — is bought to distinguish a case no user can reach. A CLI-scoped change (Alternative 3) would have delivered the user-visible narrowing with none of the veto risk. | §2, §3, premortem, AT-10 |
| C-18 | `0011:C1` — "carrying its undecided atoms in `Refusal.Undecided` — the same payload"; `0011:A8` — "admits `BlockMatch` entries" | `Refusal.Undecided` is `[]UndecidedRow` (`resolve.go:375`), not `[]UndecidedAtom`; atoms ride at `Undecided[i].Atoms`. `UndecidedRow` is keyed by `(RuleID, SourceLocator)` and, per `0007:C8`, is emitted **only** for rows whose GUARD verdict is `GuardUnevaluable` — `evaluateAtoms` returns `nil` payload early at `guard.go:162-166` for any decided row. So a row with a decided guard and an indeterminate match produces **no `UndecidedRow` at all**, and the match atoms have no carrier. | The carrier C1 designates is empty in the common case — a guard-free row (like `gateMatchRow`) with an unsupplied match key. `unknown` reports nothing for exactly the rows the RDR exists to diagnose. | §1, §2, AT-6 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The resolution-level veto — `flow resolve` dies on the flagship model

**Root cause in the RDR.** RDR 0011 cites `0007:C11` five times and `0007:C10` zero times. C11 ratifies *pruning*; C10 governs what happens to a row that is **not** pruned. C10 is unambiguous:

> "Aggregation is resolution-level: if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists. This veto is the cost the RDR accepts: one unreadable row refuses a table whose other rows decide cleanly."

0011 designs directly into that veto without naming it.

**The enabling passage.** `0011:C1`:

> "An `indeterminate` row MUST NOT be pruned there: it is by definition undecided, so dropping it upstream would take its owned-state obligation with it (D8) and reopen the masking path `0007:C11` ratifies pruning against. An `indeterminate` row therefore SURVIVES into `gate` as a non-selectable row, is carried alongside the guard-undecidable rows, and its match atoms join that row's `Undecided` entry."

"Carried alongside the guard-undecidable rows" is exactly membership in the set C10 vetoes on. Phase 1 repeats it verbatim.

The shipped code makes the consequence mechanical. `internal/resolve/resolve.go:540-565`:

```go
	for i, row := range survivors {
		switch verdicts[i] {
		case GuardTrue:
			selected = append(selected, row)
		case GuardUnevaluable:
			undecidable = append(undecidable, row)
			...
		}
	}

	if len(undecidable) > 0 {
		slices.SortFunc(undecided, compareUndecidedRows)
		return nil, &Refusal{
			Kind:      KindGuardUnevaluable,
			...
		}
	}
```

`selected` is built and then **thrown away** — `return nil, &Refusal{...}`. `Resolve` then short-circuits at `resolve.go:479-481` before the `switch len(selected)` and before `escapeOrRefuse` is reachable at all.

`missingOwned` (`resolve.go:532-538`) is table-scoped in the same way, and runs *first*: it unions `RequiresOwned` across every survivor. Under 0011, an indeterminate row that carries `RequiresOwned` keys the snapshot lacks makes the whole resolve refuse `owned_state_unavailable` naming keys belonging to a row that does not even match. 0011 predicts this only for the compound single-row probe (S7) and never for the multi-row real table.

**Blast radius, measured.** `models/rdr.toml` has 22 rules over 3 outcomes: 8 `advance`, 7 `abandon`, 7 `revise`. `row.Outcome != in.Recognized` (`resolve.go:469-471`) narrows 22 → 7–8. **Every rule carries `[rule.match.stage]`** (44 `rule.match` hits over 22 rules; A6 confirms "Every rule carries a `[rule.match.stage]` block"). So the residual match pattern after the Recognized filter is *exactly* the stage key — the one key whose indeterminacy is the RDR's motivating scenario. Seven to eight co-resident rows per resolve, all keyed on the same thing, any one of which now vetoes the table.

The codebase already documents this hazard, in a file A7's census scopes out. `internal/resolve/reserved_key_sameview_0008_test.go:76-80`:

> "any other value reaches the seam's default and answers GuardUnevaluable, **which would refuse the resolve rather than prune the row**."

**Symptom the user sees.** A skill author runs `intrastate flow resolve --model models/rdr.toml --outcome advance --artifact rdr=./0011.md` on a state that resolved cleanly yesterday. Today the `rdr-status` reader is misconfigured — or a new rule was added whose `stage` value is spelled differently — and the command refuses `guard_unevaluable`, naming `resolve-abandon`, a rule the user did not request and whose outcome is not even `advance`. The plan that `finalize-pass` would have produced never appears. Because `guard_unevaluable` is non-escapable by design, there is no escape row that rescues it and no flag that recovers it. Every advance, abandon, and revise on the model is dead until the reader is fixed. The verb the RDR was not changing is the verb that breaks.

### 1.2 `gate` has no channel for "survives but is non-selectable" — so the implementer invents one

**Root cause in the RDR.** `0011:C1` asserts a state that the shipped kernel cannot represent:

> "an ordinary row whose match verdict is `indeterminate` MUST NOT be selected as a plan, and MUST NOT count toward the `len(selected)` that yields `ambiguous_match` against a row that genuinely matched. It instead refuses `guard_unevaluable`."

`gate`'s partition is total over three values and has no fourth: `GuardFalse` → `continue` (pruned, `resolve.go:524-526`); `GuardTrue` → `selected`; `GuardUnevaluable` → veto. There is no "carried, non-selectable, non-vetoing" bucket, no field on `Row` to mark one, and no signature on `gate` that could return one. The RDR proposes no new type, no new field, and no new function signature for `gate` anywhere in C1, Phase 1, or the Existing Infrastructure Audit (whose `gate` row does not exist).

**Why the obvious implementation is worse.** The natural move is to fold the indeterminate match into `GuardUnevaluable` so it rides the existing `Undecided` payload — which is what "its match atoms join that row's `Undecided` entry" literally says. That is C-1. The alternative — a genuinely new bucket — is unspecified work the RDR does not scope, requiring `gate` to return a fourth partition and `Resolve` to reason about it before `len(selected)`.

And the fixtures make the gap concrete. `gateMatchRow()` (`match_conflicted_test.go:45-58`) is **guard-free**:

```go
	return resolve.Row{
		RuleID:        "A",
		Match: []resolve.Tag{
			{Key: "status", Value: "Draft"},
			{Key: dupKey, Value: "closed"},
		},
		RequiresOwned: []string{"status"},
		...
	}
```

Its `Guard` slice is empty. Under `0007:C5`, "A row with no atoms is decided TRUE without consulting the seam." So `evaluateAtoms` returns `GuardTrue` for row A no matter what its match verdict is. The RDR's "carried alongside the guard-undecidable rows" has no mechanism to reach row A at all — the row 0011 names by name in S5 as one of its five changed oracles.

**Symptom the user sees.** Nondeterminism across builds. Whichever reading the implementer picks, `flow resolve`'s disposition on an indeterminate match differs, and `S5`'s "any OTHER oracle going red is a defect" cannot be evaluated because the expected dispositions were never derivable from the record. The user sees `flow resolve` refuse differently after an unrelated refactor, with no changelog entry, because the behaviour was never pinned.

### 1.3 0010's "otherwise" row stops firing — the coupling JC1 declares impossible

**Root cause in the RDR.** `0011:JC1` states:

> "the coupling stays one-way — C1 makes an undecidable match atom a named fact instead of an escapable `no_match`, which strictly widens what a zero-owned decision table can report and cannot turn a 0010 acceptance into a refusal."

`0010:A9` states the shape that claim ignores:

> "The 'otherwise' row of a decision table is an **escape row rescuing `no_match`** that carries `[rule.emit]`."

The mechanism is settled in shipped code. `escapeOrRefuse` is only reachable from `Resolve` at `resolve.go:483` — strictly **below** the `if blocked != nil { return refuse(in, *blocked), nil }` short-circuit at `resolve.go:479-481`. A `guard_unevaluable` from `gate` never reaches the escape phase. And `0007:C8` puts `guard_unevaluable` outside escape lists by construction: "The kernel maps GuardUnevaluable to `guard_unevaluable`, which RDR 0002 excludes from escape lists."

So the transformation JC1 calls "strictly widening" is precisely: **`no_match` (escapable, rescued by the otherwise row) → `guard_unevaluable` (non-escapable, never reaches the escape phase)**. That converts a 0010 acceptance into a refusal. JC1 asserts the negation.

There is a second regression in the same path. `escapeOrRefuse` (`resolve.go:620-623`) drops the original refusal:

```go
	viable, blocked := gate(escapes, in.Guards, view)
	if blocked != nil {
		return refuse(in, *blocked)
	}
```

`*blocked`, not `r`. The `case 0: return refuse(in, r)` line that preserves the caller's original `no_match` is unreachable once `blocked != nil`. Under 0011, an escape row with an indeterminate match — and `guard_fixtures_test.go::conformingEscapeRow` carries `Match: status=Draft`, as A7 itself notes — replaces the actionable `no_match` with a `guard_unevaluable` about a row the caller never asked about.

**Symptom the user sees.** A user builds a 0010 decision table with an `otherwise` catch-all — the entire point of the feature. They run `flow resolve` supplying two of three tags. Previously: no ordinary rule matched, `no_match`, the otherwise row rescued, `emit` fired, `escaped = true`. Now: `guard_unevaluable`, no rescue, no emit, and a refusal naming the unsupplied tag. The default row that exists specifically to handle incomplete input is the one thing incomplete input can no longer reach. 0010's MVV was written against the old behaviour and `0010:A11` is still `Pending`.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`0011:C1`, and specifically the paragraph beginning "The two verdicts are NOT applied at the same place, and the split is forced by `0007:C11`."**

That paragraph is the entire kernel design. It will be rewritten because it is the only place C-1, C-2, C-6, and C-16 all originate, and because it is normatively wrong about the clause it cites.

The paragraph's argument runs: pruning an indeterminate row upstream "would take its owned-state obligation with it (D8) and reopen the masking path `0007:C11` ratifies pruning against", therefore the row must survive into `gate`. But `0007:C10` — the immediately preceding contract, which 0011 never cites — specifies what surviving into `gate` *costs*: a resolution-wide veto. The RDR reaches for C11 to justify survival and never reads C10 to price it. Six weeks in, the first user whose `stage` reader breaks will file a bug that `flow resolve` refuses the entire model, and the fix requires either:

1. **Row-scoping `gate`'s blocking** — which contradicts `0007:C10` outright ("a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists") and therefore requires a new RDR overriding C10, not an amendment here; or
2. **Reverting to pruning indeterminate rows at the match filter** — which is Alternative 1, which this RDR rejects, and which re-collapses `absent` into exclusion; or
3. **A fourth partition in `gate`** — the unspecified mechanism of C-2, which is real new kernel design and needs its own contract.

All three are RDR-scale changes to the paragraph that is the sole normative home of the selection-site rule. None can be done by amendment; the house rule is that RDR content is never amended. So this becomes RDR 0012, and 0011's C1 becomes an overridden clause within six weeks of shipping — the same fate 0011 hands `0005:C1`, on a shorter clock and for the same reason: the record settled a predicate without pricing the seam it fires into.

Three secondary rewrites ride along. `0011:S5`'s oracle enumeration will be rewritten because it is arithmetically inconsistent with MVV 8 and Consequences (C-5, C-12), and because it names five oracles from a six-file census over a nine-file population (C-4). `0011:A8`'s carrier design will be rewritten because widening `isGuardBlock` does not do what A8 says it does (C-7) — `guard.go:120-141`'s own comment states that admitting a non-guard block "would let it be DECIDED, and a decided-FALSE one would prune its row along with the row's owned-state obligation (D8)", which is the defect A8 cites `0007:C11` to avoid — and because the carrier it designates, `[]UndecidedRow`, is structurally empty for a guard-free row (C-18).

**And there is a version of six-weeks-later in which C1 is not rewritten but deleted.** If C-17 lands before the veto does — if someone tries to build C3's mandatory conflicted fixture and hits `load.go:371` — then `0011:A3` falls, and with it the stated ground for kernel placement in `§approach`, `D-placement`, and the `ALT3` rejection. What replaces C1 is not a corrected kernel contract but Alternative 3: a CLI-scoped presence rule over `assembledView`, which `0011:A2` already establishes is sufficient for every case a user can produce. That rewrite is larger than the first one and lands on the same paragraph, which is why the paragraph is the answer either way.

---

## 3. The one assumption that will not survive first contact with a real user

**`0011:A3` — "The conflicted exception is real and user-reachable."** Stamped **Verified**, method **Source Search**. It is false, and it is the assumption the entire kernel-seam design is bought with.

A3's argument for reachability is:

> "The reachable producer is the OWNED provenance — `internal/cli/flow_exec.go::flowRequest.runReaders` appends each invoked reader's `OwnedSnapshot()` into one slice with no dedup, so two readers writing the same key with differing values conflicts it."

The absence of dedup is real. The premise it is attached to — that two readers can write the same owned key — is refused at model load, in a file A3's Source Search never opened. `internal/table/load.go:360-375`:

```go
	for _, key := range slices.Sorted(maps.Keys(l.model.Tags)) {
		decl := l.model.Tags[key]
		switch decl.Provenance {
		case ProvenanceOwned:
			if readerCount[key] != 1 {
				return fail(CatMalformedAccessorBinding,
					fmt.Sprintf("owned tag %s is served by %d readers; want exactly one",
						key, readerCount[key]))
			}
```

**Exactly one reader per owned tag, enforced at load.** A model in which two readers serve one owned key does not load. The second route is closed too: `internal/accessor/executor.go:87-105` classifies a reader's output against `def.RequestedKeys()` — "Deriving it from what came back makes completeness self-fulfilling" — so a reader cannot emit a key it did not declare, and `OwnedSnapshot` renders only those. The third route, a repeated `--tag`, is refused by `parseTags` with `flow-tag-duplicate`, which the RDR documents itself. On this build, `tv.conflicted` is **unreachable from the CLI**.

Now trace what rests on it. A3 is the load-bearing citation in four places:

- **`§approach`**: "Three facts force that placement. (i) The CLI cannot implement the rule: conflictedness is `internal/resolve`'s `TagSet.conflicted` … a CLI-side presence test would retain a conflicted atom and the kernel would exclude the row as `no_match` with nothing in `unknown`, which is the masking `P1` forbids."
- **`D-placement`**: "Rejected: the CLI-side pre-strip (Alternative 3), because conflictedness is not visible to the CLI and a presence-only test silently excludes a conflicted row."
- **`ALT3` Reason for rejection**: "Verified in Stage 4 (A3): the conflicted-key case is user-reachable … so the failure is real, not theoretical."
- **`§decision-rationale`**: "reachable by a caller who repeats a `--tag`" — an inaccuracy the record corrects elsewhere but leaves standing here.

Strip the conflicted case and the discriminator between the chosen kernel design and Alternative 3 is the **absent**-key case alone. But `0011:A2` establishes, Verified, that the CLI already handles absence:

> "`summarize` walks `Row.Atoms` … so a match key the view lacks lands in the candidate's undecided list today."

So the CLI *can* see every case a user can actually produce. The one case it cannot see is the one no user can reach. The RDR spent its `large` Profile, its kernel-seam blast radius, its `flow resolve` behaviour change, and every defect in §1 of this critique to distinguish a condition that does not occur.

**First contact.** The implementer sits down to build C3's mandatory fixture — "two invoked readers writing the same owned key with differing values", with the `--tag` route explicitly banned — and the model will not load. They will do one of three things: write the banned `--tag` version (which C3 says "MUST NOT be accepted as this fixture"), construct a `TagSet` directly in-package (which tests the kernel but not the user journey, and leaves `uncomparable` with no end-to-end oracle), or file the RDR as unimplementable. None of these is a path the record anticipates, and the third is the correct one.

**The runner-up, which fails the same way.** `0011:A6`'s narrowing is verified on `models/rdr.toml`, and A6 states plainly what that model cannot test:

> "Every rule carries a `[rule.match.stage]` block and `stage` is `required = true` and served by the invoked `rdr-status` reader, so **C1's absent-key provision never fires for this model**."

The 21 → 3 result, MVV 3's pinned cardinality, the Decision Rationale, and the trace table's witnesses are all validated on the one model where the RDR's central new machinery is provably unreachable. Where it *is* reachable, A5 records that no narrowing happens ("every row remains a candidate"), and the Problem Statement pre-absolves that outcome as "intended behaviour, not a residual defect" — making the success criterion unfalsifiable by the user who filed the original defect. A5's own fixture is likewise a substitute: "a genuine zero-owned 0010 decision table is **unloadable on this build** … the fixture carries one dummy owned tag and `required` is `['answer']`, never empty."

Between them, A3 and A6 mean the RDR is validated exclusively where it cannot fail, and justified exclusively by a case that cannot occur.

---

## 3b. The runner-up assumption, stated in full

**That the user's match keys are owned and reader-served — the assumption `0011:A6` verifies and every other section then treats as the general case.**

A6 is verified, live, on `models/rdr.toml`, and its result is genuine: 21 → 3. But read what A6 actually establishes:

> "Every rule carries a `[rule.match.stage]` block and `stage` is `required = true` and served by the invoked `rdr-status` reader, so **C1's absent-key provision never fires for this model**."

The load-bearing narrowing claim, the MVV's pinned cardinality of `3`, the Decision Rationale, and the trace table's witnesses are all validated on the one model where the RDR's central new machinery — the three-valued verdict, the `indeterminate` value, the `unknown` list, the whole reason vocabulary — is provably unreachable. The RDR validates the happy path on a fixture that cannot exercise the unhappy one.

Where the machinery *does* fire, A5 is candid that the narrowing does not happen:

> "where a match key is observed and unsupplied (0010's decision-table class with no `--tag`, A5) every row remains a candidate and each names the undecided key under `unknown`."

And the Problem Statement pre-declares that outcome acceptable:

> "That second case is the old symptom's shape with a diagnosis attached, and it is the intended behaviour, not a residual defect."

So the RDR's success criterion — "no row is listed that the supplied state excludes" — is satisfied by an output identical to the reported defect. The first real user with an observed-key model gets 21 candidates and a per-row `unknown` list, which is strictly more output for the same non-answer. They will report the bug again, and the record already contains the sentence that closes it as working-as-designed.

Compounding it, A5's own verification is on a fixture the RDR admits is not the real shape:

> "a genuine zero-owned 0010 decision table is **unloadable on this build** … so the fixture carries one dummy owned tag and `required` is `['answer']`, never empty."

The one class where the absent-key provision earns its keep is verified against a shape that cannot exist, and the shape that will exist is deferred to "downstream of 0010 shipping."

First contact: a user writes a model whose match keys come from `--tag` rather than a reader, runs `flow next`, sees the full alphabet, and concludes the RDR did not ship. They are looking at correct, contracted, intended behaviour. There is no flag that helps, because `--all` widens rather than narrows. The remedy in the payload — "bind the reader, supply the tag" — requires them to supply, by hand, the state whose derivation was the reason they ran `flow next`.

---

## 4. Premortem — written from six weeks after the ship

The rollout looked clean for four days. `flow next --model models/rdr.toml` narrowed 21 → 3 exactly as `evidence/spikes/a6-rdr-toml-narrowing.md` predicted, MVV steps 2–4 were green in CI, and the `--all` escape hatch meant nobody complained about the override of `0005:C1`. The `internal/resolve` suite went red on five oracles, we re-expected all five per `S5`, and `make check` came back green.

**Day 5 — the `flow resolve` outage.** A skill author updated `models/rdr.toml` to add a `stage = "reconciled"` rule and, in the same commit, renamed a key in the `rdr-status` reader's `keys` list. `runReaders` returned a snapshot without `stage`. Under the old build this was harmless in the way the record predicted: `TagSet.matches` returned false on `!ok`, every row was pruned at `resolve.go:472`, `escapeOrRefuse` fired, and the modeled escape rescued with `escaped = true`.

Under the new build, `Resolve` kept all eight `advance` rows past the match filter as `indeterminate`, they all survived `gate`'s `GuardFalse` prune, and `len(undecidable) > 0` at `resolve.go:558` fired. `gate` returned `nil, &Refusal{Kind: KindGuardUnevaluable}` — discarding the `selected` slice it had just built, which contained the one row that would have produced the correct plan. `Resolve` short-circuited at `resolve.go:479-481`, and `escapeOrRefuse` was never called.

Every `intrastate flow resolve --outcome advance|abandon|revise` on the model refused `guard_unevaluable`, listing eight rule ids in `Undecided` — including three the caller's outcome could not have selected. `guard_unevaluable` is non-escapable by `0007:C8`, so no escape row helped and no flag recovered it. Twelve RDRs in flight were unable to advance a stage for six hours. The author's first instinct was to add an escape row for `guard_unevaluable`; RDR 0002 refuses that class in escape lists, so the model would not load.

We searched RDR 0011 for the veto and found `0007:C11` cited five times and `0007:C10` cited zero. The paragraph that put us here — "An `indeterminate` row therefore SURVIVES into `gate` as a non-selectable row, is carried alongside the guard-undecidable rows" — read as safe because C11 is about pruning. C10, one contract earlier in the same record, is the one that prices survival: "one unreadable row refuses a table whose other rows decide cleanly." Nobody in three pre-lock lenses read C10.

**Day 6 — the mechanism was never specified.** Fixing it meant deciding what "non-selectable row" means to `gate`, and C1 never said. `gate`'s partition is `{GuardFalse, GuardTrue, GuardUnevaluable}` and none of them is it. Worse, `match_conflicted_test.go::gateMatchRow` — one of the five oracles `S5` names by name — is guard-free, so `evaluateAtoms` returns `GuardTrue` for it under `0007:C5` regardless of its match verdict; the RDR's "carried alongside the guard-undecidable rows" had no way to reach row `A` at all. The implementer had folded indeterminate-match into `GuardUnevaluable` at `evaluateAtoms` because that was the only thing that compiled against "its match atoms join that row's `Undecided` entry."

**Day 8 — decision tables lost their default.** A user on the new 0010 decision-table path ran `flow resolve` supplying two of three tags. The `otherwise` row — `0010:A9`, an escape row rescuing `no_match` — did not fire. The refusal was `guard_unevaluable`. `escapeOrRefuse` had never been reached, because `gate`'s block returns above it. `0011:JC1` says in as many words that this coupling "cannot turn a 0010 acceptance into a refusal." It had. `0010:A11` was still `Pending`, so nothing in either record had ever run the journey end to end.

**Day 9 — the escape control was gone.** Reviewing `S5`, we found we had re-expected `TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues` to a refusal. Its docstring reads "the **positive control** for the test above. Pruning conflicted-key matches must prune ONLY those." We had deleted the only oracle proving a clean escape row still rescues, in the same change that broke escape rescue for decision tables. Nothing caught it because we had re-expected it on purpose.

**Day 11 — the CLI half was Alternative 3 after all.** Chasing a report that conflicted keys said `absent`, we read `flow_next.go`. `excluded(row, owned, observed) bool` throws the `Result` away; `summarize(row, view, gatesRan)` never sees the kernel. C1 said "the CLI reads the match facts off it exactly as it reads guard facts today, and no new field on `Plan` is required" — true of `Plan`, and silent on the two signatures that actually had to change. The implementer had kept both signatures and re-derived match facts from `assembledView`, which is Alternative 3, the design the RDR rejects across three sections. Conflicted keys reported `absent` universally, not only in C1's named compound case. Users bound the reader as the payload instructed, re-ran, got the identical entry, and filed it as a display bug.

**Day 12 — the green suite that tested nothing.** `flow_adversarial_0005_test.go:446` had been passing since day one. It reads `unresolved, _ := stringsAt(c, "unresolved")` — the ok bool discarded — and asserts `!containsString(unresolved, "flag")`. After the rename to a `{key, reason}` list, `stringsAt` hit `e.(string)` on an object, returned `nil, false`, and the negative assertion passed on an empty slice. `flow_mvv_0005_test.go:66` had the same shape. C3 had called the six re-homings "mechanical" and had itself written "A green 0005 suite MUST NOT be cited as evidence that this contract was implemented" — which we honoured for C1 and forgot for the payload rename that C3 said was independent of it.

**Day 13 — the carrier was empty for the rows that mattered.** A user reported that `flow next` over their decision table listed every rule with `"unknown": []` — no diagnosis at all, the pre-0011 symptom with the new field name. We traced it. `Refusal.Undecided` is `[]UndecidedRow` (`resolve.go:375`), not a flat atom list, and `evaluateAtoms` returns `nil` payload for any row whose guard is decided (`guard.go:162-166`). Their ordinary rules were guard-free, so `0007:C5` made every one of them `GuardTrue`, so none produced an `UndecidedRow`, so the match atoms C1 designates as riding "in `Refusal.Undecided`" had nowhere to ride. The carrier the RDR chose specifically so that "no new field on `Plan` is required" is structurally absent for exactly the row shape 0010's decision tables are made of — and `match_conflicted_test.go::gateMatchRow`, one of S5's five named oracles, is the same shape.

**Day 13, later — the fixture that could not be built.** The `uncomparable` path had never shipped an end-to-end test, and we found out why when someone finally tried to write C3's mandatory fixture. Two invoked readers writing one owned key: the model would not load. `internal/table/load.go:371` — "owned tag `stage` is served by 2 readers; want exactly one". The second route was closed too: `internal/accessor/executor.go:89` classifies a reader's output against `def.RequestedKeys()`, so no reader can smuggle a key it did not declare. The third, a repeated `--tag`, the RDR itself documents as refused. **`tv.conflicted` was unreachable from the CLI the whole time.**

That was the day the postmortem stopped being about bugs. `0011:A3` — Status **Verified**, Method **Source Search** — asserts "the conflicted exception is real and user-reachable", and it is the first of the "three facts force that placement" in the Approach, the stated ground in `D-placement`, and the closing sentence of the Alternative 3 rejection ("the failure is real, not theoretical"). The Source Search never opened `internal/table/load.go`. And `0011:A2`, also Verified, had already recorded that the CLI *can* see the absent case: "a match key the view lacks lands in the candidate's undecided list today."

So the case the CLI could not handle could not occur, and the case that occurs the CLI already handled. We had moved the verdict into the kernel, taken the `large` Profile, changed `flow resolve`, vetoed a live model for six hours, broken 0010's default row, and deleted an escape-rescue control — to distinguish a condition no user can produce. Alternative 3, rejected across three sections, would have delivered the narrowing the original defect asked for with none of it.

**Day 14 — the census.** Five files with `Match` blocks that `A7`'s six-file census never listed — `adversarial_test.go`, `fixup_test.go`, `reserved_key_sameview_0008_test.go`, `escape_shape_adv_0009_test.go`, `resolve_test.go` — had gone red on the branch. We had re-expected them to get green, which is what `S5` says is a defect ("any OTHER oracle going red is a defect, not an expected inversion"). One of them, `reserved_key_sameview_0008_test.go:76-80`, carried a comment stating the veto in plain language: "any other value reaches the seam's default and answers GuardUnevaluable, which would refuse the resolve rather than prune the row." The failure mode was documented in the test suite, in the package the RDR changed, and the census that would have surfaced it was `Pending` at lock while three other sections cited its conclusion as fact.

**What it cost.** Six hours of blocked RDR flow, two verbs regressed (`flow resolve` and 0010's `otherwise` path), one deleted regression control, one silently-dead oracle, and RDR 0012 to override `0011:C1` — the clause that had been locked for six weeks.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time tests: each is executable against the *record* plus the shipped source, before any code is written.

### AT-1 — the multi-row veto (catches C-1, C-16)
```gherkin
Scenario: an indeterminate match row must not veto a sibling's plan
  Given a table with two ordinary rows for outcome "advance"
    And row X has match {stage: "resolved"} and the view carries stage=resolved
    And row Y has match {stage: "reconciled"} over a key ABSENT from the view
    And row X's guard is decided TRUE
  When flow resolve --outcome advance is run
  Then the result MUST be a Plan naming row X
   And the result MUST NOT be a refusal of any kind
```
**Review-time form.** Grep the record for `0007:C10`. Zero hits ⇒ BLOCK. Then read `resolve.go:540-565` and answer in writing: with `undecidable` non-empty and `selected` non-empty, what does `gate` return? The answer is `nil, &Refusal{...}` — the plan is discarded. The RDR must either cite C10 and accept the veto explicitly in Consequences (with the `models/rdr.toml` 7–8-row blast radius stated), or specify the row-scoping change and its own override of `0007:C10`. It does neither.

### AT-2 — the missing partition (catches C-2)
```gherkin
Scenario: "non-selectable survivor" must be a nameable kernel state
  Given RDR 0011 C1 requires a row that survives into gate, is not selected,
        and does not count toward len(selected)
  When gate's partition is enumerated from resolve.go:521-556
  Then one of {GuardFalse, GuardTrue, GuardUnevaluable} MUST denote that state
    Or the RDR MUST name the new Row field / gate signature that does
```
**Review-time form.** Enumerate `gate`'s three arms; show none is the state C1 describes. Then take `match_conflicted_test.go::gateMatchRow` — guard-free, therefore `GuardTrue` under `0007:C5` — and ask which arm it lands in under C1. The record has no answer. BLOCK until Phase 1 names the mechanism.

### AT-3 — 0010's otherwise row (catches C-3)
```gherkin
Scenario: an unsupplied tag still reaches the decision table's default
  Given a 0010 decision table with an "otherwise" escape row rescuing no_match
    And an ordinary rule matching on tag "size", which is NOT supplied
  When flow resolve is run
  Then the otherwise row MUST rescue
   And the plan MUST report escaped = true and carry its emit block
```
**Review-time form.** Trace: does a `guard_unevaluable` from `gate` ever reach `escapeOrRefuse`? Read `resolve.go:479-483` — the block returns above the call. Cross-read `0010:A9` ("escape row rescuing `no_match`") against `0011:C1` ("`indeterminate` … refuses `guard_unevaluable`, the non-escapable class"). This falsifies `0011:JC1`'s "cannot turn a 0010 acceptance into a refusal" on the record alone, with no code run. BLOCK: JC1 must be re-dispositioned, and 0010 owes the class-specific clause `0010:A11`'s "If wrong" already anticipates.

### AT-4 — the census and the arithmetic (catches C-4, C-5, C-12)
```gherkin
Scenario: the changed-oracle set is enumerable and consistently counted
  Given RDR 0011 claims a bounded set of internal/resolve oracles change
  When every file in internal/resolve building a row with a Match block is listed
  Then A7's census MUST cover all of them
   And S5, MVV step 8, and Consequences MUST state the SAME count
```
**Review-time form.** Run `grep -rln 'Match:' internal/resolve/*_test.go` → nine files. A7 names six. Then grep the record for the count: S5 says "Five", MVV 8 says "the four oracles S5 names", Consequences says "Three shipped `internal/resolve` oracles". Three different numbers for one set, in one record, one of which is the reviewer-facing blast-radius figure. BLOCK on both arms: A7 is `Pending` while three sections state its conclusion as settled, which `§Assumption Verification` forbids ("no assumption marked `Pending` … may have settled-fact prose elsewhere in the RDR depending on it").

### AT-5 — do not re-expect a control (catches C-6)
```gherkin
Scenario: a positive control is not re-expected into the thing it controls for
  Given S5 re-expects TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues
  When that test's docstring is read
  Then it MUST NOT declare itself the positive control for another test
    Or S5 MUST name the replacement oracle preserving the controlled property
```
**Review-time form.** Read `match_conflicted_test.go:190-194`: "is the positive control for the test above. Pruning conflicted-key matches must prune ONLY those." Under 0011 the property "a clean escape row still rescues in the presence of a conflicted key" becomes unobservable — `guard_unevaluable` precedes the escape phase. BLOCK: S5 must add a replacement oracle for over-pruning, or the RDR must accept in writing that escape rescue loses its regression control.

### AT-6 — isGuardBlock is not a payload filter (catches C-7)
```gherkin
Scenario: widening isGuardBlock changes only the payload
  Given A8 claims widening isGuardBlock admits BlockMatch to Refusal.Undecided
        while leaving every existing guard payload byte-identical
  When evaluateAtoms is read at guard.go:120-155
  Then isGuardBlock MUST gate the payload append ONLY
   And MUST NOT gate the seam call or the kleeneAnd operand contribution
```
**Review-time form.** Read the function. `isGuardBlock` gates the whole per-atom body: `continue` skips the seam call, the payload append, *and* the `switch atom.Block` that feeds `allResult`/`unlessConj`. The comment above it forbids exactly A8's edit: sweeping such an atom into `all_result` "would let it be DECIDED, and a decided-FALSE one would prune its row along with the row's owned-state obligation (D8) — reopening the masking path `0007:C11` ratifies pruning against." A8 cites C11 in support of the edit C11's own implementation comment forbids. BLOCK: A8 is `Pending` and its stated mechanism is refuted by the source it cites.

### AT-7 — the rename is not mechanical (catches C-8)
```gherkin
Scenario: no oracle passes vacuously after the unresolved → unknown rename
  Given C3 mandates six mechanical re-homings of the "unresolved" key
  When each of the six call sites is read
  Then no site may discard stringsAt's ok bool
   And no negative assertion may be satisfied by an empty slice
```
**Review-time form.** Read `flow_harness_0005_test.go:426-443` — `e.(string)` returns `nil, false` on an object element. Then read the six sites: `flow_adversarial_0005_test.go:446` and `flow_mvv_0005_test.go:66` are both `unresolved, _ :=`. The adversarial one is a negative assertion, so it passes vacuously post-rename. BLOCK: C3 must mandate replacing `stringsAt` with an `objectsAt`-based reader at all six sites and must require each re-homed assertion to fail against the pre-rename payload — a negative control for the re-homing itself.

### AT-8 — the carrier requires signature changes (catches C-9)
```gherkin
Scenario: the kernel's match facts can reach summarize
  Given C1 makes Refusal.Undecided the CARRIER for match facts
  When excluded() and summarize() signatures are read
  Then some parameter MUST convey the probe's Refusal to summarize
```
**Review-time form.** `excluded(row, owned, observed) bool` (`flow_next.go:276`) discards `result`; `summarize(row, view, gatesRan)` (`flow_next.go:168`) has no kernel input; `runFlowNext` (`flow_next.go:114-129`) calls them independently. C1's "no new field on `Plan` is required" is true and irrelevant — the missing plumbing is two function signatures and the loop between them, which neither Phase 2 nor the Infrastructure Audit names. BLOCK: Phase 2 must state the signatures, or the implementer will keep them and re-derive facts CLI-side, which is Alternative 3.

### AT-9 — reason tokens must carry distinct remedies (catches C-10)
```gherkin
Scenario: a conflicted key never reports a remedy the user cannot act on
  Given D-undecided-vocabulary states absent and uncomparable have different remedies
  When a conflicted match key is reported as {key, absent}
  Then the user following that remedy (supply the tag) MUST change the output
```
**Review-time form.** Follow the remedy on C1's compound case: the key is conflicted because two readers write it; supplying a tag does not deconflict an owned-provenance collision (`merge` conflicts within one provenance, `resolve.go:238-240`). The remedy is inert. Combined with AT-8, the collapse is not confined to the compound case. BLOCK: either the compound case reports a third token, or C1 states in Consequences that a user-visible remedy can be wrong and why that is acceptable.

### AT-10 — the mandatory fixture must be constructible (catches C-11, C-17)
```gherkin
Scenario: two invoked readers can write the same owned key with differing values
  Given C3 makes this fixture MANDATORY and bans the repeated --tag form
  When a model declaring two readers over one owned key is loaded
  Then the model MUST load
   And flow next over it MUST report {key, uncomparable}
```
**Review-time form — this is the test that should have stopped the lock.** A3 argues reachability from the *absence* of dedup in `runReaders`. That is necessary, not sufficient. Write the model. It does not load: `internal/table/load.go:371` fails `CatMalformedAccessorBinding` — "owned tag %s is served by %d readers; want exactly one". Close the second route: `internal/accessor/executor.go:89` classifies against `def.RequestedKeys()`, so a reader cannot emit an undeclared key. Close the third: `parseTags` refuses a repeated `--tag`, as the RDR states. Conclusion: `tv.conflicted` is unreachable from the CLI.

Then re-run the Alternative 3 rejection with that fact in hand. `§approach`(i), `D-placement`, and `ALT3`'s "the failure is real, not theoretical" all collapse — and `0011:A2` already establishes, Verified, that the CLI *can* see the absent case ("a match key the view lacks lands in the candidate's undecided list today"). BLOCK, and re-open Stage 2: the kernel-seam placement, the `large` Profile, and the entire `flow resolve` blast radius were bought to distinguish a case no user can produce. Every other load-bearing claim in this RDR got a spike; the one that decides the placement got an inference over an unopened file.

### AT-14 — the carrier must be non-empty for the motivating row (catches C-18)
```gherkin
Scenario: match facts reach the payload on a guard-free row
  Given a row with an EMPTY guard block and one match atom over an absent key
  When flow next probes that row
  Then Refusal.Undecided MUST contain an UndecidedRow for it
   And that row's Atoms MUST name the match key with reason "absent"
```
**Review-time form.** Read the type: `Undecided []UndecidedRow` (`resolve.go:375`) — a per-ROW wrapper, not a flat atom list, so C1's repeated "atoms ride in `Refusal.Undecided`" is one indirection off. Then read `guard.go:162-166`: `if verdict != GuardUnevaluable { return verdict, nil }` — a decided row emits no payload at all. Now apply `0007:C5`: "A row with no atoms is decided TRUE." A guard-free row is `GuardTrue`, therefore emits no `UndecidedRow`, therefore carries no match atoms — and `match_conflicted_test.go::gateMatchRow` is exactly that shape, as is the ordinary row in every decision table 0010 describes. BLOCK: C1 must specify how a `GuardTrue` row acquires an `UndecidedRow`, which is the same missing mechanism as AT-2.

### AT-11 — the negative flag surface must be pinned (catches C-13)
```gherkin
Scenario: --all is refused by every selection verb
  Given C2 satisfies this by NON-REGISTRATION
  When flow resolve --all is invoked
  Then the command MUST fail
   And a test MUST assert the failure for resolve, read-state, and set-state
```
**Review-time form.** Non-registration is the absence of a contract. `registerSelectionFlags` is the shared registrar any future verb-add reaches for; nothing prevents `--all` migrating there. Given C1's rule that a selection site must not act on an undecidable precondition, `flow resolve --all` would emit a write-performing plan on an unchecked match — the exact masking C1 forbids. S2 asserts it, but on cobra's third-party error class rather than on a `flow-*` code. BLOCK: C2 must state the positive obligation ("`--all` MUST NOT be registered on `registerSelectionFlags`") so a future edit trips a named contract.

### AT-12 — validate on the shape that will exist (catches C-14)
```gherkin
Scenario: the zero-owned decision table is exercised before lock
  Given A5 records that a genuine zero-owned decision table is unloadable
    And the fixture therefore carries a dummy owned tag
  When 0011's absent-key provision is validated
  Then it MUST be validated on the shape 0010 will ship, not the substitute
```
**Review-time form.** A5 admits the substitution and defers the real case to "downstream of 0010 shipping." JC1 accepts that deferral. The result is that 0011's coupling to 0010 is verified against a shape that cannot exist. BLOCK: either 0010 ships first and 0011 re-verifies A5, or 0011 states in Consequences that the decision-table class is unvalidated on this build.

### AT-13 — the narrowing must be shown where it can fail (catches C-15)
```gherkin
Scenario: narrowing is demonstrated on a model whose match keys are NOT reader-served
  Given A6 verifies 21 → 3 on models/rdr.toml
    And A6 records that on that model "C1's absent-key provision never fires"
  When the MVV is run on a model whose match keys are observed and unsupplied
  Then the candidate count MUST narrow
    Or the RDR MUST state that it does not, in the Problem Statement's success criterion
```
**Review-time form.** A6's own text concedes the machinery is unreachable on the validating model. A5 confirms every row remains a candidate where it *is* reachable. The Problem Statement pre-absolves this ("intended behaviour, not a residual defect"), which makes the success criterion unfalsifiable by the user who filed the original defect. BLOCK: MVV must add a step on an observed-key model and state the expected cardinality — even if that cardinality is "unchanged" — so the limit is a pinned contract rather than a paragraph a future reader has to find.
